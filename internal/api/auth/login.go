package auth

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Oleg-Pro/auth/internal/model"
	desc "github.com/Oleg-Pro/auth/pkg/auth_v1"
)

// Login usr login
func (i *Implemenation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {

	refreshToken, err := i.authenticationService.Login(ctx, model.LoginParams{Email: req.GetUsername(), Password: req.GetPassword()})
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &desc.LoginResponse{RefreshToken: refreshToken}, nil
}

/*GetRefreshToken(context.Context, *GetRefreshTokenRequest) (*GetRefreshTokenResponse, error)
GetAccessToken(context.Context, *GetAccessTokenRequest) (*GetAccessTokenResponse, error)*/
