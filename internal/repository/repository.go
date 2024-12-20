package repository

import (
	"context"

	"github.com/Oleg-Pro/auth/internal/model"
)

// UserRepository User Repository interface
type UserRepository interface {
	Create(ctx context.Context, info *model.UserInfo) (int64, error)
	Update(ctx context.Context, id int64, info *model.UserUpdateInfo) (int64, error)
	Get(ctx context.Context, filter UserFilter) (*model.User, error)
	Delete(ctx context.Context, id int64) (int64, error)
}

// UserCacheRepository UserCacheRepository
type UserCacheRepository interface {
	Create(ctx context.Context, id int64, info *model.UserInfo) (int64, error)
	Get(ctx context.Context, filter UserFilter) (*model.User, error)
	Update(ctx context.Context, id int64, info *model.UserUpdateInfo) (int64, error)
	Delete(ctx context.Context, id int64) (int64, error)
}

// UserFilter filter to search users
type UserFilter struct {
	ID    *int64
	Email *string
}
