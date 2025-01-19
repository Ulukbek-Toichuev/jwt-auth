package api

import (
	"database/sql"
	"errors"
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

var ErrFobidden error

func NewTodoHandler(db *sql.DB) *TodoHandler {
	return &TodoHandler{service.NewUserService(db), service.NewTodoService(db)}
}

func (th *TodoHandler) GetAllByUser(w http.ResponseWriter, r *http.Request) {
	res, err := th.getUserByEmailFromCTX(r)
	if err != nil {
		var errCustom util.ErrorCustom
		if errors.As(err, &errCustom) {
			util.WriteResponseWithMssg(w, errCustom.Code, errCustom.Message)
		}
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
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

func (th *TodoHandler) GetTodoByIdAndUserId(w http.ResponseWriter, r *http.Request) {
	result := r.Context().Value(middleware.ResultCtxKey).(map[string]interface{})

	currUserEmail := ""
	if value, ok := result["email"]; !ok {
		util.WriteResponseWithMssg(w, http.StatusForbidden, "incorrect token claims, please check your token")
		return
	} else {
		currUserEmail = value.(string)
	}

	res, err := th.userService.GetUserByEmail(currUserEmail)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	todoId := r.PathValue("id")
	todoIdInt, err := strconv.Atoi(todoId)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusBadRequest, "todos id must be digit")
		return
	}

	todo, err := th.todoService.GetByIdAndByUserId(todoIdInt, res.UserId)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusNotFound, fmt.Sprintf("%v", err))
		return
	}

	util.WriteResponse(w, http.StatusOK, todo)
}

func (th *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	res, err := th.getUserByEmailFromCTX(r)
	if err != nil {
		var errCustom util.ErrorCustom
		if errors.As(err, &errCustom) {
			util.WriteResponseWithMssg(w, errCustom.Code, errCustom.Message)
		}
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
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
	res, err := th.getUserByEmailFromCTX(r)
	if err != nil {
		var errCustom util.ErrorCustom
		if errors.As(err, &errCustom) {
			util.WriteResponseWithMssg(w, errCustom.Code, errCustom.Message)
			return
		}
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
	}

	todoId := r.PathValue("id")
	todoIdInt, err := strconv.Atoi(todoId)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusBadRequest, "todos id must be digit")
		return
	}

	var payload *model.Todo_statusRequest
	payload, err = util.ParsePayloadWithValidator[model.Todo_statusRequest](w, r)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedRowsCount, err := th.todoService.UpdateStatus(todoIdInt, res.UserId, entity.Todo_status(payload.Status))
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	util.WriteResponseWithMssg(w, http.StatusOK, fmt.Sprintf("todo succesfully update: %d", updatedRowsCount))
}

func (th *TodoHandler) DeleteTodoById(w http.ResponseWriter, r *http.Request) {
	res, err := th.getUserByEmailFromCTX(r)
	if err != nil {
		var errCustom util.ErrorCustom
		if errors.As(err, &errCustom) {
			util.WriteResponseWithMssg(w, errCustom.Code, errCustom.Message)
			return
		}
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
	}

	todoId := r.PathValue("id")
	todoIdInt, err := strconv.Atoi(todoId)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusBadRequest, "todos id must be digit")
		return
	}
	updatedRowsCount, err := th.todoService.DeleteById(todoIdInt, res.UserId)
	if err != nil {
		util.WriteResponseWithMssg(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	util.WriteResponseWithMssg(w, http.StatusOK, fmt.Sprintf("todo succesfully delete: %d", updatedRowsCount))
}

func (th *TodoHandler) getUserByEmailFromCTX(r *http.Request) (model.UserResponse, error) {
	result := r.Context().Value(middleware.ResultCtxKey).(map[string]interface{})

	currUserEmail := ""
	if value, ok := result["email"]; !ok {
		return model.UserResponse{}, &util.ErrorCustom{Code: http.StatusForbidden, Message: "incorrect token claims, please check your token"}
	} else {
		currUserEmail = value.(string)
	}

	res, err := th.userService.GetUserByEmail(currUserEmail)
	if err != nil {
		return model.UserResponse{}, &util.ErrorCustom{Code: http.StatusForbidden, Message: fmt.Sprintf("%v", err)}
	}

	return res, nil
}
