package services

import (
	"context"
	"math"
	"time"

	"github.com/mwdev22/CarRental/internal/store"
	"github.com/mwdev22/CarRental/internal/types"
)

type BookService struct {
	bookingStore store.BookingStore
	carStore     store.CarStore
	userStore    store.UserStore
}

func NewBookService(bookStore store.BookingStore, carStore store.CarStore, userStore store.UserStore) *BookService {
	return &BookService{
		bookingStore: bookStore,
		carStore:     carStore,
		userStore:    userStore,
	}
}

func (s *BookService) Create(userId int, payload *types.CreateBookingPayload) error {
	user, err := s.userStore.GetByID(context.Background(), userId)
	if err != nil {
		return types.DatabaseError(err)
	} else if user == nil {
		return types.NotFound("user")

	}

	car, err := s.carStore.GetByID(context.Background(), payload.CarID)
	if err != nil {
		return types.DatabaseError(err)
	} else if car == nil {
		return types.NotFound("car")
	}

	startDate, err := time.Parse(time.RFC3339, payload.StartDate)
	if err != nil {
		return types.InternalServerError(err.Error())
	}
	endDate, err := time.Parse(time.RFC3339, payload.EndDate)
	if err != nil {
		return types.InternalServerError(err.Error())
	}

	if startDate.After(endDate) {
		return types.BadRequest("start date cannot be after end date")
	}

	// estimate price based on days
	totalPrice := math.Ceil(endDate.Sub(startDate).Hours()) / 24 * car.PricePerDay

	book := &types.Booking{
		CarID:     payload.CarID,
		UserID:    userId,
		StartDate: startDate,
		EndDate:   endDate,
		Total:     totalPrice,
	}

	if err := s.bookingStore.Create(context.Background(), book); err != nil {
		return types.DatabaseError(err)
	}

	return nil
}

func (s *BookService) GetByID(id int) (*types.Booking, error) {
	book, err := s.bookingStore.GetByID(context.Background(), id)
	if err != nil {
		return nil, types.DatabaseError(err)
	}
	return book, nil
}

// user can only extend or shorten the booking
func (s *BookService) Update(id int, payload *types.UpdateBookingPayload) error {
	book, err := s.bookingStore.GetByID(context.Background(), id)
	if err != nil {
		return types.DatabaseError(err)
	} else if book == nil {
		return types.NotFound("booking")
	}

	car, err := s.carStore.GetByID(context.Background(), book.CarID)
	if err != nil {
		return types.DatabaseError(err)
	} else if car == nil {
		return types.NotFound("car")
	}

	startDate, err := time.Parse(time.RFC3339, payload.StartDate)
	if err != nil {
		return types.InternalServerError(err.Error())
	}

	endDate, err := time.Parse(time.RFC3339, payload.EndDate)
	if err != nil {
		return types.InternalServerError(err.Error())
	}

	if startDate.After(endDate) {
		return types.BadRequest("start date cannot be after end date")
	}

	// estimate price based on days once again
	totalPrice := math.Ceil(endDate.Sub(startDate).Hours()) / 24 * car.PricePerDay

	// update booking
	book.Total = totalPrice
	book.StartDate = startDate
	book.EndDate = endDate

	if err := s.bookingStore.Update(context.Background(), book); err != nil {
		return types.DatabaseError(err)
	}

	return nil
}

func (s *BookService) Delete(id int) error {
	if err := s.bookingStore.Delete(context.Background(), id); err != nil {
		return types.DatabaseError(err)
	}
	return nil
}
