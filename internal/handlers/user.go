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
	"github.com/mwdev22/CarRental/internal/utils"
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

// @Summary Register a new user
// @Description Registers a user using the provided payload
// @Accept json
// @Produce json
// @Param payload body types.CreateUserPayload true "User registration details"
// @Success 200 {object} map[string]string
// @Router /register [post]
func (h *UserHandler) handleRegister(w http.ResponseWriter, r *http.Request) error {
	var payload types.CreateUserPayload
	if err := types.ParseJSON(r, &payload); err != nil {
		return types.InvalidJSON(err)
	}

	if errors := utils.ValidateStruct(&payload); len(errors) > 0 {
		return types.ValidationError(errors)
	}

	err := h.user.Register(&payload)
	if err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "user registered successfully!",
	})
}

// @Summary User login
// @Description Authenticates a user and returns a token
// @Accept json
// @Produce json
// @Param payload body types.LoginPayload true "User login details"
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} map[string]string
// @Router /login [post]
func (h *UserHandler) handleLogin(w http.ResponseWriter, r *http.Request) error {
	var payload types.LoginPayload
	if err := types.ParseJSON(r, &payload); err != nil {
		return types.InvalidJSON(err)
	}

	if errors := utils.ValidateStruct(&payload); len(errors) > 0 {
		return types.ValidationError(errors)
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

// @Summary Get user by ID
// @Description Retrieves user details by ID
// @Produce json
// @Param id path int true "User ID"
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} types.User
// @Router /user/{id} [get]
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

// @Summary Delete a user
// @Description Deletes a user by ID
// @Param id path int true "User ID"
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} map[string]string
// @Router /user/{id} [delete]
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

// @Summary Update a user
// @Description Updates a user's profile
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param Authorization header string true "Bearer Token"
// @Param payload body types.UpdateUserPayload true "Updated user details"
// @Success 200 {object} map[string]string
// @Router /user/{id} [put]
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
	if errors := utils.ValidateStruct(&payload); len(errors) > 0 {
		return types.ValidationError(errors)
	}

	err = h.user.Update(userIDInt, &payload)
	if err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("user with id %v updated successfully", userIDInt),
	})
}

// @Summary Check Token
// @Description Retrieve user's token claims
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} map[string]string
// @Router /check-token [post]
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
