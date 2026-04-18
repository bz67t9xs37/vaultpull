package config

import "time"

// PromoteConfig holds settings for secret promotion between Vault paths.
type PromoteConfig struct {
	Enabled     bool          `yaml:"enabled"`
	DryRun      bool          `yaml:"dry_run"`
	Timeout     time.Duration `yaml:"timeout"`
	SourceMount string        `yaml:"source_mount"`
	DestMount   string        `yaml:"dest_mount"`
}

func DefaultPromoteConfig() *PromoteConfig {
	return &PromoteConfig{
		Enabled:     false,
		DryRun:      false,
		Timeout:     30 * time.Second,
		SourceMount: "secret",
		DestMount:   "secret",
	}
}

func ApplyPromoteDefaults(c *PromoteConfig) *PromoteConfig {
	if c == nil {
		return DefaultPromoteConfig()
	}
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}
	if c.SourceMount == "" {
		c.SourceMount = "secret"
	}
	if c.DestMount == "" {
		c.DestMount = "secret"
	}
	return c
}

func (c *PromoteConfig) IsEnabled() bool {
	return c != nil && c.Enabled
}

func (c *PromoteConfig) IsDryRun() bool {
	return c != nil && c.DryRun
}
