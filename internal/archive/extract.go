package archive

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
	"time"

	"github.com/iwonz/courier/internal/fsx"
	"github.com/iwonz/courier/internal/progress"
	"github.com/iwonz/courier/internal/safety"
	"github.com/iwonz/courier/internal/selection"
)

const (
	MaxExtractionEntries = 100_000
	MaxExtractionDepth   = 64
	MaxExpansionRatio    = int64(100)
	DefaultMaxExtracted  = int64(100 << 30)
	extractionBufferSize = 128 * 1024
	maxArchivePathBytes  = 4096
)

var (
	// ErrCodec identifies registry and format resolution failures.
	ErrCodec = errors.New("archive codec error")
	// ErrArchiveLimit identifies an archive-bomb safety limit.
	ErrArchiveLimit   = errors.New("archive extraction limit exceeded")
	extractRandomRead = rand.Read
	countArchiveEntry = nextEntryCount
	validatePayload   = validateEntryPayload
	discardPayload    = func(reader io.Reader, buffer []byte) (int64, error) { return io.CopyBuffer(io.Discard, reader, buffer) }
)

// Limits contains the configurable expanded-size policy. Fixed count, depth,
// and ratio limits always remain active.
type Limits struct {
	MaxBytes   int64
	Unlimited  bool
	Configured bool
}

// Entry is a validated normalized archive object.
type Entry struct {
	Name     string
	Mode     fs.FileMode
	Size     int64
	ModTime  time.Time
	Linkname string
	Selected bool
}

// Summary describes a complete archive inspection or extraction pass.
type Summary struct {
	Entries       int
	ExpandedBytes int64
	SelectedBytes int64
}

// InspectRequest describes a side-effect-free codec inspection.
type InspectRequest struct {
	SourceFS       fsx.Backend
	SourcePath     string
	CompressedSize int64
	Selector       selection.Selector
	Limits         Limits
	Visit          func(Entry) error
}

// CodecExtractRequest describes extraction into an already private stage.
type CodecExtractRequest struct {
	SourceFS       fsx.Backend
	SourcePath     string
	CompressedSize int64
	DestinationFS  fsx.Backend
	StageRoot      string
	Selector       selection.Selector
	Limits         Limits
	Tracker        *progress.Tracker
}

// Codec is an archive format implementation registered by extension.
type Codec interface {
	Name() string
	Extensions() []string
	Inspect(context.Context, InspectRequest) (Summary, error)
	Extract(context.Context, CodecExtractRequest) (Summary, error)
}

// Registry is an immutable extension-to-codec dispatch table.
type Registry struct {
	codecs     []Codec
	extensions map[string]Codec
}

// NewRegistry validates and registers codecs.
func NewRegistry(codecs ...Codec) (Registry, error) {
	result := Registry{codecs: append([]Codec(nil), codecs...), extensions: map[string]Codec{}}
	names := map[string]bool{}
	for _, codec := range codecs {
		if codec == nil || strings.TrimSpace(codec.Name()) == "" || names[codec.Name()] || len(codec.Extensions()) == 0 {
			return Registry{}, fmt.Errorf("%w: invalid or duplicate codec", ErrCodec)
		}
		names[codec.Name()] = true
		for _, extension := range codec.Extensions() {
			extension = strings.ToLower(extension)
			if !strings.HasPrefix(extension, ".") || len(extension) < 2 || result.extensions[extension] != nil {
				return Registry{}, fmt.Errorf("%w: invalid or duplicate extension %q", ErrCodec, extension)
			}
			result.extensions[extension] = codec
		}
	}
	return result, nil
}

// DefaultRegistry returns the shipped tar.gz registry.
func DefaultRegistry() Registry {
	registry, _ := NewRegistry(TarGzipCodec{})
	return registry
}

// Resolve selects a codec by the longest registered filename extension.
func (r Registry) Resolve(name string) (Codec, error) {
	lower := strings.ToLower(name)
	var matched string
	for extension := range r.extensions {
		if strings.HasSuffix(lower, extension) && len(extension) > len(matched) {
			matched = extension
		}
	}
	if matched == "" {
		return nil, fmt.Errorf("%w: unsupported archive %q", ErrCodec, name)
	}
	return r.extensions[matched], nil
}

