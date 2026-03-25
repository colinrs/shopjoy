# Project Rules

> **Full documentation:** See [AGENTS.md](./AGENTS.md) for complete project guide.

## Project Directory Convention

> **详细规范:** See [`.claude/rules/document/README.md`](.claude/rules/document/README.md)

### 三合一对应原则

文档、SQL、代码按相同领域组织：

| 领域 | 文档目录 | SQL 目录 | 代码目录 |
|-----|---------|---------|---------|
| user | `docs/domains/user/` | `sql/user/` | `domain/{user,adminuser,role,tenant}/` |
| product | `docs/domains/product/` | `sql/product/` | `domain/{product,market}/` |
| order | `docs/domains/order/` | `sql/order/` | `domain/{order,cart}/` |
| promotion | `docs/domains/promotion/` | `sql/promotion/` | `domain/{promotion,coupon}/` |
| points | `docs/domains/points/` | `sql/points/` | 待创建 |
| shop | `docs/domains/shop/` | `sql/shop/` | `handler/shop/` |
| storefront | `docs/domains/storefront/` | `sql/storefront/` | `domain/storefront/` |
| fulfillment | `docs/domains/fulfillment/` | `sql/fulfillment/` | `domain/fulfillment/` |
| payment | `docs/domains/payment/` | `sql/payment/` | `domain/payment/` |
| review | `docs/domains/review/` | `sql/review/` | `domain/review/` |

### 文档命名规范

```
{YYYY-MM-DD}-{领域}-{文档类型}.md
```

| 文档类型 | 命名格式 | 示例 |
|---------|---------|------|
| PRD | `{日期}-{领域}-prd.md` | `2026-03-24-order-prd.md` |
| Schema | `{日期}-{领域}-schema.md` | `2026-03-24-user-schema.md` |
| UI 设计 | `{日期}-{领域}-ui-design.md` | `2026-03-24-payment-ui-design.md` |
| 技术设计 | `{日期}-{名称}-design.md` | `2026-03-22-sku-design.md` |

### SQL 迁移命名规范

```
{YYYYMMDD}{序号}_{动作}_{对象}.sql
```

示例：`2026032401_create_reviews.sql`, `2026032201_alter_promotions_add_scope.sql`

## API Definition Changes

When modifying HTTP API definitions:

1. **Edit `.api` files** in the corresponding project's `desc/` directory:
   - `admin/desc/*.api` for admin service
   - `shop/desc/*.api` for shop service

2. **Run `make api`** to regenerate code:
   ```bash
   cd admin && make api
   # or
   cd shop && make api
   ```

3. **DO NOT edit** auto-generated files:
   - `internal/types/types.go`
   - `internal/handler/routes.go`
   - `internal/handler/*.go`

## Middleware Configuration

**Define middleware in `.api` files, NOT by editing `routes.go`:**

```go
@server (
    group:      users
    middleware: AuthMiddleware  // Add middleware here
)
service admin-api {
    @handler ListUsersHandler
    get /api/v1/users (ListUsersRequest) returns (ListUsersResponse)
}
```

Then implement in `internal/middleware/auth_middleware.go`.

## Build Commands

**ALWAYS use `make build` for compilation:**

```bash
# Build specific service
cd admin && make build
cd shop && make build

# Build all services from root
make build
```

**DO NOT use `go build` directly.**

## After Code Changes

**MUST run `make build`** after any code modifications to verify compilation succeeds.

## Error Handling

**Use custom errors from `pkg/code`, NOT `errors.New()`:**

```go
// BAD: Using errors.New()
var ErrOrderNotFound = errors.New("order not found")
return errors.New("something went wrong")

// GOOD: Use pkg/code errors
return code.ErrOrderNotFound
```

1. **Define all business errors** in `pkg/code/code.go` with proper HTTP status codes and error codes
2. **Error code ranges** by module (see `pkg/code/code.go` for full list):
   - Admin User: 10xxx
   - User: 11xxx
   - Product: 30xxx
   - Order: 40xxx
   - Payment: 50xxx
   - Cart: 60xxx
   - Coupon: 70xxx
   - Promotion: 80xxx
   - Tenant: 90xxx
   - Role: 100xxx
   - Shop: 110xxx
   - Fulfillment: 120xxx
3. **DO NOT create local error variables** with `errors.New()` in application or domain layers

## Requirement Development Workflow

> **详细规范:** See [`.claude/skills/requirement/SKILL.md`](.claude/skills/requirement/SKILL.md)

### 流程阶段

| Phase | 名称 | 执行者 | 输出物 |
|-------|------|--------|--------|
| 1 | 需求分析 | `product-manager` + `shopify-expert` | PRD 文档 |
| 2 | 设计阶段 | `backend-developer` + `api-designer` + `ui-designer` + `frontend-developer` | API设计 + UI设计 + 前端技术设计 |
| 3 | 计划阶段 | `writing-plans` | 开发计划 |
| 4 | 实施阶段 | `subagent-driven-development` | 代码实现 |

### 关键原则

| 类型 | 原则 |
|-----|------|
| MUST | 全栈覆盖：每个需求必须同时考虑 Frontend + Backend + UI + Database |
| MUST | 文档先行：所有设计文档完成并审批后，方可进入开发 |
| MUST | 一致性保证：最终实现与 PRD、UI 设计完全一致，不得遗漏功能 |
| MUST | 两轮审查：每个阶段输出物需经过至少 2 轮评审 |
| MUST | 用户确认：每个阶段完成必须获得用户明确同意 |
| MUST NOT | 文档未审批即开始开发 |
| MUST NOT | 只考虑后端忽略前端/UI |

### 触发条件

当用户提出新功能需求时，使用 `/requirement` skill 启动完整开发流程。