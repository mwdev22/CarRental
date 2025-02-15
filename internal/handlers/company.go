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

// @Summary Create a new company
// @Description Creates a company using the provided payload
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param payload body types.CreateCompanyPayload true "Company details"
// @Success 200 {object} map[string]string
// @Router /companies [post]
func (h *CompanyHandler) handleCreateCompany(w http.ResponseWriter, r *http.Request) error {
	var payload types.CreateCompanyPayload
	if err := types.ParseJSON(r, &payload); err != nil {
		return types.InvalidJSON(err)
	}

	if errors := utils.ValidateStruct(&payload); len(errors) > 0 {
		return types.ValidationError(errors)
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

// @Summary Get a company by ID
// @Description Retrieves a company based on the provided ID
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "Company ID"
// @Success 200 {object} types.Company
// @Router /companies/{id} [get]
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

// @Summary Update a company
// @Description Updates an existing company with new data
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "Company ID"
// @Param payload body types.UpdateCompanyPayload true "Updated company details"
// @Success 200 {object} map[string]string
// @Router /companies/{id} [put]
func (h *CompanyHandler) handleUpdateCompany(w http.ResponseWriter, r *http.Request) error {
	companyId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return types.BadPathParameter("id")
	}

	var payload types.UpdateCompanyPayload
	if err := types.ParseJSON(r, &payload); err != nil {
		return types.InvalidJSON(err)
	}

	if errors := utils.ValidateStruct(&payload); len(errors) > 0 {
		return types.ValidationError(errors)
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

// @Summary Delete a company
// @Description Deletes a company by ID
// @Param id path int true "Company ID"
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} map[string]string
// @Router /companies/{id} [delete]
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

// @Summary Get a list of companies
// @Description Retrieves a list of companies based on query filters and options
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param filters query object false "Query filters"
// @Param options query object false "Query options"
// @Success 200 {array} types.Company
// @Router /companies [get]
func (h *CompanyHandler) handleGetCopmanies(w http.ResponseWriter, r *http.Request) error {
	filters, err := utils.ParseQueryFilters(r)
	if err != nil {
		return err
	}

	opts, err := utils.ParseQueryOptions(r)
	if err != nil {
		return err
	}

	companies, err := h.company.GetBatch(filters, opts)
	if err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, companies)
}
