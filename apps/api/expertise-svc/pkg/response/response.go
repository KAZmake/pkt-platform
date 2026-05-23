package response

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func OK(w http.ResponseWriter, payload any)      { JSON(w, http.StatusOK, payload) }
func Created(w http.ResponseWriter, payload any) { JSON(w, http.StatusCreated, payload) }

func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, map[string]string{"error": msg})
}

func BadRequest(w http.ResponseWriter, msg string) { Error(w, http.StatusBadRequest, msg) }
func Unauthorized(w http.ResponseWriter)           { Error(w, http.StatusUnauthorized, "unauthorized") }
func Forbidden(w http.ResponseWriter)              { Error(w, http.StatusForbidden, "forbidden") }
func NotFound(w http.ResponseWriter)               { Error(w, http.StatusNotFound, "not found") }
func InternalError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "internal server error")
}
