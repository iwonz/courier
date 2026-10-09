package app

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/admin"
	"github.com/iwonz/courier/internal/control"
	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/endpoint"
	"github.com/iwonz/courier/internal/operation"
	"github.com/iwonz/courier/internal/terminalui"
	"github.com/iwonz/courier/internal/update"
	"github.com/spf13/cobra"
)

func forceTerminalMode(t *testing.T, width int) {
	t.Helper()
	original := terminalMode
	terminalMode = func(io.Writer) terminalui.Mode { return terminalui.Mode{Interactive: true, Color: false, Width: width} }
	t.Cleanup(func() { terminalMode = original })
}

func TestStyledServers(t *testing.T) {
	var output bytes.Buffer
	for _, width := range []int{280, 70} {
		got := renderServersStyled(&output, terminalui.Mode{Interactive: true, Width: width}, commandServerViews())
		for _, wanted := range []string{string(commandServerID), string(commandDeliveryID), "127.0.0.1:8080", "confirmed=1", "auth=basic"} {
			if !strings.Contains(got, wanted) {
				t.Fatalf("width=%d missing %q in %q", width, wanted, got)
			}
		}
		if width >= 220 && !strings.Contains(got, "DATA SERVERS") {
			t.Fatalf("wide table=%q", got)
		}
		if width >= 220 {
			for _, heading := range []string{"Started/Updated", "Source/Destination", "Read/Sent/Confirmed", "Created/Updated"} {
				if !strings.Contains(got, heading) {
					t.Fatalf("wide table missing %q: %q", heading, got)
				}
			}
		}
		if width < 220 && (!strings.Contains(got, "SERVER ") || !strings.Contains(got, "DELIVERY ")) {
			t.Fatalf("narrow tables=%q", got)
		}
		if strings.Contains(got, "private.sock") || strings.Contains(got, "web-v1/test") {
			t.Fatalf("sensitive server details leaked: %q", got)
		}
	}
	if got := renderServersStyled(&output, terminalui.Mode{Interactive: true, Width: 80}, nil); !strings.Contains(got, "No Courier data servers found") {
		t.Fatalf("empty=%q", got)
	}
	withoutDeliveries := commandServerViews()
	withoutDeliveries[0].Live = false
	withoutDeliveries[0].Deliveries = nil
	if got := renderServersStyled(&output, terminalui.Mode{Interactive: true, Width: 280}, withoutDeliveries); !strings.Contains(got, "unreachable") {
		t.Fatalf("unreachable=%q", got)
	}

	forceTerminalMode(t, 100)
	dependencies := Dependencies{
		ListServers: func(context.Context) ([]control.ServerView, error) { return commandServerViews(), nil },
		StopServers: func(_ context.Context, request control.StopRequest) (control.StopResult, error) {
			if request.All {
				return control.StopResult{StoppedServers: 2}, nil
			}
			return control.StopResult{ID: request.ID, Kind: delivery.TargetDelivery}, nil
		},
	}
	for _, test := range []struct {
		args []string
		want string
	}{{[]string{"servers"}, "SERVER " + string(commandServerID)}, {[]string{"servers", "stop", "--all"}, "DATA SERVERS STOPPED"}, {[]string{"servers", "stop", string(commandDeliveryID)}, "DATA SERVER STOPPED"}} {
		output.Reset()
		if code := Execute(context.Background(), NewRoot(dependencies), test.args, &output, io.Discard); code != ExitOK || !strings.Contains(output.String(), test.want) {
			t.Fatalf("args=%v code=%d output=%q", test.args, code, output.String())
		}
	}

	command := &cobra.Command{}
	command.SetOut(&output)
	output.Reset()
	if err := renderStopResult(command, control.StopRequest{}, control.StopResult{ID: commandDeliveryID, Kind: delivery.TargetDelivery, AlreadyStopped: true}); err != nil || !strings.Contains(output.String(), "ALREADY STOPPED") {
		t.Fatalf("already stopped=%q err=%v", output.String(), err)
	}
	if err := renderStopResult(command, control.StopRequest{}, control.StopResult{}); err != nil {
		t.Fatal(err)
	}
	command.SetOut(failingWriter{err: errors.New("write")})
	if err := renderStopResult(command, control.StopRequest{All: true}, control.StopResult{StoppedServers: 1}); err == nil {
		t.Fatal("expected rich stop write error")
	}
}

