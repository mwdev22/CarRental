package services

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mwdev22/FileStorage/internal/config"
	"github.com/mwdev22/FileStorage/internal/store"
	"github.com/mwdev22/FileStorage/internal/types"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userStore store.UserStore
}

func NewUserService(userStore store.UserStore) *UserService {
	return &UserService{
		userStore: userStore,
	}
}

func (s *UserService) Register(payload *types.CreateUserPayload) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return types.ServiceError(fmt.Errorf("failed to hash password: %v", err))
	}

	user := &store.User{
		Username: payload.Username,
		Password: hashedPassword,
		Email:    payload.Email,
		Created:  time.Now(),
	}

	if err := s.userStore.Create(context.Background(), user); err != nil {
		return types.DatabaseError(err)
	}

	return nil
}

func (s *UserService) Login(payload *types.LoginPayload) (string, error) {
	user, err := s.userStore.GetByUsername(context.Background(), payload.Username)
	if err != nil {
		return "", types.DatabaseError(err)
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(payload.Password)); err != nil {
		return "", types.Unauthorized("invalid password")
	}

	token, err := generateJWT(user)
	if err != nil {
		return "", types.ServiceError(err)
	}

	return token, nil
}

func (s *UserService) GetByID(id int) (*store.User, error) {
	user, err := s.userStore.GetByID(context.Background(), id)
	if err != nil {
		return nil, types.DatabaseError(err)
	}
	return user, nil
}

func generateJWT(user *store.User) (string, error) {
	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(config.SecretKey)
	if err != nil {
		return "", types.ServiceError(err)
	}

	return signedToken, nil
}

func (s *UserService) Delete(id int) error {
	if err := s.userStore.Delete(context.Background(), id); err != nil {
		return types.DatabaseError(err)
	}

	return nil
}

func (s *UserService) Update(payload *types.UpdateUserPayload, id int) error {
	user, err := s.userStore.GetByID(context.Background(), id)
	if err != nil {
		return types.DatabaseError(fmt.Errorf("failed to get user: %v", err))
	}

	if payload.Email != "" {
		user.Email = payload.Email
	}

	if payload.Username != "" {
		user.Username = payload.Username
	}

	if err := s.userStore.Update(context.Background(), user); err != nil {
		return types.DatabaseError(fmt.Errorf("failed to update user: %v", err))
	}

	return nil
}
