package app

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/archive"
	"github.com/iwonz/courier/internal/endpoint"
	"github.com/iwonz/courier/internal/fsx"
	"github.com/iwonz/courier/internal/operation"
	"github.com/iwonz/courier/internal/policy"
	"github.com/iwonz/courier/internal/progress"
	"github.com/iwonz/courier/internal/report"
	"github.com/iwonz/courier/internal/selection"
	"github.com/iwonz/courier/internal/webhook"
)

type failingWriter struct{ err error }

func (writer failingWriter) Write([]byte) (int, error) { return 0, writer.err }

type denySelector struct{}

func (denySelector) Include(string, bool) bool { return false }

func outgoingPlan(t *testing.T, source, destination string, options ...operation.Option) operation.Plan {
	t.Helper()
	plan, err := operation.Build(operation.Request{Source: source, Destination: destination, Options: options})
	if err != nil {
		t.Fatal(err)
	}
	return plan
}

func outgoingDependencies(t *testing.T, sourcePath string) Dependencies {
	t.Helper()
	return Dependencies{
		Open: func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
			return &Resource{Endpoint: value, Backend: fsx.Local{}, Path: sourcePath}, nil
		},
		OpenArtifact: func(path string) (fsx.Backend, string, func() error, error) {
			return fsx.Local{}, path, func() error { return nil }, nil
		},
		Archive: archive.CreateSelected,
		Webhook: func(_ context.Context, request webhook.Request) (webhook.Result, error) {
			info, err := request.SourceFS.Lstat(request.SourcePath)
			return webhook.Result{StatusCode: 204, Bytes: info.Size(), Elapsed: time.Second}, err
		},
		DeliveryCredentials: func(context.Context, operation.AuthMode) (policy.Credentials, error) {
			return policy.Credentials{}, nil
		},
		Reporter: report.New,
		Terminal: func(io.Writer) bool { return false },
		TempDir:  t.TempDir(),
	}
}

func TestOutgoingWebhookDirectAndCommand(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "résumé.txt")
	if err := os.WriteFile(path, []byte("payload"), 0o600); err != nil {
		t.Fatal(err)
	}
	plan := outgoingPlan(t, path, "https://example.test/hook", operation.Option{Name: operation.OptionAuth, Value: "basic"})
	dependencies := outgoingDependencies(t, path)
	credentials := policy.Credentials{BasicUsername: "courier", BasicPassword: []byte("secret")}
	dependencies.DeliveryCredentials = func(_ context.Context, mode operation.AuthMode) (policy.Credentials, error) {
		if mode != operation.AuthBasic {
			t.Fatalf("mode=%s", mode)
		}
		return credentials, nil
	}
	dependencies.Webhook = func(_ context.Context, request webhook.Request) (webhook.Result, error) {
		if request.URL != plan.Destination.Raw || request.Name != "résumé.txt" || request.Username != "courier" || string(request.Password) != "secret" || !request.Unlimited {
			t.Fatalf("request=%+v", request)
		}
		return webhook.Result{StatusCode: 202, Bytes: 7, Elapsed: time.Second}, nil
	}
	var stdout, stderr bytes.Buffer
	if err := runOutgoingWebhook(context.Background(), dependencies, plan, selection.All(), &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "http-status: 202") || !strings.Contains(stdout.String(), "transferred: 7 bytes") || !bytes.Equal(credentials.BasicPassword, make([]byte, len(credentials.BasicPassword))) {
		t.Fatalf("stdout=%q credentials=%+v", stdout.String(), credentials)
	}

	dependencies = outgoingDependencies(t, path)
	var commandOutput bytes.Buffer
	if code := Execute(context.Background(), NewRoot(dependencies), []string{"from", path, "to", "https://example.test/hook"}, &commandOutput, io.Discard); code != ExitOK || !strings.Contains(commandOutput.String(), "http-status: 204") {
		t.Fatalf("code=%d output=%q", code, commandOutput.String())
	}

	remotePlan := outgoingPlan(t, "alias:/srv/résumé.txt", "https://example.test/hook")
	dependencies = outgoingDependencies(t, path)
	dependencies.Open = func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
		if !value.Remote || value.Host != "alias" {
			t.Fatalf("remote endpoint=%+v", value)
		}
		return &Resource{Endpoint: value, Backend: fsx.Local{}, Path: path}, nil
	}
	if err := runOutgoingWebhook(context.Background(), dependencies, remotePlan, selection.All(), io.Discard, io.Discard); err != nil {
		t.Fatalf("remote source=%v", err)
	}
}

