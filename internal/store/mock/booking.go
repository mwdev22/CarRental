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

func (bs *BookingStore) Create(ctx context.Context, book *types.Booking) error {
	id := len(bs.books) + 1
	book.ID = id
	bs.books[id] = book
	return nil
}

func (bs *BookingStore) GetByID(ctx context.Context, id int) (*types.Booking, error) {
	if book, ok := bs.books[id]; ok {
		return book, nil
	}
	return nil, types.NotFound("booking")
}

func (bs *BookingStore) GetByUserID(ctx context.Context, userID int) ([]*types.Booking, error) {
	var books []*types.Booking
	for _, book := range bs.books {
		if book.UserID == userID {
			books = append(books, book)
		}
	}
	return books, nil
}

func (bs *BookingStore) Update(ctx context.Context, book *types.Booking) error {
	if _, ok := bs.books[book.ID]; !ok {
		return types.NotFound("booking")
	}
	bs.books[book.ID] = book
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
	for _, book := range bs.books {
		if book.CarID == carID {
			if (book.StartDate.Before(startDate) || book.StartDate.Equal(startDate)) && (book.EndDate.After(startDate) || book.EndDate.Equal(startDate)) {
				return false
			}
		}
	}
	return true
}
