package config

// MaskConfig holds configuration for secret masking behaviour.
type MaskConfig struct {
	Enabled      bool     `yaml:"enabled"`
	MaskChar     string   `yaml:"mask_char"`
	CustomKeys   []string `yaml:"custom_keys"`
}

// DefaultMaskConfig returns sensible defaults for masking.
func DefaultMaskConfig() MaskConfig {
	return MaskConfig{
		Enabled:    true,
		MaskChar:   "*",
		CustomKeys: []string{},
	}
}

// ApplyMaskDefaults fills zero-value fields with defaults.
func ApplyMaskDefaults(c *MaskConfig) {
	if c == nil {
		return
	}
	if c.MaskChar == "" {
		c.MaskChar = DefaultMaskConfig().MaskChar
	}
	if c.CustomKeys == nil {
		c.CustomKeys = []string{}
	}
}

// HasCustomKeys reports whether any custom key patterns are configured.
func (m *MaskConfig) HasCustomKeys() bool {
	return m != nil && len(m.CustomKeys) > 0
}
