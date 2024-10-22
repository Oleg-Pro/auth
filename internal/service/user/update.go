package user

import (
	"context"

	"github.com/Oleg-Pro/auth/internal/model"
)

func (s *serv) Update(ctx context.Context, id int64, name *string, email *string, role model.Role) (int64, error) {
	return s.userRepository.Update(ctx, id, name, email, role)
}
