package config

import "time"

// RateLimitConfig holds rate limiting configuration for Vault API calls.
type RateLimitConfig struct {
	// MaxRequests is the maximum number of requests allowed per Window.
	MaxRequests int `yaml:"max_requests"`
	// Window is the duration over which MaxRequests is enforced.
	Window time.Duration `yaml:"window"`
	// Enabled controls whether rate limiting is active.
	Enabled bool `yaml:"enabled"`
}

// DefaultRateLimitConfig returns a RateLimitConfig with sensible defaults.
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		MaxRequests: 60,
		Window:      time.Minute,
		Enabled:     true,
	}
}

// ApplyRateLimitDefaults fills zero-value fields with defaults.
func ApplyRateLimitDefaults(cfg *RateLimitConfig) {
	if cfg == nil {
		return
	}
	def := DefaultRateLimitConfig()
	if cfg.MaxRequests == 0 {
		cfg.MaxRequests = def.MaxRequests
	}
	if cfg.Window == 0 {
		cfg.Window = def.Window
	}
}
