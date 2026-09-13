package delivery

import (
	"fmt"
	"math"
	"sync"
)

type Counters struct {
	mutex    sync.Mutex
	snapshot CounterSnapshot
}

func NewCounters(snapshot CounterSnapshot) (*Counters, error) {
	if err := snapshot.Validate(); err != nil {
		return nil, err
	}
	return &Counters{snapshot: snapshot}, nil
}

func (c *Counters) AddRead(delta int64) error      { return c.add(delta, 0, 0) }
func (c *Counters) AddSent(delta int64) error      { return c.add(0, delta, 0) }
func (c *Counters) AddConfirmed(delta int64) error { return c.add(0, 0, delta) }

func (c *Counters) Snapshot() CounterSnapshot {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.snapshot
}

func (c *Counters) add(read, sent, confirmed int64) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if read < 0 || sent < 0 || confirmed < 0 {
		return fmt.Errorf("%w: counter delta must be non-negative", ErrInvalid)
	}
	if c.snapshot.Read > math.MaxInt64-read || c.snapshot.Sent > math.MaxInt64-sent || c.snapshot.Confirmed > math.MaxInt64-confirmed {
		return fmt.Errorf("%w: counter overflow", ErrInvalid)
	}
	next := CounterSnapshot{Read: c.snapshot.Read + read, Sent: c.snapshot.Sent + sent, Confirmed: c.snapshot.Confirmed + confirmed}
	if err := next.Validate(); err != nil {
		return err
	}
	c.snapshot = next
	return nil
}
