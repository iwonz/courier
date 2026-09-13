package archive

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/fsx"
	"github.com/iwonz/courier/internal/progress"
)

func TestCreateVerifyAndCleanup(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	if err := os.MkdirAll(filepath.Join(source, "nested"), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "nested", "file"), []byte("payload"), 0o640); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("nested/file", filepath.Join(source, "link")); err != nil {
		t.Fatal(err)
	}
	var events []progress.Event
	artifact, err := Create(context.Background(), fsx.Local{}, source, "source", root, func(event progress.Event) { events = append(events, event) })
	if err != nil {
		t.Fatal(err)
	}
	if artifact.Name != "source.tar.gz" || artifact.Bytes == 0 || events[len(events)-1].Stage != progress.StageComplete {
		t.Fatalf("artifact=%+v events=%v", artifact, events)
	}
	if info, err := os.Stat(artifact.Path); err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("info=%v err=%v", info, err)
	}
	if err := Verify(artifact.Path); err != nil {
		t.Fatal(err)
	}
	if err := artifact.Cleanup(); err != nil {
		t.Fatal(err)
	}
	if err := artifact.Cleanup(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(artifact.Path); !os.IsNotExist(err) {
		t.Fatalf("archive remains: %v", err)
	}
}

func TestCreateInputAndCancellationErrors(t *testing.T) {
	for _, test := range []struct {
		backend fsx.Backend
		path    string
		name    string
	}{
		{nil, "source", "source"},
		{fsx.Local{}, "", "source"},
		{fsx.Local{}, "source", ""},
		{fsx.Local{}, "source", "../source"},
		{fsx.Local{}, "source", `bad\name`},
	} {
		if _, err := Create(context.Background(), test.backend, test.path, test.name, t.TempDir(), nil); err == nil {
			t.Fatal("expected input error")
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := Create(ctx, fsx.Local{}, "source", "source", t.TempDir(), nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
}

func TestCreateOperationFailures(t *testing.T) {
	originalCreate, originalRemove, originalStat, originalVerify := createTemporary, removeTemporary, statTemporary, verifyTemporary
	originalTar, originalGzip := closeTar, closeGzip
	t.Cleanup(func() {
		createTemporary, removeTemporary, statTemporary, verifyTemporary = originalCreate, originalRemove, originalStat, originalVerify
		closeTar, closeGzip = originalTar, originalGzip
	})
	backend := archiveBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return archiveInfo{mode: 0o600}, nil }, open: func(string) (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(nil)), nil }}
	if _, err := Create(context.Background(), backend, "source", "source", filepath.Join(t.TempDir(), "missing"), nil); err == nil {
		t.Fatal("expected create-temp error")
	}
	for _, test := range []struct {
		name      string
		configure func(*memoryTemp)
	}{
		{"chmod", func(file *memoryTemp) { file.chmodErr = errors.New("chmod") }},
		{"tar close", func(*memoryTemp) { closeTar = func(*tar.Writer) error { return errors.New("tar close") } }},
		{"gzip close", func(*memoryTemp) { closeGzip = func(*gzip.Writer) error { return errors.New("gzip close") } }},
		{"sync", func(file *memoryTemp) { file.syncErr = errors.New("sync") }},
		{"close", func(file *memoryTemp) { file.closeErr = errors.New("close") }},
		{"verify", func(*memoryTemp) { verifyTemporary = func(string) error { return errors.New("verify") } }},
		{"stat", func(*memoryTemp) {
			statTemporary = func(string) (fs.FileInfo, error) { return nil, errors.New("stat") }
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			closeTar = originalTar
			closeGzip = originalGzip
			verifyTemporary = func(string) error { return nil }
			statTemporary = func(string) (fs.FileInfo, error) { return archiveInfo{}, nil }
			file := &memoryTemp{name: "temporary"}
			test.configure(file)
			createTemporary = func(string, string) (temporaryFile, error) { return file, nil }
			removeTemporary = func(string) error { return nil }
			if _, err := Create(context.Background(), backend, "source", "source", t.TempDir(), nil); err == nil {
				t.Fatal("expected operation error")
			}
		})
	}
}

