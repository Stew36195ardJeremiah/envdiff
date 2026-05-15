// Package validator checks .env entries for common issues such as
// empty values, suspicious characters, and key naming conventions.
package validator

import (
	"fmt"
	"strings"
	"unicode"
)

// Issue represents a single validation problem found in an env map.
type Issue struct {
	Key     string
	Message string
}

func (i Issue) String() string {
	return fmt.Sprintf("%s: %s", i.Key, i.Message)
}

// Options controls which checks are performed.
type Options struct {
	WarnEmptyValues    bool
	WarnInvalidKeyName bool
	WarnWhitespace     bool
}

// DefaultOptions returns a sensible default set of validation options.
func DefaultOptions() Options {
	return Options{
		WarnEmptyValues:    true,
		WarnInvalidKeyName: true,
		WarnWhitespace:     true,
	}
}

// Validate inspects the provided env map and returns a slice of Issues.
// An empty slice means no problems were found.
func Validate(env map[string]string, opts Options) []Issue {
	var issues []Issue

	for k, v := range env {
		if opts.WarnInvalidKeyName && !validKeyName(k) {
			issues = append(issues, Issue{Key: k, Message: "key contains invalid characters (expected A-Z, 0-9, _)"})
		}
		if opts.WarnEmptyValues && strings.TrimSpace(v) == "" {
			issues = append(issues, Issue{Key: k, Message: "value is empty"})
		}
		if opts.WarnWhitespace && v != strings.TrimSpace(v) {
			issues = append(issues, Issue{Key: k, Message: "value has leading or trailing whitespace"})
		}
	}

	return issues
}

// validKeyName returns true when every rune in name is an uppercase letter,
// digit, or underscore, and the name is non-empty.
func validKeyName(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if !unicode.IsUpper(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}
