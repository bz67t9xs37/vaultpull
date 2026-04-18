package scrub_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/scrub"
)

func TestScrub_UsesConfigMinLength(t *testing.T) {
	cfg := &config.ScrubConfig{
		Enabled:   true,
		MinLength: 8,
	}
	s := scrub.New(cfg.MinLength)

	result := s.Line("connect with token=abcdefgh here", []string{"abcdefgh"})
	if result == "connect with token=abcdefgh here" {
		t.Error("expected secret to be scrubbed")
	}
}

func TestScrub_ShortSecretNotScrubbed_ViaConfig(t *testing.T) {
	cfg := &config.ScrubConfig{
		Enabled:   true,
		MinLength: 8,
	}
	s := scrub.New(cfg.MinLength)

	result := s.Line("value=abc here", []string{"abc"})
	if result != "value=abc here" {
		t.Errorf("expected short secret to remain, got: %s", result)
	}
}

func TestScrub_ConfigDefaultsApplied(t *testing.T) {
	cfg := &config.Config{}
	config.ApplyScrubDefaults(cfg)

	if cfg.Scrub == nil {
		t.Fatal("expected Scrub config to be initialized")
	}

	s := scrub.New(cfg.Scrub.MinLength)
	secret := "supersecretvalue"
	result := s.Line("msg: supersecretvalue end", []string{secret})
	if result == "msg: supersecretvalue end" {
		t.Error("expected secret to be redacted using default config")
	}
}
