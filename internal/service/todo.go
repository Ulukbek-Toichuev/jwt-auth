package service

import (
	"database/sql"
	store "jwt-auth/internal/db"
	"jwt-auth/internal/entity"
	"jwt-auth/internal/model"
	"time"
)

type Todo interface {
	CreateTodo(entity.TodoEntity) (int, error)
	GetAll() ([]entity.TodoEntity, error)
	GetAllByUserId(userId int) ([]entity.TodoEntity, error)
	GetById(id int) (entity.TodoEntity, error)
	UpdateStatus(id int, status entity.Todo_status) (int, error)
	DeleteById(id int) (int, error)
}

type TodoService struct {
	todoStore Todo
}

func NewTodoService(db *sql.DB) *TodoService {
	todoStore := store.NewTodoStore(db)
	return &TodoService{todoStore}
}

func (ts *TodoService) CreateTodo(model model.TodoCreateRequest) (int, error) {
	entity := entity.TodoEntity{
		Id:          0,
		UserId:      model.UserId,
		Title:       model.Title,
		Description: model.Description,
		Status:      entity.CREATED,
		CreatedDate: time.Now(),
		DeletedDate: sql.NullTime{},
	}
	result, err := ts.todoStore.CreateTodo(entity)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (ts *TodoService) GetAllByUserId(userId int) ([]model.TodoModelResponse, error) {
	tmp, err := ts.todoStore.GetAllByUserId(userId)
	if err != nil {
		return []model.TodoModelResponse{}, err
	}
	result := make([]model.TodoModelResponse, len(tmp))
	for i, v := range tmp {
		result[i] = *model.NewTodoModelResponse(v.Id, v.UserId, v.Title, v.Description, string(v.Status), v.CreatedDate)
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

func (ts *TodoService) GetById(id int) (model.TodoModelResponse, error) {
	tmp, err := ts.todoStore.GetById(id)
	if err != nil {
		return model.TodoModelResponse{}, err
	}

	result := model.TodoModelResponse{
		Id:          tmp.Id,
		UserId:      tmp.UserId,
		Title:       tmp.Title,
		Description: tmp.Description,
		Status:      string(tmp.Status),
		CreatedDate: tmp.CreatedDate,
	}

	return result, nil
}

func (ts *TodoService) UpdateStatus(id int, status entity.Todo_status) (int, error) {
	res, err := ts.todoStore.UpdateStatus(id, status)
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (ts *TodoService) DeleteById(id int) (int, error) {
	res, err := ts.todoStore.DeleteById(id)
	if err != nil {
		return 0, err
	}

	return res, nil
}
