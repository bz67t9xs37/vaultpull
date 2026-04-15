package transform

import (
	"fmt"
	"strings"
)

// Rule defines a transformation to apply to secret values or keys.
type Rule struct {
	Type  string // "prefix", "suffix", "uppercase", "lowercase", "replace"
	Key   string // optional: apply only to this key; empty means all keys
	From  string // used by "replace"
	To    string // used by "prefix", "suffix", "replace"
}

// Transformer applies a set of rules to a map of secrets.
type Transformer struct {
	rules []Rule
}

// New creates a new Transformer with the given rules.
func New(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// Apply applies all transformation rules to the provided secrets map,
// returning a new map with transformed values. Keys are never mutated.
func (t *Transformer) Apply(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	for _, rule := range t.rules {
		for k, v := range out {
			if rule.Key != "" && rule.Key != k {
				continue
			}
			result, err := applyRule(rule, v)
			if err != nil {
				return nil, fmt.Errorf("transform rule %q on key %q: %w", rule.Type, k, err)
			}
			out[k] = result
		}
	}
	return out, nil
}

func applyRule(rule Rule, value string) (string, error) {
	switch rule.Type {
	case "prefix":
		return rule.To + value, nil
	case "suffix":
		return value + rule.To, nil
	case "uppercase":
		return strings.ToUpper(value), nil
	case "lowercase":
		return strings.ToLower(value), nil
	case "replace":
		return strings.ReplaceAll(value, rule.From, rule.To), nil
	default:
		return "", fmt.Errorf("unknown rule type %q", rule.Type)
	}
}
