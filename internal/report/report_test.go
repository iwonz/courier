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
	reporter.Finish()
	if got := output.String(); !strings.Contains(got, "stage=transfer bytes=512 total=1024 speed=256B/s elapsed=1.5s") {
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
}

func TestInteractiveReporter(t *testing.T) {
	var output bytes.Buffer
	reporter := New(&output, true)
	reporter.Handle(progress.Event{Stage: progress.StagePreflight})
	reporter.Handle(progress.Event{Stage: progress.StageTransfer, Current: 4, Total: 8, Elapsed: time.Second})
	reporter.Handle(progress.Event{Stage: progress.StageComplete, Current: 8, Total: 8, Elapsed: 2 * time.Second})
	if prefix := reporter.stagePrefix(decor.Statistics{}); prefix != "complete " {
		t.Fatalf("prefix=%q", prefix)
	}
	reporter.Finish()

	empty := New(&output, true)
	empty.Finish()
}
