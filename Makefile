include ./.dev/misc/make/tools.Makefile

# Konfigurasi DSN database PostgreSQL
DB_DSN=postgres://ppms_user:ppms123@localhost:5432/ppms?sslmode=disable
# Konfigurasi path migrations
MIGRATIONS=db/migrations

docs-init:
	@swag init -g cmd/api/main.go
run:
	@go run cmd/api/main.go
migrate-up:
	@goose -dir $(MIGRATIONS) postgres "$(DB_DSN)" up
migrate-down:
	@goose -dir $(MIGRATIONS) postgres "$(DB_DSN)" down
migrate-down-n:
	@goose -dir $(MIGRATIONS) postgres "$(DB_DSN)" down $(n)
migrate-reset:
	@goose -dir $(MIGRATIONS) postgres "$(DB_DSN)" reset

service-up:
	@docker compose -f compose.yaml up

service-down:
	@docker compose -f compose.yaml down

clean: clean-artifacts clean-docker

clean-docker: ## Removes dangling docker images
	@ docker image prune -f