package app

import (
	"bytes"
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/control"
	"github.com/iwonz/courier/internal/delivery"
)

const (
	commandServerID   = delivery.ID("00000000-0000-4000-8000-000000000501")
	commandDeliveryID = delivery.ID("00000000-0000-4000-8000-000000000502")
)

func commandServerViews() []control.ServerView {
	at := time.Date(2026, 9, 14, 12, 13, 14, 15, time.FixedZone("test", 3600))
	configured := delivery.DefaultPolicy()
	configured.Auth = delivery.AuthBasic
	configured.AllowIP = []string{"2001:db8::/32", "127.0.0.1/32"}
	configured.DeliveryLimit = delivery.Limit{Value: 2}
	configured.NoUI = true
	return []control.ServerView{{
		Live: true,
		Server: delivery.Server{
			ID: commandServerID, Bind: "127.0.0.1:8080", ControlEndpoint: "private.sock", Compatibility: "web-v1/test",
			ProcessID: 42, State: delivery.StateActive, StartedAt: at, UpdatedAt: at,
		},
		Deliveries: []delivery.Delivery{
			{
				ID: commandDeliveryID, ServerID: commandServerID, Route: delivery.RoutePathToWeb,
				Source: "./data", Destination: "web://", State: delivery.StateActive, Policy: configured,
				Counters: delivery.CounterSnapshot{Read: 3, Sent: 2, Confirmed: 1}, CreatedAt: at, UpdatedAt: at,
			},
			{
				ID: delivery.ID("00000000-0000-4000-8000-000000000503"), ServerID: commandServerID,
				Route: delivery.RouteWebToPath, State: delivery.StateStarting, Policy: delivery.DefaultPolicy(), CreatedAt: at, UpdatedAt: at,
			},
		},
	}}
}

func TestServersCommand(t *testing.T) {
	want := errors.New("control failure")
	dependencies := Dependencies{
		ListServers: func(context.Context) ([]control.ServerView, error) { return commandServerViews(), nil },
		StopServers: func(_ context.Context, request control.StopRequest) (control.StopResult, error) {
			if request.All {
				return control.StopResult{Kind: delivery.TargetServer, StoppedServers: 2}, nil
			}
			return control.StopResult{ID: request.ID, Kind: delivery.TargetDelivery}, nil
		},
	}
	for _, test := range []struct {
		name       string
		args       []string
		code       int
		wantOutput string
		wantError  string
	}{
		{"list", []string{"servers"}, ExitOK, "delivery " + string(commandDeliveryID), ""},
		{"list args", []string{"servers", "extra"}, ExitCLI, "", "unknown command"},
		{"stop delivery", []string{"servers", "stop", string(commandDeliveryID)}, ExitOK, "Stopped delivery", ""},
		{"stop all", []string{"servers", "stop", "--all"}, ExitOK, "Stopped data servers: 2", ""},
		{"stop all conflict", []string{"servers", "stop", string(commandDeliveryID), "--all"}, ExitCLI, "", "cannot be combined"},
		{"stop missing", []string{"servers", "stop"}, ExitCLI, "", "expected one canonical UUID"},
		{"stop extra", []string{"servers", "stop", string(commandDeliveryID), string(commandServerID)}, ExitCLI, "", "expected one canonical UUID"},
		{"stop invalid UUID", []string{"servers", "stop", "not-a-uuid"}, ExitCLI, "", "canonical UUID"},
		{"stop repeated all", []string{"servers", "stop", "--all", "--all"}, ExitCLI, "", "value may only be set once"},
	} {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := Execute(context.Background(), NewRoot(dependencies), test.args, &stdout, &stderr)
			if code != test.code || !strings.Contains(stdout.String(), test.wantOutput) || !strings.Contains(stderr.String(), test.wantError) {
				t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
		})
	}

	withoutList := dependencies
	withoutList.ListServers = nil
	if code := Execute(context.Background(), NewRoot(withoutList), []string{"servers"}, io.Discard, io.Discard); code != ExitControl {
		t.Fatalf("missing list code=%d", code)
	}
	withoutStop := dependencies
	withoutStop.StopServers = nil
	if code := Execute(context.Background(), NewRoot(withoutStop), []string{"servers", "stop", "--all"}, io.Discard, io.Discard); code != ExitControl {
		t.Fatalf("missing stop code=%d", code)
	}
	failingList := dependencies
	failingList.ListServers = func(context.Context) ([]control.ServerView, error) { return nil, want }
	if code := Execute(context.Background(), NewRoot(failingList), []string{"servers"}, io.Discard, io.Discard); code != ExitControl {
		t.Fatalf("failed list code=%d", code)
	}
	failingStop := dependencies
	failingStop.StopServers = func(context.Context, control.StopRequest) (control.StopResult, error) {
		return control.StopResult{Kind: delivery.TargetServer, StoppedServers: 1}, want
	}
	var partial bytes.Buffer
	if code := Execute(context.Background(), NewRoot(failingStop), []string{"servers", "stop", "--all"}, &partial, io.Discard); code != ExitControl || !strings.Contains(partial.String(), "1") {
		t.Fatalf("partial stop code=%d output=%q", code, partial.String())
	}
}

