package config_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/config"
)

func TestDefaultExpireConfig_Values(t *testing.T) {
	cfg := config.DefaultExpireConfig()
	if cfg == nil {
		t.Fatal("expected non-nil default expire config")
	}
	if cfg.WarningWindow == 0 {
		t.Error("expected non-zero default warning window")
	}
}

func TestApplyExpireDefaults_NilSafe(t *testing.T) {
	cfg := &config.Config{Expire: nil}
	config.ApplyExpireDefaults(cfg)
	if cfg.Expire == nil {
		t.Fatal("expected Expire to be initialized after ApplyExpireDefaults")
	}
}

func TestApplyExpireDefaults_FillsZeroWarningWindow(t *testing.T) {
	cfg := &config.Config{Expire: &config.ExpireConfig{}}
	config.ApplyExpireDefaults(cfg)
	if cfg.Expire.WarningWindow == 0 {
		t.Error("expected WarningWindow to be filled with default")
	}
}

func TestApplyExpireDefaults_PreservesExistingValues(t *testing.T) {
	import_duration := int64(999)
	cfg := &config.Config{
		Expire: &config.ExpireConfig{
			WarningWindow: import_duration,
		},
	}
	config.ApplyExpireDefaults(cfg)
	if cfg.Expire.WarningWindow != import_duration {
		t.Errorf("expected WarningWindow to be preserved, got %d", cfg.Expire.WarningWindow)
	}
}

func TestIsEnabled_Expire_True(t *testing.T) {
	cfg := &config.ExpireConfig{Enabled: true}
	if !cfg.IsEnabled() {
		t.Error("expected IsEnabled to return true")
	}
}

func TestIsEnabled_Expire_False(t *testing.T) {
	cfg := &config.ExpireConfig{Enabled: false}
	if cfg.IsEnabled() {
		t.Error("expected IsEnabled to return false")
	}
}
