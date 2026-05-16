package profiler_test

import (
	"testing"

	"github.com/user/envdiff/internal/profiler"
)

func sampleEnv() map[string]string {
	return map[string]string{
		"DB_HOST":       "localhost",
		"DB_PORT":       "5432",
		"DB_PASSWORD":   "secret",
		"AWS_ACCESS_KEY": "AKIA123",
		"AWS_SECRET_KEY": "topsecret",
		"APP_NAME":      "envdiff",
		"APP_ENV":       "",
		"STANDALONE":    "yes",
	}
}

func TestAnalyze_TotalKeys(t *testing.T) {
	p := profiler.Analyze(sampleEnv())
	if p.TotalKeys != 8 {
		t.Errorf("expected 8 total keys, got %d", p.TotalKeys)
	}
}

func TestAnalyze_EmptyValues(t *testing.T) {
	p := profiler.Analyze(sampleEnv())
	if p.EmptyValues != 1 {
		t.Errorf("expected 1 empty value, got %d", p.EmptyValues)
	}
}

func TestAnalyze_SensitiveKeys(t *testing.T) {
	p := profiler.Analyze(sampleEnv())
	// DB_PASSWORD, AWS_ACCESS_KEY, AWS_SECRET_KEY should be flagged
	if p.SensitiveKeys < 2 {
		t.Errorf("expected at least 2 sensitive keys, got %d", p.SensitiveKeys)
	}
}

func TestAnalyze_PrefixCounts(t *testing.T) {
	p := profiler.Analyze(sampleEnv())
	if p.PrefixCounts["DB"] != 3 {
		t.Errorf("expected DB prefix count 3, got %d", p.PrefixCounts["DB"])
	}
	if p.PrefixCounts["AWS"] != 2 {
		t.Errorf("expected AWS prefix count 2, got %d", p.PrefixCounts["AWS"])
	}
	if _, ok := p.PrefixCounts["STANDALONE"]; ok {
		t.Error("STANDALONE should not appear in prefix counts")
	}
}

func TestAnalyze_EmptyMap(t *testing.T) {
	p := profiler.Analyze(map[string]string{})
	if p.TotalKeys != 0 || p.EmptyValues != 0 || p.SensitiveKeys != 0 {
		t.Error("expected all zero counts for empty map")
	}
}

func TestTopPrefixes_Order(t *testing.T) {
	p := profiler.Analyze(sampleEnv())
	top := profiler.TopPrefixes(p, 2)
	if len(top) != 2 {
		t.Fatalf("expected 2 top prefixes, got %d", len(top))
	}
	if top[0] != "DB" {
		t.Errorf("expected DB as top prefix, got %s", top[0])
	}
}

func TestTopPrefixes_LimitRespected(t *testing.T) {
	p := profiler.Analyze(sampleEnv())
	top := profiler.TopPrefixes(p, 1)
	if len(top) != 1 {
		t.Errorf("expected 1 result, got %d", len(top))
	}
}
