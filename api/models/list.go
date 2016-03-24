package models

type List struct {
	ID      int64          `db:"id,primarykey,autoincrement" json:"id"`
	Name    string         `db:"name" json:"name"`
	OwnerID int            `db:"owner_id" json:"-"`
	Owner   *User          `db:"-" json:"owner"`
	Pins    int            `db:"pins" json:"pins"`
	Public  bool           `db:"public" json:"public"`
	Users   []*UserHasList `db:"-" json:"users"`
}

type ListRole byte

const (
	ListOwner ListRole = 1 << iota
	ListAdder
	ListGuest
)

type UserHasList struct {
	ID     int64    `id:"id,primarykey,autoincrement" json:"-"`
	Role   ListRole `db:"role" json:"role"`
	ListID int64    `db:"list_id" json:"-"`
	UserID int64    `db:"user_id" json:"-"`
	User   *User    `db:"-" json:"user"`
}
