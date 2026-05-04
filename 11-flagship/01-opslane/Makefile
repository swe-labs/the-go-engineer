# Opslane Makefile
# Common development and operations commands

# Colors for output
RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[0;33m
NC = \033[0m # No Color

.PHONY: help build up down restart logs health test clean db-reset

# Default target
help:
	@echo "$(GREEN)Opslane Makefile$(NC)"
	@echo "======================"
	@echo ""
	@echo "Available commands:"
	@echo "  make build        - Build the Docker image"
	@echo "  make up           - Start all services (db + api)"
	@echo "  make down         - Stop all services"
	@echo "  make restart      - Restart all services"
	@echo "  make logs        - View API logs (follow mode)"
	@echo "  make db-logs     - View database logs"
	@echo "  make health      - Check service health"
	@echo "  make test        - Run tests inside container"
	@echo "  make shell       - Get shell inside API container"
	@echo "  make db-reset    - Reset database (danger!)"
	@echo "  make clean       - Remove all containers, volumes, images"
	@echo ""
	@echo "Quick access:"
	@echo "  make up-detached - Start in detached mode"
	@echo "  make ps          - Show running containers"
	@echo ""

build:
	@echo "$(YELLOW)Building Opslane...$(NC)"
	docker compose build --no-cache

up:
	@echo "$(YELLOW)Starting Opslane...$(NC)"
	docker compose up -d
	@echo "$(GREEN)Waiting for services to be healthy...$(NC)"
	@sleep 3
	@docker compose ps

up-detached:
	@echo "$(YELLOW)Starting Opslane in background...$(NC)"
	docker compose up -d --build

down:
	@echo "$(YELLOW)Stopping Opslane...$(NC)"
	docker compose down

restart: down up

logs:
	docker compose logs -f api

db-logs:
	docker compose logs -f db

health:
	@echo "$(YELLOW)Checking health...$(NC)"
	@curl -s http://localhost:8080/health && echo "$(GREEN) - API OK$(NC)" || echo "$(RED) - API FAILED$(NC)"
	@docker compose ps

test:
	@echo "$(YELLOW)Running tests...$(NC)"
	docker compose exec api go test -v ./...

shell:
	docker compose exec api sh

ps:
	docker compose ps

db-reset:
	@echo "$(RED)WARNING: This will delete all data!$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "$(YELLOW)Resetting database...$(NC)"; \
		docker compose down -v; \
		docker compose up -d; \
	else \
		echo "Aborted."; \
	fi

clean:
	@echo "$(YELLOW)Cleaning up Docker resources...$(NC)"
	docker compose down -v --rmi local
	@echo "$(GREEN)Done!$(NC)"