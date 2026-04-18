package config_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/config"
)

func TestDefaultPinConfig_Values(t *testing.T) {
	d := config.DefaultPinConfig()
	if d.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if d.StoreDir == "" {
		t.Error("expected non-empty StoreDir")
	}
}

func TestApplyPinDefaults_FillsStoreDir(t *testing.T) {
	c := &config.PinConfig{}
	config.ApplyPinDefaults(c)
	if c.StoreDir == "" {
		t.Error("expected StoreDir to be filled")
	}
}

func TestApplyPinDefaults_PreservesExistingValues(t *testing.T) {
	c := &config.PinConfig{StoreDir: "/custom/pins", Enabled: true}
	config.ApplyPinDefaults(c)
	if c.StoreDir != "/custom/pins" {
		t.Errorf("StoreDir changed to %q", c.StoreDir)
	}
	if !c.Enabled {
		t.Error("Enabled should remain true")
	}
}

func TestApplyPinDefaults_NilSafe(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panicked on nil: %v", r)
		}
	}()
	config.ApplyPinDefaults(nil)
}

func TestIsEnabled_Pin_True(t *testing.T) {
	c := &config.PinConfig{Enabled: true}
	if !c.IsEnabled() {
		t.Error("expected IsEnabled=true")
	}
}

func TestIsEnabled_Pin_False(t *testing.T) {
	c := &config.PinConfig{Enabled: false}
	if c.IsEnabled() {
		t.Error("expected IsEnabled=false")
	}
}
