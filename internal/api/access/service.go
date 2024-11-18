package access

import (
	"github.com/Oleg-Pro/auth/internal/service"
	desc "github.com/Oleg-Pro/auth/pkg/access_v1"
)

// Implemenation AuthenticationServicer implementation
type Implemenation struct {
	desc.UnimplementedAccessV1Server
	accessService service.AccessService

	// authenticationService service.AuthenticationService
}

// NewImplementation create Access Api implementation
func NewImplementation(accessService service.AccessService) *Implemenation {
	return &Implemenation{accessService: accessService}
}
