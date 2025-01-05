package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/mwdev22/CarRental/internal/types"
)

func TestCreateCar(t *testing.T) {
	url := testServer.URL + "/car"

	payload := &types.CreateCarPayload{
		Make:           "Toyota",
		Model:          "Corolla",
		Year:           2021,
		Color:          "Red",
		RegistrationNo: "ABC123",
		PricePerDay:    100,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	resp, err := testServer.Client().Post(url, "application/json", bytes.NewReader(payloadBytes))
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
