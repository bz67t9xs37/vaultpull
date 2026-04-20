package rollback_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/rollback"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "rollback-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func writeBackup(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeBackup: %v", err)
	}
	return p
}

func TestLatest_ReturnsNilWhenNoBackup(t *testing.T) {
	dir := tempDir(t)
	rb := rollback.New(dir)
	entry, err := rb.Latest(filepath.Join(dir, ".env"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Errorf("expected nil, got %+v", entry)
	}
}

func TestLatest_ReturnsMostRecentBackup(t *testing.T) {
	dir := tempDir(t)
	writeBackup(t, dir, ".env.2024-01-01T00:00:00.bak", "OLD=1")
	time.Sleep(10 * time.Millisecond)
	writeBackup(t, dir, ".env.2024-01-02T00:00:00.bak", "NEW=2")

	rb := rollback.New(dir)
	entry, err := rb.Latest(filepath.Join(dir, ".env"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == nil {
		t.Fatal("expected entry, got nil")
	}
	if filepath.Base(entry.BackupFile) != ".env.2024-01-02T00:00:00.bak" {
		t.Errorf("unexpected backup file: %s", entry.BackupFile)
	}
}

func TestRestore_WritesBackupContentToTarget(t *testing.T) {
	dir := tempDir(t)
	const original = "SECRET=abc123\n"
	writeBackup(t, dir, ".env.2024-01-01T12:00:00.bak", original)

	target := filepath.Join(dir, ".env")
	if err := os.WriteFile(target, []byte("SECRET=changed\n"), 0o600); err != nil {
		t.Fatalf("setup target: %v", err)
	}

	rb := rollback.New(dir)
	entry, err := rb.Restore(target)
	if err != nil {
		t.Fatalf("Restore: %v", err)
	}
	if entry == nil {
		t.Fatal("expected entry, got nil")
	}
	got, _ := os.ReadFile(target)
	if string(got) != original {
		t.Errorf("want %q, got %q", original, string(got))
	}
}

func TestRestore_ErrorWhenNoBackup(t *testing.T) {
	dir := tempDir(t)
	rb := rollback.New(dir)
	_, err := rb.Restore(filepath.Join(dir, ".env"))
	if err == nil {
		t.Error("expected error, got nil")
	}
}
