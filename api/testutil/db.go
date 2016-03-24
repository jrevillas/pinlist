package testutil

import (
	"database/sql"
	"io/ioutil"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mvader/pinlist/api/models"
	"gopkg.in/gorp.v1"
)

func CreateTestDatabase() (string, *gorp.DbMap, error) {
	f, err := ioutil.TempFile("", "db_test_")
	if err != nil {
		return "", nil, err
	}
	f.Close()

	db, err := sql.Open("sqlite3", f.Name())
	if err != nil {
		return "", nil, err
	}

	dbmap := models.NewDBMap(db, gorp.SqliteDialect{})
	if err := dbmap.CreateTables(); err != nil {
		return "", nil, err
	}

	return f.Name(), dbmap, nil
}
