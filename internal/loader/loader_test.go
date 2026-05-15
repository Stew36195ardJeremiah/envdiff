package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/loader"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestLoad_SingleFile(t *testing.T) {
	p := writeTempEnv(t, "KEY=value\nFOO=bar\n")
	files, err := loader.Load([]string{p}, loader.LoadOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	if files[0].Entries["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", files[0].Entries["KEY"])
	}
}

func TestLoad_MultipleFiles(t *testing.T) {
	p1 := writeTempEnv(t, "A=1\n")
	p2 := writeTempEnv(t, "B=2\n")
	files, err := loader.Load([]string{p1, p2}, loader.LoadOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
}

func TestLoad_FileNotFound_Error(t *testing.T) {
	_, err := loader.Load([]string{"/nonexistent/.env"}, loader.LoadOptions{})
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_FileNotFound_SkipMissing(t *testing.T) {
	p := writeTempEnv(t, "X=1\n")
	files, err := loader.Load([]string{p, "/nonexistent/.env"}, loader.LoadOptions{SkipMissing: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file after skipping missing, got %d", len(files))
	}
}

func TestLoad_NoPaths(t *testing.T) {
	_, err := loader.Load([]string{}, loader.LoadOptions{})
	if err == nil {
		t.Fatal("expected error for empty path list, got nil")
	}
}
