package http

import (
	"net/http"

	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/infrastructure/validation"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/pkg/httputils"
)

// ValidateStruct validates a struct using the centralized validator.
// Returns an error if validation fails.
func ValidateStruct(s any) error {
	return validation.Validate(s)
}

// RespondValidationError writes a validation error response with field details.
// Uses the new httputils package for consistent response format.
func RespondValidationError(w http.ResponseWriter, err error) {
	details := validation.ValidationErrors(err)
	if details != nil {
		httputils.RespondError(w, http.StatusBadRequest, "Validation failed", details)
		return
	}
	httputils.RespondError(w, http.StatusBadRequest, "Validation failed", nil)
}
