package config

import (
	"errors"
	"os"
	"strings"

	"github.com/IBM/sarama"
)

const (
	brokersEnvName = "KAFKA_BROKERS"
	groupIDEnvName = "KAFKA_GROUP_ID"
	topicNameName  = "KAFKA_CREATE_USER_TOPIC"
)

type KafkaConsumerConfig interface {
	Brokers() []string
	GroupID() string
	TopicName() string
	Config() *sarama.Config
}

type kafkaConsumerConfig struct {
	brokers   []string
	groupID   string
	topicName string
}

func NewKafkaConsumerConfig() (*kafkaConsumerConfig, error) {
	brokersStr := os.Getenv(brokersEnvName)
	if len(brokersStr) == 0 {
		return nil, errors.New("kafka brokers address not found")
	}

	brokers := strings.Split(brokersStr, ",")

	groupID := os.Getenv(groupIDEnvName)
	if len(groupID) == 0 {
		return nil, errors.New("kafka group id not found")
	}

	topicName := os.Getenv(topicNameName)
	if len(topicName) == 0 {
		return nil, errors.New("kafka topic name not found")
	}

	return &kafkaConsumerConfig{
		brokers:   brokers,
		groupID:   groupID,
		topicName: topicName,
	}, nil
}

func (cfg *kafkaConsumerConfig) Brokers() []string {
	return cfg.brokers
}

func (cfg *kafkaConsumerConfig) GroupID() string {
	return cfg.groupID
}

func (cfg *kafkaConsumerConfig) TopicName() string {
	return cfg.topicName
}

// Config возвращает конфигурацию для sarama consumer
func (cfg *kafkaConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