func TestWriteNodeFailures(t *testing.T) {
	tracker := progress.New(0, nil, nil)
	regular := archiveInfo{mode: 0o600, size: 1}
	directory := archiveInfo{mode: fs.ModeDir | 0o700}
	symlink := archiveInfo{mode: fs.ModeSymlink | 0o777}
	newWriter := func() *tar.Writer { return tar.NewWriter(io.Discard) }
	lstatFailure := archiveBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return nil, errors.New("lstat") }}
	if err := writeNode(context.Background(), lstatFailure, "x", "x", newWriter(), tracker); err == nil {
		t.Fatal("expected lstat error")
	}
	readlinkFailure := archiveBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return symlink, nil }, readlink: func(string) (string, error) { return "", errors.New("readlink") }}
	if err := writeNode(context.Background(), readlinkFailure, "x", "x", newWriter(), tracker); err == nil {
		t.Fatal("expected readlink error")
	}
	headerFailure := archiveBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return archiveInfo{mode: fs.ModeSocket}, nil }}
	if err := writeNode(context.Background(), headerFailure, "x", "x", newWriter(), tracker); err == nil {
		t.Fatal("expected header error")
	}
	if err := writeNode(context.Background(), archiveBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return regular, nil }}, "x", "../x", newWriter(), tracker); err == nil {
		t.Fatal("expected unsafe entry error")
	}
	closedWriter := newWriter()
	_ = closedWriter.Close()
	if err := writeNode(context.Background(), archiveBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return regular, nil }}, "x", "x", closedWriter, tracker); err == nil {
		t.Fatal("expected header writer error")
	}
	openFailure := archiveBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return regular, nil }, open: func(string) (io.ReadCloser, error) { return nil, errors.New("open") }}
	if err := writeNode(context.Background(), openFailure, "x", "x", newWriter(), tracker); err == nil {
		t.Fatal("expected open error")
	}
	for _, reader := range []io.ReadCloser{
		&failingReader{err: errors.New("read")},
		&failingReader{data: []byte("x"), closeErr: errors.New("close")},
		io.NopCloser(bytes.NewReader(nil)),
	} {
		backend := archiveBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return regular, nil }, open: func(string) (io.ReadCloser, error) { return reader, nil }}
		if err := writeNode(context.Background(), backend, "x", "x", newWriter(), tracker); err == nil {
			t.Fatal("expected regular payload error")
		}
	}
	unsupported := archiveBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return archiveInfo{mode: fs.ModeDevice}, nil }}
	if err := writeNode(context.Background(), unsupported, "x", "x", newWriter(), tracker); err == nil {
		t.Fatal("expected unsupported object")
	}
	readDirFailure := archiveBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return directory, nil }, readDir: func(string) ([]fs.DirEntry, error) { return nil, errors.New("read-dir") }}
	if err := writeNode(context.Background(), readDirFailure, "x", "x", newWriter(), tracker); err == nil {
		t.Fatal("expected read-dir error")
	}
	entry := archiveEntry{name: "child"}
	childFailure := archiveBackend{Backend: fsx.Local{}, lstat: func(name string) (fs.FileInfo, error) {
		if name == "x" {
			return directory, nil
		}
		return nil, errors.New("child")
	}, readDir: func(string) ([]fs.DirEntry, error) { return []fs.DirEntry{entry}, nil }, join: filepath.Join}
	if err := writeNode(context.Background(), childFailure, "x", "x", newWriter(), tracker); err == nil {
		t.Fatal("expected child error")
	}
}

