package utils

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/skpr/mtk/internal/mysql/mock"
	"github.com/skpr/mtk/internal/mysql/provider"
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

func TestMySQLGetColumnsForSelect(t *testing.T) {
	db, mock := mock.GetDB(t)
	mock.ExpectQuery("SELECT \\* FROM `table` LIMIT 1").WillReturnRows(
		sqlmock.NewRows([]string{"col1", "col2", "col3"}).AddRow("a", "b", "c"))
	columns, err := QueryColumnsForTable(db, "table", provider.DumpParams{
		SelectMap: map[string]map[string]string{"table": {"col2": "NOW()"}},
	})
	assert.Nil(t, err)
	assert.Equal(t, []string{"`col1`", "NOW() AS `col2`", "`col3`"}, columns)
}

func TestMySQLGetColumnsForSelectWithBackticks(t *testing.T) {
	db, mock := mock.GetDB(t)
	mock.ExpectQuery("SELECT \\* FROM `table` LIMIT 1").WillReturnRows(
		sqlmock.NewRows([]string{"col1", "col`2", "col3"}).AddRow("a", "b", "c"))
	columns, err := QueryColumnsForTable(db, "table", provider.DumpParams{})
	assert.Nil(t, err)
	assert.Equal(t, []string{"`col1`", "`col``2`", "`col3`"}, columns)
}

func TestMySQLGetColumnsForSelectHandlingErrorWhenQuerying(t *testing.T) {
	db, mock := mock.GetDB(t)
	e := errors.New("broken")
	mock.ExpectQuery("SELECT \\* FROM `table` LIMIT 1").WillReturnError(e)
	columns, err := QueryColumnsForTable(db, "table", provider.DumpParams{
		SelectMap: map[string]map[string]string{"table": {"col2": "NOW()"}},
	})
	assert.Equal(t, err, e)
	assert.Empty(t, columns)
}