// ExtractionRequest is a complete registry-backed extraction operation.
type ExtractionRequest struct {
	SourceFS        fsx.Backend
	SourcePath      string
	SourceName      string
	DestinationFS   fsx.Backend
	DestinationRoot string
	Selector        selection.Selector
	Limits          Limits
	Progress        progress.Sink
}

// ExtractionResult reports selected bytes and committed top-level objects.
type ExtractionResult struct {
	Bytes     int64
	Entries   int
	Committed int
	Elapsed   time.Duration
}

// ExtractionError reports the phase and bytes already committed to the final
// extraction root when an extraction operation fails.
type ExtractionError struct {
	Stage     progress.Stage
	Confirmed int64
	Cause     error
}

func (e *ExtractionError) Error() string { return e.Cause.Error() }
func (e *ExtractionError) Unwrap() error { return e.Cause }

// Extract inspects, stages, and commits an archive into an extraction root.
func (r Registry) Extract(ctx context.Context, request ExtractionRequest) (result ExtractionResult, resultErr error) {
	failureStage := progress.StagePreflight
	defer func() {
		if resultErr != nil {
			var typed *ExtractionError
			if !errors.As(resultErr, &typed) {
				resultErr = &ExtractionError{Stage: failureStage, Confirmed: result.Bytes, Cause: resultErr}
			}
		}
	}()
	if request.SourceFS == nil || request.DestinationFS == nil || request.SourcePath == "" || request.SourceName == "" || request.DestinationRoot == "" {
		return ExtractionResult{}, fmt.Errorf("%w: extraction paths and backends are required", ErrCodec)
	}
	codec, err := r.Resolve(request.SourceName)
	if err != nil {
		return ExtractionResult{}, err
	}
	sourceInfo, err := request.SourceFS.Lstat(request.SourcePath)
	if err != nil {
		return ExtractionResult{}, err
	}
	if !sourceInfo.Mode().IsRegular() {
		return ExtractionResult{}, fmt.Errorf("%w: extraction source must be a regular file", ErrCodec)
	}
	selector := request.Selector
	if selector == nil {
		all := selection.All()
		selector = all
	}
	rootExists := false
	if info, statErr := request.DestinationFS.Lstat(request.DestinationRoot); statErr == nil {
		if !info.IsDir() {
			return ExtractionResult{}, fmt.Errorf("extraction root is not a directory")
		}
		rootExists = true
	} else if !errors.Is(statErr, fs.ErrNotExist) {
		return ExtractionResult{}, statErr
	}
	topLevelBytes := map[string]int64{}
	checkedTopLevels := map[string]bool{}
	inspect := InspectRequest{
		SourceFS: request.SourceFS, SourcePath: request.SourcePath, CompressedSize: sourceInfo.Size(), Selector: selector, Limits: request.Limits,
		Visit: func(entry Entry) error {
			if !entry.Selected {
				return nil
			}
			topLevel := strings.Split(entry.Name, "/")[0]
			topLevelBytes[topLevel] += entry.Size
			if !rootExists || checkedTopLevels[topLevel] {
				return nil
			}
			checkedTopLevels[topLevel] = true
			name := request.DestinationFS.Join(request.DestinationRoot, topLevel)
			if _, collisionErr := request.DestinationFS.Lstat(name); collisionErr == nil {
				return fmt.Errorf("%w: %q", fsx.ErrDestinationExists, name)
			} else if !errors.Is(collisionErr, fs.ErrNotExist) {
				return collisionErr
			}
			return nil
		},
	}
	summary, err := codec.Inspect(ctx, inspect)
	if err != nil {
		return ExtractionResult{}, err
	}
	failureStage = progress.StageExtract
	token, err := extractionToken()
	if err != nil {
		return ExtractionResult{}, err
	}
	stageRoot := extractionStage(request.DestinationRoot, token)
	cleanupStage := true
	defer func() {
		if cleanupStage {
			resultErr = errors.Join(resultErr, request.DestinationFS.RemoveAll(stageRoot))
		}
	}()
	if err := request.DestinationFS.RemoveAll(stageRoot); err != nil {
		return ExtractionResult{}, err
	}
	if err := request.DestinationFS.MkdirAll(request.DestinationFS.Dir(stageRoot), 0o700); err != nil {
		return ExtractionResult{}, err
	}
	if err := request.DestinationFS.MkdirAll(stageRoot, 0o700); err != nil {
		return ExtractionResult{}, err
	}
	tracker := progress.New(summary.SelectedBytes, nil, request.Progress)
	tracker.Stage(progress.StageExtract)
	extracted, err := codec.Extract(ctx, CodecExtractRequest{
		SourceFS: request.SourceFS, SourcePath: request.SourcePath, CompressedSize: sourceInfo.Size(), DestinationFS: request.DestinationFS,
		StageRoot: stageRoot, Selector: selector, Limits: request.Limits, Tracker: tracker,
	})
	if err != nil {
		return ExtractionResult{}, err
	}
	tracker.Stage(progress.StageCommit)
	failureStage = progress.StageCommit
	committed := 0
	if !rootExists {
		if err := request.DestinationFS.CommitAbsent(stageRoot, request.DestinationRoot); err != nil {
			return ExtractionResult{}, err
		}
		cleanupStage = false
		committed = 1
		result.Bytes = summary.SelectedBytes
	} else {
		entries, err := request.DestinationFS.ReadDir(stageRoot)
		if err != nil {
			return ExtractionResult{}, err
		}
		for _, entry := range entries {
			finalName := request.DestinationFS.Join(request.DestinationRoot, entry.Name())
			if _, collisionErr := request.DestinationFS.Lstat(finalName); collisionErr == nil {
				return ExtractionResult{}, fmt.Errorf("%w: %q", fsx.ErrDestinationExists, finalName)
			} else if !errors.Is(collisionErr, fs.ErrNotExist) {
				return ExtractionResult{}, collisionErr
			}
		}
		for _, entry := range entries {
			stageName := request.DestinationFS.Join(stageRoot, entry.Name())
			finalName := request.DestinationFS.Join(request.DestinationRoot, entry.Name())
			if err := request.DestinationFS.CommitAbsent(stageName, finalName); err != nil {
				return ExtractionResult{Bytes: result.Bytes, Entries: extracted.Entries, Committed: committed, Elapsed: tracker.Snapshot().Elapsed}, err
			}
			committed++
			result.Bytes += topLevelBytes[entry.Name()]
		}
	}
	tracker.Stage(progress.StageCleanup)
	failureStage = progress.StageCleanup
	if cleanupStage {
		if err := request.DestinationFS.RemoveAll(stageRoot); err != nil {
			return ExtractionResult{}, err
		}
		cleanupStage = false
	}
	tracker.Stage(progress.StageComplete)
	return ExtractionResult{Bytes: result.Bytes, Entries: extracted.Entries, Committed: committed, Elapsed: tracker.Snapshot().Elapsed}, nil
}

