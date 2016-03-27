package models

import "gopkg.in/gorp.v1"

// Tag is a word or set of words that allow better indexing
// of a pin.
type Tag struct {
	ID    int64  `db:"id" json:"id"`
	PinID int64  `db:"pin_id" json:"-"`
	Name  string `db:"name" json:"name"`
}

// NewTag creates a new tag with the given name.
func NewTag(name string) *Tag {
	return &Tag{Name: name}
}

// TagStore is a service to access the tag database table.
type TagStore struct {
	*gorp.DbMap
}

// TagCount represents a tag and the number of appearances per user.
type TagCount struct {
	Name  string `db:"name" json:"name"`
	Count int    `db:"count" json:"count"`
}

const allTagsQuery = `SELECT t.name, COUNT(*) as count FROM tag t
INNER JOIN pin p ON t.pin_id = p.id
LEFT JOIN user_has_list l ON l.list_id = p.list_id
WHERE (p.creator_id = :user OR l.user_id = :user)
GROUP BY t.name ORDER BY t.name ASC LIMIT :limit OFFSET :offset`

// All retrieves the all tags of the user with their count.
func (s TagStore) All(user int64, limit, offset int) ([]*TagCount, error) {
	var tags []*TagCount
	_, err := s.Select(&tags, allTagsQuery, map[string]interface{}{
		"user":   user,
		"limit":  limit,
		"offset": offset,
	})
	if err != nil {
		return nil, err
	}
	return tags, nil
}

const countTagsQuery = `SELECT COUNT(*) FROM (SELECT t.name FROM tag t
INNER JOIN pin p ON t.pin_id = p.id
LEFT JOIN user_has_list l ON l.list_id = p.list_id
WHERE (p.creator_id = :user OR l.user_id = :user)
GROUP BY t.name)`

// Count retrieves the number of unique tags an user has.
func (s TagStore) Count(user int64) (int64, error) {
	return s.SelectInt(countTagsQuery, map[string]interface{}{
		"user": user,
	})
}
