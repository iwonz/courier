// Package update implements verified GitHub release self-updates.
package update

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/iwonz/courier/internal/safety"
	"golang.org/x/mod/semver"
)

var maxDownload int64 = 512 << 20

type stagedFile interface {
	io.Writer
	Name() string
	Chmod(os.FileMode) error
	Sync() error
	Close() error
}

var (
	runtimeOS             = runtime.GOOS
	runtimeArch           = runtime.GOARCH
	makeUpdateDirectory   = os.MkdirTemp
	removeUpdateDirectory = os.RemoveAll
	currentExecutable     = os.Executable
	evaluateSymlinks      = filepath.EvalSymlinks
	createStagedFile      = func(directory, pattern string) (stagedFile, error) { return os.CreateTemp(directory, pattern) }
	openUpdateSource      = func(name string) (io.ReadCloser, error) { return os.Open(name) }
	createDownloadFile    = func(name string) (io.WriteCloser, error) {
		return os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	}
	newHTTPRequest      = http.NewRequestWithContext
	readChecksumFile    = os.ReadFile
	openChecksumFile    = func(name string) (io.ReadCloser, error) { return os.Open(name) }
	openReleaseArchive  = func(name string) (io.ReadCloser, error) { return os.Open(name) }
	createExtractedFile = func(name string) (io.WriteCloser, error) {
		return os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o700)
	}
)

// HTTPClient is the request boundary used by Updater.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Updater checks, verifies, and replaces one Courier executable.
type Updater struct {
	Repository string
	Version    string
	Token      string
	Client     HTTPClient
	GOOS       string
	GOARCH     string
	Executable func() (string, error)
	Rename     func(string, string) error
	Handoff    func(string, string, int) error
}

// Result describes an update check or installation.
type Result struct {
	Current bool
	From    string
	To      string
	Notes   string
}

// BinaryArtifact is a verified Courier executable extracted from one exact
// GitHub Release. Call Cleanup when the executable is no longer needed.
type BinaryArtifact struct {
	Path    string
	SHA256  string
	Cleanup func() error
}

type githubRelease struct {
	TagName string `json:"tag_name"`
	Body    string `json:"body"`
	Assets  []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

// Run installs a newer stable release when one exists.
func (u Updater) Run(ctx context.Context) (Result, error) {
	repository := u.Repository
	if repository == "" {
		repository = "iwonz/courier"
	}
	client := u.Client
	if client == nil {
		client = http.DefaultClient
	}
	var release githubRelease
	if err := u.getJSON(ctx, client, "https://api.github.com/repos/"+repository+"/releases/latest", &release); err != nil {
		return Result{}, err
	}
	latest := normalizeVersion(release.TagName)
	current := normalizeVersion(u.Version)
	if latest == "" {
		return Result{}, fmt.Errorf("GitHub returned invalid release tag %q", release.TagName)
	}
	result := Result{From: u.Version, To: release.TagName, Notes: release.Body}
	if current != "" && semver.Compare(current, latest) >= 0 {
		result.Current = true
		return result, nil
	}
	goos, goarch := u.GOOS, u.GOARCH
	if goos == "" {
		goos = runtimeOS
	}
	if goarch == "" {
		goarch = runtimeArch
	}
	versionText := strings.TrimPrefix(latest, "v")
	assetName := fmt.Sprintf("courier_%s_%s_%s.tar.gz", versionText, goos, goarch)
	assetURL, checksumURL := "", ""
	for _, asset := range release.Assets {
		switch asset.Name {
		case assetName:
			assetURL = asset.URL
		case "checksums.txt":
			checksumURL = asset.URL
		}
	}
	if assetURL == "" || checksumURL == "" {
		return Result{}, fmt.Errorf("release %s has no verified asset %s", release.TagName, assetName)
	}
	temporaryDirectory, err := makeUpdateDirectory("", "courier-update-*")
	if err != nil {
		return Result{}, err
	}
	defer removeUpdateDirectory(temporaryDirectory)
	archivePath := filepath.Join(temporaryDirectory, assetName)
	if err := u.download(ctx, client, assetURL, archivePath); err != nil {
		return Result{}, err
	}
	checksumPath := filepath.Join(temporaryDirectory, "checksums.txt")
	if err := u.download(ctx, client, checksumURL, checksumPath); err != nil {
		return Result{}, err
	}
	if err := verifyChecksum(archivePath, checksumPath, assetName); err != nil {
		return Result{}, err
	}
	binaryName := "courier"
	if goos == "windows" {
		binaryName += ".exe"
	}
	extracted := filepath.Join(temporaryDirectory, binaryName)
	if err := extractBinary(archivePath, extracted, binaryName); err != nil {
		return Result{}, err
	}
	executable := u.Executable
	if executable == nil {
		executable = currentExecutable
	}
	target, err := executable()
	if err != nil {
		return Result{}, err
	}
	target, err = evaluateSymlinks(target)
	if err != nil {
		return Result{}, err
	}
	partial, err := createStagedFile(filepath.Dir(target), ".courier-update-partial-*")
	if err != nil {
		return Result{}, err
	}
	if err := partial.Chmod(0o700); err != nil {
		_ = partial.Close()
		_ = os.Remove(partial.Name())
		return Result{}, err
	}
	partialPath := partial.Name()
	committed := false
	defer func() {
		_ = partial.Close()
		if !committed {
			_ = os.Remove(partialPath)
		}
	}()
	source, err := openUpdateSource(extracted)
	if err != nil {
		return Result{}, err
	}
	_, copyErr := io.Copy(partial, source)
	closeSourceErr := source.Close()
	if copyErr != nil || closeSourceErr != nil {
		return Result{}, errors.Join(copyErr, closeSourceErr)
	}
	if err := partial.Sync(); err != nil {
		return Result{}, err
	}
	if err := partial.Close(); err != nil {
		return Result{}, err
	}
	if goos == "windows" {
		handoff := u.Handoff
		if handoff == nil {
			handoff = launchUpdateHandoff
		}
		if err := handoff(partialPath, target, currentProcessID()); err != nil {
			return Result{}, fmt.Errorf("start Windows update handoff: %w", err)
		}
		committed = true
		return result, nil
	}
	rename := u.Rename
	if rename == nil {
		rename = replaceUpdateFile
	}
	if err := rename(partialPath, target); err != nil {
		return Result{}, fmt.Errorf("replace executable: %w", err)
	}
	committed = true
	return result, nil
}

// AcquireVersion downloads and verifies the Courier executable for one exact
// semantic version and platform. It does not replace the running executable.
func (u Updater) AcquireVersion(ctx context.Context, version, goos, goarch string) (*BinaryArtifact, error) {
	tag := normalizeVersion(version)
	if tag == "" {
		return nil, fmt.Errorf("cannot acquire helper for non-release version %q", version)
	}
	repository := u.Repository
	if repository == "" {
		repository = "iwonz/courier"
	}
	client := u.Client
	if client == nil {
		client = http.DefaultClient
	}
	versionText := strings.TrimPrefix(tag, "v")
	assetName := fmt.Sprintf("courier_%s_%s_%s.tar.gz", versionText, goos, goarch)
	baseURL := "https://github.com/" + repository + "/releases/download/" + tag + "/"
	temporaryDirectory, err := makeUpdateDirectory("", "courier-helper-download-*")
	if err != nil {
		return nil, err
	}
	cleanup := func() error { return removeUpdateDirectory(temporaryDirectory) }
	fail := func(cause error) (*BinaryArtifact, error) {
		return nil, errors.Join(cause, cleanup())
	}
	archivePath := filepath.Join(temporaryDirectory, assetName)
	if err := u.download(ctx, client, baseURL+assetName, archivePath); err != nil {
		return fail(err)
	}
	checksumPath := filepath.Join(temporaryDirectory, "checksums.txt")
	if err := u.download(ctx, client, baseURL+"checksums.txt", checksumPath); err != nil {
		return fail(err)
	}
	if err := verifyChecksum(archivePath, checksumPath, assetName); err != nil {
		return fail(err)
	}
	binaryName := "courier"
	if goos == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join(temporaryDirectory, binaryName)
	if err := extractBinary(archivePath, binaryPath, binaryName); err != nil {
		return fail(err)
	}
	digest, err := fileDigest(binaryPath)
	if err != nil {
		return fail(err)
	}
	return &BinaryArtifact{Path: binaryPath, SHA256: digest, Cleanup: cleanup}, nil
}

func (u Updater) getJSON(ctx context.Context, client HTTPClient, url string, destination any) error {
	response, err := u.request(ctx, client, url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return json.NewDecoder(io.LimitReader(response.Body, 2<<20)).Decode(destination)
}

func (u Updater) download(ctx context.Context, client HTTPClient, url, destination string) error {
	response, err := u.request(ctx, client, url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.ContentLength > maxDownload {
		return fmt.Errorf("download exceeds %d bytes", maxDownload)
	}
	file, err := createDownloadFile(destination)
	if err != nil {
		return err
	}
	written, copyErr := io.Copy(file, io.LimitReader(response.Body, maxDownload+1))
	closeErr := file.Close()
	if copyErr != nil || closeErr != nil {
		return errors.Join(copyErr, closeErr)
	}
	if written > maxDownload {
		return fmt.Errorf("download exceeds %d bytes", maxDownload)
	}
	return nil
}

func (u Updater) request(ctx context.Context, client HTTPClient, url string) (*http.Response, error) {
	request, err := newHTTPRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("User-Agent", "courier/"+u.Version)
	if u.Token != "" {
		request.Header.Set("Authorization", "Bearer "+u.Token)
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		defer response.Body.Close()
		message, _ := io.ReadAll(io.LimitReader(response.Body, 4<<10))
		return nil, fmt.Errorf("GitHub request failed (%s): %s", response.Status, strings.TrimSpace(string(message)))
	}
	return response, nil
}

func normalizeVersion(value string) string {
	if value == "" {
		return ""
	}
	if !strings.HasPrefix(value, "v") {
		value = "v" + value
	}
	if !semver.IsValid(value) {
		return ""
	}
	return value
}

func verifyChecksum(archivePath, checksumPath, assetName string) error {
	manifest, err := readChecksumFile(checksumPath)
	if err != nil {
		return err
	}
	expected := ""
	for _, line := range strings.Split(string(manifest), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && strings.TrimPrefix(fields[1], "*") == assetName {
			expected = strings.ToLower(fields[0])
			break
		}
	}
	if expected == "" {
		return fmt.Errorf("checksums.txt has no entry for %s", assetName)
	}
	file, err := openChecksumFile(archivePath)
	if err != nil {
		return err
	}
	hash := sha256.New()
	_, copyErr := io.Copy(hash, file)
	closeErr := file.Close()
	if copyErr != nil || closeErr != nil {
		return errors.Join(copyErr, closeErr)
	}
	actual := hex.EncodeToString(hash.Sum(nil))
	if actual != expected {
		return fmt.Errorf("checksum mismatch for %s", assetName)
	}
	return nil
}

func fileDigest(name string) (string, error) {
	file, err := openChecksumFile(name)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	_, copyErr := io.Copy(hash, file)
	closeErr := file.Close()
	if copyErr != nil || closeErr != nil {
		return "", errors.Join(copyErr, closeErr)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func extractBinary(archivePath, destination, binaryName string) error {
	file, err := openReleaseArchive(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()
	reader := tar.NewReader(gzipReader)
	found := false
	for {
		header, err := reader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		if _, err := safety.SafeArchiveJoin("release-root", header.Name); err != nil {
			return err
		}
		if filepath.Base(filepath.FromSlash(header.Name)) != binaryName {
			continue
		}
		if found || header.Typeflag != tar.TypeReg {
			return fmt.Errorf("release archive has invalid %s entry", binaryName)
		}
		output, err := createExtractedFile(destination)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(output, reader)
		closeErr := output.Close()
		if copyErr != nil || closeErr != nil {
			return errors.Join(copyErr, closeErr)
		}
		found = true
	}
	if !found {
		return fmt.Errorf("release archive has no %s", binaryName)
	}
	return nil
}
