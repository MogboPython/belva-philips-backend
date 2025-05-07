include .env.development

.PHONY: build run dev test swagger lint migrate-status migrate-create migrate-create-sql migrate-up migrate-down help deploy

BINARY_NAME=belvaphilips_backend
MIGRATIONS_DIR=./internal/database/migrations
DB_DRIVER=postgres


all: swagger build

build:
	@echo "Building application..."
	go build -o bin/app cmd/app/main.go

run: build
	@echo "Running application..."
	./bin/app

dev:
	@echo "Running in development mode..."
	go install github.com/air-verse/air@latest
	~/go/bin/air -c .air.toml

test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
# test-coverage:
# 	@echo "Running tests with coverage..."
# 	go test -v -coverprofile=coverage.out ./...
# 	go tool cover -html=coverage.out

swagger:
	@echo "Generating Swagger documentation..."
	go install github.com/swaggo/swag/cmd/swag@latest
	~/go/bin/swag init -g cmd/app/main.go -o ./cmd/app/docs

git:
	git add .
	git commit -m "$m"
	git push -u origin main
	
# swag init -g cmd/api/main.go -o ./docs

check-migration-name:
ifndef name
	$(error name is required)
endif

# Show migration status
migrate-status:
	@echo "Checking migration status..."
	goose -dir ${MIGRATIONS_DIR} postgres "$(DIRECT_URL)" status

migrate-create: check-migration-name
	@echo "Creating new Go migration '$(name)'..."
	goose -dir ${MIGRATIONS_DIR} create $(name) go

migrate-create-sql: check-migration-name
	@echo "Creating new SQL migration '$(name)'..."
	goose -dir ${MIGRATIONS_DIR} create $(name) sql

migrate-up:
	@echo "Running migrations up..."
	goose -dir ${MIGRATIONS_DIR} ${DB_DRIVER} "$(DIRECT_URL)" up

migrate-down:
	@echo "Running migrations down..."
	goose -dir ${MIGRATIONS_DIR} ${DB_DRIVER} "$(DIRECT_URL)" down

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

deploy:
	@echo "Deploying application..."
	fly deploy --local-only

# Print help information
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  dev            - Run in development mode with hot reload"
	@echo "  swagger        - Generate Swagger documentation"
	@echo "  lint           - Run linter"
	@echo "  help           - Print this help information"
	@echo "  deploy     	- Deploy the application to Fly.io"