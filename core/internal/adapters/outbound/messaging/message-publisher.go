package messaging

import "context"

type MessagePublisher struct {
}

func (mp *MessagePublisher) PublishOrder(ctx context.Context, orderID string, data any) error {
	return nil
}
