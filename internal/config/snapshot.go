package config

import "path/filepath"

// SnapshotConfig controls where and how secret snapshots are persisted.
type SnapshotConfig struct {
	// Enabled toggles snapshot capture on each sync.
	Enabled bool `yaml:"enabled"`

	// Dir is the directory where snapshot files are stored.
	// Defaults to ".vaultpull/snapshots" relative to the working directory.
	Dir string `yaml:"dir"`

	// MaxPerPath is the maximum number of snapshots to retain per vault path.
	// Zero means unlimited.
	MaxPerPath int `yaml:"max_per_path"`
}

// DefaultSnapshotConfig returns a SnapshotConfig with sensible defaults.
func DefaultSnapshotConfig() SnapshotConfig {
	return SnapshotConfig{
		Enabled:    true,
		Dir:        filepath.Join(".vaultpull", "snapshots"),
		MaxPerPath: 10,
	}
}

// ApplySnapshotDefaults fills in zero-value fields with defaults.
func ApplySnapshotDefaults(c *SnapshotConfig) {
	d := DefaultSnapshotConfig()
	if c.Dir == "" {
		c.Dir = d.Dir
	}
	if c.MaxPerPath == 0 {
		c.MaxPerPath = d.MaxPerPath
	}
}
