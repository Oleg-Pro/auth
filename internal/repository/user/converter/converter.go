package converter

import (
	"github.com/Oleg-Pro/auth/internal/model"
	modelRepo "github.com/Oleg-Pro/auth/internal/repository/user/model"
)

// ToUserFromRepo converts repo User to model User
func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Info:      ToUserInfoFromRepo(user.Info),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUserInfoFromRepo converts repo UserInfo to model UserInfo
func ToUserInfoFromRepo(user modelRepo.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name:        user.Name,
		Email:       user.Email,
		PaswordHash: user.PaswordHash,
		Role:        model.Role(user.Role),
	}
}
