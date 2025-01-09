package auth

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
