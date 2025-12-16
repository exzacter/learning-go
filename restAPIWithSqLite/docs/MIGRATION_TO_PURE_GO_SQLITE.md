# Migration Guide: CGO SQLite to Pure Go SQLite

This guide covers all changes needed to migrate from `mattn/go-sqlite3` (CGO) to `modernc.org/sqlite` (pure Go).

## Why Migrate?

| CGO Version | Pure Go Version |
|-------------|-----------------|
| Requires GCC on Windows | Works everywhere |
| Complex Docker builds | Simple Docker builds |
| Cross-compilation is hard | Cross-compilation is trivial |
| Faster performance | ~2x slower (still fast enough) |
| C code security risks | Pure Go memory safety |

---

## Files That Need Changes

```
restAPIWithSqLite/
├── go.mod                    ← Add new dependency
├── cmd/
│   ├── api/
│   │   └── main.go           ← Change import + driver name
│   └── migrate/
│       └── main.go           ← Change import + driver name + migrate driver
```

---

## Step 1: Update go.mod

Run these commands in your project root:

```bash
# Remove the CGO driver (optional, keeps go.mod clean)
go get github.com/mattn/go-sqlite3@none

# Add the pure Go driver
go get modernc.org/sqlite

# The golang-migrate sqlite3 driver also uses CGO, so we need a workaround
# (see Step 3 for details)
go mod tidy
```

Your `go.mod` should change from:
```diff
module learning/go/restAPIWithSqLite

go 1.21

require (
-   github.com/mattn/go-sqlite3 v1.14.22
+   modernc.org/sqlite v1.29.1
    github.com/gin-gonic/gin v1.9.1
    github.com/golang-migrate/migrate v3.5.4+incompatible
    github.com/joho/godotenv v1.5.1
)
```

---

## Step 2: Update cmd/api/main.go

### Before (CGO)
```go
package main

import (
    "database/sql"
    "learning/go/restAPIWithSqLite/internal/database"
    "learning/go/restAPIWithSqLite/internal/env"
    "log"

    _ "github.com/joho/godotenv/autoload"
    _ "github.com/mattn/go-sqlite3"
)

type application struct {
    port      int
    jwtSecret string
    models    database.Models
}

func main() {
    db, err := sql.Open("sqlite3", "./data.db")
    if err != nil {
        log.Fatal(err)
    }

    defer db.Close()

    models := database.NewModels(db)

    app := &application{
        port:      env.GetEnvInt("PORT", 8080),
        jwtSecret: env.GetEnvString("JWT_SECRET", "some-secret-123"),
        models:    models,
    }

    if err := app.serve(); err != nil {
        log.Fatal(err)
    }
}
```

### After (Pure Go)
```go
package main

import (
    "database/sql"
    "learning/go/restAPIWithSqLite/internal/database"
    "learning/go/restAPIWithSqLite/internal/env"
    "log"

    _ "github.com/joho/godotenv/autoload"
    _ "modernc.org/sqlite"  // ← CHANGED: Pure Go SQLite
)

type application struct {
    port      int
    jwtSecret string
    models    database.Models
}

func main() {
    db, err := sql.Open("sqlite", "./data.db")  // ← CHANGED: "sqlite" not "sqlite3"
    if err != nil {
        log.Fatal(err)
    }

    defer db.Close()

    // Optional: Enable WAL mode for better performance
    db.Exec("PRAGMA journal_mode=WAL")
    db.Exec("PRAGMA busy_timeout=5000")

    models := database.NewModels(db)

    app := &application{
        port:      env.GetEnvInt("PORT", 8080),
        jwtSecret: env.GetEnvString("JWT_SECRET", "some-secret-123"),
        models:    models,
    }

    if err := app.serve(); err != nil {
        log.Fatal(err)
    }
}
```

### Diff Summary
```diff
- _ "github.com/mattn/go-sqlite3"
+ _ "modernc.org/sqlite"

- db, err := sql.Open("sqlite3", "./data.db")
+ db, err := sql.Open("sqlite", "./data.db")
```

