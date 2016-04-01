package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mvader/pinlist/api/log"
	"github.com/mvader/pinlist/api/models"
	"gopkg.in/gorp.v1"
)

const (
	// TokenKey is the context key of the session.
	TokenKey = "session_token"
	// FullTokenKey retrieves the full token object.
	FullTokenKey = "session_full_token"
	// UserKey is the context key of the user.
	UserKey = "session_user"
	// AuthKey is the Authorization header key.
	AuthKey = "Authorization"
	// TokenParam is the query string token param key.
	TokenParam = "token"
)

// Session is a middleware in charge of handling session
// and authentication.
type Session struct {
	db    *gorp.DbMap
	store models.UserStore
}

// NewSession returns a new session middleware.
func NewSession(db *gorp.DbMap) *Session {
	return &Session{db: db, store: models.UserStore{DbMap: db}}
}

// Guest with ensure the user is not logged in.
func (s *Session) Guest(c *gin.Context) {
	if retrieveToken(c) != "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}

// Auth will ensure the user is logged in and set the user in
// the context.
func (s *Session) Auth(c *gin.Context) {
	token := retrieveToken(c)
	if token == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, fullToken, err := s.store.ByToken(token)
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
	c.Set(FullTokenKey, fullToken)
	c.Set(UserKey, user)
	c.Next()
}

// MaybeAuth will set the user in the context if the user
// is logged in but will not abort if the user is a guest.
func (s *Session) MaybeAuth(c *gin.Context) {
	token := retrieveToken(c)
	user, fullToken, err := s.store.ByToken(token)
	if err != nil {
		log.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Set(TokenKey, token)
	c.Set(FullTokenKey, fullToken)
	c.Set(UserKey, user)
	c.Next()
}

func retrieveToken(c *gin.Context) string {
	auth := c.Request.Header.Get(AuthKey)
	if auth != "" {
		parts := strings.SplitN(auth, " ", 2)
		return parts[len(parts)-1]
	}

	return c.Query(TokenParam)
}
