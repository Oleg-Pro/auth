package access

import (
	"context"
	"log"
	"slices"

	"github.com/Oleg-Pro/auth/internal/config"
	"github.com/Oleg-Pro/auth/internal/model"
	"github.com/Oleg-Pro/auth/internal/service"
)

var accessibleRoles map[string][]string

type srv struct {
	userTokenService service.UserTokenService
	authConfig       config.AuthConfig
}

func (s *srv) Allow(ctx context.Context, endpointAddress string, accessToken string) bool {

	authConfig := s.authConfig
	log.Printf("authConfig %#v\n", authConfig)

	claims, err := s.userTokenService.VerifyToken(accessToken, []byte(s.authConfig.AccessTokenSecretKey()))
	if err != nil {
		log.Printf("Verify token error: %s\n", err.Error())
		return false
	}

	// Admin has access eveywhere
	if claims.Role == string(model.RoleADMIN) {
		log.Println("ADMIN role - verified!")
		return true
	}

	accessibleMap, err := s.accessibleRolesMap(ctx)
	if err != nil {
		return false
	}

	roles, ok := accessibleMap[endpointAddress]
	if !ok {
		log.Println("No in map - verified!")
		return false
	}

	if slices.Contains(roles, claims.Role) {
		log.Println("Role has access to endpoint!")
		return true
	}

	return false
}

func (s *srv) accessibleRolesMap(_ context.Context) (map[string][]string, error) {
	if accessibleRoles == nil {
		accessibleRoles = make(map[string][]string)

		accessibleRoles["chat_v1.ChatV1/SendMessage"] = []string{string(model.RoleUSER)}
	}

	return accessibleRoles, nil
}

// New AuthenticationService constructor
func New(userTokenService service.UserTokenService,
	authConfig config.AuthConfig,
) *srv {
	return &srv{userTokenService: userTokenService, authConfig: authConfig}
}
