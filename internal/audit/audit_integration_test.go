package audit_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/audit"
	"github.com/your-org/vaultpull/internal/diff"
)

// EntryFromDiff builds an audit.Entry from a diff.Summary result.
func EntryFromDiff(op, target string, s diff.Summary, errMsg string) audit.Entry {
	return audit.Entry{
		Timestamp: time.Now().UTC(),
		Operation: op,
		Target:    target,
		Added:     s.Added,
		Removed:   s.Removed,
		Modified:  s.Modified,
		Unchanged: s.Unchanged,
		Error:     errMsg,
	}
}

func TestAudit_IntegratesWithDiffSummary(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.log")
	l, err := audit.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	old := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	new_ := map[string]string{
		"DB_HOST": "prod.db.internal",
		"API_KEY": "secret123",
	}

	changes := diff.Compute(old, new_)
	summary := diff.Summary(changes)

	entry := EntryFromDiff("sync", ".env.production", summary, "")
	if err := l.Record(entry); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	got := entries[0]
	if got.Added != 1 {
		t.Errorf("Added: got %d, want 1", got.Added)
	}
	if got.Removed != 1 {
		t.Errorf("Removed: got %d, want 1", got.Removed)
	}
	if got.Modified != 1 {
		t.Errorf("Modified: got %d, want 1", got.Modified)
	}
	if got.Target != ".env.production" {
		t.Errorf("Target: got %q, want \".env.production\"", got.Target)
	}
	if got.Operation != "sync" {
		t.Errorf("Operation: got %q, want \"sync\"", got.Operation)
	}
}
