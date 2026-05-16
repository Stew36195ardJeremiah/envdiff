package templater_test

import (
	"strings"
	"testing"

	"github.com/nicholas-fedor/envdiff/internal/parser"
	"github.com/nicholas-fedor/envdiff/internal/templater"
)

func TestGenerate_WithParsedFile(t *testing.T) {
	env, err := parser.ParseFile("../parser/testdata/sample.env")
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	opts := templater.DefaultOptions()
	out := templater.Generate(env, opts)

	if out == "" {
		t.Fatal("expected non-empty template output")
	}

	// Every key should map to the default placeholder
	for key := range env {
		expected := key + "=<value>"
		if !strings.Contains(out, expected) {
			t.Errorf("expected line %q in template output", expected)
		}
	}
}

func TestGenerate_WithParsedFile_DescribeTypes(t *testing.T) {
	env, err := parser.ParseFile("../parser/testdata/sample.env")
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	opts := templater.DefaultOptions()
	opts.DescribeTypes = true
	out := templater.Generate(env, opts)

	if out == "" {
		t.Fatal("expected non-empty template output")
	}

	// No original values should appear in the output
	for _, val := range env {
		if strings.Contains(out, "="+val+"\n") {
			t.Errorf("original value %q leaked into template", val)
		}
	}
}
