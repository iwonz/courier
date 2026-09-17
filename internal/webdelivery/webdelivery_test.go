package webdelivery

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/endpoint"
	"github.com/iwonz/courier/internal/fsx"
	"github.com/iwonz/courier/internal/operation"
	"github.com/iwonz/courier/internal/policy"
	"github.com/iwonz/courier/internal/selection"
)

const (
	testServerID   delivery.ID = "00000000-0000-4000-8000-000000000031"
	testDeliveryID delivery.ID = "00000000-0000-4000-8000-000000000032"
)

var webTestTime = time.Unix(1_700_000_000, 0).UTC()

type testFailReader struct{}

func (testFailReader) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }

func webDependencies() policy.Dependencies {
	return policy.Dependencies{
		Random: bytes.NewReader(bytes.Repeat([]byte{8}, 8192)), Now: func() time.Time { return webTestTime },
		Argon2:          policy.Argon2Parameters{Time: 1, MemoryKiB: 8, Threads: 1, KeyBytes: 16, SaltBytes: 16},
		AuthConcurrency: 1, Sessions: policy.SessionConfig{IdleTimeout: time.Minute, AbsoluteTimeout: time.Hour},
	}
}

func webRecord(id delivery.ID, route delivery.Route, configured delivery.Policy) delivery.Delivery {
	return delivery.Delivery{
		ID: id, ServerID: testServerID, Route: route, State: delivery.StateStarting,
		Policy: configured, CreatedAt: webTestTime, UpdatedAt: webTestTime,
	}
}

func fixedToken(value byte) string {
	token, err := NewToken(bytes.NewReader(bytes.Repeat([]byte{value}, ResourceTokenBytes)))
	if err != nil {
		panic(err)
	}
	return token
}

func webDefinition(route delivery.Route, localPath string) Definition {
	definition := Definition{Version: DefinitionVersion, Token: fixedToken(1)}
	if route == delivery.RouteWebToPath {
		definition.Source = "web://"
		definition.Destination = localPath
	} else if route == delivery.RouteWebhookToPath {
		definition.Source = "webhook://"
		definition.Destination = localPath
	} else {
		definition.Source = localPath
		definition.Destination = "web://"
	}
	return definition
}

func marshalDefinition(t *testing.T, definition Definition) json.RawMessage {
	t.Helper()
	data, err := json.Marshal(definition)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func localOpen(_ context.Context, value endpoint.Endpoint, _ EndpointRuntime) (*Resource, error) {
	return &Resource{Endpoint: value, Backend: fsx.Local{}, Path: value.Path}, nil
}

func newTestHost(t *testing.T, stop func(delivery.ID)) *Host {
	t.Helper()
	host, err := NewHost(HostOptions{
		Open: localOpen, Select: func(rules []operation.SelectionRule) (selection.Selector, error) {
			return selection.Compile(rules, func(string) (io.ReadCloser, error) { return nil, fs.ErrNotExist })
		},
		PolicyDependencies: webDependencies(), StopRequested: stop,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = host.Close(context.Background()) })
	return host
}

func perform(host *Host, method, target string, body io.Reader, headers map[string]string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, target, body)
	request.RemoteAddr = "192.0.2.10:4321"
	for name, value := range headers {
		request.Header.Set(name, value)
	}
	response := httptest.NewRecorder()
	host.routes().ServeHTTP(response, request)
	return response
}

