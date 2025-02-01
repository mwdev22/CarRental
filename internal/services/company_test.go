package services

import (
	"testing"

	"github.com/mwdev22/CarRental/internal/store/mock"
	"github.com/mwdev22/CarRental/internal/types"
)

func TestCompanyService(t *testing.T) {
	companyService := NewCompanyService(mock.NewCompanyRepository())
	companyOwnerID := 1

	t.Run("CreateCompany", func(t *testing.T) {
		tests := []struct {
			name        string
			payload     *types.CreateCompanyPayload
			expectError bool
		}{
			{
				name: "successful company creation",
				payload: &types.CreateCompanyPayload{
					Name: "testcompany",
				},
				expectError: false,
			},
			{
				name: "duplicate company name",
				payload: &types.CreateCompanyPayload{
					Name: "testcompany",
				},
				expectError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := companyService.Create(tt.payload, companyOwnerID)

				if tt.expectError && err == nil {
					t.Errorf("expected an error, got nil")
				}
				if !tt.expectError && err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			})
		}
	})

	t.Run("DeleteCompany", func(t *testing.T) {
		tests := []struct {
			name        string
			companyName string
			expectError bool
		}{
			{
				name:        "successful company deletion",
				companyName: "deletecompany",
				expectError: false,
			},
			{
				name:        "company not found",
				companyName: "nonexistentcompany",
				expectError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := companyService.Delete(1, companyOwnerID)

				if tt.expectError && err == nil {
					t.Errorf("expected an error, got nil")
				}
				if !tt.expectError && err != nil {
					t.Errorf("expected no error, got: %v", err)
				}

				if !tt.expectError {
					_, err := companyService.GetByID(1)
					if err == nil {
						t.Errorf("expected company to be deleted, but it still exists")
					}
				}
			})
		}
	})
}
