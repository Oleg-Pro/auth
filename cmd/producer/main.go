package main

import (
	"encoding/json"
	"flag"
	"log"

	"github.com/IBM/sarama"
	"github.com/Oleg-Pro/auth/internal/config"
	"github.com/Oleg-Pro/auth/internal/model"
	"github.com/brianvoe/gofakeit/v6"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
	flag.Parse()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Load confi error: %v", err)
	}

	kafkaConfig, err := config.NewKafkaConsumerConfig()
	if err != nil {
		log.Fatalf("failed to get kafka consumer config: %s", err.Error())
	}

	log.Printf("Kafka Config: %#v", kafkaConfig)

	producer, err := newSyncProducer(kafkaConfig.Brokers())
	if err != nil {
		log.Fatalf("failed to start producer: %v\n", err.Error())
	}

	defer func() {
		if err = producer.Close(); err != nil {
			log.Fatalf("failed to close producer: %v\n", err.Error())
		}
	}()

	info := model.UserInfo{
		Name:        gofakeit.Name(),
		Email:       gofakeit.Email(),
		PaswordHash: "$2a$10$ovdAkan0WZ4LSNkOrd1hLuYQpjq6Ree1GES/6GPU3GcEO1XIzjjFG",
		Role:        model.RoleUSER,
	}

	data, err := json.Marshal(info)
	if err != nil {
		log.Fatalf("failed to marshal data: %v\n", err.Error())
	}

	msg := &sarama.ProducerMessage{
		Topic: kafkaConfig.TopicName(),
		Value: sarama.StringEncoder(data),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("failed to send message in Kafka: %v\n", err.Error())
		return
	}

	log.Printf("message sent to partition %d with offset %d\n", partition, offset)
}

func newSyncProducer(brokerList []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
