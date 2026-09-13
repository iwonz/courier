package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/iwonz/courier/internal/app"
)

func TestRunAndMain(t *testing.T) {
	originalFactory, originalExit, originalArgs := dependencyFactory, exitProcess, os.Args
	t.Cleanup(func() { dependencyFactory, exitProcess, os.Args = originalFactory, originalExit, originalArgs })
	dependencyFactory = func(*os.File, io.Writer) (app.Dependencies, error) {
		return app.Dependencies{Build: app.BuildIdentity{Version: "dev", Commit: "unknown", Date: "unknown"}}, nil
	}
	var stdout, stderr bytes.Buffer
	if code := run(context.Background(), []string{"version"}, nil, &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), "courier dev") {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	os.Args = []string{"courier", "version"}
	code := -1
	exitProcess = func(value int) { code = value }
	main()
	if code != 0 {
		t.Fatalf("main code=%d", code)
	}
}

func TestDependencyFailure(t *testing.T) {
	original := dependencyFactory
	t.Cleanup(func() { dependencyFactory = original })
	dependencyFactory = func(*os.File, io.Writer) (app.Dependencies, error) { return app.Dependencies{}, errors.New("config") }
	var stderr bytes.Buffer
	if code := run(context.Background(), nil, nil, io.Discard, &stderr); code != app.ExitConnection || !strings.Contains(stderr.String(), "config") {
		t.Fatalf("code=%d stderr=%q", code, stderr.String())
	}
}

func TestInternalFailureRedactionAndInterruption(t *testing.T) {
	var stderr bytes.Buffer
	if code := internalFailure(&stderr, app.ExitControl, "control", fmt.Errorf("token=unsafe: %w", context.Canceled)); code != app.ExitInterrupted {
		t.Fatalf("code=%d", code)
	}
	if output := stderr.String(); strings.Contains(output, "unsafe") || !strings.Contains(output, "token=[REDACTED]") {
		t.Fatalf("stderr=%q", output)
	}
}

func TestInternalModes(t *testing.T) {
	originalUpdate, originalAdmin, originalHelper, originalWorker := runInternalUpdate, runInternalAdmin, serveHelperSFTP, runInternalWorker
	t.Cleanup(func() {
		runInternalUpdate, runInternalAdmin, serveHelperSFTP, runInternalWorker = originalUpdate, originalAdmin, originalHelper, originalWorker
	})
	runInternalUpdate = func([]string) (bool, error) { return true, nil }
	if code := run(context.Background(), []string{"internal"}, nil, io.Discard, io.Discard); code != app.ExitOK {
		t.Fatalf("internal success code=%d", code)
	}
	runInternalUpdate = func([]string) (bool, error) { return true, errors.New("handoff") }
	var stderr bytes.Buffer
	if code := run(context.Background(), []string{"internal"}, nil, io.Discard, &stderr); code != app.ExitUpdate || !strings.Contains(stderr.String(), "handoff") {
		t.Fatalf("internal failure code=%d stderr=%q", code, stderr.String())
	}
	runInternalUpdate = func([]string) (bool, error) { return false, nil }
	runInternalAdmin = func(context.Context, []string) (bool, error) { return true, nil }
	if code := run(context.Background(), []string{"_admin"}, nil, io.Discard, io.Discard); code != app.ExitOK {
		t.Fatalf("admin success code=%d", code)
	}
	runInternalAdmin = func(context.Context, []string) (bool, error) { return true, errors.New("admin") }
	stderr.Reset()
	if code := run(context.Background(), []string{"_admin"}, nil, io.Discard, &stderr); code != app.ExitControl || !strings.Contains(stderr.String(), "admin") {
		t.Fatalf("admin failure code=%d stderr=%q", code, stderr.String())
	}
	runInternalAdmin = originalAdmin
	serveHelperSFTP = func(io.Reader, io.Writer) error { return nil }
	if code := run(context.Background(), []string{"_helper-sftp"}, nil, io.Discard, io.Discard); code != app.ExitOK {
		t.Fatalf("helper success code=%d", code)
	}
	serveHelperSFTP = func(io.Reader, io.Writer) error { return errors.New("helper") }
	stderr.Reset()
	if code := run(context.Background(), []string{"_helper-sftp"}, nil, io.Discard, &stderr); code != app.ExitTransfer || !strings.Contains(stderr.String(), "helper") {
		t.Fatalf("helper failure code=%d stderr=%q", code, stderr.String())
	}
	serveHelperSFTP = originalHelper
	runInternalWorker = func(context.Context, []string) (bool, error) { return true, nil }
	if code := run(context.Background(), []string{"_worker"}, nil, io.Discard, io.Discard); code != app.ExitOK {
		t.Fatalf("worker success code=%d", code)
	}
	runInternalWorker = func(context.Context, []string) (bool, error) { return true, errors.New("worker") }
	stderr.Reset()
	if code := run(context.Background(), []string{"_worker"}, nil, io.Discard, &stderr); code != app.ExitControl || !strings.Contains(stderr.String(), "worker") {
		t.Fatalf("worker failure code=%d stderr=%q", code, stderr.String())
	}
}
