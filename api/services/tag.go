package services

import (
	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
	"github.com/mvader/pinlist/api/middlewares"
)

type Tag struct {
	db *gorp.DbMap
	*middlewares.Session
}

func NewTag(db *gorp.DbMap) *Tag {
	return &Tag{db: db, Session: middlewares.NewSession(db)}
}

func (t *Tag) Register(r *gin.RouterGroup) {
	r.GET("/tags", t.Auth, t.List)
}

func (t *Tag) List(c *gin.Context) {

}
