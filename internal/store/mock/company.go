package mock

import (
	"context"
	"sync"
	"time"

	"github.com/mwdev22/CarRental/internal/types"
)

type CompanyRepository struct {
	companies map[int]types.Company
	mu        sync.RWMutex
	nextID    int
}

func NewCompanyRepository() *CompanyRepository {
	return &CompanyRepository{
		companies: make(map[int]types.Company),
		nextID:    1,
	}
}

func (r *CompanyRepository) Create(ctx context.Context, company *types.Company) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	company.ID = r.nextID
	r.nextID++
	company.Created = time.Now()
	company.Updated = company.Created

	r.companies[company.ID] = *company
	return nil
}

func (r *CompanyRepository) GetByID(ctx context.Context, id int) (*types.Company, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	company, exists := r.companies[id]
	if !exists {
		return nil, types.NotFound("company not found")
	}
	return &company, nil
}

func (r *CompanyRepository) Update(ctx context.Context, company *types.Company) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existingCompany, exists := r.companies[company.ID]
	if !exists {
		return types.NotFound("company not found")
	}

	existingCompany.Name = company.Name
	existingCompany.Email = company.Email
	existingCompany.Phone = company.Phone
	existingCompany.Address = company.Address
	existingCompany.Updated = time.Now()

	r.companies[company.ID] = existingCompany
	return nil
}

func (r *CompanyRepository) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.companies[id]; !exists {
		return types.NotFound("company not found")
	}

	delete(r.companies, id)
	return nil
}

func (r *CompanyRepository) GetBatch(ctx context.Context, filters []*types.QueryFilter, opts *types.QueryOptions) ([]types.Company, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var companies []types.Company

	for _, c := range r.companies {
		companies = append(companies, c)
	}

	if opts != nil && opts.Limit > 0 && opts.Offset < len(companies) {
		end := opts.Offset + opts.Limit
		if end > len(companies) {
			end = len(companies)
		}
		companies = companies[opts.Offset:end]
	}

	return companies, nil
}
