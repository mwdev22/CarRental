package store

import (
	"context"
)

type UserStore interface {
	// user methods
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id int) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, u *User) error
}
