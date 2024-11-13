package model

import (
	"database/sql"
	"time"
)

// Role user role
type Role int32

const (
	// RoleUNKNOWN Unknown role
	RoleUNKNOWN Role = 0

	// RoleUSER user role
	RoleUSER Role = 1

	// RoleADMIN admin role
	RoleADMIN Role = 2
)

// UserInfo info about user
type UserInfo struct {
	Name        string
	Email       string
	PaswordHash string
	Role        Role
}

// UserUpdateInfo data to update for user
type UserUpdateInfo struct {
	Name  *string
	Email *string
	Role  Role
}

// User entity
type User struct {
	ID        int64
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type  UserTokenParams struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type LoginParams struct {
	Email string
	Password string
}
