package user_saver

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/Oleg-Pro/auth/internal/model"
)

type service struct {
	producer  sarama.SyncProducer
	topicName string
}

// NewUserSaverProducer UserSaverProducer constructor
func NewUserSaverProducer(producer sarama.SyncProducer, topicName string) *service {
	return &service{
		producer:  producer,
		topicName: topicName,
	}
}

// UserSaverProducer interface
type UserSaverProducer interface {
	Send(ctx context.Context, info *model.UserInfo) error
}

func (a *service) Send(_ context.Context, info *model.UserInfo) error {

	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: a.topicName,
		Value: sarama.StringEncoder(data),
	}

	partition, offset, err := a.producer.SendMessage(msg)
	if err != nil {
		log.Printf("failed to send message in Kafka: %v\n", err.Error())
		return err
	}

	log.Printf("message sent to partition %d with offset %d\n", partition, offset)

	return nil
}
