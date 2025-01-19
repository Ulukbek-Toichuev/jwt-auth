package api

import (
	"database/sql"
	"fmt"
	"jwt-auth/config"
	"jwt-auth/internal/model"
	"jwt-auth/internal/service"
	"jwt-auth/internal/util"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userService service.UserService
	secretKey   string
	tokenExpiry time.Duration
}

func NewAuthHandler(db *sql.DB, cg *config.Config) *AuthHandler {
	as := service.NewUserService(db)
	return &AuthHandler{
		userService: *as,
		secretKey:   cg.GetSecretKey(),
		tokenExpiry: time.Duration(cg.GetTokenExpiry()) * time.Minute,
	}
}

// Хендлер для аутентификации пользователя
func (ah *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var payload *model.UserSignInRequest

	//Валидируем тело запроса
	payload, err := util.ParsePayloadWithValidator[model.UserSignInRequest](w, r)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusBadRequest, err.Error())
		return
	}

	//Запрашиваем пользователя по email из тела запроса
	res, err := ah.userService.GetUserByEmailWithPasswd(payload.Email)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	//Сравниваем пароли
	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(payload.Password))
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, "incorrect password, please check your credential")
		return
	}

	//Генерируем JWT токен
	generatedToken, err := util.GenerateToken(&res, ah.secretKey, ah.tokenExpiry)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	util.WriteResponseWithMssg(w, http.StatusOK, generatedToken)
}

// Хендлер для регистрации нового пользователя
func (ah *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var payload *model.UserSignUpRequest

	//Валидируем тело запроса
	payload, err := util.ParsePayloadWithValidator[model.UserSignUpRequest](w, r)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusBadRequest, err.Error())
		return
	}

	//Создаем пользователя
	id, err := ah.userService.CreateUser(*payload)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, err.Error())
		return

	}
	util.WriteResponseWithMssg(w, http.StatusCreated, fmt.Sprintf("successfully create user with id %d", id))
}
