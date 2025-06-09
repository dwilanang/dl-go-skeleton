# 🧾 Go Skeleton API

A robust API with **Golang** and **PostgreSQL**.
---

## 📌 Table of Contents

- [About](#about)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Database Structure](#database-structure)
- [Installation](#installation)
- [Usage](#usage)
- [API Documentation](#api-documentation)
- [Seeding Fake Data](#seeding-fake-data)
- [Development Notes](#development-notes)
- [License](#license)

## 🚀 Features

- ✅ JWT Authentication with role-based access (`SUPERADMIN`, `ADMIN`, `EMPLOYEE`)
- ✅ Swagger API Documentation
- ✅ Database Seeding with Fake Users for Development

---

## 🛠️ Tech Stack

- **Backend**: Go (Golang)
- **Database**: PostgreSQL
- **ORM**: SQLX
- **Migration Tool**: Goose
- **Router**: Gin
- **Authentication**: JWT
- **Documentation**: Swagger (`swaggo/gin-swagger`)
- **Testing**: Go test, sqlmock, testify

---

## 🗂️ Project Structure

```
psp/
├── cmd/                # Application entry points (main.go)
│   └── api/
│       └── main.go
├── config/             # Configuration loader
│   └── config.go
├── db/
│   └── migrations/     # Database migration scripts
├── docs/               # Swagger docs (auto-generated)
│   └── swagger.yaml
├── infrastructure/
│   └── db/
│       └── postgres/   # Postgres connection and helpers
├── internal/
│   ├── auth/           # Authentication (JWT, login, middleware)
│   ├── middleware/     # Custom Gin middleware
│   ├── registry/       # Dependency injection registry
│   ├── role/           # Role management
│   ├── user/           # User management
│   └── ...             # Other business domains
├── pkg/
│   └── logger/         # Logging utilities
├── utils/              # Utility functions (request, response, etc)
├── go.mod
├── go.sum
└── README.md
```

**Key folders:**
- `cmd/`: Application entry point.
- `config/`: Loads application configuration.
- `db/migrations/`: Goose migration scripts.
- `docs/`: Swagger/OpenAPI documentation.
- `infrastructure/`: Database connection and low-level utilities.
- `internal/`: Main business logic, organized per domain (auth, user, role, etc).
- `pkg/`: Shared packages (e.g., logger).
- `utils/`: Shared utility functions.

---

## ⚙️ Installation

```bash
# Clone the project
git clone https://github.com/dwilanang/go-skeleton.git

cd go-skeleton

# Setup environment variables
cp .env.example .env

# Setup database
# create database

# Install dependencies
go mod tidy

# Run database migrations and seed dummy data
# design database dbdiagram.io > export to postgresql

go install github.com/pressly/goose/v3/cmd/goose@latest

goose -dir db/migrations postgres "postgres://[DB_USER]:[DB_PASSWORD]@[DB_HOST]:[DB_PORT]/[DB_NAME]?sslmode=disable" up

# Generate Swagger docs
swag init -g cmd/api/main.go

# Run the application
go run cmd/api/main.go
```

---

## 🚦 Usage

- Access the API at `http://localhost:8000/api/v1`
- Use Swagger UI for API documentation and endpoint testing at `http://localhost:8000/swagger/index.html`

---

## 📚 API Documentation

Swagger docs are auto-generated.  
After running the app, open:  
`http://localhost:8000/swagger/index.html`

---

## 🌱 Seeding Fake Data

To seed development data, use the provided Goose migration scripts or custom seeder in `/db/migrations`.

---

## 📝 Development Notes

- Use Go modules (`go mod tidy`) to manage dependencies.
- Use `goose` for database migrations.
- Run tests with `go test ./...`
- Update Swagger docs with `swag init` (if using swaggo).
- Clean architecture: business logic is separated by domain in `internal/`.
- Dependency injection is managed via the registry pattern.

---

## 📄 License

MIT License. See [LICENSE](LICENSE) for details.

---

**Contributions, bug reports, and suggestions are welcome!**