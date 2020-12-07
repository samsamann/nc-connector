package mssql

import (
	"database/sql"
	"fmt"
	"strconv"

	pip "github.com/samsamann/nc-connector/pkg/pipeline"
)

type mssqlSource struct {
}

// NewMssqlSource returns a new instance of mssqlSource.
func NewMssqlSource() pip.FileImporter {
	return new(mssqlSource)
}

func (ms mssqlSource) Import() <-chan pip.File {
	channel := make(chan pip.File)
	defer close(channel)
	processQueryResult(nil)
	return channel
}

func prepareDataSQLQuery(query, idColumnName string) string {
	return prepareSQLQuery(
		"select file_data.%s, file_data.* from (%s) as file_data;",
		query,
		idColumnName,
	)
}

func prepareCountSQLQuery(query, idColumnName string) string {
	return prepareSQLQuery(
		"select count(file_data.%s) from (%s) as file_data;",
		query,
		idColumnName,
	)
}

func prepareSQLQuery(mainQuery, subQuery, idColumnName string) string {
	if subQuery == "" {
		subQuery = fmt.Sprintf("select 1 as %s", idColumnName)
	}
	return fmt.Sprintf(mainQuery, idColumnName, subQuery)
}

func processQueryResult(rows *sql.Rows) []map[string]interface{} {
	defer rows.Close()

	result := make([]map[string]interface{}, 0)

	hasRow := rows.Next()
	dbCols, _ := rows.Columns()
	for hasRow {
		// Prepare
		columns := make([]interface{}, len(dbCols))
		columnPointers := make([]interface{}, len(dbCols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			// TODO: log and or continue
		}
		m := make(map[string]interface{})
		for i, colName := range dbCols {
			switch val := (*columnPointers[i].(*interface{})).(type) {
			case []uint8:
				if cint, err := strconv.Atoi(string(val)); err == nil {
					m[colName] = cint
				} else {
					m[colName] = string(val)
				}
			default:
				m[colName] = val
			}
		}
		result = append(result, m)

		hasRow = rows.Next()
	}
	return result
}
