package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mwdev22/CarRental/internal/types"
	"github.com/mwdev22/CarRental/internal/utils"
)

type CarRepositorySQL struct {
	DB *sqlx.DB
}

func NewCarRepository(db *sqlx.DB) *CarRepositorySQL {
	return &CarRepositorySQL{
		DB: db,
	}
}

func (r *CarRepositorySQL) Create(ctx context.Context, car *types.Car) error {
	query := `INSERT INTO car (company_id, make, model, year, color, registration_no, price_per_day) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.DB.Exec(query, car.CompanyID, car.Make, car.Model, car.Year, car.Color, car.RegistrationNo, car.PricePerDay)
	if err != nil {
		return err
	}
	return nil
}

func (r *CarRepositorySQL) GetByID(ctx context.Context, id int) (*types.Car, error) {
	var car types.Car
	query := `SELECT id, make, model, year, color, registration_no, price_per_day, created_at, updated_at FROM car WHERE id = $1`
	err := r.DB.Get(&car, query, id)
	if err != nil {
		return nil, err
	}
	return &car, nil
}

func (r *CarRepositorySQL) Update(ctx context.Context, id int, car *types.Car) error {
	query := `UPDATE car SET make = $1, model = $2, year = $3, color = $4, registration_no = $5, price_per_day = $6, updated = CURRENT_TIMESTAMP WHERE id = $7`
	_, err := r.DB.Exec(query, car.Make, car.Model, car.Year, car.Color, car.RegistrationNo, car.PricePerDay, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *CarRepositorySQL) GetBatch(ctx context.Context, filters []*types.QueryFilter, opts *types.QueryOptions) ([]types.Car, error) {
	query := `SELECT id, make, model, year, color, registration_no, price_per_day, created_at, updated_at FROM car WHERE 1 = 1`

	query, args := utils.BuildBatchQuery(query, filters, opts)
	query = r.DB.Rebind(query)

	var cars []types.Car
	err := r.DB.Select(&cars, query, args...)
	if err != nil {
		return nil, err
	}
	return cars, nil
}

func (r *CarRepositorySQL) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM car WHERE id = $1`
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
