package store

import (
	"context"

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
		return nil, err
	}

	return &company, nil
}

func (r *CompanyRepository) Update(ctx context.Context, id int, company *Company) error {
	query := `UPDATE company SET name = $1, email = $2, phone = $3, address = $4 WHERE id = $5`

	_, err := r.DB.Exec(query, company.Name, company.Email, company.Phone, company.Address, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *CompanyRepository) GetBatch(ctx context.Context, filters []*types.QueryFilter, opts *types.QueryOptions) ([]Company, error) {
	query := `SELECT id, owner_id, name, email, phone, address FROM company WHERE 1 = 1`

	query, args := utils.BuildBatchQuery(query, filters, opts)
	query = r.DB.Rebind(query)

	var companies []Company
	err := r.DB.Select(&companies, query, args...)
	if err != nil {
		return nil, err
	}

	return companies, nil
}
