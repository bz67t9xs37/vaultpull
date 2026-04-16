package config

// FilterConfig holds include/exclude rules for secret key filtering.
type FilterConfig struct {
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
}

// DefaultFilterConfig returns a FilterConfig with no rules (all keys allowed).
func DefaultFilterConfig() FilterConfig {
	return FilterConfig{
		Include: []string{},
		Exclude: []string{},
	}
}

// ApplyFilterDefaults fills in zero-value fields with defaults.
func ApplyFilterDefaults(f *FilterConfig) {
	if f == nil {
		return
	}
	if f.Include == nil {
		f.Include = []string{}
	}
	if f.Exclude == nil {
		f.Exclude = []string{}
	}
}

// HasRules reports whether any include or exclude rules are configured.
func (f FilterConfig) HasRules() bool {
	return len(f.Include) > 0 || len(f.Exclude) > 0
}
