package service

import (
	"jwt-auth/internal/entity"
	"jwt-auth/internal/model"
)

type Todo interface {
	CreateTodo(model.TodoCreateRequest) (int, error)
	GetAll() ([]entity.TodoEntity, error)
	GetById(id int) (entity.TodoEntity, error)
	UpdateStatus(id int, status string) (int, error)
	DeleteTodo(id int) (int, error)
}

type TodoService struct {
	todoStore Todo
}

// func NewTodoService(db *sql.DB) *TodoService {
// 	todoStore := store.NewTodoStore(db)
// 	return &TodoService{todoStore}
// }

func (ts *TodoService) CreateTodo(model model.TodoCreateRequest) (int, error) {
	result, err := ts.todoStore.CreateTodo(model)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (ts *TodoService) GetAll() ([]model.TodoModelResponse, error) {
	tmp, err := ts.todoStore.GetAll()
	if err != nil {
		return []model.TodoModelResponse{}, err
	}

	result := make([]model.TodoModelResponse, len(tmp))
	for i, v := range tmp {
		result[i] = *model.NewTodoModelResponse(v.Id, v.UserId, v.Title, v.Description, string(v.Status), v.CreatedDate)
	}

	return result, nil
}
