package store

import (
	"context"
	"time"

	"github.com/mwdev22/CarRental/internal/types"
)

type UserStore interface {
	// user methods
	Create(ctx context.Context, u *types.User) error
	GetByID(ctx context.Context, id int) (*types.User, error)
	GetByUsername(ctx context.Context, username string) (*types.User, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, u *types.User) error
}

type CompanyStore interface {
	Create(ctx context.Context, c *types.Company) error
	GetByID(ctx context.Context, id int) (*types.Company, error)
	Update(ctx context.Context, c *types.Company) error
	Delete(ctx context.Context, id int) error
	GetBatch(ctx context.Context, filters []*types.QueryFilter, opts *types.QueryOptions) ([]types.Company, error)
}

type CarStore interface {
	Create(ctx context.Context, car *types.Car) error
	GetByID(ctx context.Context, id int) (*types.Car, error)
	Update(ctx context.Context, id int, car *types.Car) error
	Delete(ctx context.Context, id int) error
	GetBatch(ctx context.Context, filters []*types.QueryFilter, opts *types.QueryOptions) ([]types.Car, error)
}

type BookingStore interface {
	Create(ctx context.Context, booking *types.Booking) error
	GetByID(ctx context.Context, id int) (*types.Booking, error)
	Update(ctx context.Context, booking *types.Booking) error
	Delete(ctx context.Context, id int) error
	GetByUserID(ctx context.Context, userID int) ([]*types.Booking, error)
	CheckDateAvailability(ctx context.Context, carID int, startDate, endDate time.Time) bool
}