---

## Step 3: Update cmd/migrate/main.go

The `golang-migrate` library's SQLite3 driver also uses CGO. You have two options:

### Option A: Use a Pure Go Migration Approach (Recommended)

Replace the CGO-dependent migrate driver with a simple pure Go solution:

### Before (CGO)
```go
package main

import (
    "database/sql"
    "log"
    "os"

    "github.com/golang-migrate/migrate"
    "github.com/golang-migrate/migrate/database/sqlite3"
    "github.com/golang-migrate/migrate/source/file"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Please provide a migration direction: 'up' or 'down'")
    }

    direction := os.Args[1]

    db, err := sql.Open("sqlite3", "./data.db")
    if err != nil {
        log.Fatal(err)
    }

    defer db.Close()

    instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
    if err != nil {
        log.Fatal(err)
    }

    fSrc, err := (&file.File{}).Open("./cmd/migrate/migrations")
    if err != nil {
        log.Fatal(err)
    }

    m, err := migrate.NewWithInstance("file", fSrc, "sqlite3", instance)
    if err != nil {
        log.Fatal(err)
    }

    switch direction {
    case "up":
        if err := m.Up(); err != nil && err != migrate.ErrNoChange {
            log.Fatal(err)
        }
    case "down":
        if err := m.Down(); err != nil && err != migrate.ErrNoChange {
            log.Fatal(err)
        }
    default:
        log.Fatal("Invalid direction. Use 'up' or 'down'")
    }
}
```

### After (Pure Go - Simple Custom Migrator)
```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "sort"
    "strings"

    _ "modernc.org/sqlite"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Please provide a migration direction: 'up' or 'down'")
    }

    direction := os.Args[1]

    db, err := sql.Open("sqlite", "./data.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create migrations tracking table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version TEXT PRIMARY KEY,
            applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        log.Fatal(err)
    }

    migrationsDir := "./cmd/migrate/migrations"

    switch direction {
    case "up":
        if err := migrateUp(db, migrationsDir); err != nil {
            log.Fatal(err)
        }
        fmt.Println("Migrations applied successfully")
    case "down":
        if err := migrateDown(db, migrationsDir); err != nil {
            log.Fatal(err)
        }
        fmt.Println("Migrations rolled back successfully")
    default:
        log.Fatal("Invalid direction. Use 'up' or 'down'")
    }
}

func migrateUp(db *sql.DB, dir string) error {
    files, err := filepath.Glob(filepath.Join(dir, "*.up.sql"))
    if err != nil {
        return err
    }
    sort.Strings(files)

    for _, file := range files {
        version := extractVersion(file)

        // Check if already applied
        var count int
        db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", version).Scan(&count)
        if count > 0 {
            continue
        }

        // Read and execute migration
        content, err := os.ReadFile(file)
        if err != nil {
            return fmt.Errorf("failed to read %s: %w", file, err)
        }

        _, err = db.Exec(string(content))
        if err != nil {
            return fmt.Errorf("failed to execute %s: %w", file, err)
        }

        // Record migration
        _, err = db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version)
        if err != nil {
            return fmt.Errorf("failed to record migration %s: %w", version, err)
        }

        fmt.Printf("Applied: %s\n", filepath.Base(file))
    }

    return nil
}

func migrateDown(db *sql.DB, dir string) error {
    // Get the latest applied migration
    var version string
    err := db.QueryRow("SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1").Scan(&version)
    if err != nil {
        if err == sql.ErrNoRows {
            fmt.Println("No migrations to roll back")
            return nil
        }
        return err
    }

    // Find corresponding down file
    files, err := filepath.Glob(filepath.Join(dir, version+"*.down.sql"))
    if err != nil || len(files) == 0 {
        return fmt.Errorf("down migration not found for version %s", version)
    }

    // Read and execute migration
    content, err := os.ReadFile(files[0])
    if err != nil {
        return fmt.Errorf("failed to read %s: %w", files[0], err)
    }

    _, err = db.Exec(string(content))
    if err != nil {
        return fmt.Errorf("failed to execute %s: %w", files[0], err)
    }

    // Remove migration record
    _, err = db.Exec("DELETE FROM schema_migrations WHERE version = ?", version)
    if err != nil {
        return fmt.Errorf("failed to remove migration record: %w", err)
    }

    fmt.Printf("Rolled back: %s\n", filepath.Base(files[0]))

    return nil
}

func extractVersion(filename string) string {
    base := filepath.Base(filename)
    // Extract "000001" from "000001_create_users_table.up.sql"
    parts := strings.Split(base, "_")
    if len(parts) > 0 {
        return parts[0]
    }
    return base
}
```

