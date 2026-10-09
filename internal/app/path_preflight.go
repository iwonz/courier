package app

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/iwonz/courier/internal/endpoint"
)

var errDirectoryConfirmationUnavailable = errors.New("interactive terminal is unavailable")

// DirectoryConfirm authorizes creation of one missing directory.
type DirectoryConfirm func(context.Context, string) (bool, error)

type pathExpectation uint8

const (
	pathFlexibleSource pathExpectation = iota
	pathDirectory
	pathRegularFile
	pathRequiredFile
	pathExactDestination
)

type pathEntry struct {
	role        string
	display     string
	resource    *Resource
	expectation pathExpectation
	info        fs.FileInfo
	missing     bool
}

type pathPreflightError struct {
	message string
}

func (failure *pathPreflightError) Error() string { return failure.message }

func inspectPath(entry *pathEntry) error {
	if entry == nil || entry.resource == nil || entry.resource.Backend == nil {
		return errors.New("path preflight resource is unavailable")
	}
	info, err := entry.resource.Backend.Lstat(entry.resource.Path)
	if err == nil {
		entry.info = info
		switch entry.expectation {
		case pathDirectory:
			if !info.IsDir() {
				return pathTypeError(entry, "is not a directory")
			}
		case pathRegularFile:
			if !info.Mode().IsRegular() {
				return pathTypeError(entry, "is not a regular file")
			}
		}
		return nil
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	entry.missing = true
	if entry.expectation == pathRegularFile || entry.expectation == pathRequiredFile {
		return pathTypeError(entry, "does not exist")
	}
	return nil
}

func authorizeAndCreatePaths(ctx context.Context, confirm DirectoryConfirm, force bool, entries ...*pathEntry) error {
	missing := make([]*pathEntry, 0, len(entries))
	for _, entry := range entries {
		if entry == nil || !entry.missing || entry.expectation == pathExactDestination {
			continue
		}
		missing = append(missing, entry)
	}
	if !force {
		for _, entry := range missing {
			if confirm == nil {
				return missingDirectoryAuthorizationError(entry)
			}
			accepted, err := confirm(ctx, fmt.Sprintf("%s directory %s does not exist. Create it?", entry.role, entry.display))
			if errors.Is(err, errDirectoryConfirmationUnavailable) {
				return missingDirectoryAuthorizationError(entry)
			}
			if err != nil {
				return err
			}
			if !accepted {
				return &pathPreflightError{message: fmt.Sprintf("%s directory %q was not created", entry.role, entry.display)}
			}
		}
	}
	for _, entry := range missing {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := entry.resource.Backend.MkdirAll(entry.resource.Path, 0o700); err != nil {
			info, statErr := entry.resource.Backend.Lstat(entry.resource.Path)
			if statErr == nil {
				if !info.IsDir() {
					return pathTypeError(entry, "is not a directory")
				}
				entry.info = info
				entry.missing = false
				continue
			}
			return fmt.Errorf("create %s directory %q: %w", strings.ToLower(entry.role), entry.display, err)
		}
		info, err := entry.resource.Backend.Lstat(entry.resource.Path)
		if err != nil {
			return fmt.Errorf("revalidate %s directory %q: %w", strings.ToLower(entry.role), entry.display, err)
		}
		if !info.IsDir() {
			return pathTypeError(entry, "is not a directory")
		}
		entry.info = info
		entry.missing = false
	}
	return nil
}

func missingDirectoryAuthorizationError(entry *pathEntry) error {
	return &pathPreflightError{message: fmt.Sprintf("%s directory %q does not exist; rerun with --force-source-creation to create it", entry.role, entry.display)}
}

func pathTypeError(entry *pathEntry, problem string) error {
	kind := "path"
	if entry.expectation == pathDirectory || entry.expectation == pathFlexibleSource && entry.missing {
		kind = "directory"
	} else if entry.expectation == pathRegularFile || entry.expectation == pathRequiredFile {
		kind = "file"
	}
	return &pathPreflightError{message: fmt.Sprintf("%s %s %q %s", entry.role, kind, entry.display, problem)}
}

func preflightTransferPaths(ctx context.Context, confirm DirectoryConfirm, planForce bool, source, destination *Resource, archive, extract bool) (fs.FileInfo, fs.FileInfo, error) {
	sourceExpectation := pathFlexibleSource
	if extract {
		sourceExpectation = pathRegularFile
	}
	sourceEntry := &pathEntry{role: "Source", display: source.Endpoint.Raw, resource: source, expectation: sourceExpectation}
	if err := inspectPath(sourceEntry); err != nil {
		return nil, nil, err
	}
	destinationExpectation := pathExactDestination
	if extract || !archive && (sourceEntry.missing || sourceEntry.info != nil && sourceEntry.info.IsDir()) || hasDirectoryHint(destination.Endpoint) {
		destinationExpectation = pathDirectory
	}
	destinationEntry := &pathEntry{role: "Destination", display: destination.Endpoint.Raw, resource: destination, expectation: destinationExpectation}
	if err := inspectPath(destinationEntry); err != nil {
		return nil, nil, err
	}
	if err := authorizeAndCreatePaths(ctx, confirm, planForce, sourceEntry, destinationEntry); err != nil {
		return nil, nil, err
	}
	return sourceEntry.info, destinationEntry.info, nil
}

func preflightSinglePath(ctx context.Context, confirm DirectoryConfirm, force bool, resource *Resource, role string, expectation pathExpectation) (fs.FileInfo, error) {
	entry := &pathEntry{role: role, display: resource.Endpoint.Raw, resource: resource, expectation: expectation}
	if err := inspectPath(entry); err != nil {
		return nil, err
	}
	if err := authorizeAndCreatePaths(ctx, confirm, force, entry); err != nil {
		return nil, err
	}
	return entry.info, nil
}

func hasDirectoryHint(value endpoint.Endpoint) bool {
	return strings.HasSuffix(value.Path, "/") || strings.HasSuffix(value.Path, `\`)
}

func isPathPreflightError(err error) bool {
	var failure *pathPreflightError
	return errors.As(err, &failure)
}
