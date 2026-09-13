// Package safety implements transfer preflight path checks.
package safety

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/iwonz/courier/internal/endpoint"
)

// ErrUnsafe identifies a path relationship that could overwrite source data.
var ErrUnsafe = errors.New("unsafe transfer path")

// Transformation identifies whether transfer output differs from the source.
type Transformation uint8

const (
	TransformNone Transformation = iota
	TransformArchive
	TransformExtract
)

// Disposition tells orchestration whether a safe transfer should run.
type Disposition uint8

const (
	Proceed Disposition = iota
	NoOp
)

var (
	absLocal             = filepath.Abs
	evalLocal            = filepath.EvalSymlinks
	statLocal            = os.Lstat
	relLocal             = filepath.Rel
	localCaseInsensitive = runtime.GOOS == "windows"
	relArchive           = filepath.Rel
)

// EvaluateTransfer returns a successful no-op for plain endpoint identity and
// rejects transformed identity and directory descendants.
func EvaluateTransfer(source, destination endpoint.Endpoint, sourceIsDirectory bool, transformation Transformation) (Disposition, error) {
	same, descendant, err := relationship(source, destination)
	if err != nil {
		return Proceed, fmt.Errorf("%w: canonicalize paths: %v", ErrUnsafe, err)
	}
	if same {
		if transformation == TransformNone {
			return NoOp, nil
		}
		return Proceed, fmt.Errorf("%w: transformed output collides with source", ErrUnsafe)
	}
	if sourceIsDirectory && descendant {
		return Proceed, fmt.Errorf("%w: destination is inside source", ErrUnsafe)
	}
	return Proceed, nil
}

func relationship(source, destination endpoint.Endpoint) (bool, bool, error) {
	if source.Remote != destination.Remote {
		return false, false, nil
	}
	if source.Remote {
		if !sameRemote(source, destination) {
			return false, false, nil
		}
		sourcePath, destinationPath := path.Clean(source.Path), path.Clean(destination.Path)
		same := sourcePath == destinationPath
		if source.PathFlavor == endpoint.PathWindows || destination.PathFlavor == endpoint.PathWindows {
			same = strings.EqualFold(sourcePath, destinationPath)
			destinationPath = strings.ToLower(destinationPath)
			sourcePath = strings.ToLower(sourcePath)
		}
		prefix := sourcePath + "/"
		if sourcePath == "/" {
			prefix = "/"
		}
		return same, !same && strings.HasPrefix(destinationPath, prefix), nil
	}
	sourcePath, err := canonicalLocal(source.Path)
	if err != nil {
		return false, false, err
	}
	destinationPath, err := canonicalLocal(destination.Path)
	if err != nil {
		return false, false, err
	}
	if localCaseInsensitive {
		sourcePath, destinationPath = strings.ToLower(sourcePath), strings.ToLower(destinationPath)
	}
	relative, err := relLocal(sourcePath, destinationPath)
	return relative == ".", isDescendant(relative), err
}

func sameRemote(left, right endpoint.Endpoint) bool {
	return strings.EqualFold(left.Host, right.Host) && left.User == right.User
}

func isDescendant(relative string) bool {
	return relative != "." && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) && !strings.HasPrefix(relative, "../")
}

func canonicalLocal(value string) (string, error) {
	absolute, err := absLocal(value)
	if err != nil {
		return "", err
	}
	absolute = filepath.Clean(absolute)
	current, suffix := absolute, ""
	for {
		if _, err := statLocal(current); err == nil {
			resolved, err := evalLocal(current)
			if err != nil {
				return "", err
			}
			if suffix != "" {
				resolved = filepath.Join(resolved, suffix)
			}
			return filepath.Clean(resolved), nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
		parent, base := filepath.Dir(current), filepath.Base(current)
		if parent == current {
			return absolute, nil
		}
		if suffix == "" {
			suffix = base
		} else {
			suffix = filepath.Join(base, suffix)
		}
		current = parent
	}
}

// SafeArchiveJoin resolves a slash-separated archive entry below root.
func SafeArchiveJoin(root, name string) (string, error) {
	if name == "" || strings.ContainsRune(name, 0) || path.IsAbs(name) || filepath.IsAbs(name) || strings.Contains(name, `\`) {
		return "", fmt.Errorf("%w: invalid archive entry %q", ErrUnsafe, name)
	}
	clean := path.Clean(name)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
		return "", fmt.Errorf("%w: archive entry escapes root", ErrUnsafe)
	}
	joined := filepath.Join(root, filepath.FromSlash(clean))
	relative, err := relArchive(root, joined)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("%w: archive entry escapes root", ErrUnsafe)
	}
	return joined, nil
}
