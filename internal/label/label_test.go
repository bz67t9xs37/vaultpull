package label_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/label"
)

func tempStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "labels.json")
}

func TestSet_And_Get(t *testing.T) {
	s, err := label.New(tempStore(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if err := s.Set("secret/app", label.Labels{"env": "prod", "team": "platform"}); err != nil {
		t.Fatalf("Set: %v", err)
	}

	got := s.Get("secret/app")
	if got["env"] != "prod" || got["team"] != "platform" {
		t.Errorf("unexpected labels: %v", got)
	}
}

func TestGet_UnknownPath_ReturnsNil(t *testing.T) {
	s, _ := label.New(tempStore(t))
	if got := s.Get("secret/missing"); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestSet_MergesLabels(t *testing.T) {
	s, _ := label.New(tempStore(t))
	_ = s.Set("secret/app", label.Labels{"env": "staging"})
	_ = s.Set("secret/app", label.Labels{"team": "ops"})

	got := s.Get("secret/app")
	if got["env"] != "staging" || got["team"] != "ops" {
		t.Errorf("merge failed: %v", got)
	}
}

func TestRemove_DeletesKey(t *testing.T) {
	s, _ := label.New(tempStore(t))
	_ = s.Set("secret/app", label.Labels{"env": "prod", "owner": "alice"})
	_ = s.Remove("secret/app", "owner")

	got := s.Get("secret/app")
	if _, ok := got["owner"]; ok {
		t.Error("expected owner to be removed")
	}
	if got["env"] != "prod" {
		t.Error("expected env to remain")
	}
}

func TestRemove_LastKey_RemovesPath(t *testing.T) {
	s, _ := label.New(tempStore(t))
	_ = s.Set("secret/app", label.Labels{"env": "prod"})
	_ = s.Remove("secret/app", "env")

	if all := s.All(); len(all) != 0 {
		t.Errorf("expected empty store, got %v", all)
	}
}

func TestAll_ReturnsAllPaths(t *testing.T) {
	s, _ := label.New(tempStore(t))
	_ = s.Set("secret/a", label.Labels{"x": "1"})
	_ = s.Set("secret/b", label.Labels{"y": "2"})

	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}

func TestPersistence_ReloadsFromDisk(t *testing.T) {
	path := tempStore(t)
	s1, _ := label.New(path)
	_ = s1.Set("secret/db", label.Labels{"tier": "critical"})

	s2, err := label.New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	got := s2.Get("secret/db")
	if got["tier"] != "critical" {
		t.Errorf("expected tier=critical after reload, got %v", got)
	}
}

func TestNew_MissingDir_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "labels.json")
	s, _ := label.New(path)
	_ = s.Set("secret/x", label.Labels{"k": "v"})

	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
