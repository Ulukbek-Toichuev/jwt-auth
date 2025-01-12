package model

import "time"

// Авторизация
type UserSignInRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func NewUserSignInRequest() *UserSignInRequest {
	return &UserSignInRequest{}
}

// Регистрация
type UserSignUpRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=12"`
}

func NewUserSignUpRequest() *UserSignUpRequest {
	return &UserSignUpRequest{}
}

type UserAuthResponse struct {
	UserId      int
	Username    string
	Role        string
	Email       string
	Password    string
	CreatedDate time.Time
}

func NewUserAuthResponse(userId int, username, role, email, password string, createdDate time.Time) *UserAuthResponse {
	return &UserAuthResponse{
		userId,
		username,
		role,
		email,
		password,
		createdDate,
	}
}

func (u *UserAuthResponse) GetId() int {
	return u.UserId
}

func (u *UserAuthResponse) GetEmail() string {
	return u.Email
}

func (u *UserAuthResponse) GetUserName() string {
	return u.Username
}

func (u *UserAuthResponse) GetRole() string {
	return u.Role
}

func (u *UserAuthResponse) GetPassword() string {
	return u.Password
}
