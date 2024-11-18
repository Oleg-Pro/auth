package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Oleg-Pro/auth/internal/client/cache"
	"github.com/Oleg-Pro/auth/internal/model"
	"github.com/Oleg-Pro/auth/internal/repository"
	"github.com/Oleg-Pro/auth/internal/repository/user/redis/converter"
	modelRepo "github.com/Oleg-Pro/auth/internal/repository/user/redis/model"
	redigo "github.com/gomodule/redigo/redis"
)

const keyPrefix = "USER_"

type repo struct {
	cl cache.RedisClient
}

// NewRepository UserCacheRepository constructor
func NewRepository(cl cache.RedisClient) repository.UserCacheRepository {
	return &repo{cl: cl}
}

func userKey(id int64) string {
	return fmt.Sprintf("%s%d", keyPrefix, id)
}

func (r *repo) Create(ctx context.Context, id int64, info *model.UserInfo) (int64, error) {
	user := modelRepo.User{
		ID:          id,
		Name:        info.Name,
		Email:       info.Email,
		PaswordHash: info.PaswordHash,
		Role:        int32(info.Role),
		CreatedAtNs: time.Now().UnixNano(),
	}

	idStr := userKey(id)
	err := r.cl.HashSet(ctx, idStr, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Update(ctx context.Context, id int64, info *model.UserUpdateInfo) (int64, error) {

	user := modelRepo.UserUpdateInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  int64(info.Role),
	}

	idStr := userKey(id)
	err := r.cl.HashSet(ctx, idStr, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, filter repository.UserFilter) (*model.User, error) {
	idStr := userKey(*filter.ID)
	values, err := r.cl.HGetAll(ctx, idStr)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, model.ErrorNoteNotFound
	}

	var user modelRepo.User
	err = redigo.ScanStruct(values, &user)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

func (r *repo) Delete(ctx context.Context, id int64) (int64, error) {
	idStr := userKey(id)
	err := r.cl.Expire(ctx, idStr, 0)
	var numberOfRows int64
	if err != nil {
		numberOfRows = 1
	}

	return numberOfRows, nil
}
