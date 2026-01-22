package constants

import "net/http"

// APISuccess represents a standardized API success response with code and HTTP status.
// Use these predefined success constants for consistent API responses across the application.
type APISuccess struct {
	Code   string
	Status int
}

// Order-related success responses
var (
	SuccessOrderProcessed = APISuccess{
		Code:   CodeOrderProcessed,
		Status: http.StatusOK,
	}
)
