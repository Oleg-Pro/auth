package authentication

import (
	"context"
	"time"
	"log"

	"github.com/Oleg-Pro/auth/internal/model"
	"github.com/Oleg-Pro/auth/internal/repository"
	"github.com/Oleg-Pro/auth/internal/service"
)

// Use config
const (
	//	authPrefix = "Bearer "
	
		refreshTokenSecretKey = "W4/X+LLjehdxptt4YgGFCvMpq5ewptpZZYRHY6A72g0="
		accessTokenSecretKey  = "VqvguGiffXILza1f44TWXowDT4zwf03dtXmqWW4SYyE="
	
		refreshTokenExpiration = 60 * time.Minute
		accessTokenExpiration  = 5 * time.Minute
	)

type srv struct {
	userTokenService service.UserTokenService
	userRepository repository.UserRepository
	passwordVerificator service.PasswordVerificator
}

func (s *srv)Login(ctx context.Context, info model.LoginParams) (refereshToken string, err error) {
	userInfo, err := s.userRepository.GetByEmail(ctx, info.Email)
	if err != nil {
		return "", model.ErrorFailToGenerateToken
	}	

	log.Printf("User Info By Email %#v", userInfo)

	if !s.passwordVerificator.VerifyPassword(userInfo.Info.PaswordHash, info.Password) {
		log.Printf("Password does not correspond to hash")
		return "", model.ErrorFailToGenerateToken
	}

refreshToken, err := s.userTokenService.GenerateToken(&model.UserTokenParams{
	Username: userInfo.Info.Email,
	Role: string(userInfo.Info.Role),
},
	[]byte(refreshTokenSecretKey),
	refreshTokenExpiration,
)
if err != nil {
	return "", model.ErrorFailToGenerateToken
}

return refreshToken, nil		
}

func New(
	userTokenService service.UserTokenService,
	 userRepository repository.UserRepository,
	  passwordVerificator service.PasswordVerificator,
	  ) *srv {
	return &srv{userTokenService: userTokenService, userRepository: userRepository, passwordVerificator: passwordVerificator}
}

