package model

import (
	"time"
)

type TodoCreateRequest struct {
	UserId      int
	Title       string
	Description string
}

func NewTodoCreateRequest() *TodoCreateRequest {
	return &TodoCreateRequest{}
}

type TodoModelResponse struct {
	Id          int
	UserId      int
	Title       string
	Description string
	Status      string
	CreatedDate time.Time
}

func NewTodoModelResponse(id, userId int, title, description, status string, cdt time.Time) *TodoModelResponse {
	return &TodoModelResponse{
		id, userId, title, description, status, cdt,
	}
}
