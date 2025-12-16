# SQLite Implementation Guide

## How SQLite is Implemented in This API

### 1. Database Driver

The API uses the `github.com/mattn/go-sqlite3` driver, which is a CGO-based SQLite3 driver for Go.

```go
// cmd/api/main.go
import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"  // Blank import registers the driver
)

func main() {
    db, err := sql.Open("sqlite3", "./data.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
}
```

**Key points:**
- `sql.Open("sqlite3", "./data.db")` opens (or creates) the database file
- The `_` blank import registers the SQLite driver with Go's `database/sql` package
- `defer db.Close()` ensures the connection closes when the app exits

---

### 2. Database Connection Flow

```
main.go                    models.go                  users.go/events.go/attendees.go
┌─────────────┐           ┌─────────────┐            ┌─────────────┐
│ sql.Open()  │──► db ───►│ NewModels() │───► db ───►│ UserModel   │
│             │           │             │            │ EventModel  │
│             │           │             │            │ AttendeeModel│
└─────────────┘           └─────────────┘            └─────────────┘
```

The `*sql.DB` connection is passed through the Models struct to each individual model, giving them database access.

---

### 3. Migration System

Migrations are handled by `github.com/golang-migrate/migrate`. SQL files live in `cmd/migrate/migrations/`:

```
migrations/
├── 000001_create_users_table.up.sql      (creates users table)
├── 000001_create_users_table.down.sql    (drops users table)
├── 000002_create_events_table.up.sql
├── 000002_create_events_table.down.sql
├── 000003_create_attendees_table.up.sql
└── 000003_create_attendees_table.down.sql
```

**Run migrations:**
```bash
go run cmd/migrate/main.go up    # Apply all migrations
go run cmd/migrate/main.go down  # Rollback all migrations
```

---

### 4. Model Pattern

Each database entity follows this pattern:

```go
// internal/database/users.go
type UserModel struct {
    DB *sql.DB  // Database connection
}

type User struct {
    Id       int    `json:"id"`
    Email    string `json:"email"`
    Name     string `json:"name"`
    Password string `json:"-"`  // "-" hides from JSON output
}

func (m *UserModel) Insert(user *User) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    query := "INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id"
    return m.DB.QueryRowContext(ctx, query, user.Email, user.Password, user.Name).Scan(&user.Id)
}
```

**Pattern components:**
- **Model struct** - Holds the DB connection
- **Data struct** - Represents a row with JSON tags
- **Methods** - CRUD operations with context timeouts

---

## Adding New Tables/Features

### Step 1: Create Migration Files

Create two new files in `cmd/migrate/migrations/`:

```sql
-- 000004_create_venues_table.up.sql
CREATE TABLE IF NOT EXISTS venues (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    capacity INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 000004_create_venues_table.down.sql
DROP TABLE IF EXISTS venues;
```

### Step 2: Create the Model File

Create `internal/database/venues.go`:

```go
package database

import (
    "context"
    "database/sql"
    "time"
)

type VenueModel struct {
    DB *sql.DB
}

type Venue struct {
    Id        int    `json:"id"`
    Name      string `json:"name"`
    Address   string `json:"address"`
    Capacity  int    `json:"capacity"`
    CreatedAt string `json:"createdAt"`
}

func (m *VenueModel) Insert(venue *Venue) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    query := "INSERT INTO venues (name, address, capacity) VALUES ($1, $2, $3) RETURNING id"
    return m.DB.QueryRowContext(ctx, query, venue.Name, venue.Address, venue.Capacity).Scan(&venue.Id)
}

func (m *VenueModel) Get(id int) (*Venue, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    query := "SELECT id, name, address, capacity, created_at FROM venues WHERE id = $1"
    var venue Venue
    err := m.DB.QueryRowContext(ctx, query, id).Scan(
        &venue.Id, &venue.Name, &venue.Address, &venue.Capacity, &venue.CreatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return &venue, err
}

func (m *VenueModel) GetAll() ([]*Venue, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    query := "SELECT id, name, address, capacity, created_at FROM venues"
    rows, err := m.DB.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var venues []*Venue
    for rows.Next() {
        var venue Venue
        err := rows.Scan(&venue.Id, &venue.Name, &venue.Address, &venue.Capacity, &venue.CreatedAt)
        if err != nil {
            return nil, err
        }
        venues = append(venues, &venue)
    }
    return venues, nil
}
```

### Step 3: Register in Models

Update `internal/database/models.go`:

```go
type Models struct {
    Users     UserModel
    Events    EventModel
    Attendees AttendeeModel
    Venues    VenueModel  // Add this
}

func NewModels(db *sql.DB) Models {
    return Models{
        Users:     UserModel{DB: db},
        Events:    EventModel{DB: db},
        Attendees: AttendeeModel{DB: db},
        Venues:    VenueModel{DB: db},  // Add this
    }
}
```

