package model

import (
	"time"
)

type TodoCreateRequest struct {
	UserId      int
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func NewTodoCreateRequest() *TodoCreateRequest {
	return &TodoCreateRequest{}
}

type TodoStatusRequest struct {
	Status string `json:"status" validate:"required"`
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
