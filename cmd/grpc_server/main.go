package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/Oleg-Pro/auth/internal/config"
	"github.com/Oleg-Pro/auth/internal/converter"
	"github.com/Oleg-Pro/auth/internal/model"
	"github.com/Oleg-Pro/auth/internal/repository"
	"github.com/Oleg-Pro/auth/internal/repository/user"
	desc "github.com/Oleg-Pro/auth/pkg/user_v1"
	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedUserV1Server
	pool           *pgxpool.Pool
	userRepository repository.UserRepository
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	if req.GetPasword() != req.PasswordConfirm {
		err := fmt.Errorf("passwords are not equal")
		log.Printf("Error: %v", err)
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.GetPasword()), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Could not create password hash: %v", err)
		return nil, err
	}

	userID, err := s.userRepository.Create(ctx, &model.UserInfo{
		Name:        req.GetName(),
		Email:       req.GetEmail(),
		PaswordHash: string(passwordHash),
		Role:        model.Role(req.GetRole()),
	})

	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return nil, err
	}

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := s.userRepository.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return converter.ToUserGetResponseFromModelUser(user), nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	/*	builderUpdate := sq.Update(userTable).
			PlaceholderFormat(sq.Dollar).
			Set(userColumnUpdateAt, time.Now()).
			Set(userColumnRoleID, req.GetRole()).
			Where(sq.Eq{fmt.Sprintf(`"%s"`, userColumnID): req.GetId()})

		if req.GetName() != nil {
			builderUpdate = builderUpdate.Set(userColumnName, req.GetName().Value)
		}

		if req.GetEmail() != nil {
			log.Printf("Email: %v", req.GetEmail().Value)

			builderUpdate = builderUpdate.Set(userColumnEmail, req.GetEmail().Value)
		}

		query, args, err := builderUpdate.ToSql()

		if err != nil {
			log.Printf("Failed to build update query: %v", err)
			return nil, err
		}

		res, err := s.pool.Exec(ctx, query, args...)
		if err != nil {
			log.Printf("Failed to update user with id %d: %v", req.GetId(), err)
			return nil, err
		}

		log.Printf("updated %d rows", res.RowsAffected())*/

	var name, email *string
	if req.GetName() != nil {
		name = &req.GetName().Value
	}

	if req.GetEmail() != nil {
		email = &req.GetEmail().Value
	}

	_, err := s.userRepository.Update(ctx, req.GetId(), name, email, model.Role(req.GetRole()))
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	log.Printf("Deleting User req=%v", req)
	_, err := s.userRepository.Delete(ctx, req.GetId())
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		return nil, err
	}

	return &empty.Empty{}, nil
}

func main() {
	flag.Parse()

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

	desc.RegisterUserV1Server(s, &server{pool: pool, userRepository: userRepository})
	log.Printf("server listening at %v", listener.Addr())

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
