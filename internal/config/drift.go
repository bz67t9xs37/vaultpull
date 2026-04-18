package config

// DriftConfig controls drift detection behaviour.
type DriftConfig struct {
	Enabled  bool   `yaml:"enabled"`
	FailFast bool   `yaml:"fail_fast"`
	LogPath  string `yaml:"log_path"`
}

// DefaultDriftConfig returns sensible defaults for drift detection.
func DefaultDriftConfig() *DriftConfig {
	return &DriftConfig{
		Enabled:  false,
		FailFast: false,
		LogPath:  "",
	}
}

// ApplyDriftDefaults fills zero-value fields with defaults.
func ApplyDriftDefaults(c *DriftConfig) {
	if c == nil {
		return
	}
	def := DefaultDriftConfig()
	if c.LogPath == "" {
		c.LogPath = def.LogPath
	}
}

// IsEnabled returns true when drift detection is active.
func (c *DriftConfig) IsEnabled() bool {
	return c != nil && c.Enabled
}
