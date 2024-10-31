package app

import (
	"context"
	"log"

	redigo "github.com/gomodule/redigo/redis"
	userAPI "github.com/Oleg-Pro/auth/internal/api/user"
	"github.com/Oleg-Pro/auth/internal/config"
	"github.com/Oleg-Pro/auth/internal/repository"
	userRepository "github.com/Oleg-Pro/auth/internal/repository/user"
	userCacheRepository "github.com/Oleg-Pro/auth/internal/repository/user/redis"	
	"github.com/Oleg-Pro/auth/internal/client/cache/redis"		
	"github.com/Oleg-Pro/auth/internal/service"
	userService "github.com/Oleg-Pro/auth/internal/service/user"
	"github.com/Oleg-Pro/platform-common/pkg/closer"
	"github.com/Oleg-Pro/platform-common/pkg/db"
	"github.com/Oleg-Pro/platform-common/pkg/db/pg"	
	"github.com/Oleg-Pro/platform-common/pkg/db/transaction"
	"github.com/Oleg-Pro/auth/internal/client/cache"	
)

type serviceProvider struct {
	pgConfig    config.PGConfig
	grpcConfig  config.GRPCConfig
	redisConfig config.RedisConfig

	dbClient       db.Client
	txManager      db.TxManager
	redisPool   *redigo.Pool	
	redisClient cache.RedisClient	

	userRepository repository.UserRepository

	userCacheRepository repository.UserCacheRepository
	userService       service.UserService
	userImplemenation *userAPI.Implementation
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

func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := config.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %s", err.Error())
		}

		s.redisConfig = cfg
	}

	log.Printf("Redis Config: %#v\n", s.redisConfig)

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

	log.Printf("Redis Pool: %#v\n", s.redisPool)	

	return s.redisPool
}

func (s *serviceProvider) RedisClient() cache.RedisClient {
	if s.redisClient == nil {
		s.redisClient = redis.NewClient(s.RedisPool(), s.RedisConfig())
	}

	log.Printf("Redis Client: %#v\n", s.redisClient)		

	return s.redisClient
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserCacheRepository(ctx context.Context) repository.UserCacheRepository {
	if s.userCacheRepository == nil {
		s.userCacheRepository = userCacheRepository.NewRepository(s.RedisClient())
	}

	log.Printf("UserCacheRepository: %#v\n", s.userCacheRepository)			
	return s.userCacheRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.New(s.UserRepository(ctx))
	}

	return s.userService
}

func (s *serviceProvider) UserImplementation(ctx context.Context) *userAPI.Implementation {

	if s.userImplemenation == nil {
		s.userImplemenation = userAPI.NewImplementation(s.UserService(ctx))
	}
	return s.userImplemenation
}
