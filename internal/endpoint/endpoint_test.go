package endpoint

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"
)

func TestKindNames(t *testing.T) {
	for kind, want := range map[Kind]string{
		KindLocal:   "local",
		KindSSH:     "ssh",
		KindWeb:     "web",
		KindWebhook: "webhook",
		KindHTTP:    "http",
		Kind(255):   "unknown",
	} {
		if got := kind.String(); got != want {
			t.Errorf("kind %d=%q want %q", kind, got, want)
		}
	}
}

func TestParse(t *testing.T) {
	for _, test := range []struct {
		input string
		want  Endpoint
	}{
		{"./data", Endpoint{Raw: "./data", Path: "./data"}},
		{"web", Endpoint{Raw: "web", Path: "web"}},
		{"webhook", Endpoint{Raw: "webhook", Path: "webhook"}},
		{`C:\Users\me\data`, Endpoint{Raw: `C:\Users\me\data`, Path: `C:\Users\me\data`, PathFlavor: PathWindows}},
		{"c:/Users/me/data", Endpoint{Raw: "c:/Users/me/data", Path: "c:/Users/me/data", PathFlavor: PathWindows}},
		{"name:part", Endpoint{Raw: "name:part", Path: "name:part"}},
		{"1bad://local", Endpoint{Raw: "1bad://local", Kind: KindSSH, Host: "1bad", Path: "//local", PathFlavor: PathPOSIX, Remote: true}},
		{"bad_name://local", Endpoint{Raw: "bad_name://local", Kind: KindSSH, Host: "bad_name", Path: "//local", PathFlavor: PathPOSIX, Remote: true}},
		{"server:/opt/data 🚚 with spaces", Endpoint{Raw: "server:/opt/data 🚚 with spaces", Kind: KindSSH, Host: "server", Path: "/opt/data 🚚 with spaces", PathFlavor: PathPOSIX, Remote: true}},
		{"root@server:/opt/data", Endpoint{Raw: "root@server:/opt/data", Kind: KindSSH, User: "root", Host: "server", Path: "/opt/data", PathFlavor: PathPOSIX, Remote: true}},
		{"root@[2001:db8::1]:/opt/data", Endpoint{Raw: "root@[2001:db8::1]:/opt/data", Kind: KindSSH, User: "root", Host: "[2001:db8::1]", Path: "/opt/data", PathFlavor: PathPOSIX, Remote: true}},
		{"root@[fe80::1%25en0]:/opt/data", Endpoint{Raw: "root@[fe80::1%25en0]:/opt/data", Kind: KindSSH, User: "root", Host: "[fe80::1%25en0]", Path: "/opt/data", PathFlavor: PathPOSIX, Remote: true}},
		{"user@server:C:/Users/me/data", Endpoint{Raw: "user@server:C:/Users/me/data", Kind: KindSSH, User: "user", Host: "server", Path: "C:/Users/me/data", PathFlavor: PathWindows, Remote: true}},
		{"web://", Endpoint{Raw: "web://", Kind: KindWeb, Scheme: "web"}},
		{"webhook://", Endpoint{Raw: "webhook://", Kind: KindWebhook, Scheme: "webhook"}},
		{"https://example.test/hooks/a%20b?token=value", Endpoint{Raw: "https://example.test/hooks/a%20b?token=value", Kind: KindHTTP, Host: "example.test", Path: "/hooks/a%20b", Scheme: "https"}},
		{"HTTP://example.test", Endpoint{Raw: "HTTP://example.test", Kind: KindHTTP, Host: "example.test", Scheme: "http"}},
	} {
		got, err := Parse(test.input)
		if err != nil || !reflect.DeepEqual(got, test.want) {
			t.Errorf("Parse(%q)=%+v,%v want %+v", test.input, got, err, test.want)
		}
	}
}

func TestParseErrors(t *testing.T) {
	for _, input := range []string{
		"", "host\n:/x", ":/x", "/bad:/x", "@host:/x", "user@@host:/x", "bad user@host:/x", "user:bad@host:/x",
		"host name:/x", "host:8080:/x", "[broken:/x", "broken]:/x", "[]:/x", "[127.0.0.1]:/x", "[not-ip]:/x",
		"ftp://example.test/file", "web://host", "webhook://path", "http://", "http://[::1", "https://user@example.test/hook", "https://example.test/hook#secret",
	} {
		if _, err := Parse(input); !errors.Is(err, ErrInvalid) {
			t.Errorf("Parse(%q) error=%v", input, err)
		}
	}
}

func TestBaseAndWithPath(t *testing.T) {
	local, _ := Parse(`C:\data\photos\`)
	if got := local.Base(); got != "photos" {
		t.Fatalf("base=%q", got)
	}
	if got := local.WithPath("next"); got.Raw != "next" || got.Path != "next" {
		t.Fatalf("local with path=%+v", got)
	}
	remote, _ := Parse("me@host:/data/photos/")
	if remote.Base() != "photos" || remote.WithPath("/next").Raw != "me@host:/next" {
		t.Fatalf("remote helpers failed: %+v", remote)
	}
	legacyRemote := Endpoint{Remote: true, Host: "host", Path: "/old"}
	if !legacyRemote.IsPath() || legacyRemote.WithPath("/new").Raw != "host:/new" {
		t.Fatalf("legacy remote helpers failed: %+v", legacyRemote)
	}
	if (Endpoint{Kind: KindWeb}).IsPath() {
		t.Fatal("web endpoint reported as a path")
	}
}

func TestResolveDestination(t *testing.T) {
	source, _ := Parse("host:/data/photos")
	for _, test := range []struct {
		destination string
		isDir       bool
		name        string
		want        string
	}{
		{"./backup/", false, "", filepath.Join("backup", "photos")},
		{`C:\backup\`, false, "", `C:\backup\photos`},
		{"other:/backup/", false, "photos.tar.gz", "/backup/photos.tar.gz"},
		{"windows:C:/backup/", false, "", "C:/backup/photos"},
		{"./exact", false, "", "./exact"},
		{"./existing", true, "", filepath.Join("existing", "photos")},
	} {
		destination, _ := Parse(test.destination)
		got, err := ResolveDestination(source, destination, test.isDir, test.name)
		if err != nil || got.Path != test.want {
			t.Errorf("destination=%q got=%q err=%v want=%q", test.destination, got.Path, err, test.want)
		}
	}
	empty := Endpoint{Path: "/"}
	destination, _ := Parse("./backup/")
	if _, err := ResolveDestination(empty, destination, false, ""); !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid source name, got %v", err)
	}
	if _, err := ResolveDestination(Endpoint{Kind: KindWeb}, destination, false, "x"); !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected non-path source error, got %v", err)
	}
	if _, err := ResolveDestination(source, Endpoint{Kind: KindWeb}, false, "x"); !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected non-path destination error, got %v", err)
	}
}
