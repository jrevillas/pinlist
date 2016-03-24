package models

import "time"

type Pin struct {
	ID        int64
	Title     string
	URL       string
	Tags      []Tag
	CreatorID int
	ListID    int
	CreatedAt time.Time
}
