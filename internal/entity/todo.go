package entity

import (
	"database/sql"
	"time"
)

type TodoStatus string

const (
	CREATED     TodoStatus = "CREATED"
	PENDING     TodoStatus = "PENDING"
	IN_PROGRESS TodoStatus = "IN_PROGRESS"
	DONE        TodoStatus = "DONE"
	CANCELED    TodoStatus = "CANCELED"
)

type TodoEntity struct {
	Id          int
	UserId      int
	Title       string
	Description string
	Status      TodoStatus
	CreatedDate time.Time
	DeletedDate sql.NullTime
}

func NewTodoEntity() *TodoEntity {
	return &TodoEntity{}
}
