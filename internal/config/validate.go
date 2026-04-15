package config

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// ValidationError holds a list of validation issues found in a Config.
type ValidationError struct {
	Errors []string
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("config validation failed:\n  - %s", strings.Join(v.Errors, "\n  - "))
}

// Validate checks that the Config is complete and well-formed.
// It returns a *ValidationError listing all problems, or nil if valid.
func Validate(cfg *Config) error {
	var errs []string

	if strings.TrimSpace(cfg.Address) == "" {
		errs = append(errs, "vault address is required")
	} else if _, err := url.ParseRequestURI(cfg.Address); err != nil {
		errs = append(errs, fmt.Sprintf("vault address %q is not a valid URL: %v", cfg.Address, err))
	}

	if strings.TrimSpace(cfg.Token) == "" {
		errs = append(errs, "vault token is required (set VAULT_TOKEN or token in config)")
	}

	if strings.TrimSpace(cfg.Mount) == "" {
		errs = append(errs, "secret mount path is required")
	}

	if len(cfg.Targets) == 0 {
		errs = append(errs, "at least one sync target must be defined")
	}

	for i, t := range cfg.Targets {
		if strings.TrimSpace(t.SecretPath) == "" {
			errs = append(errs, fmt.Sprintf("target[%d]: secret_path is required", i))
		}
		if strings.TrimSpace(t.EnvFile) == "" {
			errs = append(errs, fmt.Sprintf("target[%d]: env_file is required", i))
		}
	}

	if len(errs) > 0 {
		return &ValidationError{Errors: errs}
	}
	return nil
}

// IsValidationError reports whether err is a *ValidationError.
func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}
