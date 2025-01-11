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
	DeleteUser(id int)
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

func (us *UserService) GetUserByEmail(email string) (entity.UserEntity, error) {
	user, err := us.userStore.GetUserByEmail(email)
	if err != nil {
		return entity.UserEntity{}, err
	}
	return user, nil
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
