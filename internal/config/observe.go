package config

// ObserveConfig controls lifecycle event observation during sync.
type ObserveConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Level    string `yaml:"level"`
	Target   string `yaml:"target"` // "stderr", "stdout", or a file path
}

// DefaultObserveConfig returns a safe default ObserveConfig.
func DefaultObserveConfig() ObserveConfig {
	return ObserveConfig{
		Enabled: true,
		Level:   "INFO",
		Target:  "stderr",
	}
}

// ApplyObserveDefaults fills zero-value fields with defaults.
func ApplyObserveDefaults(c *ObserveConfig) {
	if c == nil {
		return
	}
	if c.Level == "" {
		c.Level = "INFO"
	}
	if c.Target == "" {
		c.Target = "stderr"
	}
}

// IsEnabled returns true when observation is active.
func (c *ObserveConfig) IsEnabled() bool {
	return c != nil && c.Enabled
}
