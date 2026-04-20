package config

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
