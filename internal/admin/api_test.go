package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/control"
	"github.com/iwonz/courier/internal/delivery"
)

const (
	adminServerID   = delivery.ID("00000000-0000-4000-8000-000000000601")
	adminDeliveryID = delivery.ID("00000000-0000-4000-8000-000000000602")
)

var adminTestTime = time.Date(2026, 9, 14, 1, 2, 3, 4, time.UTC)

type fakeController struct {
	views       []control.ServerView
	listErr     error
	stopResult  control.StopResult
	stopErr     error
	updateErr   error
	stopRequest control.StopRequest
	updatedID   delivery.ID
	expected    uint64
	policy      delivery.Policy
}

func (fake *fakeController) List(context.Context) ([]control.ServerView, error) {
	return fake.views, fake.listErr
}

func (fake *fakeController) Stop(_ context.Context, request control.StopRequest) (control.StopResult, error) {
	fake.stopRequest = request
	return fake.stopResult, fake.stopErr
}

func (fake *fakeController) UpdatePolicy(_ context.Context, id delivery.ID, expected uint64, policy delivery.Policy) error {
	fake.updatedID, fake.expected, fake.policy = id, expected, policy
	return fake.updateErr
}

func adminViews() []control.ServerView {
	policy := delivery.DefaultPolicy()
	policy.Auth = delivery.AuthPassword
	return []control.ServerView{
		{
			Live: true,
			Server: delivery.Server{
				ID: adminServerID, Bind: "127.0.0.1:8080", ControlEndpoint: "secret.sock", Compatibility: "secret-compatible",
				ProcessID: 42, State: delivery.StateActive, StartedAt: adminTestTime, UpdatedAt: adminTestTime,
			},
			Deliveries: []delivery.Delivery{{
				ID: adminDeliveryID, ServerID: adminServerID, Route: delivery.RoutePathToWeb, Source: "./data", Destination: "web://",
				State: delivery.StateActive, Policy: policy, Counters: delivery.CounterSnapshot{Read: 1, Sent: 2, Confirmed: 3}, CreatedAt: adminTestTime, UpdatedAt: adminTestTime,
			}},
		},
		{Server: delivery.Server{ID: delivery.ID("00000000-0000-4000-8000-000000000603"), Bind: "127.0.0.1:8081", ProcessID: 43, State: delivery.StateFailed, StartedAt: adminTestTime, UpdatedAt: adminTestTime}},
	}
}

func newTestAPI(t *testing.T, controller Controller) *API {
	t.Helper()
	api, err := NewAPI(APIConfig{Host: "127.0.0.1:9090", Control: controller, EventInterval: time.Millisecond, MaxSubscribers: 1})
	if err != nil {
		t.Fatal(err)
	}
	return api
}

func adminRequest(method, target string, body io.Reader) *http.Request {
	request := httptest.NewRequest(method, target, body)
	request.Host = "127.0.0.1:9090"
	request.RemoteAddr = "127.0.0.1:54321"
	if method != http.MethodGet && method != http.MethodHead {
		request.Header.Set("Origin", "http://127.0.0.1:9090")
	}
	return request
}

func TestAPIInventoryAndAssets(t *testing.T) {
	controller := &fakeController{views: adminViews()}
	api := newTestAPI(t, controller)
	response := httptest.NewRecorder()
	api.ServeHTTP(response, adminRequest(http.MethodGet, "/api/v1/servers", nil))
	if response.Code != http.StatusOK || response.Header().Get("Access-Control-Allow-Origin") != "" || response.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("code=%d headers=%v", response.Code, response.Header())
	}
	body := response.Body.String()
	for _, expected := range []string{string(adminServerID), string(adminDeliveryID), `"status":"live"`, `"status":"unreachable"`, `"auth":"password"`, `"confirmed":3`} {
		if !strings.Contains(body, expected) {
			t.Fatalf("missing %q in %s", expected, body)
		}
	}
	for _, forbidden := range []string{"secret.sock", "secret-compatible", "password-value", "resource-token"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("leaked %q in %s", forbidden, body)
		}
	}

	for _, target := range []string{"/", "/assets/admin.js"} {
		response = httptest.NewRecorder()
		api.ServeHTTP(response, adminRequest(http.MethodGet, target, nil))
		if response.Code != http.StatusOK || response.Body.Len() == 0 {
			t.Fatalf("asset %s code=%d body=%q", target, response.Code, response.Body.String())
		}
	}

	if _, err := NewAPI(APIConfig{}); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("invalid API=%v", err)
	}
	defaulted, err := NewAPI(APIConfig{Host: "127.0.0.1:9090", Control: controller})
	if err != nil || defaulted.eventInterval != time.Second || cap(defaulted.subscribers) != 16 {
		t.Fatalf("defaults=%+v err=%v", defaulted, err)
	}
}

