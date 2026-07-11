# SHOPJOY PROJECT KNOWLEDGE BASE

**Generated:** 2026-03-18
**Commit:** Latest
**Branch:** main
**Go Version:** 1.25.0  (per go.mod)
**Node Version:** 18+

## OVERVIEW

E-commerce platform with admin management and shop APIs. Built with go-zero microservices framework using DDD (Domain Driven Design) and repository pattern. Includes two modern Vue3 frontend applications.

## STRUCTURE

```
./
├── admin/              # Admin management API service (go-zero)
│   ├── desc/           # API definition files (.api)
│   └── internal/
│       ├── application/    # Application layer (use-case orchestration)
│       ├── domain/         # Domain layer (entities, value objects, repo interfaces)
│       ├── infrastructure/ # Infrastructure layer (repository implementations)
│       ├── handler/        # HTTP handlers (auto-generated + custom)
│       ├── logic/          # Business logic (scaffolded by goctl, safe to edit)
│       ├── middleware/     # Middleware implementations
│       ├── svc/            # Service context (dependency injection)
│       ├── config/         # Configuration structs
│       ├── types/          # Request/response types (auto-generated)
│       └── utils/
├── shop/               # Shop/e-commerce API service (go-zero)
│   ├── desc/           # API definition files (.api)
│   └── internal/
│       ├── application/    # Application layer (cart, order, payment services)
│       ├── domain/         # Domain layer (cart, order, payment entities)
│       ├── infrastructure/ # Infrastructure layer (persistence repositories)
│       ├── handler/        # HTTP handlers
│       ├── logic/          # Business logic
│       ├── svc/            # Service context
│       ├── config/         # Configuration structs
│       └── types/          # Request/response types (auto-generated)
├── pkg/                # Shared packages
│   ├── application/   # Shared application-layer helpers
│   ├── asyncq/        # Async queue (asynq, Redis-based)
│   ├── auth/          # JWT authentication
│   ├── cache/         # Cache abstraction (Redis, Ristretto)
│   ├── client/        # Client utilities (etcd)
│   ├── code/          # Error codes and definitions
│   ├── codec/         # Serialization (sonic JSON)
│   ├── contextx/      # Context utilities
│   ├── domain/        # Shared domain primitives (Money, TenantID, DomainEvent)
│   ├── gosafe/        # Safe goroutine utilities
│   ├── httpc/         # HTTP client utilities
│   ├── httpy/         # HTTP parsing utilities
│   ├── infra/         # Infrastructure (DB, Redis, metrics)
│   ├── rest/          # REST helpers
│   ├── response/      # HTTP response handlers
│   ├── sku/           # SKU code generation
│   ├── snowflake/     # ID generation
│   ├── tenant/        # Multi-tenant context management
│   └── utils/         # General utilities
├── shop-admin/         # Admin frontend (Vue3 + Element Plus)
│   ├── src/
│   │   ├── views/      # Page components
│   │   │   ├── login/
│   │   │   ├── dashboard/
│   │   │   ├── products/
│   │   │   ├── orders/
│   │   │   ├── users/
│   │   │   ├── promotions/
│   │   │   └── shop/
│   │   ├── layouts/    # Layout components
│   │   ├── components/ # Shared components
│   │   ├── stores/     # Pinia stores
│   │   ├── api/        # API client
│   │   └── router/     # Vue Router
│   └── package.json
├── joy/                # Shop frontend (Vue3 + Tailwind CSS)
│   ├── src/
│   │   ├── views/      # Page components
│   │   │   ├── home/
│   │   │   ├── login/
│   │   │   ├── products/
│   │   │   ├── cart/
│   │   │   ├── checkout/
│   │   │   ├── orders/
│   │   │   └── user/
│   │   ├── components/ # Shared components
│   │   ├── stores/     # Pinia stores
│   │   └── router/     # Vue Router
│   └── package.json
├── go.mod              # Go module definition
└── Makefile            # Build automation
```

## COMMANDS

