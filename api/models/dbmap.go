package models

import (
	"database/sql"

	"github.com/go-gorp/gorp"
)

func NewDBMap(db *sql.DB, dialect gorp.Dialect) (dbmap *gorp.DbMap) {
	dbmap = &gorp.DbMap{Db: db, Dialect: dialect}
	dbmap.AddTableWithName(&Pin{}, `pin`).SetKeys(true, "id")
	dbmap.AddTableWithName(&List{}, `list`).SetKeys(true, "id")
	dbmap.AddTableWithName(&User{}, `user`).SetKeys(true, "id")
	dbmap.AddTableWithName(&UserHasList{}, `user_has_list`).SetKeys(true, "id")
	dbmap.AddTableWithName(&Token{}, `token`).SetKeys(true, "id")
	dbmap.AddTableWithName(&Tag{}, `tag`).SetKeys(true, "id")
	return
}
