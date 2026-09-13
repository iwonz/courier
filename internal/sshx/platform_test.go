package sshx

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type probeResponse struct {
	data []byte
	err  error
}

type probeRunner struct {
	responses map[string]probeResponse
	commands  []string
}

func (p *probeRunner) Run(_ context.Context, command string) ([]byte, error) {
	p.commands = append(p.commands, command)
	response := p.responses[command]
	return response.data, response.err
}

func TestDetectPlatform(t *testing.T) {
	posix := &probeRunner{responses: map[string]probeResponse{
		posixPlatformProbe: {data: []byte("Darwin\narm64\n")},
		posixArchiverProbe: {data: []byte("/usr/bin/tar\n/usr/local/bin/bsdtar\n")},
	}}
	platform, archiver, err := DetectPlatform(context.Background(), posix)
	if err != nil || platform != (Platform{OS: "darwin", Arch: "arm64"}) || archiver.Name != "/usr/local/bin/bsdtar" || archiver.BuiltIn {
		t.Fatalf("platform=%+v archiver=%+v err=%v", platform, archiver, err)
	}
	if !reflect.DeepEqual(posix.commands, []string{posixPlatformProbe, posixArchiverProbe}) {
		t.Fatalf("commands=%v", posix.commands)
	}
	windows := &probeRunner{responses: map[string]probeResponse{
		posixPlatformProbe:   {err: errors.New("no uname")},
		windowsPlatformProbe: {data: []byte("Win32NT\nAMD64\n")},
		windowsArchiverProbe: {data: []byte(`C:\Windows\System32\tar.exe`)},
	}}
	platform, archiver, err = DetectPlatform(context.Background(), windows)
	if err != nil || platform != (Platform{OS: "windows", Arch: "amd64"}) || archiver.BuiltIn {
		t.Fatalf("platform=%+v archiver=%+v err=%v", platform, archiver, err)
	}
}

func TestDetectPlatformErrors(t *testing.T) {
	if _, _, err := DetectPlatform(context.Background(), nil); err == nil {
		t.Fatal("expected nil runner error")
	}
	runner := &probeRunner{responses: map[string]probeResponse{
		posixPlatformProbe:   {data: []byte("invalid"), err: errors.New("posix")},
		windowsPlatformProbe: {err: errors.New("windows")},
	}}
	if _, _, err := DetectPlatform(context.Background(), runner); err == nil {
		t.Fatal("expected both probes to fail")
	}
	runner.responses[windowsPlatformProbe] = probeResponse{data: []byte("unknown")}
	if _, _, err := DetectPlatform(context.Background(), runner); err == nil {
		t.Fatal("expected invalid Windows output")
	}
}

func TestParsePlatform(t *testing.T) {
	for _, test := range []struct {
		input string
		want  Platform
	}{
		{"Linux x86_64", Platform{OS: "linux", Arch: "amd64"}},
		{"FreeBSD aarch64", Platform{OS: "freebsd", Arch: "arm64"}},
		{"OpenBSD x64", Platform{OS: "openbsd", Arch: "amd64"}},
		{"NetBSD arm64", Platform{OS: "netbsd", Arch: "arm64"}},
		{"DragonFly amd64", Platform{OS: "dragonfly", Arch: "amd64"}},
		{"Windows_NT AMD64", Platform{OS: "windows", Arch: "amd64"}},
	} {
		got, err := parsePlatform(test.input)
		if err != nil || got != test.want {
			t.Errorf("parsePlatform(%q)=%+v,%v", test.input, got, err)
		}
	}
	for _, input := range []string{"", "Linux", "Solaris amd64", "Linux mips"} {
		if _, err := parsePlatform(input); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
}

func TestSelectArchiver(t *testing.T) {
	for _, test := range []struct {
		platform  Platform
		available []string
		want      Archiver
	}{
		{Platform{OS: "linux"}, []string{"/bin/gtar", "/bin/tar"}, Archiver{Name: "/bin/tar"}},
		{Platform{OS: "freebsd"}, []string{"/bin/tar", "/bin/bsdtar"}, Archiver{Name: "/bin/bsdtar"}},
		{Platform{OS: "dragonfly"}, []string{"/bin/gtar"}, Archiver{Name: "/bin/gtar"}},
		{Platform{OS: "windows"}, []string{`C:\Windows\bsdtar.exe`}, Archiver{Name: `C:\Windows\bsdtar.exe`}},
		{Platform{OS: "linux"}, nil, Archiver{Name: "builtin", BuiltIn: true}},
	} {
		if got := SelectArchiver(test.platform, test.available); got != test.want {
			t.Errorf("SelectArchiver(%+v,%v)=%+v", test.platform, test.available, got)
		}
	}
	if lastField("") != "" || lastField("one two") != "two" {
		t.Fatal("lastField failed")
	}
}
