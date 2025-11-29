# REST API - Go Learning Project

A REST API built with Go's standard library to learn web development fundamentals, database integration, and API design patterns.

## Overview

This project is a learning-focused REST API that demonstrates:
- Building HTTP servers with Go's standard library (`net/http`)
- PostgreSQL database integration
- Type-safe SQL queries using `sqlc`
- Password hashing and security best practices
- Clean architecture with separation of concerns
- Environment-based configuration
- JWT token generation (prepared for authentication)

## Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL
- sqlc (for regenerating queries)

### Environment Setup
Create a `.env` file in the `rest-api/` directory:
```env
SERVER_PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
ENVIRONMENT=development
LOG_LEVEL=info
```

### Running the Server
```bash
cd rest-api
go run main.go
```

Server will start on the configured port (default: `:8080`)

### Testing Endpoints

**Health Check:**
```bash
curl http://localhost:8080/health
```

**User Registration:**
```bash
curl -X POST http://localhost:8080/user/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "securepassword"
  }'
```

## Project Structure

```
rest-api/
├── main.go                          # Application entry point
├── serverconfig/                    # Server configuration
│   ├── config.go
│   └── README.md                    # → Configuration details
├── dbconfig/                        # Database configuration
│   ├── dbconfig.go
│   └── README.md                    # → Database connection details
├── internal/
│   ├── handlers/                    # HTTP request handlers
│   │   ├── core_handler.go
│   │   ├── health.go
│   │   ├── test.go
│   │   ├── user.go
│   │   └── README.md                # → Handler pattern explained
│   ├── routes/                      # Route definitions
│   │   ├── setup_routes.go
│   │   ├── health_routes.go
│   │   ├── test_routes.go
│   │   ├── user_rotues.go
│   │   └── README.md                # → Routing system explained
│   ├── store/                       # Database layer (sqlc generated)
│   │   ├── db.go
│   │   ├── models.go
│   │   ├── queries.sql.go
│   │   └── README.md                # → sqlc and queries explained
│   ├── dtos/                        # Data Transfer Objects
│   │   ├── request.go
│   │   └── README.md                # → DTOs vs Models explained
│   ├── utils/                       # Utility functions
│   │   ├── passwordutil.go
│   │   ├── errorresponse.go
│   │   ├── successresponse.go
│   │   ├── jwt.go
│   │   └── README.md                # → Utilities explained
│   └── migrations/                  # Database migrations
│       ├── schema.sql               # Database schema
│       └── queries.sql              # SQL queries for sqlc
├── models/                          # Domain models
│   ├── user.go
│   ├── blog.go
│   └── README.md                    # → Domain models explained
├── TODO.md                          # Future improvements
└── README.md                        # This file
```

**📖 Each folder contains a detailed README** - click through to understand how each component works!

## Architecture Overview

### Complete Application Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Application Startup                         │
│                            (main.go)                                │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ├─> 1. Load Configuration
                                  │     serverconfig.LoadConfig()
                                  │     └─> Reads .env file
                                  │         Returns Config struct
                                  │
                                  ├─> 2. Connect to Database
                                  │     dbconfig.ConnectDB(config.DatabaseURL)
                                  │     └─> Opens PostgreSQL connection
                                  │         Verifies with Ping()
                                  │         Returns *sql.DB
                                  │
                                  ├─> 3. Initialize Queries
                                  │     store.New(db)
                                  │     └─> Creates type-safe query interface
                                  │         Returns *store.Queries
                                  │
                                  ├─> 4. Create Handler
                                  │     handlers.NewHandlers(db, queries)
                                  │     └─> Injects dependencies into Handler
                                  │         Returns *handlers.Handler
                                  │
                                  ├─> 5. Setup Router
                                  │     http.NewServeMux()
                                  │     routes.SetupRoutes(mux, handler)
                                  │     └─> Registers all endpoints
                                  │
                                  └─> 6. Start Server
                                        server.ListenAndServe()
                                        └─> Listens on configured port
                                            Handles incoming requests
