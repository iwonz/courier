package worker

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/endpoint"
	"github.com/iwonz/courier/internal/fsx"
	"github.com/iwonz/courier/internal/policy"
	"github.com/iwonz/courier/internal/sshx"
	"github.com/iwonz/courier/internal/webdelivery"
)

type fakeDeliveryHost struct {
	mutex        sync.Mutex
	registerErr  error
	updateErr    error
	stopErr      error
	serveErr     error
	registered   []delivery.ID
	temps        map[delivery.ID][]delivery.OwnedTemp
	validated    []delivery.ID
	updated      []delivery.ID
	canceled     []delivery.ID
	stopped      []delivery.ID
	closed       chan struct{}
	closeOnce    sync.Once
	serveStarted chan struct{}
	onServe      func()
}

func (host *fakeDeliveryHost) Register(_ context.Context, item delivery.Delivery, _ json.RawMessage) error {
	host.mutex.Lock()
	defer host.mutex.Unlock()
	host.registered = append(host.registered, item.ID)
	return host.registerErr
}
func (host *fakeDeliveryHost) OwnedTemps(id delivery.ID) []delivery.OwnedTemp {
	host.mutex.Lock()
	defer host.mutex.Unlock()
	return append([]delivery.OwnedTemp(nil), host.temps[id]...)
}
func (host *fakeDeliveryHost) PreparePolicy(_ context.Context, id delivery.ID, _ uint64, _ delivery.Policy) (func(), func(), error) {
	host.mutex.Lock()
	host.validated = append(host.validated, id)
	err := host.updateErr
	host.mutex.Unlock()
	if err != nil {
		return nil, nil, err
	}
	var once sync.Once
	return func() {
			once.Do(func() {
				host.mutex.Lock()
				host.updated = append(host.updated, id)
				host.mutex.Unlock()
			})
		}, func() {
			once.Do(func() {
				host.mutex.Lock()
				host.canceled = append(host.canceled, id)
				host.mutex.Unlock()
			})
		}, nil
}
func (host *fakeDeliveryHost) Stop(_ context.Context, id delivery.ID) error {
	host.mutex.Lock()
	defer host.mutex.Unlock()
	host.stopped = append(host.stopped, id)
	return host.stopErr
}
func (host *fakeDeliveryHost) Serve(net.Listener) error {
	if host.serveStarted != nil {
		close(host.serveStarted)
	}
	if host.onServe != nil {
		host.onServe()
	}
	if host.closed != nil {
		<-host.closed
	}
	return host.serveErr
}
func (host *fakeDeliveryHost) Close(context.Context) error {
	if host.closed != nil {
		host.closeOnce.Do(func() { close(host.closed) })
	}
	return nil
}

