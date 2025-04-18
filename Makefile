include .env.development

.PHONY: build run dev swagger lint migrate-status migrate-create migrate-create-sql migrate-up migrate-down help
# build run dev test clean swagger lint
BINARY_NAME=belvaphilips_backend
MIGRATIONS_DIR=./internal/database/migrations

# Default target
all: swagger build

# Build the application
build:
	@echo "Building application..."
	go build -o bin/app cmd/app/main.go

# Run the application
run: build
	@echo "Running application..."
	./bin/app

# Run the application in development mode with hot reload
dev:
	@echo "Running in development mode..."
	go install github.com/air-verse/air@latest
	~/go/bin/air -c .air.toml

# Run tests
# test:
# 	@echo "Running tests..."
# 	go test -v ./...

# Run tests with coverage
# test-coverage:
# 	@echo "Running tests with coverage..."
# 	go test -v -coverprofile=coverage.out ./...
# 	go tool cover -html=coverage.out

# Generate swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	go install github.com/swaggo/swag/cmd/swag@latest
	~/go/bin/swag init -g cmd/app/main.go -o ./cmd/app/docs

git:
	git add .
	git commit -m "$m"
	git push -u origin main
	
# swag init -g cmd/api/main.go -o ./docs

# Create database migrations
check-migration-name:
ifndef name
	$(error name is required)
endif

# Show migration status
migrate-status:
	@echo "Checking migration status..."
	goose -dir ${MIGRATIONS_DIR} postgres "$(DIRECT_URL)" status


# Create a new Go migration
migrate-create: check-migration-name
	@echo "Creating new Go migration '$(name)'..."
	goose -dir ${MIGRATIONS_DIR} create $(name) go

# Create a new SQL migration
migrate-create-sql: check-migration-name
	@echo "Creating new SQL migration '$(name)'..."
	goose -dir ${MIGRATIONS_DIR} create $(name) sql

# migrate-create:
# 	@echo "Creating migration..."
# 	@read -p "Enter migration name: " name; \
# 	migrate create -ext sql -dir migrations -seq $${name}

# Run database migrations up
migrate-up:
	@echo "Running migrations up..."
	goose -dir ${MIGRATIONS_DIR} postgres "$(DIRECT_URL)" up

# Run database migrations down
migrate-down:
	@echo "Running migrations down..."
	goose -dir ${MIGRATIONS_DIR} postgres "$(DIRECT_URL)" down

# Run linter
lint:
	@echo "Running linter..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run

# Clean build artifacts
# clean:
# 	@echo "Cleaning build artifacts..."
# 	rm -rf bin/
# 	rm -rf docs/

# Run the application with Docker
# docker-run:
# 	@echo "Building and running with Docker..."
# 	docker-compose up --build

# Print help information
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  dev            - Run in development mode with hot reload"
	@echo "  swagger        - Generate Swagger documentation"
	@echo "  lint           - Run linter"
	@echo "  help           - Print this help information"

# @echo "  test           - Run tests"
# @echo "  test-coverage  - Run tests with coverage report"
# @echo "  migrate-create - Create a new migration file"
# @echo "  migrate-up     - Run migrations up"
# @echo "  migrate-down   - Run migrations down"
# @echo "  clean          - Clean build artifacts"
# @echo "  docker-run     - Run with Docker"