package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/mwdev22/CarRental/internal/store"
	"github.com/mwdev22/CarRental/internal/types"
	"github.com/mwdev22/CarRental/internal/utils"
)

// test Register user route
func TestRegister(t *testing.T) {
	url := testServer.URL + "/register"

	testUsername = utils.GenerateUniqueString("testuser")
	testPassword = "testpassword"

	payload := &types.CreateUserPayload{
		Username: testUsername,
		Password: testPassword,
		Email:    utils.GenerateUniqueString("email@blabla.com"),
		Role:     types.UserTypeAdmin,
	}

	resp := sendPostRequest(url, payload, t)

	checkResponse(resp, http.StatusOK, t)
}

// test Login user route
func TestLogin(t *testing.T) {
	registerURL := testServer.URL + "/register"

	registerPayload := &types.CreateUserPayload{
		Username: testUsername,
		Password: testPassword,
		Email:    utils.GenerateUniqueString("email@blabla.com"),
		Role:     types.UserTypeAdmin,
	}

	registerResp := sendPostRequest(registerURL, registerPayload, t)
	defer registerResp.Body.Close()

	if registerResp.StatusCode != http.StatusOK {
		t.Fatalf("failed to register test user: expected status 200, got %d", registerResp.StatusCode)
	}

	loginURL := testServer.URL + "/login"

	loginPayload := &types.LoginPayload{
		Username: testUsername,
		Password: testPassword,
	}

	loginResp := sendPostRequest(loginURL, loginPayload, t)
	defer loginResp.Body.Close()

	if loginResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(loginResp.Body)
		t.Fatalf("failed to log in: expected status 200, got %d, body: %s", loginResp.StatusCode, body)
	}

	var responseBody map[string]string
	if err := json.NewDecoder(loginResp.Body).Decode(&responseBody); err != nil {
		t.Fatalf("failed to parse login response body: %v", err)
	}

	token, ok := responseBody["token"]
	if !ok {
		t.Fatalf("token not found in login response")
	}

	authHeader = "Bearer " + token
	t.Logf("login successful, token obtained: %s", token)
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

	resp := sendPutRequest(url, payload, t)
	defer resp.Body.Close()
	body := checkResponse(resp, http.StatusOK, t)

	t.Logf("update User response: %s", string(body))
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

	resp := sendGetRequest(url, t)
	defer resp.Body.Close()

	body := checkResponse(resp, http.StatusOK, t)

	t.Logf("get User response: %s", resp.Body)

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
	resp := sendDeleteRequest(url, t)

	defer resp.Body.Close()
	respBody := checkResponse(resp, http.StatusOK, t)

	t.Logf("delete User response: %s", respBody)
}
