package progress

import (
	"testing"
	"time"
)

func TestTracker(t *testing.T) {
	current := time.Unix(100, 0)
	var events []Event
	tracker := New(20, func() time.Time { return current }, func(event Event) { events = append(events, event) })
	tracker.Stage(StageTransfer)
	current = current.Add(2 * time.Second)
	tracker.Add(10)
	tracker.Stage(StageComplete)
	snapshot := tracker.Snapshot()
	if len(events) != 4 || snapshot.Stage != StageComplete || snapshot.Current != 10 || snapshot.Total != 20 || snapshot.Elapsed != 2*time.Second || snapshot.Speed != 5 {
		t.Fatalf("events=%v snapshot=%+v", events, snapshot)
	}
}

func TestTrackerDefaultsAndZeroElapsed(t *testing.T) {
	tracker := New(0, nil, nil)
	tracker.Add(0)
	if event := tracker.Snapshot(); event.Speed < 0 || event.Total != 0 {
		t.Fatalf("event=%+v", event)
	}
}
