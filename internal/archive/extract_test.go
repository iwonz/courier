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
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/fsx"
	"github.com/iwonz/courier/internal/operation"
	"github.com/iwonz/courier/internal/progress"
	"github.com/iwonz/courier/internal/selection"
)

func TestCodecRegistry(t *testing.T) {
	good := codecStub{name: "tar.gz", extensions: []string{".TAR.GZ"}}
	registry, err := NewRegistry(good)
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := registry.Resolve("BACKUP.TAR.GZ")
	if err != nil || resolved.Name() != good.name {
		t.Fatalf("resolved=%v err=%v", resolved, err)
	}
	if _, err := registry.Resolve("backup.tgz"); !errors.Is(err, ErrCodec) {
		t.Fatalf("unsupported error=%v", err)
	}
	if DefaultRegistry().codecs[0].Name() != "tar.gz" {
		t.Fatal("unexpected default registry")
	}

	for _, test := range []struct {
		name   string
		codecs []Codec
	}{
		{"nil", []Codec{nil}},
		{"empty name", []Codec{codecStub{extensions: []string{".x"}}}},
		{"no extension", []Codec{codecStub{name: "x"}}},
		{"duplicate name", []Codec{codecStub{name: "x", extensions: []string{".x"}}, codecStub{name: "x", extensions: []string{".y"}}}},
		{"invalid extension", []Codec{codecStub{name: "x", extensions: []string{"x"}}}},
		{"empty extension", []Codec{codecStub{name: "x", extensions: []string{"."}}}},
		{"duplicate extension", []Codec{codecStub{name: "x", extensions: []string{".x"}}, codecStub{name: "y", extensions: []string{".X"}}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, err := NewRegistry(test.codecs...); !errors.Is(err, ErrCodec) {
				t.Fatalf("error=%v", err)
			}
		})
	}

	long, err := NewRegistry(codecStub{name: "gz", extensions: []string{".gz"}}, codecStub{name: "tar.gz", extensions: []string{".tar.gz"}})
	if err != nil {
		t.Fatal(err)
	}
	resolved, err = long.Resolve("data.tar.gz")
	if err != nil || resolved.Name() != "tar.gz" {
		t.Fatalf("longest resolved=%v err=%v", resolved, err)
	}
}

func TestExtractTarGzipTransactions(t *testing.T) {
	t.Run("absent root", func(t *testing.T) {
		root := t.TempDir()
		archivePath := filepath.Join(root, "source.tar.gz")
		modTime := time.Unix(1_700_000_000, 0)
		writeTarGzip(t, archivePath, gzip.NoCompression, []tarEntry{
			{header: tar.Header{Name: "source/", Typeflag: tar.TypeDir, Mode: 0o750, ModTime: modTime}},
			{header: tar.Header{Name: "source/file.txt", Typeflag: tar.TypeReg, Mode: 0o640, ModTime: modTime}, data: []byte("payload")},
		})
		destination := filepath.Join(root, "output")
		var events []progress.Event
		result, err := DefaultRegistry().Extract(context.Background(), ExtractionRequest{
			SourceFS: fsx.Local{}, SourcePath: archivePath, SourceName: "source.tar.gz",
			DestinationFS: fsx.Local{}, DestinationRoot: destination,
			Progress: func(event progress.Event) { events = append(events, event) },
		})
		data, readErr := os.ReadFile(filepath.Join(destination, "source", "file.txt"))
		if err != nil || readErr != nil || string(data) != "payload" || result.Bytes != 7 || result.Entries != 2 || result.Committed != 1 {
			t.Fatalf("result=%+v data=%q readErr=%v err=%v", result, data, readErr, err)
		}
		if len(events) == 0 || events[len(events)-1].Stage != progress.StageComplete {
			t.Fatalf("events=%v", events)
		}
		assertNoExtractionStage(t, destination)
	})

	t.Run("existing root preserves unrelated and applies selection", func(t *testing.T) {
		root := t.TempDir()
		archivePath := filepath.Join(root, "items.tar.gz")
		writeTarGzip(t, archivePath, gzip.NoCompression, []tarEntry{
			{header: tar.Header{Name: "keep.txt", Typeflag: tar.TypeReg, Mode: 0o600}, data: []byte("keep")},
			{header: tar.Header{Name: "drop.tmp", Typeflag: tar.TypeReg, Mode: 0o600}, data: []byte("drop")},
		})
		destination := filepath.Join(root, "destination")
		if err := os.Mkdir(destination, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(destination, "unrelated"), []byte("safe"), 0o600); err != nil {
			t.Fatal(err)
		}
		selector, err := selection.Compile([]operation.SelectionRule{{Kind: operation.SelectionGitignore, Value: "*.tmp"}}, nil)
		if err != nil {
			t.Fatal(err)
		}
		result, err := DefaultRegistry().Extract(context.Background(), ExtractionRequest{
			SourceFS: fsx.Local{}, SourcePath: archivePath, SourceName: "items.tar.gz",
			DestinationFS: fsx.Local{}, DestinationRoot: destination, Selector: selector,
			Limits: Limits{Unlimited: true},
		})
		keep, keepErr := os.ReadFile(filepath.Join(destination, "keep.txt"))
		unrelated, unrelatedErr := os.ReadFile(filepath.Join(destination, "unrelated"))
		_, dropErr := os.Lstat(filepath.Join(destination, "drop.tmp"))
		if err != nil || keepErr != nil || unrelatedErr != nil || string(keep) != "keep" || string(unrelated) != "safe" || !errors.Is(dropErr, fs.ErrNotExist) || result.Bytes != 4 || result.Committed != 1 {
			t.Fatalf("result=%+v keep=%q,%v unrelated=%q,%v drop=%v err=%v", result, keep, keepErr, unrelated, unrelatedErr, dropErr, err)
		}
		assertNoExtractionStage(t, destination)
	})

	t.Run("preflight collision", func(t *testing.T) {
		root := t.TempDir()
		archivePath := filepath.Join(root, "items.tar.gz")
		writeTarGzip(t, archivePath, gzip.NoCompression, []tarEntry{
			{header: tar.Header{Name: "folder/", Typeflag: tar.TypeDir, Mode: 0o700}},
			{header: tar.Header{Name: "folder/file", Typeflag: tar.TypeReg, Mode: 0o600}, data: []byte("new")},
		})
		destination := filepath.Join(root, "destination")
		if err := os.MkdirAll(filepath.Join(destination, "folder"), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(destination, "folder", "old"), []byte("keep"), 0o600); err != nil {
			t.Fatal(err)
		}
		_, err := DefaultRegistry().Extract(context.Background(), ExtractionRequest{
			SourceFS: fsx.Local{}, SourcePath: archivePath, SourceName: "items.tar.gz", DestinationFS: fsx.Local{}, DestinationRoot: destination,
		})
		data, readErr := os.ReadFile(filepath.Join(destination, "folder", "old"))
		if !errors.Is(err, fsx.ErrDestinationExists) || readErr != nil || string(data) != "keep" {
			t.Fatalf("data=%q readErr=%v err=%v", data, readErr, err)
		}
		assertNoExtractionStage(t, destination)
	})
}

