package fingerprint

import (
	"testing"
	"time"
)

var fixedClock = func() time.Time {
	return time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
}

func TestCompute_ProducesHash(t *testing.T) {
	f := newWithClock(fixedClock)
	entry := f.Compute("secret/app", map[string]string{"KEY": "value"})
	if entry.Hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if entry.Path != "secret/app" {
		t.Errorf("expected path secret/app, got %s", entry.Path)
	}
	if entry.KeyCount != 1 {
		t.Errorf("expected key count 1, got %d", entry.KeyCount)
	}
}

func TestCompute_DeterministicOutput(t *testing.T) {
	f := newWithClock(fixedClock)
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	a := f.Compute("path", secrets)
	b := f.Compute("path", secrets)
	if a.Hash != b.Hash {
		t.Errorf("expected same hash, got %s vs %s", a.Hash, b.Hash)
	}
}

func TestCompute_KeyOrderIndependent(t *testing.T) {
	f := newWithClock(fixedClock)
	a := f.Compute("path", map[string]string{"X": "1", "Y": "2"})
	b := f.Compute("path", map[string]string{"Y": "2", "X": "1"})
	if a.Hash != b.Hash {
		t.Errorf("key order should not affect hash: %s vs %s", a.Hash, b.Hash)
	}
}

func TestEqual_SameHash(t *testing.T) {
	f := newWithClock(fixedClock)
	secrets := map[string]string{"K": "v"}
	a := f.Compute("p", secrets)
	b := f.Compute("p", secrets)
	if !Equal(a, b) {
		t.Error("expected Equal to return true for identical secrets")
	}
}

func TestEqual_DifferentHash(t *testing.T) {
	f := newWithClock(fixedClock)
	a := f.Compute("p", map[string]string{"K": "v1"})
	b := f.Compute("p", map[string]string{"K": "v2"})
	if Equal(a, b) {
		t.Error("expected Equal to return false for different secrets")
	}
}

func TestChanged_DetectsModification(t *testing.T) {
	f := newWithClock(fixedClock)
	original := map[string]string{"DB_PASS": "old"}
	stored := f.Compute("secret/db", original)

	if f.Changed(stored, "secret/db", original) {
		t.Error("expected Changed=false for same secrets")
	}

	updated := map[string]string{"DB_PASS": "new"}
	if !f.Changed(stored, "secret/db", updated) {
		t.Error("expected Changed=true after value update")
	}
}

func TestCompute_SetsTimestamp(t *testing.T) {
	f := newWithClock(fixedClock)
	entry := f.Compute("path", map[string]string{})
	if !entry.ComputedAt.Equal(fixedClock()) {
		t.Errorf("expected timestamp %v, got %v", fixedClock(), entry.ComputedAt)
	}
}
