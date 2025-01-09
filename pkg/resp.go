package pkg

import (
	"encoding/json"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	response := NewGeneralMessageResponse(message)
	json.NewEncoder(w).Encode(response)
}
