package user

import (
	"context"
	"log"

	"github.com/Oleg-Pro/auth/internal/model"
)

func (s *serv) Update(ctx context.Context, id int64, info *model.UserUpdateInfo) (int64, error) {

	numberOfRows, err := s.userRepository.Update(ctx, id, info)

	if err == nil && numberOfRows != 0 {
		_, errRedis := s.userCacheRepository.Update(ctx, id, info)
		if errRedis != nil {
			log.Printf("Failed to update user in cache: %#v", errRedis)
		}
	}

	return s.userRepository.Update(ctx, id, info)
}
