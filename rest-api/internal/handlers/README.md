# HTTP Handlers

This package contains all HTTP request handlers for the API.

## Files

- `core_handler.go` - Core Handler struct with dependencies
- `health.go` - Health check endpoint handler
- `test.go` - Test endpoint handler
- `user.go` - User-related endpoint handlers

## Core Handler Structure

### Handler Struct (`core_handler.go`)

The `Handler` struct is the foundation for all HTTP handlers:

```go
type Handler struct {
    DB      *sql.DB         // Database connection pool
    Queries store.Queries   // sqlc generated queries
}
```

**Why this pattern?**
- **Dependency Injection**: All handlers have access to DB and queries
- **Testability**: Easy to mock DB and queries for testing
- **Encapsulation**: All dependencies in one place

### Initialization

Created in `main.go` and passed to routes:
```go
handler := handlers.NewHandlers(db, queries)
routes.SetupRoutes(mux, handler)
```

## Handler Pattern Explained

All handlers follow this pattern:

```go
func (h *Handler) HandlerName() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Handler logic here
    }
}
```

**Why return `http.HandlerFunc`?**
- Creates a closure over the `Handler` struct
- Each handler gets access to `h.DB` and `h.Queries`
- Follows Go's HTTP handler interface

## Individual Handlers

### Health Handler (`health.go`)

**Endpoint**: `GET /health`

**Purpose**: Simple health check to verify server is running

**Flow**:
1. Sets `Content-Type: application/json` header
2. Creates response map with message
3. Encodes to JSON and sends to client

**Response**:
```json
{
    "message": "Server is Okay"
}
```

**No database access** - just confirms the HTTP server is responding.

---

### Test Handler (`test.go`)

**Endpoint**: `GET /test`

**Purpose**: Verify routing and handler system is working

**Flow**:
1. Sets `Content-Type: application/json` header
2. Creates response map with message
3. Encodes to JSON and sends to client

**Response**:
```json
{
    "message": "Test curl has worked, handler is working for that function"
}
```

**No database access** - tests the request/response pipeline.

---

### Create User Handler (`user.go`)

**Endpoint**: `POST /user/register`

**Purpose**: Register a new user with hashed password

**Flow**:
1. **Get Context**: `ctx := r.Context()` for database operations
2. **Decode Request**: Parse JSON body into `CreateUserRequest` DTO
3. **Hash Password**: Use `utils.HashPassword()` to hash plain text password
4. **Insert to DB**: Call `h.Queries.CreateUser()` with hashed password
5. **Send Response**: Return success or error using utils functions

**Request Body**:
```json
{
    "username": "testuser",
    "email": "test@example.com",
    "password": "securepassword"
}
```

**Success Response** (201 Created):
```json
{
    "message": "user created",
    "data": "testuser"
}
```

**Error Responses**:
- **400 Bad Request**: Invalid JSON payload
- **500 Internal Server Error**: Password hashing failed or DB error

### Detailed User Handler Flow

```
POST /user/register
  └─> CreateUserHandler()
        ├─> Get request context
        ├─> Decode JSON into CreateUserRequest DTO
        │     ├─ If error: RespondWithError(400, "Invalid request payload")
        │     └─ Return
        ├─> Hash password with bcrypt
        │     ├─ If error: RespondWithError(500, "error while hashing password")
        │     └─ Return
        ├─> Call h.Queries.CreateUser() with:
        │     ├─ Username from request
        │     ├─ Email from request
        │     └─ Hashed password
        ├─> If DB error: RespondWithError(500, "error creating user")
        └─> Success: RespondWithSuccess(201, "user created", username)
```

## How Handlers Connect to Other Packages

### To Routes (`internal/routes/`)
Routes register handlers with the router:
```go
// In routes/user_routes.go
mux.HandleFunc("POST /user/register", handler.CreateUserHandler())
```

### To Store (`internal/store/`)
Handlers call database queries:
```go
result, err := h.Queries.CreateUser(ctx, store.CreateUserParams{...})
```

### To Utils (`internal/utils/`)
Handlers use utility functions:
```go
hashedPassword, err := utils.HashPassword(req.Password)
utils.RespondWithError(w, http.StatusBadRequest, "message")
utils.RespondWithSuccess(w, http.StatusCreated, "message", data)
```

### To DTOs (`internal/dtos/`)
Handlers decode requests into DTOs:
```go
var req dtos.CreateUserRequest
json.NewDecoder(r.Body).Decode(&req)
```

## Key Learning Points

1. **Method Receivers**: `(h *Handler)` gives all methods access to dependencies
2. **Closures**: Returning `http.HandlerFunc` creates closure over `h`
3. **Context**: Always pass `r.Context()` to database operations
4. **Error Handling**: Return early on errors, use appropriate status codes
5. **JSON Encoding**: `json.NewEncoder(w).Encode()` for responses
6. **JSON Decoding**: `json.NewDecoder(r.Body).Decode()` for requests
7. **Password Security**: Never store plain text passwords; always hash
8. **Separation**: Handlers orchestrate, they don't implement business logic
