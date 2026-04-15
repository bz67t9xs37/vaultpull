package filter

import (
	"strings"
)

// Rule defines a single filter rule for secret keys.
type Rule struct {
	Prefix  string
	Suffix  string
	Contains string
}

// Filter applies inclusion/exclusion rules to secret maps.
type Filter struct {
	include []Rule
	exclude []Rule
}

// New creates a Filter with the given include and exclude rules.
func New(include, exclude []Rule) *Filter {
	return &Filter{
		include: include,
		exclude: exclude,
	}
}

// Apply returns a new map containing only the secrets that pass the filter rules.
// If no include rules are defined, all keys are considered included by default.
// Exclude rules are always applied after include rules.
func (f *Filter) Apply(secrets map[string]string) map[string]string {
	result := make(map[string]string)

	for key, val := range secrets {
		if len(f.include) > 0 && !matchesAny(key, f.include) {
			continue
		}
		if matchesAny(key, f.exclude) {
			continue
		}
		result[key] = val
	}

	return result
}

// matchesAny reports whether key matches at least one of the given rules.
func matchesAny(key string, rules []Rule) bool {
	for _, r := range rules {
		if r.Prefix != "" && !strings.HasPrefix(key, r.Prefix) {
			continue
		}
		if r.Suffix != "" && !strings.HasSuffix(key, r.Suffix) {
			continue
		}
		if r.Contains != "" && !strings.Contains(key, r.Contains) {
			continue
		}
		return true
	}
	return false
}
