package model

import (
	"database/sql"
	"time"
)

// Role user role
type Role int32

/*const (
	Role_USER  Role = 1
	Role_ADMIN Role = 2
)*/

// UserInfo info about user
type UserInfo struct {
	Name        string
	Email       string
	PaswordHash string
	Role        Role
}

// User entity
type User struct {
	ID        int64
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
