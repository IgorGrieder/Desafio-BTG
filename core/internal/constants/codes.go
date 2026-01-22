package constants

// Error codes used in API responses.
// These are the machine-readable codes returned in the "error" field.
const (
	// Common error codes
	CodeInvalidRequest = "INVALID_REQUEST"
	CodeInternalError  = "INTERNAL_ERROR"
	CodeForbidden      = "FORBIDDEN"
	CodeNotFound       = "NOT_FOUND"

	// Order-specific codes
	CodeOrderNotFound       = "ORDER_NOT_FOUND"
	CodeInvalidOrderCode    = "INVALID_ORDER_CODE"
	CodeInvalidCustomerCode = "INVALID_CUSTOMER_CODE"

	// Success codes - Order operations
	CodeOrderCreated = "ORDER_CREATED"
	CodeOrderFound   = "ORDER_FOUND"
	CodeOrdersListed = "ORDERS_LISTED"
	CodeOrderCounted = "ORDER_COUNTED"
)
