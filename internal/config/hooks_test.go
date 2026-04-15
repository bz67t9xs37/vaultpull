package config_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/config"
)

func TestToRunner_BothHooks(t *testing.T) {
	hc := config.HookConfig{
		PreSync:  "echo pre",
		PostSync: "echo post",
	}
	r := hc.ToRunner()
	if !r.Has("pre-sync") {
		t.Error("expected pre-sync hook to be registered")
	}
	if !r.Has("post-sync") {
		t.Error("expected post-sync hook to be registered")
	}
}

func TestToRunner_NoHooks(t *testing.T) {
	hc := config.HookConfig{}
	r := hc.ToRunner()
	if r.Has("pre-sync") {
		t.Error("expected pre-sync hook to be absent")
	}
	if r.Has("post-sync") {
		t.Error("expected post-sync hook to be absent")
	}
}

func TestToRunner_OnlyPreSync(t *testing.T) {
	hc := config.HookConfig{
		PreSync: "echo before",
	}
	r := hc.ToRunner()
	if !r.Has("pre-sync") {
		t.Error("expected pre-sync hook to be registered")
	}
	if r.Has("post-sync") {
		t.Error("expected post-sync hook to be absent")
	}
}

func TestToRunner_OnlyPostSync(t *testing.T) {
	hc := config.HookConfig{
		PostSync: "echo after",
	}
	r := hc.ToRunner()
	if r.Has("pre-sync") {
		t.Error("expected pre-sync hook to be absent")
	}
	if !r.Has("post-sync") {
		t.Error("expected post-sync hook to be registered")
	}
}
