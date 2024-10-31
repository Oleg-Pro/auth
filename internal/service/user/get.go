package user

import (
	"context"
	"log"

	"github.com/Oleg-Pro/auth/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
	//Почему здесь паника?
	user, err := s.userCacheRepository.Get(ctx, id)
	if err == nil {
		log.Printf("Data from cache: %#v\n", user)
		return user, err
	}

	user, err = s.userRepository.Get(ctx, id)

	if err != nil {
		s.userCacheRepository.Create(ctx, user.ID, &user.Info)
	}

	return user, err

//	return s.userRepository.Get(ctx, id)	
}
