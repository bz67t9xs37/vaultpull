package config

// TemplateConfig holds configuration for secret template rendering.
type TemplateConfig struct {
	Enabled    bool   `yaml:"enabled"`
	Dir        string `yaml:"dir"`
	LeftDelim  string `yaml:"left_delim"`
	RightDelim string `yaml:"right_delim"`
}

// DefaultTemplateConfig returns a TemplateConfig with sensible defaults.
func DefaultTemplateConfig() *TemplateConfig {
	return &TemplateConfig{
		Enabled:    false,
		Dir:        "templates",
		LeftDelim:  "{{",
		RightDelim: "}}",
	}
}

// ApplyTemplateDefaults fills zero-value fields with defaults.
func ApplyTemplateDefaults(c *TemplateConfig) {
	if c == nil {
		return
	}
	def := DefaultTemplateConfig()
	if c.Dir == "" {
		c.Dir = def.Dir
	}
	if c.LeftDelim == "" {
		c.LeftDelim = def.LeftDelim
	}
	if c.RightDelim == "" {
		c.RightDelim = def.RightDelim
	}
}

// IsEnabled returns true when template rendering is active.
func (c *TemplateConfig) IsEnabled() bool {
	return c != nil && c.Enabled
}
