package delivery

import (
	"fmt"
	"math"
	"sync/atomic"
)

type Counters struct {
	read      atomic.Int64
	sent      atomic.Int64
	confirmed atomic.Int64
}

func NewCounters(snapshot CounterSnapshot) (*Counters, error) {
	if err := snapshot.Validate(); err != nil {
		return nil, err
	}
	result := &Counters{}
	result.read.Store(snapshot.Read)
	result.sent.Store(snapshot.Sent)
	result.confirmed.Store(snapshot.Confirmed)
	return result, nil
}

func (c *Counters) AddRead(delta int64) error      { return addCounter(&c.read, delta) }
func (c *Counters) AddSent(delta int64) error      { return addCounter(&c.sent, delta) }
func (c *Counters) AddConfirmed(delta int64) error { return addCounter(&c.confirmed, delta) }

func (c *Counters) Snapshot() CounterSnapshot {
	return CounterSnapshot{Read: c.read.Load(), Sent: c.sent.Load(), Confirmed: c.confirmed.Load()}
}

func addCounter(counter *atomic.Int64, delta int64) error {
	if delta < 0 {
		return fmt.Errorf("%w: counter delta must be non-negative", ErrInvalid)
	}
	for {
		current := counter.Load()
		if current > math.MaxInt64-delta {
			return fmt.Errorf("%w: counter overflow", ErrInvalid)
		}
		if counter.CompareAndSwap(current, current+delta) {
			return nil
		}
	}
}
