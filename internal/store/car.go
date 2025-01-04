package store

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mwdev22/FileStorage/internal/types"
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
	query := `INSERT INTO cars (make, model, year, color, registration_no, price_per_day) VALUES ($1, $2, $3, $4, $5, $6)`
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

func (r *CarRepository) Update(ctx context.Context, id int, car *Car) error {
	query := `UPDATE cars SET make = $1, model = $2, year = $3, color = $4, registration_no = $5, price_per_day = $6, updated = CURRENT_TIMESTAMP WHERE id = $7`
	_, err := r.DB.Exec(query, car.Make, car.Model, car.Year, car.Color, car.RegistrationNo, car.PricePerDay, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *CarRepository) GetBatch(ctx context.Context, filters []*types.QueryFilter, opts *types.QueryOptions) ([]Car, error) {
	query := `SELECT id, make, model, year, color, registration_no, price_per_day, created_at, updated_at FROM cars WHERE 1 = 1`
	args := make([]interface{}, 0)

	// i use question marks because filters and their count are dynamic
	// could be 1, 2, 3 etc through loop, but its not necessary complexity i think
	for _, filter := range filters {
		query += fmt.Sprintf(" AND %s %s ?", filter.Field, filter.Operator)
		args = append(args, filter.Value)
	}

	if opts != nil {
		query += ` ORDER BY ? ?`
		query += ` LIMIT ? OFFSET ?`
		args = append(args, opts.SortField, opts.SortDiretion, opts.Limit, opts.Offset)
	} else {
		query += ` ORDER BY id DESC`
	}

	// rebind query to postgres syntax
	query = r.DB.Rebind(query)

	var cars []Car
	err := r.DB.Select(&cars, query, args...)
	if err != nil {
		return nil, err
	}
	return cars, nil
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
