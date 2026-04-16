package redact

import (
	"regexp"
	"strings"
)

// Rule defines a redaction rule for output sanitisation.
type Rule struct {
	Pattern     *regexp.Regexp
	Replacement string
}

// Redactor applies a set of redaction rules to strings.
type Redactor struct {
	rules  []Rule
	values map[string]struct{}
}

// New creates a Redactor with the provided literal secret values to redact.
// Values shorter than 4 characters are ignored to avoid over-redaction.
func New(secrets []string) *Redactor {
	r := &Redactor{
		values: make(map[string]struct{}),
	}
	for _, s := range secrets {
		if len(s) >= 4 {
			r.values[s] = struct{}{}
		}
	}
	return r
}

// AddPattern registers a regex-based redaction rule.
func (r *Redactor) AddPattern(pattern, replacement string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	r.rules = append(r.rules, Rule{Pattern: re, Replacement: replacement})
	return nil
}

// Redact replaces all known secret values and pattern matches in s with "[REDACTED]".
func (r *Redactor) Redact(s string) string {
	result := s
	for v := range r.values {
		result = strings.ReplaceAll(result, v, "[REDACTED]")
	}
	for _, rule := range r.rules {
		result = rule.Pattern.ReplaceAllString(result, rule.Replacement)
	}
	return result
}

// RedactMap returns a copy of m with all values redacted.
func (r *Redactor) RedactMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = r.Redact(v)
	}
	return out
}