### Root Makefile (build all services)
```bash
# Generate API code for all services (runs `make api` in each service)
make api

# Build all services
make build

# Run linter
golangci-lint run --timeout=10m
```

**After any code change, you MUST run `make build` to verify compilation succeeds.** Do not use `go build` directly — the Makefile is the canonical build entry point.

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

### Frontend Development

#### shop-admin (Admin Dashboard)
```bash
cd shop-admin

# Install dependencies
npm install

# Start development server (port 3000)
npm run dev

# Build for production
npm run build

# Run linter
npm run lint

# Format code
npm run format
```

**Key Dependencies:**
- `vue@^3.4.0` - Vue framework
- `element-plus@^2.5.0` - UI component library
- `@element-plus/icons-vue@^2.3.0` - Icons
- `echarts@^6.0.0` - Charts and data visualization
- `vue-router@^4.2.0` - Routing
- `pinia@^2.1.0` - State management

#### joy (Shop Frontend)
```bash
cd joy

# Install dependencies
npm install

# Start development server (port 3001)
npm run dev

# Build for production
npm run build
```

**Key Dependencies:**
- `vue@^3.4.0` - Vue framework
- `vue-router@^4.2.0` - Routing
- `pinia@^2.1.0` - State management
- `@heroicons/vue` - Heroicons icon library
- `tailwindcss` - CSS framework (via CDN)

### Database Setup

```bash
# Start MySQL + Redis via Docker
docker-compose up -d mysql redis

# Bootstrap schemas + seed data (sql/ is organized by business domain)
bash sql/init.sh
# or run the combined init script directly:
# mysql -u root -p < sql/init.sql
```

**Note:** README references `scripts/migrate.go`, but that file does not currently exist in `scripts/` (only `gen_bcrypt.go` is present). Use `sql/init.sh` / `sql/init.sql` for DB setup.

### Testing
```bash
# Backend tests
go test ./...

# Run specific test
go test ./pkg/cache/... -v
go test -run TestRistrettoCache ./pkg/cache/... -v

# Run with race detection
go test -race ./...
```

### Linting
```bash
# Backend
# Install golangci-lint (if needed)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run --timeout=10m

# Run on specific package
golangci-lint run ./pkg/...

# Frontend (shop-admin)
cd shop-admin && npm run lint
```

### go-zero Code Generation

**IMPORTANT: After updating `.api` definition files, ALWAYS use `make api` to regenerate code. Do NOT use `goctl` directly.**

```bash
# Correct way - use Makefile (runs format + gen)
cd shop && make api
cd admin && make api

# Or from root to regenerate all
make api
```

**What `make api` does:**
1. Formats `.api` definition files (`goctl api format`)
2. Generates Go code from API definitions (`goctl api go`)
3. Auto-generates Swagger documentation

**Auto-generated files (DO NOT edit manually):**
- `internal/types/types.go` - Request/response types
- `internal/handler/*.go` - HTTP handlers
- `internal/handler/routes.go` - Route registration
- `swagger/*.json` - Swagger API documentation
- `internal/middleware/*_middleware.go` - Middleware scaffolds

**Why:** The Makefile ensures consistent code generation with proper flags and style settings.

### MANDATORY: Frontend Synchronization After API Changes

When backend `.api` definitions change, you **MUST** also update the corresponding frontend code. This step is mandatory and must not be skipped — any backend change that omits frontend updates is incomplete and will cause runtime bugs.

- **`src/api/{module}.ts`** — Update TypeScript type definitions to match backend request/response types
- **`src/views/{module}/**/*.vue`** — Update all enum value comparisons in template logic (`v-if`, `v-show`, status checks)
- **`src/components/**`** — Update any components that use the affected enum values

### Enum Conventions (Backend is source of truth)

| Type | Rule |
|------|------|
| MUST | Enum values must be defined in backend `.api` files with inline comments |
| MUST | Frontend TypeScript types must match backend enum values exactly |
| MUST | When backend enum changes, frontend must be updated accordingly |
| MUST NOT | Frontend define its own enum values independent of backend |
| MUST NOT | Modify backend enum values to match frontend existing usage |

