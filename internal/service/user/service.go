package user

import (
	"github.com/Oleg-Pro/auth/internal/repository"
	"github.com/Oleg-Pro/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	userCacheRepository repository.UserCacheRepository
}

// New create UserService
func New(userRepository repository.UserRepository, userCacheRepository repository.UserCacheRepository) service.UserService {
	return &serv{userRepository: userRepository, userCacheRepository: userCacheRepository}
}
