package lineage

import (
	"fmt"
	"sort"
	"time"
)

// Entry records a single propagation event for a secret key.
type Entry struct {
	Key       string    `json:"key"`
	Source    string    `json:"source"`
	Target    string    `json:"target"`
	Operation string    `json:"operation"` // "sync", "promote", "rotate"
	Timestamp time.Time `json:"timestamp"`
}

// Tracker maintains an in-memory lineage log for secret keys.
type Tracker struct {
	entries []Entry
	clock   func() time.Time
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{clock: time.Now}
}

// newWithClock returns a Tracker with a custom clock (for testing).
func newWithClock(clock func() time.Time) *Tracker {
	return &Tracker{clock: clock}
}

// Record appends a lineage entry for the given key.
func (t *Tracker) Record(key, source, target, operation string) {
	ts := t.clock()
	t.entries = append(t.entries, Entry{
		Key:       key,
		Source:    source,
		Target:    target,
		Operation: operation,
		Timestamp: ts,
	})
}

// History returns all recorded entries for a given key, ordered by time.
func (t *Tracker) History(key string) []Entry {
	var result []Entry
	for _, e := range t.entries {
		if e.Key == key {
			result = append(result, e)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})
	return result
}

// All returns every recorded entry ordered by timestamp.
func (t *Tracker) All() []Entry {
	out := make([]Entry, len(t.entries))
	copy(out, t.entries)
	sort.Slice(out, func(i, j int) bool {
		return out[i].Timestamp.Before(out[j].Timestamp)
	})
	return out
}

// Summary returns a human-readable summary of lineage for a key.
func (t *Tracker) Summary(key string) string {
	h := t.History(key)
	if len(h) == 0 {
		return fmt.Sprintf("no lineage recorded for %q", key)
	}
	return fmt.Sprintf("%s: %d event(s), last %s from %s → %s",
		key, len(h), h[len(h)-1].Operation, h[len(h)-1].Source, h[len(h)-1].Target)
}
