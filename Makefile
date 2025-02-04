build:
	@go build -o bin/restapi cmd/main.go

run: build
	@./bin/restapi

build-worker:
	@go build -o bin/worker cmd/worker/main.go

run-worker: build-worker
	@./bin/worker

air:
	@air

migration: 
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@, $(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down