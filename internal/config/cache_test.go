package config

import (
	"testing"
	"time"
)

func TestDefaultCacheConfig_Values(t *testing.T) {
	c := DefaultCacheConfig()
	if !c.Enabled {
		t.Error("expected Enabled to be true")
	}
	if c.TTL != 5*time.Minute {
		t.Errorf("expected TTL 5m, got %v", c.TTL)
	}
	if c.Dir == "" {
		t.Error("expected non-empty Dir")
	}
}

func TestApplyCacheDefaults_FillsZeroTTL(t *testing.T) {
	c := &CacheConfig{Enabled: true}
	ApplyCacheDefaults(c)
	if c.TTL != 5*time.Minute {
		t.Errorf("expected TTL 5m, got %v", c.TTL)
	}
}

func TestApplyCacheDefaults_FillsEmptyDir(t *testing.T) {
	c := &CacheConfig{Enabled: true, TTL: time.Minute}
	ApplyCacheDefaults(c)
	if c.Dir == "" {
		t.Error("expected Dir to be filled")
	}
}

func TestApplyCacheDefaults_PreservesExistingValues(t *testing.T) {
	c := &CacheConfig{
		Enabled: true,
		TTL:     10 * time.Minute,
		Dir:     "/tmp/mycache",
	}
	ApplyCacheDefaults(c)
	if c.TTL != 10*time.Minute {
		t.Errorf("expected TTL preserved, got %v", c.TTL)
	}
	if c.Dir != "/tmp/mycache" {
		t.Errorf("expected Dir preserved, got %s", c.Dir)
	}
}

func TestApplyCacheDefaults_NilSafe(t *testing.T) {
	ApplyCacheDefaults(nil) // must not panic
}

func TestHasCache_Enabled(t *testing.T) {
	c := &CacheConfig{Enabled: true}
	if !HasCache(c) {
		t.Error("expected HasCache true")
	}
}

func TestHasCache_Disabled(t *testing.T) {
	c := &CacheConfig{Enabled: false}
	if HasCache(c) {
		t.Error("expected HasCache false")
	}
}

func TestHasCache_Nil(t *testing.T) {
	if HasCache(nil) {
		t.Error("expected HasCache false for nil")
	}
}
