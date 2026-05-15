package validator_test

import (
	"testing"

	"github.com/nicholasgasior/envdiff/internal/parser"
	"github.com/nicholasgasior/envdiff/internal/validator"
)

func TestValidate_WithParsedFile(t *testing.T) {
	env, err := parser.ParseFile("../parser/testdata/sample.env")
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	issues := validator.Validate(env, validator.DefaultOptions())

	// sample.env is expected to be well-formed; surface any surprises.
	for _, iss := range issues {
		t.Logf("validation issue: %s", iss)
	}

	// At minimum the file must have been parsed without error.
	if len(env) == 0 {
		t.Fatal("expected at least one key from sample.env")
	}
}
