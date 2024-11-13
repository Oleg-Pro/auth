package auth

import (
	"context"

	desc "github.com/Oleg-Pro/auth/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAccessToken get access token
func (i *Implemenation) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	accessToken, err := i.authenticationService.GetRefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	return &desc.GetAccessTokenResponse{AccessToken: accessToken}, nil
}
