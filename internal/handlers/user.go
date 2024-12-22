package handlers

import (
	"log"
	"net/http"

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
	h.mux.HandleFunc("GET /login", makeHandler(h.handleLogin, logger))
	h.mux.HandleFunc("POST /check-token", authMiddleware(h.handleCheckToken, logger))

}

func (h *UserHandler) handleRegister(w http.ResponseWriter, r *http.Request) error {

	return types.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "user registered successfully!",
	})
}

func (h *UserHandler) handleLogin(w http.ResponseWriter, r *http.Request) error {
	return types.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "user logged in successfully!",
	})
}

func (h *UserHandler) handleCheckToken(w http.ResponseWriter, r *http.Request) error {
	return types.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "token is valid!",
	})
}
