package linter_test

import (
	"testing"

	"envdiff/internal/linter"
)

func TestLint_NoIssues(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"APP_PORT":     "8080",
	}
	issues := linter.Lint(env, linter.DefaultOptions())
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %+v", len(issues), issues)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	env := map[string]string{"db_host": "localhost"}
	issues := linter.Lint(env, linter.DefaultOptions())
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != "warning" {
		t.Errorf("expected warning severity, got %q", issues[0].Severity)
	}
}

func TestLint_KeyTooLong(t *testing.T) {
	key := "A_VERY_LONG_KEY_THAT_EXCEEDS_THE_MAXIMUM_ALLOWED_LENGTH_FOR_ENV_KEYS_XYZ"
	env := map[string]string{key: "value"}
	opts := linter.DefaultOptions()
	opts.MaxKeyLength = 64
	issues := linter.Lint(env, opts)
	found := false
	for _, issue := range issues {
		if issue.Severity == "error" {
			found = true
		}
	}
	if !found {
		t.Error("expected an error-level issue for key length")
	}
}

func TestLint_EmptyValue(t *testing.T) {
	env := map[string]string{"API_KEY": ""}
	opts := linter.DefaultOptions()
	opts.RequireNonEmpty = true
	issues := linter.Lint(env, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "API_KEY" {
		t.Errorf("expected issue on API_KEY, got %q", issues[0].Key)
	}
}

func TestLint_SpaceInValue(t *testing.T) {
	env := map[string]string{"APP_NAME": "my app"}
	opts := linter.DefaultOptions()
	opts.ForbidSpaceInValue = true
	issues := linter.Lint(env, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestLint_DisabledUppercase(t *testing.T) {
	env := map[string]string{"lower_key": "value"}
	opts := linter.DefaultOptions()
	opts.RequireUppercase = false
	issues := linter.Lint(env, opts)
	if len(issues) != 0 {
		t.Errorf("expected no issues with uppercase check disabled, got %d", len(issues))
	}
}
