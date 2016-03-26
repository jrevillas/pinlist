package models

import (
	"os"

	"gopkg.in/check.v1"
	"gopkg.in/gorp.v1"
)

type PinSuite struct {
	db    *gorp.DbMap
	file  string
	store PinStore
}

var _ = check.Suite(&PinSuite{})

func (s *PinSuite) SetUpSuite(c *check.C) {
	file, db, err := createTestDatabase()
	c.Assert(err, check.IsNil)
	s.file = file
	s.db = db
	s.store = PinStore{DbMap: s.db}

	s.insertFixtures(c)
}

func (s *PinSuite) insertFixtures(c *check.C) {
	user1 := NewUser("foo", "foo", "foo")
	user2 := NewUser("bar", "bar", "bar")
	c.Assert(s.store.Insert(user1, user2), check.IsNil)

	list := &List{Name: "foo"}
	list2 := &List{Name: "bar"}
	c.Assert(s.store.Insert(list, list2), check.IsNil)

	uhl1 := &UserHasList{ListID: list.ID, UserID: user1.ID}
	uhl2 := &UserHasList{ListID: list.ID, UserID: user2.ID}
	c.Assert(s.store.Insert(uhl1, uhl2), check.IsNil)

	pins := []*Pin{
		NewPin(user1, "pin1", "url1", nil, 0),
		NewPin(user1, "pin2", "url2", nil, 0),
		NewPin(user1, "pin3", "url3", nil, list.ID),
		NewPin(user2, "pin4", "url4", nil, 0),
		NewPin(user2, "pin5", "url5", nil, list.ID),
		NewPin(user2, "pin6", "url6", []string{"a", "b"}, list.ID),
	}

	for _, p := range pins {
		c.Assert(s.store.Create(p), check.IsNil)
	}
}

func (s *PinSuite) TearDownSuite(c *check.C) {
	c.Assert(os.Remove(s.file), check.IsNil)
}

func (s *PinSuite) TestNewPin(c *check.C) {
	p := NewPin(&User{ID: 3}, "title", "URL", []string{"a", "b", "a"}, 1)
	c.Assert(p.Title, check.Equals, "title")
	c.Assert(p.URL, check.Equals, "URL")
	c.Assert(p.ListID, check.Equals, int64(1))
	c.Assert(p.CreatorID, check.Equals, int64(3))
	c.Assert(p.Tags, check.HasLen, 2)
	c.Assert(p.Tags[0].Name, check.Equals, "a")
	c.Assert(p.Tags[1].Name, check.Equals, "b")
}

func (s *PinSuite) TestCreate(c *check.C) {
	p := NewPin(&User{ID: 2}, "title", "URL", []string{"a", "b", "a"}, 2)
	c.Assert(s.store.Create(p), check.IsNil)
	p1, err := s.store.ByID(p.ID)
	c.Assert(err, check.IsNil)
	c.Assert(p1, check.Not(check.IsNil))

	var tags []*Tag
	_, err = s.store.Select(&tags, "SELECT * FROM tag WHERE pin_id = $1", p.ID)
	c.Assert(err, check.IsNil)
	c.Assert(tags, check.HasLen, 2)
	c.Assert(tags[0].Name, check.Equals, "a")
	c.Assert(tags[1].Name, check.Equals, "b")

	l, err := s.store.Get(List{}, p.ListID)
	c.Assert(err, check.IsNil)
	c.Assert(l.(*List).Pins, check.Equals, 1)
}

func (s *PinSuite) TestAllForUser(c *check.C) {
	pins, err := s.store.AllForUser(&User{ID: 1, Username: "foo"}, 4, 0)
	c.Assert(err, check.IsNil)
	c.Assert(pins, check.HasLen, 4)
	c.Assert(pins[0].Title, check.Equals, "pin6")
	c.Assert(pins[0].Creator.Username, check.Equals, "bar")
	c.Assert(pins[0].Tags, check.HasLen, 2)
	c.Assert(pins[0].Tags[0].Name, check.Equals, "a")
	c.Assert(pins[0].Tags[1].Name, check.Equals, "b")
	c.Assert(pins[1].Title, check.Equals, "pin5")
	c.Assert(pins[1].Creator.Username, check.Equals, "bar")
	c.Assert(pins[2].Title, check.Equals, "pin3")
	c.Assert(pins[2].Creator.Username, check.Equals, "foo")
	c.Assert(pins[3].Title, check.Equals, "pin2")
	c.Assert(pins[3].Creator.Username, check.Equals, "foo")

	pins, err = s.store.AllForUser(&User{ID: 1, Username: "foo"}, 10, pins[3].ID)
	c.Assert(err, check.IsNil)
	c.Assert(pins, check.HasLen, 1)
	c.Assert(pins[0].Title, check.Equals, "pin1")
	c.Assert(pins[0].Creator.Username, check.Equals, "foo")
}

func (s *PinSuite) TestAllForList(c *check.C) {
	pins, err := s.store.AllForList(&User{ID: 1, Username: "foo"}, 1, 2, 0)
	c.Assert(err, check.IsNil)
	c.Assert(pins, check.HasLen, 2)
	c.Assert(pins[0].Title, check.Equals, "pin6")
	c.Assert(pins[0].Creator.Username, check.Equals, "bar")
	c.Assert(pins[0].Tags, check.HasLen, 2)
	c.Assert(pins[0].Tags[0].Name, check.Equals, "a")
	c.Assert(pins[0].Tags[1].Name, check.Equals, "b")
	c.Assert(pins[1].Title, check.Equals, "pin5")
	c.Assert(pins[1].Creator.Username, check.Equals, "bar")

	pins, err = s.store.AllForList(&User{ID: 1, Username: "foo"}, 1, 10, pins[1].ID)
	c.Assert(err, check.IsNil)
	c.Assert(pins, check.HasLen, 1)
	c.Assert(pins[0].Title, check.Equals, "pin3")
	c.Assert(pins[0].Creator.Username, check.Equals, "foo")
}