```

### Request Handling Flow

```
┌──────────────────────────────────────────────────────────────────────┐
│                    HTTP Request (e.g., POST /user/register)          │
└──────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
                    ┌─────────────────────────────┐
                    │   HTTP Server (net/http)    │
                    └─────────────────────────────┘
                                  │
                                  ▼
                    ┌─────────────────────────────┐
                    │   Router (ServeMux)         │
                    │   Pattern: "POST /user/..."│
                    └─────────────────────────────┘
                                  │
                                  ▼
                    ┌─────────────────────────────┐
                    │   Route Handler             │
                    │   (routes/user_rotues.go)   │
                    └─────────────────────────────┘
                                  │
                                  ▼
                    ┌─────────────────────────────┐
                    │   HTTP Handler              │
                    │   (handlers/user.go)        │
                    │   CreateUserHandler()       │
                    └─────────────────────────────┘
                                  │
                    ┌─────────────┴──────────────┐
                    │                             │
                    ▼                             ▼
        ┌───────────────────────┐    ┌───────────────────────┐
        │  DTO (dtos/)          │    │  Utils (utils/)       │
        │  Decode request body  │    │  Hash password        │
        │  into CreateUserReq   │    │  Response helpers     │
        └───────────────────────┘    └───────────────────────┘
                    │
                    ▼
        ┌───────────────────────────────────┐
        │  Database Queries (store/)        │
        │  h.Queries.CreateUser(ctx, ...)   │
        │  Type-safe sqlc-generated code    │
        └───────────────────────────────────┘
                    │
                    ▼
        ┌───────────────────────────────────┐
        │  PostgreSQL Database              │
        │  INSERT INTO users(...)           │
        └───────────────────────────────────┘
                    │
                    ▼
        ┌───────────────────────────────────┐
        │  Response (utils/)                │
        │  RespondWithSuccess(...)          │
        │  or RespondWithError(...)         │
        └───────────────────────────────────┘
                    │
                    ▼
        ┌───────────────────────────────────┐
        │  JSON Response to Client          │
        │  {"message": "...", "data": "..."} │
        └───────────────────────────────────┘
```

## Component Communication

### How Packages Talk to Each Other

```
main.go
  │
  ├──> serverconfig    (loads .env → Config)
  │
  ├──> dbconfig        (connects to PostgreSQL → *sql.DB)
  │
  ├──> store           (creates query interface → *Queries)
  │      │
  │      └──> Uses *sql.DB to execute queries
  │
  ├──> handlers        (creates Handler with DB + Queries)
  │      │
  │      ├──> Uses dtos      (for request/response structures)
  │      ├──> Uses utils     (for password, JWT, responses)
  │      └──> Uses store     (for database operations)
  │
  └──> routes          (registers handlers with ServeMux)
         │
         └──> Uses handlers  (to register HTTP endpoints)
```

### Dependency Injection Pattern

The `Handler` struct demonstrates dependency injection:

```go
// In handlers/core_handler.go
type Handler struct {
    DB      *sql.DB         // Injected database connection
    Queries store.Queries   // Injected query interface
}

// Created in main.go
handler := handlers.NewHandlers(db, queries)
```

**Benefits**:
- Testability: Can mock DB and Queries for testing
- Centralization: All dependencies in one place
- Flexibility: Easy to swap implementations

### Closure Pattern for Handlers

All handlers return `http.HandlerFunc` to create closures:

```go
func (h *Handler) CreateUserHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Access h.DB and h.Queries here
    }
}
```

**Why?**
- The returned function has access to `h` (closure)
- Each handler can access shared dependencies
- Follows `http.HandlerFunc` signature

## Current Features

### Endpoints

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/health` | `HealthHandler` | Health check - returns server status |
| GET | `/test` | `TestHandler` | Test endpoint - verifies routing works |
| POST | `/user/register` | `CreateUserHandler` | Create new user with hashed password |

### Implemented Functionality

