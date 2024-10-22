package user

import (
	"context"

	"github.com/Oleg-Pro/auth/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
	return s.userRepository.Get(ctx, id)
}
