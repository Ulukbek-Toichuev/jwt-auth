package pkg

import (
	"encoding/json"
	"net/http"
)

const (
	headerKey string = "Content-Type"
	headerVal string = "application/json"
)

func WriteResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set(headerKey, headerVal)
	w.WriteHeader(statusCode)
	response := NewGeneralMessageResponse(message)
	json.NewEncoder(w).Encode(response)
}
