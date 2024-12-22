package handlers

import (
	"fmt"
	"net/http"
)

type ApiError struct {
	StatusCode int `json:"status_code"`
	Msg        any `json:"msg"`
}

func (e ApiError) Error() string {
	return fmt.Sprintf("%v", e.Msg)
}

func newApiError(statusCode int, e error) ApiError {
	return ApiError{
		StatusCode: statusCode,
		Msg:        e.Error(),
	}
}

func invalidJSON(err error) ApiError {
	return newApiError(http.StatusUnprocessableEntity, err)
}

func invalidFormData(err error) ApiError {
	return newApiError(http.StatusUnprocessableEntity, err)
}

func serviceError(err error) ApiError {
	return newApiError(http.StatusInternalServerError, err)
}

func externalServiceErr(err error) ApiError {
	return newApiError(http.StatusBadRequest, err)
}

func unauthorized(msg any) ApiError {
	return newApiError(http.StatusUnauthorized, fmt.Errorf("unauthorized, %s", msg))
}

func badQueryParameter(name string) ApiError {
	return newApiError(http.StatusBadRequest, fmt.Errorf("bad query param, %s", name))
}

func badPathParameter(name string) ApiError {
	return newApiError(http.StatusBadRequest, fmt.Errorf("bad path param, %s", name))
}

func notFound(id int, name string) ApiError {
	return newApiError(http.StatusNotFound, fmt.Errorf("%s with id %v not found in database", name, id))
}
