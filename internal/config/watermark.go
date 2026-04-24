package config

// WatermarkConfig controls hash-based change detection across syncs.
type WatermarkConfig struct {
	// Enabled turns watermark tracking on or off.
	Enabled bool `yaml:"enabled"`

	// HashAlgorithm is the algorithm used to hash secret values.
	// Supported: "sha256" (default), "sha1".
	HashAlgorithm string `yaml:"hash_algorithm"`
}

// DefaultWatermarkConfig returns the baseline watermark configuration.
func DefaultWatermarkConfig() WatermarkConfig {
	return WatermarkConfig{
		Enabled:       true,
		HashAlgorithm: "sha256",
	}
}

// ApplyWatermarkDefaults fills zero-value fields with defaults.
// A nil pointer is treated as a no-op.
func ApplyWatermarkDefaults(c *WatermarkConfig) {
	if c == nil {
		return
	}
	if c.HashAlgorithm == "" {
		c.HashAlgorithm = "sha256"
	}
}

// IsEnabled returns true when watermark tracking is active.
func (c *WatermarkConfig) IsEnabled() bool {
	return c != nil && c.Enabled
}
