package access

import (
	"context"
	"log"

	"github.com/Oleg-Pro/auth/internal/config"
	"github.com/Oleg-Pro/auth/internal/model"
	"github.com/Oleg-Pro/auth/internal/service"
)

var accessibleRoles map[string]string

type srv struct {
	userTokenService service.UserTokenService
	authConfig       config.AuthConfig
}

func (s *srv) Allow(ctx context.Context, endpointAddress string, accessToken string) bool {

	authConfig := s.authConfig
	log.Printf("authConfig %#v\n", authConfig)
	accessTokenKeySecret := s.authConfig.AccessTokenSecretKey()
	log.Printf("accessTokenKeySecret %s\n", accessTokenKeySecret)

	log.Printf("Allow Access Token: %s\n", accessToken)

	claims, err := s.userTokenService.VerifyToken(accessToken, []byte(s.authConfig.AccessTokenSecretKey()))	
	if err != nil {
		log.Printf("Verify token error: %s\n", err.Error())
		return false
	}

	accessibleMap, err := s.accessibleRolesMap(ctx)
	if err != nil {
		return false
	}

	role, ok := accessibleMap[endpointAddress]
	if !ok {
		log.Println("No in map - verified!")
		return true
	}

	if role == claims.Role {
		log.Println("Correct role - verifed!")		
		return true
	}

	return false
}

func (s *srv) accessibleRolesMap(_ context.Context) (map[string]string, error) {
	if accessibleRoles == nil {
		accessibleRoles = make(map[string]string)

		accessibleRoles["/user_v1.UserV1/Get"] = string(model.RoleADMIN)
	}

	return accessibleRoles, nil
}

// New AuthenticationService constructor
func New(userTokenService service.UserTokenService,
	authConfig config.AuthConfig,
) *srv {
	return &srv{userTokenService: userTokenService, authConfig: authConfig}
}