Example of proper enum documentation in `.api` files:
```go
type OrderStatus int // 0=pending_payment, 1=paid, 2=pending_shipment, 3=shipped, 4=completed, 5=cancelled
```

### VAT/GST Rates

VAT/GST rates are returned as **strings** from the backend. Frontend handlers must implement string↔number conversion when displaying or computing with rate values.

### Middleware Configuration in API Files

**IMPORTANT: Define middleware in `.api` files, NOT by editing `routes.go`.**

go-zero supports middleware declaration directly in API definition files using the `middleware` keyword in `@server` block:

```go
// Example: admin/desc/admin_user.api
@server (
    group:      admin_users
    middleware: AuthMiddleware  // Middleware name (must match scaffolded file)
)
service admin-api {
    @doc "获取管理员列表"
    @handler ListAdminUsersHandler
    get /api/v1/admin-users (ListAdminUsersRequest) returns (ListAdminUsersResponse)
}
```

**Workflow for adding middleware:**
1. Add `middleware: YourMiddleware` in `@server` block of `.api` file
2. Run `make api` to generate scaffold
3. Implement logic in `internal/middleware/your_middleware.go`
4. Add middleware field to `ServiceContext` in `internal/svc/service_context.go`

**Example - Auth Middleware:**
```go
// internal/middleware/auth_middleware.go
func NewAuthMiddleware(jwtSecret string) rest.Middleware {
    secret := []byte(jwtSecret)
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            // Validate JWT, set user context, etc.
            next(w, r)
        }
    }
}

// internal/svc/service_context.go
type ServiceContext struct {
    AuthMiddleware rest.Middleware
    // ...
}

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        AuthMiddleware: middleware.NewAuthMiddleware(c.JWT.Secret),
        // ...
    }
}
```

**Key Points:**
- Middleware must return `rest.Middleware` (function type)
- Use `rest.WithMiddlewares()` in generated routes.go (auto-generated)
- Routes without `middleware:` in `.api` will have no middleware applied

```bash
# Install goctl (required for code generation)
go install github.com/zeromicro/go-zero/tools/goctl@latest

# Manual goctl commands (NOT recommended - use Makefile instead)
# goctl api go --api ./desc/shop.api --dir ./ --style=go_zero
# goctl api format --dir ./desc

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

### Frontend Code Organization

#### shop-admin (Element Plus)
- `src/views/` - Page components (each folder = one page)
- `src/layouts/` - Layout components (MainLayout)
- `src/components/` - Shared reusable components
- `src/stores/` - Pinia stores (user, app state)
- `src/api/` - API client functions
- `src/router/` - Route definitions

#### joy (Tailwind CSS)
- `src/views/` - Page components
- `src/components/` - Shared components
- `src/stores/` - Pinia stores
- `src/router/` - Route definitions

### Naming
- **Go Files**: `snake_case.go` for handlers, `camelCase.go` for logic
- **Vue Files**: `PascalCase.vue` for components, `index.vue` for pages
- **Go Types**: PascalCase (e.g., `ShopLogic`, `ServiceContext`)
- **Vue Components**: PascalCase (e.g., `ProductCard.vue`)
- **Interfaces**: Verb-like with `er` suffix where appropriate
- **Private**: lowercase or underscore prefix

### Error Handling (Go)
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

**Standard HTTP Error Codes:**

For specific scenarios, use the predefined errors in `pkg/code/code.go` with correct HTTP status codes:

| Scenario | Error | HTTP Code | Business Code |
|----------|-------|-----------|---------------|
| Token expired | `code.ErrTokenExpired` | 401 | 40101 |
| Invalid token | `code.ErrTokenInvalid` | 401 | 40102 |
| Unauthorized | `code.ErrUnauthorized` | 401 | 40100 |
| Forbidden | `code.ErrForbidden` | 403 | 40300 |
| Not found | `code.ErrNotFound` | 404 | 40400 |
| Rate limited | `code.ErrTooManyRequests` | 429 | 42900 |
| Internal error | `code.ErrInternalServer` | 500 | 50000 |
| Service unavailable | `code.ErrServiceUnavailable` | 503 | 50300 |

```go
// Example: Return token expired error with correct HTTP status
func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token, err := jwt.Parse(tokenString, ...)
        if err != nil {
            // Returns 401 status code automatically
            return nil, code.ErrTokenExpired
        }
        next(w, r)
    }
}

