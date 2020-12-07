package mssql

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestQueryCreation(t *testing.T) {
	tests := []struct {
		subSQLQuery    string
		idColumnName   string
		expectedResult string
	}{
		{
			subSQLQuery:    "select 1",
			idColumnName:   "id",
			expectedResult: "select file_data.id, file_data.* from (select 1) as file_data;",
		},
		{
			idColumnName:   "id_col",
			expectedResult: "select file_data.id_col, file_data.* from (select 1 as id_col) as file_data;",
		},
	}

	for _, test := range tests {
		assert.Equal(
			t,
			test.expectedResult,
			prepareDataSQLQuery(test.subSQLQuery, test.idColumnName),
		)
	}
}

func TestProcessQueryResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{"id"}
	mock.ExpectQuery("select 1 as id from;").
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1\nfoo")).
		RowsWillBeClosed()
	rows, _ := db.Query("select 1 as id from;")

	resultMap := processQueryResult(rows)
	if assert.Len(t, resultMap, 2) {
		assert.Contains(t, resultMap[0], "id")
		assert.Equal(t, resultMap[0]["id"], 1)
		assert.Contains(t, resultMap[1], "id")
		assert.Equal(t, resultMap[1]["id"], "foo")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
