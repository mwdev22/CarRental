package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/mwdev22/CarRental/internal/config"
	"github.com/mwdev22/CarRental/internal/services"
	"github.com/mwdev22/CarRental/internal/types"
)

type UserHandler struct {
	mux    *http.ServeMux
	user   *services.UserService
	logger *log.Logger
}

func NewUserHandler(mux *http.ServeMux, user *services.UserService, logger *log.Logger) *UserHandler {
	h := &UserHandler{
		mux:    mux,
		user:   user,
		logger: logger,
	}

	h.mux.HandleFunc("POST /register", makeHandler(h.handleRegister, logger))
	h.mux.HandleFunc("POST /login", makeHandler(h.handleLogin, logger))
	h.mux.HandleFunc("POST /check-token", makeHandler(h.handleCheckToken, logger))

	h.mux.HandleFunc("GET /user/{id}", authMiddleware(h.handleGetUser, logger))
	h.mux.HandleFunc("DELETE /user/{id}", authMiddleware(h.handleDeleteUser, logger))
	h.mux.HandleFunc("PUT /user/{id}", authMiddleware(h.handleUpdateUser, logger))

	return h
}

func (h *UserHandler) RegisterRoutes() {

}

func (h *UserHandler) handleRegister(w http.ResponseWriter, r *http.Request) error {
	var payload types.CreateUserPayload
	if err := types.ParseJSON(r, &payload); err != nil {
		return types.InvalidJSON(err)
	}

	err := h.user.Register(&payload)
	if err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "user registered successfully!",
	})
}

func (h *UserHandler) handleLogin(w http.ResponseWriter, r *http.Request) error {
	var payload types.LoginPayload
	if err := types.ParseJSON(r, &payload); err != nil {
		return types.InvalidJSON(err)
	}
	token, err := h.user.Login(&payload)
	if err != nil {
		return err
	}
	return types.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "user logged in successfully!",
		"token":   token,
	})
}

func (h *UserHandler) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	userID := r.PathValue("id")
	if userID == "" {
		return types.BadPathParameter("id")
	}
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return types.BadPathParameter("id")
	}
	user, err := h.user.GetByID(userIDInt)
	if err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	idFromToken := r.Context().Value(userIdKey)
	userID := r.PathValue("id")
	if userID == "" {
		return types.BadPathParameter("id")
	}
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return types.BadPathParameter("id")
	}

	if idFromToken != nil && idFromToken != userIDInt {
		return types.Unauthorized("user can only delete their own profile")
	}

	err = h.user.Delete(userIDInt)
	if err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "user deleted successfully!"})
}

func (h *UserHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {

	idFromToken := r.Context().Value(userIdKey)

	userID := r.PathValue("id")
	if userID == "" {
		return types.BadPathParameter("id")
	}
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return types.BadPathParameter("id")
	}

	if idFromToken != nil && idFromToken != userIDInt {
		return types.Unauthorized("user can only update their own profile")
	}

	var payload types.UpdateUserPayload
	if err := types.ParseJSON(r, &payload); err != nil {
		return types.InvalidJSON(err)
	}

	err = h.user.Update(userIDInt, &payload)
	if err != nil {
		return err
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
