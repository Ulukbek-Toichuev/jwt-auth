package util

import "fmt"

type ErrorCustom struct {
	Code    int
	Message string
}

func (e ErrorCustom) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}
