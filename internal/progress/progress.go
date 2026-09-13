// Package progress produces renderer-neutral transfer progress events.
package progress

import "time"

// Stage identifies the currently executing operation phase.
type Stage string

const (
	StagePreflight Stage = "preflight"
	StageArchive   Stage = "archive"
	StageExtract   Stage = "extract"
	StageTransfer  Stage = "transfer"
	StageCommit    Stage = "commit"
	StageVerify    Stage = "verify"
	StageCleanup   Stage = "cleanup"
	StageComplete  Stage = "complete"
)

// Event is a point-in-time progress snapshot.
type Event struct {
	Stage   Stage
	Current int64
	Total   int64
	Speed   float64
	Elapsed time.Duration
}

// Sink consumes structured events.
type Sink func(Event)

// Tracker maintains monotonic byte accounting.
type Tracker struct {
	started time.Time
	now     func() time.Time
	sink    Sink
	event   Event
}

// New creates a tracker and emits its initial stage.
func New(total int64, now func() time.Time, sink Sink) *Tracker {
	if now == nil {
		now = time.Now
	}
	tracker := &Tracker{started: now(), now: now, sink: sink, event: Event{Stage: StagePreflight, Total: total}}
	tracker.emit()
	return tracker
}

// Stage changes the phase without changing byte accounting.
func (t *Tracker) Stage(stage Stage) { t.event.Stage = stage; t.emit() }

// Add records confirmed bytes.
func (t *Tracker) Add(bytes int64) { t.event.Current += bytes; t.emit() }

// Snapshot returns the latest computed event.
func (t *Tracker) Snapshot() Event {
	t.updateTime()
	return t.event
}

func (t *Tracker) emit() {
	t.updateTime()
	if t.sink != nil {
		t.sink(t.event)
	}
}

func (t *Tracker) updateTime() {
	t.event.Elapsed = t.now().Sub(t.started)
	if t.event.Elapsed > 0 {
		t.event.Speed = float64(t.event.Current) / t.event.Elapsed.Seconds()
	}
}
