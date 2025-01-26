package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

type JSONResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
}

func ParseJSON(r *http.Request, value any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}

	return json.NewDecoder(r.Body).Decode(value)
}

func WriteJSON(w http.ResponseWriter, status int, message string, value interface{}, errors interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if errors == nil {
		errors = map[string]string{}
	}

	response := JSONResponse{
		Status:  status,
		Message: "testing",
		Data:    value,
		Errors:  errors,
	}
	return json.NewEncoder(w).Encode(response)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, "error", nil, err.Error())
}
