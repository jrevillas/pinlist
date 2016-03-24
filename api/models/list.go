package models

type List struct {
	ID          int64
	Name        string
	Description string
	OwnerID     int
	Pins        int
	Public      bool
}

type ListRole byte

const (
	ListOwner ListRole = 1 << iota
	ListAdder
	ListGuest
)

type UserHasList struct {
	ID     int64
	Role   ListRole
	ListID int64
	UserId int64
}
