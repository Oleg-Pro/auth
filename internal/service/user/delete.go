package user

import (
	"context"
	"log"
)

func (s *serv) Delete(ctx context.Context, id int64) (int64, error) {
	numberOfRows, err := s.userRepository.Delete(ctx, id)
	_, errRedis := s.userCacheRepository.Delete(ctx, id)
	if errRedis != nil {
		log.Printf("Failed to delete user in cache: %#v", errRedis)
	}

	return numberOfRows, err
}
