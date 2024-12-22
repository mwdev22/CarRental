package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/mwdev22/FileStorage/internal/config"
	"github.com/mwdev22/FileStorage/internal/services"
	"github.com/mwdev22/FileStorage/internal/types"
	"github.com/mwdev22/FileStorage/internal/utils"
)

type UserHandler struct {
	mux         *http.ServeMux
	userService *services.UserService
}

func NewUserHandler(mux *http.ServeMux, userService *services.UserService) *UserHandler {
	return &UserHandler{
		mux:         mux,
		userService: userService,
	}
}

func (h *UserHandler) RegisterRoutes() {
	logger, err := utils.MakeLogger("user")
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	h.mux.HandleFunc("POST /register", makeHandler(h.handleRegister, logger))
	h.mux.HandleFunc("POST /login", makeHandler(h.handleLogin, logger))
	h.mux.HandleFunc("POST /check-token", makeHandler(h.handleCheckToken, logger))

}

func (h *UserHandler) handleRegister(w http.ResponseWriter, r *http.Request) error {
	var payload types.CreateUserRequest
	if err := types.ParseJSON(r, &payload); err != nil {
		return invalidJSON(err)
	}

	err := h.userService.Register(&payload)
	if err != nil {
		return serviceError(err)
	}

	return types.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "user registered successfully!",
	})
}

func (h *UserHandler) handleLogin(w http.ResponseWriter, r *http.Request) error {
	var payload types.LoginRequest
	if err := types.ParseJSON(r, &payload); err != nil {
		return invalidJSON(err)
	}
	token, err := h.userService.Login(&payload)
	if err != nil {
		return serviceError(err)
	}
	return types.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "user logged in successfully!",
		"token":   token,
	})
}

func (h *UserHandler) handleCheckToken(w http.ResponseWriter, r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return types.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "missing Authorization header",
		})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return types.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "invalid Authorization header format",
		})
	}
	tokenStr := parts[1]

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return config.SecretKey, nil
	})
	if err != nil || !token.Valid {
		return types.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "invalid token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return types.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "invalid token claims",
		})
	}

	return types.WriteJSON(w, http.StatusOK, map[string]any{
		"data": claims,
	})
}
