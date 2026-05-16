package sorter_test

import (
	"testing"

	"github.com/nicholasgasior/envdiff/internal/sorter"
)

func TestSort_Alphabetical(t *testing.T) {
	env := map[string]string{
		"ZEBRA": "1",
		"APPLE": "2",
		"MANGO": "3",
	}
	opts := sorter.DefaultOptions()
	_, keys := sorter.Sort(env, opts)

	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("position %d: got %q, want %q", i, k, expected[i])
		}
	}
}

func TestSort_NoAlphabetical(t *testing.T) {
	env := map[string]string{"B": "1", "A": "2"}
	opts := sorter.Options{Alphabetical: false, GroupByPrefix: false, PrefixSeparator: "_"}
	_, keys := sorter.Sort(env, opts)
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}

func TestSort_GroupByPrefix(t *testing.T) {
	env := map[string]string{
		"DB_HOST":     "localhost",
		"AWS_KEY":     "abc",
		"DB_PORT":     "5432",
		"AWS_SECRET":  "xyz",
		"APP_VERSION": "1.0",
	}
	opts := sorter.Options{Alphabetical: true, GroupByPrefix: true, PrefixSeparator: "_"}
	_, keys := sorter.Sort(env, opts)

	// Expect groups: APP_, AWS_, DB_
	if keys[0] != "APP_VERSION" {
		t.Errorf("expected APP_VERSION first, got %q", keys[0])
	}
	// AWS group should come before DB group
	awsStart := indexOf(keys, "AWS_KEY")
	dbStart := indexOf(keys, "DB_HOST")
	if awsStart == -1 || dbStart == -1 {
		t.Fatal("missing expected keys")
	}
	if awsStart > dbStart {
		t.Errorf("expected AWS group before DB group")
	}
}

func TestGroupKeys_SplitsByPrefix(t *testing.T) {
	keys := []string{"DB_HOST", "DB_PORT", "APP_NAME", "NOPREFIXKEY"}
	groups := sorter.GroupKeys(keys, "_")

	if len(groups["DB"]) != 2 {
		t.Errorf("expected 2 DB keys, got %d", len(groups["DB"]))
	}
	if len(groups["APP"]) != 1 {
		t.Errorf("expected 1 APP key, got %d", len(groups["APP"]))
	}
	if len(groups[""]) != 1 {
		t.Errorf("expected 1 unprefixed key, got %d", len(groups[""]))
	}
}

func indexOf(slice []string, val string) int {
	for i, v := range slice {
		if v == val {
			return i
		}
	}
	return -1
}
