package helper

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/iwonz/courier/internal/sshx"
	"github.com/pkg/sftp"
)

func TestManagerFallbackPolicyAndSuccess(t *testing.T) {
	capability := &sshx.CapabilityError{Capability: "sftp", Cause: errors.New("subsystem")}
	connection := &sshx.Connection{Target: sshx.Target{Alias: "server"}}
	detect := func(context.Context, sshx.Runner) (sshx.Platform, sshx.Archiver, error) {
		return sshx.Platform{OS: "linux", Arch: "amd64"}, sshx.Archiver{}, nil
	}
	if _, _, err := (*Manager)(nil).Fallback(context.Background(), connection, capability); err == nil {
		t.Fatal("expected unavailable consent")
	}
	manager := &Manager{Confirm: func(context.Context, string) (bool, error) { return false, nil }, Detect: detect}
	if _, _, err := manager.Fallback(context.Background(), connection, &sshx.CapabilityError{Capability: "exec"}); err == nil {
		t.Fatal("expected capability rejection")
	}
	if _, _, err := manager.Fallback(context.Background(), connection, capability); err == nil || !strings.Contains(err.Error(), "declined") {
		t.Fatalf("expected declined error: %v", err)
	}
	originalDetect := detectRemotePlatform
	t.Cleanup(func() { detectRemotePlatform = originalDetect })
	probe := 0
	detectRemotePlatform = func(context.Context, sshx.Runner) (sshx.Platform, sshx.Archiver, error) {
		probe++
		return sshx.Platform{OS: "linux", Arch: "amd64"}, sshx.Archiver{}, nil
	}
	defaultConnection := &sshx.Connection{Target: sshx.Target{Alias: "server"}}
	defaultManager := &Manager{Confirm: func(context.Context, string) (bool, error) { return false, nil }}
	if _, _, err := defaultManager.Fallback(context.Background(), defaultConnection, capability); err == nil || probe != 1 {
		t.Fatalf("default detector probes=%d err=%v", probe, err)
	}
	manager.Confirm = func(context.Context, string) (bool, error) { return false, errors.New("prompt") }
	if _, _, err := manager.Fallback(context.Background(), connection, capability); err == nil || !strings.Contains(err.Error(), "consent") {
		t.Fatalf("expected consent error: %v", err)
	}
	manager.Confirm = func(context.Context, string) (bool, error) { return true, nil }
	manager.Detect = func(context.Context, sshx.Runner) (sshx.Platform, sshx.Archiver, error) {
		return sshx.Platform{}, sshx.Archiver{}, errors.New("detect")
	}
	if _, _, err := manager.Fallback(context.Background(), connection, capability); err == nil || !strings.Contains(err.Error(), "detect") {
		t.Fatalf("expected detection error: %v", err)
	}
	manager.Detect = detect
	if _, _, err := manager.Fallback(context.Background(), connection, capability); err == nil || !strings.Contains(err.Error(), "incomplete") {
		t.Fatalf("expected dependency error: %v", err)
	}
	manager.Acquire = func(context.Context, sshx.Platform) (*Artifact, error) { return nil, errors.New("acquire") }
	manager.Deploy = sshx.DeployHelper
	if _, _, err := manager.Fallback(context.Background(), connection, capability); err == nil || !strings.Contains(err.Error(), "acquire") {
		t.Fatalf("expected acquisition error: %v", err)
	}
	manager.Acquire = func(context.Context, sshx.Platform) (*Artifact, error) { return nil, nil }
	if _, _, err := manager.Fallback(context.Background(), connection, capability); err == nil || !strings.Contains(err.Error(), "incomplete artifact") {
		t.Fatalf("expected artifact error: %v", err)
	}

	cleanups := 0
	runtime := &countingCloser{}
	client := testSFTPClient(t)
	manager.Acquire = func(context.Context, sshx.Platform) (*Artifact, error) {
		return &Artifact{Path: "helper", SHA256: strings.Repeat("0", 64), Cleanup: func() error { cleanups++; return nil }}, nil
	}
	manager.Deploy = func(_ context.Context, _ *sshx.Connection, platform sshx.Platform, path, digest string) (*sftp.Client, io.Closer, error) {
		if platform.OS != "linux" || path != "helper" || len(digest) != 64 {
			t.Fatalf("unexpected deployment: %+v %q %q", platform, path, digest)
		}
		return client, runtime, nil
	}
	gotClient, gotRuntime, err := manager.Fallback(context.Background(), connection, capability)
	if err != nil || gotClient != client || gotRuntime != runtime || cleanups != 1 {
		t.Fatalf("client=%v runtime=%v cleanups=%d err=%v", gotClient, gotRuntime, cleanups, err)
	}
	_ = client.Close()
}

