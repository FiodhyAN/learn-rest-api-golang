package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/FiodhyAN/learn-rest-api-golang/types"
	"github.com/FiodhyAN/learn-rest-api-golang/utils"
	"github.com/hibiken/asynq"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) HandleVerificationEmailTask(c context.Context, t *asynq.Task) error {
	var user types.User

	if err := json.Unmarshal(t.Payload(), &user); err != nil {
		return fmt.Errorf("json Unmarshal Failed: %v: %w", err, asynq.SkipRetry)
	}

	if err := utils.SendVerificationMail(h.store, &user); err != nil {
		return fmt.Errorf("Error Sending email: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Successfully sending email to %s", user.Email)
	return nil
}
