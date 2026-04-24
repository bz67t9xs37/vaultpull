package checksum_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/checksum"
)

func TestCompute_DeterministicOutput(t *testing.T) {
	c := checksum.New()
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}

	r1 := c.Compute("secret/app", secrets)
	r2 := c.Compute("secret/app", secrets)

	if r1.Checksum != r2.Checksum {
		t.Errorf("expected deterministic checksum, got %q and %q", r1.Checksum, r2.Checksum)
	}
}

func TestCompute_KeyOrderIndependent(t *testing.T) {
	c := checksum.New()
	a := map[string]string{"Z_KEY": "1", "A_KEY": "2"}
	b := map[string]string{"A_KEY": "2", "Z_KEY": "1"}

	ra := c.Compute("secret/app", a)
	rb := c.Compute("secret/app", b)

	if ra.Checksum != rb.Checksum {
		t.Errorf("expected same checksum regardless of key order, got %q and %q", ra.Checksum, rb.Checksum)
	}
}

func TestCompute_DifferentValuesProduceDifferentChecksums(t *testing.T) {
	c := checksum.New()
	original := map[string]string{"API_KEY": "secret-v1"}
	updated := map[string]string{"API_KEY": "secret-v2"}

	r1 := c.Compute("secret/app", original)
	r2 := c.Compute("secret/app", updated)

	if r1.Checksum == r2.Checksum {
		t.Error("expected different checksums for different values")
	}
}

func TestCompute_SetsKeyCount(t *testing.T) {
	c := checksum.New()
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}

	r := c.Compute("secret/app", secrets)
	if r.KeyCount != 3 {
		t.Errorf("expected KeyCount=3, got %d", r.KeyCount)
	}
}

func TestEqual_SameChecksum(t *testing.T) {
	c := checksum.New()
	secrets := map[string]string{"FOO": "bar"}

	r1 := c.Compute("secret/app", secrets)
	r2 := c.Compute("secret/app", secrets)

	if !c.Equal(r1, r2) {
		t.Error("expected Equal to return true for identical secrets")
	}
}

func TestChanged_DetectsModification(t *testing.T) {
	c := checksum.New()
	original := map[string]string{"TOKEN": "abc"}
	prev := c.Compute("secret/svc", original)

	updated := map[string]string{"TOKEN": "xyz"}
	if !c.Changed(prev, "secret/svc", updated) {
		t.Error("expected Changed to return true after value modification")
	}
}

func TestChanged_NoModification(t *testing.T) {
	c := checksum.New()
	secrets := map[string]string{"TOKEN": "abc"}
	prev := c.Compute("secret/svc", secrets)

	if c.Changed(prev, "secret/svc", secrets) {
		t.Error("expected Changed to return false when secrets are unchanged")
	}
}
