# API Documentation

Base URL: http://localhost:8080/api/v1

## Authentication Flow

### 1. Register

POST /auth/register

    {
      "email": "user@example.com",
      "name": "User Name",
      "password": "StrongPass123!"
    }

Response (201):

    {
      "id": 1,
      "email": "user@example.com",
      "name": "User Name",
      "role_id": 1,
      "created_at": "2026-01-01T00:00:00Z"
    }

### 2. Login

POST /auth/login

    {
      "email": "user@example.com",
      "password": "StrongPass123!"
    }

Response (200): Same as register. Sets app_session cookie.

### 3. Authenticated Requests

Include session cookie with each request:

    curl -b cookies.txt http://localhost:8080/api/v1/auth/me

### 4. Logout

POST /auth/logout

Response (200):

    {"message": "logged out"}

## Users

### List Users (paginated)

GET /users?page=1&limit=10&email=user

Response (200):

    {
      "data": [...],
      "meta": {
        "page": 1,
        "limit": 10,
        "total": 42,
        "total_pages": 5
      }
    }

### Get User

GET /users/:id

## Admin Endpoints

Requires role: admin.

### List Roles (paginated)

GET /admin/roles?page=1&limit=10

### Create Role

POST /admin/roles

    {
      "name": "manager",
      "description": "Manager role"
    }

### Assign Permission

POST /admin/roles/:id/permissions

    {
      "permission_id": 1
    }

### List Permissions (paginated)

GET /admin/permissions?page=1&limit=10&resource=users

### Change User Role

PUT /admin/users/:id/role

    {
      "role_id": 2
    }

## Error Responses

### 400 Bad Request

    {"error": "BAD_REQUEST", "message": "invalid request body"}

### 401 Unauthorized

    {"error": "UNAUTHORIZED", "message": "authentication required"}

### 403 Forbidden

    {"error": "FORBIDDEN", "message": "insufficient permissions"}

### 404 Not Found

    {"error": "NOT_FOUND", "message": "resource not found"}

### 409 Conflict

    {"error": "CONFLICT", "message": "email already taken"}

### 422 Validation Error

    {
      "error": "VALIDATION_ERROR",
      "message": "Key: 'RegisterRequest.Email' Error:Field validation..."
    }

### 429 Rate Limit

    {"error": "RATE_LIMIT_EXCEEDED", "message": "rate limit exceeded"}

## CSRF Protection

For POST/PUT/DELETE requests, include CSRF token:

1. GET /health - receives csrf_token cookie
2. Include header: X-CSRF-Token: <token>

Example:

    curl -c cookies.txt http://localhost:8080/health > /dev/null
    CSRF=$(grep csrf_token cookies.txt | awk '{print $NF}')

    curl -X POST http://localhost:8080/api/v1/auth/login \
      -H "Content-Type: application/json" \
      -H "X-CSRF-Token: $CSRF" \
      -b cookies.txt -c cookies.txt \
      -d '{"email": "admin@example.com", "password": "AdminPass123!"}'

## Complete Example

Full auth flow:

    # 1. Get CSRF token
    curl -c cookies.txt http://localhost:8080/health > /dev/null
    CSRF=$(grep csrf_token cookies.txt | awk '{print $NF}')

    # 2. Register
    curl -X POST http://localhost:8080/api/v1/auth/register \
      -H "Content-Type: application/json" \
      -H "X-CSRF-Token: $CSRF" \
      -b cookies.txt -c cookies.txt \
      -d '{"email": "user@example.com", "name": "User", "password": "StrongPass123!"}'

    # 3. Login
    curl -X POST http://localhost:8080/api/v1/auth/login \
      -H "Content-Type: application/json" \
      -H "X-CSRF-Token: $CSRF" \
      -b cookies.txt -c cookies.txt \
      -d '{"email": "user@example.com", "password": "StrongPass123!"}'

    # 4. Get current user
    curl -b cookies.txt http://localhost:8080/api/v1/auth/me

    # 5. List users (paginated)
    curl -b cookies.txt "http://localhost:8080/api/v1/users?page=1&limit=10"
