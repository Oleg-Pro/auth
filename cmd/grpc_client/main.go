package main

import (
	"context"
	"log"
	"time"

	desc "github.com/Oleg-Pro/auth/pkg/user_v1"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
	userID  = 1
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server #{err}")
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Failed to close connection %v", err)
		}
	}()

	client := desc.NewUserV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	r, err := client.Get(ctx, &desc.GetRequest{Id: userID})
	if err != nil {
		log.Fatalf("Failed to User %v", err)
	}

	log.Printf(color.RedString("User Info \n"), color.GreenString("%+v", r))
}
