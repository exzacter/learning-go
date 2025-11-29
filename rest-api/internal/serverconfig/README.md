# Server Configuration

This package handles loading and managing server configuration from environment variables.

## Files

- `config.go` - Configuration loading and management

## How It Works

### Configuration Structure

The `Config` struct holds all server configuration:

```go
type Config struct {
    ServerPort  string  // Port the server listens on
    DatabaseURL string  // PostgreSQL connection string
    Environment string  // Environment (development, production, etc.)
    LogLevel    string  // Logging level (info, debug, error)
}
```

### Loading Configuration

`LoadConfig()` is called at application startup in `main.go`:

1. Loads `.env` file using `godotenv.Load()`
2. Reads environment variables using `getEnv()` helper
3. Returns populated `Config` struct with defaults if env vars are missing

**Default Values:**
- `SERVER_PORT`: `8080`
- `DATABASE_URL`: `postgres`
- `ENVIRONMENT`: `development`
- `LOG_LEVEL`: `info`

### Usage in main.go

```go
config, err := serverconfig.LoadConfig()
if err != nil {
    log.Fatalf("Failed to load config %v", err)
}
```

The config is then used throughout the application:
- `config.DatabaseURL` → Passed to `dbconfig.ConnectDB()`
- `config.ServerPort` → Used to set server address

## Environment Variables

Create a `.env` file in the `rest-api/` directory:

```env
SERVER_PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
ENVIRONMENT=development
LOG_LEVEL=info
```

## Key Learning Points

1. **Separation of Concerns**: Configuration is isolated from business logic
2. **Defaults**: Always provide sensible defaults for missing environment variables
3. **godotenv**: Simplifies loading `.env` files in development
4. **Error Handling**: Returns errors if `.env` file fails to load
