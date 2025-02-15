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

func ValidationError(errors map[string]string) ApiError {
	return ApiError{
		StatusCode: http.StatusBadRequest,
		Msg:        errors,
	}
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

func BadRequest(msg any) ApiError {
	return newApiError(http.StatusBadRequest, fmt.Errorf("bad request: %s", msg))
}

func InternalServerError(msg string) ApiError {
	return newApiError(http.StatusInternalServerError, fmt.Errorf("internal server error: %s", msg))
}

func NotFound(msg string) ApiError {
	return newApiError(http.StatusNotFound, fmt.Errorf("not found: %s", msg))
}
