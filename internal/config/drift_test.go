package config_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/config"
)

func TestDefaultDriftConfig_Values(t *testing.T) {
	c := config.DefaultDriftConfig()
	if c.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if c.FailFast {
		t.Error("expected FailFast=false by default")
	}
	if c.LogPath != "" {
		t.Errorf("expected empty LogPath, got %q", c.LogPath)
	}
}

func TestApplyDriftDefaults_NilSafe(t *testing.T) {
	config.ApplyDriftDefaults(nil) // must not panic
}

func TestApplyDriftDefaults_PreservesExistingValues(t *testing.T) {
	c := &config.DriftConfig{
		Enabled:  true,
		FailFast: true,
		LogPath:  "/var/log/drift.log",
	}
	config.ApplyDriftDefaults(c)
	if !c.Enabled {
		t.Error("Enabled should be preserved")
	}
	if c.LogPath != "/var/log/drift.log" {
		t.Errorf("LogPath overwritten: %s", c.LogPath)
	}
}

func TestIsEnabled_Drift_True(t *testing.T) {
	c := &config.DriftConfig{Enabled: true}
	if !c.IsEnabled() {
		t.Error("expected IsEnabled=true")
	}
}

func TestIsEnabled_Drift_False(t *testing.T) {
	c := &config.DriftConfig{Enabled: false}
	if c.IsEnabled() {
		t.Error("expected IsEnabled=false")
	}
}

func TestIsEnabled_Drift_Nil(t *testing.T) {
	var c *config.DriftConfig
	if c.IsEnabled() {
		t.Error("nil config should return false")
	}
}
