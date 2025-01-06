package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/mwdev22/CarRental/internal/services"
	"github.com/mwdev22/CarRental/internal/types"
	"github.com/mwdev22/CarRental/internal/utils"
)

type CompanyHandler struct {
	mux     *http.ServeMux
	company *services.CompanyService
}

func NewCompanyHandler(mux *http.ServeMux, company *services.CompanyService) *CompanyHandler {
	return &CompanyHandler{
		mux:     mux,
		company: company,
	}
}

func (h *CompanyHandler) RegisterRoutes() {
	logger, err := utils.MakeLogger("company")
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	h.mux.HandleFunc("POST /company", roleMiddleware(h.handleCreateCompany, types.UserTypeCompanyOwner, logger))
	h.mux.HandleFunc("GET /company/{id}", roleMiddleware(h.handleGetCompanyByID, types.UserTypeCompanyOwner, logger))
}

func (h *CompanyHandler) handleCreateCompany(w http.ResponseWriter, r *http.Request) error {
	var payload types.CreateCompanyPayload
	if err := types.ParseJSON(r, &payload); err != nil {
		return types.InvalidJSON(err)
	}

	userId, ok := r.Context().Value(userIdKey).(int)
	if !ok {
		return types.Unauthorized("user id not found in token")
	}

	if err := h.company.Create(&payload, userId); err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "company created successfully!",
	})
}

func (h *CompanyHandler) handleGetCompanyByID(w http.ResponseWriter, r *http.Request) error {
	companyId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return types.BadPathParameter("id")
	}
	company, err := h.company.GetByID(companyId)
	if err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, company)
}
