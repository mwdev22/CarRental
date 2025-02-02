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

func (h *CarHandler) handleCreateCar(w http.ResponseWriter, r *http.Request) error {
	var payload types.CreateCarPayload
	if err := types.ParseJSON(r, &payload); err != nil {
		return types.InvalidJSON(err)
	}

	err := h.car.CreateCar(&payload)
	if err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "car created successfully!",
	})
}

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

	err = h.car.UpdateCar(idInt, &payload)
	if err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "car updated successfully!",
	})
}

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
