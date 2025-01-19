package api

import (
	"database/sql"
	"fmt"
	"jwt-auth/internal/middleware"
	"jwt-auth/internal/model"
	"jwt-auth/internal/service"
	"jwt-auth/internal/util"
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
	isHavePermission := uh.userService.VerifyUserFromCTX(w, r)
	if !isHavePermission {
		util.WriteResponseWithMssg(w, http.StatusForbidden, "the user does not meet role requirements")
		return
	}
	users, err := uh.userService.GetUsers()
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	util.WriteResponse(w, http.StatusOK, users)
}

func (uh *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	isHavePermission := uh.userService.VerifyUserFromCTX(w, r)
	if !isHavePermission {
		util.WriteResponseWithMssg(w, http.StatusForbidden, "the user does not meet role requirements")
		return
	}

	email := r.PathValue("email")
	user, err := uh.userService.GetUserByEmail(email)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	util.WriteResponse(w, http.StatusOK, user)
}

func (uh *UserHandler) GetUsersOwnDetail(w http.ResponseWriter, r *http.Request) {
	result := r.Context().Value(middleware.ResultCtxKey).(map[string]interface{})

	currUserEmail := ""
	if value, ok := result["email"]; !ok {
		return
	} else {
		currUserEmail = value.(string)
	}

	res, err := uh.userService.GetUserByEmail(currUserEmail)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}
	util.WriteResponse(w, http.StatusOK, res)
}

func (uh *UserHandler) ChangeUsersRole(w http.ResponseWriter, r *http.Request) {
	isHavePermission := uh.userService.VerifyUserFromCTX(w, r)
	if !isHavePermission {
		util.WriteResponseWithMssg(w, http.StatusForbidden, "the user does not meet role requirements")
		return
	}

	var payload *model.UserChangeRoleRequest

	payload, err := util.ParsePayloadWithValidator[model.UserChangeRoleRequest](w, r)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedRowsCount, err := uh.userService.ChangeUsersRole(*payload)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	util.WriteResponseWithMssg(w, http.StatusOK, fmt.Sprintf("updated rows count in db: %d", updatedRowsCount))
}

func (uh *UserHandler) DeleteUserByEmail(w http.ResponseWriter, r *http.Request) {
	isHavePermission := uh.userService.VerifyUserFromCTX(w, r)
	if !isHavePermission {
		util.WriteResponseWithMssg(w, http.StatusForbidden, "the user does not meet role requirements")
		return
	}

	var payload *model.UserDeleteRequest

	payload, err := util.ParsePayloadWithValidator[model.UserDeleteRequest](w, r)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedRowsCount, err := uh.userService.DeleteUser(*payload)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	util.WriteResponseWithMssg(w, http.StatusOK, fmt.Sprintf("updated rows count in db: %d", updatedRowsCount))
}
