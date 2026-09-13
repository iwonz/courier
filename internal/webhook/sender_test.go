package webhook

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/iwonz/courier/internal/fsx"
	"github.com/iwonz/courier/internal/progress"
	"golang.org/x/time/rate"
)

type doerFunc func(*http.Request) (*http.Response, error)

func (function doerFunc) Do(request *http.Request) (*http.Response, error) { return function(request) }

type backendStub struct {
	fsx.Backend
	lstat func(string) (fs.FileInfo, error)
	open  func(string) (io.ReadCloser, error)
}

func (backend backendStub) Lstat(path string) (fs.FileInfo, error) {
	if backend.lstat != nil {
		return backend.lstat(path)
	}
	return backend.Backend.Lstat(path)
}

func (backend backendStub) Open(path string) (io.ReadCloser, error) {
	if backend.open != nil {
		return backend.open(path)
	}
	return backend.Backend.Open(path)
}

type closeReader struct {
	io.Reader
	err error
}

func (reader closeReader) Close() error { return reader.err }

type responseBody struct{ err error }

func (responseBody) Read([]byte) (int, error) { return 0, io.EOF }
func (body responseBody) Close() error        { return body.err }

type sizedInfo struct {
	fs.FileInfo
	size int64
}

func (info sizedInfo) Size() int64 { return info.size }

func sourceFile(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "résumé.txt")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func validRequest(path, destination string) Request {
	return Request{URL: destination, SourceFS: fsx.Local{}, SourcePath: path, Name: filepath.Base(path), Unlimited: true}
}

func TestSenderSuccessAndBasicAuthentication(t *testing.T) {
	path := sourceFile(t, "payload")
	calls := atomic.Int32{}
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		calls.Add(1)
		username, password, ok := request.BasicAuth()
		if !ok || username != "courier" || password != "secret" || request.Header.Get("Expect") != "100-continue" {
			t.Errorf("unexpected authentication or headers: %q %q %v %v", username, password, ok, request.Header)
		}
		reader, err := request.MultipartReader()
		if err != nil {
			t.Error(err)
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		part, err := reader.NextPart()
		if err != nil || part.FormName() != "file" || part.FileName() != filepath.Base(path) {
			t.Errorf("unexpected part: %+v %v", part, err)
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		data, readErr := io.ReadAll(part)
		_, nextErr := reader.NextPart()
		if readErr != nil || !errors.Is(nextErr, io.EOF) || string(data) != "payload" {
			t.Errorf("data=%q read=%v next=%v", data, readErr, nextErr)
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		response.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()
	request := validRequest(path, server.URL+"/hook")
	request.Username = "courier"
	request.Password = []byte("secret")
	request.Rate = 1 << 20
	request.Unlimited = false
	var events []progress.Event
	request.Progress = func(event progress.Event) { events = append(events, event) }
	result, err := (Sender{}).Send(context.Background(), request)
	if err != nil || result.StatusCode != http.StatusAccepted || result.Bytes != 7 || calls.Load() != 1 {
		t.Fatalf("result=%+v calls=%d err=%v", result, calls.Load(), err)
	}
	if events[len(events)-1].Stage != progress.StageComplete || events[len(events)-1].Current != 7 {
		t.Fatalf("events=%+v", events)
	}
}

func TestSenderDoesNotFollowRedirectOrRetry(t *testing.T) {
	path := sourceFile(t, "payload")
	targetCalls := atomic.Int32{}
	target := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { targetCalls.Add(1) }))
	defer target.Close()
	sourceCalls := atomic.Int32{}
	source := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		sourceCalls.Add(1)
		http.Redirect(response, request, target.URL, http.StatusTemporaryRedirect)
	}))
	defer source.Close()
	result, err := (Sender{}).Send(context.Background(), validRequest(path, source.URL))
	var rejected *RejectedError
	if !errors.As(err, &rejected) || rejected.StatusCode != http.StatusTemporaryRedirect || result.StatusCode != http.StatusTemporaryRedirect || sourceCalls.Load() != 1 || targetCalls.Load() != 0 {
		t.Fatalf("result=%+v source=%d target=%d err=%v", result, sourceCalls.Load(), targetCalls.Load(), err)
	}
	if rejected.Error() != "webhook receiver returned HTTP 307" {
		t.Fatal(rejected)
	}
}

