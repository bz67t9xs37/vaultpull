package hook_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/yourusername/vaultpull/internal/hook"
)

// TestHook_WritesToFile verifies that a hook command can produce observable
// side effects (writing a file), confirming real execution.
func TestHook_WritesToFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping integration test on windows")
	}

	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, "hook_ran")

	r := hook.New(map[string]hook.Hook{
		"post-sync": {Command: "touch " + outFile},
	})

	if err := r.Run("post-sync"); err != nil {
		t.Fatalf("hook run failed: %v", err)
	}

	if _, err := os.Stat(outFile); os.IsNotExist(err) {
		t.Error("expected hook to create file, but it does not exist")
	}
}

// TestHook_PreAndPostSyncSequence simulates a full pre/post sync lifecycle.
func TestHook_PreAndPostSyncSequence(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping integration test on windows")
	}

	tmpDir := t.TempDir()
	preFile := filepath.Join(tmpDir, "pre_ran")
	postFile := filepath.Join(tmpDir, "post_ran")

	r := hook.New(map[string]hook.Hook{
		"pre-sync":  {Command: "touch " + preFile},
		"post-sync": {Command: "touch " + postFile},
	})

	for _, event := range []string{"pre-sync", "post-sync"} {
		if err := r.Run(event); err != nil {
			t.Fatalf("hook %q failed: %v", event, err)
		}
	}

	for _, f := range []string{preFile, postFile} {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist after hook run", f)
		}
	}
}
