package api

import (
	"database/sql"
	"fmt"
	"jwt-auth/internal/service"
	"jwt-auth/pkg"
	"net/http"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(db *sql.DB) *UserHandler {
	us := service.NewUserService(db)
	return &UserHandler{*us}
}

func (uh *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	result := r.Context().Value(pkg.ResultCtxKey).(map[string]interface{})

	currUserEmail := ""
	if value, ok := result["email"]; !ok {
		return
	} else {
		currUserEmail = value.(string)
	}

	res, err := uh.userService.GetUserByEmail(currUserEmail)
	if err != nil {
		pkg.WriteResponse(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	if res.Role != "ADMIN" {
		pkg.WriteResponse(w, http.StatusForbidden, "only for user with admin role")
		return
	}

	users, err := uh.userService.GetUsers()
	if err != nil {
		pkg.WriteResponse(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	pkg.WriteResponseWithoutMssg(w, http.StatusOK, users)
}

func (uh *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {

}

func (uh *UserHandler) GetUsersOwnDetail(w http.ResponseWriter, r *http.Request) {

}

func (uh *UserHandler) ChangeUsersRole(w http.ResponseWriter, r *http.Request) {

}

func (uh *UserHandler) DeleteUserByEmail(w http.ResponseWriter, r *http.Request) {

}
