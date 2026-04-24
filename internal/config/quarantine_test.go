package config_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/config"
)

func TestDefaultQuarantineConfig_Values(t *testing.T) {
	c := config.DefaultQuarantineConfig()
	if c.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if c.StoreDir != ".vaultpull/quarantine" {
		t.Errorf("unexpected StoreDir: %s", c.StoreDir)
	}
}

func TestApplyQuarantineDefaults_FillsStoreDir(t *testing.T) {
	c := &config.QuarantineConfig{}
	config.ApplyQuarantineDefaults(c)
	if c.StoreDir == "" {
		t.Error("expected StoreDir to be filled")
	}
}

func TestApplyQuarantineDefaults_PreservesExistingValues(t *testing.T) {
	c := &config.QuarantineConfig{
		Enabled:  true,
		StoreDir: "/custom/quarantine",
	}
	config.ApplyQuarantineDefaults(c)
	if !c.Enabled {
		t.Error("expected Enabled to remain true")
	}
	if c.StoreDir != "/custom/quarantine" {
		t.Errorf("expected StoreDir to remain /custom/quarantine, got %s", c.StoreDir)
	}
}

func TestApplyQuarantineDefaults_NilSafe(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ApplyQuarantineDefaults panicked on nil: %v", r)
		}
	}()
	config.ApplyQuarantineDefaults(nil)
}

func TestIsEnabled_Quarantine_True(t *testing.T) {
	c := &config.QuarantineConfig{Enabled: true, StoreDir: ".vaultpull/quarantine"}
	if !c.IsEnabled() {
		t.Error("expected IsEnabled=true")
	}
}

func TestIsEnabled_Quarantine_False(t *testing.T) {
	c := &config.QuarantineConfig{Enabled: false}
	if c.IsEnabled() {
		t.Error("expected IsEnabled=false")
	}
}
