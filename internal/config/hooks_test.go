package config

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/hook"
)

func TestToRunner_BothHooks(t *testing.T) {
	h := &HooksConfig{
		PreSync:  "echo pre",
		PostSync: "echo post",
	}
	r := h.ToRunner()
	if !r.Has(hook.PreSync) {
		t.Error("expected PreSync hook to be registered")
	}
	if !r.Has(hook.PostSync) {
		t.Error("expected PostSync hook to be registered")
	}
}

func TestToRunner_NoHooks(t *testing.T) {
	h := &HooksConfig{}
	r := h.ToRunner()
	if r.Has(hook.PreSync) {
		t.Error("expected no PreSync hook")
	}
	if r.Has(hook.PostSync) {
		t.Error("expected no PostSync hook")
	}
}

func TestToRunner_OnlyPreSync(t *testing.T) {
	h := &HooksConfig{PreSync: "make pre"}
	r := h.ToRunner()
	if !r.Has(hook.PreSync) {
		t.Error("expected PreSync hook")
	}
	if r.Has(hook.PostSync) {
		t.Error("expected no PostSync hook")
	}
}

func TestToRunner_OnlyPostSync(t *testing.T) {
	h := &HooksConfig{PostSync: "make post"}
	r := h.ToRunner()
	if r.Has(hook.PreSync) {
		t.Error("expected no PreSync hook")
	}
	if !r.Has(hook.PostSync) {
		t.Error("expected PostSync hook")
	}
}

func TestToRunner_NilConfig(t *testing.T) {
	var h *HooksConfig
	r := h.ToRunner()
	if r == nil {
		t.Fatal("expected non-nil runner from nil config")
	}
}

func TestApplyHookDefaults_NilSafe(t *testing.T) {
	ApplyHookDefaults(nil) // should not panic
}