- ✅ **User Registration**: Creates users with bcrypt-hashed passwords
- ✅ **Database Integration**: PostgreSQL with sqlc-generated queries
- ✅ **Error Handling**: Standardized JSON error responses
- ✅ **Success Responses**: Standardized JSON success responses
- ✅ **Configuration Management**: Environment-based config with `.env` file
- ✅ **Password Security**: Bcrypt hashing with default cost factor
- ✅ **Password Comparison**: Function ready for login implementation
- ✅ **JWT Generation**: Token generation ready (not yet used)
- ✅ **Type-Safe Queries**: sqlc-generated database queries
- ✅ **Request DTOs**: Structured request validation

### Available Database Queries

Generated by sqlc from `internal/migrations/queries.sql`:

**User Queries**:
- `CreateUser(ctx, params)` - Insert new user
- `GetUser(ctx, id)` - Get user by ID
- `ListUsers(ctx)` - Get all users

**Blog Queries** (prepared, not yet exposed via API):
- `CreateBlog(ctx, params)` - Insert new blog post
- `ListBlogs(ctx)` - Get all blog posts

## Technologies & Libraries

- **Language**: Go 1.21+
- **Database**: PostgreSQL
- **SQL Generator**: [sqlc](https://sqlc.dev/) - Type-safe SQL code generation
- **Password Hashing**: [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - Industry standard password hashing
- **JWT**: [jwt-go](https://github.com/dgrijalva/jwt-go) - JSON Web Tokens
- **Environment Config**: [godotenv](https://github.com/joho/godotenv) - .env file loading
- **HTTP**: Standard library (`net/http`) - No external web framework
- **Database Driver**: [pq](https://github.com/lib/pq) - Pure Go PostgreSQL driver

## Key Learning Concepts

### 1. HTTP Router (ServeMux)
Go's built-in router that maps URL patterns to handler functions. Supports method-specific routing (`POST /path`).

### 2. Handler Pattern
Functions that return `http.HandlerFunc`, creating closures over dependencies like database connections.

### 3. Context Usage
Request context (`r.Context()`) is passed to database operations for timeout handling and cancellation.

### 4. JSON Encoding/Decoding
- `json.NewEncoder(w).Encode()` for responses
- `json.NewDecoder(r.Body).Decode()` for requests

### 5. Dependency Injection
Handler struct holds all dependencies (DB, Queries), injected at creation time.

### 6. Password Security
Never store plain text passwords. Always hash with bcrypt before database storage.

### 7. Error Handling
Proper HTTP status codes and standardized error responses for consistent API behavior.

### 8. sqlc Benefits
- Compile-time SQL validation
- Type-safe database operations
- No ORM overhead or reflection
- Full SQL control

### 9. DTOs (Data Transfer Objects)
Separate API contracts from database models for security, validation, and flexibility.

### 10. Package Organization
Clean separation of concerns: handlers, routes, store, utils, dtos, config.

## Folder Documentation

Each folder has its own detailed README explaining:
- Purpose and responsibilities
- How files connect and communicate
- Usage examples and patterns
- Key learning points

**Start exploring**:
1. [`serverconfig/README.md`](serverconfig/README.md) - Configuration management
2. [`dbconfig/README.md`](dbconfig/README.md) - Database connection
3. [`internal/store/README.md`](internal/store/README.md) - sqlc and queries
4. [`internal/handlers/README.md`](internal/handlers/README.md) - HTTP handlers
5. [`internal/routes/README.md`](internal/routes/README.md) - Route registration
6. [`internal/utils/README.md`](internal/utils/README.md) - Utility functions
7. [`internal/dtos/README.md`](internal/dtos/README.md) - Request/response structures
8. [`models/README.md`](models/README.md) - Domain models

## Next Steps

See [`TODO.md`](TODO.md) for a comprehensive list of potential improvements and features to learn next.

## Learning Resources

This project demonstrates common patterns in Go web development:
- RESTful API design
- Database integration without ORMs
- Secure password handling
- Clean code architecture
- Go standard library HTTP server
- Type-safe SQL with sqlc
