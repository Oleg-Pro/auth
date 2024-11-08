package user

import (
	"github.com/Oleg-Pro/auth/internal/service"
	userSaverProducer "github.com/Oleg-Pro/auth/internal/service/producer/user_saver"
	desc "github.com/Oleg-Pro/auth/pkg/user_v1"
)

// Implementation implementation of User API
type Implementation struct {
	desc.UnimplementedUserV1Server
	userService       service.UserService
	userSaverProducer userSaverProducer.UserSaverProducer
}

// NewImplementation create User Api implementation
func NewImplementation(userService service.UserService, userSaverProducer userSaverProducer.UserSaverProducer) *Implementation {
	return &Implementation{userService: userService, userSaverProducer: userSaverProducer}
}
