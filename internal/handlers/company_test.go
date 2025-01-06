package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/mwdev22/CarRental/internal/store"
	"github.com/mwdev22/CarRental/internal/types"
	"github.com/mwdev22/CarRental/internal/utils"
)

func TestCreateCompany(t *testing.T) {
	TestRegister(t)
	TestLogin(t)

	url := testServer.URL + "/company"

	payload := &types.CreateCompanyPayload{
		Name:    utils.GenerateUniqueString("company"),
		Email:   utils.GenerateUniqueString("company_email"),
		Phone:   utils.GenerateUniqueString("48"),
		Address: utils.GenerateUniqueString("company_address"),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("failed to read response body: %v", err)
		}
		t.Errorf("expected status 200, got %d, body: %s", resp.StatusCode, body)
	}

	var responseBody map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}
}

func TestGetCompanyByID(t *testing.T) {

	url := testServer.URL + "/company/1"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", authHeader)
	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("failed to read response body: %v", err)
		}
		t.Errorf("expected status 200, got %d, body: %s", resp.StatusCode, body)
	}

	var company store.Company
	if err := json.NewDecoder(resp.Body).Decode(&company); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}
}

func TestUpdateCompany(t *testing.T) {
	url := testServer.URL + "/company/1"

	payload := &types.UpdateCompanyPayload{
		Name:    utils.GenerateUniqueString("company"),
		Email:   utils.GenerateUniqueString("company_email"),
		Phone:   utils.GenerateUniqueString("48"),
		Address: utils.GenerateUniqueString("company_address"),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewReader(payloadBytes))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("failed to send PUT request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("failed to read response body: %v", err)
		}
		t.Errorf("expected status 200, got %d, body: %s", resp.StatusCode, body)
	}

	// after updating company, validate the changes

	url = testServer.URL + "/company/1"

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", authHeader)
	resp, err = testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("failed to read response body: %v", err)
		}
		t.Errorf("expected status 200, got %d, body: %s", resp.StatusCode, body)
	}

	var company store.Company
	if err := json.NewDecoder(resp.Body).Decode(&company); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

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
	// First, create some test company data using POST /company
	companyNames := []string{
		utils.GenerateUniqueString("company1"),
		utils.GenerateUniqueString("company2"),
		utils.GenerateUniqueString("company3"),
	}

	// Create companies by calling POST /company
	for _, name := range companyNames {
		payload := &types.CreateCompanyPayload{
			Name:    name,
			Email:   utils.GenerateUniqueString("company_email"),
			Phone:   utils.GenerateUniqueString("48"),
			Address: utils.GenerateUniqueString("company_address"),
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}

		url := testServer.URL + "/company"
		req, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader) // Assuming authHeader is set with a valid token

		resp, err := testServer.Client().Do(req)
		if err != nil {
			t.Fatalf("failed to send POST request: %v", err)
		}
		defer resp.Body.Close()

		// Check if POST request was successful
		if resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("failed to read response body: %v", err)
			}
			t.Errorf("expected status 200, got %d, body: %s", resp.StatusCode, body)
		}
	}

	// After inserting data, now test the GET /companies request
	url := testServer.URL + "/company/batch?name[sw]=company&email[ct]=company&phone[sw]=48&page=1&page_size=10&sort=name-asc"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", authHeader) // Assuming authHeader is set with a valid token

	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("failed to read response body: %v", err)
		}
		t.Errorf("expected status 200, got %d, body: %s", resp.StatusCode, body)
	}

	var companies []store.Company
	if err := json.NewDecoder(resp.Body).Decode(&companies); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	if len(companies) != 4 {
		t.Errorf("expected 4 companies, got %d", len(companies))
	}

	for i := 1; i < len(companies); i++ {
		if companies[i].Name < companies[i-1].Name {
			t.Errorf("companies are not sorted by name in ascending order: %s vs %s", companies[i-1].Name, companies[i].Name)
		}
	}

	// Check if companies match the filter (name ending with "company" and containing "company", phone starts with "48")
	for _, company := range companies {
		if !strings.HasPrefix(company.Name, "company") {
			t.Errorf("company name %s does not end with 'company'", company.Name)
		}
		if !strings.Contains(company.Email, "company") {
			t.Errorf("company email %s does not contain 'company'", company.Email)
		}
		if !strings.HasPrefix(company.Phone, "48") {
			t.Errorf("company phone %s does not start with '48'", company.Phone)
		}
	}

	url = testServer.URL + "/company/batch?name[sw]=company1"

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", authHeader) // Assuming authHeader is set with a valid token

	resp, err = testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("failed to read response body: %v", err)
		}
		t.Errorf("expected status 200, got %d, body: %s", resp.StatusCode, body)
	}

	if err := json.NewDecoder(resp.Body).Decode(&companies); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	if len(companies) != 1 {
		t.Errorf("expected 2 companies, got %d", len(companies))

		for _, company := range companies {
			if !strings.Contains(company.Name, "1") {
				t.Errorf("company name %s does not contain '1'", company.Name)
			}
		}
	}

}
