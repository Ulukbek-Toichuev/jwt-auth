package entity

import (
	"database/sql"
	"time"
)

type Role string

const (
	ADMIN Role = "ADMIN"
	USER  Role = "USER"
)

type UserEntity struct {
	UserId      int
	Username    string
	Role        Role
	Email       string
	Password    string
	CreatedDate time.Time
	DeletedDate sql.NullTime
}

func NewUserEntity() *UserEntity {
	return &UserEntity{}
}

func (u *UserEntity) GetId() int {
	return u.UserId
}

func (u *UserEntity) GetEmail() string {
	return u.Email
}

func (u *UserEntity) GetUserName() string {
	return u.Username
}

func (u *UserEntity) GetRole() Role {
	return u.Role
}

func (u *UserEntity) GetPassword() string {
	return u.Password
}