// Example: Return rate limit error
func RateLimitMiddleware() rest.Middleware {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            if isRateLimited(r) {
                // Returns 429 status code automatically
                return nil, code.ErrTooManyRequests
            }
            next(w, r)
        }
    }
}
```

**IMPORTANT:** Always use `*code.Err` type for errors that need specific HTTP status codes. The `pkg/response/response.go` handler will automatically use the correct HTTP status from `Err.HTTPCode`.

**Business Error Code Ranges by Module** (see `pkg/code/code.go` for the full list):

| Module | Code Range |
|--------|------------|
| Admin User | 10xxx |
| User | 11xxx |
| Product | 30xxx |
| Order | 40xxx |
| Payment | 50xxx |
| Cart | 60xxx |
| Coupon | 70xxx |
| Promotion | 80xxx |
| Tenant | 90xxx |
| Role | 100xxx |
| Shop | 110xxx |
| Fulfillment | 120xxx |

**DO NOT** create local error variables with `errors.New()` in application or domain layers — define all business errors centrally in `pkg/code/code.go`.

### SQL Conventions

**Schema consolidation:** Each domain keeps a single `schema.sql` containing the complete definition of all tables in that domain. When a migration alters fields/indexes, merge the change into `schema.sql` and delete the migrations directory.

**Migration file naming:** `{YYYYMMDD}{seq}_{action}_{object}.sql`
- Examples: `2026040801_create_carriers.sql`, `2026040802_create_refund_reasons.sql`

### Domain Mapping

Docs, SQL, and code are organized by the same domains:

| Domain | Docs | SQL | Code |
|--------|------|-----|------|
| user | `docs/domains/user/` | `sql/user/` | `domain/{user,adminuser,role,tenant}/` |
| product | `docs/domains/product/` | `sql/product/` | `domain/{product,market}/` |
| order | `docs/domains/order/` | `sql/order/` | `domain/{order,cart}/` |
| promotion | `docs/domains/promotion/` | `sql/promotion/` | `domain/{promotion,coupon}/` |
| points | `docs/domains/points/` | `sql/points/` | TBD |
| shop | `docs/domains/shop/` | `sql/shop/` | `handler/shop/` |
| storefront | `docs/domains/storefront/` | `sql/storefront/` | `domain/storefront/` |
| fulfillment | `docs/domains/fulfillment/` | `sql/fulfillment/` | `domain/fulfillment/` |
| payment | `docs/domains/payment/` | `sql/payment/` | `domain/payment/` |
| review | `docs/domains/review/` | `sql/review/` | `domain/review/` |

### Error Handling (Vue)
- Use try-catch for async operations
- Display user-friendly error messages
- Log errors to console in development

```typescript
// Good
try {
  await api.login(form)
  router.push('/')
} catch (error: any) {
  ElMessage.error(error.message || '登录失败')
}
```

### Imports Order (Go)
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

### Imports Order (Vue)
```typescript
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
// Third-party icons
import { ShoppingCartIcon } from '@heroicons/vue/24/outline'
// Project internal
import { useUserStore } from '@/stores/user'
import { login } from '@/api/user'
```

### Context Usage (Go)
- Always pass `context.Context` as first parameter
- Use `logx.WithContext(ctx)` for logging
- Extract trace info from context for external calls

### Vue Best Practices
- Use Composition API with `<script setup>`
- Use `ref` for reactive state, `computed` for derived state
- Use Pinia for global state management
- Use Vue Router for navigation
- Use component-based architecture

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
- **ALWAYS** use `github.com/shopspring/decimal` for all monetary calculations internally
- **NEVER** use `float64` or `int` for money - precision loss is unacceptable
- Store in database as `DECIMAL(19,4)`
- **API layer convention:** expose money as `string` representing 元 (e.g. `"1.99"` for 1.99 元). Convert string↔`decimal.Decimal` at the handler boundary.

```go
import "github.com/shopspring/decimal"

