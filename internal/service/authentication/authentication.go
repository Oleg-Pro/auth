package authentication

import (
	"context"
	"log"

	"github.com/Oleg-Pro/auth/internal/config"
	"github.com/Oleg-Pro/auth/internal/model"
	"github.com/Oleg-Pro/auth/internal/repository"
	"github.com/Oleg-Pro/auth/internal/service"
)

// Use config
/*const (
	//	authPrefix = "Bearer "

	refreshTokenSecretKey = "W4/X+LLjehdxptt4YgGFCvMpq5ewptpZZYRHY6A72g0="
	accessTokenSecretKey  = "VqvguGiffXILza1f44TWXowDT4zwf03dtXmqWW4SYyE="

	refreshTokenExpiration = 60 * time.Minute
	accessTokenExpiration  = 5 * time.Minute
)*/

type srv struct {
	userTokenService    service.UserTokenService
	userRepository      repository.UserRepository
	passwordVerificator service.PasswordVerificator
	authConfig          config.AuthConfig
}

func (s *srv) Login(ctx context.Context, info model.LoginParams) (refereshToken string, err error) {
	userInfo, err := s.userRepository.GetByEmail(ctx, info.Email)
	if err != nil {
		return "", model.ErrorFailToGenerateToken
	}

	if !s.passwordVerificator.VerifyPassword(userInfo.Info.PaswordHash, info.Password) {
		log.Printf("Password does not correspond to hash")
		return "", model.ErrorFailToGenerateToken
	}

	refreshToken, err := s.userTokenService.GenerateToken(&model.UserTokenParams{
		Username: userInfo.Info.Email,
		Role:     string(userInfo.Info.Role),
	},
		[]byte(s.authConfig.RefreshTokenSecretKey()),
		s.authConfig.RefreshTokenExpiration(),
	)
	if err != nil {
		return "", model.ErrorFailToGenerateToken
	}

	return refreshToken, nil
}

// New AuthenticationService constructor
func New(
	userTokenService service.UserTokenService,
	userRepository repository.UserRepository,
	passwordVerificator service.PasswordVerificator,
	authConfig config.AuthConfig,
) *srv {
	return &srv{
		userTokenService:    userTokenService,
		userRepository:      userRepository,
		passwordVerificator: passwordVerificator,
		authConfig:          authConfig,
	}
}
