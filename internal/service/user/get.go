package user

import (
	"context"
	"log"

	"github.com/Oleg-Pro/auth/internal/model"
	"github.com/Oleg-Pro/auth/internal/repository"
)

func (s *serv) Get(ctx context.Context, filter repository.UserFilter) (*model.User, error) {
	user, err := s.userCacheRepository.Get(ctx, filter)
	if err == nil {
		return user, err
	}

	user, err = s.userRepository.Get(ctx, filter)

	if err == nil {
		_, errRedis := s.userCacheRepository.Create(ctx, user.ID, &user.Info)
		if errRedis != nil {
			log.Printf("Failed to add user to cache: %#v", errRedis)
		}
	}

	return user, err
}
