package services

import (
	"github.com/gin-gonic/gin"
	"github.com/mvader/pinlist/api/middlewares"
	"gopkg.in/gorp.v1"
)

// List is the List service, which will handle the list related
// API requests.
type List struct {
	db *gorp.DbMap
	*middlewares.Session
}

// NewList returns a new List service given a database.
func NewList(db *gorp.DbMap) *List {
	return &List{db: db, Session: middlewares.NewSession(db)}
}

// Register autoregisters all the needed routes for the service
// on the given router.
func (l *List) Register(r *gin.RouterGroup) {
	r.GET("/lists", l.Auth, l.List)
	r.PATCH("/list/:id", l.Auth, l.Update)
	r.DELETE("/list/:id", l.Auth, l.Delete)
}

// List returns all the lists for an user.
func (l *List) List(c *gin.Context) {

}

// Update updates the settings and properties of a list.
func (l *List) Update(c *gin.Context) {

}

// Delete removes a list.
func (l *List) Delete(c *gin.Context) {

}
