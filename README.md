# Template Go Echo Squirrel

[![CI](https://github.com/ssr0016/template-go-echo-squirrel/actions/workflows/ci.yml/badge.svg)](https://github.com/ssr0016/template-go-echo-squirrel/actions/workflows/ci.yml)

Production-ready Go backend starter kit with **Echo**, **Squirrel**, **pgx**, and **Postgres**.

## 🌐 Live Demo

- **API:** https://template-go-echo-squirrel-production.up.railway.app
- **Swagger:** https://template-go-echo-squirrel-production.up.railway.app/swagger/index.html
- **Health:** https://template-go-echo-squirrel-production.up.railway.app/health

## Features

- Echo v4 - Fast web framework
- PostgreSQL + pgx - High-performance driver
- Squirrel - Type-safe SQL builder
- Goose - DB migrations
- Session auth - scs + Postgres store
- CSRF protection - Double-submit cookie
- Rate limiting - Per-IP limiter
- Password strength - 60-bit entropy min
- Structured logging - slog (JSON/text)
- Standardized errors - apperror package
- Type-safe config - Centralized loader
- Swagger UI - Interactive docs
- Docker - Multi-stage production build
- Unit tests - 20 tests, 3 packages
- CI/CD - GitHub Actions
- Air - Hot reload

## Prerequisites

- Go 1.26+
- Docker & Docker Compose
- Make
- Dev tools: swag, goose, air

See [SETUP.md](SETUP.md) for detailed installation.

## Quick Start

git clone https://github.com/ssr0016/template-go-echo-squirrel.git
cd template-go-echo-squirrel
cp .env.example .env
make db-start
go mod download
make migrate-up
make swagger
make dev

## API Endpoints

### Public
- GET /health
- POST /api/v1/auth/register
- POST /api/v1/auth/login

### Protected
- POST /api/v1/auth/logout
- GET /api/v1/auth/me
- GET /api/v1/users
- GET /api/v1/users/:id

## Makefile Commands

make help           # Show all
make dev            # Hot reload
make run            # Direct run
make build          # Build binary
make test           # Run tests
make lint           # Run linter
make swagger        # Generate Swagger
make db-start       # Start DB
make db-stop        # Stop DB
make db-reset       # Wipe + restart
make migrate-up     # Apply migrations
make migrate-down   # Rollback
make migrate-status # Check status

## Project Structure

template-go-echo-squirrel/
- cmd/api/main.go
- internal/
  - apperror/     # Error handling
  - config/       # Config loader
  - database/     # DB + migrations
  - handler/      # HTTP handlers
  - logger/       # slog wrapper
  - middleware/   # Custom middleware
  - model/        # Domain models
  - repository/   # Data access
  - router/       # Routes
  - service/      # Business logic
  - session/      # Session config
  - validator/    # Request validation
- docs/             # Swagger
- Dockerfile
- docker-compose.yaml
- docker-compose.prod.yaml
- Makefile
- .env.example
- README.md
- SETUP.md
- ROADMAP.md
- CUSTOMIZATION.md

## Security

- Session-based auth (Postgres store)
- CSRF protection (double-submit cookie)
- Rate limiting (5/min on auth)
- Password strength (60-bit entropy)
- HttpOnly + SameSite=Strict cookies
- Bcrypt password hashing
- Session renewal on login
- Security scan (gosec) in CI

## Docker Services

| Service | Port | Purpose |
|---|---|---|
| pg | 5435 | Postgres (main) |
| test-pg | 5436 | Postgres (test) |
| pg-admin | 5051 | PgAdmin UI |

## Deployment

Deployed on [Railway](https://railway.app) with:
- Auto-deploy from GitHub
- Postgres database
- SSL/TLS
- Custom domain ready

## Documentation

- [SETUP.md](SETUP.md) - Setup guide
- [ROADMAP.md](ROADMAP.md) - Project roadmap
- [CUSTOMIZATION.md](CUSTOMIZATION.md) - Customization guide

## License

MIT
