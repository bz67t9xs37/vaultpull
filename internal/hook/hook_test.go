package hook_test

import (
	"runtime"
	"testing"

	"github.com/yourusername/vaultpull/internal/hook"
)

func TestRun_NoHook(t *testing.T) {
	r := hook.New(map[string]hook.Hook{})
	if err := r.Run("pre-sync"); err != nil {
		t.Fatalf("expected no error for missing hook, got %v", err)
	}
}

func TestRun_EmptyCommand(t *testing.T) {
	r := hook.New(map[string]hook.Hook{
		"post-sync": {Command: "   "},
	})
	if err := r.Run("post-sync"); err != nil {
		t.Fatalf("expected no error for empty command, got %v", err)
	}
}

func TestRun_SuccessfulCommand(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell command test on windows")
	}
	r := hook.New(map[string]hook.Hook{
		"post-sync": {Command: "echo hello"},
	})
	if err := r.Run("post-sync"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_FailingCommand(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell command test on windows")
	}
	r := hook.New(map[string]hook.Hook{
		"pre-sync": {Command: "false"},
	})
	err := r.Run("pre-sync")
	if err == nil {
		t.Fatal("expected error for failing command, got nil")
	}
}

func TestHas_RegisteredHook(t *testing.T) {
	r := hook.New(map[string]hook.Hook{
		"pre-sync": {Command: "echo hi"},
	})
	if !r.Has("pre-sync") {
		t.Error("expected Has to return true for registered hook")
	}
	if r.Has("post-sync") {
		t.Error("expected Has to return false for unregistered hook")
	}
}

func TestHas_EmptyCommandNotRegistered(t *testing.T) {
	r := hook.New(map[string]hook.Hook{
		"pre-sync": {Command: ""},
	})
	if r.Has("pre-sync") {
		t.Error("expected Has to return false for empty command hook")
	}
}
