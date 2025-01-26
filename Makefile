build:
	@go build -o bin/restapi cmd/main.go

run: build
	@./bin/restapi

air:
	@air -c "C:\Kerja\Pengembangan Diri\Belajar Golang\learn-rest-api\.air.toml"

migration: 
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@, $(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down