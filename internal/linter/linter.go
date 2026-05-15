// Package linter provides style and convention checks for .env file entries.
// It enforces naming conventions, value formatting, and ordering rules.
package linter

import (
	"fmt"
	"strings"
	"unicode"
)

// Issue represents a single linting violation.
type Issue struct {
	Key     string
	Message string
	Severity string // "error" or "warning"
}

// Options controls which lint checks are enabled.
type Options struct {
	RequireUppercase  bool
	ForbidSpaceInValue bool
	RequireNonEmpty   bool
	MaxKeyLength      int
}

// DefaultOptions returns a sensible default lint configuration.
func DefaultOptions() Options {
	return Options{
		RequireUppercase:   true,
		ForbidSpaceInValue: false,
		RequireNonEmpty:    false,
		MaxKeyLength:       64,
	}
}

// Lint checks the provided key-value map against the given options
// and returns a slice of Issues found.
func Lint(env map[string]string, opts Options) []Issue {
	var issues []Issue

	for key, value := range env {
		if opts.RequireUppercase && !isUpperSnakeCase(key) {
			issues = append(issues, Issue{
				Key:      key,
				Message:  fmt.Sprintf("key %q should be UPPER_SNAKE_CASE", key),
				Severity: "warning",
			})
		}

		if opts.MaxKeyLength > 0 && len(key) > opts.MaxKeyLength {
			issues = append(issues, Issue{
				Key:      key,
				Message:  fmt.Sprintf("key %q exceeds max length of %d", key, opts.MaxKeyLength),
				Severity: "error",
			})
		}

		if opts.RequireNonEmpty && strings.TrimSpace(value) == "" {
			issues = append(issues, Issue{
				Key:      key,
				Message:  fmt.Sprintf("key %q has an empty or blank value", key),
				Severity: "warning",
			})
		}

		if opts.ForbidSpaceInValue && strings.Contains(value, " ") {
			issues = append(issues, Issue{
				Key:      key,
				Message:  fmt.Sprintf("key %q value contains unquoted spaces", key),
				Severity: "warning",
			})
		}
	}

	return issues
}

// isUpperSnakeCase returns true if every character in the key is
// an uppercase letter, digit, or underscore.
func isUpperSnakeCase(key string) bool {
	if key == "" {
		return false
	}
	for _, r := range key {
		if !unicode.IsUpper(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}