func TestExtractLateCommitFailureReportsConfirmedBytes(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "items.tar.gz")
	writeTarGzip(t, archivePath, gzip.NoCompression, []tarEntry{
		{header: tar.Header{Name: "a.txt", Typeflag: tar.TypeReg, Mode: 0o600}, data: []byte("aaa")},
		{header: tar.Header{Name: "b.txt", Typeflag: tar.TypeReg, Mode: 0o600}, data: []byte("bb")},
	})
	destination := filepath.Join(root, "destination")
	if err := os.Mkdir(destination, 0o700); err != nil {
		t.Fatal(err)
	}
	commits := 0
	backend := extractBackend{Backend: fsx.Local{}}
	backend.commitAbsent = func(stage, final string) error {
		commits++
		if commits == 2 {
			return fsx.ErrDestinationExists
		}
		return fsx.Local{}.CommitAbsent(stage, final)
	}
	result, err := DefaultRegistry().Extract(context.Background(), ExtractionRequest{
		SourceFS: fsx.Local{}, SourcePath: archivePath, SourceName: "items.tar.gz", DestinationFS: backend, DestinationRoot: destination,
	})
	var typed *ExtractionError
	if !errors.As(err, &typed) || typed.Stage != progress.StageCommit || typed.Confirmed != 3 || result.Bytes != 3 || result.Committed != 1 {
		t.Fatalf("result=%+v typed=%+v err=%v", result, typed, err)
	}
	data, readErr := os.ReadFile(filepath.Join(destination, "a.txt"))
	if readErr != nil || string(data) != "aaa" {
		t.Fatalf("data=%q readErr=%v", data, readErr)
	}
	assertNoExtractionStage(t, destination)
}

