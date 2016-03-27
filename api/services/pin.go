package services

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mvader/pinlist/api/middlewares"
	"github.com/mvader/pinlist/api/models"
	"gopkg.in/gorp.v1"
)

// Pin is the pin service, which has all the endpoints
// of the API that handle pin related requests.
type Pin struct {
	db *gorp.DbMap
	*middlewares.Session
	store     models.PinStore
	listStore models.ListStore
}

// NewPin returns a new Pin Service given a database.
func NewPin(db *gorp.DbMap) *Pin {
	return &Pin{
		db:        db,
		Session:   middlewares.NewSession(db),
		store:     models.PinStore{DbMap: db},
		listStore: models.ListStore{DbMap: db},
	}
}

// Register autoregisters all the needed routes for the service
// on the given router.
func (s *Pin) Register(r *gin.RouterGroup) {
	r.POST("/pin", s.Auth, s.Create)
	r.DELETE("/pin/:id", s.Auth, s.GetPin, s.Delete)
	r.GET("/pins", s.MaybeAuth, s.List)
}

func (s *Pin) user(ctx *gin.Context) *models.User {
	return ctx.MustGet(middlewares.UserKey).(*models.User)
}

func (s *Pin) id(ctx *gin.Context) int64 {
	n, _ := strconv.ParseInt(ctx.Params.ByName("id"), 10, 64)
	return n
}

func (s *Pin) list(ctx *gin.Context) int64 {
	n, _ := strconv.ParseInt(ctx.Query("list"), 10, 64)
	return n
}

// CreatePinForm is the structure of the form to create a
// new pin. It also contains the validation via binding.
type CreatePinForm struct {
	Title string   `json:"title" binding:"required,max=255"`
	URL   string   `json:"name" binding:"required,url,max=255"`
	Tags  []string `json:"tags" binding:"omitempty,lte=30,dive,required,max=30"`
	List  int64    `json:"list" binding:"omitempty,gte=1"`
}

// Create handles the creation of a pin.
func (s *Pin) Create(c *gin.Context) {
	var form CreatePinForm
	if err := c.BindJSON(&form); err != nil {
		badRequest(c)
		return
	}

	user := s.user(c)
	pin := models.NewPin(user, form.Title, form.URL, form.Tags, form.List)

	if form.List > 0 {
		ok, err := s.listStore.UserHasAccess(user, form.List)
		if err != nil {
			internalError(c, err)
			return
		}

		if !ok {
			unauthorized(c)
			return
		}
	}

	if err := s.store.Create(pin); err != nil {
		internalError(c, err)
		return
	}

	c.JSON(http.StatusCreated, pin)
}

// Delete handles the delete of a pin.
func (s *Pin) Delete(c *gin.Context) {
	if err := s.store.Delete(c.MustGet("pin").(*models.Pin)); err != nil {
		internalError(c, err)
		return
	}

	ok(c)
}

// PinListResponse is the response sent with a list of pins.
type PinListResponse struct {
	Count int           `json:"count"`
	Items []*models.Pin `json:"items"`
}

// List returns the list of pins for an user.
func (s *Pin) List(c *gin.Context) {
	user := s.user(c)
	list := s.list(c)
	limit, offset := limitAndOffset(c)

	// user is not logged in and did not request a list
	// so we can't provide a list of their pins
	if user == nil && list == 0 {
		unauthorized(c)
		return
	}

	if list > 0 {
		s.listPins(c, user, list, limit, offset)
		return
	}
	s.userPins(c, user, limit, offset)
}

func (s *Pin) listPins(c *gin.Context, user *models.User, list int64, limit int, offset int64) {
	ok, err := s.listStore.UserHasAccess(user, list)
	if err != nil {
		internalError(c, err)
		return
	}

	if !ok {
		unauthorized(c)
		return
	}

	pins, err := s.store.AllForList(user, list, limit, offset)
	if err != nil {
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, &PinListResponse{
		Count: len(pins),
		Items: pins,
	})
}

func (s *Pin) userPins(c *gin.Context, user *models.User, limit int, offset int64) {
	pins, err := s.store.AllForUser(user, limit, offset)
	if err != nil {
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, &PinListResponse{
		Count: len(pins),
		Items: pins,
	})
}

// GetPin is a middleware that checks before executing the
// handler that the requested pin exists and belongs to the
// user issuing the request.
func (s *Pin) GetPin(c *gin.Context) {
	user := s.user(c)
	ID := s.id(c)

	pin, err := s.store.ByID(ID)
	if err != nil {
		internalError(c, err)
		return
	}

	if pin == nil {
		notFound(c)
		return
	}

	if pin.CreatorID != user.ID {
		unauthorized(c)
		return
	}

	c.Set("pin", pin)
	c.Next()
}
