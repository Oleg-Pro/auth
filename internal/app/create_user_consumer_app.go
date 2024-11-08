package app

import (
	"context"
	"log"
	"sync"

	"github.com/Oleg-Pro/auth/internal/config"
	"github.com/Oleg-Pro/platform-common/pkg/closer"
)

// CreateUserConsumerApp application for kafka consumer
type CreateUserConsumerApp struct {
	serviceProvider *serviceProvider
}

// NewCreateUserConsumerApp CreateUserConsumerApp constructor
func NewCreateUserConsumerApp(ctx context.Context) (*CreateUserConsumerApp, error) {
	a := &CreateUserConsumerApp{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run run app
func (a *CreateUserConsumerApp) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	ctx, cancel := context.WithCancel(ctx)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := a.serviceProvider.UserSaverConsumer(ctx).RunConsumer(ctx)
		if err != nil {
			log.Printf("failed to run consumer: %s", err.Error())
		}
	}()

	gracefulShutdown(ctx, cancel, wg)
	return nil
}

func (a *CreateUserConsumerApp) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *CreateUserConsumerApp) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (a *CreateUserConsumerApp) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}