func TestManagerFallbackCleanupFailures(t *testing.T) {
	capability := &sshx.CapabilityError{Capability: "sftp", Cause: errors.New("subsystem")}
	platform := sshx.Platform{OS: "linux", Arch: "amd64"}
	base := &Manager{
		Confirm: func(context.Context, string) (bool, error) { return true, nil },
		Detect: func(context.Context, sshx.Runner) (sshx.Platform, sshx.Archiver, error) {
			return platform, sshx.Archiver{}, nil
		},
	}
	artifact := func(cleanup error) Acquire {
		return func(context.Context, sshx.Platform) (*Artifact, error) {
			return &Artifact{Path: "helper", SHA256: strings.Repeat("0", 64), Cleanup: func() error { return cleanup }}, nil
		}
	}

	runtime := &countingCloser{}
	base.Acquire = artifact(nil)
	base.Deploy = func(context.Context, *sshx.Connection, sshx.Platform, string, string) (*sftp.Client, io.Closer, error) {
		return nil, runtime, errors.New("deploy")
	}
	if _, _, err := base.Fallback(context.Background(), &sshx.Connection{}, capability); err == nil || runtime.count != 1 {
		t.Fatalf("expected deployment cleanup: count=%d err=%v", runtime.count, err)
	}

	runtime = &countingCloser{}
	base.Acquire = artifact(errors.New("local cleanup"))
	base.Deploy = func(context.Context, *sshx.Connection, sshx.Platform, string, string) (*sftp.Client, io.Closer, error) {
		return nil, runtime, nil
	}
	if _, _, err := base.Fallback(context.Background(), &sshx.Connection{}, capability); err == nil || runtime.count != 1 {
		t.Fatalf("expected incomplete runtime cleanup: count=%d err=%v", runtime.count, err)
	}

	runtime = &countingCloser{}
	client := testSFTPClient(t)
	base.Acquire = artifact(errors.New("local cleanup"))
	base.Deploy = func(context.Context, *sshx.Connection, sshx.Platform, string, string) (*sftp.Client, io.Closer, error) {
		return client, runtime, nil
	}
	if _, _, err := base.Fallback(context.Background(), &sshx.Connection{}, capability); err == nil || runtime.count != 1 {
		t.Fatalf("expected client and runtime cleanup: count=%d err=%v", runtime.count, err)
	}
}

func TestSourceLocalAndSelection(t *testing.T) {
	root := t.TempDir()
	binary := filepath.Join(root, "courier")
	if err := os.WriteFile(binary, []byte("helper"), 0o700); err != nil {
		t.Fatal(err)
	}
	source := Source{GOOS: "linux", GOARCH: "amd64", Executable: func() (string, error) { return binary, nil }, EvalLinks: func(name string) (string, error) { return name, nil }}
	artifact, err := source.Acquire(context.Background(), sshx.Platform{OS: "linux", Arch: "amd64"})
	if err != nil || artifact.Path != binary || len(artifact.SHA256) != 64 {
		t.Fatalf("artifact=%+v err=%v", artifact, err)
	}
	if err := artifact.Cleanup(); err != nil {
		t.Fatal(err)
	}
	defaultArtifact, err := (Source{}).Acquire(context.Background(), sshx.Platform{OS: runtime.GOOS, Arch: runtime.GOARCH})
	if err != nil || defaultArtifact.Path == "" || len(defaultArtifact.SHA256) != 64 {
		t.Fatalf("default artifact=%+v err=%v", defaultArtifact, err)
	}
	for _, platform := range []sshx.Platform{{OS: "plan9", Arch: "amd64"}, {OS: "linux", Arch: "386"}, {OS: "dragonfly", Arch: "arm64"}} {
		if _, err := source.Acquire(context.Background(), platform); err == nil {
			t.Fatalf("expected unsupported platform %+v", platform)
		}
	}
	for _, platform := range []sshx.Platform{{OS: "darwin", Arch: "arm64"}, {OS: "windows", Arch: "amd64"}, {OS: "freebsd", Arch: "arm64"}, {OS: "openbsd", Arch: "amd64"}, {OS: "netbsd", Arch: "arm64"}, {OS: "dragonfly", Arch: "amd64"}} {
		if !supportedHelperPlatform(platform) {
			t.Fatalf("expected supported platform %+v", platform)
		}
	}

	source.Executable = func() (string, error) { return "", errors.New("executable") }
	if _, err := source.Acquire(context.Background(), sshx.Platform{OS: "linux", Arch: "amd64"}); err == nil {
		t.Fatal("expected executable error")
	}
	source.Executable = func() (string, error) { return binary, nil }
	source.EvalLinks = func(string) (string, error) { return "", errors.New("links") }
	if _, err := source.Acquire(context.Background(), sshx.Platform{OS: "linux", Arch: "amd64"}); err == nil {
		t.Fatal("expected symlink error")
	}
	source.EvalLinks = func(name string) (string, error) { return name, nil }
	if _, err := (Source{GOOS: "linux", GOARCH: "arm64", Version: "dev"}).Acquire(context.Background(), sshx.Platform{OS: "windows", Arch: "amd64"}); err == nil {
		t.Fatal("expected immutable release error")
	}
	archive := sourceReleaseArchive(t, "courier.exe", []byte("windows helper"))
	archiveHash := sha256.Sum256(archive)
	asset := "courier_1.0.0_windows_amd64.tar.gz"
	base := "https://github.com/owner/repo/releases/download/v1.0.0/"
	cross := Source{
		Repository: "owner/repo", Version: "v1.0.0", GOOS: "linux", GOARCH: "arm64",
		Client: sourceHTTPClient{responses: map[string][]byte{base + asset: archive, base + "checksums.txt": []byte(hex.EncodeToString(archiveHash[:]) + "  " + asset)}},
	}
	crossArtifact, err := cross.Acquire(context.Background(), sshx.Platform{OS: "windows", Arch: "amd64"})
	if err != nil || len(crossArtifact.SHA256) != 64 {
		t.Fatalf("cross artifact=%+v err=%v", crossArtifact, err)
	}
	if err := crossArtifact.Cleanup(); err != nil {
		t.Fatal(err)
	}
}

