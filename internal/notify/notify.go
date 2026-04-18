package notify

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/vaultpull/internal/diff"
)

// Channel represents a notification destination.
type Channel string

const (
	ChannelStdout Channel = "stdout"
	ChannelStderr Channel = "stderr"
	ChannelFile   Channel = "file"
)

// Config holds notification settings.
type Config struct {
	Enabled  bool
	Channel  Channel
	FilePath string // used when Channel == ChannelFile
	OnlyDiff bool   // only notify when there are changes
}

// Notifier sends sync notifications.
type Notifier struct {
	cfg Config
	out io.Writer
}

// New creates a Notifier. If cfg.Channel is ChannelFile, the file is opened/created.
func New(cfg Config) (*Notifier, error) {
	var w io.Writer
	switch cfg.Channel {
	case ChannelStderr:
		w = os.Stderr
	case ChannelFile:
		if cfg.FilePath == "" {
			return nil, fmt.Errorf("notify: file channel requires a file path")
		}
		f, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, fmt.Errorf("notify: open file: %w", err)
		}
		w = f
	default:
		w = os.Stdout
	}
	return &Notifier{cfg: cfg, out: w}, nil
}

// Send writes a notification for the given path and diff changes.
func (n *Notifier) Send(path string, changes []diff.Change) error {
	if !n.cfg.Enabled {
		return nil
	}
	if n.cfg.OnlyDiff && len(changes) == 0 {
		return nil
	}
	summary := diff.Summary(changes)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[vaultpull] sync: %s — ", path))
	sb.WriteString(fmt.Sprintf("+%d -%d ~%d unchanged:%d",
		summary.Added, summary.Removed, summary.Modified, summary.Unchanged))
	sb.WriteString("\n")
	_, err := fmt.Fprint(n.out, sb.String())
	return err
}