func TestDefinitionValidationAndURL(t *testing.T) {
	if _, err := NewToken(nil); err == nil {
		t.Fatal("expected random dependency error")
	}
	if _, err := NewToken(testFailReader{}); err == nil {
		t.Fatal("expected token generation error")
	}
	plan, err := operation.Build(operation.Request{Source: "web://", Destination: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	definition, err := NewDefinition(plan, policy.Credentials{}, EndpointRuntime{}, bytes.NewReader(bytes.Repeat([]byte{2}, ResourceTokenBytes)))
	if err != nil || routeForPlan(plan) != delivery.RouteWebToPath || definition.Token == "" {
		t.Fatalf("unexpected definition: %#v %v", definition, err)
	}
	deliveryCredentials := policy.Credentials{BasicUsername: "user", BasicPassword: []byte("basic"), Password: []byte("password")}
	credentialDefinition, err := NewDefinition(plan, deliveryCredentials, EndpointRuntime{}, bytes.NewReader(bytes.Repeat([]byte{4}, ResourceTokenBytes)))
	if err != nil {
		t.Fatal(err)
	}
	clear(deliveryCredentials.BasicPassword)
	clear(deliveryCredentials.Password)
	if string(credentialDefinition.Credentials.BasicPassword) != "basic" || string(credentialDefinition.Credentials.Password) != "password" {
		t.Fatal("delivery credentials were not cloned")
	}
	credentialDefinition.ClearSecrets()
	if credentialDefinition.Credentials.BasicUsername != "" || len(credentialDefinition.Credentials.BasicPassword) != 0 || len(credentialDefinition.Credentials.Password) != 0 {
		t.Fatal("delivery credentials were not cleared")
	}
	if _, err := NewDefinition(plan, policy.Credentials{}, EndpointRuntime{}, testFailReader{}); err == nil {
		t.Fatal("expected definition token error")
	}
	badPlan := plan
	badPlan.Route = operation.RoutePathToPath
	if _, err := NewDefinition(badPlan, policy.Credentials{}, EndpointRuntime{}, bytes.NewReader(bytes.Repeat([]byte{2}, ResourceTokenBytes))); err == nil {
		t.Fatal("expected unsupported plan")
	}
	pathPlan, err := operation.Build(operation.Request{Source: t.TempDir(), Destination: "web://"})
	if err != nil || routeForPlan(pathPlan) != delivery.RoutePathToWeb {
		t.Fatalf("unexpected path plan: %v", err)
	}
	webhookPlan, err := operation.Build(operation.Request{Source: "webhook://", Destination: t.TempDir()})
	if err != nil || routeForPlan(webhookPlan) != delivery.RouteWebhookToPath {
		t.Fatalf("unexpected webhook plan: %v", err)
	}
	webhookDefinition, err := NewDefinition(webhookPlan, policy.Credentials{}, EndpointRuntime{}, bytes.NewReader(bytes.Repeat([]byte{5}, ResourceTokenBytes)))
	if err != nil || webhookDefinition.Validate(delivery.RouteWebhookToPath) != nil {
		t.Fatalf("unexpected webhook definition: %+v %v", webhookDefinition, err)
	}
	for index, mutate := range []func(*Definition){
		func(value *Definition) { value.Source = "web://" },
		func(value *Definition) { value.Destination = "web://" },
		func(value *Definition) { value.Archive = true },
		func(value *Definition) { value.NoUI = true },
	} {
		candidate := webhookDefinition
		mutate(&candidate)
		if err := candidate.Validate(delivery.RouteWebhookToPath); err == nil {
			t.Fatalf("webhook mutation %d accepted", index)
		}
	}

	valid := webDefinition(delivery.RouteWebToPath, t.TempDir())
	mutations := []func(*Definition){
		func(value *Definition) { value.Version = 2 },
		func(value *Definition) { value.Source = " web://" },
		func(value *Definition) { value.Token = "invalid" },
		func(value *Definition) { value.Source = "unknown://" },
		func(value *Definition) { value.Source = "local" },
		func(value *Definition) { value.Archive = true },
		func(value *Definition) { value.Selection = []operation.SelectionRule{{Value: ""}} },
	}
	for index, mutate := range mutations {
		candidate := valid
		mutate(&candidate)
		if err := candidate.Validate(delivery.RouteWebToPath); err == nil {
			t.Fatalf("mutation %d unexpectedly valid", index)
		}
	}
	pathDefinition := webDefinition(delivery.RoutePathToWeb, t.TempDir())
	pathDefinition.Extract = true
	if err := pathDefinition.Validate(delivery.RoutePathToWeb); err == nil {
		t.Fatal("expected extract rejection")
	}
	pathDefinition.Extract = false
	if err := pathDefinition.Validate(delivery.RoutePathToHTTP); err == nil {
		t.Fatal("expected route rejection")
	}
	if value, err := URL("0.0.0.0:8080", valid.Token); err != nil || value != "http://127.0.0.1:8080/d/"+valid.Token+"/" {
		t.Fatalf("unexpected wildcard URL: %q %v", value, err)
	}
	if value, err := URL("[::]:8080", valid.Token); err != nil || value != "http://127.0.0.1:8080/d/"+valid.Token+"/" {
		t.Fatalf("unexpected IPv6 wildcard URL: %q %v", value, err)
	}
	if _, err := URL("invalid", valid.Token); err == nil {
		t.Fatal("expected URL bind error")
	}
	if value, err := WebhookURL("127.0.0.1:8080", valid.Token); err != nil || value != "http://127.0.0.1:8080/d/"+valid.Token+"/upload" {
		t.Fatalf("unexpected webhook URL: %q %v", value, err)
	}
	if _, err := RandomDefinition(plan, policy.Credentials{}); err != nil {
		t.Fatal(err)
	}
	endpointRuntime := EndpointRuntime{Credentials: []EndpointCredential{{Prompt: "Password for user@host: ", Secret: []byte("secret")}}}
	remotePlan, err := operation.Build(operation.Request{Source: "host:/source", Destination: "web://"})
	if err != nil {
		t.Fatal(err)
	}
	remoteDefinition, err := NewDefinition(remotePlan, policy.Credentials{}, endpointRuntime, bytes.NewReader(bytes.Repeat([]byte{3}, ResourceTokenBytes)))
	if err != nil || string(remoteDefinition.Endpoint.Credentials[0].Secret) != "secret" {
		t.Fatalf("remote runtime=%#v err=%v", remoteDefinition.Endpoint, err)
	}
	clear(endpointRuntime.Credentials[0].Secret)
	if string(remoteDefinition.Endpoint.Credentials[0].Secret) != "secret" {
		t.Fatal("endpoint runtime was not cloned")
	}
	secret, err := remoteDefinition.Endpoint.Prompt("Password for user@host: ")
	if err != nil || string(secret) != "secret" {
		t.Fatalf("runtime prompt=%q err=%v", secret, err)
	}
	clear(secret)
	if _, err := remoteDefinition.Endpoint.Prompt("unknown"); err == nil {
		t.Fatal("unknown endpoint prompt accepted")
	}
	remoteDefinition.Endpoint.Clear()
	if len(remoteDefinition.Endpoint.Credentials) != 0 {
		t.Fatal("endpoint credentials were not cleared")
	}
	invalidRuntimes := []EndpointRuntime{
		{Credentials: make([]EndpointCredential, maxEndpointPrompts+1)},
		{Credentials: []EndpointCredential{{Prompt: ""}}},
		{Credentials: []EndpointCredential{{Prompt: "   "}}},
		{Credentials: []EndpointCredential{{Prompt: "bad\nlabel"}}},
		{Credentials: []EndpointCredential{{Prompt: strings.Repeat("p", maxPromptBytes+1)}}},
		{Credentials: []EndpointCredential{{Prompt: "prompt", Secret: bytes.Repeat([]byte("s"), maxSecretBytes+1)}}},
	}
	for index, runtime := range invalidRuntimes {
		if err := validateEndpointRuntime(runtime); err == nil {
			t.Fatalf("invalid endpoint runtime %d accepted", index)
		}
	}
	remoteWithInvalidCredentials := webDefinition(delivery.RoutePathToWeb, t.TempDir())
	remoteWithInvalidCredentials.Source = "host:/source"
	remoteWithInvalidCredentials.Endpoint = EndpointRuntime{Credentials: []EndpointCredential{{Prompt: ""}}}
	if err := remoteWithInvalidCredentials.Validate(delivery.RoutePathToWeb); err == nil {
		t.Fatal("invalid remote endpoint credentials accepted")
	}
	localWithCredentials := webDefinition(delivery.RouteWebToPath, t.TempDir())
	localWithCredentials.Endpoint = EndpointRuntime{Credentials: []EndpointCredential{{Prompt: "prompt"}}}
	if err := localWithCredentials.Validate(delivery.RouteWebToPath); err == nil {
		t.Fatal("local endpoint credentials accepted")
	}
}

func TestHostRegistrationLifecycle(t *testing.T) {
	if _, err := NewHost(HostOptions{}); err == nil {
		t.Fatal("expected dependency error")
	}
	server := &http.Server{}
	host, err := NewHost(HostOptions{Open: localOpen, Select: func([]operation.SelectionRule) (selection.Selector, error) { return selection.All(), nil }, PolicyDependencies: webDependencies(), Server: server})
	if err != nil || server.Handler == nil {
		t.Fatalf("expected configured server: %v", err)
	}
	defer host.Close(context.Background())
	configured := delivery.DefaultPolicy()
	record := webRecord(testDeliveryID, delivery.RouteWebToPath, configured)
	root := t.TempDir()
	definition := webDefinition(record.Route, root)
	if err := host.Register(context.Background(), delivery.Delivery{}, nil); err == nil {
		t.Fatal("expected record validation error")
	}
	if err := host.Register(context.Background(), record, []byte("{")); err == nil {
		t.Fatal("expected JSON error")
	}
	if err := host.Register(context.Background(), record, append(marshalDefinition(t, definition), []byte("{}")...)); err == nil {
		t.Fatal("expected trailing JSON error")
	}
	invalid := definition
	invalid.Token = "bad"
	if err := host.Register(context.Background(), record, marshalDefinition(t, invalid)); err == nil {
		t.Fatal("expected definition error")
	}
	host.selectRules = func([]operation.SelectionRule) (selection.Selector, error) { return nil, errors.New("selection") }
	if err := host.Register(context.Background(), record, marshalDefinition(t, definition)); err == nil {
		t.Fatal("expected selection error")
	}
	host.selectRules = func([]operation.SelectionRule) (selection.Selector, error) { return selection.All(), nil }
	host.open = func(context.Context, endpoint.Endpoint, EndpointRuntime) (*Resource, error) {
		return nil, errors.New("open")
	}
	if err := host.Register(context.Background(), record, marshalDefinition(t, definition)); err == nil {
		t.Fatal("expected open error")
	}
	closed := atomic.Int32{}
	host.open = func(_ context.Context, value endpoint.Endpoint, _ EndpointRuntime) (*Resource, error) {
		return &Resource{Endpoint: value, Backend: fsx.Local{}, Path: filepath.Join(root, "missing"), Close: func() error { closed.Add(1); return nil }}, nil
	}
	if err := host.Register(context.Background(), record, marshalDefinition(t, definition)); err == nil || closed.Load() != 1 {
		t.Fatalf("expected stat rollback: %v closed=%d", err, closed.Load())
	}
	file := filepath.Join(root, "file")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	host.open = func(_ context.Context, value endpoint.Endpoint, _ EndpointRuntime) (*Resource, error) {
		return &Resource{Endpoint: value, Backend: fsx.Local{}, Path: file, Close: func() error { closed.Add(1); return nil }}, nil
	}
	if err := host.Register(context.Background(), record, marshalDefinition(t, definition)); err == nil {
		t.Fatal("expected upload directory error")
	}
	host.open = localOpen
	if err := host.Register(context.Background(), record, marshalDefinition(t, definition)); err != nil {
		t.Fatal(err)
	}
	if err := host.Register(context.Background(), record, marshalDefinition(t, definition)); err == nil {
		t.Fatal("expected duplicate token")
	}
	second := definition
	second.Token = fixedToken(2)
	if err := host.Register(context.Background(), record, marshalDefinition(t, second)); !errors.Is(err, delivery.ErrDuplicate) {
		t.Fatalf("expected duplicate delivery: %v", err)
	}
	next := configured
	next.Version++
	if _, _, err := host.PreparePolicy(context.Background(), "00000000-0000-4000-8000-000000000099", configured.Version, next); !errors.Is(err, delivery.ErrNotFound) {
		t.Fatalf("expected update not found: %v", err)
	}
	if _, _, err := host.PreparePolicy(context.Background(), record.ID, 99, next); !errors.Is(err, delivery.ErrRevisionConflict) {
		t.Fatalf("expected update conflict: %v", err)
	}
	commit, cancel, err := host.PreparePolicy(context.Background(), record.ID, configured.Version, next)
	if err != nil {
		t.Fatal(err)
	}
	commit()
	cancel()
	if err := host.Stop(context.Background(), record.ID); err != nil || len(host.byID) != 0 {
		t.Fatalf("stop failed: %v", err)
	}
	if err := host.Stop(context.Background(), record.ID); err != nil {
		t.Fatal(err)
	}
	if err := host.Serve(nil); err == nil {
		t.Fatal("expected listener error")
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	result := make(chan error, 1)
	go func() { result <- host.Serve(listener) }()
	if err := host.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := <-result; err != nil {
		t.Fatal(err)
	}
}

func TestBrowserDownloadMetadataArchiveAndNoUI(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "folder"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "hello.txt"), []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "folder", "nested.txt"), []byte("nested"), 0o600); err != nil {
		t.Fatal(err)
	}
	host := newTestHost(t, nil)
	assetRequest := httptest.NewRequest(http.MethodGet, "/assets/data.js", nil)
	assetResponse := httptest.NewRecorder()
	host.Handler().ServeHTTP(assetResponse, assetRequest)
	if assetResponse.Code != http.StatusOK || !strings.Contains(assetResponse.Header().Get("Content-Type"), "javascript") || assetResponse.Header().Get("Content-Security-Policy") == "" || assetResponse.Header().Get("X-Content-Type-Options") != "nosniff" || assetResponse.Body.Len() == 0 {
		t.Fatalf("embedded asset=%d %q", assetResponse.Code, assetResponse.Header().Get("Content-Type"))
	}
	configured := delivery.DefaultPolicy()
	definition := webDefinition(delivery.RoutePathToWeb, root)
	record := webRecord(testDeliveryID, delivery.RoutePathToWeb, configured)
	if err := host.Register(context.Background(), record, marshalDefinition(t, definition)); err != nil {
		t.Fatal(err)
	}
	base := "/d/" + definition.Token
	if response := perform(host, http.MethodGet, "/d/unknown/", nil, nil); response.Code != http.StatusNotFound {
		t.Fatalf("unknown status=%d", response.Code)
	}
	if response := perform(host, http.MethodGet, "/d/unknown/api/v1/meta", nil, nil); response.Code != http.StatusNotFound {
		t.Fatalf("unknown metadata status=%d", response.Code)
	}
	response := perform(host, http.MethodGet, base+"/", nil, nil)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `<div id="root"></div>`) {
		t.Fatalf("root=%d %q", response.Code, response.Body.String())
	}
	response = perform(host, http.MethodGet, base+"/api/v1/meta", nil, nil)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "hello.txt") || !strings.Contains(response.Body.String(), "folder") {
		t.Fatalf("metadata=%d %q", response.Code, response.Body.String())
	}
	response = perform(host, http.MethodGet, base+"/api/v1/meta?path=folder", nil, nil)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "nested.txt") {
		t.Fatalf("nested metadata=%d %q", response.Code, response.Body.String())
	}
	response = perform(host, http.MethodGet, base+"/api/v1/download?path=hello.txt", nil, nil)
	if response.Code != http.StatusOK || response.Body.String() != "hello" || response.Header().Get("Content-Disposition") == "" {
		t.Fatalf("download=%d %q", response.Code, response.Body.String())
	}
	response = perform(host, http.MethodGet, base+"/api/v1/download?path=folder", nil, nil)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("directory without archive=%d", response.Code)
	}
	response = perform(host, http.MethodGet, base+"/api/v1/download?archive=tar.gz", nil, nil)
	if response.Code != http.StatusOK {
		t.Fatalf("archive=%d %q", response.Code, response.Body.String())
	}
	gzipReader, err := gzip.NewReader(bytes.NewReader(response.Body.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	tarReader := tar.NewReader(gzipReader)
	entries := map[string]string{}
	for {
		header, nextErr := tarReader.Next()
		if nextErr == io.EOF {
			break
		}
		if nextErr != nil {
			t.Fatal(nextErr)
		}
		data, _ := io.ReadAll(tarReader)
		entries[header.Name] = string(data)
	}
	if entries[filepath.Base(root)+"/hello.txt"] != "hello" || entries[filepath.Base(root)+"/folder/nested.txt"] != "nested" {
		t.Fatalf("archive entries=%v", entries)
	}
	if response := perform(host, http.MethodGet, base+"/api/v1/meta?path=../secret", nil, nil); response.Code != http.StatusBadRequest {
		t.Fatalf("traversal status=%d", response.Code)
	}
	definition2 := webDefinition(delivery.RoutePathToWeb, root)
	definition2.Token = fixedToken(3)
	definition2.NoUI = true
	record2 := webRecord("00000000-0000-4000-8000-000000000033", delivery.RoutePathToWeb, configured)
	if err := host.Register(context.Background(), record2, marshalDefinition(t, definition2)); err != nil {
		t.Fatal(err)
	}
	response = perform(host, http.MethodGet, "/d/"+definition2.Token+"/", nil, nil)
	if response.Code != http.StatusOK || !strings.Contains(response.Header().Get("Content-Type"), "application/json") || strings.Contains(response.Body.String(), "html") {
		t.Fatalf("no-ui root=%d %q", response.Code, response.Body.String())
	}
	limitedPolicy := delivery.DefaultPolicy()
	limitedPolicy.DeliveryLimit = delivery.Limit{Value: 1}
	limitedDefinition := webDefinition(delivery.RoutePathToWeb, root)
	limitedDefinition.Token = fixedToken(5)
	limitedID := delivery.ID("00000000-0000-4000-8000-000000000035")
	if err := host.Register(context.Background(), webRecord(limitedID, delivery.RoutePathToWeb, limitedPolicy), marshalDefinition(t, limitedDefinition)); err != nil {
		t.Fatal(err)
	}
	busy, err := host.byID[limitedID].policy.Authorize(context.Background(), policy.Request{Peer: peerAddress("192.0.2.20:1"), Transfer: true, DeclaredSize: -1})
	if err != nil {
		t.Fatal(err)
	}
	if response := perform(host, http.MethodGet, "/d/"+limitedDefinition.Token+"/api/v1/download?path=hello.txt", nil, nil); response.Code != http.StatusTooManyRequests {
		t.Fatalf("download reservation limit=%d", response.Code)
	}
	busy.Release()
}

func multipartBody(t *testing.T, name string, data []byte, second bool) (*bytes.Buffer, string) {
	t.Helper()
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write(data)
	if second {
		part, _ = writer.CreateFormFile("second", "second.txt")
		_, _ = part.Write([]byte("second"))
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return buffer, writer.FormDataContentType()
}

func multipartFieldBody(t *testing.T, field, name string, data []byte, second bool) (*bytes.Buffer, string) {
	t.Helper()
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	part, err := writer.CreateFormFile(field, name)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write(data)
	if second {
		part, _ = writer.CreateFormFile("second", "second.txt")
		_, _ = part.Write([]byte("second"))
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return buffer, writer.FormDataContentType()
}

func tarGzipBody(t *testing.T, entries map[string]string) []byte {
	t.Helper()
	var output bytes.Buffer
	gzipWriter := gzip.NewWriter(&output)
	tarWriter := tar.NewWriter(gzipWriter)
	for name, value := range entries {
		header := &tar.Header{Name: name, Mode: 0o600, Size: int64(len(value)), ModTime: webTestTime}
		if err := tarWriter.WriteHeader(header); err != nil {
			t.Fatal(err)
		}
		if _, err := io.WriteString(tarWriter, value); err != nil {
			t.Fatal(err)
		}
	}
	if err := tarWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gzipWriter.Close(); err != nil {
		t.Fatal(err)
	}
	return output.Bytes()
}

func TestBrowserUploadIsTransactional(t *testing.T) {
	root := t.TempDir()
	host := newTestHost(t, nil)
	configured := delivery.DefaultPolicy()
	configured.MaxFileSize = delivery.Limit{Value: 5}
	definition := webDefinition(delivery.RouteWebToPath, root)
	record := webRecord(testDeliveryID, delivery.RouteWebToPath, configured)
	if err := host.Register(context.Background(), record, marshalDefinition(t, definition)); err != nil {
		t.Fatal(err)
	}
	base := "/d/" + definition.Token + "/api/v1/upload"
	body, contentType := multipartBody(t, "hello.txt", []byte("hello"), false)
	response := perform(host, http.MethodPost, base, body, map[string]string{"Content-Type": contentType, "X-Courier-File-Size": "5"})
	data, err := os.ReadFile(filepath.Join(root, "hello.txt"))
	if response.Code != http.StatusCreated || err != nil || string(data) != "hello" || !strings.Contains(response.Body.String(), `"bytes":5`) {
		t.Fatalf("upload=%d %q data=%q err=%v", response.Code, response.Body.String(), data, err)
	}
	body, contentType = multipartBody(t, "hello.txt", []byte("other"), false)
	if response := perform(host, http.MethodPost, base, body, map[string]string{"Content-Type": contentType}); response.Code != http.StatusConflict {
		t.Fatalf("collision=%d", response.Code)
	}
	body, contentType = multipartBody(t, "large.txt", []byte("123456"), false)
	if response := perform(host, http.MethodPost, base, body, map[string]string{"Content-Type": contentType}); response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("large unknown upload=%d %q", response.Code, response.Body.String())
	}
	if _, err := os.Stat(filepath.Join(root, "large.txt")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatal("oversized upload committed")
	}
	body, contentType = multipartBody(t, "declared.txt", []byte("x"), false)
	if response := perform(host, http.MethodPost, base, body, map[string]string{"Content-Type": contentType, "X-Courier-File-Size": "6"}); response.Code != http.StatusTooManyRequests {
		t.Fatalf("declared limit=%d", response.Code)
	}
	body, contentType = multipartBody(t, "two.txt", []byte("one"), true)
	if response := perform(host, http.MethodPost, base, body, map[string]string{"Content-Type": contentType}); response.Code != http.StatusBadRequest {
		t.Fatalf("two files=%d", response.Code)
	}
	if response := perform(host, http.MethodPost, base, strings.NewReader("bad"), map[string]string{"Content-Type": "text/plain"}); response.Code != http.StatusBadRequest {
		t.Fatalf("media type=%d", response.Code)
	}
	if response := perform(host, http.MethodPost, base, strings.NewReader("bad"), map[string]string{"Sec-Fetch-Site": "cross-site"}); response.Code != http.StatusForbidden {
		t.Fatalf("cross-site upload=%d", response.Code)
	}
	if response := perform(host, http.MethodPost, base, strings.NewReader("bad"), map[string]string{"Content-Type": "multipart/form-data; boundary=x", "X-Courier-File-Size": "invalid"}); response.Code != http.StatusBadRequest {
		t.Fatalf("declared parse=%d", response.Code)
	}
}

func TestBrowserUploadExtraction(t *testing.T) {
	root := t.TempDir()
	host := newTestHost(t, nil)
	definition := webDefinition(delivery.RouteWebToPath, root)
	definition.Extract = true
	record := webRecord(testDeliveryID, delivery.RouteWebToPath, delivery.DefaultPolicy())
	if err := host.Register(context.Background(), record, marshalDefinition(t, definition)); err != nil {
		t.Fatal(err)
	}
	base := "/d/" + definition.Token + "/api/v1/upload"
	payload := tarGzipBody(t, map[string]string{"bundle/file.txt": "contents"})
	body, contentType := multipartBody(t, "bundle.tar.gz", payload, false)
	response := perform(host, http.MethodPost, base, body, map[string]string{"Content-Type": contentType, "X-Courier-File-Size": strconv.Itoa(len(payload))})
	data, err := os.ReadFile(filepath.Join(root, "bundle", "file.txt"))
	if response.Code != http.StatusCreated || err != nil || string(data) != "contents" || !strings.Contains(response.Body.String(), `"entries":1`) {
		t.Fatalf("extract=%d %q data=%q err=%v", response.Code, response.Body.String(), data, err)
	}
	body, contentType = multipartBody(t, "again.tar.gz", payload, false)
	if response := perform(host, http.MethodPost, base, body, map[string]string{"Content-Type": contentType}); response.Code != http.StatusConflict {
		t.Fatalf("extraction collision=%d %q", response.Code, response.Body.String())
	}
	body, contentType = multipartBody(t, "unsupported.zip", []byte("not an archive"), false)
	if response := perform(host, http.MethodPost, base, body, map[string]string{"Content-Type": contentType}); response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("unsupported extraction=%d %q", response.Code, response.Body.String())
	}
}

func TestIncomingWebhookProfileAndExtraction(t *testing.T) {
	root := t.TempDir()
	host := newTestHost(t, nil)
	definition := webDefinition(delivery.RouteWebhookToPath, root)
	record := webRecord(testDeliveryID, delivery.RouteWebhookToPath, delivery.DefaultPolicy())
	if err := host.Register(context.Background(), record, marshalDefinition(t, definition)); err != nil {
		t.Fatal(err)
	}
	resource := "/d/" + definition.Token
	if response := perform(host, http.MethodPost, "/d/unknown/upload", nil, nil); response.Code != http.StatusNotFound {
		t.Fatalf("unknown webhook=%d", response.Code)
	}
	browserRoot := t.TempDir()
	browserDefinition := webDefinition(delivery.RouteWebToPath, browserRoot)
	browserDefinition.Token = fixedToken(8)
	browserID := delivery.ID("00000000-0000-4000-8000-000000000038")
	if err := host.Register(context.Background(), webRecord(browserID, delivery.RouteWebToPath, delivery.DefaultPolicy()), marshalDefinition(t, browserDefinition)); err != nil {
		t.Fatal(err)
	}
	if response := perform(host, http.MethodPost, "/d/"+browserDefinition.Token+"/upload", nil, nil); response.Code != http.StatusNotFound {
		t.Fatalf("browser token accepted as webhook=%d", response.Code)
	}
	for _, surface := range []struct{ method, target string }{
		{http.MethodGet, resource + "/"},
		{http.MethodPost, resource + "/api/v1/session"},
		{http.MethodGet, resource + "/api/v1/meta"},
		{http.MethodGet, resource + "/api/v1/download"},
		{http.MethodPost, resource + "/api/v1/upload"},
	} {
		if response := perform(host, surface.method, surface.target, nil, nil); response.Code != http.StatusNotFound {
			t.Fatalf("webhook exposed browser surface %q: %d", surface.target, response.Code)
		}
	}
	body, contentType := multipartFieldBody(t, "wrong", "wrong.txt", []byte("bad"), false)
	if response := perform(host, http.MethodPost, resource+"/upload", body, map[string]string{"Content-Type": contentType}); response.Code != http.StatusBadRequest {
		t.Fatalf("wrong field=%d %q", response.Code, response.Body.String())
	}
	body, contentType = multipartFieldBody(t, "file", "many.txt", []byte("one"), true)
	if response := perform(host, http.MethodPost, resource+"/upload", body, map[string]string{"Content-Type": contentType}); response.Code != http.StatusBadRequest {
		t.Fatalf("multiple fields=%d %q", response.Code, response.Body.String())
	}
	body, contentType = multipartFieldBody(t, "file", "accepted.txt", []byte("accepted"), false)
	response := perform(host, http.MethodPost, resource+"/upload", body, map[string]string{"Content-Type": contentType, "Origin": "https://automation.example"})
	data, err := os.ReadFile(filepath.Join(root, "accepted.txt"))
	if response.Code != http.StatusCreated || err != nil || string(data) != "accepted" {
		t.Fatalf("accepted=%d %q data=%q err=%v", response.Code, response.Body.String(), data, err)
	}
	body, contentType = multipartFieldBody(t, "file", "accepted.txt", []byte("replace"), false)
	if response := perform(host, http.MethodPost, resource+"/upload", body, map[string]string{"Content-Type": contentType}); response.Code != http.StatusConflict {
		t.Fatalf("collision=%d %q", response.Code, response.Body.String())
	}

	extractRoot := t.TempDir()
	extractDefinition := webDefinition(delivery.RouteWebhookToPath, extractRoot)
	extractDefinition.Token = fixedToken(7)
	extractDefinition.Extract = true
	extractID := delivery.ID("00000000-0000-4000-8000-000000000037")
	if err := host.Register(context.Background(), webRecord(extractID, delivery.RouteWebhookToPath, delivery.DefaultPolicy()), marshalDefinition(t, extractDefinition)); err != nil {
		t.Fatal(err)
	}
	payload := tarGzipBody(t, map[string]string{"bundle/file.txt": "contents"})
	body, contentType = multipartFieldBody(t, "file", "bundle.tar.gz", payload, false)
	response = perform(host, http.MethodPost, "/d/"+extractDefinition.Token+"/upload", body, map[string]string{"Content-Type": contentType})
	data, err = os.ReadFile(filepath.Join(extractRoot, "bundle", "file.txt"))
	if response.Code != http.StatusCreated || err != nil || string(data) != "contents" {
		t.Fatalf("extract=%d %q data=%q err=%v", response.Code, response.Body.String(), data, err)
	}
}

func TestIncomingWebhookBasicAuthentication(t *testing.T) {
	root := t.TempDir()
	host := newTestHost(t, nil)
	configured := delivery.DefaultPolicy()
	configured.Auth = delivery.AuthBasic
	definition := webDefinition(delivery.RouteWebhookToPath, root)
	definition.Credentials = policy.Credentials{BasicUsername: "hook", BasicPassword: []byte("secret")}
	record := webRecord(testDeliveryID, delivery.RouteWebhookToPath, configured)
	if err := host.Register(context.Background(), record, marshalDefinition(t, definition)); err != nil {
		t.Fatal(err)
	}
	body, contentType := multipartFieldBody(t, "file", "auth.txt", []byte("accepted"), false)
	response := perform(host, http.MethodPost, "/d/"+definition.Token+"/upload", body, map[string]string{"Content-Type": contentType})
	if response.Code != http.StatusUnauthorized || response.Header().Get("WWW-Authenticate") == "" {
		t.Fatalf("unauthorized=%d headers=%v", response.Code, response.Header())
	}
	body, contentType = multipartFieldBody(t, "file", "auth.txt", []byte("accepted"), false)
	response = perform(host, http.MethodPost, "/d/"+definition.Token+"/upload", body, map[string]string{"Content-Type": contentType, "Authorization": "Basic aG9vazpzZWNyZXQ="})
	if response.Code != http.StatusCreated {
		t.Fatalf("authorized=%d %q", response.Code, response.Body.String())
	}
}

func TestExplicitBrowserArchiveIsPreparedAndCleaned(t *testing.T) {
	if _, err := prepareArchiveResource(context.Background(), &Resource{}, "", selection.All()); err == nil {
		t.Fatal("invalid archive resource accepted")
	}
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "file.txt"), []byte("archive"), 0o600); err != nil {
		t.Fatal(err)
	}
	closed := atomic.Int32{}
	host := newTestHost(t, nil)
	host.open = func(_ context.Context, value endpoint.Endpoint, _ EndpointRuntime) (*Resource, error) {
		return &Resource{Endpoint: value, Backend: fsx.Local{}, Path: value.Path, Close: func() error { closed.Add(1); return nil }}, nil
	}
	definition := webDefinition(delivery.RoutePathToWeb, root)
	definition.Archive = true
	record := webRecord(testDeliveryID, delivery.RoutePathToWeb, delivery.DefaultPolicy())
	if err := host.Register(context.Background(), record, marshalDefinition(t, definition)); err != nil {
		t.Fatal(err)
	}
	temporaries := host.OwnedTemps(record.ID)
	if len(temporaries) != 1 || temporaries[0].OwnerID != record.ID || temporaries[0].Path == "" {
		t.Fatalf("owned temporaries=%+v", temporaries)
	}
	temporaries[0].Path = "changed"
	if host.OwnedTemps(record.ID)[0].Path == "changed" || host.OwnedTemps(delivery.ID("00000000-0000-4000-8000-000000000099")) != nil {
		t.Fatal("owned temporary snapshot is not isolated")
	}
	preparedPath := host.byID[record.ID].resource.Path
	base := "/d/" + definition.Token
	response := perform(host, http.MethodGet, base+"/api/v1/meta", nil, nil)
	wantName := filepath.Base(root) + ".tar.gz"
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), wantName) || strings.Contains(response.Body.String(), "file.txt") {
		t.Fatalf("archive metadata=%d %q", response.Code, response.Body.String())
	}
	response = perform(host, http.MethodGet, base+"/api/v1/download", nil, nil)
	if response.Code != http.StatusOK || !strings.Contains(response.Header().Get("Content-Disposition"), wantName) {
		t.Fatalf("archive download=%d headers=%v", response.Code, response.Header())
	}
	reader, err := gzip.NewReader(bytes.NewReader(response.Body.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	header, err := tar.NewReader(reader).Next()
	if err != nil || header.Name != filepath.Base(root) {
		t.Fatalf("archive root=%v err=%v", header, err)
	}
	_ = reader.Close()
	if err := host.Stop(context.Background(), record.ID); err != nil || closed.Load() != 1 {
		t.Fatalf("archive stop=%v closed=%d", err, closed.Load())
	}
	if _, err := os.Stat(preparedPath); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("prepared archive remained: %v", err)
	}
}

