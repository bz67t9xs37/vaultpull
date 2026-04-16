package config

// NamespaceConfig controls how Vault secret paths are namespaced.
type NamespaceConfig struct {
	Prefix     string `yaml:"prefix"`
	StripMount bool   `yaml:"strip_mount"`
}

// DefaultNamespaceConfig returns a NamespaceConfig with default values.
func DefaultNamespaceConfig() *NamespaceConfig {
	return &NamespaceConfig{
		Prefix:     "",
		StripMount: false,
	}
}

// ApplyNamespaceDefaults ensures Namespace is initialized on the config.
func ApplyNamespaceDefaults(cfg *Config) {
	if cfg.Namespace == nil {
		cfg.Namespace = DefaultNamespaceConfig()
	}
}
