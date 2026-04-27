package config

// LabelConfig controls secret path labelling behaviour.
type LabelConfig struct {
	Enabled  bool   `yaml:"enabled"`
	StoreDir string `yaml:"store_dir"`
}

// DefaultLabelConfig returns the default label configuration.
func DefaultLabelConfig() LabelConfig {
	return LabelConfig{
		Enabled:  false,
		StoreDir: ".vaultpull/labels",
	}
}

// ApplyLabelDefaults fills zero-value fields in c with defaults.
func ApplyLabelDefaults(c *LabelConfig) {
	if c == nil {
		return
	}
	def := DefaultLabelConfig()
	if c.StoreDir == "" {
		c.StoreDir = def.StoreDir
	}
}

// IsEnabled reports whether labelling is active.
func (c *LabelConfig) IsEnabled() bool {
	return c != nil && c.Enabled
}
