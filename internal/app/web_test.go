package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/netip"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/endpoint"
	"github.com/iwonz/courier/internal/operation"
	"github.com/iwonz/courier/internal/policy"
	"github.com/iwonz/courier/internal/sshx"
	"github.com/iwonz/courier/internal/webdelivery"
	"github.com/iwonz/courier/internal/worker"
	"github.com/pkg/sftp"
)

const appWebDeliveryID delivery.ID = "00000000-0000-4000-8000-000000000091"

var appWebTime = time.Unix(1_700_000_000, 0).UTC()

type coordinatorFunc func(context.Context, worker.AcquireRequest) (worker.Acquired, error)

func (function coordinatorFunc) Acquire(ctx context.Context, request worker.AcquireRequest) (worker.Acquired, error) {
	return function(ctx, request)
}

type failureWriter struct{}

func (failureWriter) Write([]byte) (int, error) { return 0, errors.New("write") }

func browserPlan(t *testing.T, background bool) operation.Plan {
	t.Helper()
	options := []operation.Option{}
	if background {
		options = append(options, operation.Option{Name: operation.OptionBackground, Value: "true"})
	}
	plan, err := operation.Build(operation.Request{Source: t.TempDir(), Destination: "web://", Options: options})
	if err != nil {
		t.Fatal(err)
	}
	return plan
}

func TestDeliveryCredentialPrompt(t *testing.T) {
	originalTerminal, originalRead := terminalAttached, readTerminalSecret
	t.Cleanup(func() { terminalAttached, readTerminalSecret = originalTerminal, originalRead })
	provider := deliveryCredentialPrompt(nil, io.Discard)
	if credentials, err := provider(context.Background(), operation.AuthNone); err != nil || len(credentials.Password) != 0 {
		t.Fatalf("none credentials=%#v err=%v", credentials, err)
	}
	if _, err := provider(context.Background(), operation.AuthPassword); err == nil {
		t.Fatal("non-terminal password accepted")
	}

	input, err := os.CreateTemp(t.TempDir(), "credentials")
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()
	terminalAttached = func(int) bool { return true }
	readTerminalSecret = func(int) ([]byte, error) { return []byte("secret"), nil }
	if _, err := input.WriteString("user\n"); err != nil {
		t.Fatal(err)
	}
	_, _ = input.Seek(0, io.SeekStart)
	var output bytes.Buffer
	provider = deliveryCredentialPrompt(input, &output)
	credentials, err := provider(context.Background(), operation.AuthBasic)
	if err != nil || credentials.BasicUsername != "user" || string(credentials.BasicPassword) != "secret" || !strings.Contains(output.String(), "Basic username") {
		t.Fatalf("basic=%#v output=%q err=%v", credentials, output.String(), err)
	}
	_, _ = input.Seek(0, io.SeekStart)
	credentials, err = provider(context.Background(), operation.AuthPassword)
	if err != nil || string(credentials.Password) != "secret" {
		t.Fatalf("password=%#v err=%v", credentials, err)
	}

	empty, err := os.CreateTemp(t.TempDir(), "empty")
	if err != nil {
		t.Fatal(err)
	}
	defer empty.Close()
	_, _ = empty.WriteString(" \n")
	_, _ = empty.Seek(0, io.SeekStart)
	if _, err := deliveryCredentialPrompt(empty, io.Discard)(context.Background(), operation.AuthBasic); err == nil {
		t.Fatal("empty Basic username accepted")
	}
	withoutNewline, err := os.CreateTemp(t.TempDir(), "unterminated")
	if err != nil {
		t.Fatal(err)
	}
	defer withoutNewline.Close()
	_, _ = withoutNewline.WriteString("user")
	_, _ = withoutNewline.Seek(0, io.SeekStart)
	if _, err := deliveryCredentialPrompt(withoutNewline, io.Discard)(context.Background(), operation.AuthBasic); err == nil {
		t.Fatal("unterminated Basic username accepted")
	}
	_, _ = input.Seek(0, io.SeekStart)
	if _, err := deliveryCredentialPrompt(input, failureWriter{})(context.Background(), operation.AuthBasic); err == nil {
		t.Fatal("prompt output error ignored")
	}
	_, _ = input.Seek(0, io.SeekStart)
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := provider(canceled, operation.AuthPassword); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled prompt: %v", err)
	}
	readTerminalSecret = func(int) ([]byte, error) { return nil, errors.New("secret") }
	if _, err := provider(context.Background(), operation.AuthPassword); err == nil {
		t.Fatal("secret read error ignored")
	}
}

