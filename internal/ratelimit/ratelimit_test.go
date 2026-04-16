package ratelimit_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/ratelimit"
)

func newLimiter(max int, window time.Duration) *ratelimit.Limiter {
	return ratelimit.New(ratelimit.Config{Max: max, Window: window})
}

func TestAllow_WithinLimit(t *testing.T) {
	l := newLimiter(3, time.Minute)
	for i := 0; i < 3; i++ {
		ok, err := l.Allow("key1")
		if !ok || err != nil {
			t.Fatalf("expected allow on attempt %d, got err: %v", i+1, err)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	l := newLimiter(2, time.Minute)
	l.Allow("key1")
	l.Allow("key1")
	ok, err := l.Allow("key1")
	if ok {
		t.Fatal("expected deny after exceeding limit")
	}
	if err == nil {
		t.Fatal("expected error when denied")
	}
}

func TestAllow_SeparateKeys(t *testing.T) {
	l := newLimiter(1, time.Minute)
	l.Allow("keyA")
	ok, err := l.Allow("keyB")
	if !ok || err != nil {
		t.Fatalf("separate key should not be rate limited: %v", err)
	}
}

func TestReset_ClearsState(t *testing.T) {
	l := newLimiter(1, time.Minute)
	l.Allow("key1")
	l.Reset("key1")
	ok, err := l.Allow("key1")
	if !ok || err != nil {
		t.Fatalf("expected allow after reset: %v", err)
	}
}

func TestCount_ReturnsCorrectCount(t *testing.T) {
	l := newLimiter(10, time.Minute)
	l.Allow("key1")
	l.Allow("key1")
	l.Allow("key1")
	if c := l.Count("key1"); c != 3 {
		t.Fatalf("expected count 3, got %d", c)
	}
}

func TestCount_ZeroForUnknownKey(t *testing.T) {
	l := newLimiter(5, time.Minute)
	if c := l.Count("unknown"); c != 0 {
		t.Fatalf("expected 0, got %d", c)
	}
}

func TestAllow_DefaultsAppliedOnZeroConfig(t *testing.T) {
	l := ratelimit.New(ratelimit.Config{})
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
	ok, err := l.Allow("x")
	if !ok || err != nil {
		t.Fatalf("unexpected denial with default config: %v", err)
	}
}
