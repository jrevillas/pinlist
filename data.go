package main

import (
	"github.com/codegangsta/martini"
	r "github.com/dancannon/gorethink"
	"github.com/gorilla/sessions"
	"log"
	"strconv"
	"time"
)

type Connection struct {
	session *r.Session
}

func (c *Connection) GetBookmarks(userID string, page int64) ([]Bookmark, error) {
	var bookmarks []Bookmark

	cursor, err := r.DB("magnet").
		Table("bookmarks").
		OrderBy(r.OrderByOpts{r.Desc("Created")}).
		Filter(r.Row.Field("User").Eq(userID)).
		Skip(50 * page).
		Limit(50).
		Run(c.session)

	if err != nil {
		log.Print(err)
		return bookmarks, err
	}

	cursor.All(&bookmarks)
	cursor.Close()
	return bookmarks, err
}

func (c *Connection) initDatabase(connectionString string) {
	c.SetSession(connectionString, "magnet")

	c.InitDatabase()

	c.WipeExpiredSessions()
}

func (c *Connection) SetSession(address, database string) {
	session, err := r.Connect(r.ConnectOpts{
		Address:  address,
		Database: "magnet",
	})

	if err != nil {
		log.Fatal("Error connecting:", err)
	}

	c.session = session
}

func (c *Connection) NewBookmark(userID string, bookmark map[string]interface{}) (r.WriteResponse, error) {
	var response r.WriteResponse

	cursor, err := r.DB("magnet").
		Table("bookmarks").
		Insert(bookmark).
		Run(c.session)

	if err != nil {
		log.Print(err)
	}

	cursor.One(&response)
	cursor.Close()
	return response, err
}

func (c *Connection) DeleteBookmark(userID string, params martini.Params) (r.WriteResponse, error) {
	var response r.WriteResponse

	cursor, err := r.DB("magnet").
		Table("bookmarks").
		Filter(r.Row.Field("User").Eq(userID).
		And(r.Row.Field("id").Eq(params["bookmark"]))).
		Delete().
		Run(c.session)

	cursor.One(&response)
	cursor.Close()
	return response, err
}

func (c *Connection) EditBookmark(userID string, params martini.Params, bookmark map[string]interface{}) (r.WriteResponse, error) {
	var response r.WriteResponse

	cursor, err := r.DB("magnet").
		Table("bookmarks").
		Filter(r.Row.Field("User").Eq(userID).
		And(r.Row.Field("id").Eq(params["bookmark"]))).
		Update(bookmark).
		Run(c.session)

	if err != nil {
		log.Print(err)
	}

	cursor.One(&response)
	cursor.Close()
	return response, err
}

func (c *Connection) Search(userID string, params martini.Params, query string) ([]interface{}, error) {
	var response []interface{}
	page, _ := strconv.ParseInt(params["page"], 10, 16)

	cursor, err := r.DB("magnet").
		Table("bookmarks").
		OrderBy(r.OrderByOpts{r.Desc("Created")}).
		Filter(r.Row.Field("Title").Match("(?i)" + query).
		And(r.Row.Field("User").Eq(userID))).
		Skip(50 * page).
		Limit(50).
		Run(c.session)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	cursor.All(&response)
	cursor.Close()
	return response, err
}

func (c *Connection) GetTag(userID string, params martini.Params) ([]interface{}, error) {
	var response []interface{}
	page, _ := strconv.ParseInt(params["page"], 10, 16)

	cursor, err := r.DB("magnet").
		Table("bookmarks").
		OrderBy(r.OrderByOpts{r.Desc("Created")}).
		Filter(r.Row.Field("User").Eq(userID).
		And(r.Row.Field("Tags").
		Contains(params["tag"]))).
		Skip(50 * page).
		Limit(50).
		Run(c.session)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	cursor.All(&response)
	cursor.Close()
	return response, err
}

func (c *Connection) LoginPost(username, password string) ([]interface{}, error) {
	var response []interface{}

	cursor, err := r.DB("magnet").
		Table("users").
		Filter(r.Row.Field("Username").Eq(username).
		And(r.Row.Field("Password").Eq(password))).
		Run(c.session)

	if err != nil {
		return nil, err
	}

	cursor.All(&response)
	cursor.Close()
	return response, err
}

func (c *Connection) LoginPostInsertSession(session Session) (r.WriteResponse, error) {
	var response r.WriteResponse

	cursor, err := r.DB("magnet").
		Table("sessions").
		Insert(session).
		Run(c.session)

	if err != nil {
		log.Print(err)
	}

	cursor.One(&response)
	cursor.Close()
	return response, err
}

func (c *Connection) Logout(session *sessions.Session) (r.WriteResponse, error) {
	var response r.WriteResponse

	cursor, err := r.DB("magnet").
		Table("sessions").
		Get(session.Values["session_id"]).
		Delete().
		Run(c.session)

	if err != nil {
		log.Print(err)
	}

	cursor.One(&response)
	cursor.Close()
	return response, err
}

func (c *Connection) SignUp(user *User) ([]interface{}, error) {
	var response []interface{}

	cursor, err := r.DB("magnet").
		Table("users").
		Filter(r.Row.Field("Username").Eq(user.Username).
		Or(r.Row.Field("Email").Eq(user.Email))).
		Run(c.session)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	cursor.All(&response)
	cursor.Close()
	return response, err
}

func (c *Connection) SignUpInsert(user *User) (r.WriteResponse, error) {
	var response r.WriteResponse

	cursor, err := r.DB("magnet").
		Table("users").
		Insert(user).
		Run(c.session)

	if err != nil {
		log.Print(err)
	}

	cursor.One(&response)
	cursor.Close()
	return response, err
}

func (c *Connection) InitDatabase() {
	r.DBCreate("magnet").Exec(c.session)
	r.TableCreate("users").Exec(c.session)
	r.TableCreate("bookmarks").Exec(c.session)
	_, err := r.DB("magnet").Table("bookmarks").IndexCreate("Created").RunWrite(c.session)
	if err != nil {
		log.Printf("Error creating index: %s", err)
	}
	r.TableCreate("sessions").Exec(c.session)
}

func (c *Connection) WipeExpiredSessions() (r.WriteResponse, error) {
	var response r.WriteResponse

	cursor, err := r.DB("magnet").
		Table("sessions").
		Filter(r.Row.Field("Expires").Lt(time.Now().Unix())).
		Delete().
		Run(c.session)

	if err != nil {
		log.Print(err)
	}

	cursor.One(&response)
	cursor.Close()
	return response, err
}

func (c *Connection) GetTags(userID string) ([]interface{}, error) {
	var response []interface{}

	cursor, err := r.DB("magnet").
		Table("bookmarks").
		Filter(r.Row.Field("User").Eq(userID)).
		WithFields("Tags").
		Run(c.session)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	cursor.All(&response)
	cursor.Close()
	return response, err
}

func (c *Connection) GetUnexpiredSession(session *sessions.Session) (map[string]interface{}, error) {
	var response map[string]interface{}

	cursor, err := r.DB("magnet").
		Table("sessions").
		Get(session.Values["session_id"]).
		Run(c.session)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	cursor.One(&response)
	cursor.Close()
	return response, err
}
