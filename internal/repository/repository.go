package repository

import (
	"context"

	"github.com/Oleg-Pro/auth/internal/model"
)

// UserRepository interface for User Repository
type UserRepository interface {
	Create(ctx context.Context, info *model.UserInfo) (int64, error)
	Update(ctx context.Context, id int64, name *string, email *string, role model.Role) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Delete(ctx context.Context, id int64) (int64, error)
}