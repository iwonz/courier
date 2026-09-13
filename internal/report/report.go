// Package report renders progress and stable operation summaries.
package report

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/iwonz/courier/internal/diagnostic"
	"github.com/iwonz/courier/internal/progress"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

// Reporter adapts confirmed-byte events to mpb or deterministic lines.
type Reporter struct {
	output      io.Writer
	interactive bool
	container   *mpb.Progress
	bar         *mpb.Bar
	event       atomic.Value
	mutex       sync.Mutex
}

// New creates a terminal-aware progress renderer.
func New(output io.Writer, interactive bool) *Reporter {
	reporter := &Reporter{output: output, interactive: interactive}
	reporter.event.Store(progress.Event{Stage: progress.StagePreflight})
	if interactive {
		reporter.container = mpb.New(mpb.WithOutput(output), mpb.WithWidth(72))
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
		r.bar = r.container.AddBar(total,
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

func (r *Reporter) stagePrefix(decor.Statistics) string {
	event := r.event.Load().(progress.Event)
	return fmt.Sprintf("%s r=%d s=%d c=%d ", event.Stage, event.Read, event.Sent, event.Confirmed)
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
	fmt.Fprintf(output, "source: %s\ndestination: %s\ntransferred: %d bytes\nelapsed: %s\nresult: success\n", source, destination, bytes, elapsed.Round(time.Millisecond))
}

// Failure writes a stable stage-aware error summary.
func Failure(output io.Writer, stage string, reason error, confirmed int64) {
	FailureCounters(output, stage, reason, confirmed, confirmed, confirmed)
}

// FailureCounters writes a sanitized failure with each accounting boundary.
func FailureCounters(output io.Writer, stage string, reason error, read, sent, confirmed int64) {
	if sent < confirmed {
		sent = confirmed
	}
	if read < sent {
		read = sent
	}
	fmt.Fprintf(output, "stage: %s\nreason: %s\nread: %d bytes\nsent: %d bytes\nconfirmed: %d bytes\nresult: failed\n", stage, diagnostic.Redact(reason.Error()), read, sent, confirmed)
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