func TestRuntimeDeliveryHostLifecycle(t *testing.T) {
	runtime, store := newRuntimeForTest(t)
	if err := runtime.AttachHost(nil); err == nil {
		t.Fatal("nil host attached")
	}
	host := &fakeDeliveryHost{}
	if err := runtime.AttachHost(host); err != nil {
		t.Fatal(err)
	}
	if err := runtime.AttachHost(&fakeDeliveryHost{}); err == nil {
		t.Fatal("duplicate host attached")
	}
	item := validWorkerDelivery()
	temporary := delivery.OwnedTemp{
		ID: delivery.ID("00000000-0000-4000-8000-000000000081"), OwnerID: item.ID,
		Location: delivery.TempLocal, Path: "/tmp/.courier-123.tar.gz", CreatedAt: workerTime,
	}
	host.temps = map[delivery.ID][]delivery.OwnedTemp{item.ID: {temporary}}
	definition := json.RawMessage(`{"version":1}`)
	host.registerErr = errors.New("register")
	if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: item, RuntimeDefinition: definition}, nil); err == nil {
		t.Fatal("host registration error ignored")
	}
	host.registerErr = nil
	if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: item, RuntimeDefinition: definition}, nil); err != nil {
		t.Fatal(err)
	}
	if snapshot, err := store.Load(context.Background()); err != nil || len(snapshot.OwnedTemps) != 1 {
		t.Fatalf("registered temporaries=%+v err=%v", snapshot.OwnedTemps, err)
	}
	next := item.Policy
	next.Version++
	host.updateErr = errors.New("update")
	if err := runtime.UpdatePolicy(context.Background(), UpdatePolicyRequest{DeliveryID: item.ID, ExpectedVersion: item.Policy.Version, Policy: next}); err == nil {
		t.Fatal("host update error ignored")
	}
	host.updateErr = nil
	if err := runtime.UpdatePolicy(context.Background(), UpdatePolicyRequest{DeliveryID: item.ID, ExpectedVersion: item.Policy.Version, Policy: next}); err != nil {
		t.Fatal(err)
	}
	if len(host.validated) != 2 || len(host.updated) != 1 {
		t.Fatalf("unexpected policy phases: validated=%v updated=%v", host.validated, host.updated)
	}
	host.stopErr = errors.New("stop")
	if stopped, err := runtime.StopDelivery(context.Background(), item.ID, delivery.ReasonStopped); err == nil || !stopped {
		t.Fatal("host stop error ignored")
	}
	if stopped, err := runtime.StopDelivery(context.Background(), item.ID, delivery.ReasonStopped); err != nil || stopped {
		t.Fatalf("repeated stop=%v err=%v", stopped, err)
	}
	if snapshot, err := store.Load(context.Background()); err != nil || len(snapshot.OwnedTemps) != 1 {
		t.Fatalf("failed cleanup metadata=%+v err=%v", snapshot.OwnedTemps, err)
	}

	withoutHost, _ := newRuntimeForTest(t)
	second := validWorkerDelivery()
	second.ID = delivery.ID("00000000-0000-4000-8000-000000000077")
	if _, err := withoutHost.Register(context.Background(), RegisterRequest{Delivery: second, RuntimeDefinition: definition}, nil); err == nil {
		t.Fatal("runtime definition accepted without host")
	}

	rollback, rollbackStore := newRuntimeForTest(t)
	rollbackHost := &fakeDeliveryHost{}
	if err := rollback.AttachHost(rollbackHost); err != nil {
		t.Fatal(err)
	}
	if err := rollbackStore.Close(); err != nil {
		t.Fatal(err)
	}
	third := validWorkerDelivery()
	third.ID = delivery.ID("00000000-0000-4000-8000-000000000078")
	if _, err := rollback.Register(context.Background(), RegisterRequest{Delivery: third, RuntimeDefinition: definition}, nil); err == nil || len(rollbackHost.stopped) != 1 {
		t.Fatalf("expected host rollback: %v stopped=%v", err, rollbackHost.stopped)
	}

	updateRollback, updateStore := newRuntimeForTest(t)
	updateHost := &fakeDeliveryHost{}
	if err := updateRollback.AttachHost(updateHost); err != nil {
		t.Fatal(err)
	}
	fourth := validWorkerDelivery()
	fourth.ID = delivery.ID("00000000-0000-4000-8000-000000000080")
	if _, err := updateRollback.Register(context.Background(), RegisterRequest{Delivery: fourth, RuntimeDefinition: definition}, nil); err != nil {
		t.Fatal(err)
	}
	if err := updateStore.Close(); err != nil {
		t.Fatal(err)
	}
	fourthPolicy := fourth.Policy
	fourthPolicy.Version++
	if err := updateRollback.UpdatePolicy(context.Background(), UpdatePolicyRequest{DeliveryID: fourth.ID, ExpectedVersion: fourth.Policy.Version, Policy: fourthPolicy}); err == nil || len(updateHost.canceled) != 1 || len(updateHost.updated) != 0 {
		t.Fatalf("expected prepared policy rollback: %v canceled=%v updated=%v", err, updateHost.canceled, updateHost.updated)
	}

	invalidTempRuntime, _ := newRuntimeForTest(t)
	invalidTempHost := &fakeDeliveryHost{}
	if err := invalidTempRuntime.AttachHost(invalidTempHost); err != nil {
		t.Fatal(err)
	}
	fifth := validWorkerDelivery()
	fifth.ID = delivery.ID("00000000-0000-4000-8000-000000000082")
	invalidTempHost.temps = map[delivery.ID][]delivery.OwnedTemp{fifth.ID: {{
		ID: delivery.ID("00000000-0000-4000-8000-000000000083"), OwnerID: item.ID,
		Location: delivery.TempLocal, Path: "/tmp/.courier-456.tar.gz", CreatedAt: workerTime,
	}}}
	if _, err := invalidTempRuntime.Register(context.Background(), RegisterRequest{Delivery: fifth, RuntimeDefinition: definition}, nil); !errors.Is(err, delivery.ErrNotFound) || len(invalidTempHost.stopped) != 1 {
		t.Fatalf("invalid owned temporary accepted: %v stopped=%v", err, invalidTempHost.stopped)
	}
	_ = store
}

