package diff_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/diff"
)

func TestCompare_Clean(t *testing.T) {
	left := map[string]string{"FOO": "bar", "BAZ": "qux"}
	right := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := diff.Compare(left, right)
	if !result.IsClean() {
		t.Errorf("expected clean result, got %+v", result)
	}
}

func TestCompare_MissingInRight(t *testing.T) {
	left := map[string]string{"FOO": "bar", "ONLY_LEFT": "yes"}
	right := map[string]string{"FOO": "bar"}

	result := diff.Compare(left, right)
	if len(result.MissingInRight) != 1 || result.MissingInRight[0] != "ONLY_LEFT" {
		t.Errorf("unexpected MissingInRight: %v", result.MissingInRight)
	}
	if len(result.MissingInLeft) != 0 {
		t.Errorf("expected no MissingInLeft, got %v", result.MissingInLeft)
	}
}

func TestCompare_MissingInLeft(t *testing.T) {
	left := map[string]string{"FOO": "bar"}
	right := map[string]string{"FOO": "bar", "ONLY_RIGHT": "yes"}

	result := diff.Compare(left, right)
	if len(result.MissingInLeft) != 1 || result.MissingInLeft[0] != "ONLY_RIGHT" {
		t.Errorf("unexpected MissingInLeft: %v", result.MissingInLeft)
	}
}

func TestCompare_Mismatched(t *testing.T) {
	left := map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
	right := map[string]string{"DB_HOST": "prod.db.internal", "PORT": "5432"}

	result := diff.Compare(left, right)
	if len(result.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(result.Mismatched))
	}
	m := result.Mismatched[0]
	if m.Key != "DB_HOST" || m.LeftValue != "localhost" || m.RightValue != "prod.db.internal" {
		t.Errorf("unexpected mismatch entry: %+v", m)
	}
}

func TestCompare_SortedOutput(t *testing.T) {
	left := map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"}
	right := map[string]string{}

	result := diff.Compare(left, right)
	keys := result.MissingInRight
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("MissingInRight not sorted: %v", keys)
			break
		}
	}
}

func TestIsClean_False(t *testing.T) {
	r := diff.Result{MissingInRight: []string{"X"}}
	if r.IsClean() {
		t.Error("expected IsClean to return false")
	}
}
