package user

import (
	"context"
	"log"

	desc "github.com/Oleg-Pro/auth/pkg/user_v1"
	empty "github.com/golang/protobuf/ptypes/empty"
)

// Delete implementation of Create User Api Method
func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	log.Printf("Deleting User req=%v", req)
	_, err := i.userService.Delete(ctx, req.GetId())
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		return nil, err
	}

	return &empty.Empty{}, nil
}
