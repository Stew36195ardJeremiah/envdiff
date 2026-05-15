package redactor_test

import (
	"testing"

	"github.com/nicholasgasior/envdiff/internal/redactor"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME":     "myapp",
		"DB_PASSWORD":  "supersecret",
		"API_KEY":      "abc123",
		"AUTH_TOKEN":   "tok_xyz",
		"PORT":         "8080",
		"PRIVATE_KEY":  "-----BEGIN RSA",
		"LOG_LEVEL":    "info",
	}
}

func TestRedact_DefaultPatterns(t *testing.T) {
	out := redactor.Redact(baseEnv(), redactor.Options{})

	sensitiveKeys := []string{"DB_PASSWORD", "API_KEY", "AUTH_TOKEN", "PRIVATE_KEY"}
	for _, k := range sensitiveKeys {
		if out[k] != "***" {
			t.Errorf("expected %s to be redacted, got %q", k, out[k])
		}
	}

	plainKeys := []string{"APP_NAME", "PORT", "LOG_LEVEL"}
	for _, k := range plainKeys {
		if out[k] == "***" {
			t.Errorf("expected %s to be plain, but it was redacted", k)
		}
	}
}

func TestRedact_CustomPlaceholder(t *testing.T) {
	out := redactor.Redact(baseEnv(), redactor.Options{Placeholder: "[HIDDEN]"})
	if out["DB_PASSWORD"] != "[HIDDEN]" {
		t.Errorf("expected [HIDDEN], got %q", out["DB_PASSWORD"])
	}
}

func TestRedact_CustomPatterns(t *testing.T) {
	env := map[string]string{
		"STRIPE_KEY": "sk_live_abc",
		"APP_PORT":   "3000",
	}
	out := redactor.Redact(env, redactor.Options{Patterns: []string{"stripe"}})
	if out["STRIPE_KEY"] != "***" {
		t.Errorf("expected STRIPE_KEY redacted, got %q", out["STRIPE_KEY"])
	}
	if out["APP_PORT"] != "3000" {
		t.Errorf("expected APP_PORT unchanged, got %q", out["APP_PORT"])
	}
}

func TestRedact_OriginalUnmodified(t *testing.T) {
	original := baseEnv()
	redactor.Redact(original, redactor.Options{})
	if original["DB_PASSWORD"] != "supersecret" {
		t.Error("original map was modified")
	}
}

func TestIsSensitiveKey(t *testing.T) {
	patterns := redactor.DefaultSensitivePatterns
	cases := []struct {
		key       string
		expected  bool
	}{
		{"DB_PASSWORD", true},
		{"SECRET_KEY", true},
		{"APP_NAME", false},
		{"PORT", false},
		{"OAUTH_TOKEN", true},
	}
	for _, tc := range cases {
		got := redactor.IsSensitiveKey(tc.key, patterns)
		if got != tc.expected {
			t.Errorf("IsSensitiveKey(%q) = %v, want %v", tc.key, got, tc.expected)
		}
	}
}
