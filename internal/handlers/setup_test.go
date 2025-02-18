package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mwdev22/CarRental/internal/config"
	"github.com/mwdev22/CarRental/internal/services"
	"github.com/mwdev22/CarRental/internal/store"
	"github.com/mwdev22/CarRental/internal/store/mock"
	"github.com/mwdev22/CarRental/internal/utils"
	"github.com/redis/go-redis/v9"
)

var (
	testServer   *httptest.Server
	authHeader   string
	testUsername = utils.GenerateUniqueString("testuser2")
	testPassword = "testpassword"
)

func TestMain(m *testing.M) {
	_ = config.New()
	// changes the working directory to the project root, important for running migrations
	_, err := initializeTests()
	if err != nil {
		fmt.Printf("failed to initialize tests: %v", err)
		os.Exit(1)
	}

	code := m.Run()

	os.Remove("./test.db")

	os.Exit(code)
}

func initializeTests() (*httptest.Server, error) {

	// stores and services
	userStore := mock.NewUserRepository()
	userService := services.NewUserService(userStore)

	companyStore := mock.NewCompanyRepository()
	companyService := services.NewCompanyService(companyStore)

	carStore := mock.NewCarRepository()
	carService := services.NewCarService(carStore)

	bookingStore := mock.NewBookingStore()
	bookingService := services.NewBookingService(bookingStore, carStore, userStore)

	r := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   5,
	})
	c := store.NewRedisCache(r)

	// handlers
	mux := http.NewServeMux()
	_ = NewUserHandler(mux, userService, log.Default())
	_ = NewCompanyHandler(mux, companyService, log.Default())
	_ = NewCarHandler(mux, carService, log.Default())
	_ = NewBookingHandler(mux, bookingService, carService, c, log.Default())
	// setup the test server
	testServer = httptest.NewServer(mux)
	return testServer, nil
}

func checkToken() (map[string]interface{}, error) {
	url := testServer.URL + "/check-token"

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request to /check-token: %v", err)
	}

	req.Header.Set("Authorization", authHeader)

	resp, err := testServer.Client().Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request to /check-token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}
		return nil, fmt.Errorf("token validation failed. Expected status 200, got %d, body: %s", resp.StatusCode, body)
	}

	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %v", err)
	}

	claims, ok := responseBody["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'data' field in response, got: %v", responseBody)
	}

	return claims, nil
}

// helper function to send a GET request
func sendGetRequest(url string, t *testing.T) *http.Response {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("failed to create GET request: %v", err)
	}
	req.Header.Set("Authorization", authHeader)

	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("failed to send GET request: %v", err)
	}
	return resp
}

// helper function to send a POST request with JSON payload
func sendPostRequest(url string, payload interface{}, t *testing.T) *http.Response {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		t.Fatalf("failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("failed to send POST request: %v", err)
	}
	return resp
}

// helper function to send a PUT request with JSON payload
func sendPutRequest(url string, payload interface{}, t *testing.T) *http.Response {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewReader(payloadBytes))
	if err != nil {
		t.Fatalf("failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("failed to send POST request: %v", err)
	}
	return resp
}

// helper function to send a GET request
func sendDeleteRequest(url string, t *testing.T) *http.Response {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatalf("failed to create GET request: %v", err)
	}
	req.Header.Set("Authorization", authHeader)

	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("failed to send GET request: %v", err)
	}
	return resp
}

// helper function to check the response status and decode the body
func checkResponse(resp *http.Response, expectedStatusCode int, t *testing.T) []byte {
	if resp.StatusCode != expectedStatusCode {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("failed to read response body: %v", err)
		}
		t.Errorf("expected status %d, got %d, body: %s", expectedStatusCode, resp.StatusCode, body)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	return body
}
