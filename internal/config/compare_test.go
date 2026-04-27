package config_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/config"
)

func TestDefaultCompareConfig_Values(t *testing.T) {
	c := config.DefaultCompareConfig()
	if c.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if c.Paths == nil {
		t.Error("expected non-nil Paths slice")
	}
	if c.FailFast {
		t.Error("expected FailFast=false by default")
	}
}

func TestApplyCompareDefaults_NilSafe(t *testing.T) {
	config.ApplyCompareDefaults(nil) // must not panic
}

func TestApplyCompareDefaults_FillsNilPaths(t *testing.T) {
	c := &config.CompareConfig{Enabled: true}
	config.ApplyCompareDefaults(c)
	if c.Paths == nil {
		t.Error("expected Paths to be initialised")
	}
}

func TestApplyCompareDefaults_PreservesExistingValues(t *testing.T) {
	c := &config.CompareConfig{
		Enabled:  true,
		Paths:    []string{"secrets/a", "secrets/b"},
		FailFast: true,
	}
	config.ApplyCompareDefaults(c)
	if len(c.Paths) != 2 {
		t.Errorf("expected 2 paths, got %d", len(c.Paths))
	}
	if !c.FailFast {
		t.Error("expected FailFast to remain true")
	}
}

func TestIsEnabled_Compare_True(t *testing.T) {
	c := &config.CompareConfig{Enabled: true}
	if !c.IsEnabled() {
		t.Error("expected IsEnabled=true")
	}
}

func TestIsEnabled_Compare_False(t *testing.T) {
	c := &config.CompareConfig{Enabled: false}
	if c.IsEnabled() {
		t.Error("expected IsEnabled=false")
	}
}

func TestHasPaths_True(t *testing.T) {
	c := &config.CompareConfig{Paths: []string{"secrets/x"}}
	if !c.HasPaths() {
		t.Error("expected HasPaths=true")
	}
}

func TestHasPaths_False(t *testing.T) {
	c := &config.CompareConfig{}
	if c.HasPaths() {
		t.Error("expected HasPaths=false")
	}
}