func TestAuthenticationIsolationAndStop(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "private.txt"), []byte("private"), 0o600); err != nil {
		t.Fatal(err)
	}
	stopped := make(chan delivery.ID, 1)
	var host *Host
	host = newTestHost(t, func(id delivery.ID) {
		stopped <- id
		_ = host.Stop(context.Background(), id)
	})
	configured := delivery.DefaultPolicy()
	configured.Auth = delivery.AuthPassword
	definition := webDefinition(delivery.RoutePathToWeb, root)
	definition.Credentials.Password = []byte("secret")
	record := webRecord(testDeliveryID, delivery.RoutePathToWeb, configured)
	if err := host.Register(context.Background(), record, marshalDefinition(t, definition)); err != nil {
		t.Fatal(err)
	}
	base := "/d/" + definition.Token
	response := perform(host, http.MethodGet, base+"/", nil, nil)
	if response.Code != http.StatusOK || strings.Contains(response.Body.String(), "private.txt") {
		t.Fatalf("login shell leaked metadata: %d %q", response.Code, response.Body.String())
	}
	response = perform(host, http.MethodGet, base+"/api/v1/meta", nil, map[string]string{"X-Forwarded-For": "127.0.0.1"})
	if response.Code != http.StatusUnauthorized || strings.Contains(response.Body.String(), "private.txt") {
		t.Fatalf("unauthorized metadata leaked: %d %q", response.Code, response.Body.String())
	}
	response = perform(host, http.MethodPost, base+"/api/v1/session", strings.NewReader(`{"password":"secret"}`), map[string]string{"Content-Type": "application/json"})
	if response.Code != http.StatusOK || len(response.Result().Cookies()) != 1 || !strings.Contains(response.Body.String(), "csrf") {
		t.Fatalf("login=%d %q cookies=%v", response.Code, response.Body.String(), response.Result().Cookies())
	}
	cookie := response.Result().Cookies()[0]
	request := httptest.NewRequest(http.MethodGet, base+"/api/v1/meta", nil)
	request.RemoteAddr = "192.0.2.10:99"
	request.AddCookie(cookie)
	metadata := httptest.NewRecorder()
	host.routes().ServeHTTP(metadata, request)
	if metadata.Code != http.StatusOK || !strings.Contains(metadata.Body.String(), "private.txt") {
		t.Fatalf("authenticated metadata=%d %q", metadata.Code, metadata.Body.String())
	}

	basicPolicy := delivery.DefaultPolicy()
	basicPolicy.Auth = delivery.AuthBasic
	basicPolicy.AuthAttempts = 1
	basicPolicy.AuthFailAction = delivery.AuthFailStop
	basicDefinition := webDefinition(delivery.RoutePathToWeb, root)
	basicDefinition.Token = fixedToken(4)
	basicDefinition.Credentials = policy.Credentials{BasicUsername: "user", BasicPassword: []byte("pass")}
	basicID := delivery.ID("00000000-0000-4000-8000-000000000034")
	if err := host.Register(context.Background(), webRecord(basicID, delivery.RoutePathToWeb, basicPolicy), marshalDefinition(t, basicDefinition)); err != nil {
		t.Fatal(err)
	}
	response = perform(host, http.MethodGet, "/d/"+basicDefinition.Token+"/api/v1/meta", nil, map[string]string{"Authorization": "Basic dXNlcjp3cm9uZw=="})
	if response.Code != http.StatusGone || response.Header().Get("WWW-Authenticate") == "" {
		t.Fatalf("stop auth=%d headers=%v", response.Code, response.Header())
	}
	select {
	case id := <-stopped:
		if id != basicID {
			t.Fatalf("stopped wrong delivery: %s", id)
		}
	case <-time.After(time.Second):
		t.Fatal("missing stop callback")
	}
	if response := perform(host, http.MethodGet, "/d/"+basicDefinition.Token+"/", nil, nil); response.Code != http.StatusNotFound {
		t.Fatalf("stopped route=%d", response.Code)
	}
}

