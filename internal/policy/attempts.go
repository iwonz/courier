package policy

import (
	"crypto/sha256"
	"errors"
	"net/netip"
	"sync"

	"github.com/iwonz/courier/internal/delivery"
)

type FailureDecision struct {
	Attempts   uint64
	Banned     bool
	Stop       bool
	Transition bool
}

type failureState struct {
	attempts uint64
	banned   bool
	stop     bool
}

type AttemptTracker struct {
	mutex     sync.Mutex
	threshold uint64
	action    delivery.AuthFailAction
	states    map[[32]byte]failureState
}

func NewAttemptTracker(threshold uint64, action delivery.AuthFailAction) (*AttemptTracker, error) {
	if threshold == 0 {
		return nil, errors.New("authentication attempt threshold must be positive")
	}
	if action != delivery.AuthFailBan && action != delivery.AuthFailStop {
		return nil, errors.New("invalid authentication failure action")
	}
	return &AttemptTracker{threshold: threshold, action: action, states: make(map[[32]byte]failureState)}, nil
}

func failureFingerprint(id delivery.ID, address netip.Addr) [32]byte {
	return sha256.Sum256([]byte(string(id) + "\x00" + address.Unmap().String()))
}

func (tracker *AttemptTracker) Blocked(id delivery.ID, address netip.Addr) FailureDecision {
	tracker.mutex.Lock()
	defer tracker.mutex.Unlock()
	state := tracker.states[failureFingerprint(id, address)]
	return FailureDecision{Attempts: state.attempts, Banned: state.banned, Stop: state.stop}
}

func (tracker *AttemptTracker) Failure(id delivery.ID, address netip.Addr) FailureDecision {
	tracker.mutex.Lock()
	defer tracker.mutex.Unlock()
	key := failureFingerprint(id, address)
	state := tracker.states[key]
	if state.banned || state.stop {
		return FailureDecision{Attempts: state.attempts, Banned: state.banned, Stop: state.stop}
	}
	state.attempts++
	decision := FailureDecision{Attempts: state.attempts}
	if state.attempts >= tracker.threshold {
		decision.Transition = true
		if tracker.action == delivery.AuthFailBan {
			state.banned = true
			decision.Banned = true
		} else {
			state.stop = true
			decision.Stop = true
		}
	}
	tracker.states[key] = state
	return decision
}

func (tracker *AttemptTracker) Success(id delivery.ID, address netip.Addr) {
	tracker.mutex.Lock()
	defer tracker.mutex.Unlock()
	delete(tracker.states, failureFingerprint(id, address))
}

func (tracker *AttemptTracker) Update(threshold uint64, action delivery.AuthFailAction) error {
	if threshold == 0 {
		return errors.New("authentication attempt threshold must be positive")
	}
	if action != delivery.AuthFailBan && action != delivery.AuthFailStop {
		return errors.New("invalid authentication failure action")
	}
	tracker.mutex.Lock()
	defer tracker.mutex.Unlock()
	tracker.threshold = threshold
	tracker.action = action
	return nil
}
