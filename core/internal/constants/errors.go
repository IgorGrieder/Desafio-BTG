package constants

import "net/http"

// APIError represents a standardized API error with code, message, and HTTP status.
// Use these predefined errors for consistent API responses across the application.
type APIError struct {
	Code    string
	Message string
	Status  int
}

// WithMessage returns a copy of the APIError with a custom message.
// Useful for validation errors or other dynamic messages.
func (e APIError) WithMessage(message string) APIError {
	return APIError{
		Code:    e.Code,
		Message: message,
		Status:  e.Status,
	}
}

// Common errors - shared across multiple modules
var (
	ErrInvalidRequestBody = APIError{
		Code:    CodeInvalidRequest,
		Message: MsgInvalidRequestBody,
		Status:  http.StatusBadRequest,
	}
	ErrInternalError = APIError{
		Code:    CodeInternalError,
		Message: MsgInternalError,
		Status:  http.StatusInternalServerError,
	}
	ErrNotFound = APIError{
		Code:    CodeNotFound,
		Message: MsgNotFound,
		Status:  http.StatusNotFound,
	}
)

// Order-related errors
var (
	ErrOrderNotFound = APIError{
		Code:    CodeOrderNotFound,
		Message: MsgOrderNotFound,
		Status:  http.StatusNotFound,
	}
	ErrInvalidOrderCode = APIError{
		Code:    CodeInvalidOrderCode,
		Message: MsgInvalidOrderCode,
		Status:  http.StatusBadRequest,
	}
	ErrInvalidCustomerCode = APIError{
		Code:    CodeInvalidCustomerCode,
		Message: MsgInvalidCustomerCode,
		Status:  http.StatusBadRequest,
	}
	ErrFailedToCreateOrder = APIError{
		Code:    CodeInternalError,
		Message: MsgFailedToCreateOrder,
		Status:  http.StatusInternalServerError,
	}
	ErrFailedToGetOrder = APIError{
		Code:    CodeInternalError,
		Message: MsgFailedToGetOrder,
		Status:  http.StatusInternalServerError,
	}
	ErrFailedToGetOrderTotal = APIError{
		Code:    CodeInternalError,
		Message: MsgFailedToGetOrderTotal,
		Status:  http.StatusInternalServerError,
	}
	ErrFailedToListOrders = APIError{
		Code:    CodeInternalError,
		Message: MsgFailedToListOrders,
		Status:  http.StatusInternalServerError,
	}
	ErrFailedToCountOrders = APIError{
		Code:    CodeInternalError,
		Message: MsgFailedToCountOrders,
		Status:  http.StatusInternalServerError,
	}
)
