package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mvader/pinlist/api/middlewares"
	"github.com/mvader/pinlist/api/models"
	"gopkg.in/gorp.v1"
)

// Tag is the tag service, which is in charge of handling all
// the tag related API requests.
type Tag struct {
	db *gorp.DbMap
	*middlewares.Session
	store models.TagStore
}

// NewTag returns a new Tag service given a database.
func NewTag(db *gorp.DbMap) *Tag {
	return &Tag{
		db:      db,
		Session: middlewares.NewSession(db),
		store:   models.TagStore{DbMap: db},
	}
}

// Register autoregisters all the needed routes for the service
// on the given router.
func (t *Tag) Register(r *gin.RouterGroup) {
	r.GET("/tags", t.Auth, t.List)
}

func (t *Tag) user(ctx *gin.Context) *models.User {
	return ctx.MustGet(middlewares.UserKey).(*models.User)
}

// TagListResponse is the response with all tags retrieved.
type TagListResponse struct {
	Count int                `json:"count"`
	Total int64              `json:"total"`
	Items []*models.TagCount `json:"items"`
}

// List returns all the tags for an user.
func (t *Tag) List(c *gin.Context) {
	user := t.user(c)
	limit := intParamOrDefault(c, limitParam, 25)
	offset := intParamOrDefault(c, offsetParam, 0)
	tags, err := t.store.All(user.ID, limit, offset)
	if err != nil {
		internalError(c, err)
		return
	}

	total, err := t.store.Count(user.ID)
	if err != nil {
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, TagListResponse{
		Count: len(tags),
		Total: total,
		Items: tags,
	})
}
