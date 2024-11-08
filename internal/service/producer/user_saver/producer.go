package user_saver

import (
	"context"
	"encoding/json"		
	"log"
	"github.com/IBM/sarama"
	"github.com/Oleg-Pro/auth/internal/model"
)

type service struct {
	producer sarama.SyncProducer
	topicName string
}

func NewUserSaverProducer(producer sarama.SyncProducer, topicName string) *service {	
	return &service{
		producer: producer,
		topicName: topicName,
	}
}

// UserSaverProducer
type UserSaverProducer interface {
	Send(ctx context.Context, info *model.UserInfo) error
}

func (a *service) Send(ctx context.Context, info *model.UserInfo) error {

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


/*func newSyncProducer(brokerList []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}*/