func TestOutgoingWebhookArchiveAndCleanup(t *testing.T) {
	source := t.TempDir()
	if err := os.WriteFile(filepath.Join(source, "file.txt"), []byte("archive"), 0o600); err != nil {
		t.Fatal(err)
	}
	plan := outgoingPlan(t, source, "https://example.test/hook", operation.Option{Name: operation.OptionArchive, Value: "true"})
	dependencies := outgoingDependencies(t, source)
	var artifactPath string
	dependencies.Webhook = func(_ context.Context, request webhook.Request) (webhook.Result, error) {
		artifactPath = request.SourcePath
		if request.Name != filepath.Base(source)+".tar.gz" {
			t.Fatalf("name=%q", request.Name)
		}
		if err := archive.Verify(request.SourcePath); err != nil {
			t.Fatal(err)
		}
		info, err := request.SourceFS.Lstat(request.SourcePath)
		return webhook.Result{StatusCode: 200, Bytes: info.Size()}, err
	}
	if err := runOutgoingWebhook(context.Background(), dependencies, plan, selection.All(), io.Discard, io.Discard); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(artifactPath); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("archive remained: %v", err)
	}
}

func TestOutgoingWebhookPreflightFailures(t *testing.T) {
	file := filepath.Join(t.TempDir(), "file.txt")
	if err := os.WriteFile(file, []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}
	filePlan := outgoingPlan(t, file, "https://example.test/hook")
	if err := runOutgoingWebhook(context.Background(), Dependencies{}, filePlan, selection.All(), io.Discard, io.Discard); err == nil {
		t.Fatal("incomplete dependencies accepted")
	}
	dependencies := outgoingDependencies(t, file)
	want := errors.New("failure")
	dependencies.DeliveryCredentials = func(context.Context, operation.AuthMode) (policy.Credentials, error) {
		return policy.Credentials{}, want
	}
	if err := runOutgoingWebhook(context.Background(), dependencies, filePlan, selection.All(), io.Discard, io.Discard); !errors.Is(err, want) {
		t.Fatalf("credential failure=%v", err)
	}

	for _, remote := range []bool{false, true} {
		dependencies = outgoingDependencies(t, file)
		dependencies.Open = func(context.Context, endpoint.Endpoint) (*Resource, error) { return nil, want }
		plan := filePlan
		plan.Source.Remote = remote
		err := runOutgoingWebhook(context.Background(), dependencies, plan, selection.All(), io.Discard, io.Discard)
		var commandErr *commandError
		if !errors.As(err, &commandErr) || commandErr.code != map[bool]int{false: ExitTransfer, true: ExitConnection}[remote] {
			t.Fatalf("remote=%v error=%v", remote, err)
		}
	}

	dependencies = outgoingDependencies(t, filepath.Join(t.TempDir(), "missing"))
	if err := runOutgoingWebhook(context.Background(), dependencies, filePlan, selection.All(), io.Discard, io.Discard); err == nil {
		t.Fatal("lstat failure ignored")
	}
	directory := t.TempDir()
	directoryPlan := outgoingPlan(t, directory, "https://example.test/hook")
	dependencies = outgoingDependencies(t, directory)
	if err := runOutgoingWebhook(context.Background(), dependencies, directoryPlan, selection.All(), io.Discard, io.Discard); err == nil {
		t.Fatal("directory without archive accepted")
	}
	dependencies = outgoingDependencies(t, file)
	dependencies.Open = func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
		backend := lstatBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) {
			return fakeInfo{name: "pipe", directory: false, size: 0, mode: fs.ModeNamedPipe}, nil
		}}
		return &Resource{Endpoint: value, Backend: backend, Path: file}, nil
	}
	if err := runOutgoingWebhook(context.Background(), dependencies, filePlan, selection.All(), io.Discard, io.Discard); err == nil {
		t.Fatal("special source accepted")
	}
	if err := runOutgoingWebhook(context.Background(), outgoingDependencies(t, file), filePlan, denySelector{}, io.Discard, io.Discard); err == nil {
		t.Fatal("excluded source accepted")
	}
}

