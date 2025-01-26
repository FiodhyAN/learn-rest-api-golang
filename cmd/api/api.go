package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/FiodhyAN/learn-rest-api-golang/services/users"
	"github.com/go-chi/cors"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	corsHandler := c.Handler(router)

	public := router.PathPrefix("/api").Subrouter()
	userStore := users.NewUserStore(s.db)
	userHandler := users.NewHandler(userStore)
	userHandler.RegisterRoutes(public)

	log.Println("Listening on port", s.addr)

	return http.ListenAndServe(s.addr, corsHandler)
}
