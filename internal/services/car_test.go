package services

import (
	"testing"

	"github.com/mwdev22/CarRental/internal/store/mock"
	"github.com/mwdev22/CarRental/internal/types"
)

func TestCar(t *testing.T) {
	carService := NewCarService(mock.NewCarRepository())

	payloads := make([]*types.CreateCarPayload, 0)
	for i := 0; i < 10; i++ {
		payloads = append(payloads, &types.CreateCarPayload{
			Make:           "Toyota",
			Model:          "Corolla",
			Year:           2021,
			Color:          "Red",
			RegistrationNo: "ABC123",
			PricePerDay:    100,
			CompanyID:      i,
		},
		)
	}

	for _, p := range payloads {
		err := carService.CreateCar(p)
		if err != nil {
			t.Logf("failed creating car: %s", err)
		}
	}

	expectedId := 4
	car, err := carService.GetByID(expectedId)
	if err != nil {
		t.Errorf("error getting car by id: %s", err)
	}

	if car.ID != 4 {
		t.Errorf("expected car with id: %v, got %v", expectedId, car.ID)
	}

}
