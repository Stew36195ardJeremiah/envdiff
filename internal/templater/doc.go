// Package templater generates .env template files from parsed environment maps.
//
// A template replaces each value with a configurable placeholder string,
// making it safe to commit to version control as documentation of required
// environment variables without exposing sensitive data.
//
// # Basic Usage
//
//	env := map[string]string{
//		"DB_HOST":     "localhost",
//		"DB_PASSWORD": "supersecret",
//		"PORT":        "8080",
//	}
//
//	opts := templater.DefaultOptions()
//	template := templater.Generate(env, opts)
//	// DB_HOST=<value>
//	// DB_PASSWORD=<value>
//	// PORT=<value>
//
// # Type Hints
//
// When DescribeTypes is enabled, the placeholder reflects the shape of the
// original value (e.g., <bool>, <number>, <url>, <string>):
//
//	opts.DescribeTypes = true
//	// PORT=<number>
//	// DB_HOST=<string>
package templater
