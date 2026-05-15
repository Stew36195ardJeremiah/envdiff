package merger_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/merger"
)

func TestMerge_NoConflicts(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"BAZ": "3"}

	res := merger.Merge([]map[string]string{a, b}, merger.Options{})

	if len(res.Env) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(res.Env))
	}
	if res.Env["FOO"] != "1" || res.Env["BAR"] != "2" || res.Env["BAZ"] != "3" {
		t.Error("unexpected values in merged env")
	}
}

func TestMerge_StrategyFirst(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}

	res := merger.Merge([]map[string]string{a, b}, merger.Options{Strategy: merger.StrategyFirst})

	if res.Env["KEY"] != "first" {
		t.Errorf("expected 'first', got %q", res.Env["KEY"])
	}
}

func TestMerge_StrategyLast(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}

	res := merger.Merge([]map[string]string{a, b}, merger.Options{Strategy: merger.StrategyLast})

	if res.Env["KEY"] != "second" {
		t.Errorf("expected 'second', got %q", res.Env["KEY"])
	}
}

func TestMerge_TrackOrigin(t *testing.T) {
	a := map[string]string{"A": "1"}
	b := map[string]string{"B": "2", "A": "override"}

	res := merger.Merge([]map[string]string{a, b}, merger.Options{
		Strategy:    merger.StrategyFirst,
		TrackOrigin: true,
	})

	if res.Origin == nil {
		t.Fatal("expected Origin map to be populated")
	}
	if res.Origin["A"] != 0 {
		t.Errorf("expected A from source 0, got %d", res.Origin["A"])
	}
	if res.Origin["B"] != 1 {
		t.Errorf("expected B from source 1, got %d", res.Origin["B"])
	}
}

func TestMerge_TrackOriginDisabled(t *testing.T) {
	a := map[string]string{"X": "1"}
	res := merger.Merge([]map[string]string{a}, merger.Options{TrackOrigin: false})

	if res.Origin != nil {
		t.Error("expected Origin to be nil when TrackOrigin is false")
	}
}

func TestMerge_EmptySources(t *testing.T) {
	res := merger.Merge(nil, merger.Options{})
	if len(res.Env) != 0 {
		t.Errorf("expected empty env, got %d keys", len(res.Env))
	}
}
