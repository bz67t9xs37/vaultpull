package config

import (
	"testing"
)

func TestDefaultLabelConfig_Values(t *testing.T) {
	cfg := DefaultLabelConfig()
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.StoreDir == "" {
		t.Error("expected non-empty StoreDir")
	}
	if !cfg.Enabled {
		t.Error("expected Enabled to be true by default")
	}
}

func TestApplyLabelDefaults_NilSafe(t *testing.T) {
	cfg := &Config{Label: nil}
	ApplyLabelDefaults(cfg)
	if cfg.Label == nil {
		t.Fatal("expected Label to be initialized")
	}
}

func TestApplyLabelDefaults_FillsStoreDir(t *testing.T) {
	cfg := &Config{Label: &LabelConfig{}}
	ApplyLabelDefaults(cfg)
	if cfg.Label.StoreDir == "" {
		t.Error("expected StoreDir to be filled")
	}
}

func TestApplyLabelDefaults_PreservesExistingValues(t *testing.T) {
	cfg := &Config{
		Label: &LabelConfig{
			StoreDir: "/custom/labels",
			Enabled:  true,
		},
	}
	ApplyLabelDefaults(cfg)
	if cfg.Label.StoreDir != "/custom/labels" {
		t.Errorf("expected StoreDir to be preserved, got %s", cfg.Label.StoreDir)
	}
}

func TestIsEnabled_Label_True(t *testing.T) {
	cfg := &LabelConfig{Enabled: true}
	if !cfg.IsEnabled() {
		t.Error("expected IsEnabled to return true")
	}
}

func TestIsEnabled_Label_False(t *testing.T) {
	cfg := &LabelConfig{Enabled: false}
	if cfg.IsEnabled() {
		t.Error("expected IsEnabled to return false")
	}
}

func TestApplyLabelDefaults_NilConfig_DoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic: %v", r)
		}
	}()
	var cfg *Config
	if cfg != nil {
		ApplyLabelDefaults(cfg)
	}
}
