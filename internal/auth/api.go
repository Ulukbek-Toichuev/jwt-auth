package auth

import (
	"database/sql"
	"fmt"
	"jwt-auth/config"
	"jwt-auth/pkg"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	authService AuthService
	secretKey   string
	tokenExpiry time.Duration
}

func NewAuthHandler(db *sql.DB, cg *config.Config) *AuthHandler {
	as := NewAuthService(db)
	return &AuthHandler{
		authService: *as,
		secretKey:   cg.GetSecretKey(),
		tokenExpiry: time.Duration(cg.GetTokenExpiry()) * time.Minute,
	}
}

func (ah *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var payload *UserSignInRequest

	payload, err := pkg.ParsePayloadWithValidator[UserSignInRequest](w, r)
	if err != nil {
		pkg.WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := ah.authService.getUserByEmail(payload.Email)
	if err != nil {
		pkg.WriteResponse(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(payload.Password))
	if err != nil {
		pkg.WriteResponse(w, http.StatusInternalServerError, "incorrect password, please check your credential")
		return
	}

	generatedToken, err := pkg.GenerateToken(&res, ah.secretKey, ah.tokenExpiry)
	if err != nil {
		pkg.WriteResponse(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}
	pkg.WriteResponse(w, http.StatusOK, generatedToken)
}

func (ah *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var payload *UserSignUpRequest

	payload, err := pkg.ParsePayloadWithValidator[UserSignUpRequest](w, r)
	if err != nil {
		pkg.WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := ah.authService.createUser(*payload)
	if err != nil {
		pkg.WriteResponse(w, http.StatusInternalServerError, err.Error())
		return

	}
	pkg.WriteResponse(w, http.StatusCreated, fmt.Sprintf("successfully create user with id %d", id))
}
