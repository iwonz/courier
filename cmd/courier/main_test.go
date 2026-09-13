package main

import (
	"bytes"
	"context"
	"errors"
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
