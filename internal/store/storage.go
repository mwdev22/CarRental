package store

import (
	"time"
)

type User struct {
	ID       int       `db:"id" json:"id"`
	Username string    `db:"username" json:"username"`
	Password []byte    `db:"password" json:"password"`
	Email    string    `db:"email" json:"email"`
	Created  time.Time `db:"created" json:"created"`
}

type Storage interface {
	// user methods
	CreateUser(u *User) error
	GetByUsername(username string) (*User, error)
}
