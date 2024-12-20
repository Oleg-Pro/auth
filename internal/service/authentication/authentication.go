package authentication

import (
	"context"
	"log"

	"github.com/Oleg-Pro/auth/internal/config"
	"github.com/Oleg-Pro/auth/internal/model"
	"github.com/Oleg-Pro/auth/internal/repository"
	"github.com/Oleg-Pro/auth/internal/service"
)

type srv struct {
	userTokenService    service.UserTokenService
	userRepository      repository.UserRepository
	passwordVerificator service.PasswordVerificator
	authConfig          config.AuthConfig
}

func (s *srv) Login(ctx context.Context, info model.LoginParams) (refereshToken string, err error) {
	email := info.Email
	userInfo, err := s.userRepository.Get(ctx, repository.UserFilter{Email: &email})
	if err != nil {
		return "", model.ErrorFailToGenerateToken
	}

	if !s.passwordVerificator.VerifyPassword(userInfo.Info.PaswordHash, info.Password) {
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

func (s *srv) GetRefreshToken(_ context.Context, oldRefreshToken string) (string, error) {
	claims, err := s.userTokenService.VerifyToken(oldRefreshToken, []byte(s.authConfig.RefreshTokenSecretKey()))
	if err != nil {
		return "", model.ErrorInvalidRefereshToken
	}

	log.Printf("GetRefreshToken Claim: %#v\n", claims)

	refreshToken, err := s.userTokenService.GenerateToken(&model.UserTokenParams{
		Username: claims.Username,
		Role:     string(claims.Role),
	},
		[]byte(s.authConfig.RefreshTokenSecretKey()),
		s.authConfig.RefreshTokenExpiration(),
	)

	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (s *srv) GetAccessToken(_ context.Context, refreshToken string) (string, error) {
	claims, err := s.userTokenService.VerifyToken(refreshToken, []byte(s.authConfig.RefreshTokenSecretKey()))
	if err != nil {
		return "", model.ErrorInvalidRefereshToken
	}

	accessToken, err := s.userTokenService.GenerateToken(&model.UserTokenParams{
		Username: claims.Username,
		Role:     string(claims.Role),
	},
		[]byte(s.authConfig.AccessTokenSecretKey()),
		s.authConfig.AccessTokenExpiration(),
	)

	if err != nil {
		return "", err
	}

	return accessToken, nil
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
