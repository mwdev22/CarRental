package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/mwdev22/CarRental/internal/services"
	"github.com/mwdev22/CarRental/internal/store"
	"github.com/mwdev22/CarRental/internal/types"
	"github.com/mwdev22/CarRental/internal/utils"
)

type BookingHandler struct {
	mux         *http.ServeMux
	booking     *services.BookingService
	userCompany store.Cache
	car         *services.CarService
	logger      *log.Logger
}

func NewBookingHandler(mux *http.ServeMux, booking *services.BookingService, car *services.CarService, userCompanyCache store.Cache, logger *log.Logger) *BookingHandler {
	h := &BookingHandler{
		mux:         mux,
		car:         car,
		userCompany: userCompanyCache,
		booking:     booking,
		logger:      logger,
	}

	h.mux.HandleFunc("GET /booking/user/{id}", authMiddleware(h.handleGetUserBookings, logger))

	h.mux.HandleFunc("POST /booking", authMiddleware(h.handleCreateBooking, logger))
	h.mux.HandleFunc("GET /booking/{id}", authMiddleware(h.handleGetBookingByID, logger))
	h.mux.HandleFunc("DELETE /booking/{id}", authMiddleware(h.handleDeleteBookingByID, logger))
	h.mux.HandleFunc("PUT /booking/{id}", authMiddleware(h.handleUpdateBooking, logger))

	return h
}

func (h *BookingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// @Summary Create a booking
// @Description Creates a new booking record
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param payload body types.CreateBookingPayload true "Booking data"
// @Tags Booking
// @Success 200 {object} map[string]string
// @Router /booking [post]
func (h *BookingHandler) handleCreateBooking(w http.ResponseWriter, r *http.Request) error {
	var payload types.CreateBookingPayload
	if err := types.ParseJSON(r, &payload); err != nil {
		return types.InvalidJSON(err)
	}

	userId, ok := r.Context().Value(userIdKey).(int)
	if !ok {
		return types.Unauthorized("user id not found in token")
	}

	if errors := utils.ValidateStruct(&payload); len(errors) > 0 {
		return types.ValidationError(errors)
	}

	if err := h.booking.Create(userId, &payload); err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, map[string]string{"message": "booking created"})
}

// @Summary Get booking by ID
// @Description Retrieves a booking based on the provided ID
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "Booking ID"
// @Tags Booking
// @Success 200 {object} types.Booking
// @Router /booking/{id} [get]
func (h *BookingHandler) handleGetBookingByID(w http.ResponseWriter, r *http.Request) error {
	bookingID := r.PathValue("id")
	idInt, err := strconv.Atoi(bookingID)
	if err != nil {
		return types.BadPathParameter("id")
	}

	booking, err := h.booking.GetByID(idInt)
	if err != nil {
		return err
	}

	rentedCar, err := h.car.GetByID(booking.CarID)
	if err != nil {
		return err
	}

	userID := r.Context().Value(userIdKey)
	if booking.UserID != userID.(int) {
		// if user is an owner of the company, verify it via cache
		companyOwnerID, err := h.userCompany.Get(r.Context(), fmt.Sprintf("company:%d", rentedCar.CompanyID))
		if err != nil || companyOwnerID != userID {
			return types.Unauthorized("user does not have access to this booking")
		}
	}

	return types.WriteJSON(w, http.StatusOK, booking)
}

// @Summary Get bookings by user ID
// @Description Retrieves all bookings for a user
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "User ID"
// @Tags Booking
// @Success 200 {array} types.Booking
// @Router /booking/user/{id} [get]
func (h *BookingHandler) handleGetUserBookings(w http.ResponseWriter, r *http.Request) error {
	userID := r.PathValue("id")
	idInt, err := strconv.Atoi(userID)
	if err != nil {
		return types.BadPathParameter("id")
	}

	bookings, err := h.booking.GetByUserID(idInt)
	if err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, bookings)
}

// @Summary Update booking by ID
// @Description Updates an existing booking with new data
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "Booking ID"
// @Param payload body types.UpdateBookingPayload true "Updated booking data"
// @Tags Booking
// @Success 200 {object} map[string]string
// @Router /booking/{id} [put]
func (h *BookingHandler) handleUpdateBooking(w http.ResponseWriter, r *http.Request) error {

	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return types.BadPathParameter("id")
	}

	var payload types.UpdateBookingPayload
	if err := types.ParseJSON(r, &payload); err != nil {
		return types.InvalidJSON(err)
	}

	if errors := utils.ValidateStruct(&payload); len(errors) > 0 {
		return types.ValidationError(errors)
	}

	userId := r.Context().Value(userIdKey)
	booking, err := h.booking.GetByID(idInt)
	if err != nil {
		return err
	}

	if booking.UserID != userId.(int) {
		car, err := h.car.GetByID(booking.CarID)
		if err != nil {
			return err
		}
		companyOwnerID, err := h.userCompany.Get(r.Context(), fmt.Sprintf("company:%d", car.CompanyID))
		if err != nil || companyOwnerID != userId {
			return types.Unauthorized("user does not have access to this booking")
		}
	}

	if err := h.booking.Update(idInt, &payload); err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("booking %d updated", idInt)})
}

// @Summary Delete booking by ID
// @Description Deletes a booking by its ID
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "Booking ID"
// @Tags Booking
// @Success 200 {object} map[string]string
// @Router /booking/{id} [delete]
func (h *BookingHandler) handleDeleteBookingByID(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return types.BadPathParameter("id")
	}

	userId := r.Context().Value(userIdKey)
	booking, err := h.booking.GetByID(idInt)
	if err != nil {
		return err
	}

	if booking.UserID != userId.(int) {
		car, err := h.car.GetByID(booking.CarID)
		if err != nil {
			return err
		}
		companyOwnerID, err := h.userCompany.Get(r.Context(), fmt.Sprintf("company:%d", car.CompanyID))
		if err != nil || companyOwnerID != userId {
			return types.Unauthorized("user does not have access to this booking")
		}
	}

	if err := h.booking.Delete(idInt); err != nil {
		return err
	}

	return types.WriteJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("booking %d deleted", idInt)})

}