func TestAPIGuards(t *testing.T) {
	api := newTestAPI(t, &fakeController{})
	for _, configure := range []func(*http.Request){
		func(request *http.Request) { request.Host = "evil.test" },
		func(request *http.Request) { request.RemoteAddr = "192.0.2.1:80" },
		func(request *http.Request) { request.RemoteAddr = "invalid" },
		func(request *http.Request) { request.Header.Del("Origin") },
		func(request *http.Request) { request.Header.Set("Origin", "https://evil.test") },
	} {
		request := adminRequest(http.MethodPost, "/api/v1/servers/"+string(adminServerID)+"/stop", nil)
		configure(request)
		response := httptest.NewRecorder()
		api.ServeHTTP(response, request)
		if response.Code != http.StatusForbidden {
			t.Fatalf("code=%d", response.Code)
		}
	}
	request := adminRequest(http.MethodHead, "/", nil)
	response := httptest.NewRecorder()
	api.ServeHTTP(response, request)
	if response.Code == http.StatusForbidden {
		t.Fatal("same-host HEAD rejected")
	}
	if !loopbackRemote("[::1]:1234") || loopbackRemote("example.test:80") {
		t.Fatal("loopback parsing mismatch")
	}
}

func TestAPIPolicyUpdates(t *testing.T) {
	controller := &fakeController{}
	api := newTestAPI(t, controller)
	policy := delivery.DefaultPolicy()
	policy.Version = 2
	payload, err := json.Marshal(policyUpdate{ExpectedVersion: 1, Policy: policy})
	if err != nil {
		t.Fatal(err)
	}
	request := adminRequest(http.MethodPut, "/api/v1/deliveries/"+string(adminDeliveryID)+"/policy", bytes.NewReader(payload))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	api.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent || controller.updatedID != adminDeliveryID || controller.expected != 1 || controller.policy.Version != 2 {
		t.Fatalf("code=%d update=%s/%d/%+v", response.Code, controller.updatedID, controller.expected, controller.policy)
	}

	for _, test := range []struct {
		name        string
		id          string
		contentType string
		body        string
		controlErr  error
		want        int
	}{
		{"bad id", "bad", "application/json", string(payload), nil, http.StatusBadRequest},
		{"content type", string(adminDeliveryID), "text/plain", string(payload), nil, http.StatusBadRequest},
		{"malformed", string(adminDeliveryID), "application/json", "{", nil, http.StatusBadRequest},
		{"unknown field", string(adminDeliveryID), "application/json", `{"unknown":true}`, nil, http.StatusBadRequest},
		{"trailing", string(adminDeliveryID), "application/json", string(payload) + `{}`, nil, http.StatusBadRequest},
		{"too large", string(adminDeliveryID), "application/json", strings.Repeat(" ", maxAPIRequestBytes+1), nil, http.StatusBadRequest},
		{"conflict", string(adminDeliveryID), "application/json", string(payload), delivery.ErrRevisionConflict, http.StatusConflict},
		{"missing", string(adminDeliveryID), "application/json", string(payload), delivery.ErrNotFound, http.StatusNotFound},
		{"unavailable", string(adminDeliveryID), "application/json", string(payload), errors.New("private failure"), http.StatusServiceUnavailable},
	} {
		t.Run(test.name, func(t *testing.T) {
			controller.updateErr = test.controlErr
			request := adminRequest(http.MethodPut, "/api/v1/deliveries/"+test.id+"/policy", strings.NewReader(test.body))
			request.Header.Set("Content-Type", test.contentType)
			response := httptest.NewRecorder()
			api.ServeHTTP(response, request)
			if response.Code != test.want || strings.Contains(response.Body.String(), "private failure") {
				t.Fatalf("code=%d body=%q", response.Code, response.Body.String())
			}
		})
	}
}

