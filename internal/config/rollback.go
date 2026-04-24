package config

import "path/filepath"

// RollbackConfig controls rollback behaviour.
type RollbackConfig struct {
	Enabled   bool   `yaml:"enabled"`
	BackupDir string `yaml:"backup_dir"`
}

// DefaultRollbackConfig returns sensible defaults.
func DefaultRollbackConfig() *RollbackConfig {
	return &RollbackConfig{
		Enabled:   true,
		BackupDir: ".vaultpull/backups",
	}
}

// ApplyRollbackDefaults fills zero-value fields with defaults.
func ApplyRollbackDefaults(c *RollbackConfig) {
	if c == nil {
		return
	}
	d := DefaultRollbackConfig()
	if c.BackupDir == "" {
		c.BackupDir = d.BackupDir
	}
}

// IsEnabled returns true when rollback is enabled.
func (c *RollbackConfig) IsEnabled() bool {
	return c != nil && c.Enabled
}

// BackupPath returns the absolute path for a named backup entry by joining
// BackupDir with the provided name. It returns an empty string if the receiver
// is nil or name is empty.
func (c *RollbackConfig) BackupPath(name string) string {
	if c == nil || name == "" {
		return ""
	}
	return filepath.Join(c.BackupDir, name)
}
