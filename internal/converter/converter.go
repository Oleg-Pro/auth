package converter

import (
	"github.com/Oleg-Pro/auth/internal/model"
	desc "github.com/Oleg-Pro/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToUserGetResponseFromModelUser converts *model.User to *desc.GetResponse
func ToUserGetResponseFromModelUser(user *model.User) *desc.GetResponse {
	var updatedAtTime *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAtTime = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.GetResponse{
		Id:        user.ID,
		Name:      user.Info.Name,
		Email:     user.Info.Email,
		Role:      desc.Role(user.Info.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAtTime,
	}
}
