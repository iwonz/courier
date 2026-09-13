package app

import (
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

type nilProvider struct{}

func (nilProvider) Command() *cobra.Command { return nil }

func TestCommandRegistry(t *testing.T) {
	provider := ProviderFunc(func() *cobra.Command { return &cobra.Command{Use: "sample"} })
	root, err := NewRootWithProviders(provider)
	if err != nil || len(root.Commands()) != 2 || root.Commands()[1].Name() != "sample" {
		t.Fatalf("root=%v err=%v", root, err)
	}
	if !root.CompletionOptions.DisableDefaultCmd {
		t.Fatal("completion command must be disabled")
	}
	if err := root.RunE(root, nil); err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name      string
		providers []CommandProvider
		text      string
	}{
		{"nil provider", []CommandProvider{nil}, "nil command provider"},
		{"nil command", []CommandProvider{nilProvider{}}, "invalid command"},
		{"empty command", []CommandProvider{ProviderFunc(func() *cobra.Command { return &cobra.Command{} })}, "invalid command"},
		{"duplicate", []CommandProvider{provider, provider}, "duplicate command"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, err := NewRootWithProviders(test.providers...); err == nil || !strings.Contains(err.Error(), test.text) {
				t.Fatalf("error=%v", err)
			}
		})
	}
}

func TestDefaultRootHasNoCompletion(t *testing.T) {
	root := NewRoot(Dependencies{})
	for _, command := range root.Commands() {
		if command.Name() == "completion" {
			t.Fatal("unexpected completion command")
		}
	}
}

func TestMustRoot(t *testing.T) {
	root := &cobra.Command{Use: "root"}
	if mustRoot(root, nil) != root {
		t.Fatal("unexpected root")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	mustRoot(nil, errors.New("broken registry"))
}