func TestSenderClassifiesKnownAndUnknownFailures(t *testing.T) {
	path := sourceFile(t, "payload")
	want := errors.New("transport")
	unknownSender := Sender{Client: doerFunc(func(request *http.Request) (*http.Response, error) {
		_, _ = io.Copy(io.Discard, request.Body)
		return nil, want
	})}
	result, err := unknownSender.Send(context.Background(), validRequest(path, "https://example.test/hook"))
	var unknown *UnknownOutcomeError
	if !errors.As(err, &unknown) || !errors.Is(err, want) || unknown.Sent != 7 || result.Bytes != 7 || !strings.Contains(unknown.Error(), "unknown after sending 7 bytes") {
		t.Fatalf("result=%+v unknown=%+v err=%v", result, unknown, err)
	}

	zeroSender := Sender{Client: doerFunc(func(*http.Request) (*http.Response, error) { return nil, want })}
	result, err = zeroSender.Send(context.Background(), validRequest(path, "https://example.test/hook"))
	if !errors.Is(err, want) || result.Bytes != 0 || errors.As(err, &unknown) {
		t.Fatalf("zero result=%+v err=%v", result, err)
	}
	emptyPath := sourceFile(t, "")
	emptySender := Sender{Client: doerFunc(func(request *http.Request) (*http.Response, error) {
		_, _ = io.Copy(io.Discard, request.Body)
		return nil, want
	})}
	result, err = emptySender.Send(context.Background(), validRequest(emptyPath, "https://example.test/hook"))
	unknown = nil
	if !errors.As(err, &unknown) || unknown.Sent != 0 || result.Bytes != 0 {
		t.Fatalf("empty unknown result=%+v err=%v", result, err)
	}

	closeErr := errors.New("close")
	knownSender := Sender{Client: doerFunc(func(request *http.Request) (*http.Response, error) {
		_, _ = io.Copy(io.Discard, request.Body)
		return &http.Response{StatusCode: http.StatusBadRequest, Body: responseBody{err: closeErr}}, nil
	})}
	result, err = knownSender.Send(context.Background(), validRequest(path, "https://example.test/hook"))
	var rejected *RejectedError
	if !errors.As(err, &rejected) || !errors.Is(err, closeErr) || result.StatusCode != http.StatusBadRequest {
		t.Fatalf("known result=%+v err=%v", result, err)
	}

	successClose := Sender{Client: doerFunc(func(request *http.Request) (*http.Response, error) {
		_, _ = io.Copy(io.Discard, request.Body)
		return &http.Response{StatusCode: http.StatusNoContent, Body: responseBody{err: closeErr}}, nil
	})}
	if _, err := successClose.Send(context.Background(), validRequest(path, "https://example.test/hook")); !errors.Is(err, closeErr) {
		t.Fatalf("response close error=%v", err)
	}

	info, statErr := os.Stat(path)
	if statErr != nil {
		t.Fatal(statErr)
	}
	changed := validRequest(path, "https://example.test/hook")
	changed.SourceFS = backendStub{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) {
		return sizedInfo{FileInfo: info, size: info.Size() + 1}, nil
	}}
	changedSender := Sender{Client: doerFunc(func(request *http.Request) (*http.Response, error) {
		_, _ = io.Copy(io.Discard, request.Body)
		return &http.Response{StatusCode: http.StatusOK, Body: responseBody{}}, nil
	})}
	if result, err := changedSender.Send(context.Background(), changed); err == nil || result.Bytes != info.Size() {
		t.Fatalf("changed source result=%+v err=%v", result, err)
	}
}

