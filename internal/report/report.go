// Package report renders progress and stable operation summaries.
package report

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

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
	stage       atomic.Value
	mutex       sync.Mutex
}

// New creates a terminal-aware progress renderer.
func New(output io.Writer, interactive bool) *Reporter {
	reporter := &Reporter{output: output, interactive: interactive}
	reporter.stage.Store(string(progress.StagePreflight))
	if interactive {
		reporter.container = mpb.New(mpb.WithOutput(output), mpb.WithWidth(72))
	}
	return reporter
}

// Handle consumes one structured progress event.
func (r *Reporter) Handle(event progress.Event) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.stage.Store(string(event.Stage))
	if !r.interactive {
		fmt.Fprintf(r.output, "stage=%s bytes=%d total=%d speed=%.0fB/s elapsed=%s\n", event.Stage, event.Current, event.Total, event.Speed, event.Elapsed.Round(time.Millisecond))
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
	r.bar.SetCurrent(event.Current)
}

func (r *Reporter) stagePrefix(decor.Statistics) string {
	return r.stage.Load().(string) + " "
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
	fmt.Fprintf(output, "stage: %s\nreason: %v\nconfirmed: %d bytes\nresult: failed\n", stage, reason, confirmed)
}
