package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mvader/pinlist/api/middlewares"
	"github.com/mvader/pinlist/api/models"
	"gopkg.in/gorp.v1"
)

// List is the List service, which will handle the list related
// API requests.
type List struct {
	db *gorp.DbMap
	*middlewares.Session
	store models.ListStore
}

// NewList returns a new List service given a database.
func NewList(db *gorp.DbMap) *List {
	return &List{
		db:      db,
		Session: middlewares.NewSession(db),
		store:   models.ListStore{DbMap: db},
	}
}

// Register autoregisters all the needed routes for the service
// on the given router.
func (l *List) Register(r *gin.RouterGroup) {
	r.GET("/lists", l.Auth, l.List)
	r.POST("/list", l.Auth, l.Create)
	r.PATCH("/list/:id", l.Auth, l.isOwner, l.getList, l.Update)
	r.DELETE("/list/:id", l.Auth, l.isOwner, l.getList, l.Delete)
}

// CreateListForm defines the schema of data to create a list.
type CreateListForm struct {
	Name        string `json:"name" binding:"required,gt=3,max=100"`
	Description string `json:"description" binding:"omitempty,gt=3,max=255"`
}

// Create handles the creation of new lists.
func (l *List) Create(c *gin.Context) {
	var form CreateListForm
	if err := c.BindJSON(&form); err != nil {
		badRequest(c)
		return
	}

	list := models.NewList(form.Name, form.Description)
	user := userFromCtx(c)
	if err := l.store.Create(list, user); err != nil {
		internalError(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

// ListListResponse defines the response structure of a list retrieval.
type ListListResponse struct {
	Count int            `json:"count"`
	Total int64          `json:"total"`
	Items []*models.List `json:"items"`
}

// List returns all the lists for an user.
func (l *List) List(c *gin.Context) {
	user := userFromCtx(c)
	limit := intParamOrDefault(c, limitParam, 25)
	offset := intParamOrDefault(c, offsetParam, 0)
	lists, err := l.store.All(user.ID, limit, offset)
	if err != nil {
		internalError(c, err)
		return
	}

	total, err := l.store.Count(user.ID)
	if err != nil {
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, ListListResponse{
		Count: len(lists),
		Total: total,
		Items: lists,
	})
}

// UpdateListForm defines the schema of data to update a list.
type UpdateListForm struct {
	Public      bool   `json:"public"`
	Name        string `json:"name" binding:"required,gt=3,max=100"`
	Description string `json:"description" binding:"omitempty,gt=3,max=255"`
	// Even though UserHasList is ready, some thought has to be
	// put on the issue before allowing sharing between users.
	// So, no shares yet.
}

// Update updates the settings and properties of a list.
func (l *List) Update(c *gin.Context) {
	var form UpdateListForm
	if err := c.BindJSON(&form); err != nil {
		badRequest(c)
		return
	}

	list := c.MustGet("list").(*models.List)
	list.Name = form.Name
	list.Description = form.Description
	list.Public = form.Public
	if _, err := l.store.Update(list); err != nil {
		internalError(c, err)
		return
	}

	ok(c)
}

// Delete removes a list.
func (l *List) Delete(c *gin.Context) {
	if err := l.store.Delete(c.MustGet("list").(*models.List)); err != nil {
		internalError(c, err)
		return
	}

	ok(c)
}

func (l *List) isOwner(c *gin.Context) {
	ok, err := l.store.UserIsOwner(userFromCtx(c), idFromCtx(c))
	if err != nil {
		internalError(c, err)
		return
	}

	if !ok {
		unauthorized(c)
		return
	}

	c.Next()
}

func (l *List) getList(c *gin.Context) {
	list, err := l.store.ByID(idFromCtx(c))
	if err != nil {
		internalError(c, err)
		return
	}

	if list == nil {
		notFound(c)
		return
	}

	c.Set("list", list)
	c.Next()
}
