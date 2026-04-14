package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWrite_And_Parse_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"API_KEY": "supersecret",
	}

	if err := Write(path, secrets); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	parsed, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	for k, v := range secrets {
		if parsed[k] != v {
			t.Errorf("key %s: expected %q got %q", k, v, parsed[k])
		}
	}
}

func TestParse_NonExistentFile(t *testing.T) {
	result, err := Parse("/tmp/vaultpull_nonexistent_xyz.env")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestParse_IgnoresComments(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	content := "# this is a comment\nFOO=bar\n\nBAZ=qux\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	result, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if result["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", result["FOO"])
	}
	if result["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", result["BAZ"])
	}
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}
