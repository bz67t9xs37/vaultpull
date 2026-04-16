package config

// RedactConfig controls log/output redaction behaviour.
type RedactConfig struct {
	Enabled    bool     `yaml:"enabled"`
	MinLength  int      `yaml:"min_length"`
	Patterns   []string `yaml:"patterns"`
	MaskChar   string   `yaml:"mask_char"`
}

// DefaultRedactConfig returns sensible redaction defaults.
func DefaultRedactConfig() RedactConfig {
	return RedactConfig{
		Enabled:   true,
		MinLength: 4,
		Patterns:  []string{},
		MaskChar:  "*",
	}
}

// ApplyRedactDefaults fills zero-value fields with defaults.
func ApplyRedactDefaults(cfg *RedactConfig) {
	if cfg == nil {
		return
	}
	def := DefaultRedactConfig()
	if cfg.MinLength == 0 {
		cfg.MinLength = def.MinLength
	}
	if cfg.MaskChar == "" {
		cfg.MaskChar = def.MaskChar
	}
	if cfg.Patterns == nil {
		cfg.Patterns = def.Patterns
	}
}

// IsEnabled reports whether redaction is active.
func (r *RedactConfig) IsEnabled() bool {
	return r != nil && r.Enabled
}
