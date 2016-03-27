package models

import (
	"os"

	"gopkg.in/check.v1"
	"gopkg.in/gorp.v1"
)

type TagSuite struct {
	db    *gorp.DbMap
	file  string
	store TagStore
}

var _ = check.Suite(&TagSuite{})

func (s *TagSuite) SetUpSuite(c *check.C) {
	file, db, err := createTestDatabase()
	c.Assert(err, check.IsNil)
	s.file = file
	s.db = db
	s.store = TagStore{DbMap: s.db}

	s.insertFixtures(c)
}

func (s *TagSuite) insertFixtures(c *check.C) {
	u := &User{ID: 1}
	u2 := &User{ID: 2}
	l := &UserHasList{ListID: 1, UserID: 1}
	c.Assert(s.store.Insert(l), check.IsNil)
	pins := []*Pin{
		NewPin(u, "", "", []string{"a", "b", "c", "d"}, 0),
		NewPin(u, "", "", []string{"j", "b", "c", "d"}, 0),
		NewPin(u, "", "", []string{"a", "b", "k", "i"}, 0),
		NewPin(u2, "", "", []string{"a", "b", "s", "w"}, 0),
		NewPin(u2, "", "", []string{"a", "b", "c", "t"}, 1),
	}

	for _, p := range pins {
		c.Assert(PinStore{DbMap: s.db}.Create(p), check.IsNil)
	}
}

func (s *TagSuite) TearDownSuite(c *check.C) {
	c.Assert(os.Remove(s.file), check.IsNil)
}

func (s *TagSuite) TestAll(c *check.C) {
	tags, err := s.store.All(1, 5, 0)
	c.Assert(err, check.IsNil)
	c.Assert(tags, check.HasLen, 5)
	c.Assert(tags[0].Name, check.Equals, "a")
	c.Assert(tags[0].Count, check.Equals, 3)
	c.Assert(tags[1].Name, check.Equals, "b")
	c.Assert(tags[1].Count, check.Equals, 4)
	c.Assert(tags[2].Name, check.Equals, "c")
	c.Assert(tags[2].Count, check.Equals, 3)
	c.Assert(tags[3].Name, check.Equals, "d")
	c.Assert(tags[3].Count, check.Equals, 2)
	c.Assert(tags[4].Name, check.Equals, "i")
	c.Assert(tags[4].Count, check.Equals, 1)

	tags, err = s.store.All(1, 5, 5)
	c.Assert(err, check.IsNil)
	c.Assert(tags, check.HasLen, 3)
	c.Assert(tags[0].Name, check.Equals, "j")
	c.Assert(tags[0].Count, check.Equals, 1)
	c.Assert(tags[1].Name, check.Equals, "k")
	c.Assert(tags[1].Count, check.Equals, 1)
	c.Assert(tags[2].Name, check.Equals, "t")
	c.Assert(tags[2].Count, check.Equals, 1)
}
