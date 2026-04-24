package watermark

import (
	"fmt"
	"sync"
	"time"
)

// Entry records the high-water mark for a secret path: the last known
// version hash and when it was recorded.
type Entry struct {
	Path      string
	Hash      string
	RecordedAt time.Time
}

// Watermark tracks the last-seen content hash for each secret path so
// that downstream consumers can detect when a value has changed since
// the previous sync.
type Watermark struct {
	mu      sync.RWMutex
	marks   map[string]Entry
	clock   func() time.Time
}

// New returns a ready-to-use Watermark.
func New() *Watermark {
	return newWithClock(time.Now)
}

func newWithClock(clock func() time.Time) *Watermark {
	return &Watermark{
		marks: make(map[string]Entry),
		clock: clock,
	}
}

// Set records (or updates) the hash for path.
func (w *Watermark) Set(path, hash string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.marks[path] = Entry{
		Path:       path,
		Hash:       hash,
		RecordedAt: w.clock(),
	}
}

// Get returns the stored entry for path and whether it exists.
func (w *Watermark) Get(path string) (Entry, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	e, ok := w.marks[path]
	return e, ok
}

// Changed reports whether hash differs from the stored mark for path.
// Returns true when the path has never been seen before.
func (w *Watermark) Changed(path, hash string) bool {
	e, ok := w.Get(path)
	if !ok {
		return true
	}
	return e.Hash != hash
}

// Delete removes the watermark for path.
func (w *Watermark) Delete(path string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.marks, path)
}

// All returns a snapshot of every tracked entry.
func (w *Watermark) All() []Entry {
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := make([]Entry, 0, len(w.marks))
	for _, e := range w.marks {
		out = append(out, e)
	}
	return out
}

// Summary returns a human-readable line describing the current state.
func (w *Watermark) Summary() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return fmt.Sprintf("watermark: %d path(s) tracked", len(w.marks))
}
