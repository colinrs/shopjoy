# SHOPJOY PROJECT KNOWLEDGE BASE

**Generated:** 2026-03-17
**Commit:** a9169e7
**Branch:** main
**Go Version:** 1.24.0

## OVERVIEW

E-commerce platform with admin management and shop APIs. Built with go-zero microservices framework using DDD (Domain Driven Design) and repository pattern.

## STRUCTURE

```
./
├── admin/              # Admin management API service
├── shop/               # Shop/e-commerce API service
├── pkg/                # Shared packages
│   ├── cache/          # Cache abstraction (Redis, Ristretto)
│   ├── client/         # Client utilities (etcd)
│   ├── code/           # Error codes and definitions
│   ├── codec/          # Serialization (sonic JSON)
│   ├── gosafe/         # Safe goroutine utilities
│   ├── httpc/          # HTTP client utilities
│   ├── httpy/          # HTTP parsing utilities
│   ├── infra/          # Infrastructure (DB, Redis, metrics)
│   ├── response/       # HTTP response handlers
│   ├── snowflake/      # ID generation
│   └── utils/          # General utilities
├── go.mod              # Go module definition
└── Makefile            # Build automation
```

## COMMANDS

### Root Makefile (build all services)
```bash
# Generate API code for all services
make gen-go-api
make api

# Build all services
make build

# Run linter
golangci-lint run --timeout=10m
```

### Service Makefile (shop/ or admin/)
```bash
# Format API definitions
cd shop && make format

# Generate Go code from API definitions
cd shop && make gen-go-api
cd shop && make api          # format + gen

# Build service binary
cd shop && make build        # outputs bin/shop

# Run single service
cd shop && ./bin/shop -f etc/shop-api.yaml
```

### Testing
```bash
# Run all tests
go test ./...

# Run specific test
go test ./pkg/cache/... -v
go test -run TestRistrettoCache ./pkg/cache/... -v

# Run with race detection
go test -race ./...
```

### Linting
```bash
# Install golangci-lint (if needed)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run --timeout=10m

# Run on specific package
golangci-lint run ./pkg/...
```

### go-zero Code Generation
```bash
# Install goctl (required for code generation)
go install github.com/zeromicro/go-zero/tools/goctl@latest

# Generate API code from .api files
goctl api go --api ./desc/shop.api --dir ./ --style=go_zero

# Format API definitions
goctl api format --dir ./desc

# Generate Swagger docs
goctl api plugin -plugin goctl-swagger="swagger -filename shop.json" -api ./desc/shop.api -dir swagger
```

## CONVENTIONS

### Code Organization (go-zero pattern)
- `desc/*.api` - API definition files (go-zero syntax)
- `internal/handler/` - HTTP handlers (auto-generated + custom)
- `internal/logic/` - Business logic (scaffolded by goctl)
- `internal/svc/` - Service context (dependency injection)
- `internal/config/` - Configuration structs
- `internal/types/` - Request/response types (auto-generated)
- `etc/*.yaml` - Service configuration files

### Naming
- **Files**: `snake_case.go` for handlers, `camelCase.go` for logic
- **Types**: PascalCase (e.g., `ShopLogic`, `ServiceContext`)
- **Interfaces**: Verb-like with `er` suffix where appropriate
- **Private**: lowercase or underscore prefix

### Error Handling
- Use `pkg/code` error definitions
- Custom `Err` struct with HTTP code, business code, message
- Silent degradation for cache operations
- Always log errors with `logx.WithContext(ctx).Errorf()`

```go
// Return typed errors
return nil, code.ErrParam

// Create custom error
code.NewErr(code.WithMsg("custom message"), code.WithHTTPCode(400))
```

### Imports Order
```go
import (
    // Standard library
    "context"
    "net/http"
    
    // Third-party
    "github.com/zeromicro/go-zero/core/logx"
    
    // Project internal
    "github.com/colinrs/shopjoy/pkg/code"
    "github.com/colinrs/shopjoy/shop/internal/svc"
)
```

### Context Usage
- Always pass `context.Context` as first parameter
- Use `logx.WithContext(ctx)` for logging
- Extract trace info from context for external calls

### API Definition (.api files)
```go
type Request {
    Name string `path:"name,options=you|me"`  // path param with validation
    Age  int    `form:"age"`                   // query/form param
}

type Response {
    Message string `json:"message"`
}

service shop-api {
    @handler ShopHandler
    get /from/:name (Request) returns (Response)
}
```