---

### Option B: Use goose (Pure Go Migration Tool)

Alternatively, switch to [goose](https://github.com/pressly/goose) which has better pure Go support:

```bash
go get github.com/pressly/goose/v3
```

```go
package main

import (
    "database/sql"
    "log"
    "os"

    "github.com/pressly/goose/v3"
    _ "modernc.org/sqlite"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Please provide a migration direction: 'up' or 'down'")
    }

    direction := os.Args[1]

    db, err := sql.Open("sqlite", "./data.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    migrationsDir := "./cmd/migrate/migrations"

    // Set the dialect
    if err := goose.SetDialect("sqlite3"); err != nil {
        log.Fatal(err)
    }

    switch direction {
    case "up":
        if err := goose.Up(db, migrationsDir); err != nil {
            log.Fatal(err)
        }
    case "down":
        if err := goose.Down(db, migrationsDir); err != nil {
            log.Fatal(err)
        }
    default:
        log.Fatal("Invalid direction. Use 'up' or 'down'")
    }
}
```

**Note:** With goose, you may need to rename your migration files to match goose's format:
```
000001_create_users_table.up.sql    → 00001_create_users_table.sql
000001_create_users_table.down.sql  → (embedded in same file with -- +goose Down)
```

---

## Step 4: Update Docker (if applicable)

### Before (CGO - Complex)
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app

# CGO requires these
RUN apk add --no-cache gcc musl-dev
ENV CGO_ENABLED=1

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags='-s -w -extldflags "-static"' -o /api ./cmd/api

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /api /api
CMD ["/api"]
```

### After (Pure Go - Simple)
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app

# No GCC needed!
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags='-s -w' -o /api ./cmd/api

# Can even use scratch now!
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /api /api
CMD ["/api"]
```

---

## Summary of All Changes

| File | Change |
|------|--------|
| `go.mod` | Replace `github.com/mattn/go-sqlite3` with `modernc.org/sqlite` |
| `cmd/api/main.go` | Change import and `sql.Open("sqlite3"` → `sql.Open("sqlite"` |
| `cmd/migrate/main.go` | Replace with pure Go migrator or switch to goose |
| `Dockerfile` | Remove `gcc musl-dev`, set `CGO_ENABLED=0` |

---

## Quick Reference

| Aspect | Before (CGO) | After (Pure Go) |
|--------|--------------|-----------------|
| Import | `_ "github.com/mattn/go-sqlite3"` | `_ "modernc.org/sqlite"` |
| Driver | `"sqlite3"` | `"sqlite"` |
| Build | `CGO_ENABLED=1` + GCC | `CGO_ENABLED=0` |
| Docker base | `alpine` with `gcc musl-dev` | `scratch` or `alpine` |

---

## Testing After Migration

```bash
# Run migrations
go run ./cmd/migrate up

# Start API
go run ./cmd/api

# Test endpoints
curl http://localhost:8080/api/v1/events
```

---

## Rollback Plan

If you need to switch back to CGO:

1. Revert `go.mod` changes
2. Change `"sqlite"` back to `"sqlite3"`
3. Change import back to `_ "github.com/mattn/go-sqlite3"`
4. Restore original `cmd/migrate/main.go`
5. Run `go mod tidy`

Your database file (`data.db`) is compatible with both drivers - no data migration needed.