### Step 4: Create HTTP Handlers

Create `cmd/api/venues.go`:

```go
package main

import (
    "net/http"
    "strconv"
    "learning/go/restAPIWithSqLite/internal/database"
    "github.com/gin-gonic/gin"
)

func (app *application) createVenue(c *gin.Context) {
    var venue database.Venue
    if err := c.ShouldBindJSON(&venue); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := app.models.Venues.Insert(&venue); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, venue)
}

func (app *application) getAllVenues(c *gin.Context) {
    venues, err := app.models.Venues.GetAll()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, venues)
}
```

### Step 5: Add Routes

Update `cmd/api/routes.go`:

```go
func (app *application) routes() http.Handler {
    g := gin.Default()
    v1 := g.Group("/api/v1")
    {
        // ... existing routes ...

        // Venue routes
        v1.POST("/venues", app.createVenue)
        v1.GET("/venues", app.getAllVenues)
        v1.GET("/venues/:id", app.getVenue)
    }
    return g
}
```

### Step 6: Run Migration

```bash
go run cmd/migrate/main.go up
```

---

## Importing Large Data from External APIs

When importing thousands of attendees from an external API, you need to handle:
- HTTP requests to fetch data
- Batch inserts for performance
- Transaction safety
- Error handling

### Example: Import Attendees from External Event API

#### Step 1: Create Import Model Method

Add to `internal/database/attendees.go`:

```go
// BatchInsert inserts multiple attendees in a single transaction
func (m *AttendeeModel) BatchInsert(attendees []*Attendee) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Longer timeout for batch
    defer cancel()

    // Start transaction
    tx, err := m.DB.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback() // Rollback if we don't commit

    // Prepare statement for reuse
    stmt, err := tx.PrepareContext(ctx, "INSERT INTO attendees (event_id, user_id) VALUES ($1, $2)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    // Insert each attendee
    for _, attendee := range attendees {
        _, err := stmt.ExecContext(ctx, attendee.EventId, attendee.UserId)
        if err != nil {
            return err // Transaction will rollback
        }
    }

    // Commit transaction
    return tx.Commit()
}

// BatchInsertUsers inserts multiple users and returns their IDs
func (m *UserModel) BatchInsert(users []*User) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    tx, err := m.DB.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    stmt, err := tx.PrepareContext(ctx,
        "INSERT INTO users (email, name, password) VALUES ($1, $2, $3) RETURNING id")
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, user := range users {
        err := stmt.QueryRowContext(ctx, user.Email, user.Name, user.Password).Scan(&user.Id)
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}
```

#### Step 2: Create Import Service

Create `internal/importer/importer.go`:

```go
package importer

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "learning/go/restAPIWithSqLite/internal/database"
)

type Importer struct {
    Models     database.Models
    HTTPClient *http.Client
}

// ExternalAttendee represents the external API's attendee format
type ExternalAttendee struct {
    Email string `json:"email"`
    Name  string `json:"name"`
}

// ExternalEventResponse represents the external API response
type ExternalEventResponse struct {
    EventName  string             `json:"eventName"`
    Attendees  []ExternalAttendee `json:"attendees"`
    TotalCount int                `json:"totalCount"`
    Page       int                `json:"page"`
    TotalPages int                `json:"totalPages"`
}

func NewImporter(models database.Models) *Importer {
    return &Importer{
        Models: models,
        HTTPClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

// ImportEventAttendees fetches attendees from external API and imports them
func (imp *Importer) ImportEventAttendees(apiURL string, eventId int) (*ImportResult, error) {
    result := &ImportResult{}
    page := 1
    batchSize := 500 // Process in batches of 500

    for {
        // Fetch page from external API
        url := fmt.Sprintf("%s?page=%d&limit=%d", apiURL, page, batchSize)
        resp, err := imp.HTTPClient.Get(url)
        if err != nil {
            return result, fmt.Errorf("failed to fetch page %d: %w", page, err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            return result, fmt.Errorf("API returned status %d", resp.StatusCode)
        }

        var apiResponse ExternalEventResponse
        if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
            return result, fmt.Errorf("failed to decode response: %w", err)
        }

        // Convert external attendees to our user format
        users := make([]*database.User, len(apiResponse.Attendees))
        for i, ext := range apiResponse.Attendees {
            users[i] = &database.User{
                Email:    ext.Email,
                Name:     ext.Name,
                Password: "imported-user", // Set a default or generate random
            }
        }

        // Batch insert users
        if err := imp.Models.Users.BatchInsert(users); err != nil {
            return result, fmt.Errorf("failed to insert users batch: %w", err)
        }

        // Create attendee records linking users to event
        attendees := make([]*database.Attendee, len(users))
        for i, user := range users {
            attendees[i] = &database.Attendee{
                EventId: eventId,
                UserId:  user.Id, // ID was set by BatchInsert
            }
        }

        // Batch insert attendees
        if err := imp.Models.Attendees.BatchInsert(attendees); err != nil {
            return result, fmt.Errorf("failed to insert attendees batch: %w", err)
        }

        result.UsersImported += len(users)
        result.AttendeesImported += len(attendees)
        result.PagesProcessed++

        // Check if we've processed all pages
        if page >= apiResponse.TotalPages {
            break
        }
        page++
    }

    return result, nil
}

type ImportResult struct {
    UsersImported     int `json:"usersImported"`
    AttendeesImported int `json:"attendeesImported"`
    PagesProcessed    int `json:"pagesProcessed"`
}
```

