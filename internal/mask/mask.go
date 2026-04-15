package mask

import "strings"

// DefaultSensitiveKeys contains common secret key patterns to mask in output.
var DefaultSensitiveKeys = []string{
	"password",
	"secret",
	"token",
	"key",
	"api_key",
	"apikey",
	"auth",
	"credential",
	"private",
	"passphrase",
}

// Masker redacts sensitive secret values before display.
type Masker struct {
	sensitiveKeys []string
	maskChar      string
}

// New returns a Masker with the given sensitive key patterns.
// If sensitiveKeys is empty, DefaultSensitiveKeys is used.
func New(sensitiveKeys []string, maskChar string) *Masker {
	if len(sensitiveKeys) == 0 {
		sensitiveKeys = DefaultSensitiveKeys
	}
	if maskChar == "" {
		maskChar = "********"
	}
	return &Masker{
		sensitiveKeys: sensitiveKeys,
		maskChar:      maskChar,
	}
}

// IsSensitive reports whether the given key matches any sensitive pattern.
func (m *Masker) IsSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, pattern := range m.sensitiveKeys {
		if strings.Contains(lower, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

// MaskValue returns the masked string if the key is sensitive,
// otherwise returns the original value unchanged.
func (m *Masker) MaskValue(key, value string) string {
	if m.IsSensitive(key) {
		return m.maskChar
	}
	return value
}

// MaskMap returns a copy of secrets with sensitive values redacted.
func (m *Masker) MaskMap(secrets map[string]string) map[string]string {
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = m.MaskValue(k, v)
	}
	return result
}
