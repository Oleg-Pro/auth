package authentication

import (
	"time"
	"github.com/pkg/errors"		
	"github.com/Oleg-Pro/auth/internal/model"
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
}

func (s *srv)Login(info model.LoginParams) (refereshToken string, err error) {

// Лезем в базу или кэш за данными пользователя
// user = getByEmail(info.Email)


// Сверяем хэши пароля


refreshToken, err := s.userTokenService.GenerateToken(&model.UserTokenParams{
	Username: info.Email,
	Role: "admin",
},
	[]byte(refreshTokenSecretKey),
	refreshTokenExpiration,
)
if err != nil {
	return "", errors.New("failed to generate token")
}

return refreshToken, nil		
}

func New(userTokenService service.UserTokenService) *srv {
	return &srv{userTokenService: userTokenService}
}

