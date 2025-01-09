package todo

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
