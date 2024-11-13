package auth

import (
	"context"

	desc "github.com/Oleg-Pro/auth/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetRefreshToken get refresh token
func (i *Implemenation) GetRefreshToken(ctx context.Context, req *desc.GetRefreshTokenRequest) (*desc.GetRefreshTokenResponse, error) {
	refreshToken, err := i.authenticationService.GetRefreshToken(ctx, req.GetOldRefreshToken())
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	return &desc.GetRefreshTokenResponse{RefreshToken: refreshToken}, nil
}
