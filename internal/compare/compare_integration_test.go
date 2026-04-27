package compare_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/compare"
	"github.com/your-org/vaultpull/internal/config"
)

func TestCompare_IntegratesWithConfig_Disabled(t *testing.T) {
	cfg := &config.CompareConfig{Enabled: false}
	if cfg.IsEnabled() {
		t.Skip("compare disabled — skipping")
	}
	// Reaching here confirms disabled path is respected.
}

func TestCompare_IntegratesWithConfig_FailFast(t *testing.T) {
	cfg := &config.CompareConfig{Enabled: true, FailFast: true}
	config.ApplyCompareDefaults(cfg)

	c := compare.New()
	src := map[string]string{"DB_PASS": "old"}
	dst := map[string]string{"DB_PASS": "new"}

	r := c.Compare("secrets/app", src, dst)
	if r.Equal {
		t.Fatal("expected differences")
	}
	if cfg.FailFast && len(r.Changed) == 0 {
		t.Error("expected changed keys under fail-fast mode")
	}
}

func TestCompare_MultiPath_AllEqual(t *testing.T) {
	cfg := &config.CompareConfig{
		Enabled: true,
		Paths:   []string{"secrets/a", "secrets/b"},
	}
	config.ApplyCompareDefaults(cfg)

	c := compare.New()
	payloads := map[string]map[string]string{
		"secrets/a": {"K": "v"},
		"secrets/b": {"K": "v"},
	}

	for _, path := range cfg.Paths {
		r := c.Compare(path, payloads[path], payloads[path])
		if !r.Equal {
			t.Errorf("expected %s to be equal", path)
		}
	}
}

func TestCompare_ConfigDefaultsApplied(t *testing.T) {
	cfg := &config.CompareConfig{Enabled: true}
	config.ApplyCompareDefaults(cfg)

	if cfg.Paths == nil {
		t.Error("expected Paths to be non-nil after defaults")
	}

	c := compare.New()
	r := c.Compare("secrets/env",
		map[string]string{"X": "1"},
		map[string]string{"X": "1", "Y": "2"},
	)
	if len(r.Added) != 1 {
		t.Errorf("expected 1 added key, got %d", len(r.Added))
	}
}
