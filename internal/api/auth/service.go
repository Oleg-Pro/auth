package auth

import (
	desc "github.com/Oleg-Pro/auth/pkg/auth_v1"	
)


type Implemenation struct {
	desc.UnimplementedAuthV1Server
}

func NewImplementation() *Implemenation {
	return &Implemenation{}
}