package config

// PinConfig controls secret version pinning behaviour.
type PinConfig struct {
	Enabled  bool   `yaml:"enabled"`
	StoreDir string `yaml:"store_dir"`
}

// DefaultPinConfig returns sensible defaults.
func DefaultPinConfig() PinConfig {
	return PinConfig{
		Enabled:  false,
		StoreDir: ".vaultpull/pins",
	}
}

// ApplyPinDefaults fills zero-value fields with defaults.
func ApplyPinDefaults(c *PinConfig) {
	if c == nil {
		return
	}
	d := DefaultPinConfig()
	if c.StoreDir == "" {
		c.StoreDir = d.StoreDir
	}
}

// IsEnabled reports whether pinning is active.
func (c *PinConfig) IsEnabled() bool {
	return c != nil && c.Enabled
}
