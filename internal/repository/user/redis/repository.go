package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/Oleg-Pro/auth/internal/client/cache"
	"github.com/Oleg-Pro/auth/internal/model"
	"github.com/Oleg-Pro/auth/internal/repository"
	"github.com/Oleg-Pro/auth/internal/repository/user/redis/converter"
	modelRepo "github.com/Oleg-Pro/auth/internal/repository/user/redis/model"
	redigo "github.com/gomodule/redigo/redis"
)

type repo struct {
	cl cache.RedisClient
}

// NewRepository UserCacheRepository constructor
func NewRepository(cl cache.RedisClient) repository.UserCacheRepository {
	return &repo{cl: cl}
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

	idStr := strconv.FormatInt(id, 10)
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

	idStr := strconv.FormatInt(id, 10)
	err := r.cl.HashSet(ctx, idStr, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	idStr := strconv.FormatInt(id, 10)
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
	idStr := strconv.FormatInt(id, 10)
	err := r.cl.Expire(ctx, idStr, 0)
	var numberOfRows int64
	if err != nil {
		numberOfRows = 1
	}

	return numberOfRows, nil
}
