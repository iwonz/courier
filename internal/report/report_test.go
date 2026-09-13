package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/progress"
	"github.com/vbauerster/mpb/v8/decor"
)

func TestLineReporterAndSummaries(t *testing.T) {
	var output bytes.Buffer
	reporter := New(&output, false)
	reporter.Handle(progress.Event{Stage: progress.StageTransfer, Current: 512, Total: 1024, Speed: 256, Elapsed: 1500 * time.Millisecond})
	reporter.Handle(progress.Event{Stage: progress.StageTransfer, Read: 9, Sent: 8, Confirmed: 7, Total: 10})
	reporter.Finish()
	if got := output.String(); !strings.Contains(got, "stage=transfer read=512 sent=512 confirmed=512 total=1024 speed=256B/s elapsed=1.5s") || !strings.Contains(got, "read=9 sent=8 confirmed=7") {
		t.Fatalf("progress=%q", got)
	}

	output.Reset()
	Success(&output, "source", "destination", 42, 1250*time.Millisecond)
	if got := output.String(); !strings.Contains(got, "source: source") || !strings.Contains(got, "destination: destination") || !strings.Contains(got, "transferred: 42 bytes") || !strings.Contains(got, "elapsed: 1.25s") || !strings.Contains(got, "result: success") {
		t.Fatalf("success=%q", got)
	}

	output.Reset()
	Failure(&output, "commit", errors.New("failed"), 21)
	if got := output.String(); !strings.Contains(got, "stage: commit") || !strings.Contains(got, "reason: failed") || !strings.Contains(got, "confirmed: 21 bytes") || !strings.Contains(got, "result: failed") {
		t.Fatalf("failure=%q", got)
	}

	output.Reset()
	FailureCounters(&output, "connection", errors.New("https://user:pass@example.test/a?token=x"), 0, 0, 5)
	if got := output.String(); strings.Contains(got, "pass") || strings.Contains(got, "token") || !strings.Contains(got, "read: 5 bytes") || !strings.Contains(got, "sent: 5 bytes") {
		t.Fatalf("redacted failure=%q", got)
	}
}

func TestInteractiveReporter(t *testing.T) {
	var output bytes.Buffer
	reporter := New(&output, true)
	reporter.Handle(progress.Event{Stage: progress.StagePreflight})
	reporter.Handle(progress.Event{Stage: progress.StageTransfer, Current: 4, Total: 8, Elapsed: time.Second})
	reporter.Handle(progress.Event{Stage: progress.StageComplete, Current: 8, Total: 8, Elapsed: 2 * time.Second})
	if prefix := reporter.stagePrefix(decor.Statistics{}); prefix != "complete r=8 s=8 c=8 " {
		t.Fatalf("prefix=%q", prefix)
	}
	reporter.Finish()

	empty := New(&output, true)
	empty.Finish()
}
