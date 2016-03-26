package services

import (
	"github.com/gin-gonic/gin"
	"github.com/mvader/pinlist/api/middlewares"
	"gopkg.in/gorp.v1"
)

// Tag is the tag service, which is in charge of handling all
// the tag related API requests.
type Tag struct {
	db *gorp.DbMap
	*middlewares.Session
}

// NewTag returns a new Tag service given a database.
func NewTag(db *gorp.DbMap) *Tag {
	return &Tag{db: db, Session: middlewares.NewSession(db)}
}

// Register autoregisters all the needed routes for the service
// on the given router.
func (t *Tag) Register(r *gin.RouterGroup) {
	r.GET("/tags", t.Auth, t.List)
}

// List returns all the tags for an user.
func (t *Tag) List(c *gin.Context) {

}
