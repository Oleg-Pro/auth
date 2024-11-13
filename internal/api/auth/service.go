package auth

import (
	"github.com/Oleg-Pro/auth/internal/service"
	desc "github.com/Oleg-Pro/auth/pkg/auth_v1"
)



type Implemenation struct {
	desc.UnimplementedAuthV1Server
	authenticationService service.AuthenticationService
}

func NewImplementation(authenticationService service.AuthenticationService) *Implemenation {
	return &Implemenation{ authenticationService: authenticationService}
}