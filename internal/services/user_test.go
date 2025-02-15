package services

import (
	"testing"

	"github.com/mwdev22/CarRental/internal/store/mock"
	"github.com/mwdev22/CarRental/internal/types"
)

func TestUserService(t *testing.T) {
	userService := NewUserService(mock.NewUserRepository())

	t.Run("RegisterUser", func(t *testing.T) {
		tests := []struct {
			name        string
			payload     *types.CreateUserPayload
			expectError bool
		}{
			{
				name: "successful registration",
				payload: &types.CreateUserPayload{
					Username: "testuser",
					Password: "password",
					Email:    "testuser@example.com",
					Role:     types.UserTypeUser,
				},
				expectError: false,
			},
			{
				name: "duplicate username",
				payload: &types.CreateUserPayload{
					Username: "testuser",
					Password: "password2",
					Email:    "testuser2@example.com",
					Role:     types.UserTypeUser,
				},
				expectError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := userService.Register(tt.payload)
				if tt.expectError && err == nil {
					t.Errorf("expected an error, got nil")
				}
				if !tt.expectError && err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			})
		}
	})

	t.Run("LoginUser", func(t *testing.T) {
		// Register a user first
		registerPayload := &types.CreateUserPayload{
			Username: "loginuser",
			Password: "password",
			Email:    "loginuser@example.com",
			Role:     types.UserTypeUser,
		}
		err := userService.Register(registerPayload)
		if err != nil {
			t.Fatalf("failed to register user: %v", err)
		}

		tests := []struct {
			name        string
			payload     *types.LoginPayload
			expectError bool
			expectToken bool
		}{
			{
				name: "successful login",
				payload: &types.LoginPayload{
					Username: "loginuser",
					Password: "password",
				},
				expectError: false,
				expectToken: true,
			},
			{
				name: "invalid password",
				payload: &types.LoginPayload{
					Username: "loginuser",
					Password: "wrongpassword",
				},
				expectError: true,
				expectToken: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				token, err := userService.Login(tt.payload)
				if tt.expectError && err == nil {
					t.Errorf("expected an error, got nil")
				}
				if !tt.expectError && err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
				if tt.expectToken && token == "" {
					t.Errorf("expected a token, got empty string")
				}
				if !tt.expectToken && token != "" {
					t.Errorf("expected no token, got: %s", token)
				}
			})
		}
	})

	t.Run("UpdateUser", func(t *testing.T) {
		// Register a user first
		registerPayload := &types.CreateUserPayload{
			Username: "updateuser",
			Password: "password",
			Email:    "updateuser@example.com",
		}
		err := userService.Register(registerPayload)
		if err != nil {
			t.Fatalf("failed to register user: %v", err)
		}

		updatePayload := &types.UpdateUserPayload{
			Username: "updateduser",
			Email:    "updateduser@example.com",
		}

		err = userService.Update(1, updatePayload)
		if err != nil {
			t.Fatalf("failed to update user: %v", err)
		}

		user, err := userService.GetByID(1)
		if err != nil {
			t.Fatalf("failed to get user: %v", err)
		}

		if user.Username != "updateduser" {
			t.Errorf("expected username: updateduser, got: %s", user.Username)
		}
		if user.Email != "updateduser@example.com" {
			t.Errorf("expected email: updateduser@example.com, got: %s", user.Email)
		}
	})

	t.Run("DeleteUser", func(t *testing.T) {
		// Register a user first
		registerPayload := &types.CreateUserPayload{
			Username: "deleteuser",
			Password: "password",
			Email:    "deleteuser@example.com",
		}
		err := userService.Register(registerPayload)
		if err != nil {
			t.Fatalf("failed to register user: %v", err)
		}

		err = userService.Delete(1)
		if err != nil {
			t.Fatalf("failed to delete user: %v", err)
		}

		_, err = userService.GetByID(1)
		if err == nil {
			t.Errorf("expected an error when fetching deleted user, got nil")
		}
	})
}
