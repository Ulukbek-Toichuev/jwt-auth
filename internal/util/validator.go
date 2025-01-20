package util

import "github.com/go-playground/validator/v10"

// Кастомные валидаторы
var todoStatuses = []string{
	"CREATED",
	"PENDING",
	"IN_PROGRESS",
	"DONE",
	"CANCELED",
}

func ValidateTodoStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	for _, val := range todoStatuses {
		if val == status {
			return true
		}
	}
	return false
}
