package user

import (
	"context"
	"log"
	"github.com/Oleg-Pro/auth/internal/model"	
	desc "github.com/Oleg-Pro/auth/pkg/user_v1"	
	empty "github.com/golang/protobuf/ptypes/empty"	
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	var name, email *string
	if req.GetName() != nil {
		name = &req.GetName().Value
	}

	if req.GetEmail() != nil {
		email = &req.GetEmail().Value
	}

	_, err := i.userService.Update(ctx, req.GetId(), name, email, model.Role(req.GetRole()))
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return nil, err
	}

	return &empty.Empty{}, nil
}
