package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ParseJSON(r *http.Request, value any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}

	return json.NewDecoder(r.Body).Decode(value)
}

func WriteJSON(w http.ResponseWriter, status int, value any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	log.Println(value)

	return json.NewEncoder(w).Encode(value)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}
