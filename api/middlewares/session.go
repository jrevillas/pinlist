package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
)

const (
	TokenKey = "session_token"
	UserKey  = "session_user"
)

type Session struct {
	db *gorp.DbMap
}

func NewSession(db *gorp.DbMap) *Session {
	return &Session{db: db}
}

func (s *Session) Guest(c *gin.Context) {

}

func (s *Session) Auth(c *gin.Context) {

}
