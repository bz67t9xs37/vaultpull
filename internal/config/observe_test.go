package config

import "testing"

func TestDefaultObserveConfig_Values(t *testing.T) {
	c := DefaultObserveConfig()
	if !c.Enabled {
		t.Error("expected Enabled to be true")
	}
	if c.Level != "INFO" {
		t.Errorf("expected Level INFO, got %s", c.Level)
	}
	if c.Target != "stderr" {
		t.Errorf("expected Target stderr, got %s", c.Target)
	}
}

func TestApplyObserveDefaults_NilSafe(t *testing.T) {
	ApplyObserveDefaults(nil) // must not panic
}

func TestApplyObserveDefaults_FillsLevel(t *testing.T) {
	c := &ObserveConfig{Enabled: true}
	ApplyObserveDefaults(c)
	if c.Level != "INFO" {
		t.Errorf("expected Level INFO, got %s", c.Level)
	}
}

func TestApplyObserveDefaults_FillsTarget(t *testing.T) {
	c := &ObserveConfig{Enabled: true, Level: "WARN"}
	ApplyObserveDefaults(c)
	if c.Target != "stderr" {
		t.Errorf("expected Target stderr, got %s", c.Target)
	}
}

func TestApplyObserveDefaults_PreservesExistingValues(t *testing.T) {
	c := &ObserveConfig{
		Enabled: true,
		Level:   "WARN",
		Target:  "stdout",
	}
	ApplyObserveDefaults(c)
	if c.Level != "WARN" {
		t.Errorf("expected Level WARN to be preserved, got %s", c.Level)
	}
	if c.Target != "stdout" {
		t.Errorf("expected Target stdout to be preserved, got %s", c.Target)
	}
}

func TestIsEnabled_Observe_True(t *testing.T) {
	c := &ObserveConfig{Enabled: true}
	if !c.IsEnabled() {
		t.Error("expected IsEnabled to return true")
	}
}

func TestIsEnabled_Observe_False(t *testing.T) {
	c := &ObserveConfig{Enabled: false}
	if c.IsEnabled() {
		t.Error("expected IsEnabled to return false")
	}
}
