package handlers

import (
	"net/http"
	"testing"

	"github.com/mwdev22/CarRental/internal/types"
)

func TestCreateBooking(t *testing.T) {

	TestLogin(t)

	_, _ = generateCompanies(10, t)
	_, _ = generateCars(10, t)

	// test CreateBooking route
	url := testServer.URL + "/booking"

	payload := &types.CreateBookingPayload{
		CarID:     1,
		StartDate: "2021-01-01",
		EndDate:   "2021-01-02",
	}

	resp := sendPostRequest(url, payload, t)

	checkResponse(resp, http.StatusOK, t)
}

func TestGetBooking(t *testing.T) {
	url := testServer.URL + "/booking/1"

	resp := sendGetRequest(url, t)

	checkResponse(resp, 200, t)
}

func TestGetUserBookings(t *testing.T) {
	url := testServer.URL + "/booking/user/1"

	resp := sendGetRequest(url, t)

	checkResponse(resp, 200, t)
}

func TestUpdateBooking(t *testing.T) {
	url := testServer.URL + "/booking/1"

	payload := &types.UpdateBookingPayload{
		StartDate: "2021-01-01",
		EndDate:   "2021-01-05",
	}

	resp := sendPutRequest(url, payload, t)

	checkResponse(resp, 200, t)
}

func TestDeleteBooking(t *testing.T) {
	url := testServer.URL + "/booking/1"

	resp := sendDeleteRequest(url, t)

	checkResponse(resp, 200, t)
}
