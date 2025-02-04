package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/FiodhyAN/learn-rest-api-golang/config"
	"github.com/FiodhyAN/learn-rest-api-golang/db"
	"github.com/FiodhyAN/learn-rest-api-golang/services/users"
	"github.com/FiodhyAN/learn-rest-api-golang/tasks"
	"github.com/hibiken/asynq"
)

type WorkerServer struct {
	addr string
	db   *sql.DB
}

func NewWorkerServer(addr string, db *sql.DB) *WorkerServer {
	return &WorkerServer{
		addr: addr,
		db:   db,
	}
}

func (s *WorkerServer) Run() error {
	userStore := users.NewUserStore(s.db)
	handler := tasks.NewHandler(userStore)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: config.Envs.RedisHost + ":" + config.Envs.RedisPort},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeVerificationEmail, handler.HandleVerificationEmailTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
		return err
	}

	return nil
}

func main() {
	dbConn, err := db.NewPostgresSQL(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Envs.DBHost, config.Envs.DBPort, config.Envs.DBUser, config.Envs.DBPassword, config.Envs.DBName))
	if err != nil {
		log.Fatal("error connecting to postgres : ", err, config.Envs)
	}
	if err := initDB(dbConn); err != nil {
		log.Fatal("connection with db error : ", err)
	}
	workerServer := NewWorkerServer(":8080", dbConn)
	if err := workerServer.Run(); err != nil {
		log.Fatal("error running api server")
	}
}

func initDB(db *sql.DB) error {
	err := db.Ping()

	if err != nil {
		return err
	}

	log.Println("Connected to database : ", config.Envs.DBName)
	return nil
}
