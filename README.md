# Template Go Echo Squirrel

[![CI](https://github.com/ssr0016/template-go-echo-squirrel/actions/workflows/ci.yml/badge.svg)](https://github.com/ssr0016/template-go-echo-squirrel/actions/workflows/ci.yml)

Production-ready Go backend starter kit with Echo, Squirrel, pgx, Postgres, session auth, and RBAC.

## Live Demo

- API: https://template-go-echo-squirrel-production.up.railway.app
- Swagger: https://template-go-echo-squirrel-production.up.railway.app/swagger/index.html
- Health: https://template-go-echo-squirrel-production.up.railway.app/health
- Metrics: https://template-go-echo-squirrel-production.up.railway.app/metrics

## Features

### Core
- Echo v4 - Fast, minimalist web framework
- PostgreSQL + pgx - High-performance database driver
- Squirrel - Type-safe SQL query builder
- Goose - Database migrations
- Type-safe config - Centralized config loader
- Structured logging - log/slog (JSON/text)
- Standardized errors - Custom apperror package

### Authentication and Authorization
- Session-based auth - alexedwards/scs + Postgres store
- CSRF protection - Double-submit cookie pattern
- Rate limiting - Per-IP rate limiter
- Password strength - 60-bit entropy minimum
- RBAC - Roles + Permissions (dynamic)
- Role middleware - RequireRole, RequirePermission

### API Features
- Swagger UI - Interactive API docs
- Pagination - Standard data + meta format
- Filtering - Query param filters
- Health checks - /health, /ready, /live
- Prometheus metrics - /metrics endpoint

### Development
- Docker - Multi-stage production build
- Air - Hot reload for development
- Unit tests - 20 tests, 3 packages
- CI/CD - GitHub Actions (lint, test, build, security)
- Documentation - README, SETUP, ROADMAP, ARCHITECTURE

## Prerequisites

- Go 1.26+
- Docker and Docker Compose
- Make
- Dev tools: swag, goose, air

See SETUP.md for detailed installation.

## Quick Start

git clone https://github.com/ssr0016/template-go-echo-squirrel.git
cd template-go-echo-squirrel
cp .env.example .env
nano .env
make db-start
go mod download
make migrate-up
make swagger
make dev

Expected output:

    API Server:  http://localhost:8080
    Swagger UI:  http://localhost:8080/swagger/index.html
    Health:      http://localhost:8080/health
    Metrics:     http://localhost:8080/metrics

## API Endpoints

### Public
| Method | Path | Description |
|---|---|---|
| GET | /health | Health check |
| GET | /ready | Readiness (DB ping) |
| GET | /live | Liveness |
| GET | /metrics | Prometheus metrics |
| POST | /api/v1/auth/register | Register user |
| POST | /api/v1/auth/login | Login user |

### Authenticated (session required)
| Method | Path | Description |
|---|---|---|
| POST | /api/v1/auth/logout | Logout |
| GET | /api/v1/auth/me | Current user |
| GET | /api/v1/users | List users (paginated) |
| GET | /api/v1/users/:id | Get user |

### Admin (role: admin)
| Method | Path | Description |
|---|---|---|
| GET | /api/v1/admin/roles | List roles (paginated) |
| POST | /api/v1/admin/roles | Create role |
| GET | /api/v1/admin/roles/:id | Get role |
| PUT | /api/v1/admin/roles/:id | Update role |
| DELETE | /api/v1/admin/roles/:id | Delete role |
| POST | /api/v1/admin/roles/:id/permissions | Assign permission |
| DELETE | /api/v1/admin/roles/:id/permissions/:pid | Revoke permission |
| GET | /api/v1/admin/permissions | List permissions (paginated) |
| POST | /api/v1/admin/permissions | Create permission |
| GET | /api/v1/admin/users/:id | Get user with role |
| PUT | /api/v1/admin/users/:id/role | Change user role |

## RBAC (Role-Based Access Control)

### Default Roles
| Role | Permissions |
|---|---|
| admin | All (9 permissions) |
| editor | users:read, users:write |
| user | users:read |

### Default Admin
- Email: admin@example.com
- Password: AdminPass123!

### Permission Format
resource:action - e.g., users:read, roles:write

## Pagination

List endpoints support pagination via query params:

    GET /api/v1/users?page=1&limit=10
    GET /api/v1/admin/roles?page=1&limit=20

Response format:

    {
      "data": [...],
      "meta": {
        "page": 1,
        "limit": 10,
        "total": 100,
        "total_pages": 10
      }
    }

Limits:
- Default page: 1
- Default limit: 20
- Max limit: 100

## Observability

### Health Checks
| Endpoint | Purpose |
|---|---|
| /health | Basic health |
| /ready | Readiness (DB ping) |
| /live | Liveness |

### Metrics (Prometheus)

    curl http://localhost:8080/metrics

Available metrics:
- http_requests_total - Counter
- http_request_duration_seconds - Histogram
- http_requests_in_flight - Gauge
- db_connections_active - Gauge
- db_connections_idle - Gauge

## Makefile Commands

    make help           Show all
    make dev            Hot reload
    make run            Direct run
    make build          Build binary
    make test           Run tests
    make lint           Run linter
    make swagger        Generate Swagger
    make db-start       Start DB
    make db-stop        Stop DB
    make db-reset       Wipe + restart
    make migrate-up     Apply migrations
    make migrate-down   Rollback
    make migrate-status Check status

## Project Structure

    template-go-echo-squirrel/
    ├── cmd/api/main.go
    ├── internal/
    │   ├── apperror/     Error handling
    │   ├── config/       Config loader
    │   ├── database/     DB + migrations + seed
    │   ├── handler/      HTTP handlers
    │   ├── logger/       slog wrapper
    │   ├── middleware/   Middleware (auth, CSRF, RBAC, metrics)
    │   ├── model/        Domain models
    │   ├── repository/   Data access
    │   ├── router/       Routes
    │   ├── service/      Business logic
    │   ├── session/      Session config
    │   └── validator/    Custom validation
    ├── pkg/
    │   ├── metrics/      Prometheus metrics
    │   └── pagination/   Pagination helpers
    ├── docs/             Swagger + guides
    ├── Dockerfile
    ├── docker-compose.yaml
    ├── docker-compose.prod.yaml
    ├── Makefile
    ├── .env.example
    └── README.md

## Security

- Session-based auth (Postgres store)
- CSRF protection (double-submit cookie)
- Rate limiting (5/min on auth)
- Password strength (60-bit entropy)
- HttpOnly + SameSite=Lax cookies
- Bcrypt password hashing
- Session renewal on login
- RBAC (roles + permissions)
- Security scan (gosec) in CI

## Docker Services

| Service | Port | Purpose |
|---|---|---|
| pg | 5435 | Postgres (main) |
| test-pg | 5436 | Postgres (test) |
| pg-admin | 5051 | PgAdmin UI |

PgAdmin: http://localhost:5051 (admin@example.com / admin)

## Deployment

Deployed on Railway (https://railway.app) with:
- Auto-deploy from GitHub
- Postgres database
- SSL/TLS
- Custom domain ready

## Documentation

- SETUP.md - Setup guide
- ROADMAP.md - Project roadmap
- CUSTOMIZATION.md - Customization guide
- docs/ARCHITECTURE.md - Architecture rules
- Swagger UI - http://localhost:8080/swagger/index.html

## License

MIT
