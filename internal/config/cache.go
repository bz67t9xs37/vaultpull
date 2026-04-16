package config

import "time"

// CacheConfig holds configuration for the secret cache.
type CacheConfig struct {
	Enabled bool          `yaml:"enabled"`
	TTL     time.Duration `yaml:"ttl"`
	Dir     string        `yaml:"dir"`
}

// DefaultCacheConfig returns a CacheConfig with sensible defaults.
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		Enabled: true,
		TTL:     5 * time.Minute,
		Dir:     ".vaultpull/cache",
	}
}

// ApplyCacheDefaults fills zero-value fields with defaults.
func ApplyCacheDefaults(c *CacheConfig) {
	if c == nil {
		return
	}
	def := DefaultCacheConfig()
	if c.TTL == 0 {
		c.TTL = def.TTL
	}
	if c.Dir == "" {
		c.Dir = def.Dir
	}
}

// HasCache reports whether caching is enabled and configured.
func HasCache(c *CacheConfig) bool {
	return c != nil && c.Enabled
}
