package config

import (
	"testing"
)

func TestDefaultTransformConfig_Values(t *testing.T) {
	cfg := DefaultTransformConfig()
	if cfg.Rules == nil {
		t.Fatal("expected Rules to be non-nil")
	}
	if len(cfg.Rules) != 0 {
		t.Fatalf("expected 0 rules, got %d", len(cfg.Rules))
	}
}

func TestApplyTransformDefaults_FillsNilRules(t *testing.T) {
	cfg := &TransformConfig{}
	ApplyTransformDefaults(cfg)
	if cfg.Rules == nil {
		t.Fatal("expected Rules to be initialized")
	}
}

func TestApplyTransformDefaults_PreservesExistingRules(t *testing.T) {
	cfg := &TransformConfig{
		Rules: []TransformRule{
			{Key: "DB_PASS", Uppercase: true},
		},
	}
	ApplyTransformDefaults(cfg)
	if len(cfg.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(cfg.Rules))
	}
	if cfg.Rules[0].Key != "DB_PASS" {
		t.Errorf("expected key DB_PASS, got %s", cfg.Rules[0].Key)
	}
}

func TestApplyTransformDefaults_NilSafe(t *testing.T) {
	// should not panic
	ApplyTransformDefaults(nil)
}

func TestHasRules_NoRules(t *testing.T) {
	cfg := DefaultTransformConfig()
	if cfg.HasRules() {
		t.Fatal("expected HasRules to return false")
	}
}

func TestHasRules_WithRules(t *testing.T) {
	cfg := &TransformConfig{
		Rules: []TransformRule{
			{Prefix: "APP_"},
		},
	}
	if !cfg.HasRules() {
		t.Fatal("expected HasRules to return true")
	}
}

func TestHasRules_NilConfig(t *testing.T) {
	var cfg *TransformConfig
	if cfg.HasRules() {
		t.Fatal("expected HasRules on nil to return false")
	}
}
