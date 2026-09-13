package webdelivery

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/endpoint"
	"github.com/iwonz/courier/internal/fsx"
	"github.com/iwonz/courier/internal/policy"
	"github.com/iwonz/courier/internal/selection"
)

type selectorFunc func(string, bool) bool

func (function selectorFunc) Include(name string, directory bool) bool {
	return function(name, directory)
}

type overrideBackend struct {
	fsx.Backend
	lstat   func(string) (fs.FileInfo, error)
	readDir func(string) ([]fs.DirEntry, error)
	open    func(string) (io.ReadCloser, error)
	create  func(string, fs.FileMode) (fsx.Writable, error)
	commit  func(string, string) error
	remove  func(string) error
}

func (backend overrideBackend) Lstat(name string) (fs.FileInfo, error) {
	if backend.lstat != nil {
		return backend.lstat(name)
	}
	return backend.Backend.Lstat(name)
}
func (backend overrideBackend) ReadDir(name string) ([]fs.DirEntry, error) {
	if backend.readDir != nil {
		return backend.readDir(name)
	}
	return backend.Backend.ReadDir(name)
}
func (backend overrideBackend) Open(name string) (io.ReadCloser, error) {
	if backend.open != nil {
		return backend.open(name)
	}
	return backend.Backend.Open(name)
}
func (backend overrideBackend) Create(name string, mode fs.FileMode) (fsx.Writable, error) {
	if backend.create != nil {
		return backend.create(name, mode)
	}
	return backend.Backend.Create(name, mode)
}
func (backend overrideBackend) CommitAbsent(stage, final string) error {
	if backend.commit != nil {
		return backend.commit(stage, final)
	}
	return backend.Backend.CommitAbsent(stage, final)
}
func (backend overrideBackend) RemoveAll(name string) error {
	if backend.remove != nil {
		return backend.remove(name)
	}
	return backend.Backend.RemoveAll(name)
}

type errorWritable struct {
	writeErr error
	syncErr  error
	closeErr error
}

func (writer *errorWritable) Write(data []byte) (int, error) {
	if writer.writeErr != nil {
		return 0, writer.writeErr
	}
	return len(data), nil
}
func (writer *errorWritable) Sync() error  { return writer.syncErr }
func (writer *errorWritable) Close() error { return writer.closeErr }

type errorReadCloser struct {
	reader   io.Reader
	closeErr error
}

func (reader errorReadCloser) Read(data []byte) (int, error) { return reader.reader.Read(data) }
func (reader errorReadCloser) Close() error                  { return reader.closeErr }

type errorListener struct{ err error }

func (listener errorListener) Accept() (net.Conn, error) { return nil, listener.err }
func (listener errorListener) Close() error              { return nil }
func (listener errorListener) Addr() net.Addr            { return peerAddress("127.0.0.1:1") }

type errorResponse struct{ header http.Header }

func (response *errorResponse) Header() http.Header {
	if response.header == nil {
		response.header = make(http.Header)
	}
	return response.header
}
func (*errorResponse) WriteHeader(int)           {}
func (*errorResponse) Write([]byte) (int, error) { return 0, errors.New("response write") }

