package http

import (
	"encoding/json"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := SuccessResponse{Data: data}
	json.NewEncoder(w).Encode(response)
}

func RespondError(w http.ResponseWriter, statusCode int, message string, details map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Error:   message,
		Details: details,
	}
	json.NewEncoder(w).Encode(response)
}