func TestOutgoingWebhookArchiveAndRuntimeFailures(t *testing.T) {
	source := t.TempDir()
	plan := outgoingPlan(t, source, "https://example.test/hook", operation.Option{Name: operation.OptionArchive, Value: "true"})
	want := errors.New("failure")
	dependencies := outgoingDependencies(t, source)
	dependencies.Archive = func(context.Context, fsx.Backend, string, string, string, selection.Selector, progress.Sink) (*archive.Artifact, error) {
		return nil, want
	}
	if err := runOutgoingWebhook(context.Background(), dependencies, plan, selection.All(), io.Discard, io.Discard); !errors.Is(err, want) {
		t.Fatalf("archive failure=%v", err)
	}

	artifactFile := filepath.Join(t.TempDir(), "artifact.tar.gz")
	if err := os.WriteFile(artifactFile, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	dependencies = outgoingDependencies(t, source)
	dependencies.Archive = func(context.Context, fsx.Backend, string, string, string, selection.Selector, progress.Sink) (*archive.Artifact, error) {
		return &archive.Artifact{Path: artifactFile, Name: "source.tar.gz", Bytes: 1}, nil
	}
	dependencies.OpenArtifact = func(string) (fsx.Backend, string, func() error, error) { return nil, "", nil, want }
	if err := runOutgoingWebhook(context.Background(), dependencies, plan, selection.All(), io.Discard, io.Discard); !errors.Is(err, want) {
		t.Fatalf("open artifact failure=%v", err)
	}

	file := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	filePlan := outgoingPlan(t, file, "https://example.test/hook")
	dependencies = outgoingDependencies(t, file)
	dependencies.Webhook = func(context.Context, webhook.Request) (webhook.Result, error) {
		return webhook.Result{}, &url.Error{Op: "Post", URL: filePlan.Destination.Raw, Err: want}
	}
	err := runOutgoingWebhook(context.Background(), dependencies, filePlan, selection.All(), io.Discard, io.Discard)
	var commandErr *commandError
	if !errors.As(err, &commandErr) || commandErr.code != ExitConnection {
		t.Fatalf("connection failure=%v", err)
	}
	dependencies.Webhook = func(context.Context, webhook.Request) (webhook.Result, error) { return webhook.Result{Bytes: 1}, want }
	if err := runOutgoingWebhook(context.Background(), dependencies, filePlan, selection.All(), io.Discard, io.Discard); !errors.Is(err, want) {
		t.Fatalf("transfer failure=%v", err)
	}
	dependencies.Webhook = func(context.Context, webhook.Request) (webhook.Result, error) {
		return webhook.Result{StatusCode: 200, Bytes: 1}, nil
	}
	if err := runOutgoingWebhook(context.Background(), dependencies, filePlan, selection.All(), failingWriter{err: want}, io.Discard); !errors.Is(err, want) {
		t.Fatalf("output failure=%v", err)
	}

	dependencies = outgoingDependencies(t, file)
	dependencies.Open = func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
		return &Resource{Endpoint: value, Backend: fsx.Local{}, Path: file, Close: func() error { return want }}, nil
	}
	if err := runOutgoingWebhook(context.Background(), dependencies, filePlan, selection.All(), io.Discard, io.Discard); !errors.Is(err, want) {
		t.Fatalf("source cleanup failure=%v", err)
	}

	missingArtifact := filepath.Join(t.TempDir(), "missing.tar.gz")
	dependencies = outgoingDependencies(t, source)
	dependencies.Archive = func(context.Context, fsx.Backend, string, string, string, selection.Selector, progress.Sink) (*archive.Artifact, error) {
		return &archive.Artifact{Path: missingArtifact, Name: "source.tar.gz", Bytes: 1}, nil
	}
	dependencies.OpenArtifact = func(string) (fsx.Backend, string, func() error, error) {
		return fsx.Local{}, missingArtifact, func() error { return want }, nil
	}
	dependencies.Webhook = func(context.Context, webhook.Request) (webhook.Result, error) {
		return webhook.Result{StatusCode: 200, Bytes: 1}, nil
	}
	if err := runOutgoingWebhook(context.Background(), dependencies, plan, selection.All(), io.Discard, io.Discard); !errors.Is(err, want) || !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("archive cleanup failures=%v", err)
	}
}
