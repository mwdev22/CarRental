package mock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mwdev22/CarRental/internal/types"
)

type CarRepository struct {
	mu     sync.RWMutex
	cars   map[int]types.Car
	nextID int
}

func NewCarRepository() *CarRepository {
	return &CarRepository{
		cars:   make(map[int]types.Car),
		nextID: 1,
	}
}

func (r *CarRepository) Create(ctx context.Context, car *types.Car) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	car.ID = r.nextID
	r.nextID++

	for _, existingCar := range r.cars {
		if existingCar.RegistrationNo == car.RegistrationNo {
			return fmt.Errorf("car with registration number %s already exists", car.RegistrationNo)
		}
	}

	car.Created = time.Now()
	car.Updated = time.Now()

	r.cars[car.ID] = *car
	return nil
}

func (r *CarRepository) GetByID(ctx context.Context, id int) (*types.Car, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	car, exists := r.cars[id]
	if !exists {
		return nil, fmt.Errorf("car with id %d not found", id)
	}

	return &car, nil
}

func (r *CarRepository) Update(ctx context.Context, id int, car *types.Car) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.cars[id]
	if !exists {
		return fmt.Errorf("car with id %d not found", id)
	}

	for _, existingCar := range r.cars {
		if existingCar.RegistrationNo == car.RegistrationNo {
			return fmt.Errorf("car with registration number %s already exists", car.RegistrationNo)
		}
	}

	r.cars[id] = *car
	return nil
}

func (r *CarRepository) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.cars[id]
	if !exists {
		return fmt.Errorf("car with id %d not found", id)
	}

	delete(r.cars, id)
	return nil
}

func (r *CarRepository) GetBatch(ctx context.Context, filters []*types.QueryFilter, opts *types.QueryOptions) ([]types.Car, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var cars []types.Car

	for _, car := range r.cars {
		cars = append(cars, car)
	}

	if opts != nil && opts.Limit > 0 {
		end := opts.Offset + opts.Limit
		if end > len(cars) {
			end = len(cars)
		}
		cars = cars[opts.Offset:end]
	}

	return cars, nil
}
