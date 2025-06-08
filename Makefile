# Konfigurasi DSN database PostgreSQL
DB_DSN=postgres://postgres:dl1625@localhost:5432/payslip_db?sslmode=disable
# Konfigurasi path migrations
MIGRATIONS=infrastructure/db/postgres/migrations

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