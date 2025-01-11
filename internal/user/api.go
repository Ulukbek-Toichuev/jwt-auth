package user

import (
	"jwt-auth/pkg"
	"net/http"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (uh *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	result := r.Context().Value(pkg.ResultCtxKey).(map[string]interface{})

	currRole := ""

}

func (uh *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {

}

func (uh *UserHandler) GetUsersOwnDetail(w http.ResponseWriter, r *http.Request) {

}

func (uh *UserHandler) ChangeUsersRole(w http.ResponseWriter, r *http.Request) {

}

func (uh *UserHandler) DeleteUserByEmail(w http.ResponseWriter, r *http.Request) {

}
