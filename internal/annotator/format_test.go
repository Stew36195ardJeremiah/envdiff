package annotator_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/annotator"
	"github.com/user/envdiff/internal/diff"
)

func TestFormat_InlineComment(t *testing.T) {
	env := map[string]string{"APP_NAME": "myapp"}
	result := diff.Result{}
	anns := annotator.Annotate(env, result)
	opts := annotator.DefaultFormatOptions()
	out := annotator.Format(anns, opts)
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Errorf("expected key=value in output, got: %q", out)
	}
	if !strings.Contains(out, "# ok") {
		t.Errorf("expected inline comment in output, got: %q", out)
	}
}

func TestFormat_OmitOK(t *testing.T) {
	env := map[string]string{"APP_NAME": "myapp", "DB_HOST": "localhost"}
	result := diff.Result{MissingInRight: []string{"DB_HOST"}}
	anns := annotator.Annotate(env, result)
	opts := annotator.FormatOptions{InlineComment: true, OmitOK: true}
	out := annotator.Format(anns, opts)
	if strings.Contains(out, "APP_NAME") {
		t.Error("expected APP_NAME (ok) to be omitted")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST (missing) to be present")
	}
}

func TestFormat_PrefixComment(t *testing.T) {
	env := map[string]string{"SECRET": "abc"}
	result := diff.Result{
		Mismatched: []diff.Mismatch{
			{Key: "SECRET", LeftVal: "abc", RightVal: "xyz"},
		},
	}
	anns := annotator.Annotate(env, result)
	opts := annotator.FormatOptions{InlineComment: false, OmitOK: false}
	out := annotator.Format(anns, opts)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected at least 2 lines, got: %q", out)
	}
	if !strings.HasPrefix(lines[0], "#") {
		t.Errorf("expected first line to be a comment, got: %q", lines[0])
	}
	if !strings.Contains(lines[1], "SECRET=abc") {
		t.Errorf("expected second line to be key=value, got: %q", lines[1])
	}
}

func TestFormat_EmptyEnv(t *testing.T) {
	anns := annotator.Annotate(map[string]string{}, diff.Result{})
	out := annotator.Format(anns, annotator.DefaultFormatOptions())
	if out != "" {
		t.Errorf("expected empty output, got: %q", out)
	}
}
