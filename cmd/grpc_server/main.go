package main

import (
	/*	"context"
		"flag"
		"log"
		"net"

		userAPI "github.com/Oleg-Pro/auth/internal/api/user"
		"github.com/Oleg-Pro/auth/internal/config"
		"github.com/Oleg-Pro/auth/internal/repository/user"
		userService "github.com/Oleg-Pro/auth/internal/service/user"
		desc "github.com/Oleg-Pro/auth/pkg/user_v1"
		"github.com/jackc/pgx/v4/pgxpool"
		"google.golang.org/grpc"
		"google.golang.org/grpc/reflection"*/

	"context"
	"log"

	"github.com/Oleg-Pro/auth/internal/app"
)

/*var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}*/

func main() {
	/*	flag.Parse()

		err := config.Load(configPath)
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		grpcConfig, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}

		pgConfig, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}

		ctx := context.Background()

		pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		defer pool.Close()

		listener, err := net.Listen("tcp", grpcConfig.Address())
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		s := grpc.NewServer()
		reflection.Register(s)

		userRepository := user.NewRepository(pool)
		userService := userService.New(userRepository)

		desc.RegisterUserV1Server(s, userAPI.NewImplementation(userService))
		log.Printf("server listening at %v", listener.Addr())

		if err := s.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}*/

	ctx := context.Background()
	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
