package models

import (
	"fmt"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/asaskevich/govalidator"
	"github.com/go-gorp/gorp"
)

type UserStatus byte

const (
	UserWaiting UserStatus = 1 << iota
	UserActive
	UserSuspended
	UserInactive
)

type User struct {
	ID        int64      `db:"id" json:"id"`
	Status    UserStatus `db:"status" json:"status"`
	Username  string     `db:"username" json:"username"`
	Email     string     `db:"email" json:"email"`
	Password  string     `db:"password" json:"-"`
	CreatedAt time.Time  `db:"created_at" json:"-"`
}

func NewUser(username, email, password string) (*User, bool) {
	user := &User{
		Status:   UserActive,
		Username: username,
		Email:    email,
		Password: password,
	}

	if ok := user.Validate(); !ok {
		return nil, false
	}

	pwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, false
	}
	user.Password = string(pwd)

	return user, true
}

var usernameRegex = regexp.MustCompile(`[a-zA-Z0-9_]+`)

func (u *User) Validate() bool {
	return govalidator.IsEmail(u.Email) &&
		usernameRegex.MatchString(u.Username) &&
		len(u.Password) > 8
}

type Token struct {
	ID        int64     `db:"id"`
	Hash      string    `db:"hash"`
	Until     time.Time `db:"until"`
	CreatedAt time.Time `db:"created_at"`
	UserID    int64     `db:"user_id"`
}

func NewToken(userID int64) *Token {
	now := time.Now()
	hash, _ := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%d_%s", userID, now)), bcrypt.DefaultCost)
	return &Token{
		Hash:      string(hash),
		Until:     now.Add(1 * 365 * 24 * time.Hour),
		CreatedAt: now,
		UserID:    userID,
	}
}

type UserStore struct {
	*gorp.DbMap
}

const existsUserQuery = `SELECT COUNT(*) FROM user
WHERE email = :email OR username = :username`

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

const byLoginDetailsQuery = `SELECT * FROM user
WHERE (email = :login OR username :login)
AND password = :password`

func (s UserStore) ByLoginDetails(login, password string) (*User, error) {
	var user User
	err := s.SelectOne(&user, byLoginDetailsQuery, map[string]interface{}{
		"login":    login,
		"password": password,
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

const byTokenQuery = `SELECT u.* FROM user u
INNER JOIN token t ON t.user_id = u.id
WHERE t.hash = :hash AND t.until < :now`

func (s UserStore) ByToken(hash string) (*User, error) {
	var user User
	err := s.SelectOne(&user, byTokenQuery, map[string]interface{}{
		"hash": hash,
		"now":  time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

const deleteTokenQuery = `DELETE FROM token WHERE hash = :hash`

func (s UserStore) DeleteToken(hash string) error {
	_, err := s.Exec(deleteTokenQuery, map[string]interface{}{
		"hash": hash,
	})
	return err
}

const removeExpiredTokensQuery = `DELETE FROM token
WHERE until < :now`

func (s UserStore) RemoveExpiredTokens() error {
	_, err := s.Exec(removeExpiredTokensQuery, map[string]interface{}{
		"now": time.Now(),
	})
	return err
}
