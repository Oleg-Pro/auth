package password_verificator

import "golang.org/x/crypto/bcrypt"

type  srv struct {
}

func(s *srv) VerifyPassword(hashedPassword string, candidatePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
	return err == nil
}

func New() *srv {
	return &srv{}
}
