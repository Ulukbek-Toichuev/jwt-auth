package entity

import (
	"database/sql"
	"time"
)

type todoStatus string

const (
	CREATED  todoStatus = "CREATED"
	PENDING  todoStatus = "PENDING"
	DONE     todoStatus = "DONE"
	CANCELED todoStatus = "CANCELED"
)

type TodoEntity struct {
	Id          int
	UserId      int
	Title       string
	Description string
	Status      todoStatus
	CreatedDate time.Time
	DeletedDate sql.NullTime
}

func NewTodoEntity() *TodoEntity {
	return &TodoEntity{}
}
