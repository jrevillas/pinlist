package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mvader/pinlist/api/middlewares"
	"github.com/mvader/pinlist/api/models"
	"github.com/mvader/pinlist/api/testutil"
	"gopkg.in/check.v1"
	"gopkg.in/gorp.v1"
)

type AccountSuite struct {
	db      *gorp.DbMap
	account *Account
	file    string
}

var _ = check.Suite(&AccountSuite{})

func (s *AccountSuite) SetUpSuite(c *check.C) {
	file, db, err := testutil.CreateTestDatabase()
	c.Assert(err, check.IsNil)
	s.file = file
	s.db = db
	s.account = NewAccount(s.db)
}

func (s *AccountSuite) TearDownSuite(c *check.C) {
	c.Assert(os.Remove(s.file), check.IsNil)
}

func (s *AccountSuite) TestCreate(c *check.C) {
	tests := []struct {
		code                      int
		ok                        bool
		username, email, password string
	}{
		{http.StatusBadRequest, false, "", "", ""},
		{http.StatusBadRequest, false, "a b c", "foo@foo.foo", "12345678"},
		{http.StatusBadRequest, false, "abc", "foo.foo", "12345678"},
		{http.StatusBadRequest, false, "abc", "foo@foo.foo", "1234567"},
		{http.StatusOK, true, "abc", "foo@foo.foo", "12345678"},
		{http.StatusBadRequest, false, "abc", "bar@foo.foo", "12345678"},
		{http.StatusBadRequest, false, "bca", "foo@foo.foo", "12345678"},
		{http.StatusOK, false, "bca", "bar@foo.foo", "12345678"},
	}

	for _, t := range tests {
		data, err := json.Marshal(CreateAccountForm{
			Email:    t.email,
			Username: t.username,
			Password: t.password,
		})
		c.Assert(err, check.IsNil)
		r := testutil.MustRequest("POST", "/create", data)
		w := testutil.ExecuteHandler(testutil.Handler("/create", s.account.Create), r)
		c.Assert(w.Code, check.Equals, t.code)

		if t.ok {
			s.assertToken(c, w, t.email, t.username)
		}
	}
}

func (s *AccountSuite) TestLogin(c *check.C) {
	tests := []struct {
		code            int
		ok              bool
		login, password string
	}{
		{http.StatusBadRequest, false, "", ""},
		{http.StatusBadRequest, false, "abc", ""},
		{http.StatusBadRequest, false, "", "12345678"},
		{http.StatusBadRequest, false, "abc", "12345679"},
		{http.StatusBadRequest, false, "foo@foo.foo", "12345679"},
		{http.StatusOK, true, "foo@foo.foo", "12345678"},
		{http.StatusOK, true, "abc", "12345678"},
	}

	for _, t := range tests {
		data, err := json.Marshal(LoginForm{
			Login:    t.login,
			Password: t.password,
		})
		c.Assert(err, check.IsNil)
		r := testutil.MustRequest("POST", "/login", data)
		w := testutil.ExecuteHandler(testutil.Handler("/login", s.account.Login), r)
		c.Assert(w.Code, check.Equals, t.code)

		if t.ok {
			s.assertToken(c, w, t.login, t.login)
		}
	}
}

func (s *AccountSuite) TestLogout(c *check.C) {
	t, err := s.db.Get(models.Token{}, 1)
	c.Assert(err, check.IsNil)
	token := t.(*models.Token)

	handler := func(c *gin.Context) {
		c.Set(middlewares.TokenKey, token.Hash)
		s.account.Logout(c)
	}

	r := testutil.MustRequest("POST", "/logout", nil)
	r.Header.Set("Authorization", "bearer "+token.Hash)
	w := testutil.ExecuteHandler(testutil.Handler("/logout", handler), r)

	c.Assert(w.Code, check.Equals, http.StatusOK)
	t, err = s.db.Get(models.Token{}, 1)
	c.Assert(t, check.IsNil)
	c.Assert(err, check.IsNil)
}

func (s *AccountSuite) assertToken(c *check.C, w *httptest.ResponseRecorder, email, login string) {
	var resp AuthResponse
	c.Assert(json.Unmarshal(w.Body.Bytes(), &resp), check.IsNil)
	c.Assert(resp.Token.Hash, check.Not(check.Equals), "")
	c.Assert(resp.Token.Until.IsZero(), check.Equals, false)
	c.Assert(resp.User.Email, check.Not(check.Equals), "")
	c.Assert(resp.User.Username, check.Not(check.Equals), "")
	ok, err := s.account.store.ExistsUser(email, login)
	c.Assert(err, check.IsNil)
	c.Assert(ok, check.Equals, true)
}
