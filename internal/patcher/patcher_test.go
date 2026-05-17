package patcher_test

import (
	"testing"

	"github.com/nicholasgasior/envdiff/internal/patcher"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_ENV":  "staging",
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"LOG_LEVEL": "info",
	}
}

func TestApply_Upsert_NewKey(t *testing.T) {
	res, err := patcher.Apply(baseEnv(), []patcher.Patch{
		{Key: "NEW_KEY", Value: "hello", Op: patcher.OpUpsert},
	}, patcher.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q", res.Env["NEW_KEY"])
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied patch, got %d", len(res.Applied))
	}
}

func TestApply_Upsert_OverwriteKey(t *testing.T) {
	res, err := patcher.Apply(baseEnv(), []patcher.Patch{
		{Key: "APP_ENV", Value: "production", Op: patcher.OpUpsert},
	}, patcher.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", res.Env["APP_ENV"])
	}
}

func TestApply_Delete_ExistingKey(t *testing.T) {
	res, err := patcher.Apply(baseEnv(), []patcher.Patch{
		{Key: "LOG_LEVEL", Op: patcher.OpDelete},
	}, patcher.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["LOG_LEVEL"]; ok {
		t.Error("expected LOG_LEVEL to be deleted")
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied patch, got %d", len(res.Applied))
	}
}

func TestApply_Delete_MissingKey_Skip(t *testing.T) {
	res, err := patcher.Apply(baseEnv(), []patcher.Patch{
		{Key: "GHOST_KEY", Op: patcher.OpDelete},
	}, patcher.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped patch, got %d", len(res.Skipped))
	}
}

func TestApply_Delete_MissingKey_FailOnMissing(t *testing.T) {
	opts := patcher.Options{FailOnMissing: true}
	_, err := patcher.Apply(baseEnv(), []patcher.Patch{
		{Key: "GHOST_KEY", Op: patcher.OpDelete},
	}, opts)
	if err == nil {
		t.Error("expected error for missing delete target, got nil")
	}
}

func TestApply_OriginalUnmodified(t *testing.T) {
	original := baseEnv()
	patcher.Apply(original, []patcher.Patch{
		{Key: "APP_ENV", Value: "production", Op: patcher.OpUpsert},
		{Key: "DB_PORT", Op: patcher.OpDelete},
	}, patcher.DefaultOptions())
	if original["APP_ENV"] != "staging" {
		t.Error("original env was mutated")
	}
	if _, ok := original["DB_PORT"]; !ok {
		t.Error("original env was mutated: DB_PORT deleted")
	}
}
