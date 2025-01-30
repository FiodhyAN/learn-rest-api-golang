package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"

	"github.com/FiodhyAN/learn-rest-api-golang/config"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func init() {
	Validate.RegisterValidation("password", validatePassword)
}

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

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, "error", nil, err.Error())
}

func SendMail() error {
	toList := []string{"example@gmail.com"}
	subject := "Subject: Test Email\n"
	msg := "Hello geeks!!!\n"
	body := []byte("From: " + config.Envs.SMTPEmail + "\n" +
		"To: " + toList[0] + "\n" +
		subject + "\n" + msg)

	mail_host := config.Envs.SMTPHost
	mail_port := config.Envs.SMTPPort
	mail_username := config.Envs.SMTPUsername
	mail_from := config.Envs.SMTPEmail
	mail_password := config.Envs.SMTPPassword

	mail_auth := smtp.PlainAuth("", mail_username, mail_password, mail_host)

	if err := smtp.SendMail(mail_host+":"+mail_port, mail_auth, mail_from, toList, body); err != nil {
		return err
	}

	return nil
}
