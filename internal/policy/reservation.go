package policy

import (
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/iwonz/courier/internal/delivery"
)

var (
	ErrTransferLimit = errors.New("delivery transfer limit reached")
	ErrFileTooLarge  = errors.New("incoming file exceeds maximum size")
)

type Reservations struct {
	mutex       sync.Mutex
	limit       delivery.Limit
	maxFileSize delivery.Limit
	active      uint64
}

func NewReservations(limit, maxFileSize delivery.Limit) (*Reservations, error) {
	reservations := &Reservations{}
	if err := reservations.Update(limit, maxFileSize); err != nil {
		return nil, err
	}
	return reservations, nil
}

func validLimit(limit delivery.Limit) bool {
	return (limit.Unlimited && limit.Value == 0) || (!limit.Unlimited && limit.Value > 0)
}

func (reservations *Reservations) Update(limit, maxFileSize delivery.Limit) error {
	if !validLimit(limit) || !validLimit(maxFileSize) {
		return errors.New("invalid reservation limit")
	}
	reservations.mutex.Lock()
	reservations.limit = limit
	reservations.maxFileSize = maxFileSize
	reservations.mutex.Unlock()
	return nil
}

func (reservations *Reservations) Reserve(declaredSize int64) (*Reservation, error) {
	if declaredSize < -1 {
		return nil, errors.New("declared file size must be non-negative or unknown")
	}
	reservations.mutex.Lock()
	defer reservations.mutex.Unlock()
	if !reservations.limit.Unlimited && reservations.active >= uint64(reservations.limit.Value) {
		return nil, ErrTransferLimit
	}
	if declaredSize >= 0 && !reservations.maxFileSize.Unlimited && declaredSize > reservations.maxFileSize.Value {
		return nil, ErrFileTooLarge
	}
	reservations.active++
	return &Reservation{owner: reservations, declared: declaredSize}, nil
}

func (reservations *Reservations) Active() uint64 {
	reservations.mutex.Lock()
	defer reservations.mutex.Unlock()
	return reservations.active
}

type Reservation struct {
	owner    *Reservations
	mutex    sync.Mutex
	declared int64
	consumed int64
	released bool
}

func (reservation *Reservation) Consume(bytes int64) error {
	if bytes < 0 {
		return errors.New("consumed byte count must not be negative")
	}
	reservation.mutex.Lock()
	defer reservation.mutex.Unlock()
	if reservation.released {
		return errors.New("transfer reservation was released")
	}
	if bytes > math.MaxInt64-reservation.consumed {
		return ErrFileTooLarge
	}
	next := reservation.consumed + bytes
	if reservation.declared >= 0 && next > reservation.declared {
		return fmt.Errorf("received bytes exceed declared size: %w", ErrFileTooLarge)
	}
	reservation.owner.mutex.Lock()
	tooLarge := !reservation.owner.maxFileSize.Unlimited && next > reservation.owner.maxFileSize.Value
	reservation.owner.mutex.Unlock()
	if tooLarge {
		return ErrFileTooLarge
	}
	reservation.consumed = next
	return nil
}

func (reservation *Reservation) Consumed() int64 {
	reservation.mutex.Lock()
	defer reservation.mutex.Unlock()
	return reservation.consumed
}

func (reservation *Reservation) Release() {
	reservation.mutex.Lock()
	defer reservation.mutex.Unlock()
	if reservation.released {
		return
	}
	reservation.released = true
	reservation.owner.mutex.Lock()
	reservation.owner.active--
	reservation.owner.mutex.Unlock()
}
