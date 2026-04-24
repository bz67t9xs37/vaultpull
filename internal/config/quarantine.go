package config

// QuarantineConfig controls the quarantine feature.
type QuarantineConfig struct {
	Enabled  bool   `yaml:"enabled"`
	StoreDir string `yaml:"store_dir"`
}

// DefaultQuarantineConfig returns sensible defaults.
func DefaultQuarantineConfig() QuarantineConfig {
	return QuarantineConfig{
		Enabled:  false,
		StoreDir: ".vaultpull/quarantine",
	}
}

// ApplyQuarantineDefaults fills zero-value fields with defaults.
func ApplyQuarantineDefaults(c *QuarantineConfig) {
	if c == nil {
		return
	}
	def := DefaultQuarantineConfig()
	if c.StoreDir == "" {
		c.StoreDir = def.StoreDir
	}
}

// IsEnabled reports whether quarantine is active.
func (q *QuarantineConfig) IsEnabled() bool {
	return q != nil && q.Enabled
}
