package user

import (
	"context"
	"github.com/Oleg-Pro/auth/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.userCacheRepository.Get(ctx, id)
	if err == nil {
		return user, err
	}

	user, err = s.userRepository.Get(ctx, id)


	if err == nil {
		s.userCacheRepository.Create(ctx, user.ID, &user.Info)
	}

	return user, err
}
