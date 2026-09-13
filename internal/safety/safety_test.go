package safety

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/iwonz/courier/internal/endpoint"
)

func local(value string) endpoint.Endpoint { return endpoint.Endpoint{Path: value, Raw: value} }
func remote(user, host, value string) endpoint.Endpoint {
	return endpoint.Endpoint{User: user, Host: host, Path: value, Remote: true}
}

func TestValidateTransfer(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	if err := os.Mkdir(source, 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "alias")
	if err := os.Symlink(source, link); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name   string
		source endpoint.Endpoint
		dest   endpoint.Endpoint
		dir    bool
		unsafe bool
	}{
		{"different kinds", local(source), remote("", "host", "/x"), true, false},
		{"different remotes", remote("me", "one", "/x"), remote("me", "two", "/x"), true, false},
		{"different users", remote("me", "one", "/x"), remote("you", "one", "/x"), true, false},
		{"remote same", remote("me", "HOST", "/x/../x"), remote("me", "host", "/x"), false, true},
		{"remote descendant dir", remote("me", "host", "/x"), remote("me", "host", "/x/y"), true, true},
		{"remote descendant file", remote("me", "host", "/x"), remote("me", "host", "/x/y"), false, false},
		{"remote sibling", remote("me", "host", "/x"), remote("me", "host", "/xy"), true, false},
		{"local same through link", local(source), local(link), true, true},
		{"local descendant missing", local(source), local(filepath.Join(link, "new", "child")), true, true},
		{"local sibling", local(source), local(filepath.Join(root, "other")), true, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateTransfer(test.source, test.dest, test.dir)
			if errors.Is(err, ErrUnsafe) != test.unsafe {
				t.Fatalf("error=%v unsafe=%v", err, test.unsafe)
			}
		})
	}
}

func TestCanonicalErrors(t *testing.T) {
	originalAbs, originalStat, originalEval, originalRel := absLocal, statLocal, evalLocal, relLocal
	originalCase := localCaseInsensitive
	t.Cleanup(func() {
		absLocal, statLocal, evalLocal, relLocal = originalAbs, originalStat, originalEval, originalRel
		localCaseInsensitive = originalCase
	})
	absLocal = func(string) (string, error) { return "", errors.New("abs") }
	if err := ValidateTransfer(local("x"), local("y"), false); !errors.Is(err, ErrUnsafe) {
		t.Fatal(err)
	}
	calls := 0
	absLocal = func(value string) (string, error) {
		calls++
		if calls == 2 {
			return "", errors.New("destination abs")
		}
		return originalAbs(value)
	}
	if err := ValidateTransfer(local("x"), local("y"), false); !errors.Is(err, ErrUnsafe) {
		t.Fatal(err)
	}
	absLocal = originalAbs
	statLocal = func(string) (os.FileInfo, error) { return nil, errors.New("stat") }
	if err := ValidateTransfer(local("x"), local("y"), false); !errors.Is(err, ErrUnsafe) {
		t.Fatal(err)
	}
	statLocal = func(string) (os.FileInfo, error) { return nil, nil }
	evalLocal = func(string) (string, error) { return "", errors.New("eval") }
	if err := ValidateTransfer(local("x"), local("y"), false); !errors.Is(err, ErrUnsafe) {
		t.Fatal(err)
	}
	statLocal, evalLocal = originalStat, originalEval
	localCaseInsensitive = true
	if err := ValidateTransfer(local("MixedCase"), local("mixedcase"), false); !errors.Is(err, ErrUnsafe) {
		t.Fatalf("expected case-insensitive match, got %v", err)
	}
	localCaseInsensitive = false
	relLocal = func(string, string) (string, error) { return "", errors.New("rel") }
	if err := ValidateTransfer(local("x"), local("y"), false); !errors.Is(err, ErrUnsafe) {
		t.Fatalf("expected relative path failure, got %v", err)
	}
	relLocal = originalRel
	statLocal = func(string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	if got, err := canonicalLocal(string(filepath.Separator)); err != nil || got == "" {
		t.Fatalf("root fallback=%q, %v", got, err)
	}
}

func TestSafeArchiveJoin(t *testing.T) {
	root := t.TempDir()
	got, err := SafeArchiveJoin(root, "folder/file-🚚.txt")
	if err != nil || got != filepath.Join(root, "folder", "file-🚚.txt") {
		t.Fatalf("got=%q err=%v", got, err)
	}
	for _, name := range []string{"", ".", "..", "../secret", "a/../../secret", "/absolute", `C:\escape`, "bad\x00name"} {
		if _, err := SafeArchiveJoin(root, name); !errors.Is(err, ErrUnsafe) {
			t.Errorf("name=%q error=%v", name, err)
		}
	}
	original := relArchive
	t.Cleanup(func() { relArchive = original })
	relArchive = func(string, string) (string, error) { return "", errors.New("rel") }
	if _, err := SafeArchiveJoin(root, "entry"); !errors.Is(err, ErrUnsafe) {
		t.Fatalf("expected relative path failure, got %v", err)
	}
}

func TestRemoteRootRelationship(t *testing.T) {
	if err := ValidateTransfer(remote("me", "host", "/"), remote("me", "host", "/child"), true); !errors.Is(err, ErrUnsafe) {
		t.Fatalf("expected root descendant rejection, got %v", err)
	}
}
