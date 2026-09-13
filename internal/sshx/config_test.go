package sshx

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	sshconfig "github.com/kevinburke/ssh_config"
)

func TestLoadAndResolveConfig(t *testing.T) {
	root := t.TempDir()
	included := filepath.Join(root, "included.conf")
	if err := os.WriteFile(included, []byte("Host included\n  User include-user\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(root, "config")
	content := "Include " + included + `
Host target !blocked
  HostName "127.0.0.1"
  User config-user
  Port=2222
  IdentityFile ~/.ssh/id_one
  IdentityFile %d/.ssh/id_%h_%r_%p
  ProxyJump jump@bastion:2200,[::1]:2201
Host *
  User fallback
`
	if err := os.WriteFile(configPath, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	config, err := LoadConfig(configPath, root)
	if err != nil {
		t.Fatal(err)
	}
	target, err := config.Resolve("target", "explicit", "local", "~/.ssh/known_hosts")
	if err != nil {
		t.Fatal(err)
	}
	wantIdentities := []string{filepath.Join(root, ".ssh", "id_one"), filepath.Join(root, ".ssh", "id_127.0.0.1_explicit_2222")}
	if target.Host != "127.0.0.1" || target.User != "explicit" || target.Port != 2222 || target.Address() != "127.0.0.1:2222" || !reflect.DeepEqual(target.IdentityFiles, wantIdentities) || target.KnownHostsFile != filepath.Join(root, ".ssh", "known_hosts") {
		t.Fatalf("target=%+v", target)
	}
	if want := []Jump{{User: "jump", Host: "bastion", Port: 2200}, {Host: "::1", Port: 2201}}; !reflect.DeepEqual(target.ProxyJump, want) {
		t.Fatalf("jumps=%+v", target.ProxyJump)
	}
	includedTarget, err := config.Resolve("included", "", "local", "known")
	if err != nil || includedTarget.User != "include-user" || includedTarget.Host != "included" || includedTarget.Port != 22 {
		t.Fatalf("included=%+v err=%v", includedTarget, err)
	}
	empty := EmptyConfig(root)
	defaultTarget, err := empty.Resolve("plain", "", "local", "known")
	if err != nil || defaultTarget.Host != "plain" || defaultTarget.User != "local" || defaultTarget.Port != 22 {
		t.Fatalf("default=%+v err=%v", defaultTarget, err)
	}
}

func TestConfigErrors(t *testing.T) {
	root := t.TempDir()
	if _, err := LoadConfig(filepath.Join(root, "missing"), root); err == nil {
		t.Fatal("expected missing config")
	}
	if _, err := (*Config)(nil).Resolve("host", "", "", "known"); err == nil {
		t.Fatal("expected nil config")
	}
	config := configFromText(t, root, "Host host\n Port bad\n")
	if _, err := config.Resolve("", "", "", "known"); err == nil {
		t.Fatal("expected empty alias")
	}
	if _, err := config.Resolve("host", "", "", "known"); err == nil {
		t.Fatal("expected invalid port")
	}
	config = configFromText(t, root, "Host host\n Port 70000\n")
	if _, err := config.Resolve("host", "", "", "known"); err == nil {
		t.Fatal("expected port range error")
	}
	config = configFromText(t, root, "Host host\n ProxyJump host:not-a-port\n")
	if _, err := config.Resolve("host", "", "default", "known"); err == nil {
		t.Fatal("expected invalid jump")
	}
}

func TestConfigHelpers(t *testing.T) {
	for _, value := range []string{"", "none", "NONE"} {
		if jumps, err := parseJumps(value); err != nil || jumps != nil {
			t.Fatalf("parseJumps(%q)=%v,%v", value, jumps, err)
		}
	}
	if jumps, err := parseJumps("plain"); err != nil || !reflect.DeepEqual(jumps, []Jump{{Host: "plain"}}) {
		t.Fatalf("plain jump=%v,%v", jumps, err)
	}
	for _, value := range []string{":0", "bad@@host", "host:22:extra"} {
		if _, err := parseJumps(value); err == nil {
			t.Fatalf("expected invalid jump %q", value)
		}
	}
	if expandHome("~", "/home/me") != "/home/me" || expandHome(`~\file`, "/home/me") != filepath.Join("/home/me", "file") || expandHome("other", "/home/me") != "other" {
		t.Fatal("home expansion failed")
	}
	if got := expandTokens("%%-%h-%r-%p-%d", "/home", "host", "user", 22); got != "%-host-user-22-/home" {
		t.Fatalf("tokens=%q", got)
	}
}

func TestConfigInjectedErrors(t *testing.T) {
	originalRead, originalAbs, originalDecode := readConfig, absConfig, decodeConfig
	t.Cleanup(func() { readConfig, absConfig, decodeConfig = originalRead, originalAbs, originalDecode })
	absConfig = func(string) (string, error) { return "", errors.New("absolute") }
	if _, err := LoadConfig("config", t.TempDir()); err == nil {
		t.Fatal("expected absolute path error")
	}
	absConfig = originalAbs
	readConfig = func(string) ([]byte, error) { return []byte("Host host"), nil }
	decodeConfig = func([]byte) (*sshconfig.Config, error) { return nil, errors.New("decode") }
	if _, err := LoadConfig("config", t.TempDir()); err == nil {
		t.Fatal("expected decode error")
	}
	for _, key := range []string{"hostname", "user", "port", "identityfile", "proxyjump"} {
		config := &Config{parsed: failingConfigValues{key: key}, home: t.TempDir()}
		if _, err := config.Resolve("host", "", "default", "known"); err == nil {
			t.Fatalf("expected %s lookup error", key)
		}
	}
}

func configFromText(t *testing.T, root, content string) *Config {
	t.Helper()
	name := filepath.Join(root, "config-test")
	if err := os.WriteFile(name, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	config, err := LoadConfig(name, root)
	if err != nil {
		t.Fatal(err)
	}
	return config
}

type failingConfigValues struct{ key string }

func (f failingConfigValues) Get(_ string, key string) (string, error) {
	if strings.EqualFold(key, f.key) {
		return "", errors.New("get")
	}
	return "", nil
}

func (f failingConfigValues) GetAll(_ string, key string) ([]string, error) {
	if strings.EqualFold(key, f.key) {
		return nil, errors.New("get all")
	}
	return nil, nil
}