func TestTarGzipValidation(t *testing.T) {
	root := t.TempDir()
	for _, test := range []struct {
		name       string
		entries    []tarEntry
		limits     Limits
		compressed int64
		want       error
	}{
		{"unsafe path", []tarEntry{{header: tar.Header{Name: "../escape", Typeflag: tar.TypeReg, Mode: 0o600}}}, Limits{Unlimited: true}, -1, nil},
		{"normalized duplicate", []tarEntry{{header: tar.Header{Name: "a/", Typeflag: tar.TypeDir}}, {header: tar.Header{Name: "a", Typeflag: tar.TypeDir}}}, Limits{Unlimited: true}, -1, nil},
		{"depth", []tarEntry{{header: tar.Header{Name: strings.Repeat("a/", MaxExtractionDepth) + "file", Typeflag: tar.TypeReg, Mode: 0o600}}}, Limits{Unlimited: true}, -1, ErrArchiveLimit},
		{"long path", []tarEntry{{header: tar.Header{Name: strings.Repeat("a", maxArchivePathBytes+1), Typeflag: tar.TypeReg, Mode: 0o600}}}, Limits{Unlimited: true}, -1, ErrArchiveLimit},
		{"unsupported type", []tarEntry{{header: tar.Header{Name: "device", Typeflag: tar.TypeChar}}}, Limits{Unlimited: true}, -1, nil},
		{"parent after child", []tarEntry{{header: tar.Header{Name: "a/b", Typeflag: tar.TypeReg, Mode: 0o600}}, {header: tar.Header{Name: "a", Typeflag: tar.TypeReg, Mode: 0o600}}}, Limits{Unlimited: true}, -1, nil},
		{"child after file", []tarEntry{{header: tar.Header{Name: "a", Typeflag: tar.TypeReg, Mode: 0o600}}, {header: tar.Header{Name: "a/b", Typeflag: tar.TypeReg, Mode: 0o600}}}, Limits{Unlimited: true}, -1, nil},
		{"unsafe symlink", []tarEntry{{header: tar.Header{Name: "source/link", Typeflag: tar.TypeSymlink, Linkname: "../../escape"}}}, Limits{Unlimited: true}, -1, nil},
		{"configured bytes", []tarEntry{{header: tar.Header{Name: "file", Typeflag: tar.TypeReg, Mode: 0o600}, data: []byte("xx")}}, Limits{MaxBytes: 1, Configured: true}, -1, ErrArchiveLimit},
	} {
		t.Run(test.name, func(t *testing.T) {
			name := filepath.Join(root, strings.ReplaceAll(test.name, " ", "-")+".tar.gz")
			writeTarGzip(t, name, gzip.NoCompression, test.entries)
			info, err := os.Stat(name)
			if err != nil {
				t.Fatal(err)
			}
			compressed := info.Size()
			if test.compressed >= 0 {
				compressed = test.compressed
			}
			_, err = (TarGzipCodec{}).Inspect(context.Background(), InspectRequest{SourceFS: fsx.Local{}, SourcePath: name, CompressedSize: compressed, Limits: test.limits})
			if err == nil || test.want != nil && !errors.Is(err, test.want) {
				t.Fatalf("error=%v", err)
			}
		})
	}

	ratioArchive := filepath.Join(root, "ratio.tar.gz")
	writeTarGzip(t, ratioArchive, gzip.BestCompression, []tarEntry{{header: tar.Header{Name: "zeros", Typeflag: tar.TypeReg, Mode: 0o600}, data: make([]byte, 1<<20)}})
	ratioInfo, err := os.Stat(ratioArchive)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := (TarGzipCodec{}).Inspect(context.Background(), InspectRequest{SourceFS: fsx.Local{}, SourcePath: ratioArchive, CompressedSize: ratioInfo.Size(), Limits: Limits{Unlimited: true}}); !errors.Is(err, ErrArchiveLimit) {
		t.Fatalf("ratio error=%v size=%d", err, ratioInfo.Size())
	}

	if !exceedsExpansionRatio(1, 0) || exceedsExpansionRatio(0, 0) || !exceedsExpansionRatio(101, 1) || exceedsExpansionRatio(100, 1) || exceedsExpansionRatio(1, int64(^uint64(0)>>1)) {
		t.Fatal("unexpected expansion-ratio boundary")
	}
	for _, header := range []*tar.Header{
		{Name: "file", Typeflag: tar.TypeReg, Mode: 0o600},
		{Name: "legacy", Typeflag: tar.TypeRegA, Mode: 0o600},
		{Name: "dir", Typeflag: tar.TypeDir, Mode: 0o700},
		{Name: "link", Typeflag: tar.TypeSymlink, Mode: 0o777},
	} {
		if _, err := archiveMode(header); err != nil {
			t.Fatalf("mode %d: %v", header.Typeflag, err)
		}
	}
}

func TestTarGzipInspectionFailures(t *testing.T) {
	root := t.TempDir()
	missing := filepath.Join(root, "missing.tar.gz")
	for _, request := range []InspectRequest{
		{},
		{SourceFS: fsx.Local{}, SourcePath: "", CompressedSize: 0},
		{SourceFS: fsx.Local{}, SourcePath: missing, CompressedSize: -1},
	} {
		if _, err := (TarGzipCodec{}).Inspect(context.Background(), request); !errors.Is(err, ErrCodec) {
			t.Fatalf("invalid request error=%v", err)
		}
	}
	if _, err := (TarGzipCodec{}).Inspect(context.Background(), InspectRequest{SourceFS: fsx.Local{}, SourcePath: missing}); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("open error=%v", err)
	}
	corrupt := filepath.Join(root, "corrupt.tar.gz")
	if err := os.WriteFile(corrupt, []byte("bad"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := (TarGzipCodec{}).Inspect(context.Background(), InspectRequest{SourceFS: fsx.Local{}, SourcePath: corrupt, CompressedSize: 3}); err == nil {
		t.Fatal("expected gzip error")
	}
	badTar := filepath.Join(root, "bad-tar.tar.gz")
	writeGzipBytes(t, badTar, []byte("not tar"))
	if _, err := (TarGzipCodec{}).Inspect(context.Background(), InspectRequest{SourceFS: fsx.Local{}, SourcePath: badTar, CompressedSize: 7}); err == nil {
		t.Fatal("expected tar error")
	}

	valid := filepath.Join(root, "valid.tar.gz")
	writeTarGzip(t, valid, gzip.NoCompression, []tarEntry{{header: tar.Header{Name: "file", Typeflag: tar.TypeReg, Mode: 0o600}, data: []byte("data")}})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := (TarGzipCodec{}).Inspect(ctx, InspectRequest{SourceFS: fsx.Local{}, SourcePath: valid, CompressedSize: 100, Limits: Limits{Unlimited: true}}); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation error=%v", err)
	}
	visitErr := errors.New("visit")
	if _, err := (TarGzipCodec{}).Inspect(context.Background(), InspectRequest{SourceFS: fsx.Local{}, SourcePath: valid, CompressedSize: 100, Limits: Limits{Unlimited: true}, Visit: func(Entry) error { return visitErr }}); !errors.Is(err, visitErr) {
		t.Fatalf("visit error=%v", err)
	}
	readerErr := errors.New("reader")
	backend := extractBackend{Backend: fsx.Local{}, open: func(string) (io.ReadCloser, error) {
		data, err := os.ReadFile(valid)
		if err != nil {
			return nil, err
		}
		return &failingReader{data: data[:len(data)-4], err: readerErr, closeErr: errors.New("close")}, nil
	}}
	if _, err := (TarGzipCodec{}).Inspect(context.Background(), InspectRequest{SourceFS: backend, SourcePath: valid, CompressedSize: 100, Limits: Limits{Unlimited: true}}); err == nil {
		t.Fatal("expected truncated reader error")
	}

	originalCount, originalValidate, originalDiscard := countArchiveEntry, validatePayload, discardPayload
	t.Cleanup(func() {
		countArchiveEntry, validatePayload, discardPayload = originalCount, originalValidate, originalDiscard
	})
	countErr := errors.New("count")
	countArchiveEntry = func(int) (int, error) { return 0, countErr }
	if _, err := (TarGzipCodec{}).Inspect(context.Background(), InspectRequest{SourceFS: fsx.Local{}, SourcePath: valid, CompressedSize: 100, Limits: Limits{Unlimited: true}}); !errors.Is(err, countErr) {
		t.Fatalf("count error=%v", err)
	}
	countArchiveEntry = originalCount
	validateErr := errors.New("payload")
	validatePayload = func(string, fs.FileMode, int64) error { return validateErr }
	if _, err := (TarGzipCodec{}).Inspect(context.Background(), InspectRequest{SourceFS: fsx.Local{}, SourcePath: valid, CompressedSize: 100, Limits: Limits{Unlimited: true}}); !errors.Is(err, validateErr) {
		t.Fatalf("payload error=%v", err)
	}
	validatePayload = originalValidate
	discardErr := errors.New("discard")
	discardPayload = func(io.Reader, []byte) (int64, error) { return 0, discardErr }
	if _, err := (TarGzipCodec{}).Inspect(context.Background(), InspectRequest{SourceFS: fsx.Local{}, SourcePath: valid, CompressedSize: 100, Limits: Limits{Unlimited: true}}); !errors.Is(err, discardErr) {
		t.Fatalf("discard error=%v", err)
	}
}

