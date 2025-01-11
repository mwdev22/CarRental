package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/mwdev22/CarRental/internal/types"
	"github.com/mwdev22/CarRental/internal/utils"
)

type CompanyRepository struct {
	DB *sqlx.DB
}

func NewCompanyRepository(db *sqlx.DB) *CompanyRepository {
	return &CompanyRepository{
		DB: db,
	}
}

func (r *CompanyRepository) Create(ctx context.Context, company *Company) error {
	query := `INSERT INTO company (owner_id, name, email, phone, address) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.DB.Exec(query, company.OwnerID, company.Name, company.Email, company.Phone, company.Address)
	if err != nil {
		return err
	}

	return nil
}

func (r *CompanyRepository) GetByID(ctx context.Context, id int) (*Company, error) {
	var company Company
	query := `SELECT id, owner_id, name, email, phone, address FROM company WHERE id = $1`

	err := r.DB.Get(&company, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.NotFound("company not found")
		}
		return nil, err
	}

	return &company, nil
}

func (r *CompanyRepository) Update(ctx context.Context, company *Company) error {
	query := `UPDATE company SET name = $1, email = $2, phone = $3, address = $4 WHERE id = $5`

	rows, err := r.DB.Exec(query, company.Name, company.Email, company.Phone, company.Address, company.ID)
	if err != nil {
		return err
	}

	if count, _ := rows.RowsAffected(); count == 0 {
		return types.NotFound("company not found")
	}

	return nil
}

func (r *CompanyRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM company WHERE id = $1`

	rows, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	if count, _ := rows.RowsAffected(); count == 0 {
		return types.NotFound("company not found")
	}

	return nil
}

func (r *CompanyRepository) GetBatch(ctx context.Context, filters []*types.QueryFilter, opts *types.QueryOptions) ([]Company, error) {
	query := `SELECT id, owner_id, name, email, phone, address FROM company WHERE 1 = 1`

	query, args := utils.BuildBatchQuery(query, filters, opts)
	query = r.DB.Rebind(query)

	var companies []Company
	rows, err := r.DB.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var company Company
		err := rows.StructScan(&company)
		if err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}

	return companies, nil
}
