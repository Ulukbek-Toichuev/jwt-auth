package util

import (
	"encoding/json"
	"jwt-auth/internal/model"
	"net/http"
)

const (
	headerKey string = "Content-Type"
	headerVal string = "application/json"
)

func WriteResponseWithMssg(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set(headerKey, headerVal)
	w.WriteHeader(statusCode)
	response := model.NewGeneralMessageResponse(message)
	json.NewEncoder(w).Encode(response)
}

func WriteResponse(w http.ResponseWriter, statusCode int, message any) {
	w.Header().Set(headerKey, headerVal)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)
}
