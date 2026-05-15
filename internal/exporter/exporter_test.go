package exporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/exporter"
)

func cleanResult() diff.Result {
	return diff.Result{}
}

func mixedResult() diff.Result {
	return diff.Result{
		MissingInRight: []string{"DB_HOST"},
		MissingInLeft:  []string{"REDIS_URL"},
		Mismatched: []diff.Mismatch{
			{Key: "APP_ENV", LeftValue: "staging", RightValue: "production"},
		},
	}
}

func TestExport_JSONClean(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.Export(&buf, cleanResult(), exporter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out exporter.JSONOutput
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !out.Clean {
		t.Error("expected Clean to be true")
	}
}

func TestExport_JSONMixed(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.Export(&buf, mixedResult(), exporter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out exporter.JSONOutput
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.Clean {
		t.Error("expected Clean to be false")
	}
	if len(out.MissingInRight) != 1 || out.MissingInRight[0] != "DB_HOST" {
		t.Errorf("unexpected MissingInRight: %v", out.MissingInRight)
	}
	if len(out.Mismatched) != 1 || out.Mismatched[0].Key != "APP_ENV" {
		t.Errorf("unexpected Mismatched: %v", out.Mismatched)
	}
}

func TestExport_TextClean(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.Export(&buf, cleanResult(), exporter.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "OK:") {
		t.Errorf("expected OK message, got: %q", buf.String())
	}
}

func TestExport_TextMixed(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.Export(&buf, mixedResult(), exporter.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "MISSING_IN_RIGHT: DB_HOST") {
		t.Errorf("expected MISSING_IN_RIGHT line, got: %q", out)
	}
	if !strings.Contains(out, "MISMATCH: APP_ENV") {
		t.Errorf("expected MISMATCH line, got: %q", out)
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, cleanResult(), exporter.Format("xml"))
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}