// Good - decimal for internal calc
price := decimal.NewFromFloat(99.99)
total := price.Mul(decimal.NewFromInt(quantity))

// API response - string of 元
type Response {
    Total string `json:"total"`  // "99.99"
}

// Bad - precision loss
price := 99.99  // float64
total := price * float64(quantity)
```

### Time Handling
- **ALWAYS** use `time.Time` in Go and `TIMESTAMP` in the database
- **ALWAYS** store and return UTC time in backend
- **NEVER** store local time in database
- Frontend is responsible for timezone conversion and display
- Use `time.Now().UTC()` for current time

```go
// Good - time.Time in Go, UTC
now := time.Now().UTC()
createdAt := now  // store as time.Time; DB column is TIMESTAMP

// API response - UTC ISO format
type Response {
    CreatedAt time.Time `json:"created_at"`  // serialized as UTC ISO string
}

// Bad - storing local time, or pre-formatting to string in the struct
now := time.Now()  // Depends on server timezone
createdAt := now.Format("2006-01-02 15:04:05")  // loses type info
```

### Frontend Styling Guidelines

#### shop-admin (Element Plus)
- Use Element Plus components for consistency
- Use CSS variables for theming
- Follow BEM naming convention for custom CSS
- Use `scoped` styles in Vue components

#### joy (Tailwind CSS)
- Use Tailwind utility classes
- Use responsive prefixes (sm:, md:, lg:)
- Use hover:, focus: states
- Use custom colors defined in design system
- Primary: `#059669`, Secondary: `#10B981`, CTA: `#F97316`

### Icon Usage
- **shop-admin**: Use `@element-plus/icons-vue`
- **joy**: Use `@heroicons/vue` (outline or solid variants)
- NEVER use emojis as UI icons
- Always use consistent icon sizing

```vue
<!-- shop-admin -->
<el-icon><ShoppingCart /></el-icon>

<!-- joy -->
<ShoppingCartIcon class="w-6 h-6" />
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
| **Frontend: shop-admin** | `shop-admin/src/views/` | Admin page components |
| **Frontend: joy** | `joy/src/views/` | Shop page components |
| **Frontend layouts** | `*/src/layouts/` | Layout components |
| **Frontend stores** | `*/src/stores/` | Pinia state management |
| **Frontend API** | `*/src/api/` | API client functions |

## ANTI-PATTERNS

### Go Backend
- **DO NOT** edit `internal/types/types.go` or `internal/handler/routes.go` - auto-generated
- **DO NOT** add middleware by editing `routes.go` - use `.api` files with `middleware:` directive instead
- **DO NOT** hardcode cache keys - use `cache_key.go` utilities
- **DO NOT** use cache as primary storage
- **NEVER** ignore context cancellation in long operations
- **AVOID** blocking calls without timeout in handlers
- **NEVER** modify goctl-scaffolded files that have "Safe to edit" comments outside logic blocks

### Frontend
- **DO NOT** mix Tailwind and Element Plus in the same project
- **DO NOT** use inline styles - use Tailwind classes or scoped CSS
- **DO NOT** hardcode colors - use theme variables
- **AVOID** deeply nested component hierarchies
- **NEVER** ignore TypeScript errors
- **AVOID** using `any` type - use proper types/interfaces

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

### Tenant Isolation (TenantScopePlugin)

租户隔离通过 `TenantScopePlugin` 自动处理，**不要手动传递 tenantID**：

```go
// CORRECT - Plugin 自动注入 WHERE tenant_id = ?
func (r *userRepo) FindByID(ctx context.Context, db *gorm.DB, id int64) (*User, error) {
    var model userModel
    err := db.WithContext(ctx).First(&model, id).Error
    return &model, err
}

