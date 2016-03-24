package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
	"github.com/mvader/pinlist/api/log"
	"github.com/mvader/pinlist/api/models"
)

const (
	TokenKey   = "session_token"
	UserKey    = "session_user"
	AuthKey    = "Authorization"
	TokenParam = "token"
)

type Session struct {
	db    *gorp.DbMap
	store models.UserStore
}

func NewSession(db *gorp.DbMap) *Session {
	return &Session{db: db, store: models.UserStore{db}}
}

func (s *Session) Guest(c *gin.Context) {
	if retrieveToken(c) != "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}

func (s *Session) Auth(c *gin.Context) {
	token := retrieveToken(c)
	user, err := s.store.ByToken(token)
	if err != nil {
		log.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if user == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set(TokenKey, token)
	c.Set(UserKey, user)
	c.Next()
}

func retrieveToken(c *gin.Context) string {
	auth := c.Request.Header.Get(AuthKey)
	if auth != "" {
		parts := strings.SplitN(auth, " ", 1)
		return parts[len(parts)-1]
	}

	return c.Query(TokenParam)
}
