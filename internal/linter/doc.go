// Package linter implements style and convention linting for .env file entries.
//
// It operates on a parsed map[string]string (as produced by the parser package)
// and applies configurable checks including:
//
//   - Key naming convention (UPPER_SNAKE_CASE)
//   - Maximum key length
//   - Empty or whitespace-only values
//   - Unquoted spaces in values
//
// # Usage
//
//	env, _ := parser.ParseFile("production.env")
//	issues := linter.Lint(env, linter.DefaultOptions())
//	for _, issue := range issues {
//		fmt.Printf("[%s] %s: %s\n", issue.Severity, issue.Key, issue.Message)
//	}
//
// Each Issue carries a Severity of either "error" or "warning" so callers
// can decide whether to fail a CI pipeline or merely print advisories.
package linter
