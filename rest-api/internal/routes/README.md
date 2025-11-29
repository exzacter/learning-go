# Routes

This package registers all HTTP routes with the router.

## Files

- `setup_routes.go` - Main route setup that calls all sub-route files
- `health_routes.go` - Health check route registration
- `test_routes.go` - Test route registration
- `user_rotues.go` - User-related route registration

## How Routes Work

Routes connect URL paths to handler functions using `http.ServeMux`.

### Setup Routes (`setup_routes.go`)

Central function that calls all individual route setup functions:

```go
func SetupRoutes(mux *http.ServeMux, handler *handlers.Handler) {
    SetupHealthRoute(mux, handler)
    SetupTestRoute(mux, handler)
    SetupUserRoute(mux, handler)
}
```

**Called from** `main.go`:
```go
mux := http.NewServeMux()
routes.SetupRoutes(mux, handler)
```

**Pattern**: Each feature gets its own route setup function for organization.

---

### Health Routes (`health_routes.go`)

```go
func SetupHealthRoute(mux *http.ServeMux, handler *handlers.Handler) {
    mux.HandleFunc("/health", handler.HealthHandler())
}
```

**Registers**: `GET /health` → `handler.HealthHandler()`

---

### Test Routes (`test_routes.go`)

```go
func SetupTestRoute(mux *http.ServeMux, handler *handlers.Handler) {
    mux.HandleFunc("/test", handler.TestHandler())
}
```

**Registers**: `GET /test` → `handler.TestHandler()`

---

### User Routes (`user_rotues.go`)

```go
func SetupUserRoute(mux *http.ServeMux, handler *handlers.Handler) {
    mux.HandleFunc("POST /user/register", handler.CreateUserHandler())
}
```

**Registers**: `POST /user/register` → `handler.CreateUserHandler()`

**Method Restriction**: The `POST` prefix ensures only POST requests match this route.

## Router Flow

```
Client Request
  └─> HTTP Server receives request
        └─> ServeMux (mux) matches URL pattern
              └─> Calls registered handler function
                    └─> Handler processes request
                          └─> Response sent to client
```

### Example: User Registration Request

```
POST /user/register
  └─> mux.HandleFunc("POST /user/register", ...)  // Route matches
        └─> handler.CreateUserHandler()             // Handler called
              └─> [Handler logic executes]
                    └─> JSON response sent
```

## ServeMux Pattern Matching

Go's `http.ServeMux` supports:

1. **Exact paths**: `/health` matches only `/health`
2. **Method prefixes**: `POST /user/register` only matches POST requests
3. **Subtree paths**: `/api/` matches `/api/*` (not used yet in this project)

**Pattern Priority**: More specific patterns take precedence.

## How Routes Connect

### Called by main.go
```go
mux := http.NewServeMux()                        // Create router
routes.SetupRoutes(mux, handler)                 // Register all routes
server := &http.Server{Addr: ":8080", Handler: mux}  // Attach to server
```

### Calls Handlers
```go
mux.HandleFunc("/path", handler.SomeHandler())
```

The handler function is invoked when a matching request arrives.

## Key Learning Points

1. **ServeMux**: Go's built-in HTTP request multiplexer (router)
2. **HandleFunc**: Registers a handler function for a pattern
3. **Method Matching**: `POST /path` restricts to HTTP method
4. **Organization**: Each feature has its own route file
5. **Centralization**: `SetupRoutes()` is the single entry point
6. **Dependency Passing**: Handler is passed down to all route functions
7. **Pattern Matching**: More specific patterns override general ones

## Adding New Routes

To add a new endpoint:

1. **Create handler** in `internal/handlers/`
2. **Create route file** (e.g., `internal/routes/example_routes.go`):
   ```go
   func SetupExampleRoute(mux *http.ServeMux, handler *handlers.Handler) {
       mux.HandleFunc("GET /example", handler.ExampleHandler())
   }
   ```
3. **Register in setup_routes.go**:
   ```go
   func SetupRoutes(mux *http.ServeMux, handler *handlers.Handler) {
       // ... existing routes
       SetupExampleRoute(mux, handler)
   }
   ```
