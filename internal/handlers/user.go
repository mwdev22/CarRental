package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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

	h.mux.HandleFunc("GET /user/{id}", makeHandler(h.handleGetUser, logger))
	h.mux.HandleFunc("DELETE /user/{id}", makeHandler(h.handleDeleteUser, logger))
	h.mux.HandleFunc("PUT /user/{id}", makeHandler(h.handleUpdateUser, logger))

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

func (h *UserHandler) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	userID := r.PathValue("id")
	if userID == "" {
		return badPathParameter("id")
	}
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return newApiError(http.StatusBadRequest, fmt.Errorf("invalid user id: %v", err))
	}
	user, err := h.userService.GetByID(userIDInt)
	if err != nil {
		return serviceError(err)
	}

	return types.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	userID := r.PathValue("id")
	if userID == "" {
		return badPathParameter("id")
	}
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return newApiError(http.StatusBadRequest, fmt.Errorf("invalid user id: %v", err))
	}
	err = h.userService.Delete(userIDInt)
	if err != nil {
		return serviceError(err)
	}

	return types.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "user deleted successfully!"})
}

func (h *UserHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {

	userID := r.PathValue("id")
	if userID == "" {
		return badPathParameter("id")
	}
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return newApiError(http.StatusBadRequest, fmt.Errorf("invalid user id: %v", err))
	}

	var payload types.UpdateUserPayload
	if err := types.ParseJSON(r, &payload); err != nil {
		return invalidJSON(err)
	}

	err = h.userService.Update(&payload, userIDInt)
	if err != nil {
		return serviceError(err)
	}

	return types.WriteJSON(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("user with id %v updated successfully", userIDInt),
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
