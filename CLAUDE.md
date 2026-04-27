# Project Rules

> **Full documentation:** See [AGENTS.md](./AGENTS.md) for complete project guide.

## Implementation Planning Protocol

1. **NEVER start autonomous codebase exploration** without explicit user request
2. When asked for implementation plan: FIRST read provided PRD/document → ask clarifying questions if needed → THEN create plan

## Session Continuity

When user says "continue" or references a previous task:

1. Check `.claude/sessions/` directory for recent context
2. Ask user: "Which previous task?" with top 3 recent session summaries
3. **NEVER assume context** from previous sessions — require explicit confirmation
4. If no context available: "I don't have context from previous sessions. Please restate your goal."

**Before ending session:** Summarize to `.claude/sessions/YYYY-MM-DD-task-name.md` with: goal, decisions made, current state, next step.

## Code Review Requirements

ALL implementations must pass automated code review before commit:

1. Run `/review` skill or agent before `git commit`
2. Check for: duplicate type definitions, missing exports, `any` types, unsafe URL handling, calculation errors, CSS typos
3. Verify error handling uses standardized `code` package, not `errors.New()`
4. **User confirmation required** after review fixes before final commit

## Process Phases

| Phase | Name | Executors | Deliverables |
|-------|------|-----------|--------------|
| 1 | Requirements Analysis | `product-manager` + `shopify-expert` | PRD Document |
| 2 | Design Phase | `backend-developer` + `api-designer` + `ui-designer` + `frontend-developer` | API Design + UI Design + Frontend Tech Design |
| 3 | Planning Phase | `writing-plans` | Development Plan |
| 4 | Implementation Phase | `subagent-driven-development` | Code Implementation |

## Key Principles

| Type | Principle |
|------|-----------|
| MUST | Full-stack coverage: Every requirement must consider Frontend + Backend + UI + Database |
| MUST | Documentation first: All design documents must be completed and approved before development |
| MUST | Consistency guarantee: Final implementation must match PRD and UI design exactly |
| MUST | Two-round review: Each phase deliverable requires at least 2 rounds of review |
| MUST | User confirmation: Each phase completion must receive explicit user approval |
| MUST NOT | Start development before documentation approval |
| MUST NOT | Only consider backend while ignoring frontend/UI |

## Project Directory Convention

> **Full specification:** See [`.claude/rules/document/README.md`](.claude/rules/document/README.md)

### Domain Mapping

Documents, SQL, and code are organized by the same domains:

| Domain | Docs Directory | SQL Directory | Code Directory |
|--------|---------------|---------------|----------------|
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

### Document Naming Convention

```
{YYYY-MM-DD}-{domain}-{type}.md
```

| Type | Format | Example |
|------|--------|---------|
| PRD | `{date}-{domain}-prd.md` | `2026-03-24-order-prd.md` |
| Schema | `{date}-{domain}-schema.md` | `2026-03-24-user-schema.md` |
| UI Design | `{date}-{domain}-ui-design.md` | `2026-03-24-payment-ui-design.md` |
| Tech Design | `{date}-{name}-design.md` | `2026-03-22-sku-design.md` |

**Special directories without timestamps** (always kept current):
- `docs/reference/` - `api-reference.md`, `database-overview.md`, `error-codes.md`
- `docs/guides/` - `developer-guide.md`, `onboarding.md`, `user-guide.md`
- `docs/cross-cutting/api/` - `api-reference.md`, `openapi.yaml`

### SQL Migration Naming Convention

```
{YYYYMMDD}{seq}_{action}_{object}.sql
```

Examples: `2026032401_create_reviews.sql`, `2026032201_alter_promotions_add_scope.sql`

## Go Backend Conventions

- **VAT/GST rates:** Use `string` type (backend returns strings), implement string↔number conversion in handlers
- **Timestamps:** 所有时间字段统一使用 `time.Time` 类型，数据库使用 `TIMESTAMP` 类型。详见 [`.claude/rules/golang/time.md`](.claude/rules/golang/time.md)
- **Money/Price:** API 层使用 `string` 类型表示元（如 `"1.99"` 表示 1.99 元），内部使用 `decimal.Decimal`。详见 [`.claude/rules/golang/price.md`](.claude/rules/golang/price.md)
- **Migrations:** Merge into schema files, check for duplicate table definitions
- **SQL Schema Consolidation:** 每个领域只保留一个 `schema.sql` 文件，包含该领域所有表的完整定义。migrations 文件中的字段/索引变更需合并到 schema.sql 后删除 migrations 目录
- **Errors:** Always use `code` package for standardized error codes, never `errors.New()`

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

