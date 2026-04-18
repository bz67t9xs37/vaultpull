package notify

import "io"

// NewWithWriter creates a Notifier that writes to the provided writer.
// Useful for testing.
func NewWithWriter(cfg Config, w io.Writer) (*Notifier, error) {
	return &Notifier{cfg: cfg, out: w}, nil
}
