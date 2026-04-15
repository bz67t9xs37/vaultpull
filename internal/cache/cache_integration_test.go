package cache_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/cache"
	"github.com/your-org/vaultpull/internal/diff"
)

// TestCache_IntegratesWithDiff verifies that a cached entry can be used as the
// "before" snapshot when computing a diff against freshly fetched secrets.
func TestCache_IntegratesWithDiff(t *testing.T) {
	c, err := cache.New(tempDir(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	initial := map[string]string{
		"DB_HOST": "localhost",
		"DB_PASS": "old-pass",
		"LEGACY":  "remove-me",
	}
	if err := c.Set("secret/app", initial); err != nil {
		t.Fatalf("Set initial: %v", err)
	}

	entry, err := c.Get("secret/app")
	if err != nil || entry == nil {
		t.Fatalf("Get: %v", err)
	}

	updated := map[string]string{
		"DB_HOST":  "localhost",
		"DB_PASS":  "new-pass",
		"NEW_FLAG": "true",
	}

	changes := diff.Compute(entry.Secrets, updated)
	summary := diff.Summary(changes)

	if summary.Added != 1 {
		t.Errorf("Added: got %d, want 1", summary.Added)
	}
	if summary.Removed != 1 {
		t.Errorf("Removed: got %d, want 1", summary.Removed)
	}
	if summary.Modified != 1 {
		t.Errorf("Modified: got %d, want 1", summary.Modified)
	}
	if summary.Unchanged != 1 {
		t.Errorf("Unchanged: got %d, want 1", summary.Unchanged)
	}

	// Update cache with latest secrets.
	if err := c.Set("secret/app", updated); err != nil {
		t.Fatalf("Set updated: %v", err)
	}
	newEntry, _ := c.Get("secret/app")
	if newEntry.Checksum == entry.Checksum {
		t.Error("cache checksum should differ after update")
	}
}
