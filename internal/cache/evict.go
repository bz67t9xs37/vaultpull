package cache

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// EvictExpired removes all cache entries older than the given TTL.
// It scans the cache directory and deletes stale JSON files.
func (c *Cache) EvictExpired(ttl time.Duration) error {
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	cutoff := time.Now().Add(-ttl)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			_ = os.Remove(filepath.Join(c.dir, entry.Name()))
		}
	}
	return nil
}

// EvictAll removes all cache entries from the cache directory.
func (c *Cache) EvictAll() error {
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			_ = os.Remove(filepath.Join(c.dir, entry.Name()))
		}
	}
	return nil
}
