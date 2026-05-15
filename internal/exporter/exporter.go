// Package exporter provides functionality to export diff results
// to various output formats such as JSON and plain text.
package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the output format for exporting diff results.
type Format string

const (
	FormatJSON  Format = "json"
	FormatText  Format = "text"
)

// JSONOutput is the structured representation of a diff result for JSON export.
type JSONOutput struct {
	MissingInRight []string            `json:"missing_in_right,omitempty"`
	MissingInLeft  []string            `json:"missing_in_left,omitempty"`
	Mismatched     []MismatchedEntry   `json:"mismatched,omitempty"`
	Clean          bool                `json:"clean"`
}

// MismatchedEntry holds a key and its two differing values.
type MismatchedEntry struct {
	Key        string `json:"key"`
	LeftValue  string `json:"left_value"`
	RightValue string `json:"right_value"`
}

// Export writes the diff result to w in the specified format.
func Export(w io.Writer, result diff.Result, format Format) error {
	switch format {
	case FormatJSON:
		return exportJSON(w, result)
	case FormatText:
		return exportText(w, result)
	default:
		return fmt.Errorf("unsupported export format: %q", format)
	}
}

func exportJSON(w io.Writer, result diff.Result) error {
	out := JSONOutput{
		Clean:          result.IsClean(),
		MissingInRight: result.MissingInRight,
		MissingInLeft:  result.MissingInLeft,
	}
	for _, m := range result.Mismatched {
		out.Mismatched = append(out.Mismatched, MismatchedEntry{
			Key:        m.Key,
			LeftValue:  m.LeftValue,
			RightValue: m.RightValue,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func exportText(w io.Writer, result diff.Result) error {
	var sb strings.Builder
	for _, k := range result.MissingInRight {
		sb.WriteString(fmt.Sprintf("MISSING_IN_RIGHT: %s\n", k))
	}
	for _, k := range result.MissingInLeft {
		sb.WriteString(fmt.Sprintf("MISSING_IN_LEFT: %s\n", k))
	}
	for _, m := range result.Mismatched {
		sb.WriteString(fmt.Sprintf("MISMATCH: %s (%q != %q)\n", m.Key, m.LeftValue, m.RightValue))
	}
	if result.IsClean() {
		sb.WriteString("OK: no differences found\n")
	}
	_, err := fmt.Fprint(w, sb.String())
	return err
}
