package models

import (
	"fmt"
	"time"

	"gopkg.in/gorp.v1"
)

// Pin is a link pinned to a list or to an account.
type Pin struct {
	ID        int64     `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	URL       string    `db:"url" json:"url"`
	Tags      []*Tag    `db:"-" json:"tags"`
	CreatorID int64     `db:"creator_id" json:"-"`
	Creator   *User     `db:"-" json:"creator"`
	ListID    int64     `db:"list_id" json:"list_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewPin creates a new pin with all its fields.
func NewPin(creator *User, title, url string, tags []string, list int64) *Pin {
	var tagList []*Tag
	var tagMap = make(map[string]struct{})
	for _, t := range tags {
		// make sure no repeated tags make it to the db
		if _, ok := tagMap[t]; ok {
			continue
		}
		tagMap[t] = struct{}{}
		tagList = append(tagList, NewTag(t))
	}

	return &Pin{
		Title:     title,
		URL:       url,
		Tags:      tagList,
		CreatorID: creator.ID,
		Creator:   creator,
		ListID:    list,
		CreatedAt: time.Now(),
	}
}

// PinStore is the service to execute operations about pins
// that require the database.
type PinStore struct {
	*gorp.DbMap
}

const updateListPinsQuery = `UPDATE list
SET pins = pins + 1 WHERE id = %s`

// Create inserts a pin into the database and its associated tags.
func (s PinStore) Create(pin *Pin) error {
	tx, err := s.Begin()
	if err != nil {
		return err
	}

	if err := tx.Insert(pin); err != nil {
		return err
	}

	var tags []interface{}
	for _, t := range pin.Tags {
		t.PinID = pin.ID
		tags = append(tags, t)
	}
	if err := tx.Insert(tags...); err != nil {
		return err
	}

	if pin.ListID > 0 {
		q := fmt.Sprintf(updateListPinsQuery, s.Dialect.BindVar(0))
		if _, err := tx.Exec(q, pin.ListID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

const (
	decrListPinsQuery = `UPDATE list
SET pins = pins - 1 WHERE id = %s`
	deletePinTagsQuery = `DELETE FROM tag
WHERE pin_id = %s`
)

// Delete removes a pin and its tags and also decrements the
// count of pins in a list, if any.
func (s PinStore) Delete(p *Pin) error {
	tx, err := s.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Delete(p); err != nil {
		return err
	}

	if _, err := tx.Exec(fmt.Sprintf(deletePinTagsQuery, s.Dialect.BindVar(0)), p.ID); err != nil {
		return err
	}

	if p.ListID > 0 {
		q := fmt.Sprintf(decrListPinsQuery, s.Dialect.BindVar(0))
		if _, err := tx.Exec(q, p.ListID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// ByID retrieves a pin by its id.
func (s PinStore) ByID(ID int64) (*Pin, error) {
	p, err := s.Get(Pin{}, ID)
	if err != nil || p == nil {
		return nil, err
	}

	return p.(*Pin), nil
}

const (
	allForUserQuery = `SELECT DISTINCT p.* FROM pin p
LEFT JOIN user_has_list l ON l.list_id = p.list_id
WHERE (p.creator_id = :user OR l.user_id = :user)
ORDER BY p.created_at DESC LIMIT :limit`
	allForUserQueryWithOffset = `SELECT DISTINCT p.* FROM pin p
LEFT JOIN user_has_list l ON l.list_id = p.list_id
WHERE (p.creator_id = :user OR l.user_id = :user)
AND p.id < :offset
ORDER BY p.created_at DESC LIMIT :limit`
)

// AllForUser retrieves all pins for the user with the limit and offset received.
func (s PinStore) AllForUser(user *User, limit int, offset int64) ([]*Pin, error) {
	q := allForUserQuery
	params := map[string]interface{}{
		"user":  user.ID,
		"limit": limit,
	}
	if offset > 0 {
		q = allForUserQueryWithOffset
		params["offset"] = offset
	}

	return s.withTagsAndCreator(user, q, params)
}

const (
	allForListQuery = `SELECT p.* FROM pin p
WHERE p.list_id = :list
ORDER BY p.created_at DESC LIMIT :limit`
	allForListQueryWithOffset = `SELECT p.* FROM pin p
WHERE p.list_id = :list
AND p.id < :offset
ORDER BY p.created_at DESC LIMIT :limit`
)

// AllForList retrieves all pins for the list with the limit and offset received.
func (s PinStore) AllForList(user *User, list int64, limit int, offset int64) ([]*Pin, error) {
	q := allForListQuery
	params := map[string]interface{}{
		"list":  list,
		"limit": limit,
	}
	if offset > 0 {
		q = allForListQueryWithOffset
		params["offset"] = offset
	}

	return s.withTagsAndCreator(user, q, params)
}

func (s PinStore) withTagsAndCreator(creator *User, q string, params map[string]interface{}) ([]*Pin, error) {
	var pins []*Pin
	_, err := s.Select(&pins, q, params)
	if err != nil {
		return nil, err
	}

	var tagIDs []int64
	var userIDs []int64
	for _, p := range pins {
		tagIDs = append(tagIDs, p.ID)
		if p.CreatorID != creator.ID {
			userIDs = append(userIDs, p.CreatorID)
		} else {
			p.Creator = creator
		}
	}

	tags, err := s.pinTags(tagIDs)
	if err != nil {
		return nil, err
	}

	for _, t := range tags {
		for _, p := range pins {
			if p.ID == t.PinID {
				p.Tags = append(p.Tags, t)
				break
			}
		}
	}

	if len(userIDs) == 0 {
		return pins, nil
	}

	users, err := s.pinCreators(userIDs)
	if err != nil {
		return nil, err
	}

	for _, p := range pins {
		for _, u := range users {
			if p.CreatorID == creator.ID {
				break
			} else if p.CreatorID == u.ID {
				p.Creator = u
				break
			}
		}
	}

	return pins, nil
}

const pinTagsQuery = `SELECT * FROM tag
WHERE pin_id IN %s`

func (s PinStore) pinTags(ids []int64) ([]*Tag, error) {
	var tags []*Tag
	_, err := s.Select(&tags, inQuery(pinTagsQuery, ids))
	if err != nil {
		return nil, err
	}

	return tags, nil
}

const pinCreatorsQuery = `SELECT * FROM user
WHERE id IN %s`

func (s PinStore) pinCreators(ids []int64) ([]*User, error) {
	var users []*User
	_, err := s.Select(&users, inQuery(pinCreatorsQuery, ids))
	if err != nil {
		return nil, err
	}

	return users, nil
}