func TestEndpointCredentialProvider(t *testing.T) {
	originalOpen, originalDetect, originalClose := openSSHConnection, detectSSHPlatform, closeSSHConnection
	t.Cleanup(func() {
		openSSHConnection, detectSSHPlatform, closeSSHConnection = originalOpen, originalDetect, originalClose
	})
	openSSHConnection = func(ctx context.Context, factory sshx.Factory, host, username string) (*sshx.Connection, error) {
		if host == "fail" {
			return nil, errors.New("open")
		}
		if host == "prompt" {
			secret, err := factory.Prompt("Password for user@prompt: ")
			if err != nil {
				return nil, err
			}
			clear(secret)
		}
		if host == "helper" {
			client, closer, err := factory.SFTPFallback(ctx, &sshx.Connection{}, &sshx.CapabilityError{Capability: "sftp", Cause: errors.New("missing")})
			if err != nil || client == nil || closer == nil {
				return nil, errors.Join(err, errors.New("incomplete helper"))
			}
			_ = closer.Close()
		}
		return &sshx.Connection{Target: sshx.Target{Host: host, User: username}}, nil
	}
	detectSSHPlatform = func(_ context.Context, connection sshx.Runner) (sshx.Platform, sshx.Archiver, error) {
		if typed := connection.(*sshx.Connection); typed.Target.Host == "detect" {
			return sshx.Platform{}, sshx.Archiver{}, errors.New("detect")
		}
		return sshx.Platform{OS: "linux", Arch: "amd64"}, sshx.Archiver{}, nil
	}
	closeCalls := 0
	closeSSHConnection = func(connection *sshx.Connection) error {
		closeCalls++
		if connection.Target.Host == "close" {
			return errors.New("close")
		}
		return nil
	}
	provider := endpointCredentialProvider(sshx.Factory{Prompt: func(string) ([]byte, error) { return []byte("secret"), nil }})
	if runtime, err := provider(context.Background(), endpoint.Endpoint{Path: "local"}); err != nil || len(runtime.Credentials) != 0 {
		t.Fatalf("local runtime=%#v err=%v", runtime, err)
	}
	if runtime, err := provider(context.Background(), endpoint.Endpoint{Remote: true, Host: "fail"}); err == nil || len(runtime.Credentials) != 0 {
		t.Fatalf("open failure runtime=%#v err=%v", runtime, err)
	}
	withoutPrompt := endpointCredentialProvider(sshx.Factory{})
	if runtime, err := withoutPrompt(context.Background(), endpoint.Endpoint{Remote: true, Host: "prompt", User: "user"}); err == nil || len(runtime.Credentials) != 0 {
		t.Fatalf("missing prompt runtime=%#v err=%v", runtime, err)
	}
	failingPrompt := endpointCredentialProvider(sshx.Factory{Prompt: func(string) ([]byte, error) { return nil, errors.New("prompt") }})
	if runtime, err := failingPrompt(context.Background(), endpoint.Endpoint{Remote: true, Host: "prompt", User: "user"}); err == nil || len(runtime.Credentials) != 0 {
		t.Fatalf("prompt failure runtime=%#v err=%v", runtime, err)
	}
	runtime, err := provider(context.Background(), endpoint.Endpoint{Remote: true, Host: "prompt", User: "user"})
	if err != nil || len(runtime.Credentials) != 1 || runtime.Credentials[0].Prompt != "Password for user@prompt: " || string(runtime.Credentials[0].Secret) != "secret" {
		t.Fatalf("captured runtime=%#v err=%v", runtime, err)
	}
	runtime.Clear()
	helperProvider := endpointCredentialProvider(sshx.Factory{
		Prompt: func(string) ([]byte, error) { return []byte("secret"), nil },
		SFTPFallback: func(context.Context, *sshx.Connection, *sshx.CapabilityError) (*sftp.Client, io.Closer, error) {
			return &sftp.Client{}, io.NopCloser(strings.NewReader("")), nil
		},
	})
	if runtime, err := helperProvider(context.Background(), endpoint.Endpoint{Remote: true, Host: "helper"}); err != nil || !runtime.AllowHelper {
		t.Fatalf("helper runtime=%#v err=%v", runtime, err)
	} else {
		runtime.Clear()
		if runtime.AllowHelper {
			t.Fatal("helper approval was not cleared")
		}
	}
	if runtime, err := provider(context.Background(), endpoint.Endpoint{Remote: true, Host: "detect"}); err == nil || len(runtime.Credentials) != 0 {
		t.Fatalf("detect failure runtime=%#v err=%v", runtime, err)
	}
	if runtime, err := provider(context.Background(), endpoint.Endpoint{Remote: true, Host: "close"}); err == nil || len(runtime.Credentials) != 0 {
		t.Fatalf("close failure runtime=%#v err=%v", runtime, err)
	}
	if closeCalls != 4 {
		t.Fatalf("close calls=%d", closeCalls)
	}
}