func TestHostRegistrationFailuresAndCloseErrors(t *testing.T) {
	root := t.TempDir()
	host := newTestHost(t, nil)
	definition := webDefinition(delivery.RoutePathToWeb, root)
	record := webRecord(testDeliveryID, delivery.RoutePathToWeb, delivery.DefaultPolicy())

	specialBackend := overrideBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return fakeInfo{mode: fs.ModeSocket}, nil }}
	host.open = func(_ context.Context, _ endpoint.Endpoint, _ EndpointRuntime) (*Resource, error) {
		return &Resource{Backend: specialBackend, Path: root}, nil
	}
	if err := host.Register(context.Background(), record, marshalDefinition(t, definition)); err == nil {
		t.Fatal("expected special source rejection")
	}
	host.open = localOpen
	missingCredentials := record
	missingCredentials.Policy.Auth = delivery.AuthBasic
	if err := host.Register(context.Background(), missingCredentials, marshalDefinition(t, definition)); err == nil {
		t.Fatal("expected policy engine credential error")
	}
	archiveDefinition := definition
	archiveDefinition.Archive = true
	host.archive = func(context.Context, *Resource, string, selection.Selector) (*Resource, error) {
		return nil, errors.New("archive")
	}
	if err := host.Register(context.Background(), record, marshalDefinition(t, archiveDefinition)); err == nil {
		t.Fatal("archive preparation error ignored")
	}
	preparedClosed := atomic.Int32{}
	host.archive = func(context.Context, *Resource, string, selection.Selector) (*Resource, error) {
		return &Resource{Close: func() error { preparedClosed.Add(1); return nil }}, nil
	}
	if err := host.Register(context.Background(), record, marshalDefinition(t, archiveDefinition)); err == nil || preparedClosed.Load() != 1 {
		t.Fatalf("incomplete archive accepted: %v closed=%d", err, preparedClosed.Load())
	}
	host.archive = func(context.Context, *Resource, string, selection.Selector) (*Resource, error) {
		backend := overrideBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return nil, errors.New("prepared stat") }}
		return &Resource{Backend: backend, Path: "prepared", Close: func() error { preparedClosed.Add(1); return nil }}, nil
	}
	if err := host.Register(context.Background(), record, marshalDefinition(t, archiveDefinition)); err == nil {
		t.Fatal("prepared archive stat error ignored")
	}
	host.archive = func(context.Context, *Resource, string, selection.Selector) (*Resource, error) {
		backend := overrideBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return fakeInfo{mode: fs.ModeSocket}, nil }}
		return &Resource{Backend: backend, Path: "prepared", Close: func() error { preparedClosed.Add(1); return nil }}, nil
	}
	if err := host.Register(context.Background(), record, marshalDefinition(t, archiveDefinition)); err == nil || preparedClosed.Load() != 3 {
		t.Fatalf("special prepared archive accepted: %v closed=%d", err, preparedClosed.Load())
	}
	host.archive = func(context.Context, *Resource, string, selection.Selector) (*Resource, error) {
		return &Resource{
			Backend: overrideBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return fakeInfo{mode: 0o600}, nil }},
			Path:    root, Close: func() error { preparedClosed.Add(1); return nil },
			Temps: []delivery.OwnedTemp{{ID: "invalid", Location: delivery.TempLocal, Path: "invalid"}},
		}, nil
	}
	if err := host.Register(context.Background(), record, marshalDefinition(t, archiveDefinition)); err == nil || preparedClosed.Load() != 4 {
		t.Fatalf("invalid archive temporary accepted: %v closed=%d", err, preparedClosed.Load())
	}

	closeHost := newTestHost(t, nil)
	closedErr := errors.New("close")
	closeHost.open = func(_ context.Context, value endpoint.Endpoint, _ EndpointRuntime) (*Resource, error) {
		return &Resource{Endpoint: value, Backend: fsx.Local{}, Path: root, Close: func() error { return closedErr }}, nil
	}
	if err := closeHost.Register(context.Background(), record, marshalDefinition(t, definition)); err != nil {
		t.Fatal(err)
	}
	hostedToClose := closeHost.byID[record.ID]
	if err := closeHost.Stop(context.Background(), record.ID); !errors.Is(err, closedErr) {
		t.Fatalf("expected close error: %v", err)
	}
	if err := hostedToClose.close(); !errors.Is(err, closedErr) {
		t.Fatalf("idempotent close changed result: %v", err)
	}
	if err := host.Serve(errorListener{err: errors.New("accept")}); err == nil {
		t.Fatal("expected serve error")
	}
	_ = peerAddress("127.0.0.1:1").Network()
}

