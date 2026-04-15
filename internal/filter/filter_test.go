package filter_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/filter"
)

var testSecrets = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "secret",
	"APP_KEY":     "abc123",
	"APP_SECRET":  "xyz789",
	"LOG_LEVEL":   "info",
}

func TestApply_NoRules_ReturnsAll(t *testing.T) {
	f := filter.New(nil, nil)
	result := f.Apply(testSecrets)
	if len(result) != len(testSecrets) {
		t.Errorf("expected %d keys, got %d", len(testSecrets), len(result))
	}
}

func TestApply_IncludeByPrefix(t *testing.T) {
	f := filter.New([]filter.Rule{{Prefix: "DB_"}}, nil)
	result := f.Apply(testSecrets)
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
	if _, ok := result["DB_PASSWORD"]; !ok {
		t.Error("expected DB_PASSWORD in result")
	}
}

func TestApply_ExcludeByContains(t *testing.T) {
	f := filter.New(nil, []filter.Rule{{Contains: "SECRET"}})
	result := f.Apply(testSecrets)
	if _, ok := result["APP_SECRET"]; ok {
		t.Error("APP_SECRET should have been excluded")
	}
	if _, ok := result["DB_PASSWORD"]; !ok {
		t.Error("DB_PASSWORD should remain")
	}
}

func TestApply_IncludeAndExcludeCombined(t *testing.T) {
	include := []filter.Rule{{Prefix: "APP_"}}
	exclude := []filter.Rule{{Suffix: "_SECRET"}}
	f := filter.New(include, exclude)
	result := f.Apply(testSecrets)
	if len(result) != 1 {
		t.Errorf("expected 1 key, got %d", len(result))
	}
	if _, ok := result["APP_KEY"]; !ok {
		t.Error("expected APP_KEY in result")
	}
}

func TestApply_ExcludeAll(t *testing.T) {
	exclude := []filter.Rule{{Contains: "_"}}
	f := filter.New(nil, exclude)
	result := f.Apply(testSecrets)
	if len(result) != 0 {
		t.Errorf("expected 0 keys, got %d", len(result))
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	f := filter.New([]filter.Rule{{Prefix: "DB_"}}, nil)
	result := f.Apply(map[string]string{})
	if len(result) != 0 {
		t.Errorf("expected 0 keys, got %d", len(result))
	}
}
