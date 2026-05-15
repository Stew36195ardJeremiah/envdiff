package resolver

import (
	"testing"
)

func TestResolve_NoReferences(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	result, err := Resolve(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["HOST"] != "localhost" || result["PORT"] != "5432" {
		t.Errorf("values should be unchanged")
	}
}

func TestResolve_BraceStyle(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "https://${HOST}:${PORT}",
		"HOST":     "example.com",
		"PORT":     "8080",
	}
	result, err := Resolve(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["BASE_URL"] != "https://example.com:8080" {
		t.Errorf("got %q, want %q", result["BASE_URL"], "https://example.com:8080")
	}
}

func TestResolve_NoBraceStyle(t *testing.T) {
	env := map[string]string{
		"GREETING": "Hello $NAME",
		"NAME":     "World",
	}
	result, err := Resolve(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["GREETING"] != "Hello World" {
		t.Errorf("got %q", result["GREETING"])
	}
}

func TestResolve_MissingVar_NoFail(t *testing.T) {
	env := map[string]string{
		"URL": "http://${UNKNOWN_HOST}/path",
	}
	result, err := Resolve(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// unreferenced variable should remain as-is
	if result["URL"] != "http://${UNKNOWN_HOST}/path" {
		t.Errorf("got %q", result["URL"])
	}
}

func TestResolve_MissingVar_FailOnMissing(t *testing.T) {
	env := map[string]string{
		"URL": "http://${MISSING}/path",
	}
	opts := DefaultOptions()
	opts.FailOnMissing = true
	_, err := Resolve(env, opts)
	if err == nil {
		t.Fatal("expected error for missing variable, got nil")
	}
}

func TestResolve_OriginalUnmodified(t *testing.T) {
	env := map[string]string{
		"A": "${B}",
		"B": "value",
	}
	original := map[string]string{
		"A": "${B}",
		"B": "value",
	}
	_, _ = Resolve(env, DefaultOptions())
	for k, v := range original {
		if env[k] != v {
			t.Errorf("original map mutated: key %q changed to %q", k, env[k])
		}
	}
}
