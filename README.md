# Template Go Echo Squirrel

Go backend template with Echo + Squirrel + pgx + Postgres.

## Features

- Session-based auth (scs + pgxstore)
- CSRF protection
- Rate limiting
- Password strength validation
- Swagger UI
- Docker Compose (Postgres + PgAdmin)

## Quick Start

1. cp .env.example .env
2. make db-start
3. make migrate-up
4. make swagger
5. make dev

## Endpoints

- GET /health
- POST /api/v1/auth/register
- POST /api/v1/auth/login
- POST /api/v1/auth/logout
- GET /api/v1/auth/me
- GET /api/v1/users
- GET /api/v1/users/:id

## URLs

- API: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html
- PgAdmin: http://localhost:5051

## License

MIT
