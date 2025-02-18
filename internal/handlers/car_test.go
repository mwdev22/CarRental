package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/mwdev22/CarRental/internal/types"
	"github.com/mwdev22/CarRental/internal/utils"
)

func generateCars(count int, t *testing.T) ([]byte, []*types.CreateCarPayload) {
	payloads := make([]*types.CreateCarPayload, 0)
	var body []byte
	for i := 0; i < count; i++ {
		payload := &types.CreateCarPayload{
			Make:           "Fiat",
			Model:          "Punto",
			Year:           2021,
			Color:          "Blue",
			RegistrationNo: utils.GenerateUniqueString(""),
			PricePerDay:    50.0,
			CompanyID:      1,
		}
		url := testServer.URL + "/car"

		resp := sendPostRequest(url, payload, t)
		body = checkResponse(resp, http.StatusOK, t)
		payloads = append(payloads, payload)
	}
	return body, payloads
}

func TestCreateCar(t *testing.T) {

	// TestLogin(t)

	body, _ := generateCars(1, t)
	var responseBody map[string]string
	if err := json.Unmarshal(body, &responseBody); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

}

func TestGetCarByID(t *testing.T) {
	url := testServer.URL + "/car/1"

	resp := sendGetRequest(url, t)

	checkResponse(resp, 200, t)
}

func TestUpdateCar(t *testing.T) {
	url := testServer.URL + "/car/1"

	payload := &types.UpdateCarPayload{
		Make:           "Ford",
		Model:          "Focus",
		Year:           2021,
		Color:          "Blue",
		RegistrationNo: "ABC123567",
		PricePerDay:    100,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(payloadBytes))
	if err != nil {
		t.Fatalf("failed to create PUT request: %v", err)
	}

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
}

func TestGetBatchCars(t *testing.T) {

	_, payloads := generateCars(10, t)

	url := testServer.URL + "/car/batch?limit=10&offset=0"

	resp := sendGetRequest(url, t)

	body := checkResponse(resp, http.StatusOK, t)

	var companies []types.Car
	if err := json.Unmarshal(body, &companies); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	if len(companies) != len(payloads) {
		t.Errorf("expected %v companies, got %v", len(payloads), len(companies))
	}
}

func TestDeleteCar(t *testing.T) {
	url := testServer.URL + "/car/1"

	resp := sendDeleteRequest(url, t)
	checkResponse(resp, http.StatusOK, t)
}
