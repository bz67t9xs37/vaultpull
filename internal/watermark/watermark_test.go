package watermark

import (
	"testing"
	"time"
)

var fixedClock = func() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func TestSet_And_Get(t *testing.T) {
	w := newWithClock(fixedClock)
	w.Set("secret/app", "abc123")

	e, ok := w.Get("secret/app")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Hash != "abc123" {
		t.Errorf("got hash %q, want %q", e.Hash, "abc123")
	}
	if e.Path != "secret/app" {
		t.Errorf("got path %q, want %q", e.Path, "secret/app")
	}
	if !e.RecordedAt.Equal(fixedClock()) {
		t.Errorf("unexpected timestamp: %v", e.RecordedAt)
	}
}

func TestGet_UnknownPath_ReturnsFalse(t *testing.T) {
	w := New()
	_, ok := w.Get("secret/missing")
	if ok {
		t.Fatal("expected ok=false for unknown path")
	}
}

func TestChanged_NeverSeen_ReturnsTrue(t *testing.T) {
	w := New()
	if !w.Changed("secret/new", "hash1") {
		t.Error("expected Changed=true for unseen path")
	}
}

func TestChanged_SameHash_ReturnsFalse(t *testing.T) {
	w := New()
	w.Set("secret/app", "hash1")
	if w.Changed("secret/app", "hash1") {
		t.Error("expected Changed=false when hash is identical")
	}
}

func TestChanged_DifferentHash_ReturnsTrue(t *testing.T) {
	w := New()
	w.Set("secret/app", "hash1")
	if !w.Changed("secret/app", "hash2") {
		t.Error("expected Changed=true when hash differs")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	w := New()
	w.Set("secret/app", "hash1")
	w.Delete("secret/app")
	_, ok := w.Get("secret/app")
	if ok {
		t.Error("expected entry to be removed after Delete")
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	w := New()
	w.Set("secret/a", "h1")
	w.Set("secret/b", "h2")
	w.Set("secret/c", "h3")

	all := w.All()
	if len(all) != 3 {
		t.Errorf("got %d entries, want 3", len(all))
	}
}

func TestSummary_ReflectsCount(t *testing.T) {
	w := New()
	w.Set("secret/a", "h1")
	w.Set("secret/b", "h2")

	s := w.Summary()
	expected := "watermark: 2 path(s) tracked"
	if s != expected {
		t.Errorf("got %q, want %q", s, expected)
	}
}