func TestHTTPAuthenticationAndLoginFailures(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "file"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	host := newTestHost(t, nil)
	configured := delivery.DefaultPolicy()
	configured.Auth = delivery.AuthPassword
	configured.AuthAttempts = 1
	configured.AuthFailAction = delivery.AuthFailStop
	definition := webDefinition(delivery.RoutePathToWeb, root)
	definition.Credentials.Password = []byte("secret")
	record := webRecord(testDeliveryID, delivery.RoutePathToWeb, configured)
	if err := host.Register(context.Background(), record, marshalDefinition(t, definition)); err != nil {
		t.Fatal(err)
	}
	base := "/d/" + definition.Token
	if response := perform(host, http.MethodPost, "/d/unknown/api/v1/session", strings.NewReader(`{}`), nil); response.Code != http.StatusNotFound {
		t.Fatal(response.Code)
	}

	openDefinition := webDefinition(delivery.RoutePathToWeb, root)
	openDefinition.Token = fixedToken(9)
	openID := delivery.ID("00000000-0000-4000-8000-000000000039")
	if err := host.Register(context.Background(), webRecord(openID, delivery.RoutePathToWeb, delivery.DefaultPolicy()), marshalDefinition(t, openDefinition)); err != nil {
		t.Fatal(err)
	}
	if response := perform(host, http.MethodPost, "/d/"+openDefinition.Token+"/api/v1/session", strings.NewReader(`{}`), nil); response.Code != http.StatusNotFound {
		t.Fatal(response.Code)
	}
	if response := perform(host, http.MethodPost, base+"/api/v1/session", strings.NewReader(`{`), nil); response.Code != http.StatusBadRequest {
		t.Fatal(response.Code)
	}
	if response := perform(host, http.MethodPost, base+"/api/v1/session", strings.NewReader(`{"password":"secret"}{}`), nil); response.Code != http.StatusBadRequest {
		t.Fatal(response.Code)
	}
	if response := perform(host, http.MethodPost, base+"/api/v1/session", strings.NewReader(`{"password":"secret"}`), map[string]string{"Origin": "https://attacker.example"}); response.Code != http.StatusForbidden {
		t.Fatalf("cross-origin login=%d", response.Code)
	}
	response := perform(host, http.MethodPost, base+"/api/v1/session", strings.NewReader(`{"password":"wrong"}`), nil)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("wrong login=%d", response.Code)
	}
	if response := perform(host, http.MethodGet, base+"/", nil, nil); response.Code != http.StatusNotFound {
		t.Fatalf("stop must revoke route: %d", response.Code)
	}

	allowed := delivery.DefaultPolicy()
	allowed.AllowIP = []string{"198.51.100.0/24"}
	allowDefinition := webDefinition(delivery.RoutePathToWeb, root)
	allowDefinition.Token = fixedToken(10)
	allowID := delivery.ID("00000000-0000-4000-8000-000000000040")
	if err := host.Register(context.Background(), webRecord(allowID, delivery.RoutePathToWeb, allowed), marshalDefinition(t, allowDefinition)); err != nil {
		t.Fatal(err)
	}
	if response := perform(host, http.MethodGet, "/d/"+allowDefinition.Token+"/", nil, nil); response.Code != http.StatusForbidden {
		t.Fatalf("peer allowlist=%d", response.Code)
	}

	csrfPolicy := delivery.DefaultPolicy()
	csrfPolicy.Auth = delivery.AuthPassword
	csrfDefinition := webDefinition(delivery.RouteWebToPath, root)
	csrfDefinition.Token = fixedToken(11)
	csrfDefinition.Credentials.Password = []byte("secret")
	csrfID := delivery.ID("00000000-0000-4000-8000-000000000041")
	if err := host.Register(context.Background(), webRecord(csrfID, delivery.RouteWebToPath, csrfPolicy), marshalDefinition(t, csrfDefinition)); err != nil {
		t.Fatal(err)
	}
	login := perform(host, http.MethodPost, "/d/"+csrfDefinition.Token+"/api/v1/session", strings.NewReader(`{"password":"secret"}`), nil)
	cookie := login.Result().Cookies()[0]
	request := httptest.NewRequest(http.MethodPost, "/d/"+csrfDefinition.Token+"/api/v1/upload", strings.NewReader(""))
	request.RemoteAddr = "192.0.2.10:1"
	request.AddCookie(cookie)
	result := httptest.NewRecorder()
	host.routes().ServeHTTP(result, request)
	if result.Code != http.StatusForbidden {
		t.Fatalf("CSRF=%d", result.Code)
	}

	noUIPassword := csrfDefinition
	noUIPassword.Token = fixedToken(12)
	noUIPassword.NoUI = true
	noUIID := delivery.ID("00000000-0000-4000-8000-000000000042")
	if err := host.Register(context.Background(), webRecord(noUIID, delivery.RouteWebToPath, csrfPolicy), marshalDefinition(t, noUIPassword)); err != nil {
		t.Fatal(err)
	}
	if response := perform(host, http.MethodGet, "/d/"+noUIPassword.Token+"/", nil, nil); response.Code != http.StatusUnauthorized {
		t.Fatalf("no-ui password root=%d", response.Code)
	}
}

