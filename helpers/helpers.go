package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

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

	response := JSONResponse{
		Status:  status,
		Message: message,
		Data:    value,
		Errors:  errors,
	}
	return json.NewEncoder(w).Encode(response)
}

func WriteError(w http.ResponseWriter, status int, message string, err error) {
	WriteJSON(w, status, message, nil, err.Error())
}

func CreateHashPassword(plainPass string) (string, error) {
	passByte := []byte(plainPass)
	hashedPass, err := bcrypt.GenerateFromPassword(passByte, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPass), nil
}
