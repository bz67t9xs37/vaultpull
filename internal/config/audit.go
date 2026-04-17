package config

// AuditConfig holds configuration for audit logging.
type AuditConfig struct {
	Enabled  bool   `yaml:"enabled"`
	LogPath  string `yaml:"log_path"`
	MaxLines int    `yaml:"max_lines"`
}

// DefaultAuditConfig returns sensible defaults for audit logging.
func DefaultAuditConfig() *AuditConfig {
	return &AuditConfig{
		Enabled:  true,
		LogPath:  ".vaultpull_audit.log",
		MaxLines: 1000,
	}
}

// ApplyAuditDefaults fills zero-value fields with defaults.
func ApplyAuditDefaults(c *AuditConfig) {
	if c == nil {
		return
	}
	def := DefaultAuditConfig()
	if c.LogPath == "" {
		c.LogPath = def.LogPath
	}
	if c.MaxLines == 0 {
		c.MaxLines = def.MaxLines
	}
}

// IsEnabled returns true when audit logging is active.
func (c *AuditConfig) IsEnabled() bool {
	return c != nil && c.Enabled
}
