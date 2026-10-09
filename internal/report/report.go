// Package report renders progress and stable operation summaries.
package report

import (
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/iwonz/courier/internal/diagnostic"
	"github.com/iwonz/courier/internal/progress"
	"github.com/iwonz/courier/internal/terminalui"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

// Reporter adapts confirmed-byte events to mpb or deterministic lines.
type Reporter struct {
	output      io.Writer
	interactive bool
	container   *mpb.Progress
	bar         *mpb.Bar
	terminal    *terminalui.Renderer
	event       atomic.Value
	mutex       sync.Mutex
}

// New creates a terminal-aware progress renderer.
func New(output io.Writer, interactive bool) *Reporter {
	return NewWithMode(output, terminalui.Mode{Interactive: interactive, Color: interactive, Width: 72})
}

// NewWithMode creates a progress renderer with explicit terminal capabilities.
func NewWithMode(output io.Writer, mode terminalui.Mode) *Reporter {
	reporter := &Reporter{output: output, interactive: mode.Interactive, terminal: terminalui.New(output, mode)}
	reporter.event.Store(progress.Event{Stage: progress.StagePreflight})
	if mode.Interactive {
		width := mode.Width
		if width > 72 {
			width = 72
		}
		reporter.container = mpb.New(mpb.WithOutput(output), mpb.WithWidth(width))
	}
	return reporter
}

// Handle consumes one structured progress event.
func (r *Reporter) Handle(event progress.Event) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	event = normalize(event)
	r.event.Store(event)
	if !r.interactive {
		fmt.Fprintf(r.output, "stage=%s read=%d sent=%d confirmed=%d total=%d speed=%.0fB/s elapsed=%s\n", event.Stage, event.Read, event.Sent, event.Confirmed, event.Total, event.Speed, event.Elapsed.Round(time.Millisecond))
		return
	}
	if r.bar == nil {
		total := event.Total
		if total <= 0 {
			total = 1
		}
		barStyle := mpb.BarStyle().Lbound("[").Filler("■").Tip("■").Padding("·").Rbound("]").
			FillerMeta(r.successText).
			TipMeta(r.warningText)
		r.bar, _ = r.container.Add(total, barStyle.Build(),
			mpb.PrependDecorators(decor.Any(r.stagePrefix)),
			mpb.AppendDecorators(
				decor.CountersKibiByte(" % .1f / % .1f"),
				decor.AverageSpeed(decor.SizeB1024(0), " % .1f/s"),
				decor.Elapsed(decor.ET_STYLE_GO, decor.WCSyncSpaceR),
			),
		)
	}
	if event.Total > 0 {
		r.bar.SetTotal(event.Total, false)
	}
	r.bar.SetCurrent(event.Confirmed)
}

func (r *Reporter) successText(value string) string {
	return r.terminal.Text(terminalui.ToneSuccess, value)
}

func (r *Reporter) warningText(value string) string {
	return r.terminal.Text(terminalui.ToneWarning, value)
}

func (r *Reporter) stagePrefix(decor.Statistics) string {
	event := r.event.Load().(progress.Event)
	return r.terminal.Text(terminalui.ToneInfo, string(event.Stage)) + fmt.Sprintf(" r=%d s=%d c=%d ", event.Read, event.Sent, event.Confirmed)
}

// Finish flushes terminal rendering.
func (r *Reporter) Finish() {
	r.mutex.Lock()
	if r.bar != nil {
		r.bar.SetTotal(-1, true)
	}
	container := r.container
	r.mutex.Unlock()
	if container != nil {
		container.Wait()
	}
}

// Success writes the stable final transfer summary.
func Success(output io.Writer, source, destination string, bytes int64, elapsed time.Duration) {
	SuccessWithMode(output, terminalui.Mode{}, source, destination, bytes, elapsed)
}

// SuccessWithMode writes a rich or stable final transfer summary.
func SuccessWithMode(output io.Writer, mode terminalui.Mode, source, destination string, bytes int64, elapsed time.Duration) {
	if !mode.Interactive {
		fmt.Fprintf(output, "source: %s\ndestination: %s\ntransferred: %d bytes\nelapsed: %s\nresult: success\n", source, destination, bytes, elapsed.Round(time.Millisecond))
		return
	}
	panel := terminalui.New(output, mode).Panel("Transfer complete", terminalui.ToneSuccess, []terminalui.Field{
		{Label: "Source", Value: source},
		{Label: "Destination", Value: destination},
		{Label: "Transferred", Value: fmt.Sprintf("%d bytes", bytes)},
		{Label: "Elapsed", Value: elapsed.Round(time.Millisecond).String()},
	}, "✓ Confirmed by the destination")
	fmt.Fprint(output, panel)
}

// Failure writes a stable stage-aware error summary.
func Failure(output io.Writer, stage string, reason error, confirmed int64) {
	FailureCounters(output, stage, reason, confirmed, confirmed, confirmed)
}

// FailureCounters writes a sanitized failure with each accounting boundary.
func FailureCounters(output io.Writer, stage string, reason error, read, sent, confirmed int64) {
	FailureCountersWithMode(output, terminalui.Mode{}, stage, reason, read, sent, confirmed)
}

// FailureCountersWithMode writes a styled or stable sanitized failure.
func FailureCountersWithMode(output io.Writer, mode terminalui.Mode, stage string, reason error, read, sent, confirmed int64) {
	if sent < confirmed {
		sent = confirmed
	}
	if read < sent {
		read = sent
	}
	if !mode.Interactive {
		fmt.Fprintf(output, "stage: %s\nreason: %s\nread: %d bytes\nsent: %d bytes\nconfirmed: %d bytes\nresult: failed\n", stage, diagnostic.Redact(reason.Error()), read, sent, confirmed)
		return
	}
	panel := terminalui.New(output, mode).Panel("Operation failed", terminalui.ToneFailure, []terminalui.Field{
		{Label: "Stage", Value: stage},
		{Label: "Reason", Value: diagnostic.Redact(reason.Error())},
		{Label: "Read", Value: fmt.Sprintf("%d bytes", read)},
		{Label: "Sent", Value: fmt.Sprintf("%d bytes", sent)},
		{Label: "Confirmed", Value: fmt.Sprintf("%d bytes", confirmed)},
	}, "× No uncertain outcome was reported as success")
	fmt.Fprint(output, panel)
}

// InterruptedWithMode renders a user interruption while retaining exit code 130.
func InterruptedWithMode(output io.Writer, mode terminalui.Mode, stage string, read, sent, confirmed int64) {
	if !mode.Interactive {
		FailureCounters(output, stage, context.Canceled, read, sent, confirmed)
		return
	}
	panel := terminalui.New(output, mode).Panel("Stopped by user", terminalui.ToneWarning, []terminalui.Field{
		{Label: "Stage", Value: stage},
		{Label: "Read", Value: fmt.Sprintf("%d bytes", read)},
		{Label: "Sent", Value: fmt.Sprintf("%d bytes", sent)},
		{Label: "Confirmed", Value: fmt.Sprintf("%d bytes", confirmed)},
	}, "! Exit code 130")
	fmt.Fprint(output, panel)
}

func normalize(event progress.Event) progress.Event {
	if event.Confirmed == 0 && event.Current != 0 {
		event.Confirmed = event.Current
	}
	if event.Sent < event.Confirmed {
		event.Sent = event.Confirmed
	}
	if event.Read < event.Sent {
		event.Read = event.Sent
	}
	event.Current = event.Confirmed
	return event
}
