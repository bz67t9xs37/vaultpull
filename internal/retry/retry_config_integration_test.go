package retry_test

import (
	"errors"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/retry"
)

func TestRetry_UsesConfigPolicy(t *testing.T) {
	cfg := &config.RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 1 * time.Millisecond,
		MaxDelay:     5 * time.Millisecond,
		Multiplier:   2.0,
	}

	attempts, initial, maxD, mult := cfg.ToPolicy()
	policy := retry.Policy{
		MaxAttempts:  attempts,
		InitialDelay: initial,
		MaxDelay:     maxD,
		Multiplier:   mult,
	}

	calls := 0
	err := retry.Do(policy, func() error {
		calls++
		if calls < 3 {
			return errors.New("connection refused")
		}
		return nil
	})

	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestRetry_ConfigDefaultsApplied(t *testing.T) {
	cfg := &config.RetryConfig{}
	config.ApplyRetryDefaults(cfg)

	attempts, initial, maxD, mult := cfg.ToPolicy()
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
	if initial != 500*time.Millisecond {
		t.Errorf("unexpected initial delay: %v", initial)
	}
	_ = maxD
	_ = mult
}
