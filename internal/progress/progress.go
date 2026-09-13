// Package progress produces renderer-neutral transfer progress events.
package progress

import (
	"fmt"
	"math"
	"sync"
	"time"
)

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

// Valid reports whether a stage belongs to Courier's public operation
// vocabulary.
func (stage Stage) Valid() bool {
	switch stage {
	case StagePreflight, StageArchive, StageExtract, StageTransfer, StageCommit, StageVerify, StageCleanup, StageComplete:
		return true
	default:
		return false
	}
}

// Event is a point-in-time progress snapshot.
type Event struct {
	Stage     Stage
	Read      int64
	Sent      int64
	Confirmed int64
	// Current is the confirmed-byte compatibility projection used by progress
	// bar renderers. Tracker emissions always set it to Confirmed.
	Current int64
	Total   int64
	Speed   float64
	Elapsed time.Duration
}

// Sink consumes structured events.
type Sink func(Event)

// Tracker maintains monotonic byte accounting.
type Tracker struct {
	mutex   sync.Mutex
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
func (t *Tracker) Stage(stage Stage) error {
	t.mutex.Lock()
	if !stage.Valid() {
		t.mutex.Unlock()
		return fmt.Errorf("invalid progress stage %q", stage)
	}
	t.event.Stage = stage
	event, sink := t.emissionLocked()
	t.mutex.Unlock()
	if sink != nil {
		sink(event)
	}
	return nil
}

// Add records a chunk that was read, sent, and confirmed at one boundary.
func (t *Tracker) Add(bytes int64) error { return t.AddCounters(bytes, bytes, bytes) }

// AddRead records source bytes without claiming that they were sent.
func (t *Tracker) AddRead(bytes int64) error { return t.AddCounters(bytes, 0, 0) }

// AddSent records bytes submitted to a destination without claiming target
// confirmation. The caller must already have recorded at least as many reads.
func (t *Tracker) AddSent(bytes int64) error { return t.AddCounters(0, bytes, 0) }

// AddConfirmed records target-confirmed bytes. The caller must already have
// recorded at least as many sent bytes.
func (t *Tracker) AddConfirmed(bytes int64) error { return t.AddCounters(0, 0, bytes) }

// AddCounters atomically advances the three accounting legs.
func (t *Tracker) AddCounters(read, sent, confirmed int64) error {
	t.mutex.Lock()
	if read < 0 || sent < 0 || confirmed < 0 {
		t.mutex.Unlock()
		return fmt.Errorf("progress counters cannot decrease")
	}
	nextRead, okRead := checkedAdd(t.event.Read, read)
	nextSent, okSent := checkedAdd(t.event.Sent, sent)
	nextConfirmed, okConfirmed := checkedAdd(t.event.Confirmed, confirmed)
	if !okRead || !okSent || !okConfirmed {
		t.mutex.Unlock()
		return fmt.Errorf("progress counter overflow")
	}
	if nextConfirmed > nextSent || nextSent > nextRead {
		t.mutex.Unlock()
		return fmt.Errorf("progress counters require confirmed <= sent <= read")
	}
	t.event.Read, t.event.Sent, t.event.Confirmed = nextRead, nextSent, nextConfirmed
	event, sink := t.emissionLocked()
	t.mutex.Unlock()
	if sink != nil {
		sink(event)
	}
	return nil
}

// Snapshot returns the latest computed event.
func (t *Tracker) Snapshot() Event {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.updateTimeLocked()
	return t.event
}

func (t *Tracker) emit() {
	t.mutex.Lock()
	event, sink := t.emissionLocked()
	t.mutex.Unlock()
	if sink != nil {
		sink(event)
	}
}

func (t *Tracker) emissionLocked() (Event, Sink) {
	t.updateTimeLocked()
	return t.event, t.sink
}

func (t *Tracker) updateTimeLocked() {
	t.event.Current = t.event.Confirmed
	t.event.Elapsed = t.now().Sub(t.started)
	if t.event.Elapsed > 0 {
		t.event.Speed = float64(t.event.Confirmed) / t.event.Elapsed.Seconds()
	}
}

func checkedAdd(current, delta int64) (int64, bool) {
	if current > math.MaxInt64-delta {
		return 0, false
	}
	return current + delta, true
}
