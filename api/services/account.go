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
	g.POST("/create", a.Guest, a.Create)
	g.POST("/login", a.Guest, a.Login)
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
		c.AbortWithStatus(http.StatusInternalServerError)
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
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, user)
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

	user, err := a.store.ByLoginDetails(form.Login, form.Password)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token := models.NewToken(user.ID)
	if err := a.store.Insert(token); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, token)
}

func (a *Account) Logout(c *gin.Context) {
	token := c.MustGet(middlewares.TokenKey).(string)
	if err := a.store.DeleteToken(token); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
