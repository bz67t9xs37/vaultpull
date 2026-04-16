package config

import "time"

// RetryConfig holds retry policy configuration.
type RetryConfig struct {
	MaxAttempts int           `yaml:"max_attempts"`
	InitialDelay time.Duration `yaml:"initial_delay"`
	MaxDelay     time.Duration `yaml:"max_delay"`
	Multiplier   float64       `yaml:"multiplier"`
}

// DefaultRetryConfig returns a RetryConfig with sensible defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 500 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
	}
}

// ApplyRetryDefaults fills zero values in c with defaults.
func ApplyRetryDefaults(c *RetryConfig) {
	if c == nil {
		return
	}
	d := DefaultRetryConfig()
	if c.MaxAttempts == 0 {
		c.MaxAttempts = d.MaxAttempts
	}
	if c.InitialDelay == 0 {
		c.InitialDelay = d.InitialDelay
	}
	if c.MaxDelay == 0 {
		c.MaxDelay = d.MaxDelay
	}
	if c.Multiplier == 0 {
		c.Multiplier = d.Multiplier
	}
}

// ToPolicy converts RetryConfig into values suitable for retry.Policy.
func (c *RetryConfig) ToPolicy() (maxAttempts int, initialDelay, maxDelay time.Duration, multiplier float64) {
	return c.MaxAttempts, c.InitialDelay, c.MaxDelay, c.Multiplier
}
