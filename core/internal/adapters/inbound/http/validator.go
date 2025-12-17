package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Data any `json:"data"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Generic method to invoke struct validation based on parameters
func ValidateStruct(s any) error {
	return validate.Struct(s)
}

// Function to responde based on the errors
func RespondValidationError(w http.ResponseWriter, err error) {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		details := make(map[string]string)
		for _, e := range validationErrs {
			details[e.Field()] = e.Error()
		}
		RespondError(w, http.StatusBadRequest, "Validation failed", details)
		return
	}
	RespondError(w, http.StatusBadRequest, "Validation failed", nil)
}
