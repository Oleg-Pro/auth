package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	desc "github.com/Oleg-Pro/auth/pkg/user_v1"
	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grcPort = 50051

const (
	dbDSN     = "host=localhost port=54321 dbname=auth user=auth-user password=auth-password sslmode=disable"
	userTable = "users"
)

type server struct {
	desc.UnimplementedUserV1Server
	pool *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Create User req=%v\n", req)

	if req.GetPasword() != req.PasswordConfirm {
		err := fmt.Errorf("passwords are not equal")
		log.Printf("Error: %v", err)
		return &desc.CreateResponse{}, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.GetPasword()), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Could not create password hash: %v", err)
		return &desc.CreateResponse{}, err
	}

	builderInsert := sq.Insert(userTable).
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "password_hash", "role_id").
		Values(req.GetName(), req.GetEmail(), passwordHash, req.GetRole()).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("Failed to build insert query: %v", err)
		return &desc.CreateResponse{}, err
	}

	var userID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return &desc.CreateResponse{}, err
	}

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("Get User req=%v", req)

	builderSelectOne := sq.Select("id", "name", "email", "role_id", "created_at", "updated_at").
		From(userTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		log.Printf("Failed to build get query: %v", err)
		return &desc.GetResponse{}, err
	}

	var id int64
	var roleID int32
	var name, email string
	var createdAt time.Time
	var updatedAt sql.NullTime

	err = s.pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &roleID, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return &desc.GetResponse{}, err
	}

	log.Printf("id: %d, name: %s, email: %s, roleId: %d, created_at: %v, updated_at: %v\n", id, name, email, roleID, createdAt, updatedAt)

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

	log.Printf("Update User req=%v", req)

	builderUpdate := sq.Update(userTable).
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", time.Now()).
		Set("role_id", req.GetRole()).
		Where(sq.Eq{"id": req.GetId()})

	if req.GetName() != nil {
		builderUpdate.Set("name", req.GetName().Value)
	}

	if req.GetEmail() != nil {
		log.Printf("Email: %v", req.GetEmail().Value)
		builderUpdate.Set("email", req.GetEmail().Value)
	}

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Printf("Failed to build update query: %v", err)
		return &empty.Empty{}, err
	}

	res, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("Failed to update user with id %d: %v", req.GetId(), err)
		return &empty.Empty{}, err
	}

	log.Printf("updated %d rows", res.RowsAffected())

	return &empty.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {

	builderDelete := sq.Delete(userTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Printf("Failed to build delete query: %v", err)
		return &empty.Empty{}, err
	}

	log.Printf("DELETE SQL query: %s", query)

	res, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("Failed to delete user with id %d: %v", req.GetId(), err)
		return &empty.Empty{}, err
	}

	log.Printf("delete %d rows", res.RowsAffected())

	log.Printf("Deleting User req=%v", req)
	return &empty.Empty{}, nil
}

func main() {
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grcPort))
	if err != nil {
		log.Fatalf("Failed to listen #{err}")
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{pool: pool})
	log.Printf("server listening at %v", listener.Addr())

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve #{err}")
	}
}
