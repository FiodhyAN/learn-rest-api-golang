package users

import (
	"fmt"
	"net/http"

	"github.com/FiodhyAN/learn-rest-api-golang/auth"
	"github.com/FiodhyAN/learn-rest-api-golang/types"
	"github.com/FiodhyAN/learn-rest-api-golang/utils"
	"github.com/go-playground/validator/v10"
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
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("query database error"))
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

	if err := utils.SendMail(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to send email"))
		return
	}

	err = h.store.CreateUser(types.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: hashedPassword,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("query database error"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Successfully registered", map[string]string{
		"username": payload.Username,
		"email":    payload.Email,
	}, nil)
}
