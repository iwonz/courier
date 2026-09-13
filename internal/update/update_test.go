package update

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeHTTPClient struct {
	responses map[string]fakeResponse
	requests  []*http.Request
}

type fakeResponse struct {
	status int
	body   []byte
	reader io.ReadCloser
	err    error
	length int64
}

func (f *fakeHTTPClient) Do(request *http.Request) (*http.Response, error) {
	f.requests = append(f.requests, request)
	response := f.responses[request.URL.String()]
	if response.err != nil {
		return nil, response.err
	}
	status := response.status
	if status == 0 {
		status = http.StatusOK
	}
	length := response.length
	if length == 0 {
		length = int64(len(response.body))
	}
	body := response.reader
	if body == nil {
		body = io.NopCloser(bytes.NewReader(response.body))
	}
	return &http.Response{StatusCode: status, Status: http.StatusText(status), ContentLength: length, Body: body}, nil
}

func TestUpdaterSuccessAndCurrent(t *testing.T) {
	assetName := "courier_1.2.0_linux_amd64.tar.gz"
	archive := releaseArchive(t, []tar.Header{{Name: "README.md", Mode: 0o600, Size: 4, Typeflag: tar.TypeReg}, {Name: "courier", Mode: 0o700, Size: 3, Typeflag: tar.TypeReg}}, [][]byte{[]byte("docs"), []byte("new")})
	hash := sha256.Sum256(archive)
	manifest := []byte(hex.EncodeToString(hash[:]) + "  " + assetName + "\n")
	releaseJSON := []byte(`{"tag_name":"v1.2.0","body":"release notes","assets":[{"name":"` + assetName + `","browser_download_url":"https://download/asset"},{"name":"checksums.txt","browser_download_url":"https://download/checksums"}]}`)
	client := &fakeHTTPClient{responses: map[string]fakeResponse{
		"https://api.github.com/repos/iwonz/courier/releases/latest": {body: releaseJSON},
		"https://download/asset":                                     {body: archive},
		"https://download/checksums":                                 {body: manifest},
	}}
	root := t.TempDir()
	executable := filepath.Join(root, "courier")
	if err := os.WriteFile(executable, []byte("old"), 0o700); err != nil {
		t.Fatal(err)
	}
	updater := Updater{Version: "v1.0.0", Token: "token", Client: client, GOOS: "linux", GOARCH: "amd64", Executable: func() (string, error) { return executable, nil }}
	result, err := updater.Run(context.Background())
	if err != nil || result.Current || result.To != "v1.2.0" || result.Notes != "release notes" {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	data, err := os.ReadFile(executable)
	if err != nil || string(data) != "new" {
		t.Fatalf("executable=%q err=%v", data, err)
	}
	if authorization := client.requests[0].Header.Get("Authorization"); authorization != "Bearer token" {
		t.Fatalf("authorization=%q", authorization)
	}
	updater.Version = "1.2.0"
	result, err = updater.Run(context.Background())
	if err != nil || !result.Current {
		t.Fatalf("current result=%+v err=%v", result, err)
	}
}

func TestUpdaterMetadataAndDownloadFailures(t *testing.T) {
	api := "https://api.github.com/repos/iwonz/courier/releases/latest"
	for _, test := range []struct {
		name      string
		responses map[string]fakeResponse
	}{
		{"request", map[string]fakeResponse{api: {err: errors.New("network")}}},
		{"status", map[string]fakeResponse{api: {status: 500, body: []byte("failure")}}},
		{"json", map[string]fakeResponse{api: {body: []byte("not json")}}},
		{"tag", map[string]fakeResponse{api: {body: []byte(`{"tag_name":"bad"}`)}}},
		{"assets", map[string]fakeResponse{api: {body: []byte(`{"tag_name":"v1.0.0"}`)}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			updater := Updater{Version: "dev", Client: &fakeHTTPClient{responses: test.responses}, GOOS: "linux", GOARCH: "amd64"}
			if _, err := updater.Run(context.Background()); err == nil {
				t.Fatal("expected update error")
			}
		})
	}
	assetName := "courier_1.0.0_linux_amd64.tar.gz"
	releaseJSON := []byte(`{"tag_name":"v1.0.0","assets":[{"name":"` + assetName + `","browser_download_url":"asset"},{"name":"checksums.txt","browser_download_url":"checksums"}]}`)
	for _, responses := range []map[string]fakeResponse{
		{api: {body: releaseJSON}, "asset": {err: errors.New("asset")}},
		{api: {body: releaseJSON}, "asset": {body: []byte("x")}, "checksums": {err: errors.New("checksums")}},
		{api: {body: releaseJSON}, "asset": {body: []byte("x")}, "checksums": {body: []byte("bad checksum")}},
		{api: {body: releaseJSON}, "asset": {body: []byte("x"), length: maxDownload + 1}},
	} {
		updater := Updater{Version: "dev", Client: &fakeHTTPClient{responses: responses}, GOOS: "linux", GOARCH: "amd64"}
		if _, err := updater.Run(context.Background()); err == nil {
			t.Fatal("expected asset failure")
		}
	}
}

func TestUpdaterReplacementFailures(t *testing.T) {
	assetName := "courier_1.0.0_linux_amd64.tar.gz"
	archive := releaseArchive(t, []tar.Header{{Name: "courier", Mode: 0o700, Size: 3, Typeflag: tar.TypeReg}}, [][]byte{[]byte("new")})
	hash := sha256.Sum256(archive)
	api := "https://api.github.com/repos/custom/repo/releases/latest"
	releaseJSON := []byte(`{"tag_name":"1.0.0","assets":[{"name":"` + assetName + `","browser_download_url":"asset"},{"name":"checksums.txt","browser_download_url":"checksums"}]}`)
	responses := map[string]fakeResponse{api: {body: releaseJSON}, "asset": {body: archive}, "checksums": {body: []byte(hex.EncodeToString(hash[:]) + " *" + assetName)}}
	for _, updater := range []Updater{
		{Repository: "custom/repo", Version: "dev", Client: &fakeHTTPClient{responses: responses}, GOOS: "linux", GOARCH: "amd64", Executable: func() (string, error) { return "", errors.New("executable") }},
		{Repository: "custom/repo", Version: "dev", Client: &fakeHTTPClient{responses: responses}, GOOS: "linux", GOARCH: "amd64", Executable: func() (string, error) { return filepath.Join(t.TempDir(), "missing"), nil }},
	} {
		if _, err := updater.Run(context.Background()); err == nil {
			t.Fatal("expected executable failure")
		}
	}
	root := t.TempDir()
	executable := filepath.Join(root, "courier")
	if err := os.WriteFile(executable, []byte("old"), 0o700); err != nil {
		t.Fatal(err)
	}
	updater := Updater{Repository: "custom/repo", Version: "dev", Client: &fakeHTTPClient{responses: responses}, GOOS: "linux", GOARCH: "amd64", Executable: func() (string, error) { return executable, nil }, Rename: func(string, string) error { return errors.New("rename") }}
	if _, err := updater.Run(context.Background()); err == nil {
		t.Fatal("expected rename error")
	}
	data, _ := os.ReadFile(executable)
	if string(data) != "old" {
		t.Fatalf("original changed: %q", data)
	}
}

func TestChecksumAndExtraction(t *testing.T) {
	root := t.TempDir()
	asset := filepath.Join(root, "asset")
	manifest := filepath.Join(root, "checksums")
	if err := os.WriteFile(asset, []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := verifyChecksum(asset, filepath.Join(root, "missing"), "asset"); err == nil {
		t.Fatal("expected manifest read error")
	}
	if err := os.WriteFile(manifest, []byte("none"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := verifyChecksum(asset, manifest, "asset"); err == nil {
		t.Fatal("expected missing checksum")
	}
	if err := os.WriteFile(manifest, []byte(strings.Repeat("0", 64)+"  asset"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := verifyChecksum(asset, manifest, "asset"); err == nil {
		t.Fatal("expected checksum mismatch")
	}
	if err := extractBinary(filepath.Join(root, "missing"), filepath.Join(root, "out"), "courier"); err == nil {
		t.Fatal("expected archive open error")
	}
	badGzip := filepath.Join(root, "bad.gz")
	_ = os.WriteFile(badGzip, []byte("bad"), 0o600)
	if err := extractBinary(badGzip, filepath.Join(root, "out"), "courier"); err == nil {
		t.Fatal("expected gzip error")
	}
	for _, test := range []struct {
		name    string
		headers []tar.Header
		data    [][]byte
	}{
		{"missing", []tar.Header{{Name: "readme", Typeflag: tar.TypeReg}}, [][]byte{nil}},
		{"unsafe", []tar.Header{{Name: "../courier", Typeflag: tar.TypeReg}}, [][]byte{nil}},
		{"directory", []tar.Header{{Name: "courier", Typeflag: tar.TypeDir}}, [][]byte{nil}},
		{"duplicate", []tar.Header{{Name: "courier", Typeflag: tar.TypeReg}, {Name: "dir/courier", Typeflag: tar.TypeReg}}, [][]byte{nil, nil}},
	} {
		name := filepath.Join(root, test.name+".tar.gz")
		_ = os.WriteFile(name, releaseArchive(t, test.headers, test.data), 0o600)
		if err := extractBinary(name, filepath.Join(root, test.name+"-out"), "courier"); err == nil {
			t.Errorf("expected %s extraction error", test.name)
		}
	}
}

func TestNormalizeVersion(t *testing.T) {
	for input, want := range map[string]string{"": "", "bad": "", "1.2.3": "v1.2.3", "v2.0.0": "v2.0.0"} {
		if got := normalizeVersion(input); got != want {
			t.Errorf("normalizeVersion(%q)=%q", input, got)
		}
	}
}

func TestUpdaterDefaultRuntimeWindowsAndRunFailures(t *testing.T) {
	t.Run("default client runtime and executable", func(t *testing.T) {
		isolateUpdateHooks(t)
		assetName := "courier_1.0.0_" + runtimeOS + "_" + runtimeArch + ".tar.gz"
		binaryName := "courier"
		if runtimeOS == "windows" {
			binaryName += ".exe"
		}
		archive := releaseArchive(t, []tar.Header{{Name: binaryName, Mode: 0o700, Size: 3, Typeflag: tar.TypeReg}}, [][]byte{[]byte("new")})
		hash := sha256.Sum256(archive)
		releaseJSON := []byte(`{"tag_name":"v1.0.0","assets":[{"name":"` + assetName + `","browser_download_url":"https://download/asset"},{"name":"checksums.txt","browser_download_url":"https://download/checksums"}]}`)
		client := &fakeHTTPClient{responses: map[string]fakeResponse{
			"https://api.github.com/repos/iwonz/courier/releases/latest": {body: releaseJSON},
			"https://download/asset":                                     {body: archive},
			"https://download/checksums":                                 {body: []byte(hex.EncodeToString(hash[:]) + "  " + assetName)},
		}}
		originalDefault := http.DefaultClient
		http.DefaultClient = &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) { return client.Do(request) })}
		t.Cleanup(func() { http.DefaultClient = originalDefault })
		target := filepath.Join(t.TempDir(), binaryName)
		if err := os.WriteFile(target, []byte("old"), 0o700); err != nil {
			t.Fatal(err)
		}
		currentExecutable = func() (string, error) { return target, nil }
		if _, err := (Updater{Version: "dev"}).Run(context.Background()); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("windows asset", func(t *testing.T) {
		isolateUpdateHooks(t)
		updater, _, target := validUpdater(t, "windows", "amd64", "courier.exe")
		if _, err := updater.Run(context.Background()); err != nil {
			t.Fatal(err)
		}
		if data, err := os.ReadFile(target); err != nil || string(data) != "new" {
			t.Fatalf("target=%q err=%v", data, err)
		}
	})

	for _, test := range []struct {
		name   string
		mutate func(*testing.T, *Updater)
	}{
		{name: "temporary directory", mutate: func(_ *testing.T, _ *Updater) {
			makeUpdateDirectory = func(string, string) (string, error) { return "", errors.New("mkdir") }
		}},
		{name: "invalid archive", mutate: func(t *testing.T, updater *Updater) {
			assetName := "courier_1.0.0_linux_amd64.tar.gz"
			data := []byte("not gzip")
			hash := sha256.Sum256(data)
			updater.Client = releaseClient(assetName, data, []byte(hex.EncodeToString(hash[:])+"  "+assetName))
		}},
		{name: "partial create", mutate: func(_ *testing.T, _ *Updater) {
			createStagedFile = func(string, string) (stagedFile, error) { return nil, errors.New("partial") }
		}},
		{name: "partial chmod", mutate: func(_ *testing.T, _ *Updater) {
			createStagedFile = func(string, string) (stagedFile, error) {
				return &fakeStagedFile{name: "partial", chmodErr: errors.New("chmod")}, nil
			}
		}},
		{name: "source open", mutate: func(_ *testing.T, _ *Updater) {
			openUpdateSource = func(string) (io.ReadCloser, error) { return nil, errors.New("source") }
		}},
		{name: "source copy", mutate: func(_ *testing.T, _ *Updater) {
			openUpdateSource = func(string) (io.ReadCloser, error) { return &failingReadCloser{readErr: errors.New("read")}, nil }
		}},
		{name: "source close", mutate: func(_ *testing.T, _ *Updater) {
			openUpdateSource = func(string) (io.ReadCloser, error) {
				return &failingReadCloser{Reader: strings.NewReader("new"), closeErr: errors.New("close source")}, nil
			}
		}},
		{name: "partial sync", mutate: func(_ *testing.T, _ *Updater) {
			createStagedFile = func(string, string) (stagedFile, error) {
				return &fakeStagedFile{name: "partial", syncErr: errors.New("sync")}, nil
			}
		}},
		{name: "partial close", mutate: func(_ *testing.T, _ *Updater) {
			createStagedFile = func(string, string) (stagedFile, error) {
				return &fakeStagedFile{name: "partial", closeErr: errors.New("close partial")}, nil
			}
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			isolateUpdateHooks(t)
			updater, _, _ := validUpdater(t, "linux", "amd64", "courier")
			test.mutate(t, &updater)
			if _, err := updater.Run(context.Background()); err == nil {
				t.Fatal("expected updater failure")
			}
		})
	}
}

func TestDownloadRequestChecksumAndExtractFailures(t *testing.T) {
	t.Run("download file and body failures", func(t *testing.T) {
		isolateUpdateHooks(t)
		updater := Updater{}
		client := &fakeHTTPClient{responses: map[string]fakeResponse{"download": {body: []byte("data")}}}
		createDownloadFile = func(string) (io.WriteCloser, error) { return nil, errors.New("create") }
		if err := updater.download(context.Background(), client, "download", "file"); err == nil {
			t.Fatal("expected create failure")
		}
		createDownloadFile = func(string) (io.WriteCloser, error) { return &failingWriteCloser{writeErr: errors.New("write")}, nil }
		if err := updater.download(context.Background(), client, "download", "file"); err == nil {
			t.Fatal("expected write failure")
		}
		createDownloadFile = func(string) (io.WriteCloser, error) { return &failingWriteCloser{closeErr: errors.New("close")}, nil }
		if err := updater.download(context.Background(), client, "download", "file"); err == nil {
			t.Fatal("expected close failure")
		}
		createDownloadFile = func(string) (io.WriteCloser, error) { return &failingWriteCloser{}, nil }
		maxDownload = 3
		client.responses["download"] = fakeResponse{body: []byte("data"), length: -1}
		if err := updater.download(context.Background(), client, "download", "file"); err == nil {
			t.Fatal("expected streamed size failure")
		}
	})

	t.Run("malformed request", func(t *testing.T) {
		if _, err := (Updater{}).request(context.Background(), &fakeHTTPClient{}, ":"); err == nil {
			t.Fatal("expected request error")
		}
	})

	t.Run("checksum open read and close", func(t *testing.T) {
		isolateUpdateHooks(t)
		root := t.TempDir()
		manifest := filepath.Join(root, "checksums")
		if err := os.WriteFile(manifest, []byte(strings.Repeat("0", 64)+"  asset"), 0o600); err != nil {
			t.Fatal(err)
		}
		openChecksumFile = func(string) (io.ReadCloser, error) { return nil, errors.New("open") }
		if err := verifyChecksum("asset", manifest, "asset"); err == nil {
			t.Fatal("expected checksum open failure")
		}
		openChecksumFile = func(string) (io.ReadCloser, error) { return &failingReadCloser{readErr: errors.New("read")}, nil }
		if err := verifyChecksum("asset", manifest, "asset"); err == nil {
			t.Fatal("expected checksum read failure")
		}
		openChecksumFile = func(string) (io.ReadCloser, error) {
			return &failingReadCloser{Reader: strings.NewReader("data"), closeErr: errors.New("close")}, nil
		}
		if err := verifyChecksum("asset", manifest, "asset"); err == nil {
			t.Fatal("expected checksum close failure")
		}
	})

	t.Run("archive stream and output failures", func(t *testing.T) {
		isolateUpdateHooks(t)
		root := t.TempDir()
		archivePath := filepath.Join(root, "archive.tar.gz")
		if err := os.WriteFile(archivePath, releaseArchive(t, []tar.Header{{Name: "courier", Size: 3, Mode: 0o700, Typeflag: tar.TypeReg}}, [][]byte{[]byte("new")}), 0o600); err != nil {
			t.Fatal(err)
		}
		createExtractedFile = func(string) (io.WriteCloser, error) { return nil, errors.New("create output") }
		if err := extractBinary(archivePath, filepath.Join(root, "out"), "courier"); err == nil {
			t.Fatal("expected output create failure")
		}
		createExtractedFile = func(string) (io.WriteCloser, error) { return &failingWriteCloser{writeErr: errors.New("write")}, nil }
		if err := extractBinary(archivePath, filepath.Join(root, "out"), "courier"); err == nil {
			t.Fatal("expected output write failure")
		}
		createExtractedFile = func(string) (io.WriteCloser, error) { return &failingWriteCloser{closeErr: errors.New("close")}, nil }
		if err := extractBinary(archivePath, filepath.Join(root, "out"), "courier"); err == nil {
			t.Fatal("expected output close failure")
		}
		truncated := filepath.Join(root, "truncated.tar.gz")
		var truncatedData bytes.Buffer
		gzipWriter := gzip.NewWriter(&truncatedData)
		tarWriter := tar.NewWriter(gzipWriter)
		if err := tarWriter.WriteHeader(&tar.Header{Name: "readme", Size: 100, Mode: 0o600, Typeflag: tar.TypeReg}); err != nil {
			t.Fatal(err)
		}
		if _, err := tarWriter.Write([]byte("x")); err != nil {
			t.Fatal(err)
		}
		if err := gzipWriter.Close(); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(truncated, truncatedData.Bytes(), 0o600); err != nil {
			t.Fatal(err)
		}
		createExtractedFile = func(name string) (io.WriteCloser, error) {
			return os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o700)
		}
		if err := extractBinary(truncated, filepath.Join(root, "truncated-out"), "courier"); err == nil {
			t.Fatal("expected truncated archive failure")
		}
	})
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) { return f(request) }

type fakeStagedFile struct {
	bytes.Buffer
	name     string
	chmodErr error
	syncErr  error
	closeErr error
	writeErr error
}

func (f *fakeStagedFile) Name() string            { return f.name }
func (f *fakeStagedFile) Chmod(os.FileMode) error { return f.chmodErr }
func (f *fakeStagedFile) Sync() error             { return f.syncErr }
func (f *fakeStagedFile) Close() error            { return f.closeErr }
func (f *fakeStagedFile) Write(data []byte) (int, error) {
	if f.writeErr != nil {
		return 0, f.writeErr
	}
	return f.Buffer.Write(data)
}

type failingReadCloser struct {
	io.Reader
	readErr  error
	closeErr error
}

func (f *failingReadCloser) Read(data []byte) (int, error) {
	if f.readErr != nil {
		return 0, f.readErr
	}
	return f.Reader.Read(data)
}
func (f *failingReadCloser) Close() error { return f.closeErr }

type failingWriteCloser struct {
	data     []byte
	writeErr error
	closeErr error
}

func (f *failingWriteCloser) Write(data []byte) (int, error) {
	if f.writeErr != nil {
		return 0, f.writeErr
	}
	f.data = append(f.data, data...)
	return len(data), nil
}
func (f *failingWriteCloser) Close() error { return f.closeErr }

func validUpdater(t *testing.T, goos, goarch, binaryName string) (Updater, *fakeHTTPClient, string) {
	t.Helper()
	assetName := "courier_1.0.0_" + goos + "_" + goarch + ".tar.gz"
	archive := releaseArchive(t, []tar.Header{{Name: binaryName, Mode: 0o700, Size: 3, Typeflag: tar.TypeReg}}, [][]byte{[]byte("new")})
	hash := sha256.Sum256(archive)
	client := releaseClient(assetName, archive, []byte(hex.EncodeToString(hash[:])+"  "+assetName))
	target := filepath.Join(t.TempDir(), binaryName)
	if err := os.WriteFile(target, []byte("old"), 0o700); err != nil {
		t.Fatal(err)
	}
	updater := Updater{Version: "dev", Client: client, GOOS: goos, GOARCH: goarch, Executable: func() (string, error) { return target, nil }}
	return updater, client, target
}

func releaseClient(assetName string, archive, manifest []byte) *fakeHTTPClient {
	return &fakeHTTPClient{responses: map[string]fakeResponse{
		"https://api.github.com/repos/iwonz/courier/releases/latest": {body: []byte(`{"tag_name":"v1.0.0","assets":[{"name":"` + assetName + `","browser_download_url":"asset"},{"name":"checksums.txt","browser_download_url":"checksums"}]}`)},
		"asset":     {body: archive},
		"checksums": {body: manifest},
	}}
}

func isolateUpdateHooks(t *testing.T) {
	t.Helper()
	originalMax := maxDownload
	originalOS, originalArch := runtimeOS, runtimeArch
	originalMake, originalRemove := makeUpdateDirectory, removeUpdateDirectory
	originalExecutable, originalEvaluate := currentExecutable, evaluateSymlinks
	originalStaged, originalSource := createStagedFile, openUpdateSource
	originalDownload, originalRequest := createDownloadFile, newHTTPRequest
	originalReadChecksum, originalOpenChecksum := readChecksumFile, openChecksumFile
	originalOpenArchive, originalExtracted := openReleaseArchive, createExtractedFile
	t.Cleanup(func() {
		maxDownload = originalMax
		runtimeOS, runtimeArch = originalOS, originalArch
		makeUpdateDirectory, removeUpdateDirectory = originalMake, originalRemove
		currentExecutable, evaluateSymlinks = originalExecutable, originalEvaluate
		createStagedFile, openUpdateSource = originalStaged, originalSource
		createDownloadFile, newHTTPRequest = originalDownload, originalRequest
		readChecksumFile, openChecksumFile = originalReadChecksum, originalOpenChecksum
		openReleaseArchive, createExtractedFile = originalOpenArchive, originalExtracted
	})
}

func releaseArchive(t *testing.T, headers []tar.Header, payloads [][]byte) []byte {
	t.Helper()
	var output bytes.Buffer
	gzipWriter := gzip.NewWriter(&output)
	tarWriter := tar.NewWriter(gzipWriter)
	for index := range headers {
		if err := tarWriter.WriteHeader(&headers[index]); err != nil {
			t.Fatal(err)
		}
		if _, err := tarWriter.Write(payloads[index]); err != nil {
			t.Fatal(err)
		}
	}
	if err := tarWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gzipWriter.Close(); err != nil {
		t.Fatal(err)
	}
	return output.Bytes()
}
