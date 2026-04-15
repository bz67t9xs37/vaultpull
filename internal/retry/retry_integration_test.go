package retry_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/retry"
)

// TestRetry_IntegratesWithVaultLikeOperation simulates a vault fetch that
// fails transiently before succeeding, verifying that retry + backoff work
// end-to-end within a reasonable wall-clock budget.
func TestRetry_IntegratesWithVaultLikeOperation(t *testing.T) {
	var attempts int32
	result := map[string]string{}

	p := retry.Policy{
		MaxAttempts:  5,
		InitialDelay: 2 * time.Millisecond,
		MaxDelay:     20 * time.Millisecond,
		Multiplier:   2.0,
	}

	fetchSecrets := func() error {
		n := atomic.AddInt32(&attempts, 1)
		if n < 4 {
			return &retry.Retryable{Cause: errors.New("vault: 503 service unavailable")}
		}
		result["DB_PASSWORD"] = "s3cr3t"
		return nil
	}

	start := time.Now()
	err := retry.Do(context.Background(), p, fetchSecrets)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
	if attempts != 4 {
		t.Fatalf("expected 4 attempts, got %d", attempts)
	}
	if result["DB_PASSWORD"] != "s3cr3t" {
		t.Fatalf("expected secret to be populated")
	}
	// total delay: 2ms + 4ms + 8ms = 14ms; allow generous headroom
	if elapsed > 500*time.Millisecond {
		t.Fatalf("retry took too long: %v", elapsed)
	}
}

// TestRetry_PermanentErrorNotRetried ensures non-retryable vault errors
// (e.g. 403 permission denied) are surfaced immediately without retrying.
func TestRetry_PermanentErrorNotRetried(t *testing.T) {
	var calls int32
	permDenied := errors.New("vault: 403 permission denied")

	err := retry.Do(context.Background(), retry.DefaultPolicy(), func() error {
		atomic.AddInt32(&calls, 1)
		return permDenied
	})

	if !errors.Is(err, permDenied) {
		t.Fatalf("expected permission denied error, got: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected exactly 1 call for permanent error, got %d", calls)
	}
}
