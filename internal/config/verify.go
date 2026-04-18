package config

// VerifyConfig controls post-sync secret verification behaviour.
type VerifyConfig struct {
	Enabled         bool     `yaml:"enabled"`
	RequireNonEmpty bool     `yaml:"require_non_empty"`
	ExpectedKeys    []string `yaml:"expected_keys"`
}

// DefaultVerifyConfig returns sensible defaults.
func DefaultVerifyConfig() *VerifyConfig {
	return &VerifyConfig{
		Enabled:         false,
		RequireNonEmpty: true,
		ExpectedKeys:    []string{},
	}
}

// ApplyVerifyDefaults fills zero-value fields with defaults.
func ApplyVerifyDefaults(c *VerifyConfig) *VerifyConfig {
	if c == nil {
		return DefaultVerifyConfig()
	}
	if c.ExpectedKeys == nil {
		c.ExpectedKeys = []string{}
	}
	return c
}

// IsEnabled returns true when verification is active.
func (c *VerifyConfig) IsEnabled() bool {
	return c != nil && c.Enabled
}

// HasExpectedKeys returns true when at least one key is configured.
func (c *VerifyConfig) HasExpectedKeys() bool {
	return c != nil && len(c.ExpectedKeys) > 0
}
