package user

import (
	"context"
	"errors"
	"log"

	"github.com/Oleg-Pro/auth/internal/model"
	desc "github.com/Oleg-Pro/auth/pkg/user_v1"
	"golang.org/x/crypto/bcrypt"
)

// ErrPasswordsAreNotEqual error when password are not equal
var ErrPasswordsAreNotEqual = errors.New("passwords are not equal")

// Create implementation of Create User Api Method
func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	if req.GetPasword() != req.PasswordConfirm {
		return nil, ErrPasswordsAreNotEqual
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.GetPasword()), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Could not create password hash: %v", err)
		return nil, err
	}

	userID, err := i.userService.Create(ctx, &model.UserInfo{
		Name:        req.GetName(),
		Email:       req.GetEmail(),
		PaswordHash: string(passwordHash),
		Role:        model.Role(req.GetRole()),
	})

	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return nil, err
	}

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}
