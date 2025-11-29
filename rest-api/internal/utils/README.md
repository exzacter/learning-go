# Utilities

This package contains reusable utility functions used across the API.

## Files

- `errorresponse.go` - Standardized error response helper
- `successresponse.go` - Standardized success response helper
- `passwordutil.go` - Password hashing and comparison
- `jwt.go` - JWT token generation and validation

## Error Response (`errorresponse.go`)

### ErrorResponse Struct
```go
type ErrorResponse struct {
    Message string `json:"message"`
}
```

### RespondWithError Function
```go
func RespondWithError(w http.ResponseWriter, code int, message string)
```

**Purpose**: Send standardized JSON error responses

**Parameters**:
- `w` - HTTP response writer
- `code` - HTTP status code (400, 500, etc.)
- `message` - Error message

**Usage in Handlers**:
```go
utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
```

**JSON Output**:
```json
{
    "message": "Invalid request payload"
}
```

**Flow**:
1. Sets `Content-Type: application/json` header
2. Sets HTTP status code with `w.WriteHeader(code)`
3. Encodes `ErrorResponse` struct to JSON
4. Sends to client

---

## Success Response (`successresponse.go`)

### SuccessResponse Struct
```go
type SuccessResponse struct {
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`  // Optional data
}
```

**`omitempty` tag**: Excludes `data` field from JSON if it's nil/empty.

### RespondWithSuccess Function
```go
func RespondWithSuccess(w http.ResponseWriter, code int, message string, data interface{})
```

**Purpose**: Send standardized JSON success responses

**Parameters**:
- `w` - HTTP response writer
- `code` - HTTP status code (200, 201, etc.)
- `message` - Success message
- `data` - Optional data to include (can be nil)

**Usage in Handlers**:
```go
utils.RespondWithSuccess(w, http.StatusCreated, "user created", req.Username)
utils.RespondWithSuccess(w, http.StatusOK, "operation successful", nil)
```

**JSON Output** (with data):
```json
{
    "message": "user created",
    "data": "testuser"
}
```

**JSON Output** (without data):
```json
{
    "message": "operation successful"
}
```

---

## Password Utilities (`passwordutil.go`)

### HashPassword Function
```go
func HashPassword(password string) (string, error)
```

**Purpose**: Hash passwords using bcrypt before storing in database

**Parameters**:
- `password` - Plain text password

**Returns**:
- Hashed password string
- Error if hashing fails

**Usage**:
```go
hashedPassword, err := utils.HashPassword(req.Password)
```

**Implementation**:
- Uses `golang.org/x/crypto/bcrypt`
- Uses `bcrypt.DefaultCost` (cost factor 10)
- Returns base64-encoded hash string

**Security**:
- **Never store plain text passwords**
- Bcrypt is slow by design (prevents brute force)
- Cost factor can be increased for more security

---

### ComparePassword Function
```go
func ComparePassword(storedPassword, providedPassword string) bool
```

**Purpose**: Verify a password against stored hash (for login)

**Parameters**:
- `storedPassword` - Hashed password from database
- `providedPassword` - Plain text password from user

**Returns**:
- `true` if password matches
- `false` if password doesn't match

**Usage** (for future login endpoint):
```go
if !utils.ComparePassword(user.Password, loginRequest.Password) {
    return errors.New("invalid password")
}
```

**Implementation**:
- Uses `bcrypt.CompareHashAndPassword()`
- Returns `true` if `err == nil` (passwords match)
- Constant-time comparison (prevents timing attacks)

---

## JWT Utilities (`jwt.go`)

### Claims Struct
```go
type Claims struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    jwt.StandardClaims  // Embeds standard JWT claims
}
```

**StandardClaims** includes:
- `ExpiresAt` - Token expiration timestamp
- `Issuer` - Who issued the token
- `Subject`, `Audience`, etc.

---

### GenerateJWT Function
```go
func GenerateJWT(userID int64, username string, secretKey []byte) (string, error)
```

**Purpose**: Generate JWT token for authenticated users (for future login endpoint)

**Parameters**:
- `userID` - User's database ID
- `username` - User's username
- `secretKey` - Secret key for signing token (from environment variable)

**Returns**:
- Signed JWT token string
- Error if generation fails

**Token Contents**:
- Custom claims: `user_id`, `username`
- Standard claims: `expires_at`, `issuer`

**Token Expiration**: 24 hours from generation

**Issuer**: "Project Harbinger"

**Signing Method**: HS256 (HMAC-SHA256)

**Usage** (for future login endpoint):
```go
token, err := utils.GenerateJWT(user.ID, user.Username, []byte(secretKey))
// Send token to client
```

**Note**: Not yet used in the API; prepared for authentication implementation.

---

## How Utils Connect to Handlers

### Error Handling
```go
// In handler
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
    return
}
```

### Success Responses
```go
// In handler
utils.RespondWithSuccess(w, http.StatusCreated, "user created", username)
```

### Password Hashing
```go
// In CreateUserHandler
hashedPassword, err := utils.HashPassword(req.Password)
// Store hashedPassword in database
```

### JWT Generation (Future)
```go
// In future LoginHandler
token, err := utils.GenerateJWT(user.ID, user.Username, secretKey)
// Return token to client
```

## Key Learning Points

1. **DRY Principle**: Utilities prevent code duplication
2. **Standardization**: Consistent response format across all endpoints
3. **Security**: Password hashing is isolated and reusable
4. **interface{}**: Used for flexible data types in responses
5. **omitempty**: Excludes empty fields from JSON
6. **bcrypt**: Industry-standard password hashing
7. **JWT**: Stateless authentication tokens
8. **Separation**: Utilities don't know about HTTP handlers (loose coupling)