// TarGzipCodec implements the shipped POSIX tar over gzip codec.
type TarGzipCodec struct{}

func (TarGzipCodec) Name() string         { return "tar.gz" }
func (TarGzipCodec) Extensions() []string { return []string{".tar.gz"} }

func (TarGzipCodec) Inspect(ctx context.Context, request InspectRequest) (Summary, error) {
	return walkTarGzip(ctx, request.SourceFS, request.SourcePath, request.CompressedSize, request.Selector, request.Limits, func(entry Entry, _ io.Reader) error {
		if request.Visit != nil {
			return request.Visit(entry)
		}
		return nil
	})
}

func (TarGzipCodec) Extract(ctx context.Context, request CodecExtractRequest) (Summary, error) {
	directories := make([]Entry, 0)
	summary, err := walkTarGzip(ctx, request.SourceFS, request.SourcePath, request.CompressedSize, request.Selector, request.Limits, func(entry Entry, reader io.Reader) error {
		if !entry.Selected {
			return nil
		}
		destination := joinArchivePath(request.DestinationFS, request.StageRoot, entry.Name)
		if entry.Mode.IsDir() {
			if err := request.DestinationFS.MkdirAll(destination, 0o700); err != nil {
				return err
			}
			directories = append(directories, entry)
			return nil
		}
		if err := request.DestinationFS.MkdirAll(request.DestinationFS.Dir(destination), 0o700); err != nil {
			return err
		}
		if entry.Mode&fs.ModeSymlink != 0 {
			return request.DestinationFS.Symlink(entry.Linkname, destination)
		}
		writer, err := request.DestinationFS.Create(destination, 0o600)
		if err != nil {
			return err
		}
		written, copyErr := io.CopyBuffer(&extractWriter{ctx: ctx, destination: writer, tracker: request.Tracker}, reader, make([]byte, extractionBufferSize))
		syncErr := writer.Sync()
		closeErr := writer.Close()
		if copyErr != nil || syncErr != nil || closeErr != nil || written != entry.Size {
			return errors.Join(copyErr, syncErr, closeErr, sizeError(written, entry.Size))
		}
		if err := request.DestinationFS.Chmod(destination, entry.Mode.Perm()); err != nil {
			return err
		}
		return request.DestinationFS.Chtimes(destination, entry.ModTime, entry.ModTime)
	})
	if err != nil {
		return Summary{}, err
	}
	for index := len(directories) - 1; index >= 0; index-- {
		directory := directories[index]
		destination := joinArchivePath(request.DestinationFS, request.StageRoot, directory.Name)
		if err := request.DestinationFS.Chmod(destination, directory.Mode.Perm()); err != nil {
			return Summary{}, err
		}
		if err := request.DestinationFS.Chtimes(destination, directory.ModTime, directory.ModTime); err != nil {
			return Summary{}, err
		}
	}
	return summary, nil
}

