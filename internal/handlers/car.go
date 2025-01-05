package handlers

import (
	"log"
	"net/http"

	"github.com/mwdev22/CarRental/internal/services"
	"github.com/mwdev22/CarRental/internal/types"
	"github.com/mwdev22/CarRental/internal/utils"
)

type CarHandler struct {
	mux *http.ServeMux
	car *services.CarService
}

func NewCarHandler(mux *http.ServeMux, car *services.CarService) *CarHandler {
	return &CarHandler{
		mux: mux,
		car: car,
	}
}

func (h *CarHandler) RegisterRoutes() {
	logger, err := utils.MakeLogger("car")
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	h.mux.HandleFunc("POST /car", makeHandler(h.handleCreateCar, logger))

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
