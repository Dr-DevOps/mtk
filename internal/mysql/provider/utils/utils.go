package utils

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/skpr/mtk/internal/mysql/provider"
)

// escapeIdentifier escapes MySQL identifiers (table names, column names) by doubling backticks.
// This prevents SQL injection when identifiers are used within backtick-delimited contexts.
func escapeIdentifier(identifier string) string {
	return strings.ReplaceAll(identifier, "`", "``")
}

// QueryColumnsForTable for a given table.
func QueryColumnsForTable(database *sql.DB, table string, params provider.DumpParams) ([]string, error) {
	var rows *sql.Rows

	rows, err := database.Query(fmt.Sprintf("SELECT * FROM `%s` LIMIT 1", escapeIdentifier(table)))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for k, column := range columns {
		replacement, ok := params.SelectMap[strings.ToLower(table)][strings.ToLower(column)]
		if ok {
			columns[k] = fmt.Sprintf("%s AS `%s`", replacement, escapeIdentifier(column))
		} else {
			columns[k] = fmt.Sprintf("`%s`", escapeIdentifier(column))
		}
	}

	return columns, nil
}
