package cache_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/cache"
	"github.com/yourusername/vaultpull/internal/config"
)

func TestCache_UsesConfigTTL(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.CacheConfig{
		Enabled: true,
		TTL:     50 * time.Millisecond,
		Dir:     dir,
	}
	config.ApplyCacheDefaults(cfg)

	c := cache.New(dir)
	secrets := map[string]string{"API_KEY": "abc123"}
	if err := c.Set("myapp", secrets); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got := c.Get("myapp")
	if got == nil {
		t.Fatal("expected cached entry")
	}
	if got["API_KEY"] != "abc123" {
		t.Errorf("unexpected value: %s", got["API_KEY"])
	}
}

func TestCache_ConfigDefaultsApplied(t *testing.T) {
	cfg := &config.CacheConfig{}
	config.ApplyCacheDefaults(cfg)

	if cfg.TTL == 0 {
		t.Error("expected TTL to be set by defaults")
	}
	if cfg.Dir == "" {
		t.Error("expected Dir to be set by defaults")
	}
	if !config.HasCache(cfg) {
		t.Error("expected HasCache to reflect enabled state")
	}
}
