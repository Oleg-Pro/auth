package app

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	authAPI "github.com/Oleg-Pro/auth/internal/api/auth"		
	userAPI "github.com/Oleg-Pro/auth/internal/api/user"
	userToken "github.com/Oleg-Pro/auth/internal/service/user/token"	
	"github.com/Oleg-Pro/auth/internal/client/cache"
	"github.com/Oleg-Pro/auth/internal/client/cache/redis"
	"github.com/Oleg-Pro/auth/internal/client/kafka"
	kafkaConsumer "github.com/Oleg-Pro/auth/internal/client/kafka/consumer"
	"github.com/Oleg-Pro/auth/internal/config"
	"github.com/Oleg-Pro/auth/internal/repository"
	userRepository "github.com/Oleg-Pro/auth/internal/repository/user"
	userCacheRepository "github.com/Oleg-Pro/auth/internal/repository/user/redis"
	"github.com/Oleg-Pro/auth/internal/service"
	userSaverConsumer "github.com/Oleg-Pro/auth/internal/service/consumer/user_saver"
	userSaverProducer "github.com/Oleg-Pro/auth/internal/service/producer/user_saver"
	userService "github.com/Oleg-Pro/auth/internal/service/user"
	"github.com/Oleg-Pro/auth/internal/service/authentication"
	"github.com/Oleg-Pro/platform-common/pkg/closer"
	"github.com/Oleg-Pro/platform-common/pkg/db"
	"github.com/Oleg-Pro/platform-common/pkg/db/pg"
	"github.com/Oleg-Pro/platform-common/pkg/db/transaction"
	redigo "github.com/gomodule/redigo/redis"
)

const kafkaProducerRetryMax = 5

type serviceProvider struct {
	pgConfig            config.PGConfig
	kafkaConsumerConfig config.KafkaConsumerConfig
	grpcConfig          config.GRPCConfig
	httpConfig          config.HTTPConfig
	swaggerConfig       config.SwaggerConfig
	redisConfig         config.RedisConfig

	dbClient    db.Client
	txManager   db.TxManager
	redisPool   *redigo.Pool
	redisClient cache.RedisClient

	userSaverConsumer    service.ConsumerService
	consumer             kafka.Consumer
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *kafkaConsumer.GroupHandler

	producer sarama.SyncProducer

	userSaverProducer userSaverProducer.UserSaverProducer

	userRepository repository.UserRepository

	userCacheRepository repository.UserCacheRepository
	userService         service.UserService
	userImplemenation   *userAPI.Implementation
	authImplemenation   *authAPI.Implemenation

	userTokenService  service.UserTokenService
	authenticationService service.AuthenticationService

}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := config.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := config.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %s", err.Error())
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		client, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = client.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(client.Close)

		s.dbClient = client
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}
	return s.txManager
}

func (s *serviceProvider) RedisPool() *redigo.Pool {
	if s.redisPool == nil {
		s.redisPool = &redigo.Pool{
			MaxIdle:     s.RedisConfig().MaxIdle(),
			IdleTimeout: s.RedisConfig().IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", s.RedisConfig().Address())
			},
		}
	}

	return s.redisPool
}

func (s *serviceProvider) RedisClient() cache.RedisClient {
	if s.redisClient == nil {
		s.redisClient = redis.NewClient(s.RedisPool(), s.RedisConfig())
	}

	return s.redisClient
}

func (s *serviceProvider) KafkaConsumerConfig() config.KafkaConsumerConfig {
	if s.kafkaConsumerConfig == nil {
		cfg, err := config.NewKafkaConsumerConfig()
		if err != nil {
			log.Fatalf("failed to get kafka consumer config: %s", err.Error())
		}

		s.kafkaConsumerConfig = cfg
	}

	return s.kafkaConsumerConfig
}

func (s *serviceProvider) UserSaverConsumer(ctx context.Context) service.ConsumerService {
	if s.userSaverConsumer == nil {
		s.userSaverConsumer = userSaverConsumer.NewService(
			s.UserService(ctx),
			s.Consumer(),
			s.KafkaConsumerConfig().TopicName(),
		)
	}

	return s.userSaverConsumer
}

func (s *serviceProvider) Consumer() kafka.Consumer {
	if s.consumer == nil {
		s.consumer = kafkaConsumer.NewConsumer(
			s.ConsumerGroup(),
			s.ConsumerGroupHandler(),
		)
		closer.Add(s.consumer.Close)
	}

	return s.consumer
}

func (s *serviceProvider) ConsumerGroup() sarama.ConsumerGroup {
	if s.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			s.KafkaConsumerConfig().Brokers(),
			s.KafkaConsumerConfig().GroupID(),
			s.KafkaConsumerConfig().Config(),
		)
		if err != nil {
			log.Fatalf("failed to create consumer group: %v", err)
		}

		s.consumerGroup = consumerGroup
	}

	return s.consumerGroup
}

func (s *serviceProvider) ConsumerGroupHandler() *kafkaConsumer.GroupHandler {
	if s.consumerGroupHandler == nil {
		s.consumerGroupHandler = kafkaConsumer.NewGroupHandler()
	}

	return s.consumerGroupHandler
}

func (s *serviceProvider) Producer(retryMax int) sarama.SyncProducer {
	if s.producer == nil {
		producer, err := newSyncProducer(s.KafkaConsumerConfig().Brokers(), retryMax)
		if err != nil {
			log.Fatalf("failed to start producer: %v\n", err.Error())
		}

		s.producer = producer
		closer.Add(s.producer.Close)
	}

	return s.producer
}

func (s *serviceProvider) UserSaverProducer(retryMax int) userSaverProducer.UserSaverProducer {
	if s.userSaverProducer == nil {
		s.userSaverProducer = userSaverProducer.NewUserSaverProducer(s.Producer(retryMax), s.KafkaConsumerConfig().TopicName())
	}

	return s.userSaverProducer
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserCacheRepository(_ context.Context) repository.UserCacheRepository {
	if s.userCacheRepository == nil {
		s.userCacheRepository = userCacheRepository.NewRepository(s.RedisClient())
	}

	return s.userCacheRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.New(s.UserRepository(ctx), s.UserCacheRepository(ctx))
	}

	return s.userService
}

func (s *serviceProvider) UserImplementation(ctx context.Context) *userAPI.Implementation {

	if s.userImplemenation == nil {
		s.userImplemenation = userAPI.NewImplementation(s.UserService(ctx), s.UserSaverProducer(kafkaProducerRetryMax))
	}
	return s.userImplemenation
}

func (s * serviceProvider) UserTokenService() service.UserTokenService {
	if s.userTokenService == nil {
		s.userTokenService = userToken.New()
	}

	return s.userTokenService
}

func (s * serviceProvider) AuthenticationService() service.AuthenticationService {
	if s.authenticationService == nil {
		s.authenticationService = authentication.New(s.UserTokenService())
	}

	return s.authenticationService
}


func (s *serviceProvider) AuthImplementation(ctx context.Context) *authAPI.Implemenation {

	if s.authImplemenation == nil {
		s.authImplemenation = authAPI.NewImplementation(s.AuthenticationService())
	}
	return s.authImplemenation
}

func newSyncProducer(brokerList []string, retryMax int) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = retryMax
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
