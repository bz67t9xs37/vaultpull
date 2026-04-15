package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a cached snapshot of secrets for a given target.
type Entry struct {
	Path      string            `json:"path"`
	Secrets   map[string]string `json:"secrets"`
	FetchedAt time.Time         `json:"fetched_at"`
	Checksum  string            `json:"checksum"`
}

// Cache manages on-disk secret snapshots to avoid redundant Vault reads.
type Cache struct {
	dir string
}

// New returns a Cache that stores entries under dir.
func New(dir string) (*Cache, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("cache: create dir: %w", err)
	}
	return &Cache{dir: dir}, nil
}

// Get returns the cached Entry for vaultPath, or (nil, nil) if absent.
func (c *Cache) Get(vaultPath string) (*Entry, error) {
	p := c.entryPath(vaultPath)
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cache: read: %w", err)
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("cache: unmarshal: %w", err)
	}
	return &e, nil
}

// Set writes secrets for vaultPath into the cache.
func (c *Cache) Set(vaultPath string, secrets map[string]string) error {
	e := Entry{
		Path:      vaultPath,
		Secrets:   secrets,
		FetchedAt: time.Now().UTC(),
		Checksum:  checksum(secrets),
	}
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Errorf("cache: marshal: %w", err)
	}
	return os.WriteFile(c.entryPath(vaultPath), data, 0o600)
}

// Invalidate removes the cached entry for vaultPath.
func (c *Cache) Invalidate(vaultPath string) error {
	err := os.Remove(c.entryPath(vaultPath))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (c *Cache) entryPath(vaultPath string) string {
	h := sha256.Sum256([]byte(vaultPath))
	name := hex.EncodeToString(h[:]) + ".json"
	return filepath.Join(c.dir, name)
}

func checksum(secrets map[string]string) string {
	b, _ := json.Marshal(secrets)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
