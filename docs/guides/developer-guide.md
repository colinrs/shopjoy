# ShopJoy Developer Guide

> **Version:** 1.0
> **Last Updated:** 2026-03-27
> **Audience:** Backend and Frontend Developers

---

## Table of Contents

1. [Getting Started](#getting-started)
2. [Project Structure](#project-structure)
3. [Development Environment](#development-environment)
4. [API Development Workflow](#api-development-workflow)
5. [Domain-Driven Design](#domain-driven-design)
6. [Database Operations](#database-operations)
7. [Error Handling](#error-handling)
8. [Testing](#testing)
9. [Deployment](#deployment)
10. [Best Practices](#best-practices)

---

## Getting Started

### Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.24+ | Backend development |
| Node.js | 18+ | Frontend development |
| MySQL | 8.0+ | Database |
| Redis | 7.0+ | Caching |
| Docker | Latest | Containerization |
| goctl | Latest | Code generation |

### Clone and Setup

```bash
# Clone the repository
git clone https://github.com/colinrs/shopjoy.git
cd shopjoy

# Install Go dependencies
go mod download

# Install frontend dependencies
cd shop-admin && npm install
cd ../joy && npm install

# Copy environment configuration
cp .env.example .env
```

---

## Project Structure

### Backend (`admin/` and `shop/`)

```
shopjoy/
├── admin/                          # Admin API service
│   ├── desc/                       # API definition files (*.api)
│   │   ├── admin.api               # Main entry (imports others)
│   │   ├── auth.api               # Authentication
│   │   ├── product.api             # Product management
│   │   └── ...
│   │
│   └── internal/
│       ├── handler/                # HTTP handlers (auto-generated)
│       ├── logic/                  # Business logic
│       ├── types/                  # Request/Response types (auto-generated)
│       ├── svc/                    # Service context
│       ├── config/                 # Configuration
│       │
│       ├── domain/                 # Domain layer
│       │   ├── user/
│       │   │   ├── entity.go       # User entity
│       │   │   └── repository.go   # Repository interface
│       │   ├── product/
│       │   │   ├── entity.go       # Product entity
│       │   │   ├── sku.go          # SKU entity
│       │   │   └── repository.go   # Repository interface
│       │   └── ...
│       │
│       ├── application/            # Application services
│       └── infrastructure/         # Infrastructure implementations
│
├── shop/                           # Shop API service (customer-facing)
│
├── pkg/                            # Shared packages
│   ├── code/                       # Error codes
│   │   └── code.go                # All error definitions
│   ├── domain/
│   │   └── shared/                 # Shared domain components
│   │       ├── tenant.go          # TenantID value object
│   │       ├── money.go            # Money value object
│   │       └── time.go             # Time utilities
│   └── infra/                      # Shared infrastructure
│
├── sql/                            # Database schema
│   ├── init.sql                    # Database initialization
│   ├── user/                       # User domain schema
│   │   └── schema.sql
│   ├── product/
│   │   └── schema.sql
│   └── ...
│
└── docs/                           # Documentation
```

### Frontend

```
shop-admin/                         # Admin dashboard (Vue 3 + Element Plus)
├── src/
│   ├── views/                      # Page components
│   ├── components/                 # Reusable components
│   ├── stores/                     # Pinia stores
│   ├── api/                        # API clients
│   └── router/                     # Vue Router config
│
joy/                                # Storefront (Vue 3 + Tailwind)
```

---

## Development Environment

### Docker Compose Setup

```yaml
# docker-compose.yml
version: '3.8'
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: shopjoy
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  admin:
    build: ./admin
    ports:
      - "8888:8888"
    depends_on:
      - mysql
      - redis

  shop:
    build: ./shop
    ports:
      - "8889:8889"
    depends_on:
      - mysql
      - redis

volumes:
  mysql_data:
```

### Start Development Environment

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f admin

# Stop services
docker-compose down
```

### Manual Development

```bash
# Start MySQL and Redis (using Homebrew on macOS)
brew services start mysql
brew services start redis

# Initialize database
mysql -u root -p < sql/init.sql

# Start Admin API
cd admin
go run ./cmd/admin.go -f etc/admin.yaml

# Start Shop API
cd shop
go run ./cmd/shop.go -f etc/shop.yaml
```

---

## API Development Workflow

### Step 1: Define API in .api File

Create or modify API definition in `admin/desc/{module}.api`:

```go
syntax = "v1"

info (
    title:   "Product API"
    desc:    "Product management endpoints"
    version: "v1"
)

type (
    CreateProductReq {
        Name       string `json:"name"`
        Price      int64  `json:"price"`
        CategoryID int64  `json:"category_id"`
    }

    CreateProductResp {
        ID int64 `json:"id"`
    }
)

@server (
    group:      products
    middleware: AuthMiddleware
)
service admin-api {
    @doc "Create a new product"
    @handler CreateProductHandler
    post /api/v1/products (CreateProductReq) returns (CreateProductResp)
}
```

### Step 2: Generate Code

```bash
# Generate code for admin service
cd admin && make api

# Generate code for shop service
cd shop && make api
```

This auto-generates:
- `internal/types/types.go` - Request/Response structs
- `internal/handler/*.go` - Handler stubs
- `internal/handler/routes.go` - Route registration

### Step 3: Implement Business Logic

Implement the handler logic in `internal/logic/`:

```go
// admin/internal/logic/product/create_product_logic.go
package product

import (
    "context"

    "github.com/colinrs/shopjoy/admin/internal/logic"
    "github.com/colinrs/shopjoy/admin/internal/svc"
    "github.com/colinrs/shopjoy/admin/internal/types"
    "github.com/colinrs/shopjoy/pkg/code"
)

type CreateProductLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
    return &CreateProductLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

func (l *CreateProductLogic) CreateProduct(req *types.CreateProductReq) (*types.CreateProductResp, error) {
    // Business logic here
    if req.Name == "" {
        return nil, code.ErrProductEmptyName
    }

    // Create product via domain service
    id, err := l.svcCtx.ProductService.Create(l.ctx, req)
    if err != nil {
        return nil, err
    }

    return &types.CreateProductResp{ID: id}, nil
}
```

### Step 4: Update Handler to Use Logic

The generated handler will call the logic:

```go
// admin/internal/handler/product.go (auto-generated, but logic is implemented separately)
// Handler calls logic via Logic.NewXxxLogic()

func (l *CreateProductLogic) CreateProduct(req *types.CreateProductReq) (*types.CreateProductResp, error) {
    // Logic implementation
}
```

### Step 5: Verify Build

```bash
# Build to verify compilation
cd admin && make build

# Or directly
go build -o bin/admin ./cmd/admin.go
```

---

## Domain-Driven Design

### Domain Layer Structure

Each domain follows consistent structure:

```
domain/
└── {domain_name}/
    ├── entity.go          # Main entity (aggregate root)
    ├── value_object.go    # Value objects
    ├── repository.go       # Repository interface
    └── service.go         # Domain service (if needed)
```

### Entity Example

```go
// admin/internal/domain/product/entity.go
package product

import (
    "time"
    "github.com/colinrs/shopjoy/pkg/domain/shared"
)

type Status int

const (
    StatusDraft Status = iota
    StatusOnSale
    StatusOffSale
)

type Product struct {
    ID        int64
    TenantID  shared.TenantID
    Name      string
    Price     Money              // Value object
    Status    Status
    Audit     shared.AuditInfo   // Embedded audit info
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (p *Product) PutOnSale() error {
    if p.Status != StatusDraft {
        return code.ErrProductInvalidStatusTransition
    }
    if p.Stock <= 0 {
        return code.ErrProductNoStock
    }
    p.Status = StatusOnSale
    return nil
}
```

### Value Objects

```go
// Money value object
type Money struct {
    Amount   int64   // In cents
    Currency string  // ISO 4217
}

func (m Money) Add(other Money) (Money, error) {
    if m.Currency != other.Currency {
        return Money{}, code.ErrProductCurrencyMismatch
    }
    return Money{Amount: m.Amount + other.Amount, Currency: m.Currency}, nil
}

// TenantID value object
type TenantID int64

func (t TenantID) Int64() int64 {
    return int64(t)
}
```

### Repository Pattern

```go
// Repository interface (defined in domain)
type Repository interface {
    Create(ctx context.Context, db *gorm.DB, product *Product) error
    FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Product, error)
    Update(ctx context.Context, db *gorm.DB, product *Product) error
    Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
}

// Implementation (in infrastructure)
type productRepository struct {
    db *gorm.DB
}

func (r *productRepository) Create(ctx context.Context, db *gorm.DB, product *Product) error {
    return db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Product, error) {
    var product Product
    err := db.WithContext(ctx).
        Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).
        First(&product).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, code.ErrProductNotFound
        }
        return nil, err
    }
    return &product, nil
}
```

---

## Database Operations

### Schema Management

Schema files are organized by domain:

```
sql/
├── init.sql                    # Main entry
├── user/
│   └── schema.sql              # User domain tables
├── product/
│   └── schema.sql              # Product domain tables
└── ...
```

### Table Naming Convention

- All table names: lowercase with underscores
- Primary key: `id` (BIGINT, auto-increment)
- Tenant isolation: `tenant_id` (BIGINT)
- Soft delete: `deleted_at` (BIGINT, Unix timestamp)
- Timestamps: `created_at`, `updated_at` (TIMESTAMP)

### Example Schema

```sql
CREATE TABLE products (
    id              BIGINT           PRIMARY KEY AUTO_INCREMENT,
    tenant_id       BIGINT           NOT NULL,
    name            VARCHAR(255)     NOT NULL,
    price_amount    BIGINT           NOT NULL,
    price_currency  VARCHAR(3)       NOT NULL DEFAULT 'CNY',
    stock           INT              NOT NULL DEFAULT 0,
    status          TINYINT          NOT NULL DEFAULT 0,
    category_id     BIGINT           NOT NULL,
    created_at      TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at      BIGINT          NULL,

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_category_id (category_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### Running Migrations

```bash
# Initialize database
mysql -u root -p < sql/init.sql

# For specific domain
mysql -u root -p shopjoy < sql/product/schema.sql
```

---

## Error Handling

### Error Code Structure

All errors are defined in `pkg/code/code.go`:

```go
var (
    // Module-specific errors (code range: 30xxx)
    ErrProductNotFound = &Err{
        HTTPCode: http.StatusNotFound,
        Code: 30012,
        Msg: "product not found",
    }
)
```

### Error Code Ranges

| Module | Range | Example |
|--------|-------|---------|
| Admin User | 10xxx | 10001, 10002 |
| User | 11xxx | 11001, 11002 |
| Product | 30xxx | 30001, 30012 |
| Order | 40xxx | 40001, 40002 |
| Payment | 50xxx | 50001, 50002 |
| Promotion | 80xxx | 80001, 80002 |
| Fulfillment | 120xxx | 120001, 120002 |

### Using Errors

```go
// Return error from domain/logic layer
func (l *ProductLogic) GetProduct(id int64) (*Product, error) {
    product, err := l.repo.FindByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, code.ErrProductNotFound
        }
        return nil, err
    }
    return product, nil
}

// Handle error in handler
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
    product, err := h.logic.GetProduct(id)
    if err != nil {
        // Error response is auto-handled by go-zero
        return
    }
    // Success response
}
```

---

## Testing

### Unit Tests

```bash
# Run unit tests
go test ./admin/internal/domain/product/... -v

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Example Test

```go
// admin/internal/domain/product/entity_test.go
package product

import (
    "testing"
)

func TestProduct_PutOnSale(t *testing.T) {
    product := &Product{
        ID:     1,
        Status: StatusDraft,
        Stock:  10,
    }

    err := product.PutOnSale()
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }

    if product.Status != StatusOnSale {
        t.Errorf("expected status OnSale, got %v", product.Status)
    }
}

func TestProduct_PutOnSale_NoStock(t *testing.T) {
    product := &Product{
        ID:     1,
        Status: StatusDraft,
        Stock:  0,
    }

    err := product.PutOnSale()
    if err != code.ErrProductNoStock {
        t.Errorf("expected ErrProductNoStock, got %v", err)
    }
}
```

### Integration Tests

```bash
# Run with test database
go test -tags=integration ./...
```

---

## Deployment

### Build for Production

```bash
# Build admin service
cd admin && make build

# Build shop service
cd shop && make build

# Build all
make build
```

### Docker Build

```dockerfile
# admin/Dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o bin/admin ./cmd/admin.go

FROM alpine:latest
COPY --from=builder /app/bin/admin /usr/local/bin/admin
COPY --from=builder /app/etc /etc/shopjoy
EXPOSE 8888
CMD ["admin"]
```

### Kubernetes Deployment

```yaml
# k8s/admin-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: admin-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: admin-api
  template:
    metadata:
      labels:
        app: admin-api
    spec:
      containers:
      - name: admin-api
        image: shopjoy/admin:latest
        ports:
        - containerPort: 8888
        env:
        - name: MYSQL_HOST
          valueFrom:
            secretKeyRef:
              name: shopjoy-secrets
              key: mysql-host
```

---

## Best Practices

### API Design

1. **Use consistent naming**: Plural nouns for collections (`/products`)
2. **Proper HTTP methods**: GET (read), POST (create), PUT (update), DELETE (remove)
3. **Pagination**: Always use `page` and `page_size` for lists
4. **Filtering**: Support query parameters for list filtering
5. **Error responses**: Consistent error format with codes

### Code Organization

1. **Three-layer separation**: Handler -> Logic -> Domain
2. **Repository pattern**: Abstract data access
3. **Value objects**: Use for monetary amounts, IDs, etc.
4. **Tenant isolation**: Always include tenant_id in queries

### Database

1. **Use transactions**: For multi-table operations
2. **Index wisely**: Index columns used in WHERE clauses
3. **Soft delete**: Use `deleted_at` instead of hard delete
4. **Audit trail**: Always set `created_at`, `updated_at`

### Error Handling

1. **Use code.ErrXxx**: Never use `errors.New()`
2. **Proper HTTP codes**: Match error severity to HTTP status
3. **Wrap appropriately**: Preserve error chain with `fmt.Errorf("...: %w", err)`
4. **User-friendly messages**: Provide clear error messages

### Security

1. **Validate input**: Always validate at API boundary
2. **Tenant isolation**: Verify tenant_id matches request context
3. **Authentication**: Use JWT with proper expiration
4. **Authorization**: Check permissions before operations

### Performance

1. **Use caching**: Redis for frequently accessed data
2. **Batch operations**: When possible, use batch APIs
3. **Connection pooling**: Configure appropriate pool sizes
4. **Async processing**: Use message queues for heavy operations

---

## Common Tasks

### Adding a New API Endpoint

1. Edit `admin/desc/{module}.api`
2. Run `make api` to generate code
3. Implement logic in `internal/logic/{module}/`
4. Add unit tests
5. Build and verify

### Adding a New Domain

1. Create `sql/{domain}/schema.sql`
2. Create `admin/internal/domain/{domain}/entity.go`
3. Create `admin/internal/domain/{domain}/repository.go`
4. Create `admin/internal/application/{domain}/`
5. Add API endpoints in `admin/desc/{domain}.api`
6. Update `admin/desc/admin.api` imports

### Adding a New Error Code

1. Add error to appropriate section in `pkg/code/code.go`
2. Follow the module's code range
3. Document in error codes reference

---

## Troubleshooting

### Common Issues

**API returns 404:**
- Check if route is registered in `internal/handler/routes.go`
- Verify .api file syntax is correct

**Build fails:**
- Run `go mod tidy` to clean dependencies
- Check for missing imports

**Database connection fails:**
- Verify MySQL is running
- Check connection string in config
- Ensure database exists

**Token expired:**
- Check JWT expiration settings
- Verify clock sync between services

### Debugging Tips

```bash
# Enable debug logging
go run ./cmd/admin.go -v

# Check API routes
curl http://localhost:8888/api/v1/

# View compiled routes
curl http://localhost:8888/routes
```

---

## References

- [Architecture Documentation](../ARCHITECTURE.md)
- [API Reference](../cross-cutting/api/2026-03-27-api-reference.md)
- [Error Codes Reference](../reference/2026-03-22-error-codes.md)
- [Database Overview](../reference/2026-03-22-database-overview.md)
- [go-zero Documentation](https://go-zero.dev/)
- [GORM Documentation](https://gorm.io/docs/)

---

## Document History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-03-27 | Technical Team | Initial developer guide |
