package model

import (
	"database/sql"
	"time"
)

// Role user role
type Role int32

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
