# CGO Limitations & Docker Impact Guide

## What is CGO?

CGO is Go's foreign function interface (FFI) that allows Go programs to call C code. When you import `github.com/mattn/go-sqlite3`, you're using a package that wraps the C SQLite library.

```go
import _ "github.com/mattn/go-sqlite3"  // This requires CGO
```

The `go-sqlite3` package binds the C SQLite library to Go, meaning **your Go binary depends on compiled C code**.

---

## CGO Limitations

### 1. Cross-Compilation is Painful

**The Problem:**
Pure Go can cross-compile trivially:
```bash
# Pure Go - works instantly
GOOS=linux GOARCH=amd64 go build -o myapp .
GOOS=windows GOARCH=amd64 go build -o myapp.exe .
GOOS=darwin GOARCH=arm64 go build -o myapp .
```

With CGO, you need a C cross-compiler for each target platform:
```bash
# CGO - requires platform-specific C compiler
CGO_ENABLED=1 CC=x86_64-linux-gnu-gcc GOOS=linux go build  # Need Linux GCC
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows go build  # Need MinGW
```

**Impact:** Building for multiple platforms requires either:
- Multiple build machines (one per target OS)
- Complex cross-compiler toolchains
- Docker with platform-specific images
- Tools like [xgo](https://github.com/karalabe/xgo) or Zig as CC

---

### 2. Build Environment Complexity

**Required tools for CGO:**
- GCC or Clang compiler
- C standard library headers
- Platform-specific development libraries

```bash
# Alpine Linux requirements
apk add gcc musl-dev

# Debian/Ubuntu requirements
apt-get install gcc libc6-dev

# macOS requirements
xcode-select --install
```

**Without CGO (pure Go):**
- Just the Go compiler
- No additional dependencies

---

### 3. Binary Size & Dependencies

| Build Type | Binary Type | Size Impact | Runtime Dependencies |
|------------|-------------|-------------|---------------------|
| `CGO_ENABLED=0` | Static (pure Go) | Smaller | None |
| `CGO_ENABLED=1` | Dynamic | Larger | libc, libpthread, etc. |
| `CGO_ENABLED=1` + static flags | Static with C | Largest | None (but bigger binary) |

---

### 4. Security Risks

**Memory Safety Bypass:**
> "CGO must be used with extreme caution because you are trusting a completely external dependency written in an unsafe language. The Go memory safety net is not there to save you if there are bugs or malicious routines lurking in that external code."

**Specific risks:**
- C code doesn't have Go's bounds checking
- Memory corruption bugs in C libraries affect your Go app
- Buffer overflows, use-after-free, etc. become possible
- Go's garbage collector can't manage C-allocated memory

**ASLR/PIE Concerns:**
> "Anything compiled with Golang will not have ASLR/PIE. If the process imports a C library, it exposes itself to possible issues."

---

### 5. Build Time

CGO builds are significantly slower:
- C compiler must run
- C code must be compiled
- Linking is more complex

```bash
# Pure Go build: ~2-5 seconds
CGO_ENABLED=0 go build .

# CGO build: ~10-30 seconds (depends on C code size)
CGO_ENABLED=1 go build .
```

---

## Docker-Specific Impact

### Problem 1: glibc vs musl Incompatibility

**The Core Issue:**
- Most Linux systems use **glibc** (GNU C Library)
- Alpine Linux uses **musl** (smaller, different ABI)
- Binaries compiled with glibc **won't run on Alpine** (and vice versa)

```
Build on Ubuntu (glibc) → Run on Alpine (musl) = FAIL
Build on Alpine (musl) → Run on Ubuntu (glibc) = FAIL
```

**Error you'll see:**
```
/bin/sh: ./myapp: not found
# or
standard_init_linux.go: exec user process caused "no such file or directory"
```

---

### Problem 2: Multi-Stage Build Issues

A typical multi-stage Dockerfile fails with CGO:

```dockerfile
# THIS WILL FAIL WITH go-sqlite3
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp .

FROM alpine:latest  # Uses musl, not glibc!
COPY --from=builder /app/myapp /myapp
CMD ["/myapp"]
# ERROR: binary linked against glibc, but Alpine has musl
```

---

### Problem 3: Scratch/Distroless Images Don't Work

```dockerfile
# THIS WILL FAIL - no libc at all!
FROM golang:1.21 AS builder
RUN CGO_ENABLED=1 go build -o myapp .

FROM scratch
COPY --from=builder /app/myapp /myapp
CMD ["/myapp"]
# ERROR: dynamic binary needs libc which doesn't exist in scratch
```

---

## Docker Solutions for go-sqlite3

### Solution 1: Build and Run on Same Base (Simplest)

```dockerfile
# Use the same libc for build and runtime
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Install CGO requirements
RUN apk add --no-cache gcc musl-dev

ENV CGO_ENABLED=1
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o myapp .

# Use same Alpine base for runtime
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/myapp /myapp
COPY --from=builder /app/data.db /data.db
CMD ["/myapp"]
```

---

### Solution 2: Static Binary for Scratch (Smallest Image)

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Install build tools
RUN apk add --no-cache gcc musl-dev

ENV CGO_ENABLED=1
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build fully static binary
RUN go build -ldflags='-s -w -extldflags "-static"' -o myapp .

# Scratch image - smallest possible
FROM scratch
COPY --from=builder /app/myapp /myapp
COPY --from=builder /app/data.db /data.db

# Need CA certs for HTTPS calls
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/myapp"]
```

**Build flags explained:**
- `-s` - Strip symbol table (smaller binary)
- `-w` - Strip DWARF debug info (smaller binary)
- `-extldflags "-static"` - Force static linking of C code

---

### Solution 3: Debian-based (Most Compatible)

```dockerfile
FROM golang:1.21 AS builder
WORKDIR /app

ENV CGO_ENABLED=1
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o myapp .

# Debian slim has glibc - compatible with builder
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/myapp /myapp
COPY --from=builder /app/data.db /data.db
CMD ["/myapp"]
```

---

### Solution 4: Distroless (Security-Focused)

```dockerfile
FROM golang:1.21 AS builder
WORKDIR /app

ENV CGO_ENABLED=1
COPY . .
RUN go build -o myapp .

# Distroless has glibc but minimal attack surface
FROM gcr.io/distroless/base-debian12
COPY --from=builder /app/myapp /myapp
COPY --from=builder /app/data.db /data.db
CMD ["/myapp"]
```

---

## Image Size Comparison

| Base Image | Approximate Size | CGO Compatible | Notes |
|------------|------------------|----------------|-------|
| `scratch` | ~10-15 MB | Only with static build | No shell, no debugging |
| `alpine` | ~15-20 MB | Yes (musl) | Must build on Alpine |
| `distroless` | ~20-30 MB | Yes (glibc) | No shell, secure |
| `debian:slim` | ~80-100 MB | Yes (glibc) | Most compatible |
| `ubuntu` | ~80-120 MB | Yes (glibc) | Familiar, larger |

---

## Pure Go Alternative: No CGO Required

### Option 1: modernc.org/sqlite

A CGO-free SQLite implementation (C code transpiled to Go):

```go
import (
    "database/sql"
    _ "modernc.org/sqlite"  // Drop-in replacement, no CGO!
)

func main() {
    db, err := sql.Open("sqlite", "./data.db")  // Note: "sqlite" not "sqlite3"
    // ... rest of code works the same
}
```

**Pros:**
- No CGO required
- Cross-compiles trivially
- Works with scratch/Alpine/any image
- Same `database/sql` interface

**Cons:**
- ~2x slower on INSERTs
- 10-100% slower on SELECTs
- Larger binary size (SQLite compiled into Go)

---

### Option 2: ncruces/go-sqlite3 (WASM-based)

Uses WebAssembly instead of CGO:

```go
import (
    "database/sql"
    _ "github.com/ncruces/go-sqlite3/driver"
    _ "github.com/ncruces/go-sqlite3/embed"
)
```

**Pros:**
- No CGO
- Good performance (sometimes better than modernc)
- Cross-platform

---

### Performance Comparison

| Driver | INSERT Performance | SELECT Performance | CGO Required |
|--------|-------------------|-------------------|--------------|
| mattn/go-sqlite3 | Fastest (baseline) | Fastest (baseline) | Yes |
| modernc.org/sqlite | ~2x slower | 10-100% slower | No |
| ncruces/go-sqlite3 | ~1.5x slower | ~1.2x slower | No |

**Recommendation from benchmarks:**
> "If your workload has solely small datasets (i.e. small business apps) the tradeoff allowing you to avoid CGO could be worth it. Otherwise if you care strongly about performance you'll be better off with mattn/go-sqlite3."

---

## Migration: go-sqlite3 to modernc.org/sqlite

### Step 1: Update go.mod

```bash
go get modernc.org/sqlite
```

### Step 2: Change Import

```go
// Before (CGO required)
import _ "github.com/mattn/go-sqlite3"
db, _ := sql.Open("sqlite3", "./data.db")

// After (pure Go)
import _ "modernc.org/sqlite"
db, _ := sql.Open("sqlite", "./data.db")  // Note: "sqlite" not "sqlite3"
```

### Step 3: Update Dockerfile

```dockerfile
# Now works with ANY base image!
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o myapp .  # CGO disabled!

FROM scratch
COPY --from=builder /app/myapp /myapp
CMD ["/myapp"]
```

---

## Decision Matrix

| Requirement | Use mattn/go-sqlite3 (CGO) | Use modernc.org/sqlite (Pure Go) |
|-------------|---------------------------|----------------------------------|
| Maximum performance | Yes | No |
| Simple Docker builds | No | Yes |
| Cross-compilation | Difficult | Easy |
| Scratch/distroless images | Complex | Easy |
| Small business app | Either | Recommended |
| High-throughput workload | Recommended | Acceptable |
| CI/CD simplicity | No | Yes |
| Security (no C code) | No | Yes |

---

## Quick Reference: Docker Build Commands

```bash
# Build with CGO (for go-sqlite3)
docker build --platform linux/amd64 -t myapp:cgo .

# Build for ARM (requires cross-compiler or buildx)
docker buildx build --platform linux/arm64 -t myapp:arm64 .

# Test locally
docker run --rm -p 8080:8080 myapp:cgo
```

---

## Recommended Dockerfile for Your Project

Since you're using `mattn/go-sqlite3`, here's the recommended Dockerfile:

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Install CGO dependencies
RUN apk add --no-cache gcc musl-dev

# Enable CGO
ENV CGO_ENABLED=1

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN go build -ldflags='-s -w -extldflags "-static"' -o /api ./cmd/api

# Runtime stage
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /api /app/api

# Create data directory
RUN mkdir -p /app/data

EXPOSE 8080
CMD ["/app/api"]
```

---

## Sources

- [go-sqlite3 CGO Requirements](https://github.com/mattn/go-sqlite3/issues/855)
- [Golang w/SQLite3 + Docker Scratch Image](https://7thzero.com/blog/golang-w-sqlite3-docker-scratch-image)
- [Alpine Linux & Docker: glibc vs musl](https://iifx.dev/en/articles/377343174)
- [Using CGO bindings under Alpine, CentOS and Ubuntu](https://www.x-cellent.com/posts/cgo-bindings)
- [SQLite in Go, with and without CGO](https://datastation.multiprocess.io/blog/2022-05-12-sqlite-in-go-with-and-without-cgo.html)
- [modernc.org/sqlite Package](https://pkg.go.dev/modernc.org/sqlite)
- [go-sqlite-bench Benchmarks](https://github.com/cvilsmeier/go-sqlite-bench)
- [You don't need CGO to use SQLite](https://hiandrewquinn.github.io/til-site/posts/you-don-t-need-cgo-to-use-sqlite-in-your-go-binary/)
- [Go Security Cheatsheet - Snyk](https://snyk.io/blog/go-security-cheatsheet-for-go-developers/)
- [Docker Images: Details Specific to Different Languages](https://www.ardanlabs.com/blog/2020/02/docker-images-part2-details-specific-to-different-languages.html)
