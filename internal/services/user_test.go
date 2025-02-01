package services

import (
	"testing"

	"github.com/mwdev22/CarRental/internal/store/mock"
	"github.com/mwdev22/CarRental/internal/types"
)

func TestRegisterUser(t *testing.T) {
	userService := NewUserService(mock.NewUserRepo())

	payload := &types.CreateUserPayload{
		Username: "testuser",
		Password: "password",
		Email:    "afssdafsa@gmail.com",
		Role:     types.UserTypeUser,
	}

	err := userService.Register(payload)

	if err != nil {
		t.Errorf("failed to register user: %s", err)
	}

}

func TestLoginUser(t *testing.T) {
	userService := NewUserService(mock.NewUserRepo())

	payload := &types.CreateUserPayload{
		Username: "testuser2",
		Password: "password2",
		Email:    "afssdafsa@gmail.com",
		Role:     types.UserTypeUser,
	}

	err := userService.Register(payload)

	if err != nil {
		t.Errorf("failed to register user: %s", err)
	}

	loginPayload := &types.LoginPayload{
		Username: "testuser2",
		Password: "password2",
	}

	token, err := userService.Login(loginPayload)

	if err != nil {
		t.Errorf("failed to login user: %s", err)
	}

	if token == "" {
		t.Errorf("expected token, got empty string")
	}

}

func TestUpdateUser(t *testing.T) {
	userService := NewUserService(mock.NewUserRepo())

	payload := &types.CreateUserPayload{
		Username: "testuser3",
		Password: "password3",
		Email:    "fsadfsa@gmail.com",
	}

	err := userService.Register(payload)
	if err != nil {
		t.Errorf("failed to register user: %s", err)
	}

	updatePayload := &types.UpdateUserPayload{
		Username: "testuser3up",
		Email:    "newemail@fdsfa.com",
	}
	err = userService.Update(1, updatePayload)
	if err != nil {
		t.Errorf("failed to update user: %s", err)
	}

	user, err := userService.GetByID(1)
	if err != nil {
		t.Errorf("failed to get user: %s", err)
	}

	if user.Username != "testuser3up" {
		t.Errorf("expected username: testuser3up, got %s", user.Username)
	}

	if user.Email != "newemail@fdsfa.com" {
		t.Errorf("expected email: newemail@fdsfa.com, got %s", user.Email)
	}
}

func TestDeleteUser(t *testing.T) {
	userService := NewUserService(mock.NewUserRepo())

	payload := &types.CreateUserPayload{
		Username: "testuser4",
		Password: "password4",
		Email:    "usertoDel@dfsaf.com",
	}
	err := userService.Register(payload)
	if err != nil {
		t.Errorf("failed to register user: %s", err)
	}

	err = userService.Delete(1)
	if err != nil {
		t.Errorf("failed to delete user: %s", err)
	}

}
