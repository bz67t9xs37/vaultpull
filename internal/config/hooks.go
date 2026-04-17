package config

import "github.com/yourusername/vaultpull/internal/hook"

// HooksConfig holds pre/post sync hook commands.
type HooksConfig struct {
	PreSync  string `yaml:"pre_sync"`
	PostSync string `yaml:"post_sync"`
}

// DefaultHooksConfig returns a HooksConfig with zero values.
func DefaultHooksConfig() *HooksConfig {
	return &HooksConfig{}
}

// ApplyHookDefaults is a no-op for now; hooks have no required defaults.
func ApplyHookDefaults(h *HooksConfig) {
	if h == nil {
		return
	}
}

// ToRunner converts the HooksConfig into a hook.Runner.
func (h *HooksConfig) ToRunner() *hook.Runner {
	if h == nil {
		return hook.Newtr := hook.New()
	if h.PreSync != "" {
		r.Register(hook.PreSync, h.PreSync)
	}
	if h.PostSync != "" {
		r.Register(hook.PostSync, h.PostSync)
	}
	return r
}