func TestSenderValidationAndSourceFailures(t *testing.T) {
	path := sourceFile(t, "payload")
	valid := validRequest(path, "https://example.test/hook")
	for _, test := range []struct {
		name   string
		mutate func(*Request)
	}{
		{"nil backend", func(value *Request) { value.SourceFS = nil }},
		{"empty path", func(value *Request) { value.SourcePath = " " }},
		{"empty name", func(value *Request) { value.Name = "" }},
		{"unsafe name", func(value *Request) { value.Name = "a/b" }},
		{"username only", func(value *Request) { value.Username = "user" }},
		{"password only", func(value *Request) { value.Password = []byte("secret") }},
		{"zero rate", func(value *Request) { value.Unlimited = false }},
		{"relative URL", func(value *Request) { value.URL = "/hook" }},
		{"credential URL", func(value *Request) { value.URL = "https://user@example.test/hook" }},
		{"fragment URL", func(value *Request) { value.URL = "https://example.test/hook#fragment" }},
		{"wrong scheme", func(value *Request) { value.URL = "ftp://example.test/hook" }},
		{"malformed URL", func(value *Request) { value.URL = "https://example.test/%zz" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			request := valid
			test.mutate(&request)
			if _, err := (Sender{Client: doerFunc(nil)}).Send(context.Background(), request); err == nil {
				t.Fatal("invalid request accepted")
			}
		})
	}
	if _, err := (Sender{}).Send(nil, valid); err == nil {
		t.Fatal("nil context accepted")
	}
	missing := valid
	missing.SourcePath = filepath.Join(t.TempDir(), "missing")
	if _, err := (Sender{}).Send(context.Background(), missing); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("missing source=%v", err)
	}
	directory := valid
	directory.SourcePath = t.TempDir()
	if _, err := (Sender{}).Send(context.Background(), directory); err == nil {
		t.Fatal("directory source accepted")
	}
	want := errors.New("open")
	openFailure := valid
	openFailure.SourceFS = backendStub{Backend: fsx.Local{}, open: func(string) (io.ReadCloser, error) { return nil, want }}
	if _, err := (Sender{}).Send(context.Background(), openFailure); !errors.Is(err, want) {
		t.Fatalf("open failure=%v", err)
	}
}

func TestMultipartAndPayloadFailures(t *testing.T) {
	pipeReader, pipeWriter := io.Pipe()
	_ = pipeReader.Close()
	if err := writeMultipart(context.Background(), multipart.NewWriter(pipeWriter), strings.NewReader("x"), "x", 1, 0, true, progress.New(1, nil, nil), &atomic.Int64{}); err == nil {
		t.Fatal("multipart header failure ignored")
	}

	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	consumed := &atomic.Int64{}
	reader := &payloadReader{
		ctx: canceled, source: strings.NewReader("ab"), limiter: rate.NewLimiter(1, 1), burst: 1,
		tracker: progress.New(2, nil, nil), consumed: consumed,
	}
	if count, err := reader.Read(make([]byte, 2)); err == nil || count != 0 || consumed.Load() != 0 {
		t.Fatalf("canceled rate read=%d consumed=%d err=%v", count, consumed.Load(), err)
	}
	reader = &payloadReader{ctx: context.Background(), source: strings.NewReader("x"), tracker: progress.New(1, nil, nil), consumed: consumed}
	if count, err := reader.Read(make([]byte, 2)); err != nil || count != 1 || consumed.Load() != 1 {
		t.Fatalf("unlimited read=%d consumed=%d err=%v", count, consumed.Load(), err)
	}

	client := DefaultClient()
	transport, ok := client.Transport.(*http.Transport)
	if !ok || !transport.DisableKeepAlives || transport.MaxResponseHeaderBytes != 1<<20 || client.CheckRedirect(nil, nil) != http.ErrUseLastResponse {
		t.Fatalf("default client=%+v", client)
	}
}
