package config

// ResolveConfig controls how secret keys are resolved and mapped to local names.
type ResolveConfig struct {
	// StripPrefix removes a leading prefix from resolved key names.
	StripPrefix string `yaml:"strip_prefix"`
	// AddPrefix prepends a string to all resolved key names.
	AddPrefix string `yaml:"add_prefix"`
	// FlattenPath collapses nested path segments into the key name using an underscore.
	FlattenPath bool `yaml:"flatten_path"`
}

// DefaultResolveConfig returns a ResolveConfig with sensible defaults.
func DefaultResolveConfig() *ResolveConfig {
	return &ResolveConfig{
		StripPrefix: "",
		AddPrefix:   "",
		FlattenPath: false,
	}
}

// ApplyResolveDefaults fills zero-value fields with defaults.
func ApplyResolveDefaults(c *ResolveConfig) {
	if c == nil {
		return
	}
	d := DefaultResolveConfig()
	// All fields are optional strings/bools; nothing to fill unless explicitly zero.
	_ = d
}

// HasStripPrefix returns true if a strip prefix is configured.
func (c *ResolveConfig) HasStripPrefix() bool {
	return c != nil && c.StripPrefix != ""
}

// HasAddPrefix returns true if an add prefix is configured.
func (c *ResolveConfig) HasAddPrefix() bool {
	return c != nil && c.AddPrefix != ""
}