func walkTarGzip(ctx context.Context, backend fsx.Backend, sourcePath string, compressedSize int64, selector selection.Selector, limits Limits, visit func(Entry, io.Reader) error) (summary Summary, resultErr error) {
	if backend == nil || sourcePath == "" || compressedSize < 0 {
		return Summary{}, fmt.Errorf("%w: invalid inspection request", ErrCodec)
	}
	if selector == nil {
		all := selection.All()
		selector = all
	}
	file, err := backend.Open(sourcePath)
	if err != nil {
		return Summary{}, err
	}
	defer func() { resultErr = errors.Join(resultErr, file.Close()) }()
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return Summary{}, err
	}
	defer func() { resultErr = errors.Join(resultErr, gzipReader.Close()) }()
	reader := tar.NewReader(gzipReader)
	seen := map[string]fs.FileMode{}
	requiredDirectories := map[string]bool{}
	excludedDirectories := map[string]bool{}
	buffer := make([]byte, extractionBufferSize)
	for {
		if err := ctx.Err(); err != nil {
			return Summary{}, err
		}
		header, err := reader.Next()
		if errors.Is(err, io.EOF) {
			return summary, nil
		}
		if err != nil {
			return Summary{}, err
		}
		summary.Entries, err = countArchiveEntry(summary.Entries)
		if err != nil {
			return Summary{}, err
		}
		if len(header.Name) > maxArchivePathBytes {
			return Summary{}, fmt.Errorf("%w: archive path is too long", ErrArchiveLimit)
		}
		if _, err := safety.SafeArchiveJoin("archive-root", header.Name); err != nil {
			return Summary{}, err
		}
		name := path.Clean(strings.TrimSuffix(header.Name, "/"))
		if _, duplicate := seen[name]; duplicate {
			return Summary{}, fmt.Errorf("duplicate archive entry %q", name)
		}
		components := strings.Split(name, "/")
		if len(components) > MaxExtractionDepth {
			return Summary{}, fmt.Errorf("%w: entry depth exceeds %d", ErrArchiveLimit, MaxExtractionDepth)
		}
		mode, err := archiveMode(header)
		if err != nil {
			return Summary{}, err
		}
		if requiredDirectories[name] && !mode.IsDir() {
			return Summary{}, fmt.Errorf("archive entry %q conflicts with child entries", name)
		}
		for index := 1; index < len(components); index++ {
			ancestor := strings.Join(components[:index], "/")
			if ancestorMode, exists := seen[ancestor]; exists && !ancestorMode.IsDir() {
				return Summary{}, fmt.Errorf("archive parent %q is not a directory", ancestor)
			}
			requiredDirectories[ancestor] = true
		}
		if mode&fs.ModeSymlink != 0 {
			if err := safety.ValidateArchiveSymlink(name, header.Linkname); err != nil {
				return Summary{}, err
			}
		}
		seen[name] = mode
		if err := validatePayload(name, mode, header.Size); err != nil {
			return Summary{}, err
		}
		summary.ExpandedBytes, err = addExpandedSize(summary.ExpandedBytes, header.Size, compressedSize, limits)
		if err != nil {
			return Summary{}, err
		}
		selected := true
		for index := 1; index <= len(components); index++ {
			if excludedDirectories[strings.Join(components[:index], "/")] {
				selected = false
				break
			}
		}
		if selected {
			selected = selector.Include(name, mode.IsDir())
		}
		if mode.IsDir() && !selected {
			excludedDirectories[name] = true
		}
		entry := Entry{Name: name, Mode: mode, Size: header.Size, ModTime: header.ModTime, Linkname: header.Linkname, Selected: selected}
		if selected {
			summary.SelectedBytes += header.Size
		}
		if visit != nil {
			if err := visit(entry, reader); err != nil {
				return Summary{}, err
			}
		}
		if _, err := discardPayload(&contextReader{ctx: ctx, reader: reader}, buffer); err != nil {
			return Summary{}, err
		}
	}
}

