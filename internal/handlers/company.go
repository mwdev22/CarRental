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
	logger  *log.Logger
}

func NewCompanyHandler(mux *http.ServeMux, company *services.CompanyService, logger *log.Logger) *CompanyHandler {
	h := &CompanyHandler{
		mux:     mux,
		company: company,
		logger:  logger,
	}

	h.mux.HandleFunc("POST /company", roleMiddleware(h.handleCreateCompany, types.UserTypeCompanyOwner, logger))
	h.mux.HandleFunc("GET /company/{id}", roleMiddleware(h.handleGetCompanyByID, types.UserTypeCompanyOwner, logger))
	h.mux.HandleFunc("PUT /company/{id}", roleMiddleware(h.handleUpdateCompany, types.UserTypeCompanyOwner, logger))
	h.mux.HandleFunc("DELETE /company/{id}", roleMiddleware(h.handleDeleteCompany, types.UserTypeCompanyOwner, logger))

	// check for allowed operators in utils/handlers.go
	// for example: get the first 10 companies with name ends with "company" and
	// email containing "company" and phone starts with "48" order by name ascending
	// you pass params like {field}[{operator}]={value}
	// sorting like sort={field}-{direction}
	// GET /companies?page=1&page_size=10&sort=name-asc&name[eq]=company&email[ct]=company&phone[sw]=48
	h.mux.HandleFunc("GET /company/batch", roleMiddleware(h.handleGetCopmanies, types.UserTypeCompanyOwner, logger))

	return h
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

func (h *CompanyHandler) handleUpdateCompany(w http.ResponseWriter, r *http.Request) error {
	companyId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return types.BadPathParameter("id")
	}

	var payload types.UpdateCompanyPayload
	if err := types.ParseJSON(r, &payload); err != nil {
		return types.InvalidJSON(err)
	}

	userId, ok := r.Context().Value(userIdKey).(int)
	if !ok {
		return types.Unauthorized("user id not found in token")
	}

	if err := h.company.Update(companyId, userId, &payload); err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "company updated successfully!",
	})
}

func (h *CompanyHandler) handleDeleteCompany(w http.ResponseWriter, r *http.Request) error {
	companyId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return types.BadPathParameter("id")
	}

	userId, ok := r.Context().Value(userIdKey).(int)
	if !ok {
		return types.Unauthorized("user id not found in token")
	}

	if err := h.company.Delete(companyId, userId); err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "company deleted successfully!",
	})
}

func (h *CompanyHandler) handleGetCopmanies(w http.ResponseWriter, r *http.Request) error {
	filters, err := utils.ParseQueryFilters(r)
	if err != nil {
		return err
	}

	opts, err := utils.ParseQueryOptions(r)
	if err != nil {
		return err
	}

	companies, err := h.company.GetAll(filters, opts)
	if err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, companies)
}
