package main

import (
	"context"
	"log"

	"github.com/Oleg-Pro/auth/internal/app"
)

func main() {
	ctx := context.Background()
	a, err := app.NewCreateUserConsumerApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.Run(ctx)
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
