package config

import "time"

// RotateConfig holds configuration for secret rotation scheduling.
type RotateConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Interval time.Duration `yaml:"interval"`
	StateDir string        `yaml:"state_dir"`
}

// DefaultRotateConfig returns sensible defaults for rotation.
func DefaultRotateConfig() *RotateConfig {
	return &RotateConfig{
		Enabled:  false,
		Interval: 24 * time.Hour,
		StateDir: ".vaultpull/rotate",
	}
}

// ApplyRotateDefaults fills zero-value fields with defaults.
func ApplyRotateDefaults(c *RotateConfig) *RotateConfig {
	if c == nil {
		return DefaultRotateConfig()
	}
	def := DefaultRotateConfig()
	if c.Interval == 0 {
		c.Interval = def.Interval
	}
	if c.StateDir == "" {
		c.StateDir = def.StateDir
	}
	return c
}

// IsEnabled returns true if rotation is enabled.
func (c *RotateConfig) IsEnabled() bool {
	return c != nil && c.Enabled
}
