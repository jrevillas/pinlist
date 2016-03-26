package models

import (
	"database/sql"
	"io/ioutil"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/check.v1"
	"gopkg.in/gorp.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

// Copied here again because of cyclic imports :(
func createTestDatabase() (string, *gorp.DbMap, error) {
	f, err := ioutil.TempFile("", "db_test_")
	if err != nil {
		return "", nil, err
	}
	f.Close()

	db, err := sql.Open("sqlite3", f.Name())
	if err != nil {
		return "", nil, err
	}

	dbmap := NewDBMap(db, gorp.SqliteDialect{})
	if err := dbmap.CreateTables(); err != nil {
		return "", nil, err
	}

	return f.Name(), dbmap, nil
}
