package models

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
