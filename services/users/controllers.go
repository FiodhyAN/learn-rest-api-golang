package users

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FiodhyAN/learn-rest-api-golang/auth"
	"github.com/FiodhyAN/learn-rest-api-golang/config"
	"github.com/FiodhyAN/learn-rest-api-golang/helpers"
	"github.com/FiodhyAN/learn-rest-api-golang/internal/repository"
	"github.com/FiodhyAN/learn-rest-api-golang/tasks"
	"github.com/FiodhyAN/learn-rest-api-golang/types"
	"github.com/FiodhyAN/learn-rest-api-golang/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	errorMessage := "Login Failed"

	if err := helpers.ParseJSON(r, &payload); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errorMessage, fmt.Errorf("invalid json payload"))
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteError(w, http.StatusBadRequest, errorMessage, fmt.Errorf("validation error: %v", errors))
		return
	}

	user, err := h.store.GetUser(r.Context(), payload.Username)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, errorMessage, err)
		return
	}

	if user == nil {
		helpers.WriteError(w, http.StatusBadRequest, errorMessage, fmt.Errorf("user not found"))
		return
	}

	err = auth.ComparePassword(payload.Password, user.Password)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, errorMessage, err)
		return
	}

	if !user.EmailVerified {
		if time.Now().After(user.EmailVerificationTokenExpiresAt.Time) {
			task, err := tasks.NewVerificationEmail(*user)
			if err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, errorMessage, err)
				return
			}

			if err := utils.EnqueueTask(task); err != nil {
				helpers.WriteError(w, http.StatusInternalServerError, errorMessage, err)
				return
			}
		}
		helpers.WriteError(w, http.StatusBadRequest, errorMessage, fmt.Errorf("email not verified, please verify your email"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, user.ID)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, errorMessage, err)
		return
	}

	response := types.LoginResponse{
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}

	helpers.WriteJSON(w, http.StatusOK, "Login Successfully", response, nil)
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterPayload
	errorMessage := "Failed To Register User"

	if err := helpers.ParseJSON(r, &payload); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errorMessage, fmt.Errorf("invalid json payload"))
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteError(w, http.StatusBadRequest, errorMessage, fmt.Errorf("validation error: %v", errors))
		return
	}

	user, err := h.store.GetUser(r.Context(), payload.Username)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, errorMessage, err)
		return
	}

	if user != nil {
		helpers.WriteError(w, http.StatusBadRequest, errorMessage, fmt.Errorf("username or email already exist"))
		return
	}

	hashedPassword, err := auth.CreateHashPassword(payload.Password)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, errorMessage, fmt.Errorf("error creating password"))
		return
	}

	createdUser, err := h.store.CreateUser(r.Context(), repository.CreateUserParams{
		Name:     payload.Name,
		Username: payload.Username,
		Email:    payload.Email,
		Password: hashedPassword,
	})
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, errorMessage, err)
		return
	}

	task, err := tasks.NewVerificationEmail(*createdUser)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, errorMessage, err)
		return
	}

	if err := utils.EnqueueTask(task); err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, errorMessage, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, "Successfully registered", map[string]string{
		"Name":     payload.Name,
		"username": payload.Username,
		"email":    payload.Email,
		"role":     createdUser.Role,
	}, nil)
}

func (h *Handler) handleVerifyEmail(w http.ResponseWriter, r *http.Request) {
	var payload types.VerifyEmailPayload
	errorMessage := "Email Verification Failed"

	if err := helpers.ParseJSON(r, &payload); err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, errorMessage, fmt.Errorf("invalid json payload"))
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteError(w, http.StatusBadRequest, errorMessage, fmt.Errorf("validation error: %v", errors))
		return
	}

	userIdString, err := utils.DecryptText(payload.UserId)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, errorMessage, err)
		return
	}

	userId, err := uuid.Parse(userIdString)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, errorMessage, err)
		return
	}

	user, err := h.store.GetUserById(r.Context(), userId)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, errorMessage, err)
		return
	}

	if user.EmailVerified {
		helpers.WriteError(w, http.StatusInternalServerError, errorMessage, fmt.Errorf("email already verified"))
		return
	}

	if user.EmailVerificationTokenExpiresAt.Valid {
		now := time.Now()
		if user.EmailVerificationTokenExpiresAt.Time.Before(now) {
			helpers.WriteError(w, http.StatusInternalServerError, errorMessage, fmt.Errorf("token expired"))
			return
		}
	} else {
		helpers.WriteError(w, http.StatusInternalServerError, errorMessage, fmt.Errorf("token expired"))
	}

	verificationToken, err := utils.DecryptText(payload.Token)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, errorMessage, fmt.Errorf("invalid token"))
	}

	if user.EmailVerificationToken.Valid {
		err = auth.ComparePassword(verificationToken, user.EmailVerificationToken.String)
		if err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, errorMessage, err)
			return
		}
	} else {
		helpers.WriteError(w, http.StatusBadRequest, errorMessage, fmt.Errorf("invalid token"))
		return
	}

	err = h.store.VerifyEmail(r.Context(), user.ID)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, errorMessage, err)
		return
	}

	helpers.WriteJSON(w, 200, "Successfully verify email", true, nil)
}

func (h *Handler) handleTest(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIdFromContext(r.Context())
	log.Println(userID)
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