func TestBrowserCommandDependencyFailures(t *testing.T) {
	args := []string{"from", "web://", "to", t.TempDir(), "--background"}
	var stderr bytes.Buffer
	if code := Execute(context.Background(), NewRoot(Dependencies{}), args, io.Discard, &stderr); code != ExitControl || !strings.Contains(stderr.String(), "dependencies") {
		t.Fatalf("nil web dependency code=%d stderr=%q", code, stderr.String())
	}
	dependencies := Dependencies{Web: func(context.Context, operation.Plan, io.Writer) error { return errors.New("web failure") }}
	stderr.Reset()
	if code := Execute(context.Background(), NewRoot(dependencies), args, io.Discard, &stderr); code != ExitControl || !strings.Contains(stderr.String(), "web failure") {
		t.Fatalf("web failure code=%d stderr=%q", code, stderr.String())
	}
}

func TestAcquireWebDelivery(t *testing.T) {
	originalDirectory, originalStore, originalCoordinator, originalNow := defaultStateDirectory, openDeliveryStore, newWebCoordinator, webNow
	t.Cleanup(func() {
		defaultStateDirectory, openDeliveryStore, newWebCoordinator, webNow = originalDirectory, originalStore, originalCoordinator, originalNow
	})
	plan := browserPlan(t, true)
	configured := delivery.DefaultPolicy()
	defaultStateDirectory = func() (string, error) { return "", errors.New("directory") }
	if _, err := acquireWebDelivery(context.Background(), plan, configured, nil); err == nil {
		t.Fatal("state-directory error ignored")
	}
	defaultStateDirectory = func() (string, error) { return t.TempDir(), nil }
	openDeliveryStore = func(string) (*delivery.Store, error) { return nil, errors.New("store") }
	if _, err := acquireWebDelivery(context.Background(), plan, configured, nil); err == nil {
		t.Fatal("store error ignored")
	}
	directory := t.TempDir()
	defaultStateDirectory = func() (string, error) { return directory, nil }
	openDeliveryStore = delivery.OpenStore
	webNow = func() time.Time { return appWebTime }
	called := 0
	newWebCoordinator = func(store *delivery.Store, gotDirectory string) webCoordinator {
		if store == nil || gotDirectory != directory {
			t.Fatal("invalid coordinator inputs")
		}
		return coordinatorFunc(func(_ context.Context, request worker.AcquireRequest) (worker.Acquired, error) {
			called++
			if request.Route != delivery.RoutePathToWeb || request.Foreground || request.At != appWebTime {
				t.Fatalf("unexpected acquisition: %#v", request)
			}
			return worker.Acquired{DeliveryID: appWebDeliveryID}, nil
		})
	}
	if _, err := acquireWebDelivery(context.Background(), plan, configured, json.RawMessage(`{}`)); err != nil || called != 1 {
		t.Fatalf("acquire err=%v called=%d", err, called)
	}
	uploadPlan, err := operation.Build(operation.Request{Source: "web://", Destination: t.TempDir(), Options: []operation.Option{{Name: operation.OptionBackground, Value: "true"}}})
	if err != nil {
		t.Fatal(err)
	}
	newWebCoordinator = func(*delivery.Store, string) webCoordinator {
		return coordinatorFunc(func(_ context.Context, request worker.AcquireRequest) (worker.Acquired, error) {
			if request.Route != delivery.RouteWebToPath {
				t.Fatalf("unexpected upload route: %s", request.Route)
			}
			return worker.Acquired{}, errors.New("acquire")
		})
	}
	if _, err := acquireWebDelivery(context.Background(), uploadPlan, configured, nil); err == nil {
		t.Fatal("coordinator error ignored")
	}
	if originalCoordinator(nil, directory) == nil {
		t.Fatal("default coordinator factory returned nil")
	}
}

func TestWebRunnerBranches(t *testing.T) {
	originalDefinition, originalAddress, originalAcquire, originalRelease := newWebDefinition, webAddress, runWebAcquire, releaseWebLease
	t.Cleanup(func() {
		newWebDefinition, webAddress, runWebAcquire, releaseWebLease = originalDefinition, originalAddress, originalAcquire, originalRelease
	})
	plan := browserPlan(t, true)
	if err := webRunner(nil, nil)(context.Background(), plan, io.Discard); err == nil {
		t.Fatal("nil credential provider accepted")
	}
	credentialError := errors.New("credentials")
	if err := webRunner(func(context.Context, operation.AuthMode) (policy.Credentials, error) {
		return policy.Credentials{}, credentialError
	}, nil)(context.Background(), plan, io.Discard); !errors.Is(err, credentialError) {
		t.Fatalf("credential error=%v", err)
	}
	provider := func(context.Context, operation.AuthMode) (policy.Credentials, error) {
		return policy.Credentials{Password: []byte("secret")}, nil
	}
	remotePlan, err := operation.Build(operation.Request{Source: "host:/source", Destination: "web://", Options: []operation.Option{{Name: operation.OptionBackground, Value: "true"}}})
	if err != nil {
		t.Fatal(err)
	}
	if err := webRunner(provider, nil)(context.Background(), remotePlan, io.Discard); err == nil {
		t.Fatal("missing SSH credential provider accepted")
	}
	endpointError := errors.New("endpoint credentials")
	if err := webRunner(provider, func(context.Context, endpoint.Endpoint) (webdelivery.EndpointRuntime, error) {
		return webdelivery.EndpointRuntime{}, endpointError
	})(context.Background(), remotePlan, io.Discard); !errors.Is(err, endpointError) {
		t.Fatalf("endpoint credential error=%v", err)
	}
	endpointSecret := []byte("endpoint secret")
	newWebDefinition = func(_ operation.Plan, _ policy.Credentials, runtime webdelivery.EndpointRuntime, _ io.Reader) (webdelivery.Definition, error) {
		if secret, promptErr := runtime.Prompt("prompt"); promptErr != nil || string(secret) != "endpoint secret" {
			t.Fatalf("endpoint runtime prompt=%q err=%v", secret, promptErr)
		}
		return webdelivery.Definition{}, errors.New("remote definition")
	}
	if err := webRunner(provider, func(context.Context, endpoint.Endpoint) (webdelivery.EndpointRuntime, error) {
		return webdelivery.EndpointRuntime{Credentials: []webdelivery.EndpointCredential{{Prompt: "prompt", Secret: endpointSecret}}}, nil
	})(context.Background(), remotePlan, io.Discard); err == nil || !bytes.Equal(endpointSecret, make([]byte, len(endpointSecret))) {
		t.Fatalf("remote definition error=%v cleared=%v", err, endpointSecret)
	}
	newWebDefinition = func(operation.Plan, policy.Credentials, webdelivery.EndpointRuntime, io.Reader) (webdelivery.Definition, error) {
		return webdelivery.Definition{}, errors.New("definition")
	}
	if err := webRunner(provider, nil)(context.Background(), plan, io.Discard); err == nil {
		t.Fatal("definition error ignored")
	}
	newWebDefinition = originalDefinition
	overflow := plan
	overflow.Options.Limit = operation.Limit{Value: math.MaxUint64}
	if err := webRunner(provider, nil)(context.Background(), overflow, io.Discard); err == nil {
		t.Fatal("policy conversion error ignored")
	}
	runWebAcquire = func(context.Context, operation.Plan, delivery.Policy, json.RawMessage) (worker.Acquired, error) {
		return worker.Acquired{}, errors.New("acquire")
	}
	if err := webRunner(provider, nil)(context.Background(), plan, io.Discard); err == nil {
		t.Fatal("acquire error ignored")
	}
	lease := &worker.Lease{}
	runWebAcquire = func(context.Context, operation.Plan, delivery.Policy, json.RawMessage) (worker.Acquired, error) {
		return worker.Acquired{DeliveryID: appWebDeliveryID, Lease: lease}, nil
	}
	releases := 0
	releaseWebLease = func(*worker.Lease, context.Context) error { releases++; return nil }
	webAddress = func(string, string) (string, error) { return "", errors.New("address") }
	if err := webRunner(provider, nil)(context.Background(), plan, io.Discard); err == nil || releases != 1 {
		t.Fatalf("address error=%v releases=%d", err, releases)
	}
	webAddress = func(string, string) (string, error) { return "http://delivery/", nil }
	if err := webRunner(provider, nil)(context.Background(), plan, failureWriter{}); err == nil || releases != 2 {
		t.Fatalf("output error=%v releases=%d", err, releases)
	}
	var output bytes.Buffer
	if err := webRunner(provider, nil)(context.Background(), plan, &output); err != nil || !strings.Contains(output.String(), string(appWebDeliveryID)) {
		t.Fatalf("background output=%q err=%v", output.String(), err)
	}
	foreground := plan
	foreground.Options.Background = false
	runWebAcquire = func(context.Context, operation.Plan, delivery.Policy, json.RawMessage) (worker.Acquired, error) {
		return worker.Acquired{DeliveryID: appWebDeliveryID}, nil
	}
	if err := webRunner(provider, nil)(context.Background(), foreground, io.Discard); err == nil {
		t.Fatal("missing foreground lease accepted")
	}
	runWebAcquire = func(context.Context, operation.Plan, delivery.Policy, json.RawMessage) (worker.Acquired, error) {
		return worker.Acquired{DeliveryID: appWebDeliveryID, Lease: lease}, nil
	}
	releaseError := errors.New("release")
	releaseWebLease = func(*worker.Lease, context.Context) error { return releaseError }
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if err := webRunner(provider, nil)(canceled, foreground, io.Discard); !errors.Is(err, context.Canceled) || !errors.Is(err, releaseError) {
		t.Fatalf("foreground result=%v", err)
	}
}

