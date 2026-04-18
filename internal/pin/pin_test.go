package pin_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/pin"
)

func tempStore(t *testing.T) *pin.Store {
	t.Helper()
	dir := t.TempDir()
	s, err := pin.New(filepath.Join(dir, "pins.json"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func TestPin_And_Get(t *testing.T) {
	s := tempStore(t)
	if err := s.Pin("secret/app", 3); err != nil {
		t.Fatalf("Pin: %v", err)
	}
	e := s.Get("secret/app")
	if e == nil {
		t.Fatal("expected entry, got nil")
	}
	if e.Version != 3 {
		t.Errorf("version = %d, want 3", e.Version)
	}
}

func TestGet_UnpinnedPath_ReturnsNil(t *testing.T) {
	s := tempStore(t)
	if s.Get("secret/missing") != nil {
		t.Error("expected nil for unpinned path")
	}
}

func TestUnpin_RemovesEntry(t *testing.T) {
	s := tempStore(t)
	_ = s.Pin("secret/app", 1)
	if err := s.Unpin("secret/app"); err != nil {
		t.Fatalf("Unpin: %v", err)
	}
	if s.Get("secret/app") != nil {
		t.Error("expected nil after unpin")
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	s := tempStore(t)
	_ = s.Pin("secret/a", 1)
	_ = s.Pin("secret/b", 2)
	if len(s.All()) != 2 {
		t.Errorf("All() = %d entries, want 2", len(s.All()))
	}
}

func TestPin_PersistsAcrossReload(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")
	s1, _ := pin.New(path)
	_ = s1.Pin("secret/reload", 7)

	s2, err := pin.New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	e := s2.Get("secret/reload")
	if e == nil || e.Version != 7 {
		t.Errorf("expected version 7 after reload, got %v", e)
	}
}

func TestNew_MissingDir_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "pins.json")
	s, err := pin.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = s.Pin("secret/x", 1)
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}
