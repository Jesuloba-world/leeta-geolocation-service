# Database migration commands using Goose

.PHONY: migrate-up migrate-down migrate-status migrate-reset migrate-version docker-up docker-down docker-logs

# Default database URL (can be overridden)
DATABASE_URL ?= postgres://leeta_user:leeta_password@localhost:5432/leeta_db?sslmode=disable
MIGRATIONS_DIR = scripts/migrations

# Docker commands
docker-up: ## Start PostgreSQL container
	docker compose up -d

docker-down: ## Stop PostgreSQL container
	docker compose down

docker-logs: ## View PostgreSQL container logs
	docker compose logs -f postgres

docker-clean: ## Remove PostgreSQL container and volumes
	docker compose down -v

docker-rebuild: ## Rebuild and restart containers
	docker compose down && docker compose up -d --force-recreate

migrate-up:
	@echo "Running database migrations..."
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up

migrate-down:
	@echo "Rolling back last migration..."
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" down

migrate-status:
	@echo "Checking migration status..."
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" status

migrate-reset:
	@echo "Resetting all migrations..."
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" reset

migrate-version:
	@echo "Checking current migration version..."
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" version

migrate-create: ## Create new migration (usage: make migrate-create name=migration_name)
	@$(if $(name),,$(error Usage: make migrate-create name=migration_name))
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

help:
	@echo "Available commands:"
	@echo ""
	@echo "Docker commands:"
	@echo "  make docker-up       - Start PostgreSQL with PostGIS"
	@echo "  make docker-down     - Stop PostgreSQL container"
	@echo "  make docker-logs     - Show PostgreSQL logs"
	@echo "  make docker-clean    - Remove container and volumes"
	@echo "  make docker-rebuild  - Rebuild and restart containers"
	@echo ""
	@echo "Migration commands:"
	@echo "  make migrate-up      - Run all pending migrations"
	@echo "  make migrate-down    - Rollback the last migration"
	@echo "  make migrate-status  - Show migration status"
	@echo "  make migrate-reset   - Reset all migrations"
	@echo "  make migrate-version - Show current version"
	@echo "  make migrate-create  - Create a new migration"
	@echo ""
	@echo "Set DATABASE_URL environment variable or override:"
	@echo "  make migrate-up DATABASE_URL='your_db_url'"