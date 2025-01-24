package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	m, err := migrate.New(
		"file://cmd/migrate/migrations",
		"postgres://postgres:dbnyafio@localhost:5432/learn_go_rest_api?sslmode=disable",
	)

	if err != nil {
		log.Fatal("error initialize migration : ", err)
	}

	cmd := os.Args[len(os.Args)-1]

	if cmd == "up" {
		if err := m.Up(); err != nil {
			log.Fatal("error migrate up : ", err)
		} else {
			log.Println("success migrate up")
		}
	}
	if cmd == "down" {
		if err := m.Down(); err != nil {
			log.Fatal("error migrate down : ", err)
		} else {
			log.Println("success migrate down")
		}
	}
}
