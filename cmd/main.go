package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/FiodhyAN/learn-rest-api-golang/cmd/api"
	"github.com/FiodhyAN/learn-rest-api-golang/config"
	"github.com/FiodhyAN/learn-rest-api-golang/db"
)

func main() {
	dbConn, err := db.NewPostgresSQL(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Envs.DBHost, config.Envs.DBPort, config.Envs.DBUser, config.Envs.DBPassword, config.Envs.DBName))
	if err != nil {
		log.Fatal("error connecting to postgres : ", err)
	}
	if err := initDB(dbConn); err != nil {
		log.Fatal("connection with db error : ", err)
	}
	apiServer := api.NewAPIServer(":8080")
	if err := apiServer.Run(); err != nil {
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
