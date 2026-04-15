package config

import (
	"fmt"
	"time"
)

// TTLConfig holds TTL settings for secret caching within vaultpull.
type TTLConfig struct {
	// Enabled controls whether TTL tracking is active.
	Enabled bool `yaml:"enabled"`

	// DefaultTTL is the duration after which a fetched secret is considered stale.
	// Accepts Go duration strings, e.g. "10m", "1h".
	DefaultTTL string `yaml:"default_ttl"`
}

// DefaultTTLConfig returns a TTLConfig with sensible defaults.
func DefaultTTLConfig() TTLConfig {
	return TTLConfig{
		Enabled:    false,
		DefaultTTL: "30m",
	}
}

// ParsedTTL returns the DefaultTTL as a time.Duration.
// Returns an error if the value is not a valid duration string.
func (c TTLConfig) ParsedTTL() (time.Duration, error) {
	if c.DefaultTTL == "" {
		return 30 * time.Minute, nil
	}
	d, err := time.ParseDuration(c.DefaultTTL)
	if err != nil {
		return 0, fmt.Errorf("invalid ttl %q: %w", c.DefaultTTL, err)
	}
	return d, nil
}

// Validate checks that the TTLConfig is well-formed.
func (c TTLConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	_, err := c.ParsedTTL()
	return err
}
