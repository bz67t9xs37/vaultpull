package tag

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Tag represents a named marker attached to a vault path at a point in time.
type Tag struct {
	Name    string            `json:"name"`
	Path    string            `json:"path"`
	Meta    map[string]string `json:"meta,omitempty"`
}

// Manager handles tagging of vault secret paths.
type Manager struct {
	storeDir string
}

// New creates a new tag Manager backed by storeDir.
func New(storeDir string) (*Manager, error) {
	if err := os.MkdirAll(storeDir, 0o700); err != nil {
		return nil, fmt.Errorf("tag: create store dir: %w", err)
	}
	return &Manager{storeDir: storeDir}, nil
}

func (m *Manager) storePath(name string) string {
	return filepath.Join(m.storeDir, name+".json")
}

// Set creates or replaces a tag with the given name for the given vault path.
func (m *Manager) Set(name, path string, meta map[string]string) error {
	if name == "" {
		return fmt.Errorf("tag: name must not be empty")
	}
	if path == "" {
		return fmt.Errorf("tag: path must not be empty")
	}
	t := Tag{Name: name, Path: path, Meta: meta}
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Errorf("tag: marshal: %w", err)
	}
	if err := os.WriteFile(m.storePath(name), data, 0o600); err != nil {
		return fmt.Errorf("tag: write: %w", err)
	}
	return nil
}

// Get retrieves the tag with the given name. Returns nil if not found.
func (m *Manager) Get(name string) (*Tag, error) {
	data, err := os.ReadFile(m.storePath(name))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("tag: read: %w", err)
	}
	var t Tag
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("tag: unmarshal: %w", err)
	}
	return &t, nil
}

// Delete removes a tag by name. No-op if the tag does not exist.
func (m *Manager) Delete(name string) error {
	err := os.Remove(m.storePath(name))
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("tag: delete: %w", err)
	}
	return nil
}

// All returns all stored tags sorted by name.
func (m *Manager) All() ([]Tag, error) {
	entries, err := os.ReadDir(m.storeDir)
	if err != nil {
		return nil, fmt.Errorf("tag: list: %w", err)
	}
	var tags []Tag
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		name := e.Name()[:len(e.Name())-5]
		t, err := m.Get(name)
		if err != nil || t == nil {
			continue
		}
		tags = append(tags, *t)
	}
	sort.Slice(tags, func(i, j int) bool { return tags[i].Name < tags[j].Name })
	return tags, nil
}
