package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/mvader/pinlist/api/models"
	"github.com/mvader/pinlist/api/testutil"
	"gopkg.in/check.v1"
	"gopkg.in/gorp.v1"
)

type PinSuite struct {
	db   *gorp.DbMap
	pin  *Pin
	file string

	users []*models.User
	lists []*models.List
}

var _ = check.Suite(&PinSuite{})

func (s *PinSuite) SetUpSuite(c *check.C) {
	file, db, err := testutil.CreateTestDatabase()
	c.Assert(err, check.IsNil)
	s.file = file
	s.db = db
	s.pin = NewPin(s.db)
	s.insertFixtures(c)
}

func (s *PinSuite) insertFixtures(c *check.C) {
	user1 := models.NewUser("foo", "foo", "foo")
	user2 := models.NewUser("bar", "bar", "bar")
	c.Assert(s.pin.store.Insert(user1, user2), check.IsNil)

	list := &models.List{Name: "foo"}
	list2 := &models.List{Name: "bar"}
	c.Assert(s.pin.store.Insert(list, list2), check.IsNil)

	uhl1 := &models.UserHasList{ListID: list.ID, UserID: user1.ID}
	uhl2 := &models.UserHasList{ListID: list.ID, UserID: user2.ID}
	c.Assert(s.pin.store.Insert(uhl1, uhl2), check.IsNil)

	pins := []*models.Pin{
		models.NewPin(user1, "pin1", "url1", nil, 0),
		models.NewPin(user1, "pin2", "url2", nil, 0),
		models.NewPin(user1, "pin3", "url3", nil, list.ID),
		models.NewPin(user2, "pin4", "url4", nil, 0),
		models.NewPin(user2, "pin5", "url5", nil, list.ID),
		models.NewPin(user2, "pin6", "url6", []string{"a", "b"}, list.ID),
	}

	for _, p := range pins {
		c.Assert(s.pin.store.Create(p), check.IsNil)
	}

	s.users = []*models.User{user1, user2}
	s.lists = []*models.List{list, list2}
}

func (s *PinSuite) TearDownSuite(c *check.C) {
	c.Assert(os.Remove(s.file), check.IsNil)
}

func (s *PinSuite) TestCreate(c *check.C) {
	tests := []struct {
		code int
		ok   bool
		user *models.User
		form CreatePinForm
	}{
		{http.StatusBadRequest, false, s.users[0], CreatePinForm{}},
		{http.StatusBadRequest, false, s.users[0], CreatePinForm{
			URL: "http://google.es",
		}},
		{http.StatusBadRequest, false, s.users[0], CreatePinForm{
			Title: "title",
		}},
		{http.StatusBadRequest, false, s.users[0], CreatePinForm{
			Title: "title", URL: "foo",
		}},
		{http.StatusBadRequest, false, s.users[0], CreatePinForm{
			Title: "title", URL: "http://google.es",
			Tags: []string{
				"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
				"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
				"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
				"31",
			},
		}},
		{http.StatusBadRequest, false, s.users[0], CreatePinForm{
			Title: "title", URL: "http://google.es",
			Tags: []string{"123456789123456789123456789123456789"},
		}},
		{http.StatusBadRequest, false, s.users[0], CreatePinForm{
			Title: "title", URL: "http://google.es",
			Tags: []string{"", "hi"},
		}},
		{http.StatusBadRequest, false, s.users[0], CreatePinForm{
			Title: "title", URL: "http://google.es",
			List: -1,
		}},
		{http.StatusUnauthorized, false, s.users[0], CreatePinForm{
			Title: "title", URL: "http://google.es",
			List: s.lists[1].ID,
		}},
		{http.StatusCreated, true, s.users[0], CreatePinForm{
			Title: "title", URL: "http://google.es",
			Tags: []string{"a", "b"},
			List: s.lists[0].ID,
		}},
		{http.StatusCreated, true, s.users[0], CreatePinForm{
			Title: "title", URL: "http://google.es",
			Tags: []string{"a", "b"},
		}},
	}

	for _, t := range tests {
		data, err := json.Marshal(&t.form)
		c.Assert(err, check.IsNil)
		r := testutil.MustRequest("POST", "/create", data)
		w := testutil.ExecuteHandler(
			testutil.Handler("/create", testutil.WithAuth(t.user), s.pin.Create),
			r,
		)
		c.Assert(w.Code, check.Equals, t.code)

		if t.ok {
			s.assertPin(c, w, t.form)
		}
	}
}

func (s *PinSuite) TestDelete(c *check.C) {
	pin := models.NewPin(s.users[0], "title", "url", []string{"a", "b"}, s.lists[1].ID)
	c.Assert(s.pin.store.Create(pin), check.IsNil)

	tests := []struct {
		code int
		user *models.User
		pin  int64
		ok   bool
	}{
		{http.StatusNotFound, s.users[0], 700, false},
		{http.StatusUnauthorized, s.users[1], pin.ID, false},
		{http.StatusOK, s.users[0], pin.ID, true},
	}

	for _, t := range tests {
		r := testutil.MustRequest("DELETE", fmt.Sprintf("/pin/%d", t.pin), nil)
		w := testutil.ExecuteHandler(
			testutil.Handler("/pin/:id", testutil.WithAuth(t.user), s.pin.GetPin, s.pin.Delete),
			r,
		)

		c.Assert(w.Code, check.Equals, t.code)

		if t.ok {
			p, err := s.pin.store.Get(models.Pin{}, pin.ID)
			c.Assert(err, check.IsNil)
			c.Assert(p, check.IsNil)
			n, err := s.pin.store.SelectInt(fmt.Sprintf("SELECT COUNT(*) FROM tag WHERE pin_id = %d", pin.ID))
			c.Assert(err, check.IsNil)
			c.Assert(n, check.Equals, int64(0))

			l, err := s.pin.store.Get(models.List{}, s.lists[1].ID)
			c.Assert(err, check.IsNil)
			c.Assert(l.(*models.List).Pins, check.Equals, 0)
		}
	}
}

func (s *PinSuite) assertPin(c *check.C, w *httptest.ResponseRecorder, form CreatePinForm) {
	var pin models.Pin
	err := json.Unmarshal(w.Body.Bytes(), &pin)
	c.Assert(err, check.IsNil)

	c.Assert(pin.Title, check.Equals, form.Title)
	c.Assert(pin.URL, check.Equals, form.URL)
	c.Assert(pin.ListID, check.Equals, form.List)
	c.Assert(pin.Tags, check.HasLen, len(form.Tags))
	for i, t := range pin.Tags {
		c.Assert(t.Name, check.Equals, form.Tags[i])
	}
}