func TestRuntimeOwnedTemporaryCleanup(t *testing.T) {
	runtime, store := newRuntimeForTest(t)
	host := &fakeDeliveryHost{}
	if err := runtime.AttachHost(host); err != nil {
		t.Fatal(err)
	}
	item := validWorkerDelivery()
	temporary := delivery.OwnedTemp{
		ID: delivery.ID("00000000-0000-4000-8000-000000000084"), OwnerID: item.ID,
		Location: delivery.TempLocal, Path: "/tmp/.courier-789.tar.gz", CreatedAt: workerTime,
	}
	host.temps = map[delivery.ID][]delivery.OwnedTemp{item.ID: {temporary}}
	if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: item, RuntimeDefinition: json.RawMessage(`{}`)}, nil); err != nil {
		t.Fatal(err)
	}
	if stopped, err := runtime.StopDelivery(context.Background(), item.ID, delivery.ReasonStopped); err != nil || !stopped {
		t.Fatalf("stop=%v err=%v", stopped, err)
	}
	if snapshot, err := store.Load(context.Background()); err != nil || len(snapshot.OwnedTemps) != 0 {
		t.Fatalf("temporary remained: %+v err=%v", snapshot.OwnedTemps, err)
	}

	runtime, store = newRuntimeForTest(t)
	host = &fakeDeliveryHost{}
	if err := runtime.AttachHost(host); err != nil {
		t.Fatal(err)
	}
	item = validWorkerDelivery()
	temporary.OwnerID = item.ID
	host.temps = map[delivery.ID][]delivery.OwnedTemp{item.ID: {temporary}}
	if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: item, RuntimeDefinition: json.RawMessage(`{}`)}, nil); err != nil {
		t.Fatal(err)
	}
	current, err := store.Load(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Update(context.Background(), current.Revision, func(registry *delivery.Registry) error {
		return registry.RemoveOwnedTemp(temporary.ID)
	}); err != nil {
		t.Fatal(err)
	}
	if err := runtime.StopServer(context.Background(), delivery.ReasonStopped); !errors.Is(err, delivery.ErrNotFound) {
		t.Fatalf("missing temporary cleanup metadata=%v", err)
	}
}

func TestRuntimeHostServerStop(t *testing.T) {
	runtime, _ := newRuntimeForTest(t)
	host := &fakeDeliveryHost{}
	if err := runtime.AttachHost(host); err != nil {
		t.Fatal(err)
	}
	first := validWorkerDelivery()
	second := first
	second.ID = delivery.ID("00000000-0000-4000-8000-000000000079")
	definition := json.RawMessage(`{"version":1}`)
	for _, item := range []delivery.Delivery{first, second} {
		if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: item, RuntimeDefinition: definition}, nil); err != nil {
			t.Fatal(err)
		}
	}
	host.stopErr = errors.New("stop")
	if err := runtime.StopServer(context.Background(), delivery.ReasonStopped); err == nil {
		t.Fatal("host server-stop error ignored")
	}
	if err := runtime.StopServer(context.Background(), delivery.ReasonStopped); err != nil || len(host.stopped) != 2 {
		t.Fatalf("repeated stop failed: %v calls=%v", err, host.stopped)
	}
}

func TestRuntimeHostAttachmentAfterStart(t *testing.T) {
	runtime, _ := newRuntimeForTest(t)
	runtime.listener = &scriptedListener{}
	if err := runtime.AttachHost(&fakeDeliveryHost{}); err == nil {
		t.Fatal("host attached after control start")
	}
	runtime.listener = nil
	if err := runtime.StopServer(context.Background(), delivery.ReasonStopped); err != nil {
		t.Fatal(err)
	}
}

