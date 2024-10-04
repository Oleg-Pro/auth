package main

import (
	"context"
	"log"
	"time"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	desc "github.com/Oleg-Pro/auth/pkg/user_v1"	
)

const (
	address = "localhost:50051"	
	userId = 1
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server #{err}")
	}

	defer conn.Close()

	client :=desc.NewUserV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	r, err := client.Get(ctx, &desc.GetRequest{Id: userId})
	if err != nil {				
		log.Fatalf("Failed to User #{err}")
	}

	log.Printf(color.RedString("User Info \n"), color.GreenString("%+v", r))






}