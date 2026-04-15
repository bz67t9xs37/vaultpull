package retry

import (
	"context"
	"errors"
	"math"
	"time"
)

// Policy defines the retry behaviour for vault operations.
type Policy struct {
	MaxAttempts int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// DefaultPolicy returns a sensible default retry policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts:  4,
		InitialDelay: 250 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
	}
}

// Retryable wraps an error to signal that the operation should be retried.
type Retryable struct {
	Cause error
}

func (r *Retryable) Error() string { return r.Cause.Error() }
func (r *Retryable) Unwrap() error { return r.Cause }

// IsRetryable reports whether err is a retryable error.
func IsRetryable(err error) bool {
	var r *Retryable
	return errors.As(err, &r)
}

// Do executes fn according to p, retrying on Retryable errors.
// It respects ctx cancellation between attempts.
func Do(ctx context.Context, p Policy, fn func() error) error {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 1
	}
	delay := p.InitialDelay
	var lastErr error
	for attempt := 1; attempt <= p.MaxAttempts; attempt++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if !IsRetryable(lastErr) {
			return lastErr
		}
		if attempt == p.MaxAttempts {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
		next := time.Duration(math.Min(
			float64(delay)*p.Multiplier,
			float64(p.MaxDelay),
		))
		delay = next
	}
	return lastErr
}
