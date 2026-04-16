package config

// TransformRule defines a single key transformation rule.
type TransformRule struct {
	Key       string `yaml:"key"`
	Prefix    string `yaml:"prefix"`
	Suffix    string `yaml:"suffix"`
	Uppercase bool   `yaml:"uppercase"`
	Replace   string `yaml:"replace"`
	With      string `yaml:"with"`
}

// TransformConfig holds transformation settings for a sync target.
type TransformConfig struct {
	Rules []TransformRule `yaml:"rules"`
}

// DefaultTransformConfig returns a TransformConfig with empty rules.
func DefaultTransformConfig() TransformConfig {
	return TransformConfig{
		Rules: []TransformRule{},
	}
}

// ApplyTransformDefaults fills in zero-value fields with defaults.
func ApplyTransformDefaults(cfg *TransformConfig) {
	if cfg == nil {
		return
	}
	if cfg.Rules == nil {
		cfg.Rules = []TransformRule{}
	}
}

// HasRules reports whether any transform rules are defined.
func (t *TransformConfig) HasRules() bool {
	return t != nil && len(t.Rules) > 0
}
