package config_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/config"
)

func TestDefaultPromoteConfig_Values(t *testing.T) {
	c := config.DefaultPromoteConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false")
	}
	if c.DryRun {
		t.Error("expected DryRun to be false")
	}
	if c.Timeout != 30*time.Second {
		t.Errorf("expected Timeout 30s, got %v", c.Timeout)
	}
}

func TestApplyPromoteDefaults_NilSafe(t *testing.T) {
	c := config.ApplyPromoteDefaults(nil)
	if c == nil {
		t.Fatal("expected non-nil config")
	}
	if c.Timeout != 30*time.Second {
		t.Errorf("expected default timeout, got %v", c.Timeout)
	}
}

func TestApplyPromoteDefaults_FillsZeroTimeout(t *testing.T) {
	c := &config.PromoteConfig{Enabled: true}
	result := config.ApplyPromoteDefaults(c)
	if result.Timeout != 30*time.Second {
		t.Errorf("expected timeout filled, got %v", result.Timeout)
	}
}

func TestApplyPromoteDefaults_PreservesExistingValues(t *testing.T) {
	c := &config.PromoteConfig{
		Enabled: true,
		DryRun:  true,
		Timeout: 60 * time.Second,
		Source:  "staging",
		Destination: "production",
	}
	result := config.ApplyPromoteDefaults(c)
	if result.Timeout != 60*time.Second {
		t.Errorf("expected preserved timeout, got %v", result.Timeout)
	}
	if result.Source != "staging" {
		t.Errorf("expected source preserved, got %q", result.Source)
	}
}

func TestIsEnabled_Promote(t *testing.T) {
	c := &config.PromoteConfig{Enabled: true}
	if !c.IsEnabled() {
		t.Error("expected IsEnabled true")
	}
	var nilC *config.PromoteConfig
	if nilC.IsEnabled() {
		t.Error("expected nil IsEnabled false")
	}
}

func TestIsDryRun_Promote(t *testing.T) {
	c := &config.PromoteConfig{DryRun: true}
	if !c.IsDryRun() {
		t.Error("expected IsDryRun true")
	}
	var nilC *config.PromoteConfig
	if nilC.IsDryRun() {
		t.Error("expected nil IsDryRun false")
	}
}