func nextEntryCount(current int) (int, error) {
	current++
	if current > MaxExtractionEntries {
		return 0, fmt.Errorf("%w: more than %d entries", ErrArchiveLimit, MaxExtractionEntries)
	}
	return current, nil
}

func validateEntryPayload(name string, mode fs.FileMode, size int64) error {
	if !mode.IsRegular() && size != 0 {
		return fmt.Errorf("unsupported archive entry %q with non-zero payload", name)
	}
	return nil
}

func addExpandedSize(current, size, compressed int64, limits Limits) (int64, error) {
	if size < 0 || current > int64(^uint64(0)>>1)-size {
		return 0, fmt.Errorf("%w: expanded size overflow", ErrArchiveLimit)
	}
	expanded := current + size
	maximum := limits.MaxBytes
	if !limits.Configured {
		maximum = DefaultMaxExtracted
	}
	if !limits.Unlimited && expanded > maximum {
		return 0, fmt.Errorf("%w: expanded size exceeds configured maximum", ErrArchiveLimit)
	}
	if exceedsExpansionRatio(expanded, compressed) {
		return 0, fmt.Errorf("%w: expansion ratio exceeds %d:1", ErrArchiveLimit, MaxExpansionRatio)
	}
	return expanded, nil
}

func exceedsExpansionRatio(expanded, compressed int64) bool {
	if compressed == 0 {
		return expanded > 0
	}
	if compressed > int64(^uint64(0)>>1)/MaxExpansionRatio {
		return false
	}
	return expanded > compressed*MaxExpansionRatio
}

func archiveMode(header *tar.Header) (fs.FileMode, error) {
	mode := header.FileInfo().Mode()
	switch header.Typeflag {
	case tar.TypeReg, tar.TypeRegA:
		return mode, nil
	case tar.TypeDir:
		return mode | fs.ModeDir, nil
	case tar.TypeSymlink:
		return mode | fs.ModeSymlink, nil
	default:
		return 0, fmt.Errorf("unsupported archive entry %q type %d", header.Name, header.Typeflag)
	}
}

func joinArchivePath(backend fsx.Backend, root, name string) string {
	parts := append([]string{root}, strings.Split(name, "/")...)
	return backend.Join(parts...)
}

func extractionStage(root, token string) string {
	base := strings.TrimRight(root, `/\`)
	if base == "" {
		return root + ".courier-extract-" + token
	}
	return base + ".courier-extract-" + token
}

func extractionToken() (string, error) {
	value := make([]byte, 8)
	if _, err := extractRandomRead(value); err != nil {
		return "", err
	}
	return hex.EncodeToString(value), nil
}

type extractWriter struct {
	ctx         context.Context
	destination io.Writer
	tracker     *progress.Tracker
}

func (w *extractWriter) Write(data []byte) (int, error) {
	if err := w.ctx.Err(); err != nil {
		return 0, err
	}
	written, err := w.destination.Write(data)
	if written > 0 && w.tracker != nil {
		w.tracker.Add(int64(written))
	}
	return written, err
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r *contextReader) Read(data []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.reader.Read(data)
}
