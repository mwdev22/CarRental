package services

import (
	"context"
	"fmt"
	"log"

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
	company := &types.Company{
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

func (s *CompanyService) GetByID(id int) (*types.Company, error) {
	company, err := s.companyStore.GetByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return company, nil
}

func (s *CompanyService) Update(id int, userId int, payload *types.UpdateCompanyPayload) error {
	company, err := s.companyStore.GetByID(context.Background(), id)
	if err != nil {
		return err
	}

	if company.OwnerID != userId {
		return types.Unauthorized(fmt.Sprintf("user isnt the owner of the company with id %v", id))
	}

	company.Name = payload.Name
	company.Email = payload.Email
	company.Phone = payload.Phone
	company.Address = payload.Address

	if err := s.companyStore.Update(context.Background(), company); err != nil {
		return err
	}

	return nil
}

func (cs *CompanyService) Delete(id int, userId int) error {
	company, err := cs.companyStore.GetByID(context.Background(), id)
	if err != nil {
		return err
	}

	if company.OwnerID != userId {
		return types.Unauthorized(fmt.Sprintf("user isnt the owner of the company with id %v", id))
	}

	if err := cs.companyStore.Delete(context.Background(), id); err != nil {
		return err
	}

	return nil
}

func (s *CompanyService) GetBatch(filters []*types.QueryFilter, opts *types.QueryOptions) ([]types.Company, error) {
	for _, filter := range filters {
		log.Println(filter)
	}
	companies, err := s.companyStore.GetBatch(context.Background(), filters, opts)
	if err != nil {
		return nil, types.DatabaseError(fmt.Errorf("failed to get companies, %v", err))
	}

	return companies, nil
}
