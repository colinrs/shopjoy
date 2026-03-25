# ShopJoy Developer Guide

> Getting started with ShopJoy e-commerce platform development

## Prerequisites

### Required

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.21+ | Backend development |
| Node.js | 18+ | Frontend development |
| MySQL | 8.0 | Database |
| Redis | 7.0 | Cache |

### Recommended

| Tool | Purpose |
|------|---------|
| Docker | Containerized development |
| Make | Build automation |
| golangci-lint | Go linting |

---

## Quick Start

### 1. Clone Repository

```bash
git clone https://github.com/colinrs/shopjoy.git
cd shopjoy
```

### 2. Database Setup

```bash
# Create database
mysql -u root -p -e "CREATE DATABASE shopjoy CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# Run migrations
mysql -u root -p shopjoy < sql/product.sql
mysql -u root -p shopjoy < sql/market.sql
mysql -u root -p shopjoy < sql/migrations/20260321_category_brand_inventory.sql
mysql -u root -p shopjoy < sql/migrations/20260322_product_localizations.sql
```

### 3. Configuration

```bash
# Copy example config
cp admin/etc/admin-api.yaml.example admin/etc/admin-api.yaml

# Edit with your settings
vim admin/etc/admin-api.yaml
```

Required environment variables:
```bash
export DB_PASSWORD=your_db_password
export REDIS_PASSWORD=your_redis_password
export JWT_SECRET=your_jwt_secret
```

### 4. Start Backend

```bash
# Install dependencies
go mod download

# Run admin API
cd admin && go run admin.go -f etc/admin-api.yaml
```

### 5. Start Frontend

```bash
# Admin dashboard
cd shop-admin
npm install
npm run dev

# Shop frontend (in another terminal)
cd joy
npm install
npm run dev
```

---

## Project Structure

```
shopjoy/
├── admin/                    # Admin backend service
│   ├── desc/                 # API definitions (*.api)
│   ├── etc/                  # Configuration
│   └── internal/
│       ├── handler/          # HTTP handlers (auto-generated)
│       ├── logic/            # Business logic
│       ├── domain/           # Domain entities
│       └── infrastructure/   # Repository implementations
│
├── shop/                     # Shop frontend API service
├── pkg/                      # Shared packages
├── shop-admin/               # Admin frontend (Vue 3)
├── joy/                      # Shop frontend (Vue 3)
├── sql/                      # Database schemas
└── docs/                     # Documentation
```

---

## Development Workflow

### API Development

#### 1. Define API in `.api` file

```go
// admin/desc/product.api
syntax = "v1"

type (
    CreateProductReq {
        Name  string `json:"name"`
        Price int64  `json:"price"`
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
    @doc "Create product"
    @handler CreateProductHandler
    post /api/v1/products (CreateProductReq) returns (CreateProductResp)
}
```

#### 2. Generate code

```bash
cd admin && make api
```

This generates:
- `internal/handler/products/create_product_handler.go`
- `internal/logic/products/create_product_logic.go`
- `internal/types/types.go`

#### 3. Implement business logic

```go
// admin/internal/logic/products/create_product_logic.go
func (l *CreateProductLogic) CreateProduct(req *types.CreateProductReq) (*types.CreateProductResp, error) {
    // 1. Get tenant context
    tenantID := tenant.FromContext(l.ctx)

    // 2. Create domain entity
    product := &product.Product{
        TenantID: tenantID,
        Name:     req.Name,
        Price:    shared.Money{Amount: req.Price, Currency: "CNY"},
    }

    // 3. Persist
    if err := l.svcCtx.ProductRepo.Create(l.ctx, product); err != nil {
        return nil, err
    }

    return &types.CreateProductResp{ID: product.ID}, nil
}
```

#### 4. Build and test

```bash
cd admin && make build
```

### Domain Development

#### 1. Define entity

```go
// admin/internal/domain/product/entity.go
type Product struct {
    ID          int64
    TenantID    shared.TenantID
    Name        string
    Price       shared.Money
    Status      ProductStatus
    AuditInfo   shared.AuditInfo
}

// Business method
func (p *Product) PutOnSale() error {
    if p.Status != StatusDraft && p.Status != StatusOffSale {
        return errors.New("invalid status transition")
    }
    p.Status = StatusOnSale
    return nil
}
```

#### 2. Define repository interface

```go
// admin/internal/domain/product/repository.go
type ProductRepository interface {
    Create(ctx context.Context, product *Product) error
    Update(ctx context.Context, product *Product) error
    FindByID(ctx context.Context, tenantID, id int64) (*Product, error)
}
```