func TestTarGzipExcludedDirectorySkipsDescendants(t *testing.T) {
	root := t.TempDir()
	name := filepath.Join(root, "selection.tar.gz")
	writeTarGzip(t, name, gzip.NoCompression, []tarEntry{
		{header: tar.Header{Name: "skip/", Typeflag: tar.TypeDir, Mode: 0o700}},
		{header: tar.Header{Name: "skip/file", Typeflag: tar.TypeReg, Mode: 0o600}, data: []byte("hidden")},
		{header: tar.Header{Name: "keep", Typeflag: tar.TypeReg, Mode: 0o600}, data: []byte("shown")},
	})
	info, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	var visited []Entry
	summary, err := (TarGzipCodec{}).Inspect(context.Background(), InspectRequest{
		SourceFS: fsx.Local{}, SourcePath: name, CompressedSize: info.Size(), Limits: Limits{Unlimited: true},
		Selector: selectorFunc(func(name string, _ bool) bool { return name != "skip" }),
		Visit:    func(entry Entry) error { visited = append(visited, entry); return nil },
	})
	if err != nil || summary.SelectedBytes != 5 || len(visited) != 3 || visited[0].Selected || visited[1].Selected || !visited[2].Selected {
		t.Fatalf("summary=%+v visited=%+v err=%v", summary, visited, err)
	}
}

func TestExtractionRequestFailures(t *testing.T) {
	registry := DefaultRegistry()
	root := t.TempDir()
	source := filepath.Join(root, "source.tar.gz")
	writeTarGzip(t, source, gzip.NoCompression, []tarEntry{{header: tar.Header{Name: "file", Typeflag: tar.TypeReg, Mode: 0o600}, data: []byte("x")}})
	valid := ExtractionRequest{SourceFS: fsx.Local{}, SourcePath: source, SourceName: "source.tar.gz", DestinationFS: fsx.Local{}, DestinationRoot: filepath.Join(root, "out")}
	for _, mutate := range []func(*ExtractionRequest){
		func(r *ExtractionRequest) { r.SourceFS = nil },
		func(r *ExtractionRequest) { r.SourcePath = "" },
		func(r *ExtractionRequest) { r.SourceName = "" },
		func(r *ExtractionRequest) { r.DestinationFS = nil },
		func(r *ExtractionRequest) { r.DestinationRoot = "" },
	} {
		request := valid
		mutate(&request)
		if _, err := registry.Extract(context.Background(), request); !errors.Is(err, ErrCodec) {
			t.Fatalf("invalid request error=%v", err)
		}
	}
	request := valid
	request.SourceName = "source.tgz"
	if _, err := registry.Extract(context.Background(), request); !errors.Is(err, ErrCodec) {
		t.Fatalf("resolve error=%v", err)
	}
	request = valid
	request.SourcePath = filepath.Join(root, "missing.tar.gz")
	if _, err := registry.Extract(context.Background(), request); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("stat error=%v", err)
	}
	directoryArchive := filepath.Join(root, "directory.tar.gz")
	if err := os.Mkdir(directoryArchive, 0o700); err != nil {
		t.Fatal(err)
	}
	request = valid
	request.SourcePath = directoryArchive
	if _, err := registry.Extract(context.Background(), request); !errors.Is(err, ErrCodec) {
		t.Fatalf("regular source error=%v", err)
	}
	fileDestination := filepath.Join(root, "destination-file")
	if err := os.WriteFile(fileDestination, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	request = valid
	request.DestinationRoot = fileDestination
	if _, err := registry.Extract(context.Background(), request); err == nil {
		t.Fatal("expected destination type error")
	}
	backendErr := errors.New("destination stat")
	request = valid
	request.DestinationFS = extractBackend{Backend: fsx.Local{}, lstat: func(name string) (fs.FileInfo, error) {
		if name == request.DestinationRoot {
			return nil, backendErr
		}
		return fsx.Local{}.Lstat(name)
	}}
	if _, err := registry.Extract(context.Background(), request); !errors.Is(err, backendErr) {
		t.Fatalf("destination stat error=%v", err)
	}
}

