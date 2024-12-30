package store

import "database/sql"

type CompanyRepository struct {
	DB *sql.DB
}

func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{
		DB: db,
	}
}
