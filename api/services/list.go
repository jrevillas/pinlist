package services

import (
	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
	"github.com/mvader/pinlist/api/middlewares"
)

type List struct {
	db *gorp.DbMap
	*middlewares.Session
}

func NewList(db *gorp.DbMap) *List {
	return &List{db: db, Session: middlewares.NewSession(db)}
}

func (l *List) Register(r *gin.RouterGroup) {
	r.GET("/lists", l.Auth, l.List)
	r.PATCH("/list/:id", l.Auth, l.Update)
	r.DELETE("/list/:id", l.Auth, l.Delete)
}

func (l *List) List(c *gin.Context) {

}

func (l *List) Update(c *gin.Context) {

}

func (l *List) Delete(c *gin.Context) {

}
