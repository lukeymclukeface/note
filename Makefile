# Note AI Development Makefile
# Provides commands for building, running, and developing the Note AI application

.PHONY: help build up down restart logs clean dev dev-server dev-web rebuild status install test

# Default target
.DEFAULT_GOAL := help

# Colors for output
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
RED := \033[0;31m
NC := \033[0m # No Color

# Docker Compose files
COMPOSE_FILE := docker-compose.yml
COMPOSE_DEV_FILE := docker-compose.dev.yml

help: ## Show this help message
	@echo "$(BLUE)Note AI Development Commands$(NC)"
	@echo "================================"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

build: ## Build all Docker images
	@echo "$(BLUE)Building Docker images...$(NC)"
	docker-compose build --no-cache

build-server: ## Build only the server Docker image
	@echo "$(BLUE)Building server Docker image...$(NC)"
	docker-compose build --no-cache note-server

build-web: ## Build only the web Docker image
	@echo "$(BLUE)Building web Docker image...$(NC)"
	docker-compose build --no-cache note-web

up: ## Start all services in production mode
	@echo "$(BLUE)Starting services in production mode...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)Services started!$(NC)"
	@echo "$(YELLOW)Web app: http://localhost:3000$(NC)"
	@echo "$(YELLOW)Server: http://localhost:8080$(NC)"

down: ## Stop all services
	@echo "$(BLUE)Stopping services...$(NC)"
	docker-compose down
	@echo "$(GREEN)Services stopped!$(NC)"

restart: down up ## Restart all services

logs: ## Show logs from all services
	docker-compose logs -f

logs-server: ## Show logs from server only
	docker-compose logs -f note-server

logs-web: ## Show logs from web only
	docker-compose logs -f note-web

dev: ## Start development environment with hot reloading
	@echo "$(BLUE)Starting development environment with hot reloading...$(NC)"
	@if [ ! -f $(COMPOSE_DEV_FILE) ]; then \
		echo "$(YELLOW)Creating development docker-compose file...$(NC)"; \
		$(MAKE) create-dev-compose; \
	fi
	docker-compose -f $(COMPOSE_FILE) -f $(COMPOSE_DEV_FILE) up
	@echo "$(GREEN)Development environment started!$(NC)"
	@echo "$(YELLOW)Web app (dev): http://localhost:3000$(NC)"
	@echo "$(YELLOW)Server: http://localhost:8080$(NC)"

dev-build: ## Build and start development environment
	@echo "$(BLUE)Building and starting development environment...$(NC)"
	@if [ ! -f $(COMPOSE_DEV_FILE) ]; then \
		echo "$(YELLOW)Creating development docker-compose file...$(NC)"; \
		$(MAKE) create-dev-compose; \
	fi
	docker-compose -f $(COMPOSE_FILE) -f $(COMPOSE_DEV_FILE) up --build

dev-down: ## Stop development environment
	@echo "$(BLUE)Stopping development environment...$(NC)"
	docker-compose -f $(COMPOSE_FILE) -f $(COMPOSE_DEV_FILE) down

dev-server: ## Start only server in development mode
	@echo "$(BLUE)Starting server in development mode...$(NC)"
	cd note-server && air

dev-web: ## Start only web in development mode
	@echo "$(BLUE)Starting web in development mode...$(NC)"
	cd note-web && npm run dev

rebuild: down build up ## Rebuild and restart all services

rebuild-server: ## Rebuild and restart only server
	@echo "$(BLUE)Rebuilding server...$(NC)"
	docker-compose stop note-server
	docker-compose build --no-cache note-server
	docker-compose up -d note-server

rebuild-web: ## Rebuild and restart only web
	@echo "$(BLUE)Rebuilding web...$(NC)"
	docker-compose stop note-web
	docker-compose build --no-cache note-web
	docker-compose up -d note-web

status: ## Show status of all services
	@echo "$(BLUE)Service Status:$(NC)"
	docker-compose ps

clean: ## Clean up Docker resources (containers, images, volumes)
	@echo "$(BLUE)Cleaning up Docker resources...$(NC)"
	docker-compose down --volumes --remove-orphans
	docker system prune -f
	docker volume prune -f
	@echo "$(GREEN)Cleanup completed!$(NC)"

clean-all: ## Clean up everything including images
	@echo "$(RED)Warning: This will remove all containers, images, and volumes!$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker-compose down --volumes --remove-orphans; \
		docker system prune -a -f; \
		docker volume prune -f; \
		echo "$(GREEN)Complete cleanup finished!$(NC)"; \
	else \
		echo "$(YELLOW)Cleanup cancelled.$(NC)"; \
	fi

install: ## Install dependencies for development
	@echo "$(BLUE)Installing dependencies...$(NC)"
	@echo "$(YELLOW)Installing web dependencies...$(NC)"
	cd note-web && npm install
	@echo "$(YELLOW)Installing server dependencies...$(NC)"
	cd note-server && go mod download
	@echo "$(GREEN)Dependencies installed!$(NC)"

test: ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	@echo "$(YELLOW)Running web tests...$(NC)"
	cd note-web && npm test
	@echo "$(YELLOW)Running server tests...$(NC)"
	cd note-server && go test ./...

test-web: ## Run web tests only
	@echo "$(BLUE)Running web tests...$(NC)"
	cd note-web && npm test

test-server: ## Run server tests only
	@echo "$(BLUE)Running server tests...$(NC)"
	cd note-server && go test ./...

shell-server: ## Open shell in server container
	docker-compose exec note-server sh

shell-web: ## Open shell in web container
	docker-compose exec note-web sh

create-dev-compose: ## Create development docker-compose override file
	@echo "$(BLUE)Creating development docker-compose override...$(NC)"
	@printf '%s\n' \
		'services:' \
		'  note-server:' \
		'    build:' \
		'      context: ./note-server' \
		'      dockerfile: Dockerfile.dev' \
		'    volumes:' \
		'      - ./note-server:/app' \
		'      - /app/tmp' \
		'    environment:' \
		'      - DEV_MODE=true' \
		'      - LOG_LEVEL=debug' \
		'    command: air -c .air.toml' \
		'' \
		'  note-web:' \
		'    build:' \
		'      context: ./note-web' \
		'      dockerfile: Dockerfile.dev' \
		'    volumes:' \
		'      - ./note-web:/app' \
		'      - /app/node_modules' \
		'      - /app/.next' \
		'    environment:' \
		'      - NODE_ENV=development' \
		'    command: npm run dev' \
		> $(COMPOSE_DEV_FILE)
	@echo "$(GREEN)Development compose file created!$(NC)"

create-air-config: ## Create Air config for Go hot reloading
	@echo "$(BLUE)Creating Air configuration for Go hot reloading...$(NC)"
	@echo 'root = "."' > note-server/.air.toml
	@echo 'tmp_dir = "tmp"' >> note-server/.air.toml
	@echo '' >> note-server/.air.toml
	@echo '[build]' >> note-server/.air.toml
	@echo '  cmd = "go build -o ./tmp/main ./cmd/server"' >> note-server/.air.toml
	@echo '  bin = "./tmp/main"' >> note-server/.air.toml
	@echo '  include_ext = ["go", "tpl", "tmpl", "html"]' >> note-server/.air.toml
	@echo '  exclude_dir = ["assets", "tmp", "vendor", "testdata"]' >> note-server/.air.toml
	@echo "$(GREEN)Air configuration created!$(NC)"

setup-dev: install create-dev-compose create-air-config create-dev-dockerfiles ## Setup complete development environment
	@echo "$(GREEN)Development environment setup complete!$(NC)"
	@echo "$(YELLOW)Run 'make dev' to start development with hot reloading$(NC)"

create-dev-dockerfiles: ## Create development Dockerfiles
	@echo "$(BLUE)Creating development Dockerfiles...$(NC)"
	@echo 'FROM golang:1.24-alpine AS dev' > note-server/Dockerfile.dev
	@echo '' >> note-server/Dockerfile.dev
	@echo '# Install air for hot reloading' >> note-server/Dockerfile.dev
	@echo 'RUN go install github.com/cosmtrek/air@latest' >> note-server/Dockerfile.dev
	@echo '' >> note-server/Dockerfile.dev
	@echo '# Install other dependencies' >> note-server/Dockerfile.dev
	@echo 'RUN apk add --no-cache ffmpeg' >> note-server/Dockerfile.dev
	@echo '' >> note-server/Dockerfile.dev
	@echo 'WORKDIR /app' >> note-server/Dockerfile.dev
	@echo '' >> note-server/Dockerfile.dev
	@echo '# Copy go mod files' >> note-server/Dockerfile.dev
	@echo 'COPY go.mod go.sum ./' >> note-server/Dockerfile.dev
	@echo 'RUN go mod download' >> note-server/Dockerfile.dev
	@echo '' >> note-server/Dockerfile.dev
	@echo '# Copy source code' >> note-server/Dockerfile.dev
	@echo 'COPY . .' >> note-server/Dockerfile.dev
	@echo '' >> note-server/Dockerfile.dev
	@echo '# Create tmp directory for air' >> note-server/Dockerfile.dev
	@echo 'RUN mkdir -p tmp' >> note-server/Dockerfile.dev
	@echo '' >> note-server/Dockerfile.dev
	@echo 'EXPOSE 8080' >> note-server/Dockerfile.dev
	@echo '' >> note-server/Dockerfile.dev
	@echo '# Use air for hot reloading' >> note-server/Dockerfile.dev
	@echo 'CMD ["air", "-c", ".air.toml"]' >> note-server/Dockerfile.dev
	@echo 'FROM node:20-alpine AS dev' > note-web/Dockerfile.dev
	@echo '' >> note-web/Dockerfile.dev
	@echo 'WORKDIR /app' >> note-web/Dockerfile.dev
	@echo '' >> note-web/Dockerfile.dev
	@echo '# Copy package files' >> note-web/Dockerfile.dev
	@echo 'COPY package*.json ./' >> note-web/Dockerfile.dev
	@echo '' >> note-web/Dockerfile.dev
	@echo '# Install dependencies' >> note-web/Dockerfile.dev
	@echo 'RUN npm install' >> note-web/Dockerfile.dev
	@echo '' >> note-web/Dockerfile.dev
	@echo '# Copy source code' >> note-web/Dockerfile.dev
	@echo 'COPY . .' >> note-web/Dockerfile.dev
	@echo '' >> note-web/Dockerfile.dev
	@echo 'EXPOSE 3000' >> note-web/Dockerfile.dev
	@echo '' >> note-web/Dockerfile.dev
	@echo '# Use npm run dev for hot reloading' >> note-web/Dockerfile.dev
	@echo 'CMD ["npm", "run", "dev"]' >> note-web/Dockerfile.dev
	@echo "$(GREEN)Development Dockerfiles created!$(NC)"

# Quick development commands
quick-start: build up ## Quick start: build and run
	@echo "$(GREEN)Quick start completed!$(NC)"

quick-dev: dev-build ## Quick dev: build and run in development mode

# Health check
health: ## Check if services are healthy
	@echo "$(BLUE)Checking service health...$(NC)"
	@curl -sf http://localhost:8080/healthz > /dev/null && echo "$(GREEN)✓ Server is healthy$(NC)" || echo "$(RED)✗ Server is not responding$(NC)"
	@curl -sf http://localhost:3000 > /dev/null && echo "$(GREEN)✓ Web app is healthy$(NC)" || echo "$(RED)✗ Web app is not responding$(NC)"
