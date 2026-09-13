// Package archive creates and verifies portable built-in tar.gz artifacts.
package archive

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"sync"

	"github.com/iwonz/courier/internal/fsx"
	"github.com/iwonz/courier/internal/progress"
	"github.com/iwonz/courier/internal/safety"
)

type temporaryFile interface {
	io.Writer
	Name() string
	Chmod(fs.FileMode) error
	Sync() error
	Close() error
}

var (
	createTemporary = func(directory, pattern string) (temporaryFile, error) { return os.CreateTemp(directory, pattern) }
	removeTemporary = os.Remove
	statTemporary   = os.Stat
	verifyTemporary = Verify
	closeTar        = func(writer *tar.Writer) error { return writer.Close() }
	closeGzip       = func(writer *gzip.Writer) error { return writer.Close() }
)

// Artifact is a verified private temporary archive.
type Artifact struct {
	Path  string
	Name  string
	Bytes int64
	once  sync.Once
	err   error
}

// Cleanup removes the local temporary archive and is safe to repeat.
func (a *Artifact) Cleanup() error {
	a.once.Do(func() { a.err = os.Remove(a.Path) })
	return a.err
}

// Create writes source and its root entry into a verified tar.gz file.
func Create(ctx context.Context, backend fsx.Backend, sourcePath, sourceName, tempDirectory string, sink progress.Sink) (*Artifact, error) {
	if backend == nil || sourcePath == "" || sourceName == "" || path.Base(sourceName) != sourceName {
		return nil, errors.New("archive source and base name are required")
	}
	if _, err := safety.SafeArchiveJoin(tempDirectory, sourceName); err != nil {
		return nil, err
	}
	file, err := createTemporary(tempDirectory, ".courier-*.tar.gz")
	if err != nil {
		return nil, err
	}
	temporaryPath := file.Name()
	success := false
	defer func() {
		if !success {
			_ = file.Close()
			_ = removeTemporary(temporaryPath)
		}
	}()
	if err := file.Chmod(0o600); err != nil {
		return nil, err
	}
	tracker := progress.New(0, nil, sink)
	tracker.Stage(progress.StageArchive)
	gzipWriter := gzip.NewWriter(file)
	tarWriter := tar.NewWriter(gzipWriter)
	if err := writeNode(ctx, backend, sourcePath, sourceName, tarWriter, tracker); err != nil {
		return nil, err
	}
	if err := closeTar(tarWriter); err != nil {
		return nil, err
	}
	if err := closeGzip(gzipWriter); err != nil {
		return nil, err
	}
	if err := file.Sync(); err != nil {
		return nil, err
	}
	if err := file.Close(); err != nil {
		return nil, err
	}
	tracker.Stage(progress.StageVerify)
	if err := verifyTemporary(temporaryPath); err != nil {
		return nil, err
	}
	info, err := statTemporary(temporaryPath)
	if err != nil {
		return nil, err
	}
	success = true
	tracker.Stage(progress.StageComplete)
	return &Artifact{Path: temporaryPath, Name: sourceName + ".tar.gz", Bytes: info.Size()}, nil
}

func writeNode(ctx context.Context, backend fsx.Backend, sourcePath, entryName string, writer *tar.Writer, tracker *progress.Tracker) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	info, err := backend.Lstat(sourcePath)
	if err != nil {
		return err
	}
	link := ""
	if info.Mode()&fs.ModeSymlink != 0 {
		link, err = backend.Readlink(sourcePath)
		if err != nil {
			return err
		}
	}
	header, err := tar.FileInfoHeader(info, link)
	if err != nil {
		return err
	}
	header.Name = path.Clean(entryName)
	if _, err := safety.SafeArchiveJoin("archive-root", header.Name); err != nil {
		return err
	}
	if err := writer.WriteHeader(header); err != nil {
		return err
	}
	if info.Mode().IsRegular() {
		reader, err := backend.Open(sourcePath)
		if err != nil {
			return err
		}
		written, copyErr := io.Copy(&archiveWriter{ctx: ctx, destination: writer, tracker: tracker}, reader)
		closeErr := reader.Close()
		if copyErr != nil || closeErr != nil || written != info.Size() {
			return errors.Join(copyErr, closeErr, sizeError(written, info.Size()))
		}
		return nil
	}
	if info.Mode()&fs.ModeSymlink != 0 {
		return nil
	}
	if !info.IsDir() {
		return fmt.Errorf("unsupported archive object %q (%s)", sourcePath, info.Mode())
	}
	entries, err := backend.ReadDir(sourcePath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if err := writeNode(ctx, backend, backend.Join(sourcePath, entry.Name()), path.Join(entryName, entry.Name()), writer, tracker); err != nil {
			return err
		}
	}
	return nil
}

func sizeError(actual, expected int64) error {
	if actual == expected {
		return nil
	}
	return fmt.Errorf("archive payload size changed: wrote %d, expected %d", actual, expected)
}

type archiveWriter struct {
	ctx         context.Context
	destination io.Writer
	tracker     *progress.Tracker
}

func (w *archiveWriter) Write(data []byte) (int, error) {
	if err := w.ctx.Err(); err != nil {
		return 0, err
	}
	written, err := w.destination.Write(data)
	if written > 0 {
		w.tracker.Add(int64(written))
	}
	return written, err
}

// Verify fully streams a gzip/tar archive and validates supported entry paths.
func Verify(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()
	reader := tar.NewReader(gzipReader)
	seen := make(map[string]struct{})
	for {
		header, err := reader.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if _, err := safety.SafeArchiveJoin("archive-root", header.Name); err != nil {
			return err
		}
		if _, duplicate := seen[header.Name]; duplicate {
			return fmt.Errorf("duplicate archive entry %q", header.Name)
		}
		seen[header.Name] = struct{}{}
		switch header.Typeflag {
		case tar.TypeReg, tar.TypeRegA, tar.TypeDir, tar.TypeSymlink:
		default:
			return fmt.Errorf("unsupported archive entry %q type %d", header.Name, header.Typeflag)
		}
		if _, err := io.Copy(io.Discard, reader); err != nil {
			return err
		}
	}
}
