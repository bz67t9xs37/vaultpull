package config_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/config"
)

func TestDefaultVerifyConfig_Values(t *testing.T) {
	cfg := config.DefaultVerifyConfig()
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.RequireNonEmpty {
		t.Error("expected RequireNonEmpty to default to false")
	}
	if len(cfg.Keys) != 0 {
		t.Error("expected Keys to be empty by default")
	}
}

func TestApplyVerifyDefaults_NilSafe(t *testing.T) {
	cfg := &config.VerifyConfig{}
	config.ApplyVerifyDefaults(cfg)
	if cfg.Keys == nil {
		t.Error("expected Keys to be initialized")
	}
}

func TestApplyVerifyDefaults_PreservesExistingValues(t *testing.T) {
	cfg := &config.VerifyConfig{
		RequireNonEmpty: true,
		Keys:            []string{"DB_URL", "API_KEY"},
	}
	config.ApplyVerifyDefaults(cfg)
	if !cfg.RequireNonEmpty {
		t.Error("expected RequireNonEmpty to be preserved")
	}
	if len(cfg.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(cfg.Keys))
	}
}

func TestIsEnabled_Verify_True(t *testing.T) {
	cfg := &config.VerifyConfig{
		Enabled: true,
		Keys:    []string{"SECRET"},
	}
	if !config.IsVerifyEnabled(cfg) {
		t.Error("expected verify to be enabled")
	}
}

func TestIsEnabled_Verify_False(t *testing.T) {
	cfg := &config.VerifyConfig{Enabled: false}
	if config.IsVerifyEnabled(cfg) {
		t.Error("expected verify to be disabled")
	}
}

func TestIsEnabled_Verify_NilSafe(t *testing.T) {
	if config.IsVerifyEnabled(nil) {
		t.Error("expected nil config to return false")
	}
}
