package testutil

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/mvader/pinlist/api/models"
)

// WithAuth returns a middleware to naively set the session user.
func WithAuth(user *models.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("session_user", user)
		c.Next()
	}
}

// Handler returns a new gin Engine given a path and handlers.
func Handler(path string, handlers ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Any(path, handlers...)
	return r
}

// ExecuteHandler executes the given handler calling it with the
// given request and returns a ResponseRecorder.
func ExecuteHandler(handler http.Handler, r *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	return w
}

// MustRequest creates a request with the given data, panicking if error.
func MustRequest(method, path string, data []byte) *http.Request {
	r, err := http.NewRequest(method, path, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	return r
}