3. **Auto-generated files (DO NOT edit manually):**
   - `internal/types/types.go` - request/response types
   - `internal/handler/routes.go` - route registration
   - `internal/handler/{module}/*.go` - HTTP handlers (auto-generated with stubs)
   - `internal/logic/{module}/*.go` - business logic stubs (auto-generated)

4. **Implementation workflow for new API endpoints:**
   - Add route definition in `.api` file
   - Run `make api` to generate handler and logic stubs
   - **Only implement the business logic** in the generated logic file (the stub contains `// todo: add your logic here`)
   - Handler files are already complete and should NOT be modified

5. **⚠️ MANDATORY: Frontend Synchronization (DO NOT SKIP)**
   When backend API definitions change, you MUST also update the corresponding frontend code:
   - **`src/api/{module}.ts`** - Update TypeScript type definitions to match backend
   - **`src/views/{module}/**/*.vue`** - Update all enum value comparisons in template logic (v-if, v-show, status checks)
   - **`src/components/**`** - Update any components that use the affected enum values

   **This step is MANDATORY and MUST NOT be skipped.** Any backend change that omits frontend updates is considered incomplete and will cause runtime bugs.

## Enum Conventions

**枚举定义以后端接口为准，前端必须严格跟随后端定义：**

| Type | Rule | Rationale |
|------|------|------------|
| MUST | Enum values must be defined in backend `.api` files with inline comments | Single source of truth, generated types reflect backend |
| MUST | Frontend TypeScript types must match backend enum values exactly | Avoid runtime mismatches, incorrect filtering/display |
| MUST | When backend enum changes, frontend must be updated accordingly | Maintain consistency across the full chain |
| MUST NOT | Frontend define its own enum values independent of backend | Causes inconsistencies, incorrect API calls |
| MUST NOT | Modify backend enum values to match frontend existing usage | Backend is authoritative source |

**Example of proper enum documentation in `.api` files:**
```go
type OrderStatus int // 0=pending_payment, 1=paid, 2=pending_shipment, 3=shipped, 4=completed, 5=cancelled
```

**Example of incorrect frontend usage:**
```typescript
// BAD: Frontend defining its own values
export type OrderStatus = 'pending' | 'paid' | 'shipped' | 'completed'

// GOOD: Frontend matching backend definition
export type OrderStatus = 'pending_payment' | 'paid' | 'pending_shipment' | 'shipped' | 'completed' | 'cancelled' | 'refunding' | 'refunded'
```

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

## Skill routing

When the user's request matches an available skill, invoke it via the Skill tool. The
skill has multi-step workflows, checklists, and quality gates that produce better
results than an ad-hoc answer. When in doubt, invoke the skill. A false positive is
cheaper than a false negative.

Key routing rules:
- Product ideas, "is this worth building", brainstorming → invoke /office-hours
- Strategy, scope, "think bigger", "what should we build" → invoke /plan-ceo-review
- Architecture, "does this design make sense" → invoke /plan-eng-review
- Design system, brand, "how should this look" → invoke /design-consultation
- Design review of a plan → invoke /plan-design-review
- Developer experience of a plan → invoke /plan-devex-review
- "Review everything", full review pipeline → invoke /autoplan
- Bugs, errors, "why is this broken", "wtf", "this doesn't work" → invoke /investigate
- Test the site, find bugs, "does this work" → invoke /qa (or /qa-only for report only)
- Code review, check the diff, "look at my changes" → invoke /review
- Visual polish, design audit, "this looks off" → invoke /design-review
- Developer experience audit, try onboarding → invoke /devex-review
- Ship, deploy, create a PR, "send it" → invoke /ship
- Merge + deploy + verify → invoke /land-and-deploy
- Configure deployment → invoke /setup-deploy
- Post-deploy monitoring → invoke /canary
- Update docs after shipping → invoke /document-release
- Weekly retro, "how'd we do" → invoke /retro
- Second opinion, codex review → invoke /codex
- Safety mode, careful mode, lock it down → invoke /careful or /guard
- Restrict edits to a directory → invoke /freeze or /unfreeze
- Upgrade gstack → invoke /gstack-upgrade
- Save progress, "save my work" → invoke /context-save
- Resume, restore, "where was I" → invoke /context-restore
- Security audit, OWASP, "is this secure" → invoke /cso
- Make a PDF, document, publication → invoke /make-pdf
- Launch real browser for QA → invoke /open-gstack-browser
- Import cookies for authenticated testing → invoke /setup-browser-cookies
- Performance regression, page speed, benchmarks → invoke /benchmark
- Review what gstack has learned → invoke /learn
- Tune question sensitivity → invoke /plan-tune
- Code quality dashboard → invoke /health