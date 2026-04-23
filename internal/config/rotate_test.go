package config

import (
	"testing"
	"time"
)

func TestDefaultRotateConfig_Values(t *testing.T) {
	c := DefaultRotateConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if c.Interval != 24*time.Hour {
		t.Errorf("expected 24h interval, got %v", c.Interval)
	}
	if c.StateDir == "" {
		t.Error("expected non-empty StateDir")
	}
}

func TestApplyRotateDefaults_NilSafe(t *testing.T) {
	c := ApplyRotateDefaults(nil)
	if c == nil {
		t.Fatal("expected non-nil result")
	}
	if c.Interval != 24*time.Hour {
		t.Errorf("unexpected interval: %v", c.Interval)
	}
}

func TestApplyRotateDefaults_FillsZeroInterval(t *testing.T) {
	c := &RotateConfig{Enabled: true}
	result := ApplyRotateDefaults(c)
	if result.Interval != 24*time.Hour {
		t.Errorf("expected default interval, got %v", result.Interval)
	}
	if result.StateDir == "" {
		t.Error("expected StateDir to be filled")
	}
}

func TestApplyRotateDefaults_PreservesExistingValues(t *testing.T) {
	c := &RotateConfig{
		Enabled:  true,
		Interval: 6 * time.Hour,
		StateDir: "/custom/state",
	}
	result := ApplyRotateDefaults(c)
	if result.Interval != 6*time.Hour {
		t.Errorf("expected preserved interval, got %v", result.Interval)
	}
	if result.StateDir != "/custom/state" {
		t.Errorf("expected preserved StateDir, got %s", result.StateDir)
	}
}

func TestIsEnabled_True(t *testing.T) {
	c := &RotateConfig{Enabled: true}
	if !c.IsEnabled() {
		t.Error("expected IsEnabled to return true")
	}
}

func TestIsEnabled_False(t *testing.T) {
	c := &RotateConfig{Enabled: false}
	if c.IsEnabled() {
		t.Error("expected IsEnabled to return false")
	}
}

func TestIsEnabled_Nil(t *testing.T) {
	var c *RotateConfig
	if c.IsEnabled() {
		t.Error("expected nil config to return false")
	}
}

func TestApplyRotateDefaults_FillsEmptyStateDir(t *testing.T) {
	// Ensure that a config with a non-zero interval but empty StateDir
	// still gets a default StateDir applied.
	c := &RotateConfig{
		Enabled:  true,
		Interval: 12 * time.Hour,
		StateDir: "",
	}
	result := ApplyRotateDefaults(c)
	if result.StateDir == "" {
		t.Error("expected StateDir to be filled when empty")
	}
	if result.Interval != 12*time.Hour {
		t.Errorf("expected interval to be preserved, got %v", result.Interval)
	}
}
