package annotator

import (
	"fmt"
	"strings"
)

// FormatOptions controls how annotations are rendered.
type FormatOptions struct {
	// InlineComment places the comment on the same line as the key.
	InlineComment bool
	// OmitOK skips keys with StatusOK from the output.
	OmitOK bool
}

// DefaultFormatOptions returns sensible formatting defaults.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{
		InlineComment: true,
		OmitOK:        false,
	}
}

// Format renders a slice of Annotations as a string in .env style.
func Format(annotations []Annotation, opts FormatOptions) string {
	var sb strings.Builder
	for _, a := range annotations {
		if opts.OmitOK && a.Status == StatusOK {
			continue
		}
		line := fmt.Sprintf("%s=%s", a.Key, a.Value)
		if opts.InlineComment && a.Comment != "" {
			line = fmt.Sprintf("%-30s %s", line, a.Comment)
		} else if !opts.InlineComment && a.Comment != "" {
			sb.WriteString(a.Comment + "\n")
		}
		sb.WriteString(line + "\n")
	}
	return sb.String()
}
