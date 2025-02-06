package users

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FiodhyAN/learn-rest-api-golang/auth"
	"github.com/FiodhyAN/learn-rest-api-golang/config"
	"github.com/FiodhyAN/learn-rest-api-golang/tasks"
	"github.com/FiodhyAN/learn-rest-api-golang/types"
	"github.com/FiodhyAN/learn-rest-api-golang/utils"
	"github.com/go-playground/validator/v10"
	"github.com/hibiken/asynq"
)

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid JSON payload"))
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("payload validation error: %v", errors))
		return
	}

	user, err := h.store.GetUser(payload.Username)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if user == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user does not exists"))
		return
	}

	err = auth.ComparePassword(payload.Password, user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("wrong password"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Login Successfully", map[string]string{"username": user.Username, "email": user.Email}, nil)
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid JSON payload"))
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("payload validation error: %v", errors))
		return
	}

	user, err := h.store.GetUser(payload.Username)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if user != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user already exists"))
		return
	}

	hashedPassword, err := auth.CreateHashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("hash password error"))
		return
	}

	createdUser, err := h.store.CreateUser(types.User{
		Name:     payload.Name,
		Username: payload.Username,
		Email:    payload.Email,
		Password: hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: config.Envs.RedisHost + ":" + config.Envs.RedisPort})
	defer client.Close()

	task, err := tasks.NewVerificationEmail(*createdUser)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	info, err := client.Enqueue(task)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
	log.Println("info id: " + info.ID + "info queue: " + info.Queue)

	utils.WriteJSON(w, http.StatusOK, "Successfully registered", map[string]string{
		"Name":     payload.Name,
		"username": payload.Username,
		"email":    payload.Email,
		"role":     createdUser.Role,
	}, nil)
}

func (h *Handler) handleVerifyEmail(w http.ResponseWriter, r *http.Request) {
	var payload types.VerifyEmailPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("invalid JSON payload"))
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("payload validation error: %v", errors))
		return
	}

	userId, err := utils.DecryptText(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	user, err := h.store.GetUserById(string(userId))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if user.EmailVerified {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("email already verified"))
		return
	}

	if user.EmailVerificationExpiresAt.Valid {
		now := time.Now()
		if user.EmailVerificationExpiresAt.Time.Before(now) {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("token expired"))
			return
		}
	} else {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("token expired"))
	}

	verificationToken, err := utils.DecryptText(payload.Token)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("invalid token"))
	}

	if user.EmailVerificationToken.Valid {
		err = auth.ComparePassword(verificationToken, user.EmailVerificationToken.String)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid token"))
		return
	}

	err = h.store.VerifyEmail(user)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, 200, "Successfully verify email", true, nil)
}

func (h *Handler) handleTest(w http.ResponseWriter, r *http.Request) {
	encrypted, err := utils.EncryptText(`35f8a28e-487e-45d7-95f9-f67fb56a2d76`)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(encrypted)
	decrypted, err := utils.DecryptText(encrypted)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(decrypted), decrypted)
}
