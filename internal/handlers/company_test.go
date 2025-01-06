package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/mwdev22/CarRental/internal/store"
	"github.com/mwdev22/CarRental/internal/types"
	"github.com/mwdev22/CarRental/internal/utils"
)

// helper function to create a companies
func generateCompanies(count int, t *testing.T) ([]byte, []*types.CreateCompanyPayload) {
	payloads := make([]*types.CreateCompanyPayload, 0)
	var body []byte
	for i := 0; i < count; i++ {
		payload := &types.CreateCompanyPayload{
			Name:    utils.GenerateUniqueString(fmt.Sprintf("company%v", i)),
			Email:   utils.GenerateUniqueString("company_email"),
			Phone:   utils.GenerateUniqueString("48"),
			Address: utils.GenerateUniqueString("company_address"),
		}
		url := testServer.URL + "/company"

		// Send PUT request to update company
		resp := sendPostRequest(url, payload, t)
		body = checkResponse(resp, http.StatusOK, t)
		payloads = append(payloads, payload)
	}
	return body, payloads
}

func TestCreateCompany(t *testing.T) {
	TestRegister(t)
	TestLogin(t)

	body, _ := generateCompanies(1, t)

	var responseBody map[string]string
	if err := json.Unmarshal(body, &responseBody); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}
}

func TestGetCompanyByID(t *testing.T) {
	url := testServer.URL + "/company/1"
	resp := sendGetRequest(url, t)
	body := checkResponse(resp, http.StatusOK, t)

	var company store.Company
	if err := json.Unmarshal(body, &company); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}
}

func TestUpdateCompany(t *testing.T) {
	payload := &types.UpdateCompanyPayload{
		Name:    utils.GenerateUniqueString("company"),
		Email:   utils.GenerateUniqueString("company_email"),
		Phone:   utils.GenerateUniqueString("48"),
		Address: utils.GenerateUniqueString("company_address"),
	}
	url := testServer.URL + "/company/1"

	// send PUT request to update company
	resp := sendPutRequest(url, payload, t)
	checkResponse(resp, http.StatusOK, t)

	// after updating, retrieve the company to validate the changes
	resp = sendGetRequest(url, t)
	body := checkResponse(resp, http.StatusOK, t)

	var company store.Company
	if err := json.Unmarshal(body, &company); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	// validate if the company was updated correctly
	if company.Name != payload.Name {
		t.Errorf("expected name %s, got %s", payload.Name, company.Name)
	}
	if company.Email != payload.Email {
		t.Errorf("expected email %s, got %s", payload.Email, company.Email)
	}
	if company.Phone != payload.Phone {
		t.Errorf("expected phone %s, got %s", payload.Phone, company.Phone)
	}
	if company.Address != payload.Address {
		t.Errorf("expected address %s, got %s", payload.Address, company.Address)
	}
}

func TestGetCompanies(t *testing.T) {
	// create test companies
	_, payloads := generateCompanies(10, t)

	// test GET /companies with filters and pagination
	url := testServer.URL + "/company/batch?name[sw]=company&email[ct]=company&phone[sw]=48&page=1&page_size=10&sort=name-asc"
	resp := sendGetRequest(url, t)
	body := checkResponse(resp, http.StatusOK, t)

	var companies []store.Company
	if err := json.Unmarshal(body, &companies); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	if len(companies) != len(payloads) {
		t.Errorf("expected 4 companies, got %d", len(companies))
	}

	for i := 1; i < len(companies); i++ {
		if companies[i].Name < companies[i-1].Name {
			t.Errorf("companies are not sorted by name in ascending order: %s vs %s", companies[i-1].Name, companies[i].Name)
		}
	}

	for _, company := range companies {
		if !strings.HasPrefix(company.Name, "company") {
			t.Errorf("company name %s does not start with 'company'", company.Name)
		}
		if !strings.Contains(company.Email, "company") {
			t.Errorf("company email %s does not contain 'company'", company.Email)
		}
		if !strings.HasPrefix(company.Phone, "48") {
			t.Errorf("company phone %s does not start with '48'", company.Phone)
		}
	}

	// test filtering by name
	url = testServer.URL + "/company/batch?name[sw]=company1"
	resp = sendGetRequest(url, t)
	body = checkResponse(resp, http.StatusOK, t)

	if err := json.Unmarshal(body, &companies); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	if len(companies) != 1 {
		t.Errorf("expected 1 company, got %d", len(companies))
	}
	for _, company := range companies {
		if !strings.Contains(company.Name, "company1") {
			t.Errorf("company name %s does not contain 'company1'", company.Name)
		}
	}
}
