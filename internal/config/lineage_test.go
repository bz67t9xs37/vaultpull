package config_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/config"
)

func TestDefaultLineageConfig_Values(t *testing.T) {
	cfg := config.DefaultLineageConfig()
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.StoreDir == "" {
		t.Error("expected non-empty StoreDir")
	}
	if cfg.MaxHistory <= 0 {
		t.Errorf("expected positive MaxHistory, got %d", cfg.MaxHistory)
	}
}

func TestApplyLineageDefaults_NilSafe(t *testing.T) {
	cfg := &config.Config{}
	cfg.Lineage = nil
	config.ApplyLineageDefaults(cfg)
	if cfg.Lineage == nil {
		t.Fatal("expected Lineage to be initialised")
	}
}

func TestApplyLineageDefaults_FillsStoreDir(t *testing.T) {
	cfg := &config.Config{
		Lineage: &config.LineageConfig{},
	}
	config.ApplyLineageDefaults(cfg)
	if cfg.Lineage.StoreDir == "" {
		t.Error("expected StoreDir to be filled")
	}
}

func TestApplyLineageDefaults_FillsMaxHistory(t *testing.T) {
	cfg := &config.Config{
		Lineage: &config.LineageConfig{},
	}
	config.ApplyLineageDefaults(cfg)
	if cfg.Lineage.MaxHistory <= 0 {
		t.Errorf("expected positive MaxHistory after defaults, got %d", cfg.Lineage.MaxHistory)
	}
}

func TestApplyLineageDefaults_PreservesExistingValues(t *testing.T) {
	cfg := &config.Config{
		Lineage: &config.LineageConfig{
			StoreDir:   "/custom/lineage",
			MaxHistory: 99,
			Enabled:    true,
		},
	}
	config.ApplyLineageDefaults(cfg)
	if cfg.Lineage.StoreDir != "/custom/lineage" {
		t.Errorf("expected StoreDir to be preserved, got %q", cfg.Lineage.StoreDir)
	}
	if cfg.Lineage.MaxHistory != 99 {
		t.Errorf("expected MaxHistory 99, got %d", cfg.Lineage.MaxHistory)
	}
}

func TestIsEnabled_Lineage_True(t *testing.T) {
	cfg := &config.LineageConfig{Enabled: true}
	if !cfg.IsEnabled() {
		t.Error("expected IsEnabled to return true")
	}
}

func TestIsEnabled_Lineage_False(t *testing.T) {
	cfg := &config.LineageConfig{Enabled: false}
	if cfg.IsEnabled() {
		t.Error("expected IsEnabled to return false")
	}
}
