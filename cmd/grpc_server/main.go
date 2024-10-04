package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"

	desc "github.com/Oleg-Pro/auth/pkg/user_v1"
	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	empty "github.com/golang/protobuf/ptypes/empty"	
)


const grcPort = 50051


type server struct {
	desc.UnimplementedUserV1Server
}

func(s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	return &desc.CreateResponse{
		Id: rand.Int63(),
	}, nil
}


func (s * server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %+v", req)

	return &desc.GetResponse{
		Id: req.GetId(),
		Name: gofakeit.Name(),
		Email: gofakeit.Email(),
		Role: desc.Role_USER,		
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),					

	}, nil
}

func(s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {

	log.Printf("Id=%d: Name=%s Email=%s Role=%d", req.GetId(), req.GetName(), req.GetEmail(), req.GetRole())		
	return &empty.Empty{}, nil
}


func(s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {

	log.Printf("Deleting User with Id=%d", req.GetId())		
	return &empty.Empty{}, nil
}


func main() {
	fmt.Println("Server:" , grcPort)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grcPort))
	if err != nil {
		log.Fatalf("Failed to listen #{err}")
	}

	s := grpc.NewServer()
	reflection.Register(s)
    desc.RegisterUserV1Server(s, &server{})
	log.Printf("server listening at %v", listener.Addr())	

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve #{err}")
	}
}

