package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mwdev22/FileStorage/internal/config"
	"github.com/mwdev22/FileStorage/internal/store"
	"github.com/mwdev22/FileStorage/internal/types"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo store.Storage
}

func NewUserService(userRepo store.Storage) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Register(payload *types.CreateUserRequest) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	user := &store.User{
		Username: payload.Username,
		Password: hashedPassword,
		Email:    payload.Email,
		Created:  time.Now(),
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

func (s *UserService) Login(payload *types.LoginRequest) (string, error) {
	user, err := s.userRepo.GetByUsername(payload.Username)
	if err != nil {
		return "", fmt.Errorf("failed to get user by username: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(payload.Password)); err != nil {
		return "", fmt.Errorf("invalid password: %v", err)
	}

	token, err := generateJWT(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %v", err)
	}

	return token, nil
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
		return "", err
	}

	return signedToken, nil
}