func TestVerifyFailures(t *testing.T) {
	root := t.TempDir()
	missing := filepath.Join(root, "missing")
	if err := Verify(missing); err == nil {
		t.Fatal("expected open error")
	}
	corrupt := filepath.Join(root, "corrupt.tar.gz")
	if err := os.WriteFile(corrupt, []byte("not gzip"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Verify(corrupt); err == nil {
		t.Fatal("expected gzip error")
	}
	for _, test := range []struct {
		name    string
		headers []tar.Header
	}{
		{"unsafe", []tar.Header{{Name: "../escape", Typeflag: tar.TypeReg}}},
		{"duplicate", []tar.Header{{Name: "same", Typeflag: tar.TypeReg}, {Name: "same", Typeflag: tar.TypeReg}}},
		{"unsupported", []tar.Header{{Name: "device", Typeflag: tar.TypeChar}}},
	} {
		archivePath := filepath.Join(root, test.name+".tar.gz")
		writeArchive(t, archivePath, test.headers)
		if err := Verify(archivePath); err == nil {
			t.Errorf("expected %s verification error", test.name)
		}
	}
	badTar := filepath.Join(root, "bad-tar.tar.gz")
	file, err := os.Create(badTar)
	if err != nil {
		t.Fatal(err)
	}
	gzipWriter := gzip.NewWriter(file)
	_, _ = gzipWriter.Write([]byte("not a tar stream"))
	_ = gzipWriter.Close()
	_ = file.Close()
	if err := Verify(badTar); err == nil {
		t.Fatal("expected tar stream error")
	}
	truncated := filepath.Join(root, "truncated.tar.gz")
	file, err = os.Create(truncated)
	if err != nil {
		t.Fatal(err)
	}
	gzipWriter = gzip.NewWriter(file)
	tarWriter := tar.NewWriter(gzipWriter)
	_ = tarWriter.WriteHeader(&tar.Header{Name: "payload", Typeflag: tar.TypeReg, Mode: 0o600, Size: 10})
	_, _ = tarWriter.Write([]byte("x"))
	_ = gzipWriter.Close()
	_ = file.Close()
	if err := Verify(truncated); err == nil {
		t.Fatal("expected truncated payload error")
	}
}

func TestSizeErrorAndArchiveWriter(t *testing.T) {
	if sizeError(1, 1) != nil || sizeError(1, 2) == nil {
		t.Fatal("unexpected size error behavior")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	writer := &archiveWriter{ctx: ctx, destination: io.Discard, tracker: progress.New(0, nil, nil)}
	if _, err := writer.Write([]byte("x")); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
}

func writeArchive(t *testing.T, name string, headers []tar.Header) {
	t.Helper()
	file, err := os.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	gzipWriter := gzip.NewWriter(file)
	tarWriter := tar.NewWriter(gzipWriter)
	for index := range headers {
		if err := tarWriter.WriteHeader(&headers[index]); err != nil {
			t.Fatal(err)
		}
	}
	if err := tarWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gzipWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}

type memoryTemp struct {
	bytes.Buffer
	name     string
	chmodErr error
	syncErr  error
	closeErr error
}

func (m *memoryTemp) Name() string            { return m.name }
func (m *memoryTemp) Chmod(fs.FileMode) error { return m.chmodErr }
func (m *memoryTemp) Sync() error             { return m.syncErr }
func (m *memoryTemp) Close() error            { return m.closeErr }

type archiveInfo struct {
	mode fs.FileMode
	size int64
}

func (archiveInfo) Name() string        { return "entry" }
func (a archiveInfo) Size() int64       { return a.size }
func (a archiveInfo) Mode() fs.FileMode { return a.mode }
func (archiveInfo) ModTime() time.Time  { return time.Time{} }
func (a archiveInfo) IsDir() bool       { return a.mode.IsDir() }
func (archiveInfo) Sys() any            { return nil }

type archiveBackend struct {
	fsx.Backend
	lstat    func(string) (fs.FileInfo, error)
	readDir  func(string) ([]fs.DirEntry, error)
	open     func(string) (io.ReadCloser, error)
	readlink func(string) (string, error)
	join     func(...string) string
}

func (a archiveBackend) Lstat(name string) (fs.FileInfo, error) {
	if a.lstat != nil {
		return a.lstat(name)
	}
	return a.Backend.Lstat(name)
}
func (a archiveBackend) ReadDir(name string) ([]fs.DirEntry, error) {
	if a.readDir != nil {
		return a.readDir(name)
	}
	return a.Backend.ReadDir(name)
}
func (a archiveBackend) Open(name string) (io.ReadCloser, error) {
	if a.open != nil {
		return a.open(name)
	}
	return a.Backend.Open(name)
}
func (a archiveBackend) Readlink(name string) (string, error) {
	if a.readlink != nil {
		return a.readlink(name)
	}
	return a.Backend.Readlink(name)
}
func (a archiveBackend) Join(parts ...string) string {
	if a.join != nil {
		return a.join(parts...)
	}
	return a.Backend.Join(parts...)
}

type archiveEntry struct{ name string }

func (a archiveEntry) Name() string             { return a.name }
func (archiveEntry) IsDir() bool                { return false }
func (archiveEntry) Type() fs.FileMode          { return 0 }
func (archiveEntry) Info() (fs.FileInfo, error) { return archiveInfo{}, nil }

type failingReader struct {
	data     []byte
	err      error
	closeErr error
}

func (f *failingReader) Read(buffer []byte) (int, error) {
	if len(f.data) > 0 {
		n := copy(buffer, f.data)
		f.data = f.data[n:]
		return n, nil
	}
	if f.err != nil {
		return 0, f.err
	}
	return 0, io.EOF
}
func (f *failingReader) Close() error { return f.closeErr }
