package config

import "time"

// PromoteConfig holds configuration for secret promotion between environments.
type PromoteConfig struct {
	Enabled     bool          `yaml:"enabled"`
	Source      string        `yaml:"source"`
	Destination string        `yaml:"destination"`
	DryRun      bool          `yaml:"dry_run"`
	Timeout     time.Duration `yaml:"timeout"`
}

// DefaultPromoteConfig returns a PromoteConfig with sensible defaults.
func DefaultPromoteConfig() *PromoteConfig {
	return &PromoteConfig{
		Enabled: false,
		DryRun:  false,
		Timeout: 30 * time.Second,
	}
}

// ApplyPromoteDefaults fills zero-value fields with defaults.
func ApplyPromoteDefaults(c *PromoteConfig) *PromoteConfig {
	if c == nil {
		return DefaultPromoteConfig()
	}
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}
	return c
}

// IsEnabled returns true if promotion is enabled.
func (c *PromoteConfig) IsEnabled() bool {
	if c == nil {
		return false
	}
	return c.Enabled
}

// IsDryRun returns true if dry-run mode is active.
func (c *PromoteConfig) IsDryRun() bool {
	if c == nil {
		return false
	}
	return c.DryRun
}
