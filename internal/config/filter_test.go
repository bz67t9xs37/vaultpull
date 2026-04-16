package config

import (
	"testing"
)

func TestDefaultFilterConfig_Values(t *testing.T) {
	cfg := DefaultFilterConfig()
	if cfg.Include == nil {
		t.Error("expected Include to be non-nil")
	}
	if cfg.Exclude == nil {
		t.Error("expected Exclude to be non-nil")
	}
	if len(cfg.Include) != 0 {
		t.Errorf("expected empty Include, got %v", cfg.Include)
	}
	if len(cfg.Exclude) != 0 {
		t.Errorf("expected empty Exclude, got %v", cfg.Exclude)
	}
}

func TestApplyFilterDefaults_FillsNilSlices(t *testing.T) {
	cfg := &FilterConfig{}
	ApplyFilterDefaults(cfg)
	if cfg.Include == nil {
		t.Error("expected Include to be initialized")
	}
	if cfg.Exclude == nil {
		t.Error("expected Exclude to be initialized")
	}
}

func TestApplyFilterDefaults_PreservesExistingValues(t *testing.T) {
	cfg := &FilterConfig{
		Include: []string{"DB_"},
		Exclude: []string{"TEST_"},
	}
	ApplyFilterDefaults(cfg)
	if len(cfg.Include) != 1 || cfg.Include[0] != "DB_" {
		t.Errorf("expected Include to be preserved, got %v", cfg.Include)
	}
	if len(cfg.Exclude) != 1 || cfg.Exclude[0] != "TEST_" {
		t.Errorf("expected Exclude to be preserved, got %v", cfg.Exclude)
	}
}

func TestApplyFilterDefaults_NilSafe(t *testing.T) {
	// should not panic
	ApplyFilterDefaults(nil)
}

func TestHasRules_NoRules(t *testing.T) {
	cfg := DefaultFilterConfig()
	if cfg.HasRules() {
		t.Error("expected HasRules to return false for empty config")
	}
}

func TestHasRules_WithInclude(t *testing.T) {
	cfg := FilterConfig{Include: []string{"API_"}, Exclude: []string{}}
	if !cfg.HasRules() {
		t.Error("expected HasRules to return true when Include is set")
	}
}

func TestHasRules_WithExclude(t *testing.T) {
	cfg := FilterConfig{Include: []string{}, Exclude: []string{"LEGACY_"}}
	if !cfg.HasRules() {
		t.Error("expected HasRules to return true when Exclude is set")
	}
}
