// Package endpoint parses and describes transfer endpoints.
package endpoint

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"unicode"
)

// ErrInvalid identifies malformed endpoint syntax.
var ErrInvalid = errors.New("invalid endpoint")

// Endpoint is either a local path or an SSH-backed absolute path.
type Endpoint struct {
	Raw    string
	User   string
	Host   string
	Path   string
	Remote bool
}

// Parse preserves path text while separating strict remote authority syntax.
func Parse(raw string) (Endpoint, error) {
	if raw == "" {
		return Endpoint{}, fmt.Errorf("%w: empty value", ErrInvalid)
	}
	if hasUnsafeControl(raw) {
		return Endpoint{}, fmt.Errorf("%w: control character", ErrInvalid)
	}
	separator := strings.Index(raw, ":/")
	if separator < 0 || isWindowsDrive(raw) {
		return Endpoint{Raw: raw, Path: raw}, nil
	}
	authority := raw[:separator]
	remotePath := raw[separator+1:]
	if authority == "" || strings.ContainsAny(authority, `/\\`) {
		return Endpoint{}, fmt.Errorf("%w: malformed remote authority", ErrInvalid)
	}
	user, host := "", authority
	if at := strings.LastIndex(authority, "@"); at >= 0 {
		user, host = authority[:at], authority[at+1:]
		if user == "" || strings.Contains(user, "@") {
			return Endpoint{}, fmt.Errorf("%w: empty remote user", ErrInvalid)
		}
	}
	if host == "" || strings.Contains(host, "@") || !validHost(host) {
		return Endpoint{}, fmt.Errorf("%w: malformed remote host", ErrInvalid)
	}
	return Endpoint{Raw: raw, User: user, Host: host, Path: remotePath, Remote: true}, nil
}

// Base returns the final source name using slash semantics for remote and
// either slash style for local input.
func (e Endpoint) Base() string {
	value := e.Path
	if !e.Remote {
		value = strings.ReplaceAll(value, `\`, "/")
	}
	return path.Base(strings.TrimRight(value, "/"))
}

// WithPath returns a copy with a new path and matching raw representation.
func (e Endpoint) WithPath(value string) Endpoint {
	e.Path = value
	if e.Remote {
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

// ResolveDestination applies Courier's directory-versus-exact semantics.
// outputName may override the source base name, as archive mode does.
func ResolveDestination(source, destination Endpoint, destinationIsDirectory bool, outputName string) (Endpoint, error) {
	if outputName == "" {
		outputName = source.Base()
	}
	if outputName == "" || outputName == "." || outputName == ".." {
		return Endpoint{}, fmt.Errorf("%w: source has no transferable name", ErrInvalid)
	}
	directoryHint := destinationIsDirectory || strings.HasSuffix(destination.Path, "/") || (!destination.Remote && strings.HasSuffix(destination.Path, `\`))
	if !directoryHint {
		return destination, nil
	}
	return destination.WithPath(join(destination, outputName)), nil
}

func join(destination Endpoint, name string) string {
	if destination.Remote {
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
		return strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") && len(host) > 2
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
