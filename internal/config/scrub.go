package config

// ScrubConfig controls output scrubbing of secret values.
type ScrubConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Replacement string `yaml:"replacement"`
}

// DefaultScrubConfig returns sensible defaults.
func DefaultScrubConfig() *ScrubConfig {
	return &ScrubConfig{
		Enabled:     true,
		Replacement: "[REDACTED]",
	}
}

// ApplyScrubDefaults fills zero-value fields on cfg with defaults.
// Safe to call with a nil pointer — returns a new default config.
func ApplyScrubDefaults(cfg *ScrubConfig) *ScrubConfig {
	if cfg == nil {
		return DefaultScrubConfig()
	}
	if cfg.Replacement == "" {
		cfg.Replacement = "[REDACTED]"
	}
	return cfg
}

// IsEnabled reports whether scrubbing is active.
func (c *ScrubConfig) IsEnabled() bool {
	return c != nil && c.Enabled
}
