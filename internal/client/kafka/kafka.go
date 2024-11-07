package kafka

import (
	"context"

	"github.com/Oleg-Pro/auth/internal/client/kafka/consumer"
)

type Consumer interface {
	Consume(ctx context.Context, topicName string, handler consumer.Handler) (err error)
	Close() error
}
