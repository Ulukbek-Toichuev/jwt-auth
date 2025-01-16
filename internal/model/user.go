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
		UserId:      userId,
		Username:    username,
		Role:        role,
		Email:       email,
		CreatedDate: createdDate,
	}
}

type UserChangeRoleRequest struct {
	Email string `json:"email" validate:"required"`
	Role  string `json:"role" validate:"required"`
}

func NewUserChangeRoleRequest() *UserChangeRoleRequest {
	return &UserChangeRoleRequest{}
}

type UserDeleteRequest struct {
	Email string `json:"email" validate:"required"`
}

func NewUserDeleteRequest() *UserDeleteRequest {
	return &UserDeleteRequest{}
}
