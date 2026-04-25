package lineage_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/lineage"
)

func TestLineage_UsesConfigMaxHistory(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.LineageConfig{
		StoreDir:   dir,
		MaxHistory: 3,
		Enabled:    true,
	}

	tracker := lineage.New(cfg.StoreDir, cfg.MaxHistory)

	for i := 0; i < 5; i++ {
		if err := tracker.Record("secret/app", "KEY", "value"); err != nil {
			t.Fatalf("Record failed: %v", err)
		}
	}

	history, err := tracker.History("secret/app", "KEY")
	if err != nil {
		t.Fatalf("History failed: %v", err)
	}
	if len(history) > cfg.MaxHistory {
		t.Errorf("expected at most %d history entries, got %d", cfg.MaxHistory, len(history))
	}
}

func TestLineage_ConfigDefaultsApplied(t *testing.T) {
	cfg := &config.Config{}
	config.ApplyLineageDefaults(cfg)

	if cfg.Lineage.StoreDir == "" {
		t.Fatal("expected StoreDir after defaults")
	}
	if cfg.Lineage.MaxHistory <= 0 {
		t.Fatalf("expected positive MaxHistory after defaults, got %d", cfg.Lineage.MaxHistory)
	}

	tracker := lineage.New(cfg.Lineage.StoreDir, cfg.Lineage.MaxHistory)
	if err := tracker.Record("secret/app", "DB_PASS", "s3cr3t"); err != nil {
		t.Fatalf("Record failed: %v", err)
	}

	history, err := tracker.History("secret/app", "DB_PASS")
	if err != nil {
		t.Fatalf("History failed: %v", err)
	}
	if len(history) == 0 {
		t.Error("expected at least one history entry")
	}
}

func TestLineage_DisabledViaConfig(t *testing.T) {
	cfg := &config.LineageConfig{
		Enabled: false,
	}
	if cfg.IsEnabled() {
		t.Error("expected lineage to be disabled")
	}
}
