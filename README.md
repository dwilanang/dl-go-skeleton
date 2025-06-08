# 🧾 Payroll & Payslip Management System

A full-featured payroll and payslip management application built with **Golang** and **PostgreSQL**, designed to automate salary calculations, attendance tracking, overtime, and reimbursement management.

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

---

## 📖 About

This system allows **admins** to manage employee salary components, generate payrolls based on attendance, overtime, and reimbursement data, and produce itemized payslips. It enforces role-based access, JWT authentication, and supports Swagger for API documentation.

---

## 🚀 Features

- ✅ JWT Authentication with role-based access (`admin`, `employee`)
- ✅ CRUD for Users and Salaries
- ✅ Attendance (daily records)
- ✅ Overtime tracking with hourly constraints
- ✅ Reimbursement management with optional descriptions
- ✅ Payroll period management (`pending`, `processed`)
- ✅ Automatic payroll calculation per employee per period
- ✅ Detailed payslip generation in JSON
- ✅ Swagger API documentation
- ✅ Audit timestamps (created_at, updated_at)
- ✅ Database seeding with fake users for development

---

## 🛠️ Tech Stack

- **Backend**: Go (Golang)
- **Database**: PostgreSQL
- **ORM**: GORM
- **Migration Tool**: Goose
- **Router**: Gin
- **Authentication**: JWT
- **Documentation**: Swagger (via `swaggo/gin-swagger`)
- **Testing**: Go test
- **Seeding**: Custom Goose script with Faker

---

## 🗂️ Project Structure

```
psp/
├── cmd/                # Application entry points (main.go)
│   └── api/
│       └── main.go
├── config/             # Configuration files and utilities
│   └── config.go
├── db/
│   └── migrations/     # Database migration scripts
├── docs/               # Swagger docs (auto-generated)
│   └── swagger.yaml
├── infrastructure/
│   └── db/
│       └── postgres/   # Postgres connection and helpers
├── internal/
│   ├── admin/          # Admin domain (service, handler, repository, dto, model)
│   ├── auditlogs/      # Audit logs domain
│   ├── auth/           # Authentication (JWT, login, middleware)
│   ├── employee/       # Employee domain
│   ├── middleware/     # Custom Gin middleware
│   ├── registry/       # Dependency injection registry
│   ├── role/           # Role management
│   ├── user/           # User management
│   └── ...             # Other business domains
├── utils/              # Utility functions (request, response, etc)
├── go.mod
├── go.sum
└── README.md
```

**Key folders:**
- `cmd/`: Entry point for the application.
- `config/`: Application configuration loader.
- `db/migrations/`: Goose migration scripts.
- `docs/`: Swagger/OpenAPI documentation.
- `infrastructure/`: Database connection and low-level utilities.
- `internal/`: Main business logic, organized per domain (admin, user, auth, etc).
- `utils/`: Shared utility functions.

---

## 🧱 Database Structure

Core tables include:

- `users`: Master user data with role control
- `user_salaries`: Salary history with effective dates
- `attendances`: Daily attendance per user
- `overtimes`: Overtime record (limited to 3 hours/day)
- `reimbursements`: Expense claims per period
- `attendance_periods`: Defines salary period range and status
- `payrolls`: Processed payrolls per user per period
- `payslips`: JSON-based salary breakdown

Migration scripts are stored in `/db/migrations`.

---

## ⚙️ Installation

```bash
# Clone the project
git clone https://github.com/dwilanang/hr-ppms.git
cd hr-ppms

# Setup env (edit as needed)
cp .env.example .env

go install github.com/pressly/goose/v3/cmd/goose@latest

# Setup DB
createdb ppms_db

make migrate-up

# Run application
go run cmd/api/main.go
```

---

## 🚦 Usage

- Access the API at `http://localhost:8080/api/v1`
- Use Swagger UI for API documentation and testing endpoints.

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

---

## 📄 License

MIT License. See [LICENSE](LICENSE) for details.