package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"jwt-auth/pkg"
	"net/http"
)

type TodoHandler struct {
}

func NewTodoHandler(db *sql.DB) *TodoHandler {
	return &TodoHandler{}
}

func (th *TodoHandler) GetAllByUser(w http.ResponseWriter, r *http.Request) {
	result := r.Context().Value(pkg.ResultCtxKey).(map[string]interface{})

	currUserEmail := ""
	if value, ok := result["email"]; !ok {
		return
	} else {
		currUserEmail = value.(string)
	}
	response := pkg.NewGeneralMessageResponse(fmt.Sprintf("current user email - %s", currUserEmail))
	json.NewEncoder(w).Encode(response)
}

func (th *TodoHandler) GetAllByStatus(w http.ResponseWriter, r *http.Request) {

}

func (th *TodoHandler) GetTodoById(w http.ResponseWriter, r *http.Request) {

}

func (th *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {

}

func (th *TodoHandler) UpdateTodoStatus(w http.ResponseWriter, r *http.Request) {

}

func (th *TodoHandler) DeleteTodoById(w http.ResponseWriter, r *http.Request) {

}
