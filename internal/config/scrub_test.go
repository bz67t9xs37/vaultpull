package config

import (
	"testing"
)

func TestDefaultScrubConfig_Values(t *testing.T) {
	cfg := DefaultScrubConfig()
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if !cfg.Enabled {
		t.Error("expected Enabled to be true by default")
	}
	if cfg.MinLength == 0 {
		t.Error("expected MinLength to be non-zero")
	}
}

func TestApplyScrubDefaults_NilSafe(t *testing.T) {
	cfg := &Config{Scrub: nil}
	ApplyScrubDefaults(cfg)
	if cfg.Scrub == nil {
		t.Fatal("expected Scrub to be initialized")
	}
}

func TestApplyScrubDefaults_FillsMinLength(t *testing.T) {
	cfg := &Config{Scrub: &ScrubConfig{}}
	ApplyScrubDefaults(cfg)
	if cfg.Scrub.MinLength == 0 {
		t.Error("expected MinLength to be filled")
	}
}

func TestApplyScrubDefaults_PreservesExistingValues(t *testing.T) {
	cfg := &Config{Scrub: &ScrubConfig{
		Enabled:   false,
		MinLength: 12,
	}}
	ApplyScrubDefaults(cfg)
	if cfg.Scrub.Enabled != false {
		t.Error("expected Enabled to remain false")
	}
	if cfg.Scrub.MinLength != 12 {
		t.Errorf("expected MinLength 12, got %d", cfg.Scrub.MinLength)
	}
}

func TestIsEnabled_Scrub_True(t *testing.T) {
	cfg := &ScrubConfig{Enabled: true}
	if !cfg.IsEnabled() {
		t.Error("expected IsEnabled to return true")
	}
}

func TestIsEnabled_Scrub_False(t *testing.T) {
	cfg := &ScrubConfig{Enabled: false}
	if cfg.IsEnabled() {
		t.Error("expected IsEnabled to return false")
	}
}
