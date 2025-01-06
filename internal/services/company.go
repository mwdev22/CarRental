package services

import (
	"context"
	"fmt"

	"github.com/mwdev22/CarRental/internal/store"
	"github.com/mwdev22/CarRental/internal/types"
)

type CompanyService struct {
	companyStore store.CompanyStore
}

func NewCompanyService(companyStore store.CompanyStore) *CompanyService {
	return &CompanyService{
		companyStore: companyStore,
	}
}

func (s *CompanyService) Create(payload *types.CreateCompanyPayload, ownerId int) error {
	company := &store.Company{
		Name:    payload.Name,
		OwnerID: ownerId,
		Email:   payload.Email,
		Phone:   payload.Phone,
		Address: payload.Address,
	}

	if err := s.companyStore.Create(context.Background(), company); err != nil {
		return types.DatabaseError(fmt.Errorf("failed to create company"))
	}

	return nil
}

func (s *CompanyService) GetByID(id int) (*store.Company, error) {
	company, err := s.companyStore.GetByID(context.Background(), id)
	if err != nil {
		return nil, types.DatabaseError(fmt.Errorf("failed to get company by id"))
	}

	return company, nil
}
