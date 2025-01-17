package db

import (
	"database/sql"
	"fmt"
	"jwt-auth/internal/entity"
	"log"
)

const (
	select_todo_query = "SELECT todo_id, user_id, title, description, status, created_date, deleted_date FROM todos"
	insert_todo_query = "INSERT INTO todos (user_id, title, description, status, created_date, deleted_date) VALUES(?, ?, ?, ?, ?, NULL)"
)

type TodoStore struct {
	db *sql.DB
}

func NewTodoStore(db *sql.DB) *TodoStore {
	return &TodoStore{db}
}

func (ts *TodoStore) CreateTodo(entity entity.TodoEntity) (int, error) {
	res, err := ts.db.Exec(insert_todo_query, entity.UserId, entity.Title, entity.Description, string(entity.Status), entity.CreatedDate)
	if err != nil {
		log.Printf("%v", err)
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("%v", err)
		return 0, err
	}
	return int(id), nil
}

func (ts *TodoStore) GetAll() ([]entity.TodoEntity, error) {
	result := make([]entity.TodoEntity, 0)
	rows, err := ts.db.Query(fmt.Sprintf("%s %s", select_todo_query, " WHERE deleted_date IS NULL;"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var todo entity.TodoEntity
		var status string
		err := rows.Scan(&todo.Id, &todo.UserId, &todo.Title, &todo.Description, &status, &todo.CreatedDate, &todo.DeletedDate)
		if err != nil {
			return nil, err
		}

		todo.Status = mapStatus(status)
		result = append(result, todo)
	}
	return result, nil
}

func (ts *TodoStore) GetById(id int) (entity.TodoEntity, error) {
	return entity.TodoEntity{}, nil
}

func (ts *TodoStore) UpdateStatus(id int, status string) (int, error) {
	return 0, nil
}

func (ts *TodoStore) DeleteTodo(id int) (int, error) {
	return 0, nil
}

func mapStatus(status string) entity.TodoStatus {
	switch status {
	case string(entity.CREATED):
		return entity.CREATED
	case string(entity.PENDING):
		return entity.PENDING
	case string(entity.DONE):
		return entity.DONE
	case string(entity.CANCELED):
		return entity.CANCELED
	default:
		return ""
	}
}
