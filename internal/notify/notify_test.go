package notify_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vaultpull/internal/diff"
	"github.com/vaultpull/internal/notify"
)

func TestSend_Disabled(t *testing.T) {
	n, err := notify.New(notify.Config{Enabled: false, Channel: notify.ChannelStdout})
	if err != nil {
		t.Fatal(err)
	}
	// Should not error even with changes
	changes := []diff.Change{{Key: "FOO", Type: diff.Added, New: "bar"}}
	if err := n.Send(".env", changes); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSend_OnlyDiff_NoChanges(t *testing.T) {
	buf := &bytes.Buffer{}
	n := notifierWithWriter(t, buf)
	if err := n.Send(".env", nil); err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got: %s", buf.String())
	}
}

func TestSend_WithChanges_WritesOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	n := notifierWithWriter(t, buf)
	changes := []diff.Change{
		{Key: "A", Type: diff.Added, New: "1"},
		{Key: "B", Type: diff.Removed, Old: "2"},
	}
	if err := n.Send(".env", changes); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, ".env") {
		t.Errorf("expected path in output, got: %s", out)
	}
	if !strings.Contains(out, "+1") {
		t.Errorf("expected added count, got: %s", out)
	}
}

func TestSend_FileChannel_WritesToFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "notify.log")
	n, err := notify.New(notify.Config{
		Enabled:  true,
		Channel:  notify.ChannelFile,
		FilePath: path,
		OnlyDiff: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := n.Send(".env", nil); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "vaultpull") {
		t.Errorf("expected log entry in file, got: %s", string(data))
	}
}

func TestNew_FileChannel_MissingPath(t *testing.T) {
	_, err := notify.New(notify.Config{Enabled: true, Channel: notify.ChannelFile})
	if err == nil {
		t.Error("expected error for missing file path")
	}
}

// notifierWithWriter returns a notifier that writes to buf via stdout override trick.
// We use a helper that builds a notifier with a custom writer via a small shim.
func notifierWithWriter(t *testing.T, buf *bytes.Buffer) *notify.Notifier {
	t.Helper()
	n, err := notify.NewWithWriter(notify.Config{Enabled: true, OnlyDiff: true}, buf)
	if err != nil {
		t.Fatal(err)
	}
	return n
}
