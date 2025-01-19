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
	"strconv"
)

type TodoHandler struct {
	userService *service.UserService
	todoService *service.TodoService
}

func NewTodoHandler(db *sql.DB) *TodoHandler {
	return &TodoHandler{service.NewUserService(db), service.NewTodoService(db)}
}

func (th *TodoHandler) GetAllByUser(w http.ResponseWriter, r *http.Request) {
	result := r.Context().Value(middleware.ResultCtxKey).(map[string]interface{})

	currUserEmail := ""
	if value, ok := result["email"]; !ok {
		util.WriteResponseWithMssg(w, http.StatusForbidden, "the user does not meet role requirements")
		return
	} else {
		currUserEmail = value.(string)
	}

	res, err := th.userService.GetUserByEmail(currUserEmail)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	todos, err := th.todoService.GetAllByUserId(res.UserId)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	util.WriteResponse(w, http.StatusOK, todos)
}

func (th *TodoHandler) GetAllByStatus(w http.ResponseWriter, r *http.Request) {

}

func (th *TodoHandler) GetTodoById(w http.ResponseWriter, r *http.Request) {
	isHavePermission := th.userService.VerifyUserFromCTX(w, r)
	if !isHavePermission {
		util.WriteResponseWithMssg(w, http.StatusForbidden, "the user does not meet role requirements")
		return
	}

	todoId := r.PathValue("id")
	todoIdInt, err := strconv.Atoi(todoId)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusBadRequest, "todos id must be digit")
		return
	}

	todo, err := th.todoService.GetById(todoIdInt)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusNotFound, fmt.Sprintf("%v", err))
		return
	}

	util.WriteResponse(w, http.StatusOK, todo)
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
	isHavePermission := th.userService.VerifyUserFromCTX(w, r)
	if !isHavePermission {
		util.WriteResponseWithMssg(w, http.StatusForbidden, "the user does not meet role requirements")
		return
	}

	todoId := r.PathValue("id")
	todoIdInt, err := strconv.Atoi(todoId)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusBadRequest, "todos id must be digit")
		return
	}

	var payload *model.TodoStatusRequest
	payload, err = util.ParsePayloadWithValidator[model.TodoStatusRequest](w, r)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedRowsCount, err := th.todoService.UpdateStatus(todoIdInt, entity.TodoStatus(payload.Status))
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	util.WriteResponseWithMssg(w, http.StatusOK, fmt.Sprintf("todo succesfully update: %d", updatedRowsCount))
}

func (th *TodoHandler) DeleteTodoById(w http.ResponseWriter, r *http.Request) {

}
