// Package sshx implements native SSH connectivity and SFTP-backed filesystems.
package sshx

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	sshconfig "github.com/kevinburke/ssh_config"
)

// Target is the effective connection configuration for one SSH alias.
type Target struct {
	Alias          string
	Host           string
	User           string
	Port           int
	IdentityFiles  []string
	ProxyJump      []Jump
	KnownHostsFile string
}

// Address returns a host:port network address.
func (t Target) Address() string {
	return net.JoinHostPort(strings.Trim(t.Host, "[]"), strconv.Itoa(t.Port))
}

// Jump is one ProxyJump hop before config resolution.
type Jump struct {
	User string
	Host string
	Port int
}

type configValues interface {
	Get(string, string) (string, error)
	GetAll(string, string) ([]string, error)
}

// Config wraps the established OpenSSH config parser with Courier resolution.
type Config struct {
	parsed configValues
	home   string
}

var (
	readConfig   = os.ReadFile
	absConfig    = filepath.Abs
	decodeConfig = sshconfig.DecodeBytes
)

// LoadConfig loads OpenSSH configuration, including Host patterns and Include.
func LoadConfig(name, home string) (*Config, error) {
	absolute, err := absConfig(expandHome(name, home))
	if err != nil {
		return nil, err
	}
	data, err := readConfig(absolute)
	if err != nil {
		return nil, err
	}
	parsed, err := decodeConfig(data)
	if err != nil {
		return nil, fmt.Errorf("parse SSH config %q: %w", absolute, err)
	}
	return &Config{parsed: parsed, home: home}, nil
}

// EmptyConfig returns valid configuration for hosts without a user config file.
func EmptyConfig(home string) *Config {
	parsed, _ := sshconfig.DecodeBytes(nil)
	return &Config{parsed: parsed, home: home}
}

// Resolve applies first-value OpenSSH semantics for required scalar fields.
func (c *Config) Resolve(alias, explicitUser, defaultUser, knownHosts string) (Target, error) {
	if c == nil || c.parsed == nil || alias == "" {
		return Target{}, errors.New("SSH host alias is required")
	}
	get := func(key string) (string, error) { return c.parsed.Get(alias, key) }
	host, err := get("HostName")
	if err != nil {
		return Target{}, err
	}
	if host == "" {
		host = alias
	}
	configuredUser, err := get("User")
	if err != nil {
		return Target{}, err
	}
	user := explicitUser
	if user == "" {
		user = configuredUser
	}
	if user == "" {
		user = defaultUser
	}
	portText, err := get("Port")
	if err != nil {
		return Target{}, err
	}
	port := 22
	if portText != "" {
		parsed, parseErr := strconv.Atoi(portText)
		if parseErr != nil || parsed < 1 || parsed > 65535 {
			return Target{}, fmt.Errorf("invalid SSH port %q", portText)
		}
		port = parsed
	}
	identities, err := c.parsed.GetAll(alias, "IdentityFile")
	if err != nil {
		return Target{}, err
	}
	for index, identity := range identities {
		identities[index] = filepath.Clean(expandTokens(identity, c.home, host, user, port))
	}
	proxyJump, err := get("ProxyJump")
	if err != nil {
		return Target{}, err
	}
	jumps, err := parseJumps(proxyJump)
	if err != nil {
		return Target{}, err
	}
	return Target{Alias: alias, Host: expandTokens(host, c.home, host, user, port), User: user, Port: port, IdentityFiles: identities, ProxyJump: jumps, KnownHostsFile: expandHome(knownHosts, c.home)}, nil
}

func parseJumps(value string) ([]Jump, error) {
	if value == "" || strings.EqualFold(value, "none") {
		return nil, nil
	}
	parts := strings.Split(value, ",")
	result := make([]Jump, 0, len(parts))
	for _, original := range parts {
		item := strings.TrimSpace(original)
		user := ""
		if at := strings.LastIndex(item, "@"); at >= 0 {
			user, item = item[:at], item[at+1:]
		}
		host, portText, err := net.SplitHostPort(item)
		port := 0
		if err == nil {
			port, err = strconv.Atoi(portText)
		} else if strings.Contains(item, ":") && !strings.HasPrefix(item, "[") {
			host, portText, _ = strings.Cut(item, ":")
			port, err = strconv.Atoi(portText)
		} else {
			host, err = item, nil
		}
		if err != nil || strings.Contains(user, "@") || user == "" && strings.Contains(original, "@@") || host == "" || port < 0 || port > 65535 {
			return nil, fmt.Errorf("invalid ProxyJump %q", original)
		}
		result = append(result, Jump{User: user, Host: strings.Trim(host, "[]"), Port: port})
	}
	return result, nil
}

func expandHome(value, home string) string {
	if value == "~" {
		return home
	}
	if strings.HasPrefix(value, "~/") || strings.HasPrefix(value, `~\`) {
		return filepath.Join(home, value[2:])
	}
	return value
}

func expandTokens(value, home, host, user string, port int) string {
	replacer := strings.NewReplacer("%d", home, "%h", host, "%r", user, "%p", strconv.Itoa(port), "%%", "%")
	return expandHome(replacer.Replace(value), home)
}
