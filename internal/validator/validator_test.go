package validator_test

import (
	"testing"

	"github.com/nicholasgasior/envdiff/internal/validator"
)

func TestValidate_NoIssues(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"PORT":         "8080",
	}
	issues := validator.Validate(env, validator.DefaultOptions())
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	env := map[string]string{
		"SECRET_KEY": "",
	}
	issues := validator.Validate(env, validator.DefaultOptions())
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "SECRET_KEY" {
		t.Errorf("unexpected key %q", issues[0].Key)
	}
}

func TestValidate_InvalidKeyName(t *testing.T) {
	env := map[string]string{
		"bad-key": "value",
	}
	issues := validator.Validate(env, validator.DefaultOptions())
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestValidate_WhitespaceValue(t *testing.T) {
	env := map[string]string{
		"API_KEY": "  abc  ",
	}
	issues := validator.Validate(env, validator.DefaultOptions())
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestValidate_DisabledChecks(t *testing.T) {
	env := map[string]string{
		"bad-key": "",
	}
	opts := validator.Options{} // all checks disabled
	issues := validator.Validate(env, opts)
	if len(issues) != 0 {
		t.Fatalf("expected no issues with all checks disabled, got %d", len(issues))
	}
}

func TestIssue_String(t *testing.T) {
	i := validator.Issue{Key: "FOO", Message: "value is empty"}
	want := "FOO: value is empty"
	if i.String() != want {
		t.Errorf("got %q, want %q", i.String(), want)
	}
}
