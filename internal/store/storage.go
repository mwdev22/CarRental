package store

import (
	"context"
	"time"

	"github.com/mwdev22/FileStorage/internal/types"
)

type User struct {
	ID       int            `db:"id" json:"id"`
	Username string         `db:"username" json:"username"`
	Password []byte         `db:"password" json:"-"`
	Email    string         `db:"email" json:"email"`
	Role     types.UserRole `db:"role" json:"role"`
	Created  time.Time      `db:"created" json:"-"`
}

type UserStore interface {
	// user methods
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id int) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, u *User) error
}
