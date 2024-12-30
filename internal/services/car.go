package services

import "github.com/mwdev22/FileStorage/internal/store"

type CarService struct {
	carRepo store.CarRepository
}

func NewCarService(carRepo store.CarRepository) *CarService {
	return &CarService{
		carRepo: carRepo,
	}
}
