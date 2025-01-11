package model

import (
	"time"
)

type UserResponse struct {
	UserId      int
	Username    string
	Role        string
	Email       string
	CreatedDate time.Time
}

func NewUserResponse(userId int, username, role, email string, createdDate time.Time) *UserResponse {
	return &UserResponse{
		userId,
		username,
		role,
		email,
		createdDate,
	}
}
