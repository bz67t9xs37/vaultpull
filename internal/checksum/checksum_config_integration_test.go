package checksum_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/checksum"
)

func TestChecksum_IntegratesWithDiff(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"API_KEY":     "abc123",
	}

	c := checksum.New()
	result, err := c.Compute(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Hash == "" {
		t.Error("expected non-empty hash")
	}
	if result.KeyCount != 3 {
		t.Errorf("expected key count 3, got %d", result.KeyCount)
	}

	// Same secrets in different order should produce the same checksum.
	reordered := map[string]string{
		"API_KEY":     "abc123",
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
	}

	result2, err := c.Compute(reordered)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !c.Equal(result, result2) {
		t.Error("expected checksums to be equal for same secrets in different order")
	}
}

func TestChecksum_DetectsModifiedValue(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PASS": "secret",
	}

	c := checksum.New()
	before, err := c.Compute(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	modified := map[string]string{
		"DB_HOST": "localhost",
		"DB_PASS": "changed-secret",
	}

	after, err := c.Compute(modified)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.Equal(before, after) {
		t.Error("expected checksums to differ after value change")
	}
}

func TestChecksum_EmptyMap(t *testing.T) {
	c := checksum.New()
	result, err := c.Compute(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.KeyCount != 0 {
		t.Errorf("expected key count 0, got %d", result.KeyCount)
	}
	if result.Hash == "" {
		t.Error("expected non-empty hash even for empty map")
	}
}
