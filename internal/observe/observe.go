package observe

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Event represents a single observation event during a sync lifecycle.
type Event struct {
	Timestamp time.Time
	Level     string
	Path      string
	Message   string
}

// Observer records and emits lifecycle events for a sync operation.
type Observer struct {
	w      io.Writer
	events []Event
	clock  func() time.Time
}

// New returns an Observer that writes to stderr by default.
func New() *Observer {
	return newWithClock(os.Stderr, time.Now)
}

func newWithClock(w io.Writer, clock func() time.Time) *Observer {
	return &Observer{w: w, clock: clock}
}

// Emit records and writes an event at the given level.
func (o *Observer) Emit(level, path, message string) {
	e := Event{
		Timestamp: o.clock(),
		Level:     level,
		Path:      path,
		Message:   message,
	}
	o.events = append(o.events, e)
	fmt.Fprintf(o.w, "[%s] %s %s: %s\n", e.Timestamp.Format(time.RFC3339), e.Level, e.Path, e.Message)
}

// Info emits an INFO-level event.
func (o *Observer) Info(path, message string) {
	o.Emit("INFO", path, message)
}

// Warn emits a WARN-level event.
func (o *Observer) Warn(path, message string) {
	o.Emit("WARN", path, message)
}

// Error emits an ERROR-level event.
func (o *Observer) Error(path, message string) {
	o.Emit("ERROR", path, message)
}

// Events returns all recorded events.
func (o *Observer) Events() []Event {
	out := make([]Event, len(o.events))
	copy(out, o.events)
	return out
}

// CountByLevel returns the number of events at a given level.
func (o *Observer) CountByLevel(level string) int {
	count := 0
	for _, e := range o.events {
		if e.Level == level {
			count++
		}
	}
	return count
}
