package mock

import (
	"context"

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

func (s *BookingStore) Create(ctx context.Context, book *types.Booking) error {
	id := len(s.books) + 1
	book.ID = id
	s.books[id] = book
	return nil
}

func (s *BookingStore) GetByID(ctx context.Context, id int) (*types.Booking, error) {
	if book, ok := s.books[id]; ok {
		return book, nil
	}
	return nil, types.NotFound("booking")
}

func (s *BookingStore) Update(ctx context.Context, book *types.Booking) error {
	if _, ok := s.books[book.ID]; !ok {
		return types.NotFound("booking")
	}
	s.books[book.ID] = book
	return nil
}

func (s *BookingStore) Delete(ctx context.Context, id int) error {
	if _, ok := s.books[id]; !ok {
		return types.NotFound("booking")
	}
	delete(s.books, id)
	return nil
}
