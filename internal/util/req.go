package util

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func ParsePayload[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	var payload T
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}

func ParsePayloadWithValidator[T any](w http.ResponseWriter, r *http.Request, v *validator.Validate) (*T, error) {
	var payload T
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}

	err = v.Struct(&payload)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Printf("Validation error: %s - %s (Tag: %s, Param: %s)", e.Field(), e.Value(), e.Tag(), e.Param())
			return nil, e
		}
	}

	return &payload, nil
}
