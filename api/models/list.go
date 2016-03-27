package models

import "gopkg.in/gorp.v1"

// List is a collection of pins under a certain name, topic, etc.
// Lists can be shared between many users but have only one
// creator.
type List struct {
	ID     int64          `db:"id" json:"id"`
	Name   string         `db:"name" json:"name"`
	Pins   int            `db:"pins" json:"pins"`
	Public bool           `db:"public" json:"public"`
	Users  []*UserHasList `db:"-" json:"users"`
}

// ListRole defines the kind of role an user has for a list.
type ListRole byte

const (
	// ListOwner gives the user permission to do anything on a list.
	ListOwner ListRole = 1 << iota
	// ListAdder gives the user permission to add pins to a list.
	ListAdder
	// ListGuest gives the user permission to see a list.
	ListGuest
)

// UserHasList defines the relationship between an user and a
// list, which also states the role of the user for that list.
type UserHasList struct {
	ID     int64    `id:"id" json:"-"`
	Role   ListRole `db:"role" json:"role"`
	ListID int64    `db:"list_id" json:"-"`
	UserID int64    `db:"user_id" json:"-"`
	User   *User    `db:"-" json:"user"`
}

// ListStore is the service to execute operations about lists
// that require the database.
type ListStore struct {
	*gorp.DbMap
}

const userHasAccessQuery = `SELECT COUNT(*) FROM user_has_list
WHERE user_id = :user AND list_id = :list`

// UserHasAccess reports whether the user has access to the given list.
func (s ListStore) UserHasAccess(user *User, list int64) (bool, error) {
	l, err := s.ByID(list)
	if err != nil {
		return false, err
	}

	if l.Public {
		return true, nil
	}

	// if the user is nil at this point, they can't access
	if user == nil {
		return false, nil
	}

	n, err := s.SelectInt(userHasAccessQuery, map[string]interface{}{
		"list": list,
		"user": user.ID,
	})
	if err != nil {
		return false, err
	}

	return n > 0, nil
}

// ByID returns the list with the given id.
func (s ListStore) ByID(id int64) (*List, error) {
	l, err := s.Get(List{}, id)
	if err != nil || l == nil {
		return nil, err
	}
	return l.(*List), nil
}
