package access

import (
	"context"
	"strings"

	"github.com/Oleg-Pro/auth/internal/model"
	desc "github.com/Oleg-Pro/auth/pkg/access_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	authPrefix = "Bearer "
)

// Check check access to endpoint
func (i *Implemenation) Check(ctx context.Context, req *desc.CheckRequest) (*emptypb.Empty, error) {
	accessToken, err := getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	if !i.accessService.Allow(ctx, req.GetEndpointAddress(), accessToken) {
		return nil, model.ErrorAccessDenied
	}

	return &emptypb.Empty{}, nil
}

func getAccessToken(ctx context.Context) (string, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", model.ErrorMetadataNotProvided
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return "", model.ErrorAuthorizationHeaderNotProvided
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return "", model.ErrorAuthorizationHeaderFormat
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	return accessToken, nil
}
