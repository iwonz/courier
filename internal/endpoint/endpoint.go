// Package endpoint parses and describes transfer endpoints.
package endpoint

import (
	"errors"
	"fmt"
	"net/netip"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"unicode"
)

// ErrInvalid identifies malformed endpoint syntax.
var ErrInvalid = errors.New("invalid endpoint")

// Kind identifies how Courier reaches an endpoint.
type Kind uint8

const (
	KindLocal Kind = iota
	KindSSH
	KindWeb
	KindWebhook
	KindHTTP
)

// String returns the stable contract name for a kind.
func (k Kind) String() string {
	switch k {
	case KindLocal:
		return "local"
	case KindSSH:
		return "ssh"
	case KindWeb:
		return "web"
	case KindWebhook:
		return "webhook"
	case KindHTTP:
		return "http"
	default:
		return "unknown"
	}
}

// PathFlavor identifies the path rules encoded by an endpoint.
type PathFlavor uint8

const (
	PathNative PathFlavor = iota
	PathPOSIX
	PathWindows
)

// Endpoint is a typed path, Courier service sentinel, or HTTP(S) URL.
// Remote is retained for the existing transfer runtime and is true exactly
// when Kind is KindSSH.
type Endpoint struct {
	Raw        string
	Kind       Kind
	User       string
	Host       string
	Path       string
	Scheme     string
	PathFlavor PathFlavor
	Remote     bool
}

// Parse preserves path text while separating strict endpoint syntax.
func Parse(raw string) (Endpoint, error) {
	if raw == "" {
		return Endpoint{}, fmt.Errorf("%w: empty value", ErrInvalid)
	}
	if hasUnsafeControl(raw) {
		return Endpoint{}, fmt.Errorf("%w: control character", ErrInvalid)
	}
	if raw == "web://" {
		return Endpoint{Raw: raw, Kind: KindWeb, Scheme: "web"}, nil
	}
	if raw == "webhook://" {
		return Endpoint{Raw: raw, Kind: KindWebhook, Scheme: "webhook"}, nil
	}
	if scheme, ok := explicitScheme(raw); ok {
		if scheme != "http" && scheme != "https" {
			return Endpoint{}, fmt.Errorf("%w: unsupported scheme %q", ErrInvalid, scheme)
		}
		parsed, err := url.Parse(raw)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" || parsed.User != nil || parsed.Fragment != "" {
			return Endpoint{}, fmt.Errorf("%w: malformed HTTP URL", ErrInvalid)
		}
		return Endpoint{Raw: raw, Kind: KindHTTP, Host: parsed.Host, Path: parsed.EscapedPath(), Scheme: strings.ToLower(parsed.Scheme)}, nil
	}
	if isWindowsDrive(raw) {
		return Endpoint{Raw: raw, Path: raw, PathFlavor: PathWindows}, nil
	}
	if separator := remoteSeparator(raw); separator >= 0 {
		authority := raw[:separator]
		remotePath := raw[separator+1:]
		user, host, err := parseAuthority(authority)
		if err != nil {
			return Endpoint{}, err
		}
		flavor := PathPOSIX
		if isWindowsDrive(remotePath) {
			flavor = PathWindows
		}
		return Endpoint{Raw: raw, Kind: KindSSH, User: user, Host: host, Path: remotePath, PathFlavor: flavor, Remote: true}, nil
	}
	return Endpoint{Raw: raw, Path: raw}, nil
}

func explicitScheme(raw string) (string, bool) {
	separator := strings.Index(raw, "://")
	if separator <= 0 {
		return "", false
	}
	scheme := raw[:separator]
	for index, character := range scheme {
		if !unicode.IsLetter(character) && (index == 0 || !unicode.IsDigit(character) && character != '+' && character != '-' && character != '.') {
			return "", false
		}
	}
	return strings.ToLower(scheme), true
}

func remoteSeparator(raw string) int {
	for index := 0; index < len(raw); index++ {
		if raw[index] == ':' && remoteAbsolutePath(raw[index+1:]) {
			return index
		}
	}
	return -1
}

