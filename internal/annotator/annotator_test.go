package annotator_test

import (
	"testing"

	"github.com/user/envdiff/internal/annotator"
	"github.com/user/envdiff/internal/diff"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME": "myapp",
		"DB_HOST":  "localhost",
		"SECRET":   "abc123",
	}
}

func TestAnnotate_AllOK(t *testing.T) {
	result := diff.Result{}
	anns := annotator.Annotate(baseEnv(), result)
	for _, a := range anns {
		if a.Status != annotator.StatusOK {
			t.Errorf("expected OK for %q, got %q", a.Key, a.Status)
		}
	}
}

func TestAnnotate_MissingKey(t *testing.T) {
	result := diff.Result{MissingInRight: []string{"DB_HOST"}}
	anns := annotator.Annotate(baseEnv(), result)
	for _, a := range anns {
		if a.Key == "DB_HOST" {
			if a.Status != annotator.StatusMissing {
				t.Errorf("expected Missing, got %q", a.Status)
			}
			return
		}
	}
	t.Error("DB_HOST not found in annotations")
}

func TestAnnotate_MismatchedKey(t *testing.T) {
	result := diff.Result{
		Mismatched: []diff.Mismatch{
			{Key: "SECRET", LeftVal: "abc123", RightVal: "xyz789"},
		},
	}
	anns := annotator.Annotate(baseEnv(), result)
	for _, a := range anns {
		if a.Key == "SECRET" {
			if a.Status != annotator.StatusMismatched {
				t.Errorf("expected Mismatched, got %q", a.Status)
			}
			if a.Comment == "" {
				t.Error("expected non-empty comment for mismatched key")
			}
			return
		}
	}
	t.Error("SECRET not found in annotations")
}

func TestAnnotate_SortedOutput(t *testing.T) {
	env := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	anns := annotator.Annotate(env, diff.Result{})
	if len(anns) != 3 {
		t.Fatalf("expected 3 annotations, got %d", len(anns))
	}
	if anns[0].Key != "A_KEY" || anns[1].Key != "M_KEY" || anns[2].Key != "Z_KEY" {
		t.Errorf("annotations not sorted: %v", anns)
	}
}
