package config

import "time"

// PromoteConfig defines settings for promoting secrets between environments.
type PromoteConfig struct {
	Enabled     bool          `yaml:"enabled"`
	Source      string        `yaml:"source"`
	Destination string        `yaml:"destination"`
	DryRun      bool          `yaml:"dry_run"`
	Timeout     time.Duration `yaml:"timeout"`
}

func DefaultPromoteConfig() *PromoteConfig {
	return &PromoteConfig{
		Enabled: false,
		DryRun:  true,
		Timeout: 30 * time.Second,
	}
}

func ApplyPromoteDefaults(c *PromoteConfig) *PromoteConfig {
	if c == nil {
		return DefaultPromoteConfig()
	}
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}
	return c
}

func (c *PromoteConfig) IsEnabled() bool {
	return c != nil && c.Enabled
}

func (c *PromoteConfig) IsDryRun() bool {
	return c == nil || c.DryRun
}
