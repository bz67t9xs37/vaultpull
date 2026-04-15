package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level vaultpull configuration.
type Config struct {
	Vault   VaultConfig    `yaml:"vault"`
	Targets []TargetConfig `yaml:"targets"`
}

// VaultConfig holds Vault connection settings.
type VaultConfig struct {
	Address string `yaml:"address"`
	Token   string `yaml:"token"`
	Mount   string `yaml:"mount"`
}

// TargetConfig maps a Vault secret path to a local .env file.
type TargetConfig struct {
	Secret  string `yaml:"secret"`
	EnvFile string `yaml:"env_file"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: reading file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parsing YAML: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Vault.Address == "" {
		return fmt.Errorf("config: vault.address is required")
	}
	if c.Vault.Token == "" {
		return fmt.Errorf("config: vault.token is required")
	}
	if len(c.Targets) == 0 {
		return fmt.Errorf("config: at least one target must be defined")
	}
	for i, t := range c.Targets {
		if t.Secret == "" {
			return fmt.Errorf("config: targets[%d].secret is required", i)
		}
		if t.EnvFile == "" {
			return fmt.Errorf("config: targets[%d].env_file is required", i)
		}
	}
	return nil
}
