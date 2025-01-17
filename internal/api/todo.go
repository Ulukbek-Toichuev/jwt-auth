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

type TodoHandler struct {
	userService *service.UserService
	todoService *service.TodoService
}

func NewTodoHandler(db *sql.DB) *TodoHandler {
	return &TodoHandler{service.NewUserService(db), service.NewTodoService(db)}
}

func (th *TodoHandler) GetAllByUser(w http.ResponseWriter, r *http.Request) {
	isHavePermission := th.userService.VerifyUserFromCTX(w, r)
	if !isHavePermission {
		util.WriteResponseWithMssg(w, http.StatusForbidden, "the user does not meet role requirements")
		return
	}

	todos, err := th.todoService.GetAll()
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	util.WriteResponse(w, http.StatusOK, todos)
}

func (th *TodoHandler) GetAllByStatus(w http.ResponseWriter, r *http.Request) {

}

func (th *TodoHandler) GetTodoById(w http.ResponseWriter, r *http.Request) {

}

func (th *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	result := r.Context().Value(middleware.ResultCtxKey).(map[string]interface{})

	currUserEmail := ""
	if value, ok := result["email"]; !ok {
		util.WriteResponseWithMssg(w, http.StatusForbidden, "Token error")
		return
	} else {
		currUserEmail = value.(string)
	}

	res, err := th.userService.GetUserByEmail(currUserEmail)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	var payload *model.TodoCreateRequest
	payload, err = util.ParsePayloadWithValidator[model.TodoCreateRequest](w, r)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusBadRequest, err.Error())
		return
	}
	payload.UserId = res.UserId
	_, err = th.todoService.CreateTodo(*payload)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}
	util.WriteResponseWithMssg(w, http.StatusOK, "Todo successfully created!")
}

func (th *TodoHandler) UpdateTodoStatus(w http.ResponseWriter, r *http.Request) {

}

func (th *TodoHandler) DeleteTodoById(w http.ResponseWriter, r *http.Request) {

}
