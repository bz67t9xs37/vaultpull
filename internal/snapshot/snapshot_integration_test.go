package snapshot_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/diff"
	"github.com/your-org/vaultpull/internal/snapshot"
)

// TestSnapshot_IntegratesWithDiff verifies that a snapshot can be used as the
// "before" state to compute a meaningful diff against new secrets.
func TestSnapshot_IntegratesWithDiff(t *testing.T) {
	store, err := snapshot.New(tempDir(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	original := map[string]string{
		"DB_HOST": "localhost",
		"DB_PASS": "old-pass",
		"API_KEY": "abc123",
	}
	if err := store.Save(snapshot.Entry{Path: "secret/app", Secrets: original}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	updated := map[string]string{
		"DB_HOST": "localhost",   // unchanged
		"DB_PASS": "new-pass",   // modified
		"NEW_KEY": "xyz",        // added
		// API_KEY removed
	}

	snap, err := store.Latest("secret/app")
	if err != nil || snap == nil {
		t.Fatalf("Latest: %v", err)
	}

	changes := diff.Compute(snap.Secrets, updated)
	summary := diff.Summary(changes)

	if summary.Added != 1 {
		t.Errorf("expected 1 added, got %d", summary.Added)
	}
	if summary.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", summary.Removed)
	}
	if summary.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", summary.Modified)
	}
	if summary.Unchanged != 1 {
		t.Errorf("expected 1 unchanged, got %d", summary.Unchanged)
	}
}
