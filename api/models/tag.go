package models

type Tag struct {
	ID    int64  `db:"id,primarykey,autoincrement" json:"id"`
	PinID int64  `db:"pin_id" json:"pin_id"`
	Name  string `db:"name" json:"name"`
}
