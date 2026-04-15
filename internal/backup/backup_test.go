package backup_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/internal/backup"
)

func TestCreate_WritesBackupFile(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(src, []byte("KEY=value\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	bakDir := filepath.Join(tmpDir, "backups")
	m := backup.New(bakDir)

	path, err := m.Create(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("backup file not found at %q: %v", path, err)
	}

	data, _ := os.ReadFile(path)
	if string(data) != "KEY=value\n" {
		t.Errorf("backup content mismatch: got %q", string(data))
	}
}

func TestCreate_MissingSource(t *testing.T) {
	m := backup.New(t.TempDir())
	_, err := m.Create("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing source, got nil")
	}
}

func TestRestore_WritesOriginalContent(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(src, []byte("ORIG=1\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	m := backup.New(filepath.Join(tmpDir, "backups"))
	bakPath, _ := m.Create(src)

	// Overwrite the original
	_ = os.WriteFile(src, []byte("NEW=2\n"), 0o600)

	if err := m.Restore(bakPath, src); err != nil {
		t.Fatalf("restore failed: %v", err)
	}

	data, _ := os.ReadFile(src)
	if string(data) != "ORIG=1\n" {
		t.Errorf("expected restored content, got %q", string(data))
	}
}

func TestList_ReturnsBackupFiles(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, ".env")
	_ = os.WriteFile(src, []byte("A=1\n"), 0o600)

	bakDir := filepath.Join(tmpDir, "backups")
	m := backup.New(bakDir)

	m.Create(src) //nolint
	m.Create(src) //nolint

	files, err := m.List(".env")
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("expected 2 backups, got %d", len(files))
	}
	for _, f := range files {
		if !strings.HasSuffix(f, ".bak") {
			t.Errorf("expected .bak suffix, got %q", f)
		}
	}
}
