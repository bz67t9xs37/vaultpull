package config

// FingerprintConfig controls fingerprint computation behaviour.
type FingerprintConfig struct {
	Enabled   bool   `yaml:"enabled"`
	StoreDir  string `yaml:"store_dir"`
}

// DefaultFingerprintConfig returns sensible defaults.
func DefaultFingerprintConfig() FingerprintConfig {
	return FingerprintConfig{
		Enabled:  true,
		StoreDir: ".vaultpull/fingerprints",
	}
}

// ApplyFingerprintDefaults fills zero-value fields with defaults.
func ApplyFingerprintDefaults(c *FingerprintConfig) {
	if c == nil {
		return
	}
	d := DefaultFingerprintConfig()
	if c.StoreDir == "" {
		c.StoreDir = d.StoreDir
	}
}

// IsEnabled returns true when fingerprinting is active.
func (c *FingerprintConfig) IsEnabled() bool {
	return c != nil && c.Enabled
}
