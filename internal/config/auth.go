package config

import (
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	refreshTokenSecretKeyEnvName = "REFRESH_TOKEN_SECRET_KEY" // #nosec G101
	accessTokenSecretKeyEnvName  = "ACCESS_TOKEN_SECRET_KEY"

	refreshTokenExpiration = "REFRESH_TOKEN_EXPIRATION" // #nosec G101
	accessTokenExpiration  = "ACCESS_TOKEN_EXPIRATION"
)

// AuthConfig to get server address
type AuthConfig interface {
	RefreshTokenSecretKey() string
	AccessTokenSecretKey() string
	RefreshTokenExpiration() time.Duration
	AccessTokenExpiration() time.Duration
}

type authConfig struct {
	refreshTokenSecretKey  string
	accessTokenSecretKey   string
	refreshTokenExpiration time.Duration
	accessTokenExpiration  time.Duration
}

// NewAuthConfig will get authentication paramteres
func NewAuthConfig() (AuthConfig, error) {
	refreshTokenSecretKey := os.Getenv(refreshTokenSecretKeyEnvName)
	if len(refreshTokenSecretKey) == 0 {
		return nil, errors.New("refreshTokenSecretKey not found")
	}

	accessTokenSecretKey := os.Getenv(accessTokenSecretKeyEnvName)
	if len(accessTokenSecretKey) == 0 {
		return nil, errors.New("accessTokenSecretKey not found")
	}

	refreshTokenExpirationMinutes, err := strconv.Atoi(os.Getenv(refreshTokenExpiration))
	if refreshTokenExpirationMinutes == 0 || err != nil {
		return nil, errors.New("refreshTokenExpiration not found")
	}

	accessTokenExpirationMinutes, err := strconv.Atoi(os.Getenv(accessTokenExpiration))
	if accessTokenExpirationMinutes == 0 || err != nil {
		return nil, errors.New("refreshTokenExpiration not found")
	}

	return &authConfig{
		refreshTokenSecretKey:  refreshTokenSecretKey,
		accessTokenSecretKey:   accessTokenSecretKey,
		refreshTokenExpiration: time.Duration(refreshTokenExpirationMinutes) * time.Minute,
		accessTokenExpiration:  time.Duration(accessTokenExpirationMinutes) * time.Minute,
	}, nil
}

func (a *authConfig) RefreshTokenSecretKey() string {
	return a.refreshTokenSecretKey
}

func (a *authConfig) AccessTokenSecretKey() string {
	return a.accessTokenSecretKey
}

func (a *authConfig) RefreshTokenExpiration() time.Duration {
	return a.refreshTokenExpiration
}

func (a *authConfig) AccessTokenExpiration() time.Duration {
	return a.accessTokenExpiration
}
