package models

import "time"

type Pin struct {
	ID        int64     `db:"id,primarykey,autoincrement" json:"id"`
	Title     string    `db:"title" json:"title"`
	URL       string    `db:"url" json:"url"`
	Tags      []Tag     `db:"-" json:"tags"`
	CreatorID int       `db:"creator_id" json:"-"`
	Creator   *User     `db:"-" json:"creator"`
	ListID    int       `db:"list_id" json:"list_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
