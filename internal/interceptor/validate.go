package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

type validator interface {
	Validate() error
}

// ValidateInterceptor validation interceptor
func ValidateInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Incerceptor FullMethod : %s\n", info.FullMethod)
	if val, ok := req.(validator); ok {
		if err := val.Validate(); err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}
