package api

import (
	"database/sql"
	"fmt"
	"jwt-auth/config"
	"jwt-auth/internal/auth"
	"jwt-auth/internal/model"
	"jwt-auth/internal/service"
	"jwt-auth/pkg"
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

func (ah *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var payload *model.UserSignInRequest

	payload, err := pkg.ParsePayloadWithValidator[model.UserSignInRequest](w, r)
	if err != nil {
		pkg.WriteResponseWithMssg(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := ah.userService.GetUserByEmailWithPasswd(payload.Email)
	if err != nil {
		pkg.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(payload.Password))
	if err != nil {
		pkg.WriteResponseWithMssg(w, http.StatusInternalServerError, "incorrect password, please check your credential")
		return
	}

	generatedToken, err := auth.GenerateToken(&res, ah.secretKey, ah.tokenExpiry)
	if err != nil {
		pkg.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}
	pkg.WriteResponseWithMssg(w, http.StatusOK, generatedToken)
}

func (ah *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var payload *model.UserSignUpRequest

	payload, err := pkg.ParsePayloadWithValidator[model.UserSignUpRequest](w, r)
	if err != nil {
		pkg.WriteResponseWithMssg(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := ah.userService.CreateUser(*payload)
	if err != nil {
		pkg.WriteResponseWithMssg(w, http.StatusInternalServerError, err.Error())
		return

	}
	pkg.WriteResponseWithMssg(w, http.StatusCreated, fmt.Sprintf("successfully create user with id %d", id))
}
