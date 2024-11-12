package token

import (
	"time"
	"github.com/pkg/errors"
	"github.com/dgrijalva/jwt-go"	
	"github.com/Oleg-Pro/auth/internal/model"
)

type serv struct {
}

func (s *serv) Token(info model.UserTokenParams, secretKey []byte, duration time.Duration) (string, error) {

	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		Username: info.Username,
		Role:     info.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)	
}

func New() *serv {
	return &serv{}
}