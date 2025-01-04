package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mwdev22/FileStorage/internal/config"
	"github.com/mwdev22/FileStorage/internal/database"
	"github.com/mwdev22/FileStorage/internal/services"
	"github.com/mwdev22/FileStorage/internal/store"
	"github.com/mwdev22/FileStorage/internal/types"
)

var testServer *httptest.Server
var authHeader string

func TestMain(m *testing.M) {
	// changes the working directory to the project root, important for running migrations
	_ = config.New()
	// setup the test db
	testDB, err := database.OpenTestSqlDB()
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}

	// setup routes
	userStore := store.NewUserRepo(testDB)
	userService := services.NewUserService(userStore)
	userHandler := NewUserHandler(http.NewServeMux(), userService)
	userHandler.RegisterRoutes()

	testServer = httptest.NewServer(userHandler.mux)
	defer testServer.Close()

	code := m.Run()

	// delete the database after tests
	os.Remove("./test.db")

	os.Exit(code)
}

// test Register user route
func TestRegister(t *testing.T) {
	url := testServer.URL + "/register"

	payload := &types.CreateUserPayload{
		Username: "testuser",
		Password: "testpassword",
		Email:    "email@blabla.com",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	// Use testServer.Client() to send the POST request
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
}

// test Login user route
func TestLogin(t *testing.T) {
	url := testServer.URL + "/login"

	payload := &types.LoginPayload{
		Username: "testuser",
		Password: "testpassword",
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

	token, ok := responseBody["token"]
	if !ok {
		t.Fatalf("token not found in response")
	}

	authHeader = "Bearer " + token

}

// test update user route
func TestUpdateUser(t *testing.T) {

	// first check token to get user ID
	claims, err := checkToken()
	if err != nil {
		t.Fatalf("failed to check token: %v", err)
	}

	userID, ok := claims["id"].(float64)
	if !ok {
		t.Fatalf("expected user ID in claims, got: %v", claims)
	}
	userIDstr := fmt.Sprintf("%.0f", userID)

	// then send PUT request to update user
	url := testServer.URL + "/user/" + userIDstr

	payload := &types.UpdateUserPayload{
		Username: "testuser2",
		Email:    "newmail@gmail.com",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(payloadBytes))
	req.Header.Set("Authorization", authHeader)
	if err != nil {
		t.Fatalf("failed to create PUT request: %v", err)
	}

	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("failed to send PUT request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", resp.StatusCode, body)
	}

	t.Logf("update User response: %s", body)
}

func TestGetUser(t *testing.T) {
	claims, err := checkToken()
	if err != nil {
		t.Fatalf("failed to check token: %v", err)
	}

	userID, ok := claims["id"].(float64)
	if !ok {
		t.Fatalf("expected user ID in claims, got: %v", claims)
	}
	userIDstr := fmt.Sprintf("%.0f", userID)

	url := testServer.URL + "/user/" + userIDstr

	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", authHeader)
	if err != nil {
		t.Fatalf("failed to create GET request: %v", err)
	}

	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", resp.StatusCode, body)
	}

	t.Logf("get User response: %s", body)

	var user *store.User
	if err := json.Unmarshal(body, &user); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if user.ID != int(userID) {
		t.Errorf("expected user ID %d, got %d", int(userID), user.ID)
	}

	if user.Username != "testuser2" {
		t.Errorf("expected username testuser2, got %s", user.Username)
	}

	if user.Email != "newmail@gmail.com" {
		t.Errorf("expected email newmail@gmail.com, got %s", user.Email)
	}
}

func TestDeleteUser(t *testing.T) {
	claims, err := checkToken()
	if err != nil {
		t.Fatalf("failed to check token: %v", err)
	}

	userID, ok := claims["id"].(float64)
	if !ok {
		t.Fatalf("expected user ID in claims, got: %v", claims)
	}
	userIDstr := fmt.Sprintf("%.0f", userID)

	url := testServer.URL + "/user/" + userIDstr

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Authorization", authHeader)
	if err != nil {
		t.Fatalf("failed to create DELETE request: %v", err)
	}

	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("failed to send DELETE request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", resp.StatusCode, body)
	}

	t.Logf("delete User response: %s", body)
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
