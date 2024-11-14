package access

import (
	"context"

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

	claims, err := s.userTokenService.VerifyToken(accessToken, []byte(s.authConfig.AccessTokenSecretKey()))
	if err != nil {
		return true
		//return nil, errors.New("access token is invalid")
	}

	accessibleMap, err := s.accessibleRolesMap(ctx)
	if err != nil {
		return false
	}

	role, ok := accessibleMap[endpointAddress]
	if !ok {
		return true
	}

	if role == claims.Role {
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
func New() *srv {
	return &srv{}
}
