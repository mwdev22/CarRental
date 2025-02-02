package services

import (
	"fmt"
	"testing"

	"github.com/mwdev22/CarRental/internal/store/mock"
	"github.com/mwdev22/CarRental/internal/types"
)

func TestCarService(t *testing.T) {
	carService := NewCarService(mock.NewCarRepository())

	t.Run("CreateCar", func(t *testing.T) {
		tests := []struct {
			name        string
			payload     *types.CreateCarPayload
			expectError bool
		}{
			{
				name: "successful car creation",
				payload: &types.CreateCarPayload{
					Make:           "Toyota",
					Model:          "Corolla",
					Year:           2021,
					Color:          "Red",
					RegistrationNo: "ABC123",
					PricePerDay:    100,
					CompanyID:      1,
				},
				expectError: false,
			},
			{
				name: "duplicate registration number",
				payload: &types.CreateCarPayload{
					Make:           "Toyota",
					Model:          "Corolla",
					Year:           2021,
					Color:          "Red",
					RegistrationNo: "ABC123",
					PricePerDay:    100,
					CompanyID:      1,
				},
				expectError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := carService.CreateCar(tt.payload)

				if tt.expectError && err == nil {
					t.Errorf("expected an error, got nil")
				}
				if !tt.expectError && err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			})
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		createPayload := &types.CreateCarPayload{
			Make:           "Honda",
			Model:          "Civic",
			Year:           2020,
			Color:          "Blue",
			RegistrationNo: "XYZ789",
			PricePerDay:    80,
			CompanyID:      2,
		}
		err := carService.CreateCar(createPayload)
		if err != nil {
			t.Fatalf("failed to create car: %v", err)
		}

		tests := []struct {
			name        string
			carID       int
			expectError bool
		}{
			{
				name:        "successful car retrieval",
				carID:       1,
				expectError: false,
			},
			{
				name:        "car not found",
				carID:       999,
				expectError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				car, err := carService.GetByID(tt.carID)

				if tt.expectError && err == nil {
					t.Errorf("expected an error, got nil")
				}
				if !tt.expectError && err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
				if !tt.expectError && car.ID != tt.carID {
					t.Errorf("expected car ID: %d, got: %d", tt.carID, car.ID)
				}
			})
		}
	})

	t.Run("UpdateCar", func(t *testing.T) {
		createPayload := &types.CreateCarPayload{
			Make:           "Ford",
			Model:          "Mustang",
			Year:           2019,
			Color:          "Black",
			RegistrationNo: "MUS123",
			PricePerDay:    150,
			CompanyID:      3,
		}
		err := carService.CreateCar(createPayload)
		if err != nil {
			t.Fatalf("failed to create car: %v", err)
		}

		updatePayload := &types.UpdateCarPayload{
			Make:           "Ford",
			Model:          "Mustang GT",
			Year:           2020,
			Color:          "Blue",
			RegistrationNo: "MUS223",
			PricePerDay:    200,
		}

		tests := []struct {
			name        string
			carID       int
			payload     *types.UpdateCarPayload
			expectError bool
		}{
			{
				name:        "successful car update",
				carID:       3,
				payload:     updatePayload,
				expectError: false,
			},
			{
				name:        "car not found",
				carID:       999,
				payload:     updatePayload,
				expectError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := carService.UpdateCar(tt.carID, tt.payload)

				if tt.expectError && err == nil {
					t.Errorf("expected an error, got nil")
				}
				if !tt.expectError && err != nil {
					t.Errorf("expected no error, got: %v", err)
				}

				if !tt.expectError {
					car, err := carService.GetByID(tt.carID)
					if err != nil {
						t.Fatalf("failed to get car: %v", err)
					}
					if car.Model != tt.payload.Model {
						t.Errorf("expected model: %s, got: %s", tt.payload.Model, car.Model)
					}
					if car.PricePerDay != tt.payload.PricePerDay {
						t.Errorf("expected price per day: %f, got: %f", tt.payload.PricePerDay, car.PricePerDay)
					}
				}
			})
		}
	})

	t.Run("GetBatch", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			payload := &types.CreateCarPayload{
				Make:           "Tesla",
				Model:          "Model 3",
				Year:           2022,
				Color:          "White",
				RegistrationNo: fmt.Sprintf("TES%d", i),
				PricePerDay:    300,
				CompanyID:      4,
			}
			err := carService.CreateCar(payload)
			if err != nil {
				t.Fatalf("failed to create car: %v", err)
			}
		}

		tests := []struct {
			name        string
			filters     []*types.QueryFilter
			opts        *types.QueryOptions
			expectCount int
			expectError bool
		}{
			{
				name: "successful batch retrieval",
				filters: []*types.QueryFilter{
					{Field: "company_id", Value: 4},
				},
				opts:        &types.QueryOptions{Limit: 5},
				expectCount: 5,
				expectError: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				cars, err := carService.GetBatch(tt.filters, tt.opts)

				if tt.expectError && err == nil {
					t.Errorf("expected an error, got nil")
				}
				if !tt.expectError && err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
				if len(cars) != tt.expectCount {
					t.Errorf("expected %d cars, got: %d", tt.expectCount, len(cars))
				}
			})
		}
	})

	t.Run("DeleteCar", func(t *testing.T) {
		tests := []struct {
			name        string
			companyName string
			expectError bool
		}{
			{
				name:        "successful car deletion",
				expectError: false,
			},
			{
				name:        "car not found",
				expectError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := carService.Delete(1)

				if tt.expectError && err == nil {
					t.Errorf("expected an error, got nil")
				}
				if !tt.expectError && err != nil {
					t.Errorf("expected no error, got: %v", err)
				}

				if !tt.expectError {
					_, err := carService.GetByID(1)
					if err == nil {
						t.Errorf("expected car to be deleted, but it still exists")
					}
				}
			})
		}
	})
}