func TestStyledUIAndVersionCommands(t *testing.T) {
	forceTerminalMode(t, 90)
	dependencies := Dependencies{
		Build: BuildIdentity{Version: "v1.2.3", Commit: "abc", Date: "today"},
		StartUI: func(_ context.Context, request admin.StartRequest) (admin.StartResult, error) {
			state := admin.State{ID: commandServerID, Bind: request.Bind, ProcessID: 42}
			return admin.StartResult{State: state}, request.Ready(state)
		},
		StopUI: func(context.Context) (admin.StopResult, error) { return admin.StopResult{AlreadyStopped: true}, nil },
	}
	for _, test := range []struct {
		args []string
		want []string
	}{
		{[]string{"version"}, []string{"COURIER", "v1.2.3", "/releases/tag/v1.2.3", "abc", "today"}},
		{[]string{"ui", "start"}, []string{"ADMINISTRATION UI READY", "foreground", "Press Ctrl+C"}},
		{[]string{"ui", "start", "--background"}, []string{"ADMINISTRATION UI READY", "background", "courier ui stop"}},
		{[]string{"ui", "stop"}, []string{"ADMINISTRATION UI ALREADY STOPPED"}},
	} {
		var output bytes.Buffer
		if code := Execute(context.Background(), NewRoot(dependencies), test.args, &output, io.Discard); code != ExitOK {
			t.Fatalf("args=%v code=%d output=%q", test.args, code, output.String())
		}
		for _, wanted := range test.want {
			if !strings.Contains(output.String(), wanted) {
				t.Fatalf("args=%v missing=%q output=%q", test.args, wanted, output.String())
			}
		}
	}

	dependencies.StopUI = func(context.Context) (admin.StopResult, error) { return admin.StopResult{}, nil }
	var output bytes.Buffer
	if code := Execute(context.Background(), NewRoot(dependencies), []string{"ui", "stop"}, &output, io.Discard); code != ExitOK || !strings.Contains(output.String(), "ADMINISTRATION UI STOPPED") {
		t.Fatalf("stop code=%d output=%q", code, output.String())
	}

	development := newVersionCommand(BuildIdentity{Version: "dev", Commit: "none", Date: "unknown"})
	development.SetOut(&output)
	output.Reset()
	if err := development.Execute(); err != nil || !strings.Contains(output.String(), "github.com/iwonz/courier/releases") || strings.Contains(output.String(), "/tag/") {
		t.Fatalf("development=%q err=%v", output.String(), err)
	}
}

func TestUpdateProgressAndResults(t *testing.T) {
	mode := terminalui.Mode{Interactive: true, Color: false, Width: 90}
	var output bytes.Buffer
	progress := newUpdateProgress(&output, mode)
	stages := []update.Stage{update.StageCheck, update.StageDownloadArchive, update.StageDownloadChecksum, update.StageVerify, update.StageExtract, update.StageInstall, update.StageComplete, update.Stage("unknown")}
	for _, stage := range stages {
		progress.Handle(update.Event{Stage: stage, Current: 5, Total: 10})
		progress.Handle(update.Event{Stage: stage, Current: 10, Total: 10, Complete: true})
	}
	progress.Handle(update.Event{Stage: update.StageDownloadArchive, Current: 7, Total: -1})
	progress.Handle(update.Event{Stage: update.StageDownloadArchive, Current: 11, Total: 10})
	progress.Finish()
	progress.Finish()
	for _, wanted := range []string{"Checking the latest release", "Downloading the release archive", "Downloading checksums", "Verifying the checksum", "Extracting the Courier binary", "Installing the verified binary", "Completing the update", "50%", "7 bytes", "100%"} {
		if !strings.Contains(output.String(), wanted) {
			t.Fatalf("missing %q in %q", wanted, output.String())
		}
	}
	plain := newUpdateProgress(&output, terminalui.Mode{})
	before := output.Len()
	plain.Handle(update.Event{Stage: update.StageCheck})
	plain.Finish()
	if output.Len() != before {
		t.Fatal("plain update progress wrote output")
	}

	output.Reset()
	renderUpdateResult(&output, terminalui.Mode{}, update.Result{Current: true, From: "v1", ReleaseURL: "release"})
	if got := output.String(); got != "courier v1 is current\nrelease: release\n" {
		t.Fatalf("plain current=%q", got)
	}
	output.Reset()
	renderUpdateResult(&output, mode, update.Result{Current: true, From: "v1", ReleaseURL: "release"})
	if got := output.String(); !strings.Contains(got, "COURIER IS UP TO DATE") || !strings.Contains(got, "release") {
		t.Fatalf("rich current=%q", got)
	}
	output.Reset()
	renderUpdateResult(&output, mode, update.Result{From: "v1", To: "v2", Notes: "changes", ReleaseURL: "release"})
	if got := output.String(); !strings.Contains(got, "RELEASE NOTES") || !strings.Contains(got, "COURIER UPDATED") || !strings.Contains(got, "changes") {
		t.Fatalf("rich update=%q", got)
	}
	output.Reset()
	renderUpdateResult(&output, mode, update.Result{From: "v1", To: "v2", ReleaseURL: "release"})
	if got := output.String(); strings.Contains(got, "RELEASE NOTES") || !strings.Contains(got, "COURIER UPDATED") {
		t.Fatalf("rich update without notes=%q", got)
	}

	forceTerminalMode(t, 90)
	dependencies := Dependencies{Update: func(_ context.Context, sink update.Sink) (update.Result, error) {
		sink(update.Event{Stage: update.StageCheck})
		return update.Result{From: "v1", To: "v2"}, nil
	}}
	output.Reset()
	var stderr bytes.Buffer
	if code := Execute(context.Background(), NewRoot(dependencies), []string{"update"}, &output, &stderr); code != ExitOK || !strings.Contains(output.String(), "/releases/tag/v2") || !strings.Contains(stderr.String(), "Checking") {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, output.String(), stderr.String())
	}
}

