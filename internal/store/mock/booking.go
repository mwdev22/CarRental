package mock

import (
	"context"
	"time"

	"github.com/mwdev22/CarRental/internal/types"
)

type BookingStore struct {
	books map[int]*types.Booking
}

func NewBookingStore() *BookingStore {
	return &BookingStore{
		books: make(map[int]*types.Booking),
	}
}

func (bs *BookingStore) Create(ctx context.Context, booking *types.Booking) error {
	id := len(bs.books) + 1
	booking.ID = id
	bs.books[id] = booking
	return nil
}

func (bs *BookingStore) GetByID(ctx context.Context, id int) (*types.Booking, error) {
	if booking, ok := bs.books[id]; ok {
		return booking, nil
	}
	return nil, types.NotFound("booking")
}

func (bs *BookingStore) GetByUserID(ctx context.Context, userID int) ([]*types.Booking, error) {
	var books []*types.Booking
	for _, booking := range bs.books {
		if booking.UserID == userID {
			books = append(books, booking)
		}
	}
	return books, nil
}

func (bs *BookingStore) Update(ctx context.Context, booking *types.Booking) error {
	if _, ok := bs.books[booking.ID]; !ok {
		return types.NotFound("booking")
	}
	bs.books[booking.ID] = booking
	return nil
}

func (bs *BookingStore) Delete(ctx context.Context, id int) error {
	if _, ok := bs.books[id]; !ok {
		return types.NotFound("booking")
	}
	delete(bs.books, id)
	return nil
}

func (bs *BookingStore) CheckDateAvailability(ctx context.Context, carID int, startDate, endDate time.Time) bool {
	for _, booking := range bs.books {
		if booking.CarID == carID {
			if (booking.StartDate.Before(startDate) || booking.StartDate.Equal(startDate)) && (booking.EndDate.After(startDate) || booking.EndDate.Equal(startDate)) {
				return false
			}
		}
	}
	return true
}
