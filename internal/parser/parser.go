// Package parser provides functionality to parse .env files into key-value maps.
package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents the parsed key-value pairs from a .env file.
type EnvMap map[string]string

// ParseFile reads a .env file at the given path and returns an EnvMap.
// It skips blank lines and comment lines (starting with '#').
// It returns an error if the file cannot be opened or contains malformed lines.
func ParseFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parser: opening file %q: %w", path, err)
	}
	defer f.Close()

	return ParseReader(bufio.NewScanner(f))
}

// ParseReader parses .env content from a bufio.Scanner.
func ParseReader(scanner *bufio.Scanner) (EnvMap, error) {
	env := make(EnvMap)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// skip blank lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// strip optional inline export keyword
		line = strings.TrimPrefix(line, "export ")

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("parser: malformed line %d: %q", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = stripQuotes(value)

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parser: scanning: %w", err)
	}

	return env, nil
}

// stripQuotes removes surrounding single or double quotes from a value.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
