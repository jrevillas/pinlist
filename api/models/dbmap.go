package models

import (
	"database/sql"

	"gopkg.in/gorp.v1"
)

func NewDBMap(db *sql.DB, dialect gorp.Dialect) (dbmap *gorp.DbMap) {
	dbmap = &gorp.DbMap{Db: db, Dialect: dialect}
	dbmap.AddTableWithName(Pin{}, `pin`).SetKeys(true, "ID")
	dbmap.AddTableWithName(List{}, `list`).SetKeys(true, "ID")
	dbmap.AddTableWithName(User{}, `user`).SetKeys(true, "ID")
	dbmap.AddTableWithName(UserHasList{}, `user_has_list`).SetKeys(true, "ID")
	dbmap.AddTableWithName(Token{}, `token`).SetKeys(true, "ID")
	dbmap.AddTableWithName(Tag{}, `tag`).SetKeys(true, "ID")
	return
}
