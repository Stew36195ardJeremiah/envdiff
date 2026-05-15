package resolver_test

import (
	"testing"

	"github.com/nicholasgasior/envdiff/internal/parser"
	"github.com/nicholasgasior/envdiff/internal/resolver"
)

func TestResolve_WithParsedFile(t *testing.T) {
	env, err := parser.ParseFile("testdata/interpolation.env")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	result, err := resolver.Resolve(env, resolver.DefaultOptions())
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}

	want := map[string]string{
		"HOST":     "db.example.com",
		"PORT":     "5432",
		"DB_URL":   "postgres://db.example.com:5432/mydb",
		"APP_NAME": "envdiff",
		"BANNER":   "Welcome to envdiff",
	}
	for k, wantVal := range want {
		if result[k] != wantVal {
			t.Errorf("key %q: got %q, want %q", k, result[k], wantVal)
		}
	}
}
