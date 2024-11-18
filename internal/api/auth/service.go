package auth

import (
	"github.com/Oleg-Pro/auth/internal/service"
	desc "github.com/Oleg-Pro/auth/pkg/auth_v1"
)

// Implemenation AuthenticationServicer implementation
type Implemenation struct {
	desc.UnimplementedAuthV1Server
	authenticationService service.AuthenticationService
}

// NewImplementation create Auth Api implementation
func NewImplementation(authenticationService service.AuthenticationService) *Implemenation {
	return &Implemenation{authenticationService: authenticationService}
}
