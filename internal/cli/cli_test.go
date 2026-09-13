package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestRegistry(t *testing.T) {
	called := false
	command := CommandFunc{CommandName: "z", CommandSynopsis: "last", Execute: func(_ context.Context, args []string, stdout, _ io.Writer) error {
		called = len(args) == 1 && args[0] == "arg"
		_, _ = io.WriteString(stdout, "done")
		return nil
	}}
	other := CommandFunc{CommandName: "a", CommandSynopsis: "first", Execute: func(context.Context, []string, io.Writer, io.Writer) error { return nil }}
	r, err := NewRegistry("courier from <source> to <destination> [flags]", command, other)
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := r.Run(context.Background(), []string{"z", "arg"}, &output, io.Discard); err != nil || !called || output.String() != "done" {
		t.Fatalf("dispatch failed: called=%v output=%q err=%v", called, output.String(), err)
	}
	output.Reset()
	r.Help(&output)
	if got := output.String(); !strings.Contains(got, "courier from") || strings.Index(got, "a ") > strings.Index(got, "z ") {
		t.Fatalf("unexpected help: %q", got)
	}
}

func TestRegistryHelpFormsAndErrors(t *testing.T) {
	r, err := NewRegistry("usage")
	if err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{nil, {"help"}, {"-h"}, {"--help"}} {
		var output bytes.Buffer
		if err := r.Run(context.Background(), args, &output, io.Discard); err != nil || output.String() != "usage\n" {
			t.Fatalf("args=%v output=%q err=%v", args, output.String(), err)
		}
	}
	if err := r.Run(context.Background(), []string{"missing"}, io.Discard, io.Discard); !errors.Is(err, ErrUsage) {
		t.Fatalf("expected usage error, got %v", err)
	}
	if _, err := NewRegistry("usage", nil); !errors.Is(err, ErrUsage) {
		t.Fatalf("expected nil command rejection, got %v", err)
	}
	empty := CommandFunc{}
	if _, err := NewRegistry("usage", empty); !errors.Is(err, ErrUsage) {
		t.Fatalf("expected empty command rejection, got %v", err)
	}
	duplicate := CommandFunc{CommandName: "same"}
	if _, err := NewRegistry("usage", duplicate, duplicate); !errors.Is(err, ErrUsage) {
		t.Fatalf("expected duplicate rejection, got %v", err)
	}
}
