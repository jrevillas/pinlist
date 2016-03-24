package services

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/gorp.v1"
	"github.com/mvader/pinlist/api/middlewares"
)

type Pin struct {
	db *gorp.DbMap
	*middlewares.Session
}

func NewPin(db *gorp.DbMap) *Pin {
	return &Pin{db: db, Session: middlewares.NewSession(db)}
}

func (b *Pin) Register(r *gin.RouterGroup) {
	r.POST("/pin", b.Create)
	r.PATCH("/pin/:id", b.Update)
	r.DELETE("/pin/:id", b.Delete)
	r.GET("/pins", b.List)
}

func (b *Pin) Create(c *gin.Context) {

}

func (b *Pin) Update(c *gin.Context) {

}

func (b *Pin) Delete(c *gin.Context) {

}

func (b *Pin) List(c *gin.Context) {

}
