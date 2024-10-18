package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Oleg-Pro/auth/internal/config"
	desc "github.com/Oleg-Pro/auth/pkg/user_v1"
	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

const (
	userTable = "users"

	userColumnID           = "id"
	userColumnName         = "name"
	userColumnEmail        = "email"
	userColumnRoleID       = "role_id"
	userColumnCreatedAt    = "created_at"
	userColumnUpdateAt     = "updated_at"
	userColumnPasswordHash = "password_hash"
)

type server struct {
	desc.UnimplementedUserV1Server
	pool *pgxpool.Pool
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

	builderInsert := sq.Insert(userTable).
		PlaceholderFormat(sq.Dollar).
		Columns(userColumnName, userColumnEmail, userColumnPasswordHash, userColumnRoleID).
		Values(req.GetName(), req.GetEmail(), passwordHash, req.GetRole()).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("Failed to build insert query: %v", err)
		return nil, err
	}

	var userID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return nil, err
	}

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	builderSelectOne := sq.Select(userColumnID, userColumnName, userColumnEmail, userColumnRoleID, userColumnCreatedAt, userColumnUpdateAt).
		From(userTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{fmt.Sprintf(`"%s"`, userColumnID): req.GetId()}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		log.Printf("Failed to build get query: %v", err)
		return nil, err
	}

	var id int64
	var roleID int32
	var name, email string
	var createdAt time.Time
	var updatedAt sql.NullTime

	err = s.pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &roleID, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return nil, err
	}

	var updateAtTime *timestamppb.Timestamp
	if updatedAt.Valid {
		updateAtTime = timestamppb.New(updatedAt.Time)
	}

	return &desc.GetResponse{
		Id:        id,
		Name:      name,
		Email:     email,
		Role:      desc.Role(roleID),
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: updateAtTime,
	}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	builderUpdate := sq.Update(userTable).
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

	log.Printf("updated %d rows", res.RowsAffected())

	return &empty.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	builderDelete := sq.Delete(userTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{fmt.Sprintf(`"%s"`, userColumnID): req.GetId()})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Printf("Failed to build delete query: %v", err)
		return nil, err
	}

	log.Printf("DELETE SQL query: %s", query)

	res, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("Failed to delete user with id %d: %v", req.GetId(), err)
		return nil, err
	}

	log.Printf("delete %d rows", res.RowsAffected())

	log.Printf("Deleting User req=%v", req)
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
	desc.RegisterUserV1Server(s, &server{pool: pool})
	log.Printf("server listening at %v", listener.Addr())

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
