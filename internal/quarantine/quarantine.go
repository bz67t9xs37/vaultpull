package quarantine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a quarantined secret key with its reason.
type Entry struct {
	Path      string    `json:"path"`
	Key       string    `json:"key"`
	Reason    string    `json:"reason"`
	QuarantinedAt time.Time `json:"quarantined_at"`
}

// Quarantine manages keys flagged for review before sync.
type Quarantine struct {
	storeDir string
}

// New returns a Quarantine backed by the given directory.
func New(storeDir string) *Quarantine {
	return &Quarantine{storeDir: storeDir}
}

func (q *Quarantine) filePath(path string) string {
	safe := filepath.Base(path)
	return filepath.Join(q.storeDir, safe+".quarantine.json")
}

// Add records a key as quarantined for a given vault path.
func (q *Quarantine) Add(path, key, reason string) error {
	if err := os.MkdirAll(q.storeDir, 0o700); err != nil {
		return fmt.Errorf("quarantine: mkdir: %w", err)
	}
	entries, _ := q.List(path)
	for _, e := range entries {
		if e.Key == key {
			return nil // already quarantined
		}
	}
	entries = append(entries, Entry{
		Path:          path,
		Key:           key,
		Reason:        reason,
		QuarantinedAt: time.Now().UTC(),
	})
	return q.save(path, entries)
}

// Remove lifts the quarantine on a specific key.
func (q *Quarantine) Remove(path, key string) error {
	entries, err := q.List(path)
	if err != nil {
		return err
	}
	filtered := entries[:0]
	for _, e := range entries {
		if e.Key != key {
			filtered = append(filtered, e)
		}
	}
	return q.save(path, filtered)
}

// List returns all quarantined entries for a vault path.
func (q *Quarantine) List(path string) ([]Entry, error) {
	data, err := os.ReadFile(q.filePath(path))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("quarantine: read: %w", err)
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("quarantine: unmarshal: %w", err)
	}
	return entries, nil
}

// IsQuarantined reports whether a key is quarantined for a vault path.
func (q *Quarantine) IsQuarantined(path, key string) bool {
	entries, _ := q.List(path)
	for _, e := range entries {
		if e.Key == key {
			return true
		}
	}
	return false
}

func (q *Quarantine) save(path string, entries []Entry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("quarantine: marshal: %w", err)
	}
	return os.WriteFile(q.filePath(path), data, 0o600)
}
