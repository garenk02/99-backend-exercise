// domain/errors.go
package domain

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// ErrorResponse is the standard error response format
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// RespondWithError writes an error response in JSON format
func RespondWithError(w http.ResponseWriter, code int, message string, err error) {
	// Log the error
	slog.Error("API error",
		"status_code", code,
		"message", message,
		"error", err,
	)

	// Create error response
	errResp := ErrorResponse{
		Error:   http.StatusText(code),
		Code:    code,
		Message: message,
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errResp)
}
