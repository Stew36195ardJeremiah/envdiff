package templater

import (
	"strings"
	"testing"
)

func TestGenerate_DefaultOptions(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	opts := DefaultOptions()
	out := Generate(env, opts)

	if !strings.Contains(out, "DB_HOST=<value>") {
		t.Errorf("expected DB_HOST=<value>, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_PORT=<value>") {
		t.Errorf("expected DB_PORT=<value>, got:\n%s", out)
	}
}

func TestGenerate_CustomPlaceholder(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "abc123"}
	opts := DefaultOptions()
	opts.Placeholder = "REPLACE_ME"
	out := Generate(env, opts)

	if !strings.Contains(out, "SECRET_KEY=REPLACE_ME") {
		t.Errorf("expected SECRET_KEY=REPLACE_ME, got:\n%s", out)
	}
}

func TestGenerate_DescribeTypes(t *testing.T) {
	env := map[string]string{
		"ENABLED":  "true",
		"PORT":     "8080",
		"BASE_URL": "https://example.com",
		"NAME":     "myapp",
	}
	opts := DefaultOptions()
	opts.DescribeTypes = true
	out := Generate(env, opts)

	cases := map[string]string{
		"ENABLED":  "<bool>",
		"PORT":     "<number>",
		"BASE_URL": "<url>",
		"NAME":     "<string>",
	}
	for key, want := range cases {
		line := key + "=" + want
		if !strings.Contains(out, line) {
			t.Errorf("expected %q in output, got:\n%s", line, out)
		}
	}
}

func TestGenerate_SortedKeys(t *testing.T) {
	env := map[string]string{
		"ZEBRA": "z",
		"ALPHA": "a",
		"MANGO": "m",
	}
	opts := DefaultOptions()
	opts.IncludeComments = false
	out := Generate(env, opts)

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "ALPHA") {
		t.Errorf("expected first key ALPHA, got %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "ZEBRA") {
		t.Errorf("expected last key ZEBRA, got %s", lines[2])
	}
}

func TestGenerate_HeaderComment(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	opts := DefaultOptions()
	opts.IncludeComments = true
	out := Generate(env, opts)

	if !strings.HasPrefix(out, "#") {
		t.Errorf("expected output to start with comment, got:\n%s", out)
	}
}

func TestGenerate_NoComment(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	opts := DefaultOptions()
	opts.IncludeComments = false
	out := Generate(env, opts)

	if strings.HasPrefix(out, "#") {
		t.Errorf("expected no comment header, got:\n%s", out)
	}
}

func TestIsNumeric(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"123", true},
		{"3.14", true},
		{"-42", true},
		{"abc", false},
		{"", false},
		{"1.2.3", false},
	}
	for _, tc := range cases {
		got := isNumeric(tc.input)
		if got != tc.want {
			t.Errorf("isNumeric(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}
