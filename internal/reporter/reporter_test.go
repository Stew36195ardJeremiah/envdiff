package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/reporter"
)

func TestReport_Clean(t *testing.T) {
	var buf bytes.Buffer
	reporter.Report(&buf, diff.Result{}, "staging", "production")
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected clean message, got: %s", buf.String())
	}
}

func TestReport_MissingInRight(t *testing.T) {
	result := diff.Result{MissingInRight: []string{"SECRET_KEY"}}
	var buf bytes.Buffer
	reporter.Report(&buf, result, "staging", "production")
	out := buf.String()
	if !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("expected SECRET_KEY in output, got: %s", out)
	}
	if !strings.Contains(out, "staging") {
		t.Errorf("expected staging label in output")
	}
}

func TestReport_Mismatched(t *testing.T) {
	result := diff.Result{
		Mismatched: []diff.Mismatch{{Key: "DB_HOST", LeftVal: "localhost", RightVal: "db.prod.example.com"}},
	}
	var buf bytes.Buffer
	reporter.Report(&buf, result, "staging", "production")
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output")
	}
	if strings.Contains(out, "localhost") {
		t.Errorf("value should be masked, not shown in plain text")
	}
}

func TestSummary_Clean(t *testing.T) {
	if got := reporter.Summary(diff.Result{}); got != "clean" {
		t.Errorf("expected 'clean', got %q", got)
	}
}

func TestSummary_Mixed(t *testing.T) {
	result := diff.Result{
		MissingInRight: []string{"A", "B"},
		Mismatched:     []diff.Mismatch{{Key: "C"}},
	}
	got := reporter.Summary(result)
	if !strings.Contains(got, "2 missing right") {
		t.Errorf("unexpected summary: %s", got)
	}
	if !strings.Contains(got, "1 mismatched") {
		t.Errorf("unexpected summary: %s", got)
	}
}
