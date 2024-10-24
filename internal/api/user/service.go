package user

import (
	"github.com/Oleg-Pro/auth/internal/service"
	desc "github.com/Oleg-Pro/auth/pkg/user_v1"
)

// Implementation implementation of User API
type Implementation struct {
	desc.UnimplementedUserV1Server
	userService service.UserService
}

// NewImplementation create User Api implementation
func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{userService: userService}
}
