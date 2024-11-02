package model

// User model for redis repository
type User struct {
	ID          int64  `redis:"id"`
	Name        string `redis:"name"`
	Email       string `redis:"email"`
	PaswordHash string `redis:"paswordhash"`
	Role        int32  `redis:"role"`
	CreatedAtNs int64  `redis:"created_at"`
	UpdatedAtNs *int64 `redis:"updated_at"`
}

// UserUpdateInfo info to update for redis repository
type UserUpdateInfo struct {
	Name  *string `redis:"name"`
	Email *string `redis:"email"`
	Role  int64   `redis:"role"`
}
