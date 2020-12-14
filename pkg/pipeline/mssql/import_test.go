package mssql

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/samsamann/nc-connector/pkg/pipeline"
	"github.com/stretchr/testify/assert"
)

func TestQueryCreation(t *testing.T) {
	tests := []struct {
		subSQLQuery    string
		expectedResult string
	}{
		{
			subSQLQuery:    "select 1",
			expectedResult: "select file_data.* from (select 1) as file_data;",
		},
		{
			expectedResult: "select file_data.* from (select 1) as file_data;",
		},
	}

	for _, test := range tests {
		assert.Equal(
			t,
			test.expectedResult,
			prepareDataSQLQuery(test.subSQLQuery),
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

	processQueryResult(
		rows,
		func(c chan<- pipeline.FileData, resultMap map[string]interface{}) {
			assert.Contains(t, resultMap, "id")
		},
		nil,
	)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestExecQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := "select 1;"
	expectedErr := errors.New("can not prepare statement")
	mock.ExpectPrepare(query).WillReturnError(expectedErr)

	sut := new(mssqlSource)
	sut.db = db

	rows, err := sut.execQuery(query)
	assert.Nil(t, rows)
	assert.EqualError(t, err, expectedErr.Error())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	query = "select 1 from foo;"
	mock.ExpectPrepare(query).
		ExpectQuery().
		WillReturnRows(sqlmock.NewRows([]string{"mock"}).FromCSVString("1"))

	sut = new(mssqlSource)
	sut.db = db

	rows, err = sut.execQuery(query)
	assert.NotNil(t, rows)
	assert.True(t, rows.Next())
	assert.Nil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
