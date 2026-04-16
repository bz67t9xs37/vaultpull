package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a point-in-time capture of secrets for a given target path.
type Entry struct {
	Path      string            `json:"path"`
	Secrets   map[string]string `json:"secrets"`
	CreatedAt time.Time         `json:"created_at"`
}

// Store manages snapshot files on disk.
type Store struct {
	dir string
}

// New creates a new Store rooted at dir.
func New(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("snapshot: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Save writes a snapshot entry to disk.
func (s *Store) Save(e Entry) error {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now().UTC()
	}
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	filename := fmt.Sprintf("%d_%s.json", e.CreatedAt.UnixNano(), sanitize(e.Path))
	return os.WriteFile(filepath.Join(s.dir, filename), data, 0600)
}

// Latest returns the most recent snapshot for the given path, or nil if none exists.
func (s *Store) Latest(path string) (*Entry, error) {
	entries, err := s.List(path)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, nil
	}
	return &entries[len(entries)-1], nil
}

// List returns all snapshots for the given path ordered by creation time (oldest first).
func (s *Store) List(path string) ([]Entry, error) {
	glob := filepath.Join(s.dir, fmt.Sprintf("*_%s.json", sanitize(path)))
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, fmt.Errorf("snapshot: glob: %w", err)
	}
	var out []Entry
	for _, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			return nil, fmt.Errorf("snapshot: read %s: %w", m, err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("snapshot: unmarshal %s: %w", m, err)
		}
		out = append(out, e)
	}
	return out, nil
}

// sanitize converts a vault path into a safe filename segment.
func sanitize(path string) string {
	out := make([]byte, len(path))
	for i := 0; i < len(path); i++ {
		c := path[i]
		if c == '/' || c == '\\' || c == ':' {
			out[i] = '_'
		} else {
			out[i] = c
		}
	}
	return string(out)
}
