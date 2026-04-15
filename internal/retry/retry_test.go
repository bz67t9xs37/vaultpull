package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/retry"
)

func fastPolicy() retry.Policy {
	return retry.Policy{
		MaxAttempts:  3,
		InitialDelay: time.Millisecond,
		MaxDelay:     5 * time.Millisecond,
		Multiplier:   2.0,
	}
}

func TestDo_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastPolicy(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnRetryableError(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastPolicy(), func() error {
		calls++
		if calls < 3 {
			return &retry.Retryable{Cause: errors.New("transient")}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_StopsOnNonRetryableError(t *testing.T) {
	calls := 0
	permanent := errors.New("permanent failure")
	err := retry.Do(context.Background(), fastPolicy(), func() error {
		calls++
		return permanent
	})
	if !errors.Is(err, permanent) {
		t.Fatalf("expected permanent error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_ExhaustsMaxAttempts(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastPolicy(), func() error {
		calls++
		return &retry.Retryable{Cause: errors.New("always fails")}
	})
	if err == nil {
		t.Fatal("expected error after exhausting attempts")
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_RespectsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	err := retry.Do(ctx, fastPolicy(), func() error {
		calls++
		cancel()
		return &retry.Retryable{Cause: errors.New("transient")}
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestIsRetryable(t *testing.T) {
	if retry.IsRetryable(errors.New("plain")) {
		t.Fatal("plain error should not be retryable")
	}
	if !retry.IsRetryable(&retry.Retryable{Cause: errors.New("x")}) {
		t.Fatal("Retryable error should be retryable")
	}
}
