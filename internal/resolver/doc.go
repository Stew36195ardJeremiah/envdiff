// Package resolver expands variable interpolation within parsed .env maps.
//
// Many .env files use shell-style variable references such as:
//
//	BASE_URL=https://${HOST}:${PORT}/api
//
// The resolver walks each value and substitutes references using the same
// map, supporting both ${VAR} and $VAR syntax. Recursive references are
// resolved up to a configurable MaxDepth to prevent infinite cycles.
//
// # Usage
//
//	env, _ := parser.ParseFile(".env")
//	resolved, err := resolver.Resolve(env, resolver.DefaultOptions())
//
// # Options
//
//   - MaxDepth (default 10): maximum recursion depth for nested references.
//   - FailOnMissing (default false): return an error when a referenced
//     variable is not present in the map; otherwise leave the reference
//     as-is.
package resolver