func TestMetadataAndDownloadFailures(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "file.txt")
	if err := os.WriteFile(file, []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("file.txt", filepath.Join(root, "link")); err != nil {
		t.Fatal(err)
	}
	host := newTestHost(t, nil)
	definition := webDefinition(delivery.RoutePathToWeb, root)
	record := webRecord(testDeliveryID, delivery.RoutePathToWeb, delivery.DefaultPolicy())
	if err := host.Register(context.Background(), record, marshalDefinition(t, definition)); err != nil {
		t.Fatal(err)
	}
	base := "/d/" + definition.Token
	uploadDefinition := webDefinition(delivery.RouteWebToPath, root)
	uploadDefinition.Token = fixedToken(13)
	uploadID := delivery.ID("00000000-0000-4000-8000-000000000043")
	if err := host.Register(context.Background(), webRecord(uploadID, delivery.RouteWebToPath, delivery.DefaultPolicy()), marshalDefinition(t, uploadDefinition)); err != nil {
		t.Fatal(err)
	}
	if response := perform(host, http.MethodGet, "/d/"+uploadDefinition.Token+"/api/v1/meta", nil, nil); response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"type":"upload"`) || strings.Contains(response.Body.String(), root) {
		t.Fatalf("upload metadata=%d %q", response.Code, response.Body.String())
	}
	if response := perform(host, http.MethodGet, base+"/api/v1/meta?path=missing", nil, nil); response.Code != http.StatusNotFound {
		t.Fatal(response.Code)
	}
	if response := perform(host, http.MethodGet, base+"/api/v1/meta?path=link", nil, nil); response.Code != http.StatusNotFound {
		t.Fatal(response.Code)
	}
	if response := perform(host, http.MethodGet, base+"/api/v1/meta?path=%5Cescape", nil, nil); response.Code != http.StatusBadRequest {
		t.Fatal(response.Code)
	}

	hosted := host.byID[record.ID]
	originalBackend := hosted.resource.Backend
	hosted.resource.Backend = overrideBackend{Backend: originalBackend, readDir: func(string) ([]fs.DirEntry, error) { return nil, errors.New("read") }}
	if response := perform(host, http.MethodGet, base+"/api/v1/meta", nil, nil); response.Code != http.StatusInternalServerError {
		t.Fatal(response.Code)
	}
	hosted.resource.Backend = originalBackend
	hosted.selector = selectorFunc(func(name string, _ bool) bool { return name != "file.txt" })
	if response := perform(host, http.MethodGet, base+"/api/v1/meta", nil, nil); response.Code != http.StatusOK || strings.Contains(response.Body.String(), "file.txt") {
		t.Fatalf("selection=%d %q", response.Code, response.Body.String())
	}
	if _, _, err := inspectRequestedPath(hosted, "file.txt"); err == nil {
		t.Fatal("expected excluded path")
	}
	hosted.selector = selection.All()
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "secret"), []byte("secret"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "outside")); err != nil {
		t.Fatal(err)
	}
	if response := perform(host, http.MethodGet, base+"/api/v1/meta?path=outside%2Fsecret", nil, nil); response.Code != http.StatusNotFound || strings.Contains(response.Body.String(), "secret") {
		t.Fatalf("intermediate symlink=%d %q", response.Code, response.Body.String())
	}
	if response := perform(host, http.MethodGet, base+"/api/v1/meta?path=file.txt%2Fchild", nil, nil); response.Code != http.StatusNotFound {
		t.Fatalf("non-directory component=%d", response.Code)
	}
	hosted.resource.Backend = overrideBackend{Backend: originalBackend, lstat: func(string) (fs.FileInfo, error) { return nil, errors.New("root stat") }}
	if _, _, err := inspectRequestedPath(hosted, ""); err == nil {
		t.Fatal("root inspection error ignored")
	}
	hosted.resource.Backend = originalBackend

	if response := perform(host, http.MethodGet, "/d/"+uploadDefinition.Token+"/api/v1/download", nil, nil); response.Code != http.StatusNotFound {
		t.Fatal(response.Code)
	}
	basicPolicy := delivery.DefaultPolicy()
	basicPolicy.Auth = delivery.AuthBasic
	basicDefinition := webDefinition(delivery.RoutePathToWeb, root)
	basicDefinition.Token = fixedToken(15)
	basicDefinition.Credentials = policy.Credentials{BasicUsername: "u", BasicPassword: []byte("p")}
	basicID := delivery.ID("00000000-0000-4000-8000-000000000046")
	if err := host.Register(context.Background(), webRecord(basicID, delivery.RoutePathToWeb, basicPolicy), marshalDefinition(t, basicDefinition)); err != nil {
		t.Fatal(err)
	}
	if response := perform(host, http.MethodGet, "/d/"+basicDefinition.Token+"/api/v1/download", nil, nil); response.Code != http.StatusUnauthorized {
		t.Fatalf("unauthorized download=%d", response.Code)
	}
	if response := perform(host, http.MethodGet, base+"/api/v1/download?path=../x", nil, nil); response.Code != http.StatusBadRequest {
		t.Fatal(response.Code)
	}
	if response := perform(host, http.MethodGet, base+"/api/v1/download?path=missing", nil, nil); response.Code != http.StatusNotFound {
		t.Fatal(response.Code)
	}
	if response := perform(host, http.MethodGet, base+"/api/v1/download?path=link", nil, nil); response.Code != http.StatusNotFound {
		t.Fatal(response.Code)
	}

	hosted.resource.Backend = overrideBackend{Backend: originalBackend, lstat: func(string) (fs.FileInfo, error) { return fakeInfo{mode: fs.ModeSocket}, nil }}
	if response := perform(host, http.MethodGet, base+"/api/v1/download", nil, nil); response.Code != http.StatusUnprocessableEntity {
		t.Fatal(response.Code)
	}
	hosted.resource.Backend = overrideBackend{Backend: originalBackend, open: func(string) (io.ReadCloser, error) { return nil, errors.New("open") }}
	if response := perform(host, http.MethodGet, base+"/api/v1/download?path=file.txt", nil, nil); response.Code != http.StatusInternalServerError {
		t.Fatal(response.Code)
	}
	hosted.resource.Backend = overrideBackend{Backend: originalBackend, readDir: func(string) ([]fs.DirEntry, error) { return nil, errors.New("walk") }}
	if response := perform(host, http.MethodGet, base+"/api/v1/download?archive=tar.gz", nil, nil); response.Code != http.StatusUnprocessableEntity {
		t.Fatal(response.Code)
	}
	if err := os.Remove(filepath.Join(root, "link")); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(root, "outside")); err != nil {
		t.Fatal(err)
	}
	hosted.resource.Backend = overrideBackend{Backend: originalBackend, open: func(string) (io.ReadCloser, error) { return nil, errors.New("stream") }}
	_ = perform(host, http.MethodGet, base+"/api/v1/download?archive=tar.gz", nil, nil)
	hosted.resource.Backend = originalBackend
}

func TestManifestStreamAndUploadFailurePaths(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "file")
	if err := os.WriteFile(file, []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}
	host := newTestHost(t, nil)
	definition := webDefinition(delivery.RoutePathToWeb, root)
	record := webRecord(testDeliveryID, delivery.RoutePathToWeb, delivery.DefaultPolicy())
	if err := host.Register(context.Background(), record, marshalDefinition(t, definition)); err != nil {
		t.Fatal(err)
	}
	hosted := host.byID[record.ID]
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := buildManifest(canceled, hosted, root, "root", ""); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	backend := hosted.resource.Backend
	hosted.resource.Backend = overrideBackend{Backend: backend, lstat: func(string) (fs.FileInfo, error) { return nil, errors.New("stat") }}
	if _, err := buildManifest(context.Background(), hosted, root, "root", ""); err == nil {
		t.Fatal("expected stat error")
	}
	hosted.resource.Backend = overrideBackend{Backend: backend, lstat: func(string) (fs.FileInfo, error) { return fakeInfo{mode: fs.ModeSocket}, nil }}
	if _, err := buildManifest(context.Background(), hosted, root, "root", ""); err == nil {
		t.Fatal("expected special error")
	}
	hosted.resource.Backend = backend
	hosted.selector = selectorFunc(func(string, bool) bool { return false })
	if manifest, err := buildManifest(context.Background(), hosted, file, "root", "file"); err != nil || manifest != nil {
		t.Fatalf("excluded manifest=%v %v", manifest, err)
	}
	hosted.selector = selection.All()
	hosted.resource.Backend = overrideBackend{Backend: backend, readDir: func(string) ([]fs.DirEntry, error) { return nil, errors.New("read") }}
	if _, err := buildManifest(context.Background(), hosted, root, "root", ""); err == nil {
		t.Fatal("expected read error")
	}
	hosted.resource.Backend = overrideBackend{Backend: backend, lstat: func(name string) (fs.FileInfo, error) {
		if name == root {
			return backend.Lstat(name)
		}
		return nil, errors.New("child stat")
	}}
	if _, err := buildManifest(context.Background(), hosted, root, "root", ""); err == nil {
		t.Fatal("expected child manifest error")
	}
	hosted.resource.Backend = backend

	authorization, err := hosted.policy.Authorize(context.Background(), policy.Request{Peer: peerAddress("192.0.2.1:1")})
	if err != nil {
		t.Fatal(err)
	}
	defer authorization.Release()
	fileInfo, _ := os.Stat(file)
	if err := streamArchive(context.Background(), &errorResponse{}, backend, []manifestEntry{{path: file, name: "file", info: fileInfo}}, authorization); err == nil {
		t.Fatal("expected archive write error")
	}
	openFailure := overrideBackend{Backend: backend, open: func(string) (io.ReadCloser, error) { return nil, errors.New("open") }}
	if err := streamArchive(context.Background(), io.Discard, openFailure, []manifestEntry{{path: file, name: "file", info: fileInfo}}, authorization); err == nil {
		t.Fatal("expected archive open error")
	}
	closeFailure := overrideBackend{Backend: backend, open: func(string) (io.ReadCloser, error) {
		return errorReadCloser{reader: strings.NewReader("data"), closeErr: errors.New("close")}, nil
	}}
	if err := streamArchive(context.Background(), io.Discard, closeFailure, []manifestEntry{{path: file, name: "file", info: fileInfo}}, authorization); err == nil {
		t.Fatal("expected archive close error")
	}
	shortSource := overrideBackend{Backend: backend, open: func(string) (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader("x")), nil
	}}
	if err := streamArchive(context.Background(), io.Discard, shortSource, []manifestEntry{{path: file, name: "file", info: fileInfo}}, authorization); err == nil {
		t.Fatal("changed archive source accepted")
	}
	if err := streamArchive(context.Background(), io.Discard, backend, []manifestEntry{{path: file, name: "socket", info: fakeInfo{mode: fs.ModeSocket}}}, authorization); err == nil {
		t.Fatal("expected archive header error")
	}

	uploadDefinition := webDefinition(delivery.RouteWebToPath, root)
	uploadDefinition.Token = fixedToken(14)
	uploadID := delivery.ID("00000000-0000-4000-8000-000000000044")
	if err := host.Register(context.Background(), webRecord(uploadID, delivery.RouteWebToPath, delivery.DefaultPolicy()), marshalDefinition(t, uploadDefinition)); err != nil {
		t.Fatal(err)
	}
	uploaded := host.byID[uploadID]
	base := "/d/" + uploadDefinition.Token + "/api/v1/upload"
	if response := perform(host, http.MethodPost, "/d/"+definition.Token+"/api/v1/upload", nil, nil); response.Code != http.StatusNotFound {
		t.Fatal(response.Code)
	}
	empty := &bytes.Buffer{}
	form := multipart.NewWriter(empty)
	contentType := form.FormDataContentType()
	_ = form.Close()
	if response := perform(host, http.MethodPost, base, empty, map[string]string{"Content-Type": contentType}); response.Code != http.StatusBadRequest {
		t.Fatal(response.Code)
	}
	fieldBody := &bytes.Buffer{}
	fieldForm := multipart.NewWriter(fieldBody)
	_ = fieldForm.WriteField("name", "value")
	fieldType := fieldForm.FormDataContentType()
	_ = fieldForm.Close()
	if response := perform(host, http.MethodPost, base, fieldBody, map[string]string{"Content-Type": fieldType}); response.Code != http.StatusBadRequest {
		t.Fatal(response.Code)
	}
	body, contentType := multipartBody(t, ".", []byte("x"), false)
	if response := perform(host, http.MethodPost, base, body, map[string]string{"Content-Type": contentType}); response.Code != http.StatusBadRequest {
		t.Fatal(response.Code)
	}
	uploaded.selector = selectorFunc(func(string, bool) bool { return false })
	body, contentType = multipartBody(t, "excluded", []byte("x"), false)
	if response := perform(host, http.MethodPost, base, body, map[string]string{"Content-Type": contentType}); response.Code != http.StatusBadRequest {
		t.Fatal(response.Code)
	}
	uploaded.selector = selection.All()
	original := uploaded.resource.Backend
	uploaded.resource.Backend = overrideBackend{Backend: original, lstat: func(string) (fs.FileInfo, error) { return nil, errors.New("stat") }}
	body, contentType = multipartBody(t, "stat", []byte("x"), false)
	if response := perform(host, http.MethodPost, base, body, map[string]string{"Content-Type": contentType}); response.Code != http.StatusInternalServerError {
		t.Fatal(response.Code)
	}
	uploaded.resource.Backend = original
	originalRandom := uploadRandom
	uploadRandom = func([]byte) (int, error) { return 0, errors.New("random") }
	t.Cleanup(func() { uploadRandom = originalRandom })
	body, contentType = multipartBody(t, "random", []byte("x"), false)
	if response := perform(host, http.MethodPost, base, body, map[string]string{"Content-Type": contentType}); response.Code != http.StatusInternalServerError {
		t.Fatal(response.Code)
	}
	uploadRandom = originalRandom
	uploaded.resource.Backend = overrideBackend{Backend: original, create: func(string, fs.FileMode) (fsx.Writable, error) { return nil, errors.New("create") }}
	body, contentType = multipartBody(t, "create", []byte("x"), false)
	if response := perform(host, http.MethodPost, base, body, map[string]string{"Content-Type": contentType}); response.Code != http.StatusInternalServerError {
		t.Fatal(response.Code)
	}
	uploaded.resource.Backend = overrideBackend{Backend: original, create: func(string, fs.FileMode) (fsx.Writable, error) {
		return &errorWritable{syncErr: errors.New("sync")}, nil
	}}
	body, contentType = multipartBody(t, "sync", []byte("x"), false)
	if response := perform(host, http.MethodPost, base, body, map[string]string{"Content-Type": contentType}); response.Code != http.StatusUnprocessableEntity {
		t.Fatal(response.Code)
	}
	for _, commitErr := range []error{fsx.ErrDestinationExists, errors.New("commit")} {
		uploaded.resource.Backend = overrideBackend{Backend: original, commit: func(string, string) error { return commitErr }}
		body, contentType = multipartBody(t, "commit"+time.Now().Format("150405.000000000"), []byte("x"), false)
		response := perform(host, http.MethodPost, base, body, map[string]string{"Content-Type": contentType})
		want := http.StatusInternalServerError
		if errors.Is(commitErr, fsx.ErrDestinationExists) {
			want = http.StatusConflict
		}
		if response.Code != want {
			t.Fatalf("commit status=%d want=%d", response.Code, want)
		}
	}
	uploaded.resource.Backend = original

	limitedPolicy := delivery.DefaultPolicy()
	limitedPolicy.UploadRate = delivery.Limit{Value: 1}
	limitedEngine, err := policy.NewEngine("00000000-0000-4000-8000-000000000045", limitedPolicy, policy.Credentials{}, webDependencies())
	if err != nil {
		t.Fatal(err)
	}
	limitedAuthorization, err := limitedEngine.Authorize(context.Background(), policy.Request{Peer: peerAddress("192.0.2.1:1"), Incoming: true, DeclaredSize: -1})
	if err != nil {
		t.Fatal(err)
	}
	defer limitedAuthorization.Release()
	if err := limitedAuthorization.Wait(context.Background(), policy.Upload, 1); err != nil {
		t.Fatal(err)
	}
	cancelContext, cancelWait := context.WithCancel(context.Background())
	cancelWait()
	uploadWriter := &uploadWriter{ctx: cancelContext, destination: io.Discard, authorization: limitedAuthorization}
	if _, err := uploadWriter.Write([]byte("x")); err == nil {
		t.Fatal("expected upload wait cancellation")
	}
	limitedWriter := &limitedWriter{ctx: cancelContext, destination: io.Discard, authorization: limitedAuthorization, direction: policy.Upload}
	if _, err := limitedWriter.Write([]byte("x")); err == nil {
		t.Fatal("expected download wait cancellation")
	}
}

func TestBrowserOriginValidation(t *testing.T) {
	tests := []struct {
		target string
		host   string
		origin string
		want   bool
	}{
		{target: "http://example.test/path", host: "example.test", want: true},
		{target: "http://example.test/path", host: "EXAMPLE.test", origin: "http://example.test", want: true},
		{target: "https://example.test/path", host: "example.test", origin: "https://example.test", want: true},
		{target: "http://example.test/path", host: "example.test", origin: "https://example.test", want: false},
		{target: "http://example.test/path", host: "example.test", origin: "http://other.test", want: false},
		{target: "http://example.test/path", host: "example.test", origin: "%", want: false},
		{target: "http://example.test/path", host: "example.test", origin: "http://user@example.test", want: false},
		{target: "http://example.test/path", host: "example.test", origin: "http:///missing", want: false},
		{target: "http://example.test/path", host: "example.test", origin: "http://example.test/path", want: false},
		{target: "http://example.test/path", host: "example.test", origin: "http://example.test?query", want: false},
		{target: "http://example.test/path", host: "example.test", origin: "http://example.test#fragment", want: false},
	}
	for _, test := range tests {
		request := httptest.NewRequest(http.MethodPost, test.target, nil)
		request.Host = test.host
		if test.origin != "" {
			request.Header.Set("Origin", test.origin)
		}
		if got := sameOrigin(request); got != test.want {
			t.Fatalf("target=%q host=%q origin=%q got=%v want=%v", test.target, test.host, test.origin, got, test.want)
		}
	}
}

type fakeInfo struct {
	directory bool
	mode      fs.FileMode
}

func (info fakeInfo) Name() string { return "fake" }
func (info fakeInfo) Size() int64  { return 0 }
func (info fakeInfo) Mode() fs.FileMode {
	if info.directory {
		return fs.ModeDir
	}
	return info.mode
}
func (info fakeInfo) ModTime() time.Time { return webTestTime }
func (info fakeInfo) IsDir() bool        { return info.directory }
func (info fakeInfo) Sys() any           { return nil }
