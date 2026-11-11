package response

import (
	"encoding/json"
	"net/http"
)

// SuccessResponse defines the standard success JSON structure.
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// ErrorResponse defines the standard error JSON structure.
type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

// PaginationMeta provides pagination metadata in responses.
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

// JSON sends a JSON response with status code.
func JSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

// OK sends a standard 200 OK success response.
func OK(w http.ResponseWriter, message string, data interface{}, meta interface{}) {
	JSON(w, http.StatusOK, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Created sends a standard 201 Created success response.
func Created(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusCreated, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends a standard error response with specified HTTP status code.
func Error(w http.ResponseWriter, statusCode int, message string, errs interface{}) {
	JSON(w, statusCode, ErrorResponse{
		Success: false,
		Message: message,
		Errors:  errs,
	})
}

// BadRequest sends a 400 Bad Request error.
func BadRequest(w http.ResponseWriter, message string, errs interface{}) {
	Error(w, http.StatusBadRequest, message, errs)
}

// Unauthorized sends a 401 Unauthorized error.
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, message, nil)
}

// Forbidden sends a 403 Forbidden error.
func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, message, nil)
}

// NotFound sends a 404 Not Found error.
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message, nil)
}

// InternalServerError sends a 500 Internal Server Error.
func InternalServerError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, message, nil)
}