func TestExtractionTransactionFailureStages(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source.arc")
	if err := os.WriteFile(source, []byte("archive"), 0o600); err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(root, "destination")
	if err := os.Mkdir(destination, 0o700); err != nil {
		t.Fatal(err)
	}
	baseCodec := codecStub{
		name:       "stub",
		extensions: []string{".arc"},
		inspect: func(_ context.Context, request InspectRequest) (Summary, error) {
			if request.Visit != nil {
				if err := request.Visit(Entry{Name: "item", Mode: 0o600, Size: 4, Selected: true}); err != nil {
					return Summary{}, err
				}
			}
			return Summary{Entries: 1, ExpandedBytes: 4, SelectedBytes: 4}, nil
		},
		extract: func(_ context.Context, request CodecExtractRequest) (Summary, error) {
			writer, err := request.DestinationFS.Create(request.DestinationFS.Join(request.StageRoot, "item"), 0o600)
			if err != nil {
				return Summary{}, err
			}
			if _, err := writer.Write([]byte("data")); err != nil {
				return Summary{}, err
			}
			if err := writer.Close(); err != nil {
				return Summary{}, err
			}
			return Summary{Entries: 1, ExpandedBytes: 4, SelectedBytes: 4}, nil
		},
	}
	newRegistry := func(codec codecStub) Registry {
		registry, err := NewRegistry(codec)
		if err != nil {
			t.Fatal(err)
		}
		return registry
	}
	requestFor := func(backend fsx.Backend) ExtractionRequest {
		return ExtractionRequest{SourceFS: fsx.Local{}, SourcePath: source, SourceName: "source.arc", DestinationFS: backend, DestinationRoot: destination, Limits: Limits{Unlimited: true}}
	}
	assertStage := func(t *testing.T, result ExtractionResult, err error, stage progress.Stage, cause error) {
		t.Helper()
		var typed *ExtractionError
		if !errors.As(err, &typed) || typed.Stage != stage || cause != nil && !errors.Is(err, cause) || typed.Confirmed != result.Bytes {
			t.Fatalf("result=%+v typed=%+v err=%v", result, typed, err)
		}
	}

	t.Run("inspection destination error", func(t *testing.T) {
		want := errors.New("inspect destination")
		backend := extractBackend{Backend: fsx.Local{}, lstat: func(name string) (fs.FileInfo, error) {
			if name == filepath.Join(destination, "item") {
				return nil, want
			}
			return fsx.Local{}.Lstat(name)
		}}
		result, err := newRegistry(baseCodec).Extract(context.Background(), requestFor(backend))
		assertStage(t, result, err, progress.StagePreflight, want)
	})

	t.Run("token", func(t *testing.T) {
		original := extractRandomRead
		t.Cleanup(func() { extractRandomRead = original })
		want := errors.New("token")
		extractRandomRead = func([]byte) (int, error) { return 0, want }
		result, err := newRegistry(baseCodec).Extract(context.Background(), requestFor(fsx.Local{}))
		assertStage(t, result, err, progress.StageExtract, want)
	})

	for _, test := range []struct {
		name      string
		configure func(*extractBackend, error)
		stage     progress.Stage
	}{
		{"initial cleanup", func(backend *extractBackend, want error) { backend.removeAll = func(string) error { return want } }, progress.StageExtract},
		{"parent mkdir", func(backend *extractBackend, want error) {
			backend.mkdirAll = func(string, fs.FileMode) error { return want }
		}, progress.StageExtract},
		{"stage mkdir", func(backend *extractBackend, want error) {
			calls := 0
			backend.mkdirAll = func(name string, mode fs.FileMode) error {
				calls++
				if calls == 2 {
					return want
				}
				return fsx.Local{}.MkdirAll(name, mode)
			}
		}, progress.StageExtract},
		{"stage read", func(backend *extractBackend, want error) {
			backend.readDir = func(string) ([]fs.DirEntry, error) { return nil, want }
		}, progress.StageCommit},
		{"cleanup", func(backend *extractBackend, want error) {
			calls := 0
			backend.removeAll = func(name string) error {
				calls++
				if calls >= 2 {
					return want
				}
				return fsx.Local{}.RemoveAll(name)
			}
			backend.readDir = func(string) ([]fs.DirEntry, error) { return nil, nil }
		}, progress.StageCleanup},
	} {
		t.Run(test.name, func(t *testing.T) {
			want := errors.New(test.name)
			backend := extractBackend{Backend: fsx.Local{}}
			test.configure(&backend, want)
			result, err := newRegistry(baseCodec).Extract(context.Background(), requestFor(backend))
			assertStage(t, result, err, test.stage, want)
		})
	}

	t.Run("codec extract", func(t *testing.T) {
		want := errors.New("codec extract")
		codec := baseCodec
		codec.extract = func(context.Context, CodecExtractRequest) (Summary, error) { return Summary{}, want }
		result, err := newRegistry(codec).Extract(context.Background(), requestFor(fsx.Local{}))
		assertStage(t, result, err, progress.StageExtract, want)
	})

	t.Run("absent root commit", func(t *testing.T) {
		want := errors.New("commit")
		absent := filepath.Join(root, "absent")
		backend := extractBackend{Backend: fsx.Local{}, commitAbsent: func(string, string) error { return want }}
		request := requestFor(backend)
		request.DestinationRoot = absent
		result, err := newRegistry(baseCodec).Extract(context.Background(), request)
		assertStage(t, result, err, progress.StageCommit, want)
	})

	for _, collision := range []struct {
		name string
		err  error
	}{
		{"appeared", nil},
		{"stat error", errors.New("late stat")},
	} {
		t.Run("late "+collision.name, func(t *testing.T) {
			checks := 0
			backend := extractBackend{Backend: fsx.Local{}, lstat: func(name string) (fs.FileInfo, error) {
				if name == filepath.Join(destination, "item") {
					checks++
					if checks == 1 {
						return nil, fs.ErrNotExist
					}
					if collision.err != nil {
						return nil, collision.err
					}
					return archiveInfo{mode: 0o600}, nil
				}
				return fsx.Local{}.Lstat(name)
			}}
			result, err := newRegistry(baseCodec).Extract(context.Background(), requestFor(backend))
			want := collision.err
			if want == nil {
				want = fsx.ErrDestinationExists
			}
			assertStage(t, result, err, progress.StageCommit, want)
		})
	}
}

