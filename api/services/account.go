package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
	"github.com/mvader/pinlist/api/middlewares"
	"github.com/mvader/pinlist/api/models"
)

type Account struct {
	db *gorp.DbMap
	*middlewares.Session
	store *models.UserStore
}

func NewAccount(db *gorp.DbMap) *Account {
	return &Account{
		db:      db,
		Session: middlewares.NewSession(db),
		store:   &models.UserStore{db},
	}
}

func (a *Account) Register(r *gin.RouterGroup) {
	g := r.Group("/account")
	g.POST("/create", middlewares.GuestHourLimit, a.Guest, a.Create)
	g.POST("/login", middlewares.GuestHourLimit, a.Guest, a.Login)
	g.POST("/logout", a.Auth, a.Logout)
}

type CreateAccountForm struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (a *Account) Create(c *gin.Context) {
	var form CreateAccountForm
	if err := c.BindJSON(&form); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ok, err := a.store.ExistsUser(form.Email, form.Username)
	if err != nil {
		internalError(c, err)
		return
	}

	if ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, ok := models.NewUser(form.Username, form.Email, form.Password)
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := a.store.Insert(user); err != nil {
		internalError(c, err)
		return
	}

	a.login(c, form.Email, form.Password)
}

type LoginForm struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (a *Account) Login(c *gin.Context) {
	var form LoginForm
	if err := c.BindJSON(&form); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	a.login(c, form.Login, form.Password)
}

func (a *Account) login(c *gin.Context, login, password string) {
	user, err := a.store.ByLoginDetails(login, password)
	if err != nil {
		internalError(c, err)
		return
	}

	if user == nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token := models.NewToken(user.ID)
	if err := a.store.Insert(token); err != nil {
		internalError(c, err)
		return
	}

	c.JSON(http.StatusOK, token)
}

func (a *Account) Logout(c *gin.Context) {
	token := c.MustGet(middlewares.TokenKey).(string)
	if err := a.store.DeleteToken(token); err != nil {
		internalError(c, err)
		return
	}
	c.Status(http.StatusOK)
}