func TestPathAndStreamHelpers(t *testing.T) {
	if safeUploadName("") || safeUploadName(".") || safeUploadName("..") || safeUploadName("a/b") || safeUploadName("a\\b") || !safeUploadName("résumé.txt") {
		t.Fatal("unexpected upload-name classification")
	}
	configured := delivery.DefaultPolicy()
	configured.DownloadRate = delivery.Limit{Unlimited: true}
	engine, err := policy.NewEngine(testDeliveryID, configured, policy.Credentials{}, webDependencies())
	if err != nil {
		t.Fatal(err)
	}
	authorization, err := engine.Authorize(context.Background(), policy.Request{Peer: &net.TCPAddr{IP: netip.MustParseAddr("192.0.2.1").AsSlice(), Port: 1}})
	if err != nil {
		t.Fatal(err)
	}
	defer authorization.Release()
	var output bytes.Buffer
	if err := copyRateLimited(context.Background(), &output, strings.NewReader("value"), authorization, policy.Download); err != nil || output.String() != "value" {
		t.Fatalf("copy helper=%q %v", output.String(), err)
	}
	writer := &limitedWriter{ctx: context.Background(), destination: &output, authorization: authorization, direction: policy.Download}
	if count, err := writer.Write([]byte("!")); err != nil || count != 1 {
		t.Fatalf("limited writer=%d %v", count, err)
	}
	if objectType(fakeInfo{directory: true}) != "directory" || objectType(fakeInfo{}) != "file" {
		t.Fatal("unexpected object type")
	}
}
