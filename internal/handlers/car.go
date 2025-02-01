package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/mwdev22/CarRental/internal/services"
	"github.com/mwdev22/CarRental/internal/types"
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
