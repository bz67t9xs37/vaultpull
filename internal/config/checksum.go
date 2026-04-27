package config

// ChecksumConfig controls whether secret checksums are computed and stored
// to detect out-of-band changes between syncs.
type ChecksumConfig struct {
	Enabled  bool   `yaml:"enabled"`
	StoreDir string `yaml:"store_dir"`
}

// DefaultChecksumConfig returns a ChecksumConfig with sensible defaults.
func DefaultChecksumConfig() ChecksumConfig {
	return ChecksumConfig{
		Enabled:  false,
		StoreDir: ".vaultpull/checksums",
	}
}

// ApplyChecksumDefaults fills zero-value fields with defaults.
// It is nil-safe and returns early when cfg is nil.
func ApplyChecksumDefaults(cfg *ChecksumConfig) {
	if cfg == nil {
		return
	}
	def := DefaultChecksumConfig()
	if cfg.StoreDir == "" {
		cfg.StoreDir = def.StoreDir
	}
}

// IsEnabled returns true when checksum tracking is active.
func (c *ChecksumConfig) IsEnabled() bool {
	return c != nil && c.Enabled
}
