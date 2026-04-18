package scrub_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/scrub"
)

func TestLine_ReplacesSecret(t *testing.T) {
	s := scrub.New([]string{"supersecret"}, "")
	got := s.Line("password=supersecret here")
	if got != "password=[REDACTED] here" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestLine_MultipleSecrets(t *testing.T) {
	s := scrub.New([]string{"alpha", "beta"}, "***")
	got := s.Line("alpha and beta")
	if got != "*** and ***" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestLine_ShortSecretIgnored(t *testing.T) {
	s := scrub.New([]string{"abc"}, "")
	got := s.Line("abc in text")
	if got != "abc in text" {
		t.Errorf("short secret should not be scrubbed: %q", got)
	}
}

func TestLines_ScrubsAll(t *testing.T) {
	s := scrub.New([]string{"token1234"}, "")
	input := []string{"key=token1234", "other=plain"}
	out := s.Lines(input)
	if out[0] != "key=[REDACTED]" {
		t.Errorf("line 0: %q", out[0])
	}
	if out[1] != "other=plain" {
		t.Errorf("line 1: %q", out[1])
	}
}

func TestMap_ScrubsValues(t *testing.T) {
	s := scrub.New([]string{"mysecret"}, "")
	m := map[string]string{"KEY": "mysecret", "OTHER": "public"}
	out := s.Map(m)
	if out["KEY"] != "[REDACTED]" {
		t.Errorf("KEY: %q", out["KEY"])
	}
	if out["OTHER"] != "public" {
		t.Errorf("OTHER: %q", out["OTHER"])
	}
}

func TestMap_KeysUntouched(t *testing.T) {
	s := scrub.New([]string{"MYKEY"}, "")
	m := map[string]string{"MYKEY": "value"}
	out := s.Map(m)
	if _, ok := out["MYKEY"]; !ok {
		t.Error("key should be preserved")
	}
}
