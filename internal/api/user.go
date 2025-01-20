package api

import (
	"database/sql"
	"fmt"
	"jwt-auth/internal/entity"
	"jwt-auth/internal/middleware"
	"jwt-auth/internal/model"
	"jwt-auth/internal/service"
	"jwt-auth/internal/util"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService service.UserService
	validate    *validator.Validate
}

func NewUserHandler(db *sql.DB, v *validator.Validate) *UserHandler {
	us := service.NewUserService(db)
	return &UserHandler{*us, v}
}

// Хендлер для получения списка всех пользователей
// Права доступа ADMIN
func (uh *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var adminPermission entity.Role = entity.ADMIN
	isHavePermission := uh.userService.VerifyIsUserHavePermissionFromCTX(w, r, adminPermission)
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

// Хендлер для получения пользователя по его email
// Права доступа ADMIN
func (uh *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	var adminPermission entity.Role = entity.ADMIN
	isHavePermission := uh.userService.VerifyIsUserHavePermissionFromCTX(w, r, adminPermission)
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

// Хендлер для получение собственных данных пользователем
func (uh *UserHandler) GetUsersOwnDetail(w http.ResponseWriter, r *http.Request) {
	result := r.Context().Value(middleware.ResultCtxKey).(map[string]interface{})

	currUserEmail := ""
	if value, ok := result["email"]; !ok {
		util.WriteResponseWithMssg(w, http.StatusForbidden, "incorrect token claims, please check your token")
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

// Хендлер для смены ролей пользователя
// Права доступа ADMIN
func (uh *UserHandler) ChangeUsersRole(w http.ResponseWriter, r *http.Request) {
	var adminPermission entity.Role = entity.ADMIN
	isHavePermission := uh.userService.VerifyIsUserHavePermissionFromCTX(w, r, adminPermission)
	if !isHavePermission {
		util.WriteResponseWithMssg(w, http.StatusForbidden, "the user does not meet role requirements")
		return
	}

	var payload *model.UserChangeRoleRequest

	payload, err := util.ParsePayloadWithValidator[model.UserChangeRoleRequest](w, r, uh.validate)
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

// Хендлер для удаления пользователя
// Права доступа ADMIN
func (uh *UserHandler) DeleteUserByEmail(w http.ResponseWriter, r *http.Request) {
	var adminPermission entity.Role = entity.ADMIN
	isHavePermission := uh.userService.VerifyIsUserHavePermissionFromCTX(w, r, adminPermission)
	if !isHavePermission {
		util.WriteResponseWithMssg(w, http.StatusForbidden, "the user does not meet role requirements")
		return
	}

	var payload *model.UserDeleteRequest

	payload, err := util.ParsePayloadWithValidator[model.UserDeleteRequest](w, r, uh.validate)
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
