package filter_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/filter"
)

var sampleEnv = map[string]string{
	"APP_HOST":     "localhost",
	"APP_PORT":     "8080",
	"DB_HOST":      "db.local",
	"DB_PASSWORD":  "secret",
	"SECRET_TOKEN": "abc123",
}

func TestApply_NoOptions(t *testing.T) {
	out := filter.Apply(sampleEnv, filter.Options{})
	if len(out) != len(sampleEnv) {
		t.Fatalf("expected %d keys, got %d", len(sampleEnv), len(out))
	}
}

func TestApply_Prefix(t *testing.T) {
	out := filter.Apply(sampleEnv, filter.Options{Prefix: "APP_"})
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
}

func TestApply_ExcludeKeys(t *testing.T) {
	out := filter.Apply(sampleEnv, filter.Options{ExcludeKeys: []string{"DB_PASSWORD", "SECRET_TOKEN"}})
	if _, found := out["DB_PASSWORD"]; found {
		t.Error("DB_PASSWORD should be excluded")
	}
	if _, found := out["SECRET_TOKEN"]; found {
		t.Error("SECRET_TOKEN should be excluded")
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(out))
	}
}

func TestApply_IgnoreCase(t *testing.T) {
	out := filter.Apply(sampleEnv, filter.Options{
		Prefix:     "app_",
		IgnoreCase: true,
	})
	if len(out) != 2 {
		t.Fatalf("expected 2 keys with case-insensitive prefix, got %d", len(out))
	}
}

func TestApply_ExcludeIgnoreCase(t *testing.T) {
	out := filter.Apply(sampleEnv, filter.Options{
		ExcludeKeys: []string{"db_password"},
		IgnoreCase:  true,
	})
	if _, found := out["DB_PASSWORD"]; found {
		t.Error("DB_PASSWORD should be excluded via case-insensitive match")
	}
}

func TestMatchesPrefix(t *testing.T) {
	if !filter.MatchesPrefix("APP_PORT", "APP_", false) {
		t.Error("expected match")
	}
	if filter.MatchesPrefix("DB_HOST", "APP_", false) {
		t.Error("expected no match")
	}
	if !filter.MatchesPrefix("app_port", "APP_", true) {
		t.Error("expected case-insensitive match")
	}
}
