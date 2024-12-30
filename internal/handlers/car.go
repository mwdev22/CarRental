package handlers

import (
	"net/http"

	"github.com/mwdev22/FileStorage/internal/services"
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
