package policy

import (
	"context"
	"errors"
	"sync"

	"github.com/iwonz/courier/internal/delivery"
	"golang.org/x/time/rate"
)

const maximumRateBurst = 64 * 1024

type Direction string

const (
	Upload   Direction = "upload"
	Download Direction = "download"
)

type Rates struct {
	mutex    sync.RWMutex
	upload   *rate.Limiter
	download *rate.Limiter
}

func NewRates(upload, download delivery.Limit) (*Rates, error) {
	rates := &Rates{}
	if err := rates.Update(upload, download); err != nil {
		return nil, err
	}
	return rates, nil
}

func (rates *Rates) Update(upload, download delivery.Limit) error {
	if !validRate(upload) || !validRate(download) {
		return errors.New("invalid directional rate")
	}
	rates.mutex.Lock()
	defer rates.mutex.Unlock()
	rates.upload = configureLimiter(rates.upload, upload)
	rates.download = configureLimiter(rates.download, download)
	return nil
}

func validRate(limit delivery.Limit) bool {
	return (limit.Unlimited && limit.Value == 0) || (!limit.Unlimited && limit.Value > 0)
}

func configureLimiter(existing *rate.Limiter, configured delivery.Limit) *rate.Limiter {
	if configured.Unlimited {
		return nil
	}
	burst := configured.Value
	if burst > maximumRateBurst {
		burst = maximumRateBurst
	}
	if existing == nil {
		return rate.NewLimiter(rate.Limit(configured.Value), int(burst))
	}
	existing.SetLimit(rate.Limit(configured.Value))
	existing.SetBurst(int(burst))
	return existing
}

func (rates *Rates) Wait(ctx context.Context, direction Direction, bytes int) error {
	if ctx == nil {
		return errors.New("rate-limit context is required")
	}
	if bytes < 0 {
		return errors.New("rate-limit byte count must not be negative")
	}
	limiter, err := rates.limiter(direction)
	if err != nil || limiter == nil || bytes == 0 {
		return err
	}
	for bytes > 0 {
		chunk := bytes
		if chunk > limiter.Burst() {
			chunk = limiter.Burst()
		}
		if err := limiter.WaitN(ctx, chunk); err != nil {
			return err
		}
		bytes -= chunk
	}
	return nil
}

func (rates *Rates) limiter(direction Direction) (*rate.Limiter, error) {
	rates.mutex.RLock()
	defer rates.mutex.RUnlock()
	switch direction {
	case Upload:
		return rates.upload, nil
	case Download:
		return rates.download, nil
	default:
		return nil, errors.New("invalid transfer direction")
	}
}
