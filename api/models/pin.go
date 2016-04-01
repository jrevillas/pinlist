package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gopkg.in/gorp.v1"
)

// Pin is a link pinned to a list or to an account.
type Pin struct {
	ID        int64  `db:"id" json:"id"`
	Title     string `db:"title" json:"title"`
	URL       string `db:"url" json:"url"`
	Tags      []*Tag `db:"-" json:"tags"`
	CreatorID int64  `db:"creator_id" json:"-"`
	Creator   *User  `db:"-" json:"creator"`
	// TODO: Fill List when (without users) for output
	ListID    sql.NullInt64 `db:"list_id" json:"-"`
	CreatedAt time.Time     `db:"created_at" json:"created_at"`
}

// NewPin creates a new pin with all its fields.
func NewPin(creator *User, title, url string, tags []string, list int64) *Pin {
	var tagList []*Tag
	var tagMap = make(map[string]struct{})
	for _, t := range tags {
		t = strings.ToLower(strings.TrimSpace(t))
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
		ListID:    sql.NullInt64{Int64: list, Valid: list > 0},
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

	if pin.ListID.Int64 > 0 {
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

	if p.ListID.Int64 > 0 {
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

type pinWithUser struct {
	BasicUser
	Pin
}

const (
	allForUserQuery = `SELECT DISTINCT p.*, u.username, u.status, u.email FROM pin p
INNER JOIN "user" u ON u.id = p.creator_id
LEFT JOIN user_has_list l ON l.list_id = p.list_id
WHERE (p.creator_id = :user OR l.user_id = :user)
ORDER BY p.created_at DESC LIMIT :limit`
	allForUserQueryWithOffset = `SELECT DISTINCT p.*, u.username, u.status, u.email FROM pin p
INNER JOIN "user" u ON u.id = p.creator_id
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
	allForListQuery = `SELECT p.*, u.username, u.status, u.email FROM pin p
INNER JOIN "user" u ON u.id = p.creator_id
WHERE p.list_id = :list
ORDER BY p.created_at DESC LIMIT :limit`
	allForListQueryWithOffset = `SELECT p.*, u.username, u.status, u.email FROM pin p
INNER JOIN "user" u ON u.id = p.creator_id
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
	var pinsWithUser []*pinWithUser
	_, err := s.Select(&pinsWithUser, q, params)
	if err != nil {
		return nil, err
	}

	var pinIDs []int64
	var pins = make([]*Pin, len(pinsWithUser))
	for i, p := range pinsWithUser {
		pinIDs = append(pinIDs, p.ID)
		p.Pin.Creator = &User{
			ID:       p.Pin.CreatorID,
			Username: p.BasicUser.Username,
			Status:   p.BasicUser.Status,
			Email:    p.BasicUser.Email,
		}
		pins[i] = &p.Pin
	}

	tags, err := s.pinTags(pinIDs)
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
