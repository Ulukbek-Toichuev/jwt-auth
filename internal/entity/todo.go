package entity

import (
	"database/sql"
	"time"
)

type Todo_status string

const (
	CREATED     Todo_status = "CREATED"
	PENDING     Todo_status = "PENDING"
	IN_PROGRESS Todo_status = "IN_PROGRESS"
	DONE        Todo_status = "DONE"
	CANCELED    Todo_status = "CANCELED"
)

type TodoEntity struct {
	Id          int
	UserId      int
	Title       string
	Description string
	Status      Todo_status
	CreatedDate time.Time
	DeletedDate sql.NullTime
}

func NewTodoEntity() *TodoEntity {
	return &TodoEntity{}
}
