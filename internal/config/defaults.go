package config

const (
	// DefaultConfigFile is the default config file name vaultpull looks for.
	DefaultConfigFile = ".vaultpull.yaml"

	// DefaultMount is the default KV v2 mount path used when none is specified.
	DefaultMount = "secret"
)

// ApplyDefaults fills in optional fields with sensible defaults.
func ApplyDefaults(cfg *Config) {
	if cfg.Vault.Mount == "" {
		cfg.Vault.Mount = DefaultMount
	}

	// Allow token to be sourced from environment if not set in file.
	if cfg.Vault.Token == "" {
		if tok := lookupEnv("VAULT_TOKEN"); tok != "" {
			cfg.Vault.Token = tok
		}
	}

	// Allow address to be sourced from environment if not set in file.
	if cfg.Vault.Address == "" {
		if addr := lookupEnv("VAULT_ADDR"); addr != "" {
			cfg.Vault.Address = addr
		}
	}
}
