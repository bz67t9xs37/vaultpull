package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Target describes a single Vault path → local file mapping.
type Target struct {
	Path    string `yaml:"path"`
	Output  string `yaml:"output"`
	Backup  bool   `yaml:"backup"`
}

// Config holds the full vaultpull configuration.
type Config struct {
	Address string   `yaml:"address"`
	Token   string   `yaml:"token"`
	Mount   string   `yaml:"mount"`
	BackupDir string `yaml:"backup_dir"`
	Targets []Target `yaml:"targets"`
}

// Load reads and parses a YAML config file at path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse %q: %w", path, err)
	}

	ApplyDefaults(&cfg)

	if err := Validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
