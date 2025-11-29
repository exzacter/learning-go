# Data Transfer Objects (DTOs)

This package contains request and response data structures for the API.

## Files

- `request.go` - Request DTOs for incoming API requests

## What are DTOs?

**Data Transfer Objects** are structs that define the shape of data sent to and from the API.

**Purpose**:
- Validate incoming request data
- Separate API contracts from database models
- Control what data is exposed to clients

## Request DTOs

### CreateUserRequest

```go
type CreateUserRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

**Purpose**: Define the expected JSON structure for user registration

**Used in**: `POST /user/register` endpoint

**Example JSON**:
```json
{
    "username": "testuser",
    "email": "test@example.com",
    "password": "securepassword"
}
```

**Struct Tags**: `json:"username"` maps JSON field to Go struct field

**Usage in Handler**:
```go
var req dtos.CreateUserRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    // Handle error
}

// Access fields
req.Username  // "testuser"
req.Email     // "test@example.com"
req.Password  // "securepassword"
```

## Why DTOs vs Database Models?

### Database Model (`store.User`)
```go
type User struct {
    ID       int32        `json:"id"`
    Username string       `json:"username"`
    Email    string       `json:"email"`
    Password string       `json:"password"`  // Hashed
    Created  sql.NullTime `json:"created"`
    Updated  sql.NullTime `json:"updated"`
}
```

### DTO (`CreateUserRequest`)
```go
type CreateUserRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`  // Plain text
}
```

**Key Differences**:
1. **No ID**: Client doesn't provide ID (database generates it)
2. **No Timestamps**: Server sets created/updated times
3. **Plain Password**: DTO has plain text, model has hashed password
4. **Simpler Types**: DTO uses `string`, model uses `sql.NullTime`

**Benefits**:
- **Security**: Client can't manipulate ID or timestamps
- **Validation**: Can add validation tags to DTOs
- **Flexibility**: API contract independent from database schema
- **Clarity**: Clear separation between input and storage

## DTO → Database Flow

```
Client Request (JSON)
  └─> Decode into CreateUserRequest DTO
        └─> Extract fields from DTO
              └─> Hash password
                    └─> Create CreateUserParams for database
                          └─> Insert into database
                                └─> Return created user (without password)
```

### Example in Handler
```go
// 1. Decode request into DTO
var req dtos.CreateUserRequest
json.NewDecoder(r.Body).Decode(&req)

// 2. Process DTO data
hashedPassword, _ := utils.HashPassword(req.Password)

// 3. Create database params from DTO
result, _ := h.Queries.CreateUser(ctx, store.CreateUserParams{
    Username: req.Username,        // From DTO
    Email:    req.Email,           // From DTO
    Password: hashedPassword,      // Processed from DTO
    Created:  sql.NullTime{...},   // Set by server
    Updated:  sql.NullTime{...},   // Set by server
})
```

## Future DTOs

As the API grows, you'll add more DTOs:

### Login Request (Future)
```go
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

### Update User Request (Future)
```go
type UpdateUserRequest struct {
    Username string `json:"username,omitempty"`
    Email    string `json:"email,omitempty"`
}
```

### User Response (Future)
```go
type UserResponse struct {
    ID       int32  `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Created  string `json:"created"`
    // Password excluded for security
}
```

## Key Learning Points

1. **Separation of Concerns**: DTOs separate API from database
2. **Security**: Control what data clients can send/receive
3. **Validation**: Can add validation tags (future enhancement)
4. **JSON Tags**: Map between JSON and Go struct fields
5. **Plain Types**: DTOs use simple types; models use database types
6. **No Business Logic**: DTOs are just data structures
7. **Decode Pattern**: Always decode request body into DTO first
