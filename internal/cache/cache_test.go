package cache_test

import (
	"os"
	"testing"

	"github.com/your-org/vaultpull/internal/cache"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "vaultpull-cache-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func TestSet_And_Get_RoundTrip(t *testing.T) {
	c, err := cache.New(tempDir(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	secrets := map[string]string{"DB_PASS": "s3cr3t", "API_KEY": "abc123"}
	if err := c.Set("secret/myapp", secrets); err != nil {
		t.Fatalf("Set: %v", err)
	}
	entry, err := c.Get("secret/myapp")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if entry == nil {
		t.Fatal("expected entry, got nil")
	}
	if entry.Secrets["DB_PASS"] != "s3cr3t" {
		t.Errorf("DB_PASS: got %q, want %q", entry.Secrets["DB_PASS"], "s3cr3t")
	}
	if entry.Checksum == "" {
		t.Error("expected non-empty checksum")
	}
}

func TestGet_MissingEntry_ReturnsNil(t *testing.T) {
	c, err := cache.New(tempDir(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	entry, err := c.Get("secret/nonexistent")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if entry != nil {
		t.Errorf("expected nil, got %+v", entry)
	}
}

func TestInvalidate_RemovesEntry(t *testing.T) {
	c, err := cache.New(tempDir(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := c.Set("secret/myapp", map[string]string{"KEY": "val"}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := c.Invalidate("secret/myapp"); err != nil {
		t.Fatalf("Invalidate: %v", err)
	}
	entry, err := c.Get("secret/myapp")
	if err != nil {
		t.Fatalf("Get after invalidate: %v", err)
	}
	if entry != nil {
		t.Error("expected nil after invalidation")
	}
}

func TestInvalidate_NonExistent_NoError(t *testing.T) {
	c, err := cache.New(tempDir(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := c.Invalidate("secret/ghost"); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestChecksum_ChangesWithSecrets(t *testing.T) {
	c, err := cache.New(tempDir(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := c.Set("secret/app", map[string]string{"K": "v1"}); err != nil {
		t.Fatalf("Set v1: %v", err)
	}
	e1, _ := c.Get("secret/app")

	if err := c.Set("secret/app", map[string]string{"K": "v2"}); err != nil {
		t.Fatalf("Set v2: %v", err)
	}
	e2, _ := c.Get("secret/app")

	if e1.Checksum == e2.Checksum {
		t.Error("expected different checksums for different secrets")
	}
}
