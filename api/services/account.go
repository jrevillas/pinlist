package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mvader/pinlist/api/middlewares"
	"github.com/mvader/pinlist/api/models"
	"gopkg.in/gorp.v1"
)

// Account is the account service, which has all the endpoints
// of the API that handle account related requests.
type Account struct {
	*middlewares.Session
	store *models.UserStore
}

// NewAccount returns a new Account service given a database.
func NewAccount(db *gorp.DbMap) *Account {
	return &Account{
		Session: middlewares.NewSession(db),
		store:   &models.UserStore{DbMap: db},
	}
}

// Register autoregisters all the needed routes for the service
// on the given router.
func (a *Account) Register(r *gin.RouterGroup) {
	g := r.Group("/account")
	{
		g.POST("/create", middlewares.GuestHourLimit, a.Guest, a.Create)
		g.POST("/login", middlewares.GuestHourLimit, a.Guest, a.Login)
		g.POST("/logout", a.Auth, a.Logout)
	}
}

// CreateAccountForm is the structure of the form to create
// a new account. It also contains the validations of itself.
type CreateAccountForm struct {
	Username string `json:"username" binding:"required,alphanum,max=60"`
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,min=8"`
}

// Create is the handler in charge of creating accounts.
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

	user := models.NewUser(form.Username, form.Email, form.Password)
	if err := a.store.Insert(user); err != nil {
		internalError(c, err)
		return
	}

	a.login(c, form.Email, form.Password)
}

// LoginForm is the structure of the form to login.
type LoginForm struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login is the handler in charge of logging the user in.
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

// Logout terminates the session associated to the current token.
func (a *Account) Logout(c *gin.Context) {
	token := c.MustGet(middlewares.TokenKey).(string)
	if err := a.store.DeleteToken(token); err != nil {
		internalError(c, err)
		return
	}
	c.Status(http.StatusOK)
}
