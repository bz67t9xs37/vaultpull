package config

import (
	"testing"
	"time"
)

func TestDefaultRateLimitConfig_Values(t *testing.T) {
	cfg := DefaultRateLimitConfig()
	if cfg.MaxRequests != 60 {
		t.Errorf("expected MaxRequests=60, got %d", cfg.MaxRequests)
	}
	if cfg.Window != time.Minute {
		t.Errorf("expected Window=1m, got %s", cfg.Window)
	}
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestApplyRateLimitDefaults_FillsZeroMaxRequests(t *testing.T) {
	cfg := &RateLimitConfig{Window: 30 * time.Second}
	ApplyRateLimitDefaults(cfg)
	if cfg.MaxRequests != 60 {
		t.Errorf("expected MaxRequests=60, got %d", cfg.MaxRequests)
	}
	if cfg.Window != 30*time.Second {
		t.Errorf("expected Window unchanged, got %s", cfg.Window)
	}
}

func TestApplyRateLimitDefaults_FillsZeroWindow(t *testing.T) {
	cfg := &RateLimitConfig{MaxRequests: 10}
	ApplyRateLimitDefaults(cfg)
	if cfg.Window != time.Minute {
		t.Errorf("expected Window=1m, got %s", cfg.Window)
	}
	if cfg.MaxRequests != 10 {
		t.Errorf("expected MaxRequests=10, got %d", cfg.MaxRequests)
	}
}

func TestApplyRateLimitDefaults_PreservesExistingValues(t *testing.T) {
	cfg := &RateLimitConfig{
		MaxRequests: 100,
		Window:      5 * time.Minute,
		Enabled:     false,
	}
	ApplyRateLimitDefaults(cfg)
	if cfg.MaxRequests != 100 {
		t.Errorf("expected MaxRequests=100, got %d", cfg.MaxRequests)
	}
	if cfg.Window != 5*time.Minute {
		t.Errorf("expected Window=5m, got %s", cfg.Window)
	}
	if cfg.Enabled {
		t.Error("expected Enabled=false to be preserved")
	}
}

func TestApplyRateLimitDefaults_NilSafe(t *testing.T) {
	// Should not panic
	ApplyRateLimitDefaults(nil)
}
