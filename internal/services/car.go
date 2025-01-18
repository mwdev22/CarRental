package services

import (
	"context"
	"fmt"

	"github.com/mwdev22/CarRental/internal/store"
	"github.com/mwdev22/CarRental/internal/types"
)

type CarService struct {
	carStore store.CarStore
}

func NewCarService(carStore store.CarStore) *CarService {
	return &CarService{
		carStore: carStore,
	}
}

func (s *CarService) CreateCar(payload *types.CreateCarPayload) error {
	car := &types.Car{
		Make:           payload.Make,
		Model:          payload.Model,
		Year:           payload.Year,
		Color:          payload.Color,
		RegistrationNo: payload.RegistrationNo,
		PricePerDay:    payload.PricePerDay,
		CompanyID:      payload.CompanyID,
	}
	if err := s.carStore.Create(context.Background(), car); err != nil {
		return types.DatabaseError(fmt.Errorf("failed to create car: %v", err))
	}

	return nil
}

func (s *CarService) GetByID(id int) (*types.Car, error) {
	car, err := s.carStore.GetByID(context.Background(), id)
	if err != nil {
		return nil, types.DatabaseError(fmt.Errorf("failed to get car by id: %v", err))
	}
	return car, nil
}

func (s *CarService) UpdateCar(id int, payload *types.CreateCarPayload) error {
	car, err := s.carStore.GetByID(context.Background(), id)
	if err != nil {
		return types.DatabaseError(fmt.Errorf("failed to get car by id: %v", err))
	}

	if err := s.carStore.Update(context.Background(), id, car); err != nil {
		return types.DatabaseError(fmt.Errorf("failed to update car: %v", err))
	}

	return nil
}

func (s *CarService) GetBatch(filters []*types.QueryFilter, opts *types.QueryOptions) ([]types.Car, error) {
	cars, err := s.carStore.GetBatch(context.Background(), filters, opts)
	if err != nil {
		return nil, types.DatabaseError(fmt.Errorf("failed to get batch of cars: %v", err))
	}
	return cars, nil
}
