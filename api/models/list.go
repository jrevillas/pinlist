package models

import (
	"fmt"

	"gopkg.in/gorp.v1"
)

// List is a collection of pins under a certain name, topic, etc.
// Lists can be shared between many users but have only one
// creator.
type List struct {
	ID          int64          `db:"id" json:"id"`
	Name        string         `db:"name" json:"name"`
	Description string         `db:"description" json:"description"`
	Pins        int            `db:"pins" json:"pins"`
	Public      bool           `db:"public" json:"public"`
	Users       []*UserHasList `db:"-" json:"users"`
}

// NewList creates a new list from the name and description given.
func NewList(name, description string) *List {
	return &List{
		Name:        name,
		Description: description,
	}
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

var userIsOwnerQuery = fmt.Sprintf(`SELECT COUNT(id) FROM user_has_list
WHERE user_id = :user AND list_id = :list AND role = %d`, ListOwner)

// UserIsOwner reports whether the user is an owner of the list.
func (s ListStore) UserIsOwner(user *User, ID int64) (bool, error) {
	n, err := s.SelectInt(userIsOwnerQuery, map[string]interface{}{
		"user": user.ID,
		"list": ID,
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

const deleteUserListsQuery = `DELETE FROM user_has_list
WHERE list_id = %s`

// Delete removes a list and all its associated user lists.
func (s ListStore) Delete(list *List) error {
	tx, err := s.Begin()
	if err != nil {
		return err
	}

	if _, err = tx.Delete(list); err != nil {
		return err
	}

	_, err = tx.Exec(fmt.Sprintf(deleteUserListsQuery, s.Dialect.BindVar(0)), list.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

const allListsQuery = `SELECT l.* FROM list l
INNER JOIN user_has_list u ON u.list_id = l.id
WHERE u.user_id = :user
ORDER BY l.id DESC LIMIT :limit OFFSET :offset`

// All retrieves the all tags of the user with their count.
func (s ListStore) All(user int64, limit, offset int) ([]*List, error) {
	var lists []*List
	_, err := s.Select(&lists, allListsQuery, map[string]interface{}{
		"user":   user,
		"limit":  limit,
		"offset": offset,
	})
	if err != nil {
		return nil, err
	}

	userLists, err := s.UserLists(lists)
	if err != nil {
		return nil, err
	}

	var userListsByList = make(map[int64][]*UserHasList)
	for _, ul := range userLists {
		userListsByList[ul.ListID] = append(userListsByList[ul.ListID], ul)
	}

	for _, l := range lists {
		l.Users = userListsByList[l.ID]
	}

	return lists, nil
}

type userHasListWithUser struct {
	*UserHasList
	*User
}

const userListsQuery = `SELECT l.*, u.* FROM user_has_list l
INNER JOIN user u ON u.id = l.user_id
WHERE l.list_id IN %s`

// UserLists retrieves the users and their associations with the lists for
// the given collection of lists.
func (s ListStore) UserLists(lists []*List) ([]*UserHasList, error) {
	var IDs = make([]int64, len(lists))
	for i, l := range lists {
		IDs[i] = l.ID
	}

	var records []*userHasListWithUser
	if _, err := s.Select(&records, inQuery(userListsQuery, IDs)); err != nil {
		return nil, err
	}

	var userLists = make([]*UserHasList, len(records))
	for i, r := range records {
		l := r.UserHasList
		l.User = r.User
		userLists[i] = l
	}

	return userLists, nil
}

// Create inserts a new list on the database and associates it
// with the user.
func (s ListStore) Create(list *List, user *User) error {
	userList := &UserHasList{
		UserID: user.ID,
		User:   user,
		ListID: list.ID,
		Role:   ListOwner,
	}
	list.Users = append(list.Users, userList)
	return s.Insert(list, userList)
}

const countListsQuery = `SELECT COUNT(*) FROM user_has_list
WHERE user_id = :user`

// Count retrieves the number of unique lists an user has.
func (s ListStore) Count(user int64) (int64, error) {
	return s.SelectInt(countListsQuery, map[string]interface{}{
		"user": user,
	})
}
