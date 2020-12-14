package mssql

import (
	"database/sql"
	"fmt"
	"net/url"
	"strconv"

	mssql "github.com/denisenkom/go-mssqldb"
	pip "github.com/samsamann/nc-connector/pkg/pipeline"
)

func init() {
	pip.RegisterFileImporter("mssqlSource", NewMssqlSource)
}

type processResultFunc func(chan<- pip.FileData, map[string]interface{})

type mssqlSource struct {
	query string
	cURL  *url.URL
	db    *sql.DB
}

// NewMssqlSource returns a new instance of mssqlSource.
func NewMssqlSource(config map[string]interface{}) pip.FileImporter {
	return &mssqlSource{query: config["query"].(string), cURL: newURL(config)}
}

func (ms *mssqlSource) Connect() error {
	connector, err := mssql.NewConnector(ms.cURL.String())
	if err != nil {
		return err
	}
	ms.db = sql.OpenDB(connector)
	return nil
}

func (ms mssqlSource) Import(context pip.ImportContext, channel chan<- pip.FileData) error {
	rows, err := ms.execQuery(prepareDataSQLQuery(ms.query))
	if err != nil {
		return err
	}
	processQueryResult(
		rows,
		createFile,
		channel,
	)
	return nil
}

func (ms mssqlSource) execQuery(query string) (*sql.Rows, error) {
	err := ms.db.Ping()
	if err != nil {
		return nil, err
	}
	stmt, err := ms.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func newURL(config map[string]interface{}) *url.URL {
	return &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(config["username"].(string), config["password"].(string)),
		Host:   fmt.Sprintf("%s/%s", config["hostname"], config["instance"]),
	}
}

func prepareDataSQLQuery(query string) string {
	return prepareSQLQuery(
		"select file_data.* from (%s) as file_data;",
		query,
	)
}

func prepareCountSQLQuery(query, idColumnName string) string {
	return prepareSQLQuery(
		"select count(*) as count from (%s) as file_data;",
		query,
	)
}

func prepareSQLQuery(mainQuery, subQuery string) string {
	if subQuery == "" {
		subQuery = fmt.Sprintf("select 1")
	}
	return fmt.Sprintf(mainQuery, subQuery)
}

func processQueryResult(rows *sql.Rows, processResult processResultFunc, channel chan<- pip.FileData) {
	defer rows.Close()

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
			// TODO: log and/or continue
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
		processResult(channel, m)

		hasRow = rows.Next()
	}
}

func createFile(channel chan<- pip.FileData, result map[string]interface{}) {
	channel <- &pip.File{Name: result["name"].(string)}
}
