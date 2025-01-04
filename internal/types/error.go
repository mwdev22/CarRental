package types

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

func InvalidJSON(err error) ApiError {
	return newApiError(http.StatusUnprocessableEntity, fmt.Errorf("invalid json, %v", err))
}

func InvalidFormData(err error) ApiError {
	return newApiError(http.StatusUnprocessableEntity, err)
}

func ServiceError(err error) ApiError {
	return newApiError(http.StatusInternalServerError, err)
}

func DatabaseError(err error) ApiError {
	return newApiError(http.StatusInternalServerError, fmt.Errorf("database error: %v", err))
}

func ExternalServiceErr(err error) ApiError {
	return newApiError(http.StatusBadRequest, err)
}

func Unauthorized(msg any) ApiError {
	return newApiError(http.StatusUnauthorized, fmt.Errorf("unauthorized, %s", msg))
}

func BadQueryParameter(name string) ApiError {
	return newApiError(http.StatusBadRequest, fmt.Errorf("bad query param, %s", name))
}

func BadPathParameter(name string) ApiError {
	return newApiError(http.StatusBadRequest, fmt.Errorf("bad path param, %s", name))
}

func NotFound(id int, name string) ApiError {
	return newApiError(http.StatusNotFound, fmt.Errorf("%s with id %v not found in database", name, id))
}
