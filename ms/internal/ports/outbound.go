package ports

import "context"

// ProcessRepository defines the outbound port (driven side) for database
// This is what the application needs from the database adapter
type ProcessRepository interface {
	SaveProcessedMessage(ctx context.Context, messageID string, data any) error
	GetProcessedMessage(ctx context.Context, messageID string) (any, error)
}