func TestArtifactDigestFailures(t *testing.T) {
	original := openArtifactFile
	t.Cleanup(func() { openArtifactFile = original })
	openArtifactFile = func(string) (io.ReadCloser, error) { return nil, errors.New("open") }
	if _, err := artifactDigest("helper"); err == nil {
		t.Fatal("expected open error")
	}
	source := Source{GOOS: "linux", GOARCH: "amd64", Executable: func() (string, error) { return "helper", nil }, EvalLinks: func(name string) (string, error) { return name, nil }}
	if _, err := source.Acquire(context.Background(), sshx.Platform{OS: "linux", Arch: "amd64"}); err == nil {
		t.Fatal("expected source digest error")
	}
	openArtifactFile = func(string) (io.ReadCloser, error) { return &failingArtifactReader{readErr: errors.New("read")}, nil }
	if _, err := artifactDigest("helper"); err == nil {
		t.Fatal("expected read error")
	}
	openArtifactFile = func(string) (io.ReadCloser, error) {
		return &failingArtifactReader{Reader: bytes.NewBufferString("x"), closeErr: errors.New("close")}, nil
	}
	if _, err := artifactDigest("helper"); err == nil {
		t.Fatal("expected close error")
	}
}

type countingCloser struct{ count int }

func (c *countingCloser) Close() error { c.count++; return nil }

type failingArtifactReader struct {
	io.Reader
	readErr, closeErr error
}

func (r *failingArtifactReader) Read(data []byte) (int, error) {
	if r.readErr != nil {
		return 0, r.readErr
	}
	return r.Reader.Read(data)
}
func (r *failingArtifactReader) Close() error { return r.closeErr }

func testSFTPClient(t *testing.T) *sftp.Client {
	t.Helper()
	serverSide, clientSide := net.Pipe()
	server := sftp.NewRequestServer(serverSide, sftp.InMemHandler())
	go func() {
		_ = server.Serve()
		_ = server.Close()
	}()
	client, err := sftp.NewClientPipe(clientSide, clientSide)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.Close(); _ = serverSide.Close() })
	return client
}

type sourceHTTPClient struct{ responses map[string][]byte }

func (s sourceHTTPClient) Do(request *http.Request) (*http.Response, error) {
	data, ok := s.responses[request.URL.String()]
	if !ok {
		return nil, errors.New("unexpected release URL")
	}
	return &http.Response{StatusCode: http.StatusOK, Status: "OK", ContentLength: int64(len(data)), Body: io.NopCloser(bytes.NewReader(data))}, nil
}

func sourceReleaseArchive(t *testing.T, name string, data []byte) []byte {
	t.Helper()
	var output bytes.Buffer
	gzipWriter := gzip.NewWriter(&output)
	tarWriter := tar.NewWriter(gzipWriter)
	if err := tarWriter.WriteHeader(&tar.Header{Name: name, Mode: 0o700, Size: int64(len(data)), Typeflag: tar.TypeReg}); err != nil {
		t.Fatal(err)
	}
	if _, err := tarWriter.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := tarWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gzipWriter.Close(); err != nil {
		t.Fatal(err)
	}
	return output.Bytes()
}
