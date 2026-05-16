// Package templater provides functionality for generating .env template files
// from parsed environment maps. Templates replace values with placeholder
// descriptions, making it easy to share required keys without exposing secrets.
package templater

import (
	"fmt"
	"sort"
	"strings"
)

// Options controls template generation behavior.
type Options struct {
	// Placeholder is used as the value in the generated template.
	// Defaults to "<value>" if empty.
	Placeholder string

	// DescribeTypes annotates values with a type hint based on the original value.
	DescribeTypes bool

	// IncludeComments adds a header comment to the generated template.
	IncludeComments bool
}

// DefaultOptions returns sensible default options for template generation.
func DefaultOptions() Options {
	return Options{
		Placeholder:     "<value>",
		DescribeTypes:   false,
		IncludeComments: true,
	}
}

// Generate takes an env map and produces a template string with placeholders.
func Generate(env map[string]string, opts Options) string {
	if opts.Placeholder == "" {
		opts.Placeholder = "<value>"
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder

	if opts.IncludeComments {
		sb.WriteString("# Generated .env template — fill in the values below\n")
		sb.WriteString("# Do not commit real secrets to version control\n\n")
	}

	for _, k := range keys {
		origVal := env[k]
		placeholder := opts.Placeholder
		if opts.DescribeTypes {
			placeholder = describeType(origVal)
		}
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, placeholder))
	}

	return sb.String()
}

// describeType returns a type-hinted placeholder based on the value's shape.
func describeType(val string) string {
	if val == "true" || val == "false" {
		return "<bool>"
	}
	if isNumeric(val) {
		return "<number>"
	}
	if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
		return "<url>"
	}
	return "<string>"
}

// isNumeric reports whether s looks like an integer or float.
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	dots := 0
	for i, c := range s {
		if c == '.' {
			dots++
			if dots > 1 {
				return false
			}
			continue
		}
		if c == '-' && i == 0 {
			continue
		}
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
