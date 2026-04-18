package pin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry records a pinned secret version for a given path.
type Entry struct {
	Path      string    `json:"path"`
	Version   int       `json:"version"`
	PinnedAt  time.Time `json:"pinned_at"`
}

// Store manages pinned secret versions on disk.
type Store struct {
	filePath string
	entries  map[string]Entry
}

// New creates a Store backed by the given file path.
func New(filePath string) (*Store, error) {
	s := &Store{filePath: filePath, entries: make(map[string]Entry)}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

// Pin records a version for the given vault path.
func (s *Store) Pin(path string, version int) error {
	s.entries[path] = Entry{Path: path, Version: version, PinnedAt: time.Now().UTC()}
	return s.save()
}

// Get returns the pinned entry for a path, or nil if not pinned.
func (s *Store) Get(path string) *Entry {
	e, ok := s.entries[path]
	if !ok {
		return nil
	}
	return &e
}

// Unpin removes the pin for the given path.
func (s *Store) Unpin(path string) error {
	delete(s.entries, path)
	return s.save()
}

// All returns all pinned entries.
func (s *Store) All() []Entry {
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	return out
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.filePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("pin: read %s: %w", s.filePath, err)
	}
	return json.Unmarshal(data, &s.entries)
}

func (s *Store) save() error {
	if err := os.MkdirAll(filepath.Dir(s.filePath), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0o600)
}
