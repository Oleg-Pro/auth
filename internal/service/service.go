package service

import (
	"context"
	"time"

	"github.com/Oleg-Pro/auth/internal/model"
)

// UserService iterface for User Service
type UserService interface {
	Create(ctx context.Context, info *model.UserInfo) (int64, error)
	Update(ctx context.Context, id int64, info *model.UserUpdateInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Delete(ctx context.Context, id int64) (int64, error)	
}

type UserTokenService interface {
	GenerateToken(info *model.UserTokenParams, secretKey []byte, duration time.Duration) (string, error)
	VerifyToken(tokenStr string, secretKey []byte) (*model.UserClaims, error) 	
}

type AuthenticationService interface {
	Login(ctx context.Context, info model.LoginParams) (refereshToken string, err error)
}

type PasswordVerificator interface {
	VerifyPassword(hashedPassword string, candidatePassword string) bool
}

// ConsumerService interface
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}
