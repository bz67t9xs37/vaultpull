package config

import "github.com/yourusername/vaultpull/internal/namespace"

// NamespaceConfig holds namespace resolution settings for a sync target.
type NamespaceConfig struct {
	// Prefix is a path prefix prepended to all secret names when resolving
	// Vault paths. For example, "myapp/production".
	Prefix string `yaml:"prefix"`

	// Mount is the KV secrets engine mount path. Defaults to "secret".
	Mount string `yaml:"mount"`
}

// ToResolver converts a NamespaceConfig into a namespace.Resolver.
// If Mount is empty, the default KV mount "secret" is used.
func (n NamespaceConfig) ToResolver() *namespace.Resolver {
	mount := n.Mount
	if mount == "" {
		mount = DefaultMount
	}
	return namespace.New(n.Prefix, mount)
}

// DefaultNamespaceConfig returns a NamespaceConfig with sensible defaults.
func DefaultNamespaceConfig() NamespaceConfig {
	return NamespaceConfig{
		Mount: DefaultMount,
	}
}
