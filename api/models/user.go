package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mvader/pinlist/api/log"
	"github.com/satori/go.uuid"

	"golang.org/x/crypto/bcrypt"

	"gopkg.in/gorp.v1"
)

// UserStatus is the current status of the user account.
type UserStatus byte

const (
	// UserWaiting means the user account is registered but not active.
	UserWaiting UserStatus = 1 << iota
	// UserActive means the user account is fully active.
	UserActive
	// UserSuspended means the user was suspended by the system.
	UserSuspended
	// UserInactive means the user deleted their own account.
	UserInactive
)

// User defines the model of an user account.
type User struct {
	ID        int64      `db:"id" json:"id"`
	Status    UserStatus `db:"status" json:"status"`
	Username  string     `db:"username" json:"username"`
	Email     string     `db:"email" json:"email"`
	Password  string     `db:"password" json:"-"`
	CreatedAt time.Time  `db:"created_at" json:"-"`
}

// BasicUser are the displayable fields of user.
type BasicUser struct {
	Username string     `db:"username"`
	Email    string     `db:"email"`
	Status   UserStatus `db:"status"`
}

// NewUser creates a new user with the username, email AND
// password. The resultant user already has the password
// crypted using BCrypt. This function can panic if the BCrypt
// crypt returns an error, which shouldn't happen.
func NewUser(username, email, password string) *User {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// This should never happen, so...
		log.Err(err)
		panic(err)
	}

	user := &User{
		Status:   UserActive,
		Username: username,
		Email:    email,
		Password: string(pwd),
	}

	return user
}

// Token is an code that grants access to the app to a
// specific user for a certain amount of time.
type Token struct {
	ID        int64     `db:"id" json:"-"`
	Hash      string    `db:"hash" json:"hash"`
	Until     time.Time `db:"until" json:"until"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	UserID    int64     `db:"user_id" json:"-"`
}

// NewToken creates a new token for the given user. By default
// the expiration time of a token is a year.
func NewToken(userID int64) *Token {
	now := time.Now()
	return &Token{
		Hash:      uuid.NewV4().String(),
		Until:     now.Add(1 * 365 * 24 * time.Hour),
		CreatedAt: now,
		UserID:    userID,
	}
}

// UserStore is the service to execute operations about users
// that require the database.
type UserStore struct {
	*gorp.DbMap
}

const existsUserQuery = `SELECT COUNT(*) FROM "user"
WHERE email = :email OR username = :username`

// ExistsUser reports if there is an user with the email and
// username supplied.
func (s UserStore) ExistsUser(email, username string) (bool, error) {
	n, err := s.SelectInt(existsUserQuery, map[string]interface{}{
		"email":    email,
		"username": username,
	})
	if err != nil {
		return false, err
	}

	return n > 0, nil
}

const byLoginDetailsQuery = `SELECT * FROM "user"
WHERE (email = :login OR username = :login)`

// ByLoginDetails returns the user whose login and password matches
// the given ones. By login, it means, either email or username.
func (s UserStore) ByLoginDetails(login, password string) (*User, error) {
	var user User
	err := s.SelectOne(&user, byLoginDetailsQuery, map[string]interface{}{"login": login})
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil
	}

	return &user, nil
}

const byTokenQuery = `SELECT u.* FROM "user" u
INNER JOIN token t ON t.user_id = u.id
WHERE t.hash = :hash AND t.until > :now`

// ByToken retrieves an user that has an active token with the given hash.
func (s UserStore) ByToken(hash string) (*User, error) {
	var user User
	err := s.SelectOne(&user, byTokenQuery, map[string]interface{}{
		"hash": hash,
		"now":  time.Now(),
	})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

const deleteTokenQuery = `DELETE FROM token WHERE hash = %s`

// DeleteToken removes the token with the given hash from
// the database.
func (s UserStore) DeleteToken(hash string) error {
	q := fmt.Sprintf(deleteTokenQuery, s.Dialect.BindVar(0))
	_, err := s.Exec(q, hash)
	return err
}

const removeExpiredTokensQuery = `DELETE FROM token
WHERE until < %s`

// RemoveExpiredTokens removes all expired tokens from the database.
func (s UserStore) RemoveExpiredTokens() error {
	q := fmt.Sprintf(removeExpiredTokensQuery, s.Dialect.BindVar(0))
	_, err := s.Exec(q, time.Now())
	return err
}
