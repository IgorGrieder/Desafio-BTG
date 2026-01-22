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
	SuccessOrderCreated = APISuccess{
		Code:   CodeOrderCreated,
		Status: http.StatusCreated,
	}
	SuccessOrderFound = APISuccess{
		Code:   CodeOrderFound,
		Status: http.StatusOK,
	}
	SuccessOrdersListed = APISuccess{
		Code:   CodeOrdersListed,
		Status: http.StatusOK,
	}
	SuccessOrderCounted = APISuccess{
		Code:   CodeOrderCounted,
		Status: http.StatusOK,
	}
)
