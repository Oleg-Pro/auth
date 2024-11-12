package auth

import (
	"context"
	desc "github.com/Oleg-Pro/auth/pkg/auth_v1"		
)

func (i *Implemenation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {

	return &desc.LoginResponse{}, nil
}

/*GetRefreshToken(context.Context, *GetRefreshTokenRequest) (*GetRefreshTokenResponse, error)
GetAccessToken(context.Context, *GetAccessTokenRequest) (*GetAccessTokenResponse, error)*/


