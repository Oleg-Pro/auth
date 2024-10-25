package user

import (
	"github.com/Oleg-Pro/auth/internal/repository"
	"github.com/Oleg-Pro/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
}

// New create UserService
func New(userRepository repository.UserRepository) service.UserService {
	return &serv{userRepository: userRepository}
}
