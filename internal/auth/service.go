package auth

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth interface {
	createUser(user UserEntity) (int, error)
	getUserByEmail(email string) (UserEntity, error)
	getUsers() ([]UserEntity, error)
	deleteUser(id int)
}

type AuthService struct {
	authStore Auth
}

func NewAuthService(db *sql.DB) *AuthService {
	authStore := NewAuthStore(db)
	return &AuthService{authStore}
}

func (as *AuthService) createUser(user UserSignUpRequest) (int, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	prepareUser := UserEntity{
		Username:    user.Username,
		Role:        "USER",
		Email:       user.Email,
		Password:    string(hashedPassword),
		CreatedDate: time.Now(),
	}
	id, err := as.authStore.createUser(prepareUser)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (as *AuthService) getUserByEmail(email string) (UserEntity, error) {
	user, err := as.authStore.getUserByEmail(email)
	if err != nil {
		return UserEntity{}, err
	}
	return user, nil
}