func TestTarGzipExtractionFailures(t *testing.T) {
	root := t.TempDir()
	fileArchive := filepath.Join(root, "file.tar.gz")
	writeTarGzip(t, fileArchive, gzip.NoCompression, []tarEntry{{header: tar.Header{Name: "file", Typeflag: tar.TypeReg, Mode: 0o640}, data: []byte("data")}})
	directoryArchive := filepath.Join(root, "directory.tar.gz")
	writeTarGzip(t, directoryArchive, gzip.NoCompression, []tarEntry{{header: tar.Header{Name: "source/", Typeflag: tar.TypeDir, Mode: 0o750}}})
	symlinkArchive := filepath.Join(root, "symlink.tar.gz")
	writeTarGzip(t, symlinkArchive, gzip.NoCompression, []tarEntry{
		{header: tar.Header{Name: "source/", Typeflag: tar.TypeDir, Mode: 0o700}},
		{header: tar.Header{Name: "source/link", Typeflag: tar.TypeSymlink, Mode: 0o777, Linkname: "file"}},
	})
	request := func(source string, backend fsx.Backend) CodecExtractRequest {
		info, err := os.Stat(source)
		if err != nil {
			t.Fatal(err)
		}
		return CodecExtractRequest{SourceFS: fsx.Local{}, SourcePath: source, CompressedSize: info.Size(), DestinationFS: backend, StageRoot: filepath.Join(root, "stage"), Limits: Limits{Unlimited: true}}
	}
	for _, test := range []struct {
		name      string
		source    string
		configure func(*extractBackend, error)
	}{
		{"directory mkdir", directoryArchive, func(backend *extractBackend, want error) {
			backend.mkdirAll = func(string, fs.FileMode) error { return want }
		}},
		{"file parent mkdir", fileArchive, func(backend *extractBackend, want error) {
			backend.mkdirAll = func(string, fs.FileMode) error { return want }
		}},
		{"symlink", symlinkArchive, func(backend *extractBackend, want error) {
			backend.symlink = func(string, string) error { return want }
		}},
		{"create", fileArchive, func(backend *extractBackend, want error) {
			backend.create = func(string, fs.FileMode) (fsx.Writable, error) { return nil, want }
		}},
		{"write", fileArchive, func(backend *extractBackend, want error) {
			backend.create = func(string, fs.FileMode) (fsx.Writable, error) { return &faultWritable{writeErr: want}, nil }
		}},
		{"sync", fileArchive, func(backend *extractBackend, want error) {
			backend.create = func(string, fs.FileMode) (fsx.Writable, error) { return &faultWritable{syncErr: want}, nil }
		}},
		{"close", fileArchive, func(backend *extractBackend, want error) {
			backend.create = func(string, fs.FileMode) (fsx.Writable, error) { return &faultWritable{closeErr: want}, nil }
		}},
		{"file chmod", fileArchive, func(backend *extractBackend, want error) {
			backend.chmod = func(string, fs.FileMode) error { return want }
		}},
		{"file chtimes", fileArchive, func(backend *extractBackend, want error) {
			backend.chtimes = func(string, time.Time, time.Time) error { return want }
		}},
		{"directory chmod", directoryArchive, func(backend *extractBackend, want error) {
			backend.chmod = func(string, fs.FileMode) error { return want }
		}},
		{"directory chtimes", directoryArchive, func(backend *extractBackend, want error) {
			backend.chtimes = func(string, time.Time, time.Time) error { return want }
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			_ = os.RemoveAll(filepath.Join(root, "stage"))
			want := errors.New(test.name)
			backend := extractBackend{Backend: fsx.Local{}}
			test.configure(&backend, want)
			if _, err := (TarGzipCodec{}).Extract(context.Background(), request(test.source, backend)); !errors.Is(err, want) {
				t.Fatalf("error=%v", err)
			}
		})
	}
}

