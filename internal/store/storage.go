package store

import (
	"context"

	"github.com/mwdev22/CarRental/internal/types"
)

type UserStore interface {
	// user methods
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id int) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, u *User) error
}

type CompanyStore interface {
	Create(ctx context.Context, c *Company) error
	GetByID(ctx context.Context, id int) (*Company, error)
	Update(ctx context.Context, id int, c *Company) error
	GetBatch(ctx context.Context, filters []*types.QueryFilter, opts *types.QueryOptions) ([]Company, error)
}

type CarStore interface {
	Create(ctx context.Context, car *Car) error
	GetByID(ctx context.Context, id int) (*Car, error)
	Update(ctx context.Context, id int, car *Car) error
	GetBatch(ctx context.Context, filters []*types.QueryFilter, opts *types.QueryOptions) ([]Car, error)
}