func TestAPIStopActions(t *testing.T) {
	controller := &fakeController{}
	api := newTestAPI(t, controller)
	for _, test := range []struct {
		path   string
		kind   delivery.TargetKind
		result delivery.TargetKind
		want   int
	}{
		{"/api/v1/deliveries/" + string(adminDeliveryID) + "/stop", delivery.TargetDelivery, delivery.TargetDelivery, http.StatusNoContent},
		{"/api/v1/servers/" + string(adminServerID) + "/stop", delivery.TargetServer, delivery.TargetServer, http.StatusNoContent},
		{"/api/v1/servers/" + string(adminServerID) + "/stop", delivery.TargetServer, delivery.TargetDelivery, http.StatusNotFound},
		{"/api/v1/servers/bad/stop", delivery.TargetServer, delivery.TargetServer, http.StatusBadRequest},
	} {
		controller.stopResult = control.StopResult{Kind: test.result}
		response := httptest.NewRecorder()
		api.ServeHTTP(response, adminRequest(http.MethodPost, test.path, nil))
		if response.Code != test.want {
			t.Fatalf("path=%s code=%d", test.path, response.Code)
		}
		if test.want == http.StatusNoContent && (controller.stopRequest.ID == "" || controller.stopRequest.Kind != test.kind) {
			t.Fatalf("stop request not propagated: %+v", controller.stopRequest)
		}
	}
	controller.stopErr = delivery.ErrNotFound
	response := httptest.NewRecorder()
	api.ServeHTTP(response, adminRequest(http.MethodPost, "/api/v1/servers/"+string(adminServerID)+"/stop", nil))
	if response.Code != http.StatusNotFound {
		t.Fatalf("stop error code=%d", response.Code)
	}
}

type nonStreamingWriter struct{ header http.Header }

func (writer *nonStreamingWriter) Header() http.Header     { return writer.header }
func (*nonStreamingWriter) Write(data []byte) (int, error) { return len(data), nil }
func (*nonStreamingWriter) WriteHeader(int)                {}

type failingStreamWriter struct{ header http.Header }

func (writer *failingStreamWriter) Header() http.Header { return writer.header }
func (*failingStreamWriter) Write([]byte) (int, error)  { return 0, errors.New("write failure") }
func (*failingStreamWriter) WriteHeader(int)            {}
func (*failingStreamWriter) Flush()                     {}

func TestAPIEvents(t *testing.T) {
	controller := &fakeController{views: adminViews()}
	api := newTestAPI(t, controller)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	request := adminRequest(http.MethodGet, "/api/v1/events", nil).WithContext(ctx)
	response := httptest.NewRecorder()
	api.ServeHTTP(response, request)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "event: snapshot") || len(api.subscribers) != 0 {
		t.Fatalf("code=%d body=%q subscribers=%d", response.Code, response.Body.String(), len(api.subscribers))
	}

	nonStreaming := &nonStreamingWriter{header: http.Header{}}
	api.ServeHTTP(nonStreaming, adminRequest(http.MethodGet, "/api/v1/events", nil))

	api.subscribers <- struct{}{}
	response = httptest.NewRecorder()
	api.ServeHTTP(response, adminRequest(http.MethodGet, "/api/v1/events", nil))
	if response.Code != http.StatusTooManyRequests {
		t.Fatalf("subscriber limit=%d", response.Code)
	}
	<-api.subscribers

	controller.listErr = errors.New("unavailable")
	response = httptest.NewRecorder()
	api.ServeHTTP(response, adminRequest(http.MethodGet, "/api/v1/servers", nil))
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("list failure=%d", response.Code)
	}
	ctx, cancel = context.WithCancel(context.Background())
	cancel()
	api.ServeHTTP(httptest.NewRecorder(), adminRequest(http.MethodGet, "/api/v1/events", nil).WithContext(ctx))

	controller.listErr = nil
	ctx, cancel = context.WithCancel(context.Background())
	cancel()
	api.ServeHTTP(&failingStreamWriter{header: http.Header{}}, adminRequest(http.MethodGet, "/api/v1/events", nil).WithContext(ctx))
}
