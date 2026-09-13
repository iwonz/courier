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

func TestEvaluateTransfer(t *testing.T) {
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
		mode   Transformation
		want   Disposition
		unsafe bool
	}{
		{"different kinds", local(source), remote("", "host", "/x"), true, TransformNone, Proceed, false},
		{"different remotes", remote("me", "one", "/x"), remote("me", "two", "/x"), true, TransformNone, Proceed, false},
		{"different users", remote("me", "one", "/x"), remote("you", "one", "/x"), true, TransformNone, Proceed, false},
		{"remote same", remote("me", "HOST", "/x/../x"), remote("me", "host", "/x"), false, TransformNone, NoOp, false},
		{"remote archive collision", remote("me", "host", "/x"), remote("me", "host", "/x"), false, TransformArchive, Proceed, true},
		{"remote extract collision", remote("me", "host", "/x"), remote("me", "host", "/x"), false, TransformExtract, Proceed, true},
		{"remote descendant dir", remote("me", "host", "/x"), remote("me", "host", "/x/y"), true, TransformNone, Proceed, true},
		{"remote descendant file", remote("me", "host", "/x"), remote("me", "host", "/x/y"), false, TransformNone, Proceed, false},
		{"remote sibling", remote("me", "host", "/x"), remote("me", "host", "/xy"), true, TransformNone, Proceed, false},
		{"remote Windows same", endpoint.Endpoint{User: "me", Host: "host", Path: "C:/Data", PathFlavor: endpoint.PathWindows, Remote: true}, endpoint.Endpoint{User: "me", Host: "host", Path: "c:/data", PathFlavor: endpoint.PathWindows, Remote: true}, false, TransformNone, NoOp, false},
		{"local same through link", local(source), local(link), true, TransformNone, NoOp, false},
		{"local descendant missing", local(source), local(filepath.Join(link, "new", "child")), true, TransformNone, Proceed, true},
		{"local sibling", local(source), local(filepath.Join(root, "other")), true, TransformNone, Proceed, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := EvaluateTransfer(test.source, test.dest, test.dir, test.mode)
			if got != test.want || errors.Is(err, ErrUnsafe) != test.unsafe {
				t.Fatalf("disposition=%v error=%v unsafe=%v", got, err, test.unsafe)
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
	if _, err := EvaluateTransfer(local("x"), local("y"), false, TransformNone); !errors.Is(err, ErrUnsafe) {
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
	if _, err := EvaluateTransfer(local("x"), local("y"), false, TransformNone); !errors.Is(err, ErrUnsafe) {
		t.Fatal(err)
	}
	absLocal = originalAbs
	statLocal = func(string) (os.FileInfo, error) { return nil, errors.New("stat") }
	if _, err := EvaluateTransfer(local("x"), local("y"), false, TransformNone); !errors.Is(err, ErrUnsafe) {
		t.Fatal(err)
	}
	statLocal = func(string) (os.FileInfo, error) { return nil, nil }
	evalLocal = func(string) (string, error) { return "", errors.New("eval") }
	if _, err := EvaluateTransfer(local("x"), local("y"), false, TransformNone); !errors.Is(err, ErrUnsafe) {
		t.Fatal(err)
	}
	statLocal, evalLocal = originalStat, originalEval
	localCaseInsensitive = true
	if disposition, err := EvaluateTransfer(local("MixedCase"), local("mixedcase"), false, TransformArchive); disposition != Proceed || !errors.Is(err, ErrUnsafe) {
		t.Fatalf("expected case-insensitive collision, got %v, %v", disposition, err)
	}
	localCaseInsensitive = false
	relLocal = func(string, string) (string, error) { return "", errors.New("rel") }
	if _, err := EvaluateTransfer(local("x"), local("y"), false, TransformNone); !errors.Is(err, ErrUnsafe) {
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
	for _, name := range []string{"", ".", "..", "../secret", "a/../secret", "a//secret", "./secret", "/absolute", `C:\escape`, "C:/escape", "bad\x00name"} {
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
	relArchive = func(string, string) (string, error) { return "..", nil }
	if _, err := SafeArchiveJoin(root, "entry"); !errors.Is(err, ErrUnsafe) {
		t.Fatalf("expected relative escape failure, got %v", err)
	}
}

func TestValidateArchiveSymlink(t *testing.T) {
	for _, test := range []struct {
		name   string
		target string
		safe   bool
	}{
		{"source/link", "file", true},
		{"source/dir/link", "../file", true},
		{"source/link", "../outside", false},
		{"link", "../outside", false},
		{"source/link", "", false},
		{"source/link", "/absolute", false},
		{"source/link", `C:\outside`, false},
		{"source/link", "C:/outside", false},
		{"source/link", "bad\x00target", false},
	} {
		if err := ValidateArchiveSymlink(test.name, test.target); (err == nil) != test.safe {
			t.Errorf("name=%q target=%q safe=%t err=%v", test.name, test.target, test.safe, err)
		}
	}
}

func TestRemoteRootRelationship(t *testing.T) {
	if _, err := EvaluateTransfer(remote("me", "host", "/"), remote("me", "host", "/child"), true, TransformNone); !errors.Is(err, ErrUnsafe) {
		t.Fatalf("expected root descendant rejection, got %v", err)
	}
}
