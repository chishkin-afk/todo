# Todo List Service

A robust, production-ready Todo List API built with Go, following Clean Architecture principles. It features user authentication, task management, secure TLS communication, and persistent storage with PostgreSQL and Redis.

## 🚀 Features

- **User Authentication**: Secure registration and login using JWT sessions.
- **Task Management**: Create, read, update, and delete tasks and task groups.
- **Clean Architecture**: Strictly separated layers (Domain, Application, Infrastructure) for maintainability and testability.
- **Security**: 
  - Password hashing (bcrypt).
  - HTTPS/TLS support out of the box.
  - Input validation using `go-playground/validator`.
- **High Performance**: 
  - Caching layer with Redis for frequent data access.
  - Efficient PostgreSQL connection pooling.
- **Dockerized**: Ready-to-run with `docker-compose`, including health checks and data persistence.
- **Graceful Shutdown**: Handles SIGINT/SIGTERM signals to complete active requests before stopping.

## 🛠 Tech Stack

- **Language**: Go (1.21+)
- **Framework**: Gin (HTTP Server)
- **Database**: PostgreSQL (via `lib/pq`)
- **Cache**: Redis (via `go-redis/v9`)
- **Validation**: go-playground/validator
- **Configuration**: YAML + Environment Variables (`godotenv`)
- **Logging**: Standard `log/slog`
- **Testing**: `stretchr/testify`
- **Infrastructure**: Docker, Docker Compose, Make

## 📂 Project Structure

```
├── cmd/                      # Application entry point (main.go)
├── configs/                  # YAML configuration files
├── certs/                    # TLS certificates (server.crt, server.key)
├── internal/                 # Private application logic
│   ├── app/                  # App assembly, wiring, and management
│   ├── application/          # Use cases, DTOs (requests/responses)
│   ├── common/               # Shared utilities (config loading)
│   └── infrastructure/       # External adapters
│       ├── cache/redis/      # Redis connection and implementation
│       ├── http/             # Gin server, handlers, routes, middlewares
│       ├── persistence/postgres/ # DB connection, migrations
│       └── session/jwt/      # JWT token generation and validation
├── modules/                  # Business domains (Feature-based)
│   ├── auth/                 # Authentication module
│   │   ├── application/      # Auth services
│   │   ├── domain/           # Entities (User, Session), interfaces
│   │   └── infrastructure/   # DB/Cache repos for Auth
│   └── task/                 # Task management module
│       ├── application/      # Task services
│       ├── domain/           # Entities (Task, Group, Priority)
│       └── infrastructure/   # DB/Cache repos for Tasks
├── pkg/                      # Public library code
│   ├── consts/               # Application constants
│   ├── errors/               # Custom error types
│   └── log/                  # Logging utilities
├── docker-compose.yaml       # Container orchestration
├── Dockerfile                # App container image
├── Makefile                  # Automation commands
└── .env                      # Environment variables
```


## ⚡ Quick Start

Prerequisites:
- Docker & Docker Compose
- Make
- (Optional) `mkcert` for local TLS generation

### 1. Run the Service

Simply run the following command to set up environment variables, build images, and start all services:

```bash
make quick
```

This command will:
1. Copy `.env.example` to `.env` (if it exists).
2. Start PostgreSQL, Redis, and the Todo service via Docker Compose.

The API will be available at `https://localhost:9000` (note the HTTPS).

### 2. Generate Local TLS Certificates (Optional)

If you need to regenerate self-signed certificates for local development:

```bash
make localtls
```

## ⚙️ Configuration

The application is configured via `configs/local.yml` and environment variables defined in `.env`.

Key environment variables:
- `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DBNAME`
- `REDIS_ADDR`, `REDIS_PASSWORD`
- `HTTP_ADDR` (e.g., `:9000`)
- `SECRET_KEY` (for JWT signing)
- `APP_ENV` (local, prod)

## 🧪 Testing

Run tests using the standard Go tooling:

```bash
go test ./...
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.