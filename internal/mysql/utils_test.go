package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeIdentifier(t *testing.T) {
	// Test normal identifier
	result := escapeIdentifier("users")
	assert.Equal(t, "users", result)

	// Test identifier with backtick
	result = escapeIdentifier("orders`; DROP TABLE IF EXISTS `sensitive`; --")
	assert.Equal(t, "orders``; DROP TABLE IF EXISTS ``sensitive``; --", result)

	// Test identifier with multiple backticks
	result = escapeIdentifier("table`with`multiple`backticks")
	assert.Equal(t, "table``with``multiple``backticks", result)

	// Test empty identifier
	result = escapeIdentifier("")
	assert.Equal(t, "", result)
}

func TestGetValue(t *testing.T) {
	val, err := getValue("")
	assert.NoError(t, err)
	assert.Equal(t, "''", val)

	val, err = getValue("1")
	assert.NoError(t, err)
	assert.Equal(t, "1", val)

	val, err = getValue("foo")
	assert.NoError(t, err)
	assert.Equal(t, "'foo'", val)
}

func TestEscape(t *testing.T) {
	input := string([]byte{0, '\n', '\r', '\\', '\'', '"', '\032', 'a'})
	expected := `\0\n\r\\\'\"\Za`
	result, err := escape(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
