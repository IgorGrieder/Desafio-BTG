package constants

// Error messages used in API responses.
// These are the human-readable messages returned in the "message" field.
const (
	// Common messages
	MsgInvalidRequestBody = "Invalid request body"
	MsgInternalError      = "An internal error occurred"
	MsgNotFound           = "Resource not found"

	// Order-specific messages
	MsgOrderNotFound        = "Order not found"
	MsgInvalidOrderCode     = "Order code must be a positive integer"
	MsgInvalidCustomerCode  = "Customer code must be a positive integer"
	MsgFailedToProcessOrder = "Failed to process order"
	MsgFailedToSaveOrder    = "Failed to save order"
)
