package response

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func OK(w http.ResponseWriter, body any)      { JSON(w, http.StatusOK, body) }
func Created(w http.ResponseWriter, body any) { JSON(w, http.StatusCreated, body) }

func BadRequest(w http.ResponseWriter, msg string) {
	JSON(w, http.StatusBadRequest, map[string]string{"error": msg})
}

func TooManyRequests(w http.ResponseWriter) {
	JSON(w, http.StatusTooManyRequests, map[string]string{"error": "rate limit exceeded — try again later"})
}

func InternalError(w http.ResponseWriter) {
	JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
}

func ServiceUnavailable(w http.ResponseWriter, msg string) {
	JSON(w, http.StatusServiceUnavailable, map[string]string{"error": msg})
}
