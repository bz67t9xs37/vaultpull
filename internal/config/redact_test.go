package config

import (
	"testing"
)

func TestDefaultRedactConfig_Values(t *testing.T) {
	cfg := DefaultRedactConfig()
	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
	if cfg.MinLength != 4 {
		t.Errorf("expected MinLength 4, got %d", cfg.MinLength)
	}
	if cfg.MaskChar != "*" {
		t.Errorf("expected MaskChar '*', got %q", cfg.MaskChar)
	}
	if cfg.Patterns == nil {
		t.Error("expected Patterns to be non-nil")
	}
}

func TestApplyRedactDefaults_FillsZeroValues(t *testing.T) {
	cfg := &RedactConfig{}
	ApplyRedactDefaults(cfg)
	if cfg.MinLength != 4 {
		t.Errorf("expected MinLength 4, got %d", cfg.MinLength)
	}
	if cfg.MaskChar != "*" {
		t.Errorf("expected MaskChar '*', got %q", cfg.MaskChar)
	}
}

func TestApplyRedactDefaults_PreservesExistingValues(t *testing.T) {
	cfg := &RedactConfig{
		MinLength: 8,
		MaskChar:  "#",
		Patterns:  []string{"secret_.*"},
	}
	ApplyRedactDefaults(cfg)
	if cfg.MinLength != 8 {
		t.Errorf("expected MinLength 8, got %d", cfg.MinLength)
	}
	if cfg.MaskChar != "#" {
		t.Errorf("expected MaskChar '#', got %q", cfg.MaskChar)
	}
	if len(cfg.Patterns) != 1 {
		t.Errorf("expected 1 pattern, got %d", len(cfg.Patterns))
	}
}

func TestApplyRedactDefaults_NilSafe(t *testing.T) {
	ApplyRedactDefaults(nil) // must not panic
}

func TestIsEnabled_True(t *testing.T) {
	cfg := &RedactConfig{Enabled: true}
	if !cfg.IsEnabled() {
		t.Error("expected IsEnabled true")
	}
}

func TestIsEnabled_NilReturnsFalse(t *testing.T) {
	var cfg *RedactConfig
	if cfg.IsEnabled() {
		t.Error("expected IsEnabled false for nil")
	}
}
