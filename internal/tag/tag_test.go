package tag_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/tag"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "tag-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestSet_And_Get(t *testing.T) {
	m, err := tag.New(tempDir(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := m.Set("v1.0", "secret/app/prod", map[string]string{"env": "prod"}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := m.Get("v1.0")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got == nil {
		t.Fatal("expected tag, got nil")
	}
	if got.Path != "secret/app/prod" {
		t.Errorf("Path = %q, want %q", got.Path, "secret/app/prod")
	}
	if got.Meta["env"] != "prod" {
		t.Errorf("Meta[env] = %q, want %q", got.Meta["env"], "prod")
	}
}

func TestGet_UnknownTag_ReturnsNil(t *testing.T) {
	m, _ := tag.New(tempDir(t))
	got, err := m.Get("nonexistent")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
}

func TestSet_EmptyName_ReturnsError(t *testing.T) {
	m, _ := tag.New(tempDir(t))
	if err := m.Set("", "secret/app", nil); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestSet_EmptyPath_ReturnsError(t *testing.T) {
	m, _ := tag.New(tempDir(t))
	if err := m.Set("v1", "", nil); err == nil {
		t.Error("expected error for empty path")
	}
}

func TestDelete_RemovesTag(t *testing.T) {
	m, _ := tag.New(tempDir(t))
	_ = m.Set("v2", "secret/app", nil)
	if err := m.Delete("v2"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	got, _ := m.Get("v2")
	if got != nil {
		t.Error("expected nil after delete")
	}
}

func TestDelete_NonExistent_NoError(t *testing.T) {
	m, _ := tag.New(tempDir(t))
	if err := m.Delete("ghost"); err != nil {
		t.Errorf("Delete non-existent: %v", err)
	}
}

func TestAll_ReturnsSortedTags(t *testing.T) {
	dir := tempDir(t)
	m, _ := tag.New(dir)
	_ = m.Set("beta", "secret/beta", nil)
	_ = m.Set("alpha", "secret/alpha", nil)
	_ = m.Set("gamma", "secret/gamma", nil)

	// write a non-json file to ensure it is skipped
	_ = os.WriteFile(filepath.Join(dir, "ignore.txt"), []byte("x"), 0o600)

	tags, err := m.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(tags) != 3 {
		t.Fatalf("len = %d, want 3", len(tags))
	}
	if tags[0].Name != "alpha" || tags[1].Name != "beta" || tags[2].Name != "gamma" {
		t.Errorf("unexpected order: %v", tags)
	}
}
