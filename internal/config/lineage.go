package config

// LineageConfig controls secret lineage tracking.
type LineageConfig struct {
	Enabled bool   `yaml:"enabled"`
	LogPath string `yaml:"log_path"`
}

// DefaultLineageConfig returns the default lineage configuration.
func DefaultLineageConfig() *LineageConfig {
	return &LineageConfig{
		Enabled: false,
		LogPath: ".vaultpull/lineage.json",
	}
}

// ApplyLineageDefaults fills zero-value fields with defaults.
func ApplyLineageDefaults(c *LineageConfig) *LineageConfig {
	if c == nil {
		return DefaultLineageConfig()
	}
	d := DefaultLineageConfig()
	if c.LogPath == "" {
		c.LogPath = d.LogPath
	}
	return c
}

// IsEnabled returns true when lineage tracking is active.
func (c *LineageConfig) IsEnabled() bool {
	return c != nil && c.Enabled
}
