package parser_test

import (
	"testing"

	"github.com/envdiff/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFile_SampleEnv(t *testing.T) {
	env, err := parser.ParseFile("testdata/sample.env")
	require.NoError(t, err)

	expected := parser.EnvMap{
		"APP_NAME":  "envdiff",
		"APP_ENV":   "staging",
		"DB_HOST":   "db.staging.internal",
		"DB_PORT":   "5432",
		"DB_USER":   "admin",
		"DB_PASS":   "s3cr3t!pass",
		"REDIS_URL": "redis://localhost:6379",
		"LOG_LEVEL": "debug",
	}

	assert.Equal(t, expected, env)
}

func TestParseFile_KeyCount(t *testing.T) {
	env, err := parser.ParseFile("testdata/sample.env")
	require.NoError(t, err)
	assert.Len(t, env, 8, "expected 8 keys, comments and blanks should be excluded")
}