#### 3. Implement repository

```go
// admin/internal/infrastructure/persistence/product_repository.go
type productRepository struct{}

func NewProductRepository() ProductRepository {
    return &productRepository{}
}

func (r *productRepository) FindByID(ctx context.Context, tenantID, id int64) (*Product, error) {
    // GORM implementation
}
```

### Frontend Development

#### 1. Create API client

```typescript
// shop-admin/src/api/product.ts
import request from '@/utils/request'

export function createProduct(data: CreateProductReq) {
  return request.post('/api/v1/products', data)
}

export function getProducts(params: ListQuery) {
  return request.get('/api/v1/products', { params })
}
```

#### 2. Create page component

```vue
<!-- shop-admin/src/views/products/index.vue -->
<template>
  <div class="products-page">
    <el-table :data="products">
      <el-table-column prop="name" label="Name" />
      <el-table-column prop="price" label="Price" />
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getProducts } from '@/api/product'

const products = ref([])

onMounted(async () => {
  const { data } = await getProducts({})
  products.value = data.list
})
</script>
```

---

## Build Commands

### Backend

```bash
# Generate API code
make api

# Build binary
make build

# Run tests
make test

# Lint
make lint
```

### Frontend

```bash
# Development server
npm run dev

# Build for production
npm run build

# Lint
npm run lint
```

---

## Testing

### Unit Tests

```go
// admin/internal/domain/product/entity_test.go
func TestProduct_PutOnSale(t *testing.T) {
    product := &Product{
        Status: StatusDraft,
    }

    err := product.PutOnSale()
    assert.NoError(t, err)
    assert.Equal(t, StatusOnSale, product.Status)
}
```

Run tests:
```bash
go test ./...
```

### Integration Tests

```go
// admin/internal/infrastructure/persistence/product_repository_test.go
func TestProductRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    repo := NewProductRepository(db)
    product := &Product{
        Name:   "Test Product",
        Price:  Money{Amount: 9900},
    }

    err := repo.Create(context.Background(), product)
    assert.NoError(t, err)
    assert.NotZero(t, product.ID)
}
```

---

## Code Style

### Go

Follow [Effective Go](https://golang.org/doc/effective_go) and [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md).

Key conventions:
- Use `gofmt` for formatting
- Error messages should not be capitalized
- Use meaningful variable names
- Document exported functions

### TypeScript/Vue

Follow [Vue Style Guide](https://vuejs.org/style-guide/).

Key conventions:
- Use Composition API with `<script setup>`
- Use TypeScript for type safety
- Follow naming conventions: PascalCase for components

---

## Configuration Reference

### admin-api.yaml

```yaml
Name: admin-api
Host: 0.0.0.0
Port: 8888

Database:
  Host: localhost
  Port: 3306
  User: root
  Password: ${DB_PASSWORD}
  DBName: shopjoy

Redis:
  Host: localhost
  Port: 6379
  Password: ${REDIS_PASSWORD}

Auth:
  AccessSecret: ${JWT_SECRET}
  AccessExpire: 86400

Log:
  Mode: console
  Level: info
```

---

## Troubleshooting

### Database Connection Error

```
Error: dial tcp 127.0.0.1:3306: connect: connection refused
```

Solution: Ensure MySQL is running and credentials are correct.

### Redis Connection Error

```
Error: dial tcp 127.0.0.1:6379: connect: connection refused
```

Solution: Ensure Redis is running.

### JWT Token Invalid

```
Error: token signature is invalid
```

Solution: Check JWT_SECRET is consistent across services.

### API Code Not Generated

```
Error: undefined: types.CreateProductReq
```

Solution: Run `make api` to generate types from `.api` files.

---

## Useful Commands

```bash
# Check Go version
go version

# Format code
go fmt ./...

# Run linter
golangci-lint run

# Check dependencies
go mod tidy

# View database
mysql -u root -p shopjoy -e "SELECT * FROM products LIMIT 5;"

# Clear Redis cache
redis-cli FLUSHDB
```

---

## Contributing

1. Create a feature branch
   ```bash
   git checkout -b feature/your-feature
   ```

2. Make changes and commit
   ```bash
   git add .
   git commit -m "feat: add your feature"
   ```

3. Push and create PR
   ```bash
   git push origin feature/your-feature
   ```

4. Wait for code review

---

## Resources

- [go-zero Documentation](https://go-zero.dev/)
- [Vue 3 Documentation](https://vuejs.org/)
- [Element Plus](https://element-plus.org/)
- [Tailwind CSS](https://tailwindcss.com/)
- [GORM Documentation](https://gorm.io/)