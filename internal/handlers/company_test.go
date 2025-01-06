package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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