func TestDeliveryPolicyAndFlagCollection(t *testing.T) {
	plan := browserPlan(t, true)
	configured, err := deliveryPolicy(plan)
	if err != nil || configured.Auth != delivery.AuthNone || !configured.DeliveryLimit.Unlimited {
		t.Fatalf("default policy=%#v err=%v", configured, err)
	}
	plan.Options.Auth = operation.AuthBasic
	plan.Options.AuthAttempts = 7
	plan.Options.AuthFailAction = operation.AuthFailStop
	plan.Options.Limit = operation.Limit{Value: 3}
	plan.Options.AllowIP = []netip.Prefix{netip.MustParsePrefix("192.0.2.0/24")}
	plan.Options.MaxFileSize = operation.Quantity{Unlimited: true}
	plan.Options.MaxExtractedSize = operation.Quantity{Value: 123}
	plan.Options.UploadRate = operation.Quantity{Value: 10, Rate: true}
	plan.Options.DownloadRate = operation.Quantity{Value: 20, Rate: true}
	plan.Options.NoUI = true
	configured, err = deliveryPolicy(plan)
	if err != nil || configured.AuthAttempts != 7 || configured.DeliveryLimit.Value != 3 || configured.AllowIP[0] != "192.0.2.0/24" || !configured.MaxFileSize.Unlimited || configured.MaxExtractedSize.Value != 123 || configured.UploadRate.Value != 10 || configured.DownloadRate.Value != 20 || !configured.NoUI {
		t.Fatalf("mapped policy=%#v err=%v", configured, err)
	}
	plan.Options.Limit.Value = math.MaxUint64
	if _, err := deliveryPolicy(plan); err == nil {
		t.Fatal("overflowing limit accepted")
	}

	values := &transferFlagValues{}
	_ = values.archive.Set("true")
	_ = values.extract.Set("false")
	_ = values.listen.Set("127.0.0.1:9000")
	_ = values.background.Set("true")
	_ = values.auth.Set("basic")
	_ = values.authAttempts.Set("7")
	_ = values.authFailAction.Set("stop")
	_ = values.limit.Set("3")
	_ = values.noUI.Set("true")
	_ = values.allowIP.Set("192.0.2.0/24")
	_ = values.maxFileSize.Set("1MiB")
	_ = values.maxExtractedSize.Set("2MiB")
	_ = values.uploadRate.Set("3MiB/s")
	_ = values.downloadRate.Set("4MiB/s")
	_ = values.selection.For(operation.OptionExclude).Set("*.tmp")
	if options := values.options(); len(options) != 15 {
		t.Fatalf("unexpected option count: %d (%v)", len(options), options)
	}
	empty := (&transferFlagValues{}).options()
	if len(empty) != 0 {
		t.Fatalf("empty flags produced %v", empty)
	}
	secret := policy.Credentials{BasicUsername: "u", BasicPassword: []byte("a"), Password: []byte("b")}
	clearCredentials(&secret)
	if secret.BasicUsername != "" || len(secret.BasicPassword) != 0 || len(secret.Password) != 0 {
		t.Fatal("credentials were not cleared")
	}
}
