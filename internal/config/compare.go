package config

// CompareConfig controls cross-path secret comparison behaviour.
type CompareConfig struct {
	Enabled  bool     `yaml:"enabled"`
	Paths    []string `yaml:"paths"`
	FailFast bool     `yaml:"fail_fast"`
}

// DefaultCompareConfig returns sensible defaults for comparison.
func DefaultCompareConfig() *CompareConfig {
	return &CompareConfig{
		Enabled:  false,
		Paths:    []string{},
		FailFast: false,
	}
}

// ApplyCompareDefaults fills zero-value fields with defaults.
func ApplyCompareDefaults(c *CompareConfig) {
	if c == nil {
		return
	}
	if c.Paths == nil {
		c.Paths = []string{}
	}
}

// IsEnabled returns true when comparison is active.
func (c *CompareConfig) IsEnabled() bool {
	return c != nil && c.Enabled
}

// HasPaths returns true when at least one comparison path is configured.
func (c *CompareConfig) HasPaths() bool {
	return c != nil && len(c.Paths) > 0
}
