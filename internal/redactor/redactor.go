// Package redactor provides utilities for redacting sensitive values
// in environment variable maps before display or export.
package redactor

import "strings"

// DefaultSensitivePatterns contains common substrings that indicate a key
// holds a sensitive value (case-insensitive match).
var DefaultSensitivePatterns = []string{
	"password",
	"passwd",
	"secret",
	"token",
	"api_key",
	"apikey",
	"private",
	"credential",
	"auth",
}

// Options controls redaction behaviour.
type Options struct {
	// Patterns is the list of substrings to match against key names.
	// Matching is case-insensitive. Defaults to DefaultSensitivePatterns.
	Patterns []string

	// Placeholder is the string used in place of a redacted value.
	// Defaults to "***".
	Placeholder string
}

// Redact returns a new map with sensitive values replaced by the placeholder.
// The original map is not modified.
func Redact(env map[string]string, opts Options) map[string]string {
	if opts.Placeholder == "" {
		opts.Placeholder = "***"
	}
	patterns := opts.Patterns
	if len(patterns) == 0 {
		patterns = DefaultSensitivePatterns
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		if isSensitive(k, patterns) {
			out[k] = opts.Placeholder
		} else {
			out[k] = v
		}
	}
	return out
}

// IsSensitiveKey reports whether key matches any of the provided patterns.
func IsSensitiveKey(key string, patterns []string) bool {
	return isSensitive(key, patterns)
}

func isSensitive(key string, patterns []string) bool {
	lower := strings.ToLower(key)
	for _, p := range patterns {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}
