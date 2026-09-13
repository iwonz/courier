package endpoint

import (
	"errors"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	for _, test := range []struct {
		input string
		want  Endpoint
	}{
		{"./data", Endpoint{Raw: "./data", Path: "./data"}},
		{`C:\Users\me\data`, Endpoint{Raw: `C:\Users\me\data`, Path: `C:\Users\me\data`}},
		{"C:/Users/me/data", Endpoint{Raw: "C:/Users/me/data", Path: "C:/Users/me/data"}},
		{"name:part", Endpoint{Raw: "name:part", Path: "name:part"}},
		{"server:/opt/data 🚚 with spaces", Endpoint{Raw: "server:/opt/data 🚚 with spaces", Host: "server", Path: "/opt/data 🚚 with spaces", Remote: true}},
		{"root@server:/opt/data", Endpoint{Raw: "root@server:/opt/data", User: "root", Host: "server", Path: "/opt/data", Remote: true}},
		{"root@[2001:db8::1]:/opt/data", Endpoint{Raw: "root@[2001:db8::1]:/opt/data", User: "root", Host: "[2001:db8::1]", Path: "/opt/data", Remote: true}},
	} {
		got, err := Parse(test.input)
		if err != nil || !reflect.DeepEqual(got, test.want) {
			t.Errorf("Parse(%q)=%+v,%v want %+v", test.input, got, err, test.want)
		}
	}
}

func TestParseErrors(t *testing.T) {
	for _, input := range []string{"", "host\n:/x", ":/x", "/bad:/x", "@host:/x", "user@@host:/x", "host name:/x", "[broken:/x", "broken]:/x", "[]:/x"} {
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
}

func TestResolveDestination(t *testing.T) {
	source, _ := Parse("host:/data/photos")
	for _, test := range []struct {
		destination string
		isDir       bool
		name        string
		want        string
	}{
		{"./backup/", false, "", "backup/photos"},
		{`C:\backup\`, false, "", `C:\backup\photos`},
		{"other:/backup/", false, "photos.tar.gz", "/backup/photos.tar.gz"},
		{"./exact", false, "", "./exact"},
		{"./existing", true, "", "existing/photos"},
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
}
