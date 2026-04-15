package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Target describes a single Vault secret path → local .env file mapping.
type Target struct {
	SecretPath string `yaml:"secret_path"`
	EnvFile    string `yaml:"env_file"`
}

// Config is the top-level configuration for vaultpull.
type Config struct {
	Address string   `yaml:"address"`
	Token   string   `yaml:"token"`
	Mount   string   `yaml:"mount"`
	Targets []Target `yaml:"targets"`
}

// Load reads a YAML config file from path, applies defaults from environment
// variables, and returns the resulting Config.
// It does NOT run validation — call Validate separately.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	ApplyDefaults(&cfg)

	if cfg.Address == "" {
		return nil, fmt.Errorf("vault address is required (set 'address' in config or VAULT_ADDR)")
	}
	if len(cfg.Targets) == 0 {
		return nil, fmt.Errorf("no sync targets defined in config")
	}

	return &cfg, nil
}
