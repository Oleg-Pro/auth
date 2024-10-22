package user

import (
	"context"

	"github.com/Oleg-Pro/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, info *model.UserInfo) (int64, error) {
	return s.userRepository.Create(ctx, info)
}
