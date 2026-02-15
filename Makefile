.PHONY: all build test clean dev api fetcher scheduler migrate docker-up docker-down web lint

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_DIR=bin

# Docker Compose detection (prefer 'docker compose' over 'docker-compose')
DOCKER_COMPOSE := $(shell command -v docker-compose 2>/dev/null && echo docker-compose || echo "docker compose")

# Build all services
all: build

build: build-api build-fetcher build-scheduler build-migrator

build-api:
	$(GOBUILD) -o $(BINARY_DIR)/api ./cmd/api

build-fetcher:
	$(GOBUILD) -o $(BINARY_DIR)/fetcher ./cmd/fetcher

build-scheduler:
	$(GOBUILD) -o $(BINARY_DIR)/scheduler ./cmd/scheduler

build-migrator:
	$(GOBUILD) -o $(BINARY_DIR)/migrator ./cmd/migrator

# Run services
api:
	$(GOCMD) run ./cmd/api

fetcher:
	$(GOCMD) run ./cmd/fetcher

scheduler:
	$(GOCMD) run ./cmd/scheduler

# Development
dev:
	$(DOCKER_COMPOSE) up -d postgres redis
	@echo "Waiting for services to start..."
	@sleep 3
	$(GOCMD) run ./cmd/api

dev-all:
	$(DOCKER_COMPOSE) up

# Database
migrate:
	$(GOCMD) run ./cmd/migrator up

migrate-down:
	$(GOCMD) run ./cmd/migrator down

migrate-create:
	@read -p "Migration name: " name; \
	touch migrations/$$(date +%Y%m%d%H%M%S)_$$name.up.sql; \
	touch migrations/$$(date +%Y%m%d%H%M%S)_$$name.down.sql

# Testing
test:
	$(GOTEST) -v -race ./...

test-coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Linting
lint:
	golangci-lint run ./...

# Docker
docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-build:
	$(DOCKER_COMPOSE) build

docker-logs:
	$(DOCKER_COMPOSE) logs -f

# Frontend
web:
	cd web && npm run dev

web-build:
	cd web && npm run build

web-install:
	cd web && npm install

# Dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Clean
clean:
	rm -rf $(BINARY_DIR)
	rm -rf coverage.out coverage.html
	rm -rf web/dist

# Help
help:
	@echo "Available commands:"
	@echo "  make build        - Build all Go services"
	@echo "  make api          - Run API service"
	@echo "  make fetcher      - Run Fetcher service"
	@echo "  make scheduler    - Run Scheduler service"
	@echo "  make dev          - Start dev environment (postgres, redis + api)"
	@echo "  make dev-all      - Start all services via docker-compose"
	@echo "  make migrate      - Run database migrations"
	@echo "  make test         - Run tests"
	@echo "  make docker-up    - Start docker containers"
	@echo "  make docker-down  - Stop docker containers"
	@echo "  make web          - Run frontend dev server"
	@echo "  make clean        - Clean build artifacts"
