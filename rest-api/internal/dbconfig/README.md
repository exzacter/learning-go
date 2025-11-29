# Database Configuration

This package handles PostgreSQL database connection setup.

## Files

- `dbconfig.go` - Database connection initialization

## How It Works

### Connection Process

`ConnectDB(databaseURL string)` establishes a PostgreSQL connection:

1. **Open Connection**: Uses `sql.Open("postgres", databaseURL)`
   - Driver: `github.com/lib/pq` (PostgreSQL driver)
   - Returns `*sql.DB` connection pool

2. **Verify Connection**: Calls `db.Ping()`
   - Ensures database is reachable
   - Fails fast if connection is invalid

3. **Return Connection**: Returns `*sql.DB` for use throughout the app

### Usage in main.go

```go
db := dbconfig.ConnectDB(config.DatabaseURL)
defer db.Close()  // Close connection when main() exits
```

The database connection is then:
- Passed to `store.New(db)` to create query interface
- Stored in `handlers.Handler` struct for use in HTTP handlers

## Connection Flow

```
main.go
  └─> dbconfig.ConnectDB(url)
        └─> sql.Open("postgres", url)  // Create connection pool
        └─> db.Ping()                   // Verify connection
        └─> Returns *sql.DB
```

## Error Handling

- **Connection Failure**: `log.Fatal()` if `sql.Open()` fails
- **Ping Failure**: `log.Fatalf()` if database is unreachable

Both errors terminate the application since the database is essential.

## Key Learning Points

1. **Connection Pooling**: `*sql.DB` is a connection pool, not a single connection
2. **Defer Close**: Always `defer db.Close()` after opening
3. **Ping for Validation**: `sql.Open()` doesn't validate connection; `Ping()` does
4. **Driver Import**: `_ "github.com/lib/pq"` imports driver for side effects
5. **Fail Fast**: Database failures should stop the application immediately
