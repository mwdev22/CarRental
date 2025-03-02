package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/mwdev22/CarRental/internal/services"
	"github.com/mwdev22/CarRental/internal/types"
	"github.com/mwdev22/CarRental/internal/utils"
)

type CarHandler struct {
	mux    *http.ServeMux
	car    *services.CarService
	logger *log.Logger
}

func NewCarHandler(mux *http.ServeMux, car *services.CarService, logger *log.Logger) *CarHandler {
	h := &CarHandler{
		mux:    mux,
		car:    car,
		logger: logger,
	}

	h.mux.HandleFunc("POST /car", makeHandler(h.handleCreateCar, logger))
	h.mux.HandleFunc("GET /car/{id}", makeHandler(h.handleGetCarByID, logger))
	h.mux.HandleFunc("DELETE /car/{id}", makeHandler(h.handleDeleteCarByID, logger))
	h.mux.HandleFunc("PUT /car/{id}", makeHandler(h.handleUpdateCarByID, logger))

	// check for allowed operators in utils/handlers.go
	// you pass params like {field}[{operator}]={value}
	// sorting like sort={field}-{direction}
	// GET /car?page=1&page_size=10&sort=name-asc&make[ct]=Mercedes&model[ct]=CLA&year=2022
	h.mux.HandleFunc("GET /car/batch", makeHandler(h.handleGetCars, logger))

	return h
}

func (h *CarHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// @Summary Create a car
// @Description Creates a new car record
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param payload body types.CreateCarPayload true "Car data"
// @Tags Car
// @Success 200 {object} map[string]string
// @Router /cars [post]
func (h *CarHandler) handleCreateCar(w http.ResponseWriter, r *http.Request) error {
	var payload types.CreateCarPayload
	if err := types.ParseJSON(r, &payload); err != nil {
		return types.InvalidJSON(err)
	}

	if errors := utils.ValidateStruct(&payload); len(errors) > 0 {
		return types.ValidationError(errors)
	}

	err := h.car.CreateCar(&payload)
	if err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "car created successfully!",
	})
}

// @Summary Get car by ID
// @Description Retrieves a car by its ID
// @Produce json
// @Param id path int true "Car ID"
// @Tags Car
// @Success 200 {object} types.Car
// @Router /car/{id} [get]
func (h *CarHandler) handleGetCarByID(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return types.BadPathParameter("id")
	}
	car, err := h.car.GetByID(idInt)
	if err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, car)
}

// @Summary Delete car by ID
// @Description Deletes a car by its ID
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "Car ID"
// @Tags Car
// @Success 200 {object} map[string]string
// @Router /car/{id} [delete]
func (h *CarHandler) handleDeleteCarByID(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return types.BadPathParameter("id")
	}

	err = h.car.Delete(idInt)
	if err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "car deleted successfully!",
	})
}

// @Summary Update car by ID
// @Description Updates a car's details
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "Car ID"
// @Param payload body types.UpdateCarPayload true "Updated car data"
// @Tags Car
// @Success 200 {object} map[string]string
// @Router /car/{id} [put]
func (h *CarHandler) handleUpdateCarByID(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return types.BadPathParameter("id")
	}

	var payload types.UpdateCarPayload
	if err := types.ParseJSON(r, &payload); err != nil {
		return types.InvalidJSON(err)
	}

	if errors := utils.ValidateStruct(&payload); len(errors) > 0 {
		return types.ValidationError(errors)
	}

	err = h.car.UpdateCar(idInt, &payload)
	if err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "car updated successfully!",
	})
}

// @Summary Get cars
// @Description Retrieves a list of cars with optional filters
// @Produce json
// @Param filters query string false "Filters for car retrieval. eg. make[ct]=Mercedes&model[ct]=CLA&year=2022"
// @Param sort query string false "sort for car retrieval, eg. id-asc"
// @Param page query int false "page number for car retrieval"
// @Param page_size query int false "number of items per page"
// @Tags Car
// @Success 200 {array} types.Car
// @Router /cars [get]
func (h *CarHandler) handleGetCars(w http.ResponseWriter, r *http.Request) error {
	filters, err := utils.ParseQueryFilters(r)
	if err != nil {
		return err
	}

	opts, err := utils.ParseQueryOptions(r)
	if err != nil {
		return err
	}

	cars, err := h.car.GetBatch(filters, opts)
	if err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, cars)
}
