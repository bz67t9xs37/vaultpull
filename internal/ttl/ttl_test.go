package ttl_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/ttl"
)

func TestTrack_And_Get(t *testing.T) {
	tr := ttl.New(5 * time.Minute)
	tr.Track("API_KEY", "secret123")

	e := tr.Get("API_KEY")
	if e == nil {
		t.Fatal("expected entry, got nil")
	}
	if e.Value != "secret123" {
		t.Errorf("expected value 'secret123', got %q", e.Value)
	}
}

func TestGet_Untracked_ReturnsNil(t *testing.T) {
	tr := ttl.New(5 * time.Minute)
	if tr.Get("MISSING") != nil {
		t.Error("expected nil for untracked key")
	}
}

func TestIsExpired_NotYet(t *testing.T) {
	tr := ttl.New(10 * time.Minute)
	tr.Track("DB_PASS", "hunter2")
	e := tr.Get("DB_PASS")
	if e.IsExpired() {
		t.Error("entry should not be expired yet")
	}
}

func TestIsExpired_Past(t *testing.T) {
	tr := ttl.New(1 * time.Millisecond)
	tr.Track("OLD_KEY", "value")
	time.Sleep(5 * time.Millisecond)
	e := tr.Get("OLD_KEY")
	if !e.IsExpired() {
		t.Error("entry should be expired")
	}
}

func TestIsExpired_ZeroTTL_NeverExpires(t *testing.T) {
	tr := ttl.New(0)
	tr.Track("FOREVER", "value")
	e := tr.Get("FOREVER")
	if e.IsExpired() {
		t.Error("zero-TTL entry should never expire")
	}
}

func TestExpiredKeys_ReturnsOnlyExpired(t *testing.T) {
	tr := ttl.New(1 * time.Millisecond)
	tr.Track("FAST", "v1")
	time.Sleep(5 * time.Millisecond)

	tr2 := ttl.New(10 * time.Minute)
	_ = tr2
	tr.entries = tr.ExpiredKeysHelper(tr)

	expired := tr.ExpiredKeys()
	if len(expired) != 1 || expired[0] != "FAST" {
		t.Errorf("expected [FAST], got %v", expired)
	}
}

func TestEvict_RemovesKey(t *testing.T) {
	tr := ttl.New(5 * time.Minute)
	tr.Track("TOKEN", "abc")
	tr.Evict("TOKEN")
	if tr.Get("TOKEN") != nil {
		t.Error("expected nil after eviction")
	}
}

func TestSummary_NotTracked(t *testing.T) {
	tr := ttl.New(5 * time.Minute)
	s := tr.Summary("UNKNOWN")
	if s != "UNKNOWN: not tracked" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_ValidEntry(t *testing.T) {
	tr := ttl.New(10 * time.Minute)
	tr.Track("MY_KEY", "val")
	s := tr.Summary("MY_KEY")
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
