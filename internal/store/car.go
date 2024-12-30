package store

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type CarRepository struct {
	DB *sqlx.DB
}

func NewCarRepo(db *sqlx.DB) *CarRepository {
	return &CarRepository{
		DB: db,
	}
}

func (r *CarRepository) Create(ctx context.Context, car *Car) error {
	query := `INSERT INTO cars (make, model, year, created) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, car.Make, car.Model, car.Year, car.Created)
	if err != nil {
		return err
	}
	return nil
}

func (r *CarRepository) GetByID(ctx context.Context, id int) (*Car, error) {
	var car Car
	query := `SELECT id, make, model, year, color, registration_no, price_per_day, created_at, updated_at FROM cars WHERE id = $1`
	err := r.DB.Get(&car, query, id)
	if err != nil {
		return nil, err
	}
	return &car, nil
}

func (r *CarRepository) Update(ctx context.Context, car *Car) error {
	query := `UPDATE cars SET make = $1, model = $2, year = $3, color = $4, registration_no = $5, price_per_day = $6, updated = CURRENT_TIMESTAMP WHERE id = $7`
	_, err := r.DB.Exec(query, car.Make, car.Model, car.Year, car.Color, car.RegistrationNo, car.PricePerDay, car.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CarRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM cars WHERE id = $1`
	res, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("car with provided id not found")
	}

	return nil
}