#### Step 3: Create Import Handler

Add to `cmd/api/events.go`:

```go
import (
    "learning/go/restAPIWithSqLite/internal/importer"
)

type ImportRequest struct {
    APIURL string `json:"apiUrl" binding:"required"`
}

func (app *application) importEventAttendees(c *gin.Context) {
    // Get event ID from URL
    eventId, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
        return
    }

    // Verify event exists
    event, err := app.models.Events.Get(eventId)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if event == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
        return
    }

    // Get API URL from request body
    var req ImportRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Run import
    imp := importer.NewImporter(app.models)
    result, err := imp.ImportEventAttendees(req.APIURL, eventId)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":         err.Error(),
            "partialResult": result,
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Import completed successfully",
        "result":  result,
    })
}
```

#### Step 4: Add Import Route

```go
v1.POST("/events/:id/import-attendees", app.importEventAttendees)
```

---

## Performance Tips for Large Imports

### 1. Use Transactions
Wrap batch inserts in transactions to ensure atomicity and improve speed:

```go
tx, _ := db.BeginTx(ctx, nil)
defer tx.Rollback()
// ... do inserts ...
tx.Commit()
```

### 2. Use Prepared Statements
Reuse prepared statements within a batch:

```go
stmt, _ := tx.PrepareContext(ctx, "INSERT INTO ...")
for _, item := range items {
    stmt.ExecContext(ctx, item.Field1, item.Field2)
}
```

### 3. Disable Synchronous Writes (for bulk imports only)
```go
db.Exec("PRAGMA synchronous = OFF")
db.Exec("PRAGMA journal_mode = MEMORY")
// ... do bulk import ...
db.Exec("PRAGMA synchronous = FULL")
db.Exec("PRAGMA journal_mode = DELETE")
```

### 4. Use Chunked Processing
Process in chunks to avoid memory issues:

```go
const chunkSize = 1000
for i := 0; i < len(items); i += chunkSize {
    end := i + chunkSize
    if end > len(items) {
        end = len(items)
    }
    chunk := items[i:end]
    // Process chunk
}
```

### 5. Background Processing with Goroutines
For very large imports, process in the background:

```go
func (app *application) importEventAttendeesAsync(c *gin.Context) {
    // ... validation ...

    // Start background job
    go func() {
        imp := importer.NewImporter(app.models)
        result, err := imp.ImportEventAttendees(req.APIURL, eventId)
        // Log result or store in a jobs table
    }()

    c.JSON(http.StatusAccepted, gin.H{
        "message": "Import started in background",
    })
}
```

---

## SQLite Limitations to Consider

| Limitation | Impact | Workaround |
|------------|--------|------------|
| Single writer | Concurrent writes queue up | Use WAL mode: `PRAGMA journal_mode=WAL` |
| No native JSON | Can't query JSON fields | Store as TEXT, parse in Go |
| Max 1GB default | Large datasets | Increase with `PRAGMA max_page_count` |
| No `RETURNING` in older versions | Can't get inserted ID | Use `last_insert_rowid()` |

### Enable WAL Mode (Recommended)
Add this after opening the database in `main.go`:

```go
db, err := sql.Open("sqlite3", "./data.db")
if err != nil {
    log.Fatal(err)
}

// Enable Write-Ahead Logging for better concurrent read performance
db.Exec("PRAGMA journal_mode=WAL")
db.Exec("PRAGMA busy_timeout=5000")  // Wait 5s if database is locked
```

---

## Quick Reference

| Task | Command/Code |
|------|--------------|
| Run migrations up | `go run cmd/migrate/main.go up` |
| Run migrations down | `go run cmd/migrate/main.go down` |
| Add new table | Create `.up.sql` and `.down.sql` in migrations folder |
| Add new model | Create file in `internal/database/`, register in `models.go` |
| Batch insert | Use transactions + prepared statements |
| Import from API | Use HTTP client + batch processing + transactions |
