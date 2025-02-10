package tasks

import (
	"encoding/json"

	"github.com/FiodhyAN/learn-rest-api-golang/internal/repository"
	"github.com/hibiken/asynq"
)

const TypeVerificationEmail = "email:verification"

func NewVerificationEmail(user repository.User) (*asynq.Task, error) {
	payload, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	task := asynq.NewTask(TypeVerificationEmail, payload)
	return task, nil
}
