# TODO - Future Improvements

A list of potential improvements and features to add to the REST API as you continue learning Go.

## ✅ Completed

- ✅ User registration endpoint
- ✅ Password hashing with bcrypt
- ✅ Password comparison function (ready for login)
- ✅ Database integration with PostgreSQL
- ✅ Type-safe SQL queries with sqlc
- ✅ Standardized error responses
- ✅ Standardized success responses
- ✅ Environment-based configuration
- ✅ Health check endpoint
- ✅ Request DTOs for user creation
- ✅ JWT token generation function (prepared for auth)

---

## Authentication & Security

- [ ] Implement user login endpoint
  - Use existing `ComparePassword()` function
  - Use existing `GenerateJWT()` function
  - Return JWT token to client
- [ ] Create authentication middleware
  - Verify JWT tokens from Authorization header
  - Extract user info from token
  - Add user context to requests
- [ ] Add refresh token functionality
  - Generate refresh tokens on login
  - Store refresh tokens in database
  - Endpoint to exchange refresh token for new access token
- [ ] Implement password reset flow
  - Generate password reset token
  - Send reset email (requires email service)
  - Verify token and update password
- [ ] Add email verification
  - Generate verification token on registration
  - Send verification email
  - Verify endpoint to confirm email
- [ ] Rate limiting to prevent brute force attacks
  - Limit login attempts per IP
  - Add middleware for rate limiting
- [ ] Add CORS middleware for cross-origin requests
  - Configure allowed origins
  - Set CORS headers

---

## User Management

- [ ] Get user profile endpoint (GET `/user/:id`)
  - Use existing `GetUser()` query
  - Exclude password from response
- [ ] Update user profile endpoint (PUT `/user/:id`)
  - Create `UpdateUserRequest` DTO
  - Add `UpdateUser` SQL query
  - Regenerate sqlc code
- [ ] Delete user endpoint (DELETE `/user/:id`)
  - Add `DeleteUser` SQL query
  - Regenerate sqlc code
- [ ] List users with pagination (GET `/users?page=1&limit=10`)
  - Modify existing `ListUsers()` query to support LIMIT/OFFSET
  - Add pagination params to request
- [ ] User search functionality
  - Search by username or email
  - Add SQL query with WHERE clause
- [ ] Change password endpoint
  - Verify old password with `ComparePassword()`
  - Hash new password
  - Update in database

---

## Validation & Error Handling

- [ ] Add request validation middleware
  - Validate required fields
  - Check field types and formats
- [ ] Validate email format
  - Use regex or email validation library
  - Return 400 if invalid
- [ ] Validate password strength requirements
  - Minimum length
  - Require special characters/numbers
  - Return clear error messages
- [ ] Add custom validation errors with field details
  - Return which field failed validation
  - Include helpful error messages
- [ ] Centralized error logging
  - Log all errors to file or logging service
  - Include request context (IP, user, timestamp)
- [ ] Better database error handling
  - Detect unique constraint violations
  - Return user-friendly messages ("Email already exists")
  - Handle foreign key violations

---

## Middleware

- [ ] Request logging middleware
  - Log all incoming requests
  - Include method, path, status code, response time
- [ ] Recovery middleware (panic recovery)
  - Catch panics in handlers
  - Return 500 error instead of crashing
  - Log stack traces
- [ ] Request ID middleware for tracing
  - Generate unique ID for each request
  - Add to logs and responses
  - Track requests across services
- [ ] Response time tracking
  - Measure handler execution time
  - Add to response headers or logs
- [ ] Content-Type validation
  - Ensure POST/PUT requests have JSON content-type
  - Return 415 (Unsupported Media Type) if not

---

## Database & Queries

- [ ] Add database migrations tool (like golang-migrate)
  - Version control database schema
  - Apply migrations automatically on startup
  - Rollback support
- [ ] Implement soft deletes for users
  - Add `deleted_at` column to users table
  - Modify queries to exclude soft-deleted records
  - Create restore endpoint
- [ ] Add database indexes for performance
  - Index on `users.email` for faster lookups
  - Index on `users.username`
  - Analyze query performance
- [ ] Implement database transactions for complex operations
  - Use `store.WithTx()` for multi-step operations
  - Ensure atomicity (all or nothing)
- [ ] Add prepared statements caching
  - Use `store.Prepare()` instead of `store.New()`
  - Improves query performance
- [ ] Connection pool optimization
  - Configure `SetMaxOpenConns()`
  - Configure `SetMaxIdleConns()`
  - Configure `SetConnMaxLifetime()`

