// Package validator provides lightweight static analysis for parsed .env maps.
//
// It checks for common authoring mistakes including:
//   - Empty or blank values
//   - Keys that do not follow the conventional UPPER_SNAKE_CASE naming style
//   - Values that contain leading or trailing whitespace
//
// Usage:
//
//	env, _ := parser.ParseFile(".env")
//	issues := validator.Validate(env, validator.DefaultOptions())
//	for _, iss := range issues {
//		fmt.Println(iss)
//	}
//
// Each check can be toggled independently via [Options].
package validator
