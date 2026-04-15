package ttl_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/ttl"
)

// ExpiredKeysHelper is a test-only shim to avoid exporting internal map.
// Real integration test exercises the full lifecycle.
func (t *ttl.Tracker) ExpiredKeysHelper(_ *ttl.Tracker) map[string]*ttl.Entry {
	return nil // unused shim
}

func TestTTL_FullLifecycle(t *testing.T) {
	tr := ttl.New(20 * time.Millisecond)

	keys := map[string]string{
		"API_KEY":    "key-abc",
		"DB_PASS":    "pass-xyz",
		"AUTH_TOKEN": "tok-123",
	}

	for k, v := range keys {
		tr.Track(k, v)
	}

	// All should be valid immediately after tracking.
	for k := range keys {
		e := tr.Get(k)
		if e == nil {
			t.Fatalf("expected entry for %s", k)
		}
		if e.IsExpired() {
			t.Errorf("%s should not be expired yet", k)
		}
	}

	// Wait for TTL to elapse.
	time.Sleep(40 * time.Millisecond)

	expired := tr.ExpiredKeys()
	if len(expired) != len(keys) {
		t.Errorf("expected %d expired keys, got %d", len(keys), len(expired))
	}

	// Evict all expired keys and confirm they're gone.
	for _, k := range expired {
		tr.Evict(k)
	}
	for k := range keys {
		if tr.Get(k) != nil {
			t.Errorf("expected %s to be evicted", k)
		}
	}
}

func TestTTL_MixedExpiryStates(t *testing.T) {
	tr := ttl.New(5 * time.Millisecond)
	tr.Track("SHORT", "v1")
	time.Sleep(15 * time.Millisecond)

	// Re-track with fresh timestamp to simulate a refresh.
	tr.Track("SHORT", "v2")

	e := tr.Get("SHORT")
	if e.IsExpired() {
		t.Error("re-tracked entry should not be expired")
	}
	if e.Value != "v2" {
		t.Errorf("expected refreshed value 'v2', got %q", e.Value)
	}
}