func TestServerCommandRenderingAndOutputFailures(t *testing.T) {
	output := renderServers(commandServerViews())
	for _, expected := range []string{
		"status: live", "source: ./data", "destination: web://", "auth=basic", "limit=2", "max-file-size=10737418240",
		"allow-ip=127.0.0.1/32,2001:db8::/32", "read=3 sent=2 confirmed=1", "source: <unavailable>", "limit=unlimited", "2026-09-14T11:13:14.000000015Z",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("missing %q in %q", expected, output)
		}
	}
	for _, forbidden := range []string{"private.sock", "web-v1/test", "Authorization", "resource token"} {
		if strings.Contains(output, forbidden) {
			t.Fatalf("control output leaked %q: %q", forbidden, output)
		}
	}
	if output := renderServers(nil); output != "No Courier data servers found.\n" {
		t.Fatalf("empty output=%q", output)
	}

	want := errors.New("write failure")
	dependencies := Dependencies{ListServers: func(context.Context) ([]control.ServerView, error) { return nil, nil }}
	if code := Execute(context.Background(), NewRoot(dependencies), []string{"servers"}, failingWriter{err: want}, io.Discard); code != ExitControl {
		t.Fatalf("list write code=%d", code)
	}
	dependencies.StopServers = func(context.Context, control.StopRequest) (control.StopResult, error) {
		return control.StopResult{ID: commandServerID, Kind: delivery.TargetServer}, nil
	}
	if code := Execute(context.Background(), NewRoot(dependencies), []string{"servers", "stop", string(commandServerID)}, failingWriter{err: want}, io.Discard); code != ExitControl {
		t.Fatalf("stop write code=%d", code)
	}

	command := newServersStopCommand(func(context.Context, control.StopRequest) (control.StopResult, error) {
		return control.StopResult{ID: commandDeliveryID, Kind: delivery.TargetDelivery, AlreadyStopped: true}, nil
	})
	command.SetOut(&bytes.Buffer{})
	command.SetArgs([]string{string(commandDeliveryID)})
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if err := renderStopResult(command, control.StopRequest{}, control.StopResult{}); err != nil {
		t.Fatal(err)
	}
}

func TestDefaultServerControlFunctions(t *testing.T) {
	originalDirectory, originalStore := defaultStateDirectory, openDeliveryStore
	t.Cleanup(func() { defaultStateDirectory, openDeliveryStore = originalDirectory, originalStore })
	directory := filepath.Join(t.TempDir(), "state")
	defaultStateDirectory = func() (string, error) { return directory, nil }
	openDeliveryStore = delivery.OpenStore
	views, err := listServers(context.Background())
	if err != nil || len(views) != 0 {
		t.Fatalf("views=%v err=%v", views, err)
	}
	result, err := stopServers(context.Background(), control.StopRequest{All: true})
	if err != nil || result.StoppedServers != 0 {
		t.Fatalf("result=%+v err=%v", result, err)
	}

	want := errors.New("state failure")
	defaultStateDirectory = func() (string, error) { return "", want }
	if _, err := listServers(context.Background()); !errors.Is(err, want) {
		t.Fatalf("directory error=%v", err)
	}
	defaultStateDirectory = func() (string, error) { return directory, nil }
	openDeliveryStore = func(string) (*delivery.Store, error) { return nil, want }
	if _, err := stopServers(context.Background(), control.StopRequest{All: true}); !errors.Is(err, want) {
		t.Fatalf("store error=%v", err)
	}
}
