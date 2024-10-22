package user

import (
	"context"

	"github.com/Oleg-Pro/auth/internal/converter"
	desc "github.com/Oleg-Pro/auth/pkg/user_v1"
)

// Get implementation of Create User Api Method
func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := i.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return converter.ToUserGetResponseFromModelUser(user), nil
}
