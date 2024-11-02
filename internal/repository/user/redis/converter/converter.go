package converter

import (
	"database/sql"
	"time"

	"github.com/Oleg-Pro/auth/internal/model"
	modelRepo "github.com/Oleg-Pro/auth/internal/repository/user/redis/model"
)

// ToUserFromRepo convert from redis repo to user model
func ToUserFromRepo(user *modelRepo.User) *model.User {
	var updatedAt sql.NullTime
	if user.UpdatedAtNs != nil {
		updatedAt = sql.NullTime{
			Time:  time.Unix(0, *user.UpdatedAtNs),
			Valid: true,
		}
	}

	return &model.User{
		ID: user.ID,
		Info: model.UserInfo{
			Name:        user.Name,
			Email:       user.Email,
			PaswordHash: user.PaswordHash,
			Role:        model.Role(user.Role),
		},
		CreatedAt: time.Unix(0, user.CreatedAtNs),
		UpdatedAt: updatedAt,
	}
}
