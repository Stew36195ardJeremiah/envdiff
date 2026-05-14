package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format defines the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report writes a human-readable diff report to w.
func Report(w io.Writer, result diff.Result, leftName, rightName string) {
	if result.IsClean() {
		fmt.Fprintln(w, "✓ No differences found.")
		return
	}

	if len(result.MissingInRight) > 0 {
		fmt.Fprintf(w, "Keys present in %s but missing in %s:\n", leftName, rightName)
		for _, k := range result.MissingInRight {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(result.MissingInLeft) > 0 {
		fmt.Fprintf(w, "Keys present in %s but missing in %s:\n", rightName, leftName)
		for _, k := range result.MissingInLeft {
			fmt.Fprintf(w, "  + %s\n", k)
		}
	}

	if len(result.Mismatched) > 0 {
		fmt.Fprintln(w, "Mismatched values:")
		for _, m := range result.Mismatched {
			fmt.Fprintf(w, "  ~ %s\n", m.Key)
			fmt.Fprintf(w, "      %s: %s\n", leftName, maskValue(m.LeftVal))
			fmt.Fprintf(w, "      %s: %s\n", rightName, maskValue(m.RightVal))
		}
	}
}

// Summary returns a one-line summary string of the diff result.
func Summary(result diff.Result) string {
	if result.IsClean() {
		return "clean"
	}
	parts := []string{}
	if n := len(result.MissingInRight); n > 0 {
		parts = append(parts, fmt.Sprintf("%d missing right", n))
	}
	if n := len(result.MissingInLeft); n > 0 {
		parts = append(parts, fmt.Sprintf("%d missing left", n))
	}
	if n := len(result.Mismatched); n > 0 {
		parts = append(parts, fmt.Sprintf("%d mismatched", n))
	}
	return strings.Join(parts, ", ")
}

// maskValue replaces the value with asterisks to avoid leaking secrets.
func maskValue(v string) string {
	if len(v) == 0 {
		return "(empty)"
	}
	return strings.Repeat("*", len(v))
}
