package ports

import "context"

// MessageProcessor defines the inbound port (driving side)
// This is what the RabbitMQ consumer adapter will call
type MessageProcessor interface {
	ProcessMessage(ctx context.Context, messageData []byte) error
}
