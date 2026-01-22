package http

import (
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/pkg/httputils"
	"net/http"
)

// RespondJSON writes a JSON success response (legacy wrapper).
// Consider using httputils.WriteAPISuccess for new code.
func RespondJSON(w http.ResponseWriter, statusCode int, data any) {
	httputils.RespondJSON(w, statusCode, data)
}

// RespondError writes a JSON error response (legacy wrapper).
// Consider using httputils.WriteAPIError for new code.
func RespondError(w http.ResponseWriter, statusCode int, message string, details map[string]string) {
	httputils.RespondError(w, statusCode, message, details)
}
