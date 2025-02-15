package services

import (
	"context"
	"testing"

	"github.com/mwdev22/CarRental/internal/store/mock"
	"github.com/mwdev22/CarRental/internal/types"
	"github.com/mwdev22/CarRental/internal/utils"
)

func TestBookingService(t *testing.T) {
	bookingService := NewBookingService(mock.NewBookingStore(), mock.NewCarRepository(), mock.NewUserRepository())

	for i := 1; i <= 5; i++ {
		err := bookingService.userStore.Create(context.Background(), &types.User{
			ID:       i,
			Username: utils.GenerateUniqueString("user"),
			Password: []byte(utils.GenerateUniqueString("password")),
			Email:    utils.GenerateUniqueString("email@bllabal.com"),
			Role:     types.UserTypeUser,
		})
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}
	}

	for i := 1; i <= 5; i++ {
		err := bookingService.carStore.Create(context.Background(), &types.Car{
			ID:             i,
			Make:           utils.GenerateUniqueString("make"),
			Model:          utils.GenerateUniqueString("model"),
			Year:           2025,
			Color:          utils.GenerateUniqueString("color"),
			RegistrationNo: utils.GenerateUniqueString("regno"),
			PricePerDay:    100,
			CompanyID:      1,
		})
		if err != nil {
			t.Fatalf("failed to create car: %v", err)
		}
	}

	t.Run("CreateBooking", func(t *testing.T) {
		tests := []struct {
			name        string
			payload     *types.CreateBookingPayload
			expectError bool
		}{
			{
				name: "successful booking",
				payload: &types.CreateBookingPayload{
					CarID:     1,
					StartDate: "2025-01-01",
					EndDate:   "2025-01-04",
				},
				expectError: false,
			},
			{
				name: "successful booking",
				payload: &types.CreateBookingPayload{
					CarID:     2,
					StartDate: "2025-01-01",
					EndDate:   "2025-01-12",
				},
				expectError: false,
			},
			{
				name: "successful booking",
				payload: &types.CreateBookingPayload{
					CarID:     3,
					StartDate: "2025-01-05",
					EndDate:   "2025-01-18",
				},
				expectError: false,
			},
			{
				name: "invalid booking",
				payload: &types.CreateBookingPayload{
					CarID:     1,
					StartDate: "2025-01-02",
					EndDate:   "2025-01-01",
				},
				expectError: true,
			},
			{
				name: "car already boked",
				payload: &types.CreateBookingPayload{
					CarID:     1,
					StartDate: "2025-01-02",
					EndDate:   "2025-01-03",
				},
				expectError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := bookingService.Create(1, tt.payload)

				if tt.expectError && err == nil {
					t.Errorf("expected an error, got nil")
				}
				if !tt.expectError && err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			})
		}
	})

	t.Run("GetBooking", func(t *testing.T) {
		book, err := bookingService.GetByID(1)
		if err != nil {
			t.Fatalf("failed to get booking: %v", err)
		}

		if book == nil {
			t.Fatalf("expected booking, got nil")
		}
	})

	t.Run("GetUserBookings", func(t *testing.T) {
		books, err := bookingService.GetByUserID(1)
		if err != nil {
			t.Fatalf("failed to get user bookings: %v", err)
		}

		if len(books) != 3 {
			t.Fatalf("expected 3 bookings, got %v", len(books))
		}
	})

	t.Run("UpdateBooking", func(t *testing.T) {
		err := bookingService.Update(1, &types.UpdateBookingPayload{
			StartDate: "2025-01-01",
			EndDate:   "2025-01-10",
		})
		if err != nil {
			t.Fatalf("failed to update booking: %v", err)
		}

		book, err := bookingService.GetByID(1)
		if err != nil {
			t.Fatalf("failed to get booking: %v", err)
		}

		if book.Total != 1000 {
			t.Fatalf("expected total price 1000, got %v", book.Total)
		}

	})
}
