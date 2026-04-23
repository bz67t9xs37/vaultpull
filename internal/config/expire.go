package config

import "time"

// ExpireConfig controls secret expiry checking behaviour.
type ExpireConfig struct {
	Enabled    bool          `yaml:"enabled"`
	WarnBefore time.Duration `yaml:"warn_before"`
}

// DefaultExpireConfig returns sensible defaults for expiry checking.
func DefaultExpireConfig() *ExpireConfig {
	return &ExpireConfig{
		Enabled:    false,
		WarnBefore: 48 * time.Hour,
	}
}

// ApplyExpireDefaults fills zero-value fields with defaults.
func ApplyExpireDefaults(c *ExpireConfig) {
	if c == nil {
		return
	}
	if c.WarnBefore == 0 {
		c.WarnBefore = DefaultExpireConfig().WarnBefore
	}
}

// IsEnabled returns true when expiry checking is active.
func (c *ExpireConfig) IsEnabled() bool {
	return c != nil && c.Enabled
}