func TestExtractionHelpers(t *testing.T) {
	if count, err := nextEntryCount(MaxExtractionEntries - 1); err != nil || count != MaxExtractionEntries {
		t.Fatalf("count=%d err=%v", count, err)
	}
	if _, err := nextEntryCount(MaxExtractionEntries); !errors.Is(err, ErrArchiveLimit) {
		t.Fatalf("entry limit error=%v", err)
	}
	if err := validateEntryPayload("directory", fs.ModeDir, 1); err == nil {
		t.Fatal("expected directory payload failure")
	}
	if err := validateEntryPayload("file", 0o600, 1); err != nil {
		t.Fatalf("regular payload error=%v", err)
	}
	if value, err := addExpandedSize(1, 2, 100, Limits{Unlimited: true}); err != nil || value != 3 {
		t.Fatalf("expanded=%d err=%v", value, err)
	}
	for _, test := range []struct {
		current, size, compressed int64
		limits                    Limits
	}{
		{0, -1, 1, Limits{Unlimited: true}},
		{int64(^uint64(0) >> 1), 1, int64(^uint64(0) >> 1), Limits{Unlimited: true}},
		{0, 2, 100, Limits{MaxBytes: 1, Configured: true}},
		{0, DefaultMaxExtracted + 1, int64(^uint64(0) >> 1), Limits{}},
		{0, 101, 1, Limits{Unlimited: true}},
	} {
		if _, err := addExpandedSize(test.current, test.size, test.compressed, test.limits); err == nil {
			t.Fatalf("expected expanded-size error for %+v", test)
		}
	}
	if got := joinArchivePath(fsx.Local{}, "root", "a/b"); got != filepath.Join("root", "a", "b") {
		t.Fatalf("joined=%q", got)
	}
	if extractionStage("/", "token") != "/.courier-extract-token" || extractionStage("root/", "token") != "root.courier-extract-token" {
		t.Fatal("unexpected stage path")
	}
	originalRandom := extractRandomRead
	t.Cleanup(func() { extractRandomRead = originalRandom })
	extractRandomRead = func(data []byte) (int, error) {
		copy(data, []byte{0, 1, 2, 3, 4, 5, 6, 7})
		return len(data), nil
	}
	if token, err := extractionToken(); err != nil || token != "0001020304050607" {
		t.Fatalf("token=%q err=%v", token, err)
	}
	want := errors.New("random")
	extractRandomRead = func([]byte) (int, error) { return 0, want }
	if _, err := extractionToken(); !errors.Is(err, want) {
		t.Fatalf("random error=%v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	writer := &extractWriter{ctx: ctx, destination: io.Discard}
	if _, err := writer.Write([]byte("x")); !errors.Is(err, context.Canceled) {
		t.Fatalf("writer cancellation=%v", err)
	}
	reader := &contextReader{ctx: ctx, reader: strings.NewReader("x")}
	if _, err := reader.Read(make([]byte, 1)); !errors.Is(err, context.Canceled) {
		t.Fatalf("reader cancellation=%v", err)
	}
	ctx = context.Background()
	reader = &contextReader{ctx: ctx, reader: strings.NewReader("x")}
	data := make([]byte, 1)
	if count, err := reader.Read(data); err != nil || count != 1 || string(data) != "x" {
		t.Fatalf("read=%d,%q err=%v", count, data, err)
	}
	var tracked []progress.Event
	writer = &extractWriter{ctx: ctx, destination: &shortWriter{}, tracker: progress.New(1, nil, func(event progress.Event) { tracked = append(tracked, event) })}
	if count, err := writer.Write([]byte("xx")); err != nil || count != 1 || tracked[len(tracked)-1].Current != 1 {
		t.Fatalf("write=%d events=%v err=%v", count, tracked, err)
	}

	cause := errors.New("cause")
	typed := &ExtractionError{Stage: progress.StageExtract, Confirmed: 2, Cause: cause}
	if typed.Error() != "cause" || !errors.Is(typed, cause) {
		t.Fatal("extraction error wrapper is inconsistent")
	}
}

type tarEntry struct {
	header tar.Header
	data   []byte
}

func writeTarGzip(t *testing.T, name string, level int, entries []tarEntry) {
	t.Helper()
	file, err := os.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	gzipWriter, err := gzip.NewWriterLevel(file, level)
	if err != nil {
		t.Fatal(err)
	}
	tarWriter := tar.NewWriter(gzipWriter)
	for _, entry := range entries {
		header := entry.header
		header.Size = int64(len(entry.data))
		if err := tarWriter.WriteHeader(&header); err != nil {
			t.Fatal(err)
		}
		if len(entry.data) != 0 {
			if _, err := tarWriter.Write(entry.data); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := errors.Join(tarWriter.Close(), gzipWriter.Close(), file.Close()); err != nil {
		t.Fatal(err)
	}
}

func writeGzipBytes(t *testing.T, name string, data []byte) {
	t.Helper()
	file, err := os.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	writer := gzip.NewWriter(file)
	_, writeErr := writer.Write(data)
	if err := errors.Join(writeErr, writer.Close(), file.Close()); err != nil {
		t.Fatal(err)
	}
}

func assertNoExtractionStage(t *testing.T, destination string) {
	t.Helper()
	matches, err := filepath.Glob(strings.TrimRight(destination, `/\`) + ".courier-extract-*")
	if err != nil || len(matches) != 0 {
		t.Fatalf("stages=%v err=%v", matches, err)
	}
}

type codecStub struct {
	name       string
	extensions []string
	inspect    func(context.Context, InspectRequest) (Summary, error)
	extract    func(context.Context, CodecExtractRequest) (Summary, error)
}

type selectorFunc func(string, bool) bool

func (f selectorFunc) Include(name string, directory bool) bool { return f(name, directory) }

func (c codecStub) Name() string         { return c.name }
func (c codecStub) Extensions() []string { return append([]string(nil), c.extensions...) }
func (c codecStub) Inspect(ctx context.Context, request InspectRequest) (Summary, error) {
	if c.inspect != nil {
		return c.inspect(ctx, request)
	}
	return Summary{}, nil
}
func (c codecStub) Extract(ctx context.Context, request CodecExtractRequest) (Summary, error) {
	if c.extract != nil {
		return c.extract(ctx, request)
	}
	return Summary{}, nil
}

type extractBackend struct {
	fsx.Backend
	lstat        func(string) (fs.FileInfo, error)
	readDir      func(string) ([]fs.DirEntry, error)
	open         func(string) (io.ReadCloser, error)
	create       func(string, fs.FileMode) (fsx.Writable, error)
	mkdirAll     func(string, fs.FileMode) error
	removeAll    func(string) error
	chmod        func(string, fs.FileMode) error
	chtimes      func(string, time.Time, time.Time) error
	symlink      func(string, string) error
	commitAbsent func(string, string) error
}

func (b extractBackend) Lstat(name string) (fs.FileInfo, error) {
	if b.lstat != nil {
		return b.lstat(name)
	}
	return b.Backend.Lstat(name)
}
func (b extractBackend) ReadDir(name string) ([]fs.DirEntry, error) {
	if b.readDir != nil {
		return b.readDir(name)
	}
	return b.Backend.ReadDir(name)
}
func (b extractBackend) Open(name string) (io.ReadCloser, error) {
	if b.open != nil {
		return b.open(name)
	}
	return b.Backend.Open(name)
}
func (b extractBackend) Create(name string, mode fs.FileMode) (fsx.Writable, error) {
	if b.create != nil {
		return b.create(name, mode)
	}
	return b.Backend.Create(name, mode)
}
func (b extractBackend) MkdirAll(name string, mode fs.FileMode) error {
	if b.mkdirAll != nil {
		return b.mkdirAll(name, mode)
	}
	return b.Backend.MkdirAll(name, mode)
}
func (b extractBackend) RemoveAll(name string) error {
	if b.removeAll != nil {
		return b.removeAll(name)
	}
	return b.Backend.RemoveAll(name)
}
func (b extractBackend) Chmod(name string, mode fs.FileMode) error {
	if b.chmod != nil {
		return b.chmod(name, mode)
	}
	return b.Backend.Chmod(name, mode)
}
func (b extractBackend) Chtimes(name string, atime, mtime time.Time) error {
	if b.chtimes != nil {
		return b.chtimes(name, atime, mtime)
	}
	return b.Backend.Chtimes(name, atime, mtime)
}
func (b extractBackend) Symlink(target, name string) error {
	if b.symlink != nil {
		return b.symlink(target, name)
	}
	return b.Backend.Symlink(target, name)
}
func (b extractBackend) CommitAbsent(stage, destination string) error {
	if b.commitAbsent != nil {
		return b.commitAbsent(stage, destination)
	}
	return b.Backend.CommitAbsent(stage, destination)
}

type shortWriter struct{ bytes.Buffer }

func (w *shortWriter) Write(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}
	return w.Buffer.Write(data[:1])
}

type faultWritable struct {
	bytes.Buffer
	writeErr error
	syncErr  error
	closeErr error
}

func (w *faultWritable) Write(data []byte) (int, error) {
	if w.writeErr != nil {
		return 0, w.writeErr
	}
	return w.Buffer.Write(data)
}
func (w *faultWritable) Sync() error  { return w.syncErr }
func (w *faultWritable) Close() error { return w.closeErr }

func TestExtractBackendForwarding(t *testing.T) {
	backend := extractBackend{Backend: fsx.Local{}}
	root := t.TempDir()
	if err := backend.MkdirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := backend.Lstat(root); err != nil {
		t.Fatal(err)
	}
	if _, err := backend.ReadDir(root); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(root, "file")
	writable, err := backend.Create(file, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if err := writable.Close(); err != nil {
		t.Fatal(err)
	}
	reader, err := backend.Open(file)
	if err != nil {
		t.Fatal(err)
	}
	_ = reader.Close()
	if err := backend.Chmod(file, 0o640); err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	if err := backend.Chtimes(file, now, now); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "link")
	if err := backend.Symlink("file", link); err != nil {
		t.Fatal(err)
	}
	staged := filepath.Join(root, "staged")
	if err := os.WriteFile(staged, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := backend.CommitAbsent(staged, filepath.Join(root, "final")); err != nil {
		t.Fatal(err)
	}
	if err := backend.RemoveAll(link); err != nil {
		t.Fatal(err)
	}
	if entries, err := backend.ReadDir(root); err != nil || reflect.ValueOf(entries).IsNil() {
		t.Fatalf("entries=%v err=%v", entries, err)
	}
}
