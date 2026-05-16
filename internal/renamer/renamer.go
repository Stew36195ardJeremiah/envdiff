// Package renamer provides utilities for renaming keys across env maps,
// supporting bulk renames via a mapping and optional case normalization.
package renamer

import "strings"

// Options controls renaming behaviour.
type Options struct {
	// IgnoreCase performs case-insensitive lookup when matching keys to rename.
	IgnoreCase bool
	// FailOnMissing returns an error if a rename source key is not found.
	FailOnMissing bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		IgnoreCase:    false,
		FailOnMissing: false,
	}
}

// Result holds the output of a rename operation.
type Result struct {
	// Env is the transformed key-value map.
	Env map[string]string
	// Renamed lists keys that were successfully renamed (old -> new).
	Renamed map[string]string
	// NotFound lists source keys from the mapping that were not present in the env.
	NotFound []string
}

// Rename applies the provided mapping to env, returning a new map with keys
// replaced according to the mapping. The original map is never modified.
func Rename(env map[string]string, mapping map[string]string, opts Options) (Result, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	renamed := make(map[string]string)
	var notFound []string

	for oldKey, newKey := range mapping {
		matched := ""
		if opts.IgnoreCase {
			for k := range out {
				if strings.EqualFold(k, oldKey) {
					matched = k
					break
				}
			}
		} else {
			if _, ok := out[oldKey]; ok {
				matched = oldKey
			}
		}

		if matched == "" {
			notFound = append(notFound, oldKey)
			if opts.FailOnMissing {
				return Result{}, &MissingKeyError{Key: oldKey}
			}
			continue
		}

		out[newKey] = out[matched]
		delete(out, matched)
		renamed[matched] = newKey
	}

	return Result{Env: out, Renamed: renamed, NotFound: notFound}, nil
}

// MissingKeyError is returned when FailOnMissing is set and a source key is absent.
type MissingKeyError struct {
	Key string
}

func (e *MissingKeyError) Error() string {
	return "renamer: key not found: " + e.Key
}
