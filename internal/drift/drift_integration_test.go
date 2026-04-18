package drift_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/drift"
)

func TestDrift_DisabledViaConfig(t *testing.T) {
	cfg := &config.DriftConfig{Enabled: false}
	if cfg.IsEnabled() {
		t.Fatal("drift should be disabled")
	}
	// When disabled the caller should skip Check; verify detector still works standalone.
	d := drift.New(map[string]string{"X": "1"})
	results := d.Check(".env", map[string]string{"X": "1"})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestDrift_FailFastFlagPropagates(t *testing.T) {
	cfg := &config.DriftConfig{Enabled: true, FailFast: true}
	config.ApplyDriftDefaults(cfg)
	if !cfg.FailFast {
		t.Error("FailFast should be preserved after ApplyDriftDefaults")
	}
}

func TestDrift_SummaryWithConfig(t *testing.T) {
	cfg := &config.DriftConfig{Enabled: true}
	if !cfg.IsEnabled() {
		t.Fatal("expected drift enabled")
	}
	d := drift.New(map[string]string{"DB_PASS": "secret", "API_KEY": "key"})
	results := d.Check("/app/.env", map[string]string{"DB_PASS": "old", "API_KEY": "key"})
	summary := drift.Summary(results)
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}
