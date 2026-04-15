package hook

import (
	"fmt"
	"os/exec"
	"strings"
)

// Hook represents a shell command to run at a lifecycle point.
type Hook struct {
	Command string
}

// Runner executes lifecycle hooks.
type Runner struct {
	hooks map[string]Hook
}

// New creates a Runner with the given hooks map.
// Keys are lifecycle event names (e.g. "pre-sync", "post-sync").
func New(hooks map[string]Hook) *Runner {
	return &Runner{hooks: hooks}
}

// Run executes the hook registered for the given event.
// If no hook is registered, Run is a no-op and returns nil.
func (r *Runner) Run(event string) error {
	h, ok := r.hooks[event]
	if !ok || strings.TrimSpace(h.Command) == "" {
		return nil
	}

	parts := strings.Fields(h.Command)
	if len(parts) == 0 {
		return nil
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("hook %q failed: %w\noutput: %s", event, err, strings.TrimSpace(string(out)))
	}
	return nil
}

// Has reports whether a hook is registered for the given event.
func (r *Runner) Has(event string) bool {
	h, ok := r.hooks[event]
	return ok && strings.TrimSpace(h.Command) != ""
}
