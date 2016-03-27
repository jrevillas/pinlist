package testutil

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/mvader/pinlist/api/models"
)

func WithAuth(user *models.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("session_user", user)
		c.Next()
	}
}

func Handler(path string, handlers ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Any(path, handlers...)
	return r
}

func ExecuteHandler(handler http.Handler, r *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	return w
}

func MustRequest(method, path string, data []byte) *http.Request {
	r, err := http.NewRequest(method, path, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	return r
}
