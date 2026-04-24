package quarantine_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/quarantine"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "quarantine-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestAdd_And_List(t *testing.T) {
	q := quarantine.New(tempDir(t))
	if err := q.Add("secret/app", "DB_PASS", "contains forbidden pattern"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	entries, err := q.List("secret/app")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Key != "DB_PASS" {
		t.Errorf("expected DB_PASS, got %s", entries[0].Key)
	}
	if entries[0].Reason != "contains forbidden pattern" {
		t.Errorf("unexpected reason: %s", entries[0].Reason)
	}
	if entries[0].QuarantinedAt.IsZero() {
		t.Error("expected non-zero QuarantinedAt")
	}
}

func TestAdd_Idempotent(t *testing.T) {
	q := quarantine.New(tempDir(t))
	q.Add("secret/app", "API_KEY", "reason")
	q.Add("secret/app", "API_KEY", "reason")
	entries, _ := q.List("secret/app")
	if len(entries) != 1 {
		t.Errorf("expected 1 entry after duplicate add, got %d", len(entries))
	}
}

func TestRemove_LiftQuarantine(t *testing.T) {
	q := quarantine.New(tempDir(t))
	q.Add("secret/app", "TOKEN", "suspicious")
	q.Add("secret/app", "SECRET", "empty")
	if err := q.Remove("secret/app", "TOKEN"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	entries, _ := q.List("secret/app")
	if len(entries) != 1 || entries[0].Key != "SECRET" {
		t.Errorf("expected only SECRET remaining, got %+v", entries)
	}
}

func TestIsQuarantined_True(t *testing.T) {
	q := quarantine.New(tempDir(t))
	q.Add("secret/app", "DB_PASS", "reason")
	if !q.IsQuarantined("secret/app", "DB_PASS") {
		t.Error("expected DB_PASS to be quarantined")
	}
}

func TestIsQuarantined_False(t *testing.T) {
	q := quarantine.New(tempDir(t))
	if q.IsQuarantined("secret/app", "UNKNOWN") {
		t.Error("expected UNKNOWN to not be quarantined")
	}
}

func TestList_EmptyDir_ReturnsNil(t *testing.T) {
	q := quarantine.New(tempDir(t))
	entries, err := q.List("secret/missing")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty list, got %d", len(entries))
	}
}

func TestAdd_IsolatedByPath(t *testing.T) {
	dir := tempDir(t)
	q := quarantine.New(dir)
	q.Add("secret/app", "KEY_A", "r")
	q.Add("secret/other", "KEY_B", "r")
	a, _ := q.List("secret/app")
	b, _ := q.List("secret/other")
	if len(a) != 1 || a[0].Key != "KEY_A" {
		t.Errorf("wrong entries for secret/app: %+v", a)
	}
	if len(b) != 1 || b[0].Key != "KEY_B" {
		t.Errorf("wrong entries for secret/other: %+v", b)
	}
	_ = filepath.Join(dir, "app.quarantine.json") // sanity reference
}
