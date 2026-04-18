package config

import (
	"testing"
	"time"
)

func TestDefaultPromoteConfig_Values(t *testing.T) {
	c := DefaultPromoteConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false")
	}
	if !c.DryRun {
		t.Error("expected DryRun to be true")
	}
	if c.Timeout != 30*time.Second {
		t.Errorf("unexpected timeout: %v", c.Timeout)
	}
}

func TestApplyPromoteDefaults_NilSafe(t *testing.T) {
	c := ApplyPromoteDefaults(nil)
	if c == nil {
		t.Fatal("expected non-nil result")
	}
	if c.Timeout != 30*time.Second {
		t.Errorf("unexpected timeout: %v", c.Timeout)
	}
}

func TestApplyPromoteDefaults_FillsZeroTimeout(t *testing.T) {
	c := &PromoteConfig{Enabled: true}
	result := ApplyPromoteDefaults(c)
	if result.Timeout != 30*time.Second {
		t.Errorf("expected default timeout, got %v", result.Timeout)
	}
}

func TestApplyPromoteDefaults_PreservesExistingValues(t *testing.T) {
	c := &PromoteConfig{
		Enabled: true,
		Source:  "staging",
		Timeout: 60 * time.Second,
	}
	result := ApplyPromoteDefaults(c)
	if result.Source != "staging" {
		t.Errorf("expected source to be preserved, got %q", result.Source)
	}
	if result.Timeout != 60*time.Second {
		t.Errorf("expected timeout to be preserved, got %v", result.Timeout)
	}
}

func TestIsEnabled_Promote(t *testing.T) {
	if (*PromoteConfig)(nil).IsEnabled() {
		t.Error("nil config should not be enabled")
	}
	c := &PromoteConfig{Enabled: true}
	if !c.IsEnabled() {
		t.Error("expected enabled")
	}
}

func TestIsDryRun_Promote(t *testing.T) {
	if !(*PromoteConfig)(nil).IsDryRun() {
		t.Error("nil config should default to dry run")
	}
	c := &PromoteConfig{DryRun: false}
	if c.IsDryRun() {
		t.Error("expected dry run to be false")
	}
}
