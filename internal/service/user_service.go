package service

import (
	"database/sql"
	store "jwt-auth/internal/db"
	entity "jwt-auth/internal/entity"
	model "jwt-auth/internal/model"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User interface {
	CreateUser(user entity.UserEntity) (int, error)
	GetUserByEmail(email string) (entity.UserEntity, error)
	GetUsers() ([]entity.UserEntity, error)
	ChangeUsersRole(role, email string) (int, error)
	DeleteUser(email string) (int, error)
}

type UserService struct {
	userStore User
}

func NewUserService(db *sql.DB) *UserService {
	userStore := store.NewUserStore(db)
	return &UserService{userStore}
}

func (us *UserService) CreateUser(user model.UserSignUpRequest) (int, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	prepareUser := entity.UserEntity{
		Username:    user.Username,
		Role:        "USER",
		Email:       user.Email,
		Password:    string(hashedPassword),
		CreatedDate: time.Now(),
	}
	id, err := us.userStore.CreateUser(prepareUser)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (us *UserService) GetUserByEmail(email string) (model.UserResponse, error) {
	user, err := us.userStore.GetUserByEmail(email)
	if err != nil {
		return model.UserResponse{}, err
	}
	return *model.NewUserResponse(
		user.UserId,
		user.Username,
		user.Role,
		user.Email,
		user.CreatedDate,
	), nil
}

func (us *UserService) GetUserByEmailWithPasswd(email string) (model.UserAuthResponse, error) {
	user, err := us.userStore.GetUserByEmail(email)
	if err != nil {
		return model.UserAuthResponse{}, err
	}
	return *model.NewUserAuthResponse(
		user.UserId,
		user.Username,
		user.Role,
		user.Email,
		user.Password,
		user.CreatedDate,
	), nil
}

func (us *UserService) GetUsers() ([]model.UserResponse, error) {
	users, err := us.userStore.GetUsers()
	if err != nil {
		return []model.UserResponse{}, err
	}

	result := make([]model.UserResponse, len(users))
	for i, v := range users {
		result[i] = *model.NewUserResponse(
			v.UserId,
			v.Username,
			v.Role,
			v.Email,
			v.CreatedDate,
		)
	}
	return result, nil
}

func (us *UserService) ChangeUsersRole(user model.UserChangeRoleRequest) (int, error) {
	res, err := us.userStore.ChangeUsersRole(user.Role, user.Email)
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (us *UserService) DeleteUser(user model.UserDeleteRequest) (int, error) {
	res, err := us.userStore.DeleteUser(user.Email)
	if err != nil {
		return 0, err
	}

	return res, nil
}
