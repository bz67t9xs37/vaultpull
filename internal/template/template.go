package template

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Renderer renders .env templates by substituting ${VAR} placeholders
// with values from a provided secrets map.
type Renderer struct {
	placeholder *regexp.Regexp
}

// New returns a new Renderer.
func New() *Renderer {
	return &Renderer{
		placeholder: regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`),
	}
}

// RenderFile reads a template file and substitutes all ${VAR} placeholders
// with values from secrets. Returns the rendered content as a string.
// Unknown placeholders are left unchanged.
func (r *Renderer) RenderFile(templatePath string, secrets map[string]string) (string, error) {
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("template: read file %q: %w", templatePath, err)
	}
	return r.Render(string(data), secrets), nil
}

// Render substitutes all ${VAR} placeholders in src with values from secrets.
// Placeholders without a matching key are left unchanged.
func (r *Renderer) Render(src string, secrets map[string]string) string {
	return r.placeholder.ReplaceAllStringFunc(src, func(match string) string {
		key := match[2 : len(match)-1] // strip ${ and }
		if val, ok := secrets[key]; ok {
			return val
		}
		return match
	})
}

// ListPlaceholders returns all unique placeholder names found in src.
func (r *Renderer) ListPlaceholders(src string) []string {
	matches := r.placeholder.FindAllStringSubmatch(src, -1)
	seen := make(map[string]struct{})
	var keys []string
	for _, m := range matches {
		if _, ok := seen[m[1]]; !ok {
			seen[m[1]] = struct{}{}
			keys = append(keys, m[1])
		}
	}
	return keys
}

// MissingKeys returns placeholder names present in src that have no
// corresponding entry in secrets.
func (r *Renderer) MissingKeys(src string, secrets map[string]string) []string {
	var missing []string
	for _, key := range r.ListPlaceholders(src) {
		if _, ok := secrets[key]; !ok {
			missing = append(missing, key)
		}
	}
	return missing
}

// ValidateTemplate returns an error if the template at templatePath contains
// placeholders that are absent from secrets.
func (r *Renderer) ValidateTemplate(templatePath string, secrets map[string]string) error {
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("template: read file %q: %w", templatePath, err)
	}
	if missing := r.MissingKeys(string(data), secrets); len(missing) > 0 {
		return fmt.Errorf("template: unresolved placeholders: %s", strings.Join(missing, ", "))
	}
	return nil
}
