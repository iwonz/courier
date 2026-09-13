// Package transfer implements transport-neutral transactional tree copying.
package transfer

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"time"

	"github.com/iwonz/courier/internal/fsx"
	"github.com/iwonz/courier/internal/progress"
	"github.com/iwonz/courier/internal/selection"
)

const defaultBufferSize = 128 * 1024

var randomRead = rand.Read

// Request describes a resolved source and exact destination.
type Request struct {
	SourceFS      fsx.Backend
	SourcePath    string
	DestinationFS fsx.Backend
	Destination   string
	Selector      selection.Selector
	Progress      progress.Sink
}

// Result summarizes a committed transfer.
type Result struct {
	Bytes       int64
	Destination string
	Elapsed     time.Duration
}

// Error records the failed stage and confirmed byte count.
type Error struct {
	Stage     progress.Stage
	Read      int64
	Sent      int64
	Confirmed int64
	Cause     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s failed after read=%d sent=%d confirmed=%d bytes: %v", e.Stage, e.Read, e.Sent, e.Confirmed, e.Cause)
}

func (e *Error) Unwrap() error { return e.Cause }

// Engine controls deterministic dependencies used by the transfer algorithm.
type Engine struct {
	Token      func() (string, error)
	Now        func() time.Time
	BufferSize int
}

// Run scans, stages, and commits a transfer.
func (e Engine) Run(ctx context.Context, request Request) (result Result, resultErr error) {
	if request.SourceFS == nil || request.DestinationFS == nil || request.SourcePath == "" || request.Destination == "" {
		return Result{}, &Error{Stage: progress.StagePreflight, Cause: errors.New("source and destination are required")}
	}
	if err := ensureAbsent(request.DestinationFS, request.Destination); err != nil {
		return Result{}, &Error{Stage: progress.StagePreflight, Cause: err}
	}
	selector := request.Selector
	if selector == nil {
		all := selection.All()
		selector = all
	}
	total, err := scanSelected(ctx, request.SourceFS, request.SourcePath, "", selector)
	if err != nil {
		return Result{}, &Error{Stage: progress.StagePreflight, Cause: err}
	}
	tracker := progress.New(total, e.Now, request.Progress)
	tokenFn := e.Token
	if tokenFn == nil {
		tokenFn = randomToken
	}
	token, err := tokenFn()
	if err != nil {
		return Result{}, &Error{Stage: progress.StagePreflight, Cause: err}
	}
	stagePath := request.Destination + ".courier-partial-" + token
	cleanupStage := true
	defer func() {
		tracker.Stage(progress.StageCleanup)
		if cleanupStage {
			if cleanupErr := request.DestinationFS.RemoveAll(stagePath); cleanupErr != nil {
				snapshot := tracker.Snapshot()
				cleanupFailure := &Error{Stage: progress.StageCleanup, Read: snapshot.Read, Sent: snapshot.Sent, Confirmed: snapshot.Confirmed, Cause: cleanupErr}
				resultErr = errors.Join(resultErr, cleanupFailure)
			}
		}
		if resultErr == nil {
			tracker.Stage(progress.StageComplete)
			result = Result{Bytes: tracker.Snapshot().Current, Destination: request.Destination, Elapsed: tracker.Snapshot().Elapsed}
		}
	}()
	if err := request.DestinationFS.RemoveAll(stagePath); err != nil {
		return Result{}, transferError(tracker, progress.StageCleanup, err)
	}
	if err := request.DestinationFS.MkdirAll(request.DestinationFS.Dir(request.Destination), 0o700); err != nil {
		return Result{}, transferError(tracker, progress.StagePreflight, err)
	}
	tracker.Stage(progress.StageTransfer)
	bufferSize := e.BufferSize
	if bufferSize <= 0 {
		bufferSize = defaultBufferSize
	}
	if err := copySelected(ctx, request.SourceFS, request.SourcePath, request.DestinationFS, stagePath, "", make([]byte, bufferSize), tracker, selector); err != nil {
		return Result{}, transferError(tracker, progress.StageTransfer, err)
	}
	tracker.Stage(progress.StageCommit)
	if err := commit(request.DestinationFS, stagePath, request.Destination); err != nil {
		return Result{}, transferError(tracker, progress.StageCommit, err)
	}
	cleanupStage = false
	return result, nil
}

func transferError(tracker *progress.Tracker, stage progress.Stage, err error) *Error {
	snapshot := tracker.Snapshot()
	return &Error{Stage: stage, Read: snapshot.Read, Sent: snapshot.Sent, Confirmed: snapshot.Confirmed, Cause: err}
}