func TestRuntimeDefinitionValidationAndCoordinatorCallback(t *testing.T) {
	item := validWorkerDelivery()
	request := RegisterRequest{Delivery: item, RuntimeDefinition: json.RawMessage("{")}
	if err := request.Validate(); err == nil {
		t.Fatal("malformed definition accepted")
	}
	request.RuntimeDefinition = json.RawMessage(`"` + strings.Repeat("x", MaxRuntimeDefinition) + `"`)
	if err := request.Validate(); err == nil {
		t.Fatal("oversized definition accepted")
	}
	acquire := AcquireRequest{Bind: "127.0.0.1:8080", Compatibility: "http-v1", Route: delivery.RoutePathToWeb, Policy: delivery.DefaultPolicy(), At: workerTime, RuntimeDefinition: json.RawMessage("{")}
	if err := acquire.Validate(); err == nil {
		t.Fatal("malformed acquisition definition accepted")
	}
	acquire.RuntimeDefinition = json.RawMessage(`"` + strings.Repeat("x", MaxRuntimeDefinition) + `"`)
	if err := acquire.Validate(); err == nil {
		t.Fatal("oversized acquisition definition accepted")
	}

	store, err := delivery.OpenStore(filepath.Join(t.TempDir(), "state"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	client := Client{Endpoint: "control", ServerID: workerServerID, Compatibility: "http-v1"}
	coordinator := &Coordinator{
		Store: store, StateDirectory: store.Directory(), Locks: &fakeBindLocker{},
		Launch: func(_ context.Context, request LaunchRequest) (Client, error) {
			client = Client{Endpoint: request.ControlEndpoint, ServerID: request.ServerID, Compatibility: request.Compatibility}
			return client, nil
		},
		Hello: func(context.Context, Client) error { return nil },
		RegisterDefinition: func(_ context.Context, got Client, deliveryItem delivery.Delivery, foreground bool, definition json.RawMessage) (*Lease, error) {
			if got != client || foreground || string(definition) != `{}` || !deliveryItem.ID.Valid() {
				t.Fatalf("unexpected registration: %#v %v %s", got, foreground, definition)
			}
			return nil, nil
		},
	}
	acquire.RuntimeDefinition = json.RawMessage(`{}`)
	result, err := coordinator.Acquire(context.Background(), acquire)
	if err != nil || !result.DeliveryID.Valid() {
		t.Fatalf("acquire=%#v err=%v", result, err)
	}
}

func TestRunProcessHostFactoryPaths(t *testing.T) {
	directory := shortWorkerTempDir(t, "courier-host-process-")
	probe, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	bind := probe.Addr().String()
	_ = probe.Close()
	config := validProcessConfig(t, directory, bind)
	originalFactory := newProcessHost
	t.Cleanup(func() { newProcessHost = originalFactory })
	newProcessHost = func(func(delivery.ID)) (DeliveryHost, error) { return nil, errors.New("host factory") }
	if err := runProcess(context.Background(), config); err == nil {
		t.Fatal("host factory error ignored")
	}

	directory = shortWorkerTempDir(t, "courier-nil-host-")
	probe, _ = net.Listen("tcp", "127.0.0.1:0")
	bind = probe.Addr().String()
	_ = probe.Close()
	config = validProcessConfig(t, directory, bind)
	newProcessHost = func(func(delivery.ID)) (DeliveryHost, error) { return nil, nil }
	if err := runProcess(context.Background(), config); err == nil {
		t.Fatal("nil host attachment accepted")
	}

	directory = shortWorkerTempDir(t, "courier-live-host-")
	probe, _ = net.Listen("tcp", "127.0.0.1:0")
	bind = probe.Addr().String()
	_ = probe.Close()
	config = validProcessConfig(t, directory, bind)
	ctx, cancel := context.WithCancel(context.Background())
	host := &fakeDeliveryHost{closed: make(chan struct{}), serveStarted: make(chan struct{})}
	newProcessHost = func(stop func(delivery.ID)) (DeliveryHost, error) {
		host.onServe = func() { stop(delivery.ID("00000000-0000-4000-8000-000000000099")); cancel() }
		return host, nil
	}
	if err := runProcess(ctx, config); err != nil && !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

func TestWorkerEnvironmentAllowsAgentWithoutTokens(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "/private/agent.sock")
	t.Setenv("NPM_TOKEN", "must-not-pass")
	environment := workerEnvironment("config")
	joined := strings.Join(environment, "\n")
	if !strings.Contains(joined, "SSH_AUTH_SOCK=/private/agent.sock") || strings.Contains(joined, "must-not-pass") {
		t.Fatalf("unsafe worker environment: %s", joined)
	}
	_ = os.ErrNotExist
}

func TestNewWorkerWebHostDependenciesAndOpeners(t *testing.T) {
	if accepted, err := approveWorkerHelper(context.Background(), "approved during CLI preflight"); err != nil || !accepted {
		t.Fatalf("worker helper approval=%v err=%v", accepted, err)
	}
	originalHome, originalUser, originalLoad := workerHome, workerUser, workerLoadConfig
	originalOpen, originalDetect, originalRooted, originalResource := workerOpenSSH, workerDetectSSH, workerOpenRooted, workerRemoteResource
	t.Cleanup(func() {
		workerHome, workerUser, workerLoadConfig = originalHome, originalUser, originalLoad
		workerOpenSSH, workerDetectSSH, workerRemoteResource = originalOpen, originalDetect, originalResource
		workerOpenRooted = originalRooted
	})
	workerHome = func() (string, error) { return "", errors.New("home") }
	if _, err := newWorkerWebHost(nil); err == nil {
		t.Fatal("home error ignored")
	}
	workerHome = func() (string, error) { return t.TempDir(), nil }
	workerLoadConfig = func(string, string) (*sshx.Config, error) { return nil, errors.New("config") }
	if _, err := newWorkerWebHost(nil); err == nil {
		t.Fatal("config error ignored")
	}
	home := t.TempDir()
	workerHome = func() (string, error) { return home, nil }
	workerLoadConfig = func(string, string) (*sshx.Config, error) { return nil, os.ErrNotExist }
	workerUser = func() (*user.User, error) { return nil, errors.New("user") }
	host, err := newWorkerWebHost(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer host.Close(context.Background())
	root := t.TempDir()
	workerOpenRooted = func(string) (*fsx.RootedLocal, string, error) { return nil, "", errors.New("rooted") }
	failingDefinition := webdelivery.Definition{Version: webdelivery.DefinitionVersion, Token: fixedWorkerToken(t), Source: "web://", Destination: root}
	failingData, _ := json.Marshal(failingDefinition)
	failingItem := validWorkerDelivery()
	failingItem.Route = delivery.RouteWebToPath
	if err := host.Register(context.Background(), failingItem, failingData); err == nil {
		t.Fatal("rooted open error ignored")
	}
	workerOpenRooted = originalRooted
	definition := webdelivery.Definition{Version: webdelivery.DefinitionVersion, Token: fixedWorkerToken(t), Source: "web://", Destination: root}
	data, _ := json.Marshal(definition)
	item := validWorkerDelivery()
	item.Route = delivery.RouteWebToPath
	if err := host.Register(context.Background(), item, data); err != nil {
		t.Fatal(err)
	}
	if err := host.Stop(context.Background(), item.ID); err != nil {
		t.Fatal(err)
	}

	workerUser = func() (*user.User, error) { return &user.User{Username: "tester"}, nil }
	workerLoadConfig = func(string, string) (*sshx.Config, error) { return sshx.EmptyConfig(home), nil }
	workerOpenSSH = func(context.Context, sshx.Factory, endpoint.Endpoint) (*sshx.Connection, error) {
		return nil, errors.New("ssh")
	}
	host, err = newWorkerWebHost(nil)
	if err != nil {
		t.Fatal(err)
	}
	remoteDefinition := webdelivery.Definition{Version: webdelivery.DefinitionVersion, Token: fixedWorkerToken(t), Source: "host:/source", Destination: "web://"}
	remoteData, _ := json.Marshal(remoteDefinition)
	remoteItem := validWorkerDelivery()
	remoteItem.ID = delivery.ID("00000000-0000-4000-8000-000000000088")
	remoteItem.Route = delivery.RoutePathToWeb
	if err := host.Register(context.Background(), remoteItem, remoteData); err == nil {
		t.Fatal("SSH open error ignored")
	}
	connection := &sshx.Connection{Target: sshx.Target{Host: "canonical", User: "tester"}}
	workerOpenSSH = func(context.Context, sshx.Factory, endpoint.Endpoint) (*sshx.Connection, error) {
		return connection, nil
	}
	workerDetectSSH = func(context.Context, sshx.Runner) (sshx.Platform, sshx.Archiver, error) {
		return sshx.Platform{}, sshx.Archiver{}, errors.New("detect")
	}
	if err := host.Register(context.Background(), remoteItem, remoteData); err == nil {
		t.Fatal("SSH detection error ignored")
	}
	workerDetectSSH = func(context.Context, sshx.Runner) (sshx.Platform, sshx.Archiver, error) {
		return sshx.Platform{}, sshx.Archiver{}, nil
	}
	remoteDefinition.Endpoint = webdelivery.EndpointRuntime{Credentials: []webdelivery.EndpointCredential{{Prompt: "worker prompt", Secret: []byte("worker secret")}}, AllowHelper: true}
	remoteData, _ = json.Marshal(remoteDefinition)
	workerOpenSSH = func(_ context.Context, factory sshx.Factory, _ endpoint.Endpoint) (*sshx.Connection, error) {
		if factory.SFTPFallback == nil {
			t.Fatal("pre-approved worker helper fallback is missing")
		}
		secret, promptErr := factory.Prompt("worker prompt")
		if promptErr != nil || string(secret) != "worker secret" {
			t.Fatalf("worker prompt=%q err=%v", secret, promptErr)
		}
		return connection, nil
	}
	workerRemoteResource = func(connection *sshx.Connection, value endpoint.Endpoint) *webdelivery.Resource {
		return &webdelivery.Resource{Endpoint: value, Backend: fsx.Local{}, Path: root, Close: connection.Close}
	}
	if err := host.Register(context.Background(), remoteItem, remoteData); err != nil {
		t.Fatal(err)
	}
	if err := host.Stop(context.Background(), remoteItem.ID); err != nil {
		t.Fatal(err)
	}
	resource := originalResource(connection, endpoint.Endpoint{Host: "alias", User: "before", Path: "/path"})
	if resource.Endpoint.Host != "canonical" || resource.Endpoint.User != "tester" || resource.Path != "/path" {
		t.Fatalf("unexpected remote resource: %#v", resource)
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	_, _ = originalOpen(canceled, sshx.Factory{}, endpoint.Endpoint{Remote: true, Host: "host", Path: "/path"})

	workerOpenSSH, workerDetectSSH, workerRemoteResource = originalOpen, originalDetect, originalResource
	workerOpenRooted = originalRooted
	stopped := make(chan delivery.ID, 1)
	callbackHost, err := newWorkerWebHost(func(id delivery.ID) { stopped <- id })
	if err != nil {
		t.Fatal(err)
	}
	defer callbackHost.Close(context.Background())
	callbackDefinition := webdelivery.Definition{
		Version: webdelivery.DefinitionVersion, Token: fixedWorkerToken(t), Source: root, Destination: "web://",
		Credentials: policy.Credentials{BasicUsername: "user", BasicPassword: []byte("pass")},
	}
	callbackData, _ := json.Marshal(callbackDefinition)
	callbackItem := validWorkerDelivery()
	callbackItem.ID = delivery.ID("00000000-0000-4000-8000-000000000089")
	callbackItem.Route = delivery.RoutePathToWeb
	callbackItem.Policy.Auth = delivery.AuthBasic
	callbackItem.Policy.AuthAttempts = 1
	callbackItem.Policy.AuthFailAction = delivery.AuthFailStop
	if err := callbackHost.Register(context.Background(), callbackItem, callbackData); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/d/"+callbackDefinition.Token+"/api/v1/meta", nil)
	request.RemoteAddr = "192.0.2.1:1"
	response := httptest.NewRecorder()
	callbackHost.(*webdelivery.Host).Handler().ServeHTTP(response, request)
	if response.Code != http.StatusGone || <-stopped != callbackItem.ID {
		t.Fatalf("callback status=%d", response.Code)
	}
}

func fixedWorkerToken(t *testing.T) string {
	t.Helper()
	token, err := webdelivery.NewToken(strings.NewReader(strings.Repeat("x", webdelivery.ResourceTokenBytes)))
	if err != nil {
		t.Fatal(err)
	}
	return token
}
