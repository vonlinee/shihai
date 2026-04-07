# Shihai Backend Platform

This is the backend implementation for the Shihai (识海) ancient poetry learning platform. The project is built using Go (Golang) with the Gin web framework, GORM for database operations, and PostgreSQL.

## Directory Structure

The backend follows Standard Go Project Layout principles, carefully separating business logic, HTTP routing, and database interactions into a clear layered architecture.

```text
backend/
├── cmd/                        # Main applications for this project
│   └── server/                 # Contains the application entry point (main.go)
│
├── internal/                   # Private application and library code
│   ├── config/                 # Application configuration schemas and database initialization
│   ├── dto/                    # Data Transfer Objects (Data structures for API input/output validation)
│   ├── handlers/               # HTTP Handlers (Controllers) parsing requests and returning responses
│   ├── middleware/             # Gin Middleware (Authentication, RBAC checks, CORS)
│   ├── models/                 # Database schema definitions and Domain Models
│   ├── repository/             # Data Access Layer (Database querying logic)
│   └── services/               # Core Business Logic Layer
│
├── pkg/                        # Library code that can be exported or reused
│   └── utils/                  # Broad utility helpers (JWT tokens, password hashing, snowflakes, etc.)
│
├── .idea/                      # JetBrains IDE configuration files (local)
├── bin/                        # Compiled executables
├── Golang开发规范.md           # Go development standards and guidelines
├── config.json                 # Active configuration properties (local)
├── config.example.json         # Example structure for application configuration
├── go.mod                      # Go module dependencies file
└── go.sum                      # Go module checksums file
```

## Architecture Layers Overview

- **Handler Layer (`internal/handlers`)**: The entry point for API requests. It relies on the Gin framework to parse parameters, validates inputs via `dto`, passes logic over to the `services` layer, and formats final JSON responses. 
- **Service Layer (`internal/services`)**: The brain of the API. All business logic, system orchestration, and transaction flows happen here. It sits between `handlers` and `repositories` to keep logic decoupled from HTTP properties and database flavors.
- **Repository Layer (`internal/repository`)**: Isolates database communication. The repository runs raw DB/GORM commands strictly to load and transform data into `models`.

## Setting up and Running

1. **Install Dependencies**
   ```bash
   go mod download
   ```

2. **Configure Database**
   Copy `config.example.json` to `config.json` and configure the database credentials to match your PostgreSQL instance.

3. **Start Application**
   ```bash
   go run cmd/server/main.go
   ```
   *Note: GORM auto migrations are executed upon starting to ensure your database perfectly replicates the definitions inside `internal/models/`*.
