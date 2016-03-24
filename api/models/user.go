package models

import (
	"database/sql"
	"time"

	"github.com/satori/go.uuid"

	"golang.org/x/crypto/bcrypt"

	"gopkg.in/gorp.v1"
)

type UserStatus byte

const (
	UserWaiting UserStatus = 1 << iota
	UserActive
	UserSuspended
	UserInactive
)

type User struct {
	ID        int64      `db:"id,primarykey,autoincrement" json:"id"`
	Status    UserStatus `db:"status" json:"status"`
	Username  string     `db:"username" json:"username"`
	Email     string     `db:"email" json:"email"`
	Password  string     `db:"password" json:"-"`
	CreatedAt time.Time  `db:"created_at" json:"-"`
}

func NewUser(username, email, password string) *User {
	user := &User{
		Status:   UserActive,
		Username: username,
		Email:    email,
		Password: password,
	}

	pwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		// This should never happen, so...
		panic(err)
	}
	user.Password = string(pwd)

	return user
}

type Token struct {
	ID        int64     `db:"id,primarykey,autoincrement" json:"-"`
	Hash      string    `db:"hash" json:"hash"`
	Until     time.Time `db:"until" json:"until"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	UserID    int64     `db:"user_id" json:"-"`
}

func NewToken(userID int64) *Token {
	now := time.Now()
	return &Token{
		Hash:      uuid.NewV4().String(),
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
WHERE (email = :login OR username = :login)`

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

const byTokenQuery = `SELECT u.* FROM user u
INNER JOIN token t ON t.user_id = u.id
WHERE t.hash = :hash AND t.until < :now`

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

const deleteTokenQuery = `DELETE FROM token WHERE hash = ?`

func (s UserStore) DeleteToken(hash string) error {
	_, err := s.Exec(deleteTokenQuery, hash)
	return err
}

const removeExpiredTokensQuery = `DELETE FROM token
WHERE until < ?`

func (s UserStore) RemoveExpiredTokens() error {
	_, err := s.Exec(removeExpiredTokensQuery, time.Now())
	return err
}
