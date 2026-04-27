package label

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Labels holds arbitrary key-value metadata attached to a secret path.
type Labels map[string]string

// Store persists labels for secret paths to a JSON file.
type Store struct {
	mu      sync.RWMutex
	path    string
	entries map[string]Labels
}

// New creates a Store backed by the given file path.
func New(storePath string) (*Store, error) {
	s := &Store{
		path:    storePath,
		entries: make(map[string]Labels),
	}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

// Set assigns labels to a secret path, merging with any existing labels.
func (s *Store) Set(secretPath string, labels Labels) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.entries[secretPath]; !ok {
		s.entries[secretPath] = make(Labels)
	}
	for k, v := range labels {
		s.entries[secretPath][k] = v
	}
	return s.save()
}

// Get returns the labels for a secret path, or nil if none exist.
func (s *Store) Get(secretPath string) Labels {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.entries[secretPath]
}

// Remove deletes a specific label key from a secret path.
func (s *Store) Remove(secretPath, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if lbls, ok := s.entries[secretPath]; ok {
		delete(lbls, key)
		if len(lbls) == 0 {
			delete(s.entries, secretPath)
		}
	}
	return s.save()
}

// All returns a copy of all labeled paths and their labels.
func (s *Store) All() map[string]Labels {
	s.mu.RLock()
	defer s.mu.RUnlock()

	copy := make(map[string]Labels, len(s.entries))
	for k, v := range s.entries {
		copy[k] = v
	}
	return copy
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("label: read store: %w", err)
	}
	return json.Unmarshal(data, &s.entries)
}

func (s *Store) save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return fmt.Errorf("label: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("label: marshal: %w", err)
	}
	return os.WriteFile(s.path, data, 0o600)
}