func remoteAbsolutePath(value string) bool {
	return strings.HasPrefix(value, "/") || isWindowsDrive(value) && value[2] == '/'
}

func parseAuthority(authority string) (string, string, error) {
	if authority == "" || strings.ContainsAny(authority, `/\`) {
		return "", "", fmt.Errorf("%w: malformed remote authority", ErrInvalid)
	}
	user, host := "", authority
	if at := strings.LastIndex(authority, "@"); at >= 0 {
		user, host = authority[:at], authority[at+1:]
		if user == "" || strings.ContainsAny(user, "@:/\\ \t") {
			return "", "", fmt.Errorf("%w: malformed remote user", ErrInvalid)
		}
	}
	if host == "" || strings.Contains(host, "@") || !validHost(host) {
		return "", "", fmt.Errorf("%w: malformed remote host", ErrInvalid)
	}
	return user, host, nil
}

// Base returns the final source name using slash semantics for SSH and either
// slash style for local input.
func (e Endpoint) Base() string {
	value := e.Path
	if e.Kind != KindSSH && !e.Remote {
		value = strings.ReplaceAll(value, `\`, "/")
	}
	return path.Base(strings.TrimRight(value, "/"))
}

// WithPath returns a copy with a new path and matching raw representation.
func (e Endpoint) WithPath(value string) Endpoint {
	e.Path = value
	if e.Kind == KindSSH || e.Remote {
		authority := e.Host
		if e.User != "" {
			authority = e.User + "@" + authority
		}
		e.Raw = authority + ":" + value
	} else {
		e.Raw = value
	}
	return e
}

// IsPath reports whether the endpoint participates in filesystem transfers.
func (e Endpoint) IsPath() bool {
	return e.Kind == KindLocal || e.Kind == KindSSH || e.Remote
}

// ResolveDestination applies Courier's directory-versus-exact semantics.
// outputName may override the source base name, as archive mode does.
func ResolveDestination(source, destination Endpoint, destinationIsDirectory bool, outputName string) (Endpoint, error) {
	if !source.IsPath() || !destination.IsPath() {
		return Endpoint{}, fmt.Errorf("%w: destination resolution requires path endpoints", ErrInvalid)
	}
	if outputName == "" {
		outputName = source.Base()
	}
	if outputName == "" || outputName == "." || outputName == ".." {
		return Endpoint{}, fmt.Errorf("%w: source has no transferable name", ErrInvalid)
	}
	directoryHint := destinationIsDirectory || strings.HasSuffix(destination.Path, "/") || (destination.Kind != KindSSH && !destination.Remote && strings.HasSuffix(destination.Path, `\`))
	if !directoryHint {
		return destination, nil
	}
	return destination.WithPath(join(destination, outputName)), nil
}

func join(destination Endpoint, name string) string {
	if destination.Kind == KindSSH || destination.Remote {
		return path.Join(destination.Path, name)
	}
	if strings.HasSuffix(destination.Path, `\`) {
		return strings.TrimRight(destination.Path, `\`) + `\` + name
	}
	return filepath.Join(destination.Path, name)
}

func isWindowsDrive(value string) bool {
	return len(value) >= 3 && ((value[0] >= 'a' && value[0] <= 'z') || (value[0] >= 'A' && value[0] <= 'Z')) && value[1] == ':' && (value[2] == '/' || value[2] == '\\')
}

func validHost(host string) bool {
	if strings.HasPrefix(host, "[") || strings.HasSuffix(host, "]") {
		if !strings.HasPrefix(host, "[") || !strings.HasSuffix(host, "]") || len(host) <= 2 {
			return false
		}
		address, err := netip.ParseAddr(host[1 : len(host)-1])
		return err == nil && address.Is6()
	}
	return !strings.ContainsAny(host, "[]: \t")
}

func hasUnsafeControl(value string) bool {
	for _, r := range value {
		if r == 0 || (unicode.IsControl(r) && r != '\t') {
			return true
		}
	}
	return false
}
