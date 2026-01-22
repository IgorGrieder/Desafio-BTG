package validation

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// Get returns the singleton validator instance.
// Uses sync.Once to ensure thread-safe initialization.
func Get() *validator.Validate {
	once.Do(func() {
		validate = validator.New(validator.WithRequiredStructEnabled())

		// Register custom validators here if needed
		// Example:
		// validate.RegisterValidation("custom_rule", customValidator)
	})
	return validate
}

// Validate validates a struct and returns an error if invalid.
// This is a convenience function that wraps Get().Struct(s).
func Validate(s any) error {
	return Get().Struct(s)
}

// ValidationErrors extracts field-level error details from a validation error.
// Returns a map of field names to error messages.
func ValidationErrors(err error) map[string]string {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		details := make(map[string]string)
		for _, e := range validationErrs {
			details[e.Field()] = formatValidationError(e)
		}
		return details
	}
	return nil
}

// formatValidationError formats a single validation error into a human-readable message.
func formatValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Value is too short or too small"
	case "max":
		return "Value is too long or too large"
	case "gt":
		return "Value must be greater than " + e.Param()
	case "gte":
		return "Value must be greater than or equal to " + e.Param()
	case "lt":
		return "Value must be less than " + e.Param()
	case "lte":
		return "Value must be less than or equal to " + e.Param()
	case "email":
		return "Invalid email format"
	case "url":
		return "Invalid URL format"
	case "uuid":
		return "Invalid UUID format"
	default:
		return e.Error()
	}
}
