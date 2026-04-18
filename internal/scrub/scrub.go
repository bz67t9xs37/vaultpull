package scrub

import (
	"strings"
)

// Scrubber removes or replaces secrets from free-form text output.
type Scrubber struct {
	secrets  []string
	replacement string
}

// New creates a Scrubber that will replace any of the given secret values
// with the replacement string (default "[REDACTED]").
func New(secrets []string, replacement string) *Scrubber {
	if replacement == "" {
		replacement = "[REDACTED]"
	}
	// filter out short/empty values to avoid over-scrubbing
	filtered := make([]string, 0, len(secrets))
	for _, s := range secrets {
		if len(s) >= 4 {
			filtered = append(filtered, s)
		}
	}
	return &Scrubber{secrets: filtered, replacement: replacement}
}

// Line scrubs a single line of text.
func (s *Scrubber) Line(line string) string {
	for _, secret := range s.secrets {
		line = strings.ReplaceAll(line, secret, s.replacement)
	}
	return line
}

// Lines scrubs each line in the slice and returns a new slice.
func (s *Scrubber) Lines(lines []string) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = s.Line(l)
	}
	return out
}

// Map scrubs all values in a map[string]string, leaving keys intact.
func (s *Scrubber) Map(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = s.Line(v)
	}
	return out
}
