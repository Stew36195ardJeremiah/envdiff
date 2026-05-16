package renamer_test

import (
	"testing"

	"github.com/nicholasgasior/envdiff/internal/renamer"
)

var baseEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "secret",
	"APP_PORT":    "8080",
}

func TestRename_BasicMapping(t *testing.T) {
	mapping := map[string]string{"DB_HOST": "DATABASE_HOST"}
	res, err := renamer.Rename(baseEnv, mapping, renamer.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["DATABASE_HOST"]; !ok {
		t.Error("expected DATABASE_HOST in result")
	}
	if _, ok := res.Env["DB_HOST"]; ok {
		t.Error("old key DB_HOST should be removed")
	}
	if res.Renamed["DB_HOST"] != "DATABASE_HOST" {
		t.Errorf("renamed map mismatch: %v", res.Renamed)
	}
}

func TestRename_OriginalUnmodified(t *testing.T) {
	mapping := map[string]string{"APP_PORT": "PORT"}
	_, err := renamer.Rename(baseEnv, mapping, renamer.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := baseEnv["APP_PORT"]; !ok {
		t.Error("original map should not be modified")
	}
}

func TestRename_MissingKey_NoFail(t *testing.T) {
	mapping := map[string]string{"NONEXISTENT": "NEW_KEY"}
	res, err := renamer.Rename(baseEnv, mapping, renamer.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.NotFound) != 1 || res.NotFound[0] != "NONEXISTENT" {
		t.Errorf("expected NotFound=[NONEXISTENT], got %v", res.NotFound)
	}
}

func TestRename_MissingKey_FailOnMissing(t *testing.T) {
	opts := renamer.DefaultOptions()
	opts.FailOnMissing = true
	mapping := map[string]string{"NONEXISTENT": "NEW_KEY"}
	_, err := renamer.Rename(baseEnv, mapping, opts)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	var mke *renamer.MissingKeyError
	if me, ok := err.(*renamer.MissingKeyError); !ok || me.Key != "NONEXISTENT" {
		_ = mke
		t.Errorf("unexpected error type or key: %v", err)
	}
}

func TestRename_IgnoreCase(t *testing.T) {
	opts := renamer.DefaultOptions()
	opts.IgnoreCase = true
	mapping := map[string]string{"db_host": "DATABASE_HOST"}
	res, err := renamer.Rename(baseEnv, mapping, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["DATABASE_HOST"]; !ok {
		t.Error("expected DATABASE_HOST after case-insensitive rename")
	}
	if _, ok := res.Env["DB_HOST"]; ok {
		t.Error("old key should be removed after case-insensitive rename")
	}
}

func TestRename_MultipleKeys(t *testing.T) {
	mapping := map[string]string{
		"DB_HOST":     "DATABASE_HOST",
		"DB_PASSWORD": "DATABASE_PASSWORD",
	}
	res, err := renamer.Rename(baseEnv, mapping, renamer.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Renamed) != 2 {
		t.Errorf("expected 2 renames, got %d", len(res.Renamed))
	}
	if res.Env["DATABASE_HOST"] != "localhost" {
		t.Errorf("value mismatch for DATABASE_HOST: %s", res.Env["DATABASE_HOST"])
	}
}
