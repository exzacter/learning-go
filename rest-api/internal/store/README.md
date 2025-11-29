# Database Store (sqlc Generated)

This package contains **sqlc-generated** type-safe database queries and models.

## Files

- `db.go` - Query interface and database abstraction (generated)
- `models.go` - Database models matching schema (generated)
- `queries.sql.go` - Type-safe query functions (generated)

## What is sqlc?

**sqlc** generates type-safe Go code from SQL queries:
- Write SQL in `internal/migrations/queries.sql`
- Define schema in `internal/migrations/schema.sql`
- Run `sqlc generate` to create Go code
- Get compile-time type safety with zero runtime overhead

## Generated Models

### User Model
```go
type User struct {
    ID       int32        `json:"id"`
    Username string       `json:"username"`
    Email    string       `json:"email"`
    Password string       `json:"password"`
    Created  sql.NullTime `json:"created"`
    Updated  sql.NullTime `json:"updated"`
}
```

### Blog Model
```go
type Blog struct {
    ID      int32        `json:"id"`
    Title   string       `json:"title"`
    Content string       `json:"content"`
    UserID  int32        `json:"user_id"`
    Created sql.NullTime `json:"created"`
    Updated sql.NullTime `json:"updated"`
}
```

## Available Queries

All queries are methods on the `Queries` struct:

### User Queries

**CreateUser** - Insert new user
```go
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error)
```
- **Input**: `CreateUserParams` (username, email, password, created, updated)
- **Output**: `CreateUserRow` (id, username, email, created, updated)
- **Returns**: Newly created user data (excludes password)

**GetUser** - Get user by ID
```go
func (q *Queries) GetUser(ctx context.Context, id int32) (GetUserRow, error)
```
- **Input**: User ID
- **Output**: `GetUserRow` (id, username, email, created, updated)

**ListUsers** - Get all users
```go
func (q *Queries) ListUsers(ctx context.Context) ([]ListUsersRow, error)
```
- **Output**: Slice of `ListUsersRow` (all users, ordered by ID)

### Blog Queries

**CreateBlog** - Insert new blog post
```go
func (q *Queries) CreateBlog(ctx context.Context, arg CreateBlogParams) (CreateBlogRow, error)
```

**ListBlogs** - Get all blog posts
```go
func (q *Queries) ListBlogs(ctx context.Context) ([]Blog, error)
```

## How It Connects

### Initialization (main.go)
```go
queries := store.New(db)  // Create Queries instance
handler := handlers.NewHandlers(db, queries)  // Pass to handlers
```

### Usage in Handlers (user.go)
```go
// Inside CreateUserHandler
result, err := h.Queries.CreateUser(ctx, store.CreateUserParams{
    Username: req.Username,
    Email:    req.Email,
    Password: hashedPassword,
})
```

## Query Interface (`Queries` struct)

The `Queries` struct is the central interface:
```go
type Queries struct {
    db             DBTX           // Database connection
    tx             *sql.Tx        // Optional transaction
    createUserStmt *sql.Stmt      // Prepared statement cache
    getUserStmt    *sql.Stmt
    // ... more prepared statements
}
```

### Key Methods

- `New(db DBTX)` - Create new Queries instance
- `WithTx(tx *sql.Tx)` - Create Queries scoped to transaction
- `Close()` - Close all prepared statements

## DBTX Interface

Generic database interface that works with both `*sql.DB` and `*sql.Tx`:
```go
type DBTX interface {
    ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
    PrepareContext(context.Context, string) (*sql.Stmt, error)
    QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
    QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}
```

This allows queries to work in both regular and transaction contexts.

## Regenerating Code

When you modify SQL queries or schema:

1. Edit `internal/migrations/queries.sql` or `schema.sql`
2. Run: `sqlc generate`
3. All code in this package is regenerated

**DO NOT manually edit generated files!**

## Benefits of sqlc

1. **Type Safety**: Compile-time validation of SQL queries
2. **No Reflection**: Zero runtime overhead
3. **No ORM Magic**: Full control over SQL
4. **IDE Support**: Auto-completion for query parameters
5. **Catch Errors Early**: SQL errors found at compile time, not runtime
