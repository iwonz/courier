package selection

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/iwonz/courier/internal/operation"
)

func TestCompileOrderedRules(t *testing.T) {
	rules := []operation.SelectionRule{
		{Kind: operation.SelectionGitignore, Value: "*.tmp", Position: 0},
		{Kind: operation.SelectionFile, Value: "rules", Position: 1},
		{Kind: operation.SelectionGitignore, Value: "!later.tmp", Position: 2},
		{Kind: operation.SelectionRegex, Value: `(^|/)secret-[0-9]+$`, Position: 3},
	}
	opener := func(name string) (io.ReadCloser, error) {
		if name != "rules" {
			t.Fatalf("name=%q", name)
		}
		return io.NopCloser(strings.NewReader("# comment\n\n!keep.tmp\nprivate/\n")), nil
	}
	matcher, err := Compile(rules, opener)
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		path      string
		directory bool
		included  bool
	}{
		{"", true, true},
		{".", true, true},
		{"ordinary.txt", false, true},
		{"drop.tmp", false, false},
		{"keep.tmp", false, true},
		{"later.tmp", false, true},
		{"private", true, false},
		{"nested/secret-42", false, false},
		{"../escape", false, false},
		{"/absolute", false, false},
	} {
		if got := matcher.Include(test.path, test.directory); got != test.included {
			t.Errorf("Include(%q,%t)=%t want %t", test.path, test.directory, got, test.included)
		}
	}
}

func TestAllMatcher(t *testing.T) {
	matcher := All()
	if !matcher.Include("folder/file 🚚", false) || matcher.Include("a/../../escape", false) {
		t.Fatal("all matcher did not enforce relative containment")
	}
	compiled, err := Compile(nil, nil)
	if err != nil || !compiled.Include("anything", false) {
		t.Fatalf("compiled all=%+v err=%v", compiled, err)
	}
}

func TestCompileFailures(t *testing.T) {
	failingOpen := func(string) (io.ReadCloser, error) { return nil, errors.New("open") }
	for _, test := range []struct {
		name string
		rule operation.SelectionRule
		open func(string) (io.ReadCloser, error)
	}{
		{"regex", operation.SelectionRule{Kind: operation.SelectionRegex, Value: "[", Position: 1}, nil},
		{"missing opener", operation.SelectionRule{Kind: operation.SelectionFile, Value: "rules", Position: 2}, nil},
		{"open", operation.SelectionRule{Kind: operation.SelectionFile, Value: "rules", Position: 3}, failingOpen},
		{"unknown", operation.SelectionRule{Kind: "unknown", Position: 4}, nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := Compile([]operation.SelectionRule{test.rule}, test.open)
			if !errors.Is(err, ErrInvalidRule) || !strings.Contains(err.Error(), "position") {
				t.Fatalf("error=%v", err)
			}
		})
	}
}

func TestRuleFileFailures(t *testing.T) {
	t.Run("line too long", func(t *testing.T) {
		opener := readerOpener(strings.Repeat("x", maxRuleLineBytes+1), nil)
		if _, err := Compile([]operation.SelectionRule{{Kind: operation.SelectionFile, Value: "rules"}}, opener); !errors.Is(err, ErrInvalidRule) {
			t.Fatalf("error=%v", err)
		}
	})
	t.Run("file too large", func(t *testing.T) {
		line := strings.Repeat("x", maxRuleLineBytes-1) + "\n"
		opener := readerOpener(strings.Repeat(line, maxRuleFileBytes/len(line)+1), nil)
		if _, err := Compile([]operation.SelectionRule{{Kind: operation.SelectionFile, Value: "rules"}}, opener); !errors.Is(err, ErrInvalidRule) {
			t.Fatalf("error=%v", err)
		}
	})
	t.Run("close", func(t *testing.T) {
		opener := readerOpener("*.tmp\n", errors.New("close"))
		if _, err := Compile([]operation.SelectionRule{{Kind: operation.SelectionFile, Value: "rules"}}, opener); !errors.Is(err, ErrInvalidRule) || !strings.Contains(err.Error(), "close") {
			t.Fatalf("error=%v", err)
		}
	})
}

type closeReader struct {
	io.Reader
	err error
}

func (r closeReader) Close() error { return r.err }

func readerOpener(value string, closeErr error) func(string) (io.ReadCloser, error) {
	return func(string) (io.ReadCloser, error) {
		return closeReader{Reader: strings.NewReader(value), err: closeErr}, nil
	}
}

func TestOpenFile(t *testing.T) {
	name := filepath.Join(t.TempDir(), "rules")
	if err := os.WriteFile(name, []byte("rule"), 0o600); err != nil {
		t.Fatal(err)
	}
	file, err := OpenFile(name)
	if err != nil {
		t.Fatal(err)
	}
	_ = file.Close()
}
