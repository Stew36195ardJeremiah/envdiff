package parser_test

import (
	"bufio"
	"strings"
	"testing"

	"github.com/envdiff/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func scan(s string) *bufio.Scanner {
	return bufio.NewScanner(strings.NewReader(s))
}

func TestParseReader_BasicKeyValue(t *testing.T) {
	input := `
DB_HOST=localhost
DB_PORT=5432
`
	env, err := parser.ParseReader(scan(input))
	require.NoError(t, err)
	assert.Equal(t, "localhost", env["DB_HOST"])
	assert.Equal(t, "5432", env["DB_PORT"])
}

func TestParseReader_SkipsCommentsAndBlanks(t *testing.T) {
	input := `
# this is a comment

APP_ENV=production
`
	env, err := parser.ParseReader(scan(input))
	require.NoError(t, err)
	assert.Len(t, env, 1)
	assert.Equal(t, "production", env["APP_ENV"])
}

func TestParseReader_StripQuotes(t *testing.T) {
	input := `
SECRET="my secret value"
TOKEN='bearer-token'
`
	env, err := parser.ParseReader(scan(input))
	require.NoError(t, err)
	assert.Equal(t, "my secret value", env["SECRET"])
	assert.Equal(t, "bearer-token", env["TOKEN"])
}

func TestParseReader_ExportPrefix(t *testing.T) {
	input := `export PATH=/usr/local/bin`
	env, err := parser.ParseReader(scan(input))
	require.NoError(t, err)
	assert.Equal(t, "/usr/local/bin", env["PATH"])
}

func TestParseReader_MalformedLine(t *testing.T) {
	input := `BADLINE`
	_, err := parser.ParseReader(scan(input))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "malformed line")
}

func TestParseFile_FileNotFound(t *testing.T) {
	_, err := parser.ParseFile("/nonexistent/.env")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "opening file")
}