// WRONG - 手动传递 tenantID（已废弃）
func (r *userRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*User, error) {
    // ...
}
```

**关键规则：**
- Repository 方法**不要**接收 `tenantID` 参数
- Logic/Application 层**不要**从 ctx 提取 tenantID 传给 Repository
- Plugin 自动从 ctx 读取 tenantID（middleware 已设置）
- Create 操作 Plugin 自动填充 `tenant_id` 列
- 新增表需在 `pkg/infra/tenant_tables.go` 注册
- 跳过租户过滤：`db.WithContext(infra.SkipTenantScope(ctx))`

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

### Backend
- **go-zero**: Microservices framework with API code generation
- **goctl**: CLI tool for generating boilerplate from .api files
- **GORM**: ORM for database operations
- **go-redis**: Redis client
- **Ristretto**: In-memory cache library
- **sonic**: High-performance JSON library (ByteDance)

### Frontend
- **Vue 3**: Progressive JavaScript framework
- **Vite**: Next generation frontend tooling
- **Element Plus**: Vue 3 based component library
- **Tailwind CSS**: Utility-first CSS framework
- **Pinia**: Vue Store, intuitive and type-safe
- **Vue Router**: Official router for Vue.js
- **Heroicons**: Beautiful hand-crafted SVG icons

## DESIGN SYSTEM

### Colors
- **Primary**: `#059669` (emerald-600) - Success, primary actions
- **Secondary**: `#10B981` (emerald-500) - Secondary actions
- **CTA**: `#F97316` (orange-500) - Call to action, urgency
- **Background**: `#ECFDF5` (emerald-50) - Light backgrounds
- **Text**: `#064E3B` (emerald-900) - Primary text

### Typography
- **shop-admin**: Fira Code / Fira Sans (dashboard, data)
- **joy**: Rubik / Nunito Sans (ecommerce, clean)

### Common Rules
- No emoji icons - use SVG icons only
- cursor-pointer on all clickable elements
- Hover states with smooth transitions (150-300ms)
- Light mode: text contrast 4.5:1 minimum
- Focus states visible for keyboard navigation
- Responsive: 375px, 768px, 1024px, 1440px

## CHILD AGENTS

- `pkg/cache/AGENTS.md` - Cache abstraction layer

---

*Generated for agentic coding assistance*

<!-- code-review-graph MCP tools -->
## MCP Tools: code-review-graph

**IMPORTANT: This project has a knowledge graph. ALWAYS use the
code-review-graph MCP tools BEFORE using Grep/Glob/Read to explore
the codebase.** The graph is faster, cheaper (fewer tokens), and gives
you structural context (callers, dependents, test coverage) that file
scanning cannot.

### When to use graph tools FIRST

- **Exploring code**: `semantic_search_nodes` or `query_graph` instead of Grep
- **Understanding impact**: `get_impact_radius` instead of manually tracing imports
- **Code review**: `detect_changes` + `get_review_context` instead of reading entire files
- **Finding relationships**: `query_graph` with callers_of/callees_of/imports_of/tests_for
- **Architecture questions**: `get_architecture_overview` + `list_communities`

Fall back to Grep/Glob/Read **only** when the graph doesn't cover what you need.

### Key Tools

| Tool | Use when |
| ------ | ---------- |
| `detect_changes` | Reviewing code changes — gives risk-scored analysis |
| `get_review_context` | Need source snippets for review — token-efficient |
| `get_impact_radius` | Understanding blast radius of a change |
| `get_affected_flows` | Finding which execution paths are impacted |
| `query_graph` | Tracing callers, callees, imports, tests, dependencies |
| `semantic_search_nodes` | Finding functions/classes by name or keyword |
| `get_architecture_overview` | Understanding high-level codebase structure |
| `refactor_tool` | Planning renames, finding dead code |

### Workflow

1. The graph auto-updates on file changes (via hooks).
2. Use `detect_changes` for code review.
3. Use `get_affected_flows` to understand impact.
4. Use `query_graph` pattern="tests_for" to check coverage.
