package config

import "github.com/yourusername/vaultpull/internal/hook"

// HookConfig holds lifecycle hook commands for a sync run.
type HookConfig struct {
	PreSync  string `yaml:"pre_sync"`
	PostSync string `yaml:"post_sync"`
}

// ToRunner converts a HookConfig into a hook.Runner.
func (h HookConfig) ToRunner() *hook.Runner {
	hooks := map[string]hook.Hook{}
	if h.PreSync != "" {
		hooks["pre-sync"] = hook.Hook{Command: h.PreSync}
	}
	if h.PostSync != "" {
		hooks["post-sync"] = hook.Hook{Command: h.PostSync}
	}
	return hook.New(hooks)
}
