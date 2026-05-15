package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/snapshot"
)

func sampleResult() diff.Result {
	return diff.Result{
		MissingInRight: []string{"DB_HOST"},
		MissingInLeft:  []string{"API_KEY"},
		Mismatched: []diff.Mismatch{
			{Key: "PORT", LeftValue: "8080", RightValue: "9090"},
		},
	}
}

func TestSave_And_Load_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	result := sampleResult()
	if err := snapshot.Save(path, "test-label", result); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Label != "test-label" {
		t.Errorf("expected label 'test-label', got %q", loaded.Label)
	}
	if loaded.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if len(loaded.Result.MissingInRight) != 1 || loaded.Result.MissingInRight[0] != "DB_HOST" {
		t.Errorf("unexpected MissingInRight: %v", loaded.Result.MissingInRight)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestSave_CreatedAt_IsRecent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	before := time.Now().UTC()

	_ = snapshot.Save(path, "", sampleResult())

	loaded, _ := snapshot.Load(path)
	if loaded.CreatedAt.Before(before) {
		t.Error("CreatedAt should be after test start time")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	err := snapshot.Save(filepath.Join(os.TempDir(), "nonexistent", "dir", "snap.json"), "", sampleResult())
	if err == nil {
		t.Error("expected error writing to invalid path")
	}
}
