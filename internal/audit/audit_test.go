package audit_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/audit"
)

func tempLogPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit.log")
}

func TestRecord_And_ReadAll(t *testing.T) {
	path := tempLogPath(t)
	l, err := audit.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	e := audit.Entry{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Operation: "sync",
		Target:    ".env",
		Added:     3,
		Modified:  1,
	}
	if err := l.Record(e); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Operation != "sync" {
		t.Errorf("operation: got %q, want %q", entries[0].Operation, "sync")
	}
	if entries[0].Added != 3 {
		t.Errorf("added: got %d, want 3", entries[0].Added)
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	path := tempLogPath(t)
	l, err := audit.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	for i := 0; i < 3; i++ {
		if err := l.Record(audit.Entry{Operation: "sync", Target: ".env"}); err != nil {
			t.Fatalf("Record %d: %v", i, err)
		}
	}

	entries, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestReadAll_EmptyFile(t *testing.T) {
	path := tempLogPath(t)
	l, err := audit.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	entries, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestRecord_SetsTimestampIfZero(t *testing.T) {
	path := tempLogPath(t)
	l, _ := audit.New(path)
	before := time.Now().UTC()
	_ = l.Record(audit.Entry{Operation: "rotate", Target: ".env"})
	after := time.Now().UTC()

	entries, _ := l.ReadAll()
	if len(entries) != 1 {
		t.Fatal("expected 1 entry")
	}
	ts := entries[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", ts, before, after)
	}
}

func TestNew_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "logs")
	path := filepath.Join(dir, "audit.log")
	_, err := audit.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("expected directory %q to be created", dir)
	}
}
