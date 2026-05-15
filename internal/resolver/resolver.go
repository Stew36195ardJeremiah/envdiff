// Package resolver provides utilities for resolving placeholder references
// and variable interpolation within parsed .env maps.
package resolver

import (
	"fmt"
	"regexp"
	"strings"
)

// varPattern matches ${VAR_NAME} and $VAR_NAME style references.
var varPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

// Options controls resolver behaviour.
type Options struct {
	// MaxDepth limits recursive resolution to prevent cycles.
	MaxDepth int
	// FailOnMissing returns an error when a referenced variable is absent.
	FailOnMissing bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaxDepth:      10,
		FailOnMissing: false,
	}
}

// Resolve expands variable references inside env values using the same map.
// It returns a new map with all resolvable references expanded.
func Resolve(env map[string]string, opts Options) (map[string]string, error) {
	result := make(map[string]string, len(env))
	for k, v := range env {
		result[k] = v
	}

	for k := range result {
		resolved, err := resolveValue(result[k], result, opts.MaxDepth, opts.FailOnMissing)
		if err != nil {
			return nil, fmt.Errorf("resolver: key %q: %w", k, err)
		}
		result[k] = resolved
	}
	return result, nil
}

func resolveValue(val string, env map[string]string, depth int, failOnMissing bool) (string, error) {
	if depth <= 0 {
		return val, fmt.Errorf("max resolution depth exceeded (possible cycle)")
	}
	var resolveErr error
	result := varPattern.ReplaceAllStringFunc(val, func(match string) string {
		submatches := varPattern.FindStringSubmatch(match)
		name := submatches[1]
		if name == "" {
			name = submatches[2]
		}
		replacement, ok := env[name]
		if !ok {
			if failOnMissing {
				resolveErr = fmt.Errorf("undefined variable %q", name)
			}
			return match
		}
		if strings.ContainsAny(replacement, "$") {
			inner, err := resolveValue(replacement, env, depth-1, failOnMissing)
			if err != nil {
				resolveErr = err
				return match
			}
			return inner
		}
		return replacement
	})
	if resolveErr != nil {
		return "", resolveErr
	}
	return result, nil
}