---

## Testing

- [ ] Unit tests for handlers
  - Mock database queries
  - Test success and error cases
  - Use `httptest` package
- [ ] Unit tests for utilities
  - Test password hashing and comparison
  - Test JWT generation and parsing
  - Test response helpers
- [ ] Integration tests for API endpoints
  - Test full request/response cycle
  - Use test database
  - Test authentication flow
- [ ] Mock database for testing
  - Create mock implementation of `store.Queries`
  - Use interfaces for dependency injection
- [ ] Test coverage reporting
  - Run `go test -cover`
  - Aim for >80% coverage
  - Generate coverage HTML reports
- [ ] Load testing / benchmarks
  - Use `go test -bench`
  - Test concurrent requests
  - Identify bottlenecks

---

## Documentation

- [ ] Add API documentation (Swagger/OpenAPI)
  - Generate OpenAPI spec
  - Use tools like swaggo
  - Host interactive docs
- [ ] Add example requests/responses for all endpoints
  - Document in README or separate API.md
  - Include curl examples
  - Show error responses
- [ ] Create Postman collection
  - Export collection for testing
  - Include environment variables
  - Share with team members

---

## Code Quality

- [ ] Add linting configuration (golangci-lint)
  - Configure linters (.golangci.yml)
  - Run in CI/CD pipeline
  - Fix linter warnings
- [ ] Set up CI/CD pipeline
  - GitHub Actions or GitLab CI
  - Run tests on push
  - Check code coverage
  - Deploy on merge to main
- [ ] Add pre-commit hooks
  - Run linters before commit
  - Run tests before push
  - Format code automatically
- [ ] Improve error messages for debugging
  - Include context in errors
  - Use structured logging
  - Add error codes
- [ ] Add structured logging (replace fmt.Println)
  - Use `log/slog` (Go 1.21+) or `logrus`
  - JSON-formatted logs
  - Different log levels (debug, info, warn, error)
  - Include request context
- [ ] Add health check for database connectivity
  - Enhance `/health` endpoint
  - Check database connection
  - Return database status
  - Check other dependencies

---

## Performance & Scalability

- [ ] Implement caching (Redis)
  - Cache frequently accessed data
  - Reduce database queries
  - Set TTL for cache entries
- [ ] Add database query optimization
  - Analyze slow queries
  - Add indexes where needed
  - Use EXPLAIN to understand query plans
- [ ] Implement pagination helpers
  - Reusable pagination logic
  - Return total count
  - Include page metadata in responses
- [ ] Add response compression
  - Gzip responses
  - Reduce bandwidth usage
  - Add compression middleware
- [ ] Query result caching
  - Cache expensive queries
  - Invalidate on updates
  - Use Redis or in-memory cache

---

## DevOps & Deployment

- [ ] Create Dockerfile
  - Multi-stage build for smaller images
  - Use Alpine Linux base
  - Include only necessary files
- [ ] Docker Compose for local development
  - Include PostgreSQL service
  - Include Redis (if added)
  - Easy local setup
- [ ] Add health check endpoint improvements
  - Check database status
  - Check memory usage
  - Check disk space
  - Return detailed health info
- [ ] Environment-specific configurations
  - Separate configs for dev, staging, prod
  - Use config files or environment variables
  - Validate required configs on startup
- [ ] Add graceful shutdown handling
  - Listen for OS signals (SIGTERM, SIGINT)
  - Stop accepting new requests
  - Finish processing existing requests
  - Close database connections cleanly

---

## Advanced Features

- [ ] File upload for user avatars
  - Accept multipart/form-data
  - Store files in S3 or local storage
  - Validate file types and sizes
  - Generate thumbnails
- [ ] Email service integration
  - Use SendGrid, Mailgun, or AWS SES
  - Send verification emails
  - Send password reset emails
  - Email templates
- [ ] API versioning strategy (v1, v2)
  - Add `/api/v1/` prefix
  - Support multiple versions
  - Deprecation notices
- [ ] Role-based access control (RBAC)
  - Add `role` field to users table
  - Create roles (admin, user, moderator)
  - Middleware to check permissions
  - Protect endpoints by role
- [ ] OAuth integration (Google, GitHub login)
  - OAuth 2.0 flow
  - Link OAuth accounts to users
  - Generate JWT after OAuth login
  - Support multiple providers

---

## Notes

Items are organized by category for easier navigation. Start with the features that interest you most or follow the guide you're learning from.

Each completed item strengthens your understanding of Go, APIs, and backend development patterns.
