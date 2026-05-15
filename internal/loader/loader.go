// Package loader provides utilities for loading and validating one or more
// .env files before passing them to the diff pipeline.
package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/envdiff/internal/parser"
)

// EnvFile represents a loaded environment file with its path and parsed entries.
type EnvFile struct {
	Path    string
	Entries map[string]string
}

// LoadOptions controls optional behaviour when loading files.
type LoadOptions struct {
	// SkipMissing silences errors when a file does not exist on disk.
	SkipMissing bool
}

// Load reads one or more .env file paths and returns a slice of EnvFile.
// If a path does not exist and SkipMissing is false, an error is returned.
func Load(paths []string, opts LoadOptions) ([]EnvFile, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("loader: no file paths provided")
	}

	results := make([]EnvFile, 0, len(paths))

	for _, p := range paths {
		clean := filepath.Clean(p)

		if _, err := os.Stat(clean); os.IsNotExist(err) {
			if opts.SkipMissing {
				continue
			}
			return nil, fmt.Errorf("loader: file not found: %s", clean)
		}

		entries, err := parser.ParseFile(clean)
		if err != nil {
			return nil, fmt.Errorf("loader: failed to parse %s: %w", clean, err)
		}

		results = append(results, EnvFile{
			Path:    clean,
			Entries: entries,
		})
	}

	return results, nil
}

// MustLoad is like Load but panics on error. Useful in tests or CLI init.
func MustLoad(paths []string, opts LoadOptions) []EnvFile {
	files, err := Load(paths, opts)
	if err != nil {
		panic(err)
	}
	return files
}