func scan(ctx context.Context, backend fsx.Backend, name string) (int64, error) {
	return scanSelected(ctx, backend, name, "", selection.All())
}

func scanSelected(ctx context.Context, backend fsx.Backend, name, relative string, selector selection.Selector) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	info, err := backend.Lstat(name)
	if err != nil {
		return 0, err
	}
	if relative != "" && !selector.Include(relative, info.IsDir()) {
		return 0, nil
	}
	if info.Mode().IsRegular() {
		return info.Size(), nil
	}
	if info.Mode()&fs.ModeSymlink != 0 {
		return 0, nil
	}
	if !info.IsDir() {
		return 0, fmt.Errorf("unsupported source object %q (%s)", name, info.Mode())
	}
	entries, err := backend.ReadDir(name)
	if err != nil {
		return 0, err
	}
	var total int64
	for _, entry := range entries {
		size, err := scanSelected(ctx, backend, backend.Join(name, entry.Name()), path.Join(relative, entry.Name()), selector)
		if err != nil {
			return 0, err
		}
		total += size
	}
	return total, nil
}

func copyNode(ctx context.Context, source fsx.Backend, sourcePath string, destination fsx.Backend, destinationPath string, buffer []byte, tracker *progress.Tracker) error {
	return copySelected(ctx, source, sourcePath, destination, destinationPath, "", buffer, tracker, selection.All())
}

func copySelected(ctx context.Context, source fsx.Backend, sourcePath string, destination fsx.Backend, destinationPath, relative string, buffer []byte, tracker *progress.Tracker, selector selection.Selector) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	info, err := source.Lstat(sourcePath)
	if err != nil {
		return err
	}
	if relative != "" && !selector.Include(relative, info.IsDir()) {
		return nil
	}
	if info.Mode().IsRegular() {
		return copyFile(ctx, source, sourcePath, destination, destinationPath, info, buffer, tracker)
	}
	if info.Mode()&fs.ModeSymlink != 0 {
		target, err := source.Readlink(sourcePath)
		if err != nil {
			return err
		}
		return destination.Symlink(target, destinationPath)
	}
	if !info.IsDir() {
		return fmt.Errorf("unsupported source object %q (%s)", sourcePath, info.Mode())
	}
	if err := destination.MkdirAll(destinationPath, 0o700); err != nil {
		return err
	}
	entries, err := source.ReadDir(sourcePath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if err := copySelected(ctx, source, source.Join(sourcePath, entry.Name()), destination, destination.Join(destinationPath, entry.Name()), path.Join(relative, entry.Name()), buffer, tracker, selector); err != nil {
			return err
		}
	}
	if err := destination.Chmod(destinationPath, info.Mode().Perm()); err != nil {
		return err
	}
	return destination.Chtimes(destinationPath, info.ModTime(), info.ModTime())
}

func copyFile(ctx context.Context, source fsx.Backend, sourcePath string, destination fsx.Backend, destinationPath string, info fs.FileInfo, buffer []byte, tracker *progress.Tracker) (resultErr error) {
	reader, err := source.Open(sourcePath)
	if err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, reader.Close()) }()
	writer, err := destination.Create(destinationPath, 0o600)
	if err != nil {
		return err
	}
	closed := false
	defer func() {
		if !closed {
			resultErr = errors.Join(resultErr, writer.Close())
		}
	}()
	progressWriter := &countingWriter{ctx: ctx, destination: writer, tracker: tracker}
	if _, err := io.CopyBuffer(progressWriter, reader, buffer); err != nil {
		return err
	}
	if err := writer.Sync(); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	closed = true
	if err := destination.Chmod(destinationPath, info.Mode().Perm()); err != nil {
		return err
	}
	return destination.Chtimes(destinationPath, info.ModTime(), info.ModTime())
}

type countingWriter struct {
	ctx         context.Context
	destination io.Writer
	tracker     *progress.Tracker
}

func (w *countingWriter) Write(data []byte) (int, error) {
	if err := w.ctx.Err(); err != nil {
		return 0, err
	}
	written, err := w.destination.Write(data)
	if written > 0 {
		err = errors.Join(err, w.tracker.Add(int64(written)))
	}
	return written, err
}

func ensureAbsent(backend fsx.Backend, destination string) error {
	if _, err := backend.Lstat(destination); err == nil {
		return fmt.Errorf("%w: %q", fsx.ErrDestinationExists, destination)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}

func commit(backend fsx.Backend, stagePath, destination string) error {
	return backend.CommitAbsent(stagePath, destination)
}

func randomToken() (string, error) {
	value := make([]byte, 8)
	if _, err := randomRead(value); err != nil {
		return "", err
	}
	return hex.EncodeToString(value), nil
}
