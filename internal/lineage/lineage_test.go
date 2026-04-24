package lineage

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestRecord_And_History(t *testing.T) {
	tr := newWithClock(fixedClock(time.Unix(1000, 0)))
	tr.Record("DB_PASSWORD", "vault/prod", ".env", "sync")

	h := tr.History("DB_PASSWORD")
	if len(h) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(h))
	}
	if h[0].Operation != "sync" {
		t.Errorf("expected operation sync, got %s", h[0].Operation)
	}
	if h[0].Source != "vault/prod" {
		t.Errorf("unexpected source: %s", h[0].Source)
	}
}

func TestHistory_UnknownKey_ReturnsEmpty(t *testing.T) {
	tr := New()
	h := tr.History("NONEXISTENT")
	if len(h) != 0 {
		t.Errorf("expected empty history, got %d entries", len(h))
	}
}

func TestHistory_OrderedByTime(t *testing.T) {
	now := time.Unix(2000, 0)
	calls := 0
	clock := func() time.Time {
		calls++
		return now.Add(time.Duration(calls) * time.Second)
	}
	tr := newWithClock(clock)
	tr.Record("API_KEY", "vault/prod", ".env", "sync")
	tr.Record("API_KEY", "vault/prod", ".env.staging", "promote")

	h := tr.History("API_KEY")
	if len(h) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(h))
	}
	if !h[0].Timestamp.Before(h[1].Timestamp) {
		t.Error("expected entries ordered by time")
	}
}

func TestAll_ReturnsAllKeys(t *testing.T) {
	tr := New()
	tr.Record("KEY_A", "src", "dst", "sync")
	tr.Record("KEY_B", "src", "dst", "rotate")
	tr.Record("KEY_A", "src", "dst2", "promote")

	all := tr.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(all))
	}
}

func TestSummary_WithEntries(t *testing.T) {
	tr := New()
	tr.Record("SECRET", "vault/prod", ".env", "sync")

	s := tr.Summary("SECRET")
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestSummary_NoEntries(t *testing.T) {
	tr := New()
	s := tr.Summary("MISSING")
	expected := "no lineage recorded for \"MISSING\""
	if s != expected {
		t.Errorf("expected %q, got %q", expected, s)
	}
}
