package progress

import (
	"math"
	"testing"
	"time"
)

func TestTracker(t *testing.T) {
	current := time.Unix(100, 0)
	var events []Event
	tracker := New(20, func() time.Time { return current }, func(event Event) { events = append(events, event) })
	if err := tracker.Stage(StageTransfer); err != nil {
		t.Fatal(err)
	}
	current = current.Add(2 * time.Second)
	if err := tracker.Add(10); err != nil {
		t.Fatal(err)
	}
	if err := tracker.Stage(StageComplete); err != nil {
		t.Fatal(err)
	}
	snapshot := tracker.Snapshot()
	if len(events) != 4 || snapshot.Stage != StageComplete || snapshot.Read != 10 || snapshot.Sent != 10 || snapshot.Confirmed != 10 || snapshot.Current != 10 || snapshot.Total != 20 || snapshot.Elapsed != 2*time.Second || snapshot.Speed != 5 {
		t.Fatalf("events=%v snapshot=%+v", events, snapshot)
	}
}

func TestTrackerDefaultsAndZeroElapsed(t *testing.T) {
	tracker := New(0, nil, nil)
	if err := tracker.Add(0); err != nil {
		t.Fatal(err)
	}
	if event := tracker.Snapshot(); event.Speed < 0 || event.Total != 0 {
		t.Fatalf("event=%+v", event)
	}
}

func TestTrackerIndependentCountersAndValidation(t *testing.T) {
	tracker := New(4, func() time.Time { return time.Unix(1, 0) }, nil)
	if err := tracker.AddRead(4); err != nil {
		t.Fatal(err)
	}
	if err := tracker.AddSent(3); err != nil {
		t.Fatal(err)
	}
	if err := tracker.AddConfirmed(2); err != nil {
		t.Fatal(err)
	}
	if got := tracker.Snapshot(); got.Read != 4 || got.Sent != 3 || got.Confirmed != 2 || got.Current != 2 {
		t.Fatalf("snapshot=%+v", got)
	}
	if err := tracker.Stage("unknown"); err == nil {
		t.Fatal("invalid stage accepted")
	}
	if Stage("unknown").Valid() {
		t.Fatal("invalid stage is valid")
	}
	if err := tracker.AddCounters(-1, 0, 0); err == nil {
		t.Fatal("negative counter accepted")
	}
	if err := New(0, nil, nil).AddSent(1); err == nil {
		t.Fatal("sent bytes without reads accepted")
	}
	if err := New(0, nil, nil).AddConfirmed(1); err == nil {
		t.Fatal("confirmed bytes without sends accepted")
	}
	overflow := New(0, nil, nil)
	if err := overflow.AddRead(math.MaxInt64); err != nil {
		t.Fatal(err)
	}
	if err := overflow.AddRead(1); err == nil {
		t.Fatal("counter overflow accepted")
	}
}