### Money Handling
- **ALWAYS** use `github.com/shopspring/decimal` for all monetary calculations
- **NEVER** use `float64` or `int` for money - precision loss is unacceptable
- Store in database as `DECIMAL(19,4)` or string if needed

```go
import "github.com/shopspring/decimal"

// Good
price := decimal.NewFromFloat(99.99)
total := price.Mul(decimal.NewFromInt(quantity))

// Bad - precision loss
price := 99.99  // float64
total := price * float64(quantity)
```

### Time Handling
- **ALWAYS** store and return UTC time in backend
- **NEVER** store local time in database
- Frontend is responsible for timezone conversion and display
- Use `time.Now().UTC()` for current time

```go
// Good - store UTC
now := time.Now().UTC()
createdAt := now.Format("2006-01-02 15:04:05")

// API response - UTC format
type Response {
    CreatedAt string `json:"created_at"`  // "2024-01-15 08:30:00"
}

// Bad - storing local time
now := time.Now()  // Depends on server timezone
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Add new API endpoint | `*/desc/*.api` | Define types and service |
| Implement business logic | `*/internal/logic/` | Scaffolded, safe to edit |
| Custom request handling | `*/internal/handler/` | Rarely needed |
| Service dependencies | `*/internal/svc/` | Wire dependencies |
| Configuration | `*/internal/config/`, `*/etc/*.yaml` | YAML + struct mapping |
| Error codes | `pkg/code/` | Centralized error definitions |
| Cache operations | `pkg/cache/` | See pkg/cache/AGENTS.md |
| Database | `pkg/infra/db.go` | GORM MySQL setup |
| Redis | `pkg/infra/redis.go` | go-redis client |
| HTTP responses | `pkg/response/` | Standard response format |

## ANTI-PATTERNS

- **DO NOT** edit `internal/types/types.go` or `internal/handler/routes.go` - auto-generated
- **DO NOT** hardcode cache keys - use `cache_key.go` utilities
- **DO NOT** use cache as primary storage
- **NEVER** ignore context cancellation in long operations
- **AVOID** blocking calls without timeout in handlers
- **NEVER** modify goctl-scaffolded files that have "Safe to edit" comments outside logic blocks

## DDD ARCHITECTURE

### Repository Pattern

**ALWAYS** pass DB/Tx as function parameters, NEVER store in Repository struct:

```go
// CORRECT - DB passed as parameter
type Repository interface {
    Create(ctx context.Context, db *gorm.DB, entity *Entity) error
    Update(ctx context.Context, db *gorm.DB, entity *Entity) error
    FindByID(ctx context.Context, db *gorm.DB, id int64) (*Entity, error)
}

// WRONG - DB stored in struct
type BadRepository struct {
    db *gorm.DB  // Don't do this!
}
```

### Transaction Control

**ALWAYS** start transactions in Application/Domain layer, NEVER in Repository:

```go
// CORRECT - Transaction started in Application Service
func (s *Service) CreateOrder(ctx context.Context, req CreateOrderRequest) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Repository operations use transaction
        if err := s.orderRepo.Create(ctx, tx, order); err != nil {
            return err
        }
        if err := s.inventoryRepo.Deduct(ctx, tx, item); err != nil {
            return err
        }
        return nil
    })
}

// WRONG - Transaction inside Repository
func (r *BadRepository) Create(ctx context.Context, entity *Entity) error {
    return r.db.Transaction(func(tx *gorm.DB) error {  // Don't do this!
        // ...
    })
}
```

### Layer Dependencies

```
Handler (Interface Layer)
    ↓ depends on
Application Service (Application Layer)
    ↓ depends on
Domain Entity/Repository Interface (Domain Layer)
    ↓ depends on
Repository Implementation (Infrastructure Layer)
```

- **Domain Layer** has NO external dependencies
- **Application Layer** orchestrates use cases
- **Infrastructure Layer** implements Repository interfaces

## FRAMEWORK NOTES

- **go-zero**: Microservices framework with API code generation
- **goctl**: CLI tool for generating boilerplate from .api files
- **GORM**: ORM for database operations
- **go-redis**: Redis client
- **Ristretto**: In-memory cache library
- **sonic**: High-performance JSON library (ByteDance)

## CHILD AGENTS

- `pkg/cache/AGENTS.md` - Cache abstraction layer

---

*Generated by init-deep for agentic coding assistance*
