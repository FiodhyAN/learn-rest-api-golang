package users

import (
	"net/http"

	"github.com/FiodhyAN/learn-rest-api-golang/auth"
	"github.com/FiodhyAN/learn-rest-api-golang/types"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/register", h.handleRegister).Methods(http.MethodPost)
	router.HandleFunc("/verify-email", h.handleVerifyEmail).Methods(http.MethodPost)
	router.HandleFunc("/test", auth.VerifyToken(h.handleTest, h.store)).Methods(http.MethodGet)
}
