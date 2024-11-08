package user_saver

import (
	"context"

	"github.com/Oleg-Pro/auth/internal/client/kafka"
	def "github.com/Oleg-Pro/auth/internal/service"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	userService def.UserService
	consumer    kafka.Consumer
	topicName   string
}

// NewService create User Saver Consumer
func NewService(
	userService def.UserService,
	consumer kafka.Consumer,
	topicName string,
) *service {
	return &service{
		userService: userService,
		consumer:    consumer,
		topicName:   topicName,
	}
}

// RunConsumer run consumer
func (s *service) RunConsumer(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-s.run(ctx):
			if err != nil {
				return err
			}

		}
	}
}

func (s *service) run(ctx context.Context) <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		errChan <- s.consumer.Consume(ctx, s.topicName, s.UserSaveHandler)
	}()

	return errChan
}