func TestHostedTerminalLifecycle(t *testing.T) {
	forceTerminalMode(t, 90)
	id := commandDeliveryID
	plan := operation.Plan{
		Route:       operation.RoutePathToWeb,
		Source:      endpoint.Endpoint{Raw: "./source"},
		Destination: endpoint.Endpoint{Raw: "web://"},
	}
	var output bytes.Buffer
	if err := renderHostedReady(&output, plan, "http://127.0.0.1:8080/d/token/", id); err != nil {
		t.Fatal(err)
	}
	for _, wanted := range []string{"DELIVERY READY", "path-to-web", "./source", "web://", string(id), "foreground", "Press Ctrl+C"} {
		if !strings.Contains(output.String(), wanted) {
			t.Fatalf("missing=%q output=%q", wanted, output.String())
		}
	}
	plan.Options.Background = true
	output.Reset()
	if err := renderHostedReady(&output, plan, "http://127.0.0.1:8080/d/token/", id); err != nil || !strings.Contains(output.String(), "background") || !strings.Contains(output.String(), "courier servers stop "+string(id)) {
		t.Fatalf("background=%q err=%v", output.String(), err)
	}
	if err := renderHostedReady(failingWriter{err: errors.New("write")}, plan, "address", id); err == nil {
		t.Fatal("expected hosted ready write failure")
	}

	originalTicker := webStatusTicker
	t.Cleanup(func() { webStatusTicker = originalTicker })
	defaultTicks, stopDefaultTicker := originalTicker()
	if defaultTicks == nil {
		t.Fatal("default ticker returned nil")
	}
	stopDefaultTicker()
	stopped := false
	webStatusTicker = func() (<-chan time.Time, func()) {
		return make(chan time.Time), func() { stopped = true }
	}
	for _, route := range []operation.Route{operation.RoutePathToWeb, operation.RouteWebToPath, operation.RouteWebhookToPath} {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		output.Reset()
		if err := waitForHostedStop(ctx, &output, route, time.Now()); !errors.Is(err, context.Canceled) {
			t.Fatalf("route=%s err=%v", route, err)
		}
	}
	if !stopped {
		t.Fatal("hosted ticker was not stopped")
	}

	tick := make(chan time.Time, 1)
	tick <- time.Now().Add(time.Second)
	webStatusTicker = func() (<-chan time.Time, func()) { return tick, func() {} }
	if err := waitForHostedStop(context.Background(), failingWriter{err: errors.New("write")}, operation.RoutePathToWeb, time.Now()); err == nil {
		t.Fatal("expected initial live status write failure")
	}
	if err := waitForHostedStop(context.Background(), &failAfterWriter{remaining: 1}, operation.RoutePathToWeb, time.Now()); err == nil {
		t.Fatal("expected live status write failure")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := waitForHostedStop(ctx, &failAfterWriter{remaining: 1}, operation.RoutePathToWeb, time.Now()); !errors.Is(err, context.Canceled) || !strings.Contains(err.Error(), "write") {
		t.Fatalf("clear error=%v", err)
	}

	terminalMode = func(io.Writer) terminalui.Mode { return terminalui.Mode{} }
	ctx, cancel = context.WithCancel(context.Background())
	cancel()
	if err := waitForHostedStop(ctx, io.Discard, operation.RoutePathToWeb, time.Now()); !errors.Is(err, context.Canceled) {
		t.Fatalf("plain wait=%v", err)
	}
}

type failAfterWriter struct {
	remaining int
}

func (writer *failAfterWriter) Write(data []byte) (int, error) {
	if writer.remaining <= 0 {
		return 0, errors.New("write")
	}
	writer.remaining--
	return len(data), nil
}
