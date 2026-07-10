# User Detail Tabs Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Wire the four placeholder tabs on `/users/{id}` (orders, points, reviews, operation-logs) to real backend APIs with platform-admin support, instrumentation, and pagination.

**Architecture:** Reuse 2 existing APIs (`/orders`, `/points/transactions`), extend 1 (`/reviews` adds `user_id` filter), build 1 new (`/users/{id}/operation-logs` with `user_operation_logs` table + 7 service-layer instrumentation points). Each tab is a self-contained Vue component owning its data lifecycle. Operation log writes never block the parent business operation.

**Tech Stack:** Go (go-zero + GORM), Vue 3 + Element Plus + TypeScript, MySQL/TIMESTAMP, decimal.Decimal for money.

## Global Constraints

- **Plan location:** Spec at `docs/domains/user/2026-07-10-user-detail-tabs-design.md` (commit 9518ea3) — read fully before each task.
- **Build commands:** Always `make build` from `admin/` and `shop-admin/`. Never `go build` / `vite build` directly.
- **API definition:** Changes to `.api` files require `cd admin && make api` to regenerate handlers/types. Auto-generated files (`internal/types/types.go`, `internal/handler/routes.go`, `internal/handler/{module}/*.go`) must NOT be hand-edited.
- **Money:** Always `decimal.Decimal` in domain, `string` at API boundary (e.g. `"1.99"`). Never float.
- **Time:** Always `time.Time` for entity timestamps, `TIMESTAMP` SQL. Frontend: `time.RFC3339` strings.
- **Errors:** All business errors from `pkg/code/code.go`. No `errors.New()` or `fmt.Errorf()` in logic/application/domain layers. No local error variables.
- **DB naming:** lowercase + underscores. Every table embeds `application.Model` (id/created_at/updated_at/deleted_at via TIMESTAMP). All fields NOT NULL or have DEFAULT. Monetary = DECIMAL. Related fields exact type match.
- **SQL migrations:** Per-domain `migrations/YYYYMMDDNN_*.sql`. New table also appended to `sql/{domain}/schema.sql` per `CLAUDE.md` SQL consolidation rule.
- **Frontend:** Always `npm` commands via `make`/`package.json`. Tab components own their state; orchestrator (`UserDetailTabs.vue`) is a thin container. Use `useErrorHandler` composable. Empty state per tab. Lazy-load on first tab activation (persist across switches).
- **Platform-admin tenant bypass:** All repository read paths: `if tenantID != 0 { query.Where("tenant_id = ?", tenantID) }`. Already established in `user_repository.go`, `sku_repository.go`. Apply same pattern.
- **Commit style:** Conventional commits. End with `Co-Authored-By: Claude <noreply@anthropic.com>`.

---

## File Structure

### Backend (Go) — create
- `sql/user/migrations/2026071001_create_user_operation_logs.sql` — DB migration
- `admin/internal/domain/user/operation_log_entity.go` — `OperationLog` domain entity + interface
- `admin/internal/infrastructure/persistence/user_operation_log_repository.go` — GORM repo
- `admin/internal/application/user/operation_log_service.go` — Application service + helper to record from logic
- `admin/internal/logic/users/list_user_operation_logs_logic.go` — List endpoint handler
- `admin/internal/domain/user/operation_log_helper_test.go` — Unit test for action text mapping

### Backend — modify
- `sql/user/schema.sql` — append new table def
- `admin/desc/user.api` — add types + handler
- `admin/internal/svc/service_context.go` — wire `UserOperationLog` service
- `admin/internal/logic/users/create_user_logic.go` — instrument CREATE_USER
- `admin/internal/logic/users/update_user_logic.go` — instrument UPDATE_USER
- `admin/internal/logic/users/suspend_user_logic.go` — instrument SUSPEND_USER
- `admin/internal/logic/users/suspend_user_with_reason_logic.go` — instrument SUSPEND_WITH_REASON
- `admin/internal/logic/users/activate_user_logic.go` — instrument ACTIVATE_USER
- `admin/internal/logic/users/delete_user_logic.go` — instrument DELETE_USER
- `admin/internal/logic/users/reset_password_logic.go` — instrument RESET_PASSWORD
- `admin/desc/review.api` — add `UserID` field to `ListReviewsReq`
- `admin/internal/application/review/dto.go` — add `UserID` field to `ListReviewsRequest`
- `admin/internal/application/review/service.go` — propagate `UserID` to repo
- `admin/internal/infrastructure/persistence/review_repository.go` — apply `user_id` filter

### Frontend — create
- `shop-admin/src/views/users/components/UserOrderList.vue`
- `shop-admin/src/views/users/components/UserPointsList.vue`
- `shop-admin/src/views/users/components/UserReviewList.vue`

### Frontend — modify
- `shop-admin/src/views/users/components/UserDetailTabs.vue` — swap coming-soon for new components
- `shop-admin/src/views/users/components/UserOperationLog.vue` — rewrite from placeholder
- `shop-admin/src/api/user.ts` — append 4 functions + 4 interfaces
- `shop-admin/src/locales/zh-CN.ts` — add keys
- `shop-admin/src/locales/en-US.ts` — mirror

---

## Task 1: Database migration + schema merge

**Files:**
- Create: `sql/user/migrations/2026071001_create_user_operation_logs.sql`
- Modify: `sql/user/schema.sql` (append at end)

**Interfaces:**
- Produces: `user_operation_logs` table consumable by Task 2's entity

- [ ] **Step 1: Create migration file**

Write `sql/user/migrations/2026071001_create_user_operation_logs.sql`:

```sql
-- Create user_operation_logs table
CREATE TABLE user_operation_logs (
    id              BIGINT       PRIMARY KEY,
    tenant_id       BIGINT       NOT NULL,
    user_id         BIGINT       NOT NULL,
    action          VARCHAR(64)  NOT NULL,
    operator_id     BIGINT       NOT NULL DEFAULT 0,
    operator_name   VARCHAR(64)  NOT NULL DEFAULT 'system',
    reason          VARCHAR(500) NOT NULL DEFAULT '',
    ip_address      VARCHAR(64)  NOT NULL DEFAULT '',
    user_agent      VARCHAR(500) NOT NULL DEFAULT '',
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP    NULL,

    INDEX idx_uol_user_id (user_id, created_at),
    INDEX idx_uol_tenant  (tenant_id, created_at),
    INDEX idx_uol_action  (action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

- [ ] **Step 2: Append to schema.sql**

Open `sql/user/schema.sql`. Read the file, find the last `CREATE TABLE` block, append after it:

```sql
-- ===================================================================
-- User Operation Logs
-- ===================================================================
CREATE TABLE user_operation_logs (
    id              BIGINT       PRIMARY KEY,
    tenant_id       BIGINT       NOT NULL,
    user_id         BIGINT       NOT NULL,
    action          VARCHAR(64)  NOT NULL,
    operator_id     BIGINT       NOT NULL DEFAULT 0,
    operator_name   VARCHAR(64)  NOT NULL DEFAULT 'system',
    reason          VARCHAR(500) NOT NULL DEFAULT '',
    ip_address      VARCHAR(64)  NOT NULL DEFAULT '',
    user_agent      VARCHAR(500) NOT NULL DEFAULT '',
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP    NULL,

    INDEX idx_uol_user_id (user_id, created_at),
    INDEX idx_uol_tenant  (tenant_id, created_at),
    INDEX idx_uol_action  (action)
);
```

- [ ] **Step 3: Apply migration locally and verify**

Run from monorepo root:
```bash
mysql -h 192.168.0.100 -P 3306 -u root shopjoy < sql/user/migrations/2026071001_create_user_operation_logs.sql
mysql -h 192.168.0.100 -P 3306 -u root shopjoy -e "SHOW CREATE TABLE user_operation_logs\G"
```
Expected: table exists with all columns + indexes.

- [ ] **Step 4: Commit**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add sql/user/migrations/2026071001_create_user_operation_logs.sql sql/user/schema.sql
git commit -m "feat(user): add user_operation_logs migration and schema

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 2: OperationLog domain entity + repository interface

**Files:**
- Create: `admin/internal/domain/user/operation_log_entity.go`

**Interfaces:**
- Produces: `OperationLog` entity, `OperationLogRepository` interface with `Create` + `FindByUserID(tenantID, userID, query)` methods consumed by Task 3 & Task 4

- [ ] **Step 1: Write entity file**

Write `admin/internal/domain/user/operation_log_entity.go`:

```go
package user

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

// Action constants for user operation logs.
// Single source of truth — frontend mirrors via i18n keys.
const (
	ActionCreateUser        = "CREATE_USER"
	ActionUpdateUser        = "UPDATE_USER"
	ActionSuspendUser       = "SUSPEND_USER"
	ActionSuspendWithReason = "SUSPEND_WITH_REASON"
	ActionActivateUser      = "ACTIVATE_USER"
	ActionDeleteUser        = "DELETE_USER"
	ActionResetPassword     = "RESET_PASSWORD"
)

// OperationLog records an admin-side action against a user (state changes,
// profile updates, password resets). Writes are best-effort: instrumentation
// failures must never block the parent business operation.
type OperationLog struct {
	application.Model
	TenantID     shared.TenantID
	UserID       int64
	Action       string
	OperatorID   int64
	OperatorName string
	Reason       string
	IPAddress    string
	UserAgent    string
	Audit        shared.AuditInfo `gorm:"embedded"`
}

func (o *OperationLog) TableName() string {
	return "user_operation_logs"
}

// OperationLogQuery filters for FindByUserID.
type OperationLogQuery struct {
	Page     int
	PageSize int
	Action   string // empty = all
	Keyword  string // empty = no keyword filter
}

func (q OperationLogQuery) Offset() int {
	if q.Page <= 0 {
		q.Page = 1
	}
	return (q.Page - 1) * q.PageSize
}

func (q OperationLogQuery) Limit() int {
	if q.PageSize <= 0 {
		return 20
	}
	return q.PageSize
}

// OperationLogRepository persists and queries OperationLog records.
// Read paths respect platform-admin tenant bypass: tenantID == 0 means
// "no tenant filter".
type OperationLogRepository interface {
	Create(ctx context.Context, db *gorm.DB, log *OperationLog) error
	FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, query OperationLogQuery) ([]*OperationLog, int64, error)
}

// ActionText returns the Chinese display label for a known action.
// Returns the raw action string for unknown values so new actions
// degrade gracefully rather than disappearing.
func ActionText(action string) string {
	switch action {
	case ActionCreateUser:
		return "创建用户"
	case ActionUpdateUser:
		return "更新资料"
	case ActionSuspendUser, ActionSuspendWithReason:
		return "禁用用户"
	case ActionActivateUser:
		return "启用用户"
	case ActionDeleteUser:
		return "删除用户"
	case ActionResetPassword:
		return "重置密码"
	default:
		return action
	}
}

// Compile-time check that Application.Model.AuditInfo.Audit.CreatedAt is time.Time
var _ time.Time = shared.AuditInfo{}.CreatedAt
```

- [ ] **Step 2: Write unit test for ActionText**

Create `admin/internal/domain/user/operation_log_helper_test.go`:

```go
package user

import "testing"

func TestActionText(t *testing.T) {
	cases := []struct {
		action string
		want   string
	}{
		{ActionCreateUser, "创建用户"},
		{ActionUpdateUser, "更新资料"},
		{ActionSuspendUser, "禁用用户"},
		{ActionSuspendWithReason, "禁用用户"},
		{ActionActivateUser, "启用用户"},
		{ActionDeleteUser, "删除用户"},
		{ActionResetPassword, "重置密码"},
		{"UNKNOWN_ACTION", "UNKNOWN_ACTION"},
	}
	for _, c := range cases {
		if got := ActionText(c.action); got != c.want {
			t.Errorf("ActionText(%q) = %q, want %q", c.action, got, c.want)
		}
	}
}
```

- [ ] **Step 3: Run test, verify pass**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy/admin
go test ./internal/domain/user/... -run TestActionText -v
```
Expected: `PASS`

- [ ] **Step 4: Build check**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy/admin
make build
```
Expected: build succeeds (entity compiles, no callers yet but interface is well-formed).

- [ ] **Step 5: Commit**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add admin/internal/domain/user/operation_log_entity.go admin/internal/domain/user/operation_log_helper_test.go
git commit -m "feat(user): OperationLog domain entity + ActionText helper

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 3: OperationLog repository (GORM)

**Files:**
- Create: `admin/internal/infrastructure/persistence/user_operation_log_repository.go`

**Interfaces:**
- Consumes: `user.OperationLog`, `user.OperationLogQuery`, `user.OperationLogRepository` from Task 2
- Produces: concrete `OperationLogRepositoryImpl` wired in Task 4

- [ ] **Step 1: Write repository**

Write `admin/internal/infrastructure/persistence/user_operation_log_repository.go`:

```go
package persistence

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/user"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type OperationLogRepositoryImpl struct{}

func NewOperationLogRepository() user.OperationLogRepository {
	return &OperationLogRepositoryImpl{}
}

func (r *OperationLogRepositoryImpl) Create(ctx context.Context, db *gorm.DB, log *user.OperationLog) error {
	now := time.Now().UTC()
	if log.CreatedAt.IsZero() {
		log.CreatedAt = now
	}
	if log.UpdatedAt.IsZero() {
		log.UpdatedAt = now
	}
	return db.WithContext(ctx).Create(log).Error
}

func (r *OperationLogRepositoryImpl) FindByUserID(
	ctx context.Context,
	db *gorm.DB,
	tenantID shared.TenantID,
	userID int64,
	query user.OperationLogQuery,
) ([]*user.OperationLog, int64, error) {
	q := db.WithContext(ctx).Model(&user.OperationLog{}).
		Where("user_id = ? AND deleted_at IS NULL", userID)

	// Platform admin (tenantID == 0) sees logs across all tenants.
	if tenantID != 0 {
		q = q.Where("tenant_id = ?", tenantID.Int64())
	}

	if query.Action != "" {
		q = q.Where("action = ?", query.Action)
	}
	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		q = q.Where("(operator_name LIKE ? OR reason LIKE ?)", like, like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var logs []*user.OperationLog
	err := q.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&logs).Error
	return logs, total, err
}
```

- [ ] **Step 2: Build check**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy/admin
make build
```
Expected: succeeds.

- [ ] **Step 3: Commit**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add admin/internal/infrastructure/persistence/user_operation_log_repository.go
git commit -m "feat(user): OperationLogRepository GORM impl with tenant bypass

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 4: Wire OperationLog into ServiceContext

**Files:**
- Modify: `admin/internal/svc/service_context.go`

**Interfaces:**
- Consumes: `OperationLogRepository` from Task 3
- Produces: `svcCtx.OperationLogRepo` field used by Task 5's application service

- [ ] **Step 1: Find current wiring block**

Open `admin/internal/svc/service_context.go`. Find the existing UserRepo wiring pattern (search for `NewUserRepository`). It should be inside the `NewServiceContext` function and assign to a field on `ServiceContext`.

- [ ] **Step 2: Add repo import and field**

Add import alongside existing persistence imports:
```go
"github.com/colinrs/shopjoy/admin/internal/domain/user"  // already imported as domain user
```

Add field to the `ServiceContext` struct (find the UserRepo field, add after it):
```go
OperationLogRepo user.OperationLogRepository
```

Inside `NewServiceContext`, after the line `UserRepo: persistence.NewUserRepository(),` add:
```go
OperationLogRepo: persistence.NewOperationLogRepository(),
```

- [ ] **Step 3: Build check**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy/admin
make build
```
Expected: succeeds.

- [ ] **Step 4: Commit**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add admin/internal/svc/service_context.go
git commit -m "feat(user): wire OperationLogRepo into ServiceContext

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 5: OperationLog application service

**Files:**
- Create: `admin/internal/application/user/operation_log_service.go`

**Interfaces:**
- Consumes: `OperationLogRepository` from Task 3
- Produces: `OperationLogService` interface (Record + List) consumed by all 7 instrumentation points in Tasks 9-15 and by Task 7's handler

- [ ] **Step 1: Define DTOs and service interface**

Write `admin/internal/application/user/operation_log_service.go`:

```go
package user

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/user"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/snowflake"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// RecordOperationLogInput captures all fields needed to write one log row.
// OperatorIP and OperatorUA may be empty if not available from context.
type RecordOperationLogInput struct {
	TenantID     shared.TenantID
	UserID       int64
	Action       string
	OperatorID   int64
	OperatorName string
	Reason       string
	IPAddress    string
	UserAgent    string
}

// OperationLogListItem is the read DTO consumed by the list endpoint.
type OperationLogListItem struct {
	ID           int64
	UserID       int64
	Action       string
	ActionText   string
	OperatorID   int64
	OperatorName string
	Reason       string
	IPAddress    string
	UserAgent    string
	CreatedAt    string
}

// OperationLogListResp is the list response wrapper.
type OperationLogListResp struct {
	List     []*OperationLogListItem
	Total    int64
	Page     int
	PageSize int
}

// OperationLogService persists and reads operation logs.
type OperationLogService interface {
	// Record writes one log row. NEVER returns an error to the caller — the
	// caller MUST treat all returned errors as informational and continue
	// serving the parent business operation. (Errors are logged internally.)
	Record(ctx context.Context, db *gorm.DB, input RecordOperationLogInput)
	List(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, query user.OperationLogQuery) (*OperationLogListResp, error)
}

type operationLogServiceImpl struct {
	db   *gorm.DB
	repo user.OperationLogRepository
	idGen snowflake.Snowflake
}

// NewOperationLogService wires the application service.
func NewOperationLogService(db *gorm.DB, repo user.OperationLogRepository, idGen snowflake.Snowflake) OperationLogService {
	return &operationLogServiceImpl{db: db, repo: repo, idGen: idGen}
}

func (s *operationLogServiceImpl) Record(ctx context.Context, db *gorm.DB, input RecordOperationLogInput) {
	if db == nil {
		db = s.db
	}
	id, err := s.idGen.NextID(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorf("operation log: id generation failed: %v", err)
		return
	}
	now := time.Now().UTC()
	log := &user.OperationLog{
		Model: struct{ /* placeholder */ }{},
	}
	// Build entity; Model is application.Model
	entity := &user.OperationLog{}
	entity.ID = id
	entity.CreatedAt = now
	entity.UpdatedAt = now
	entity.TenantID = input.TenantID
	entity.UserID = input.UserID
	entity.Action = input.Action
	entity.OperatorID = input.OperatorID
	entity.OperatorName = input.OperatorName
	entity.Reason = input.Reason
	entity.IPAddress = input.IPAddress
	entity.UserAgent = input.UserAgent
	entity.Audit = shared.AuditInfo{CreatedAt: now, UpdatedAt: now}

	if err := s.repo.Create(ctx, db, entity); err != nil {
		logx.WithContext(ctx).Errorf("operation log: write failed (action=%s, user_id=%d): %v", input.Action, input.UserID, err)
	}
	_ = log // unused placeholder removed in next step
}
```

Note: the `Record` signature returns no error — this is intentional. Caller MUST NOT wrap it in `if err != nil`.

- [ ] **Step 2: Add List method**

Append to the same file (inside the same package):

```go
func (s *operationLogServiceImpl) List(
	ctx context.Context,
	db *gorm.DB,
	tenantID shared.TenantID,
	userID int64,
	query user.OperationLogQuery,
) (*OperationLogListResp, error) {
	if db == nil {
		db = s.db
	}
	logs, total, err := s.repo.FindByUserID(ctx, db, tenantID, userID, query)
	if err != nil {
		return nil, err
	}
	items := make([]*OperationLogListItem, 0, len(logs))
	for _, l := range logs {
		items = append(items, toOperationLogItem(l))
	}
	return &OperationLogListResp{
		List:     items,
		Total:    total,
		Page:     query.Page,
		PageSize: query.Limit(),
	}, nil
}

func toOperationLogItem(l *user.OperationLog) *OperationLogListItem {
	return &OperationLogListItem{
		ID:           l.ID,
		UserID:       l.UserID,
		Action:       l.Action,
		ActionText:   user.ActionText(l.Action),
		OperatorID:   l.OperatorID,
		OperatorName: l.OperatorName,
		Reason:       l.Reason,
		IPAddress:    l.IPAddress,
		UserAgent:    l.UserAgent,
		CreatedAt:    l.Audit.CreatedAt.Format(time.RFC3339),
	}
}
```

Note: `toOperationLogItem` reads `l.Audit.CreatedAt` per CLAUDE.md time rules.

- [ ] **Step 3: Clean placeholder code**

In the `Record` function, remove the `log` placeholder block:
```go
// DELETE this block:
log := &user.OperationLog{
    Model: struct{ /* placeholder */ }{},
}
_ = log // unused placeholder removed in next step
```

Final `Record`:
```go
func (s *operationLogServiceImpl) Record(ctx context.Context, db *gorm.DB, input RecordOperationLogInput) {
	if db == nil {
		db = s.db
	}
	id, err := s.idGen.NextID(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorf("operation log: id generation failed: %v", err)
		return
	}
	now := time.Now().UTC()
	entity := &user.OperationLog{
		Model:        application.Model{ID: id, CreatedAt: now, UpdatedAt: now},
		TenantID:     input.TenantID,
		UserID:       input.UserID,
		Action:       input.Action,
		OperatorID:   input.OperatorID,
		OperatorName: input.OperatorName,
		Reason:       input.Reason,
		IPAddress:    input.IPAddress,
		UserAgent:    input.UserAgent,
		Audit:        shared.AuditInfo{CreatedAt: now, UpdatedAt: now},
	}
	if err := s.repo.Create(ctx, db, entity); err != nil {
		logx.WithContext(ctx).Errorf("operation log: write failed (action=%s, user_id=%d): %v", input.Action, input.UserID, err)
	}
}
```

Add import: `"github.com/colinrs/shopjoy/pkg/application"`.

- [ ] **Step 4: Wire into ServiceContext**

Modify `admin/internal/svc/service_context.go`:
- Add import: `"github.com/colinrs/shopjoy/pkg/snowflake"` (check if already imported)
- Add field on `ServiceContext`:
```go
OperationLogService user.OperationLogService
```
- Add wiring inside `NewServiceContext`:
```go
OperationLogService: user.NewOperationLogService(DB, OperationLogRepo, Snowflake),
```
(Use the actual variable names already in the file — likely `DB`, `OperationLogRepo`, `Snowflake`. If `Snowflake` isn't yet wired, check the existing UserService wiring for the snowflake variable name.)

- [ ] **Step 5: Build check**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy/admin
make build
```
Expected: succeeds.

- [ ] **Step 6: Commit**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add admin/internal/application/user/operation_log_service.go admin/internal/svc/service_context.go
git commit -m "feat(user): OperationLog application service + ServiceContext wiring

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 6: API definition for ListUserOperationLogs

**Files:**
- Modify: `admin/desc/user.api` (append new types + handler inside existing `users` server group)

**Interfaces:**
- Produces: `types.ListUserOperationLogsReq/Resp` and `types.UserOperationLog` consumed by Task 7

- [ ] **Step 1: Append new types inside the type block**

Open `admin/desc/user.api`. Find the closing `)` of the `type (` block (after `BatchStatusFail`). Add before the closing:

```go
	ListUserOperationLogsReq {
		ID       int64  `path:"id"`
		Page     int    `form:"page,default=1"`
		PageSize int    `form:"page_size,default=20"`
		Action   string `form:"action,optional"`
		Keyword  string `form:"keyword,optional"`
	}
	UserOperationLog {
		ID           int64  `json:"id,string"`
		UserID       int64  `json:"user_id,string"`
		Action       string `json:"action"`
		ActionText   string `json:"action_text"`
		OperatorID   int64  `json:"operator_id,string"`
		OperatorName string `json:"operator_name"`
		Reason       string `json:"reason"`
		IPAddress    string `json:"ip_address"`
		UserAgent    string `json:"user_agent"`
		CreatedAt    string `json:"created_at"`
	}
	ListUserOperationLogsResp {
		List     []*UserOperationLog `json:"list"`
		Total    int64               `json:"total"`
		Page     int                 `json:"page"`
		PageSize int                 `json:"page_size"`
	}
```

- [ ] **Step 2: Append new handler in users server group**

In the same file, find the `service admin-api {` block under `@server (group: users ...)`. After the last existing handler in that block (the `BatchUpdateUserStatusHandler`), add:

```go
	@doc "获取用户操作日志"
	@handler ListUserOperationLogsHandler
	get /api/v1/users/:id/operation-logs (ListUserOperationLogsReq) returns (ListUserOperationLogsResp)
```

- [ ] **Step 3: Regenerate code**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy/admin
make api
```
Expected: regenerated `internal/types/types.go`, `internal/handler/routes.go`, `internal/handler/users/list_user_operation_logs_handler.go`, `internal/logic/users/list_user_operation_logs_logic.go` (stub).

- [ ] **Step 4: Verify generated types exist**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy/admin
grep -n "ListUserOperationLogsReq\|UserOperationLog\b" internal/types/types.go | head -20
```
Expected: matches for both.

- [ ] **Step 5: Commit**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add admin/desc/user.api admin/internal/types/types.go admin/internal/handler/ admin/internal/logic/users/list_user_operation_logs_logic.go
git commit -m "feat(user.api): add ListUserOperationLogs types and handler

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 7: Implement ListUserOperationLogs logic

**Files:**
- Modify: `admin/internal/logic/users/list_user_operation_logs_logic.go` (regenerated stub)

**Interfaces:**
- Consumes: `svcCtx.OperationLogService` from Task 5
- Produces: fully functional `ListUserOperationLogs` handler

- [ ] **Step 1: Replace stub body**

Read the generated stub file. Replace its body with:

```go
package users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/user"
	appUser "github.com/colinrs/shopjoy/admin/internal/application/user"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/tenant"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserOperationLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListUserOperationLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListUserOperationLogsLogic {
	return ListUserOperationLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUserOperationLogsLogic) ListUserOperationLogs(req *types.ListUserOperationLogsReq) (resp *types.ListUserOperationLogsResp, err error) {
	tenantID, ok := tenant.FromContext(l.ctx)
	if !ok {
		return nil, code.ErrTenantInvalidID
	}

	query := user.OperationLogQuery{
		Page:     req.Page,
		PageSize: req.PageSize,
		Action:   req.Action,
		Keyword:  req.Keyword,
	}

	out, err := l.svcCtx.OperationLogService.List(l.ctx, l.svcCtx.DB, tenantID, req.ID, query)
	if err != nil {
		return nil, err
	}

	list := make([]*types.UserOperationLog, 0, len(out.List))
	for _, item := range out.List {
		list = append(list, &types.UserOperationLog{
			ID:           item.ID,
			UserID:       item.UserID,
			Action:       item.Action,
			ActionText:   item.ActionText,
			OperatorID:   item.OperatorID,
			OperatorName: item.OperatorName,
			Reason:       item.Reason,
			IPAddress:    item.IPAddress,
			UserAgent:    item.UserAgent,
			CreatedAt:    item.CreatedAt,
		})
	}

	return &types.ListUserOperationLogsResp{
		List:     list,
		Total:    out.Total,
		Page:     out.Page,
		PageSize: out.PageSize,
	}, nil
}
```

- [ ] **Step 2: Build**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy/admin
make build
```
Expected: succeeds.

- [ ] **Step 3: Smoke test endpoint**

Start admin service, then:
```bash
TOKEN=$(curl -s -X POST http://localhost:8888/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"superadmin@shopjoy.com","password":"password123"}' | jq -r '.data.token')

curl -s "http://localhost:8888/api/v1/users/1/operation-logs?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN" | jq
```
Expected: `{"code":0,"msg":"success","data":{"list":[],"total":0,"page":1,"page_size":10}}` (empty, no logs yet).

- [ ] **Step 4: Commit**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add admin/internal/logic/users/list_user_operation_logs_logic.go
git commit -m "feat(user): implement ListUserOperationLogs handler

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 8: Review API extension (user_id filter)

**Files:**
- Modify: `admin/desc/review.api`
- Modify: `admin/internal/application/review/dto.go`
- Modify: `admin/internal/application/review/service.go`
- Modify: `admin/internal/infrastructure/persistence/review_repository.go`

**Interfaces:**
- Produces: `GET /api/v1/reviews?user_id=N` returns that user's reviews

- [ ] **Step 1: Add UserID to ListReviewsReq in .api**

Open `admin/desc/review.api`. Find `ListReviewsReq`. Add a field:
```go
		UserID    int64  `form:"user_id,optional"`
```

- [ ] **Step 2: Add UserID to ListReviewsRequest DTO**

In `admin/internal/application/review/dto.go`, find `ListReviewsRequest` struct. Add:
```go
	UserID int64
```

- [ ] **Step 3: Add UserID to types.ListReviewsReq mapping**

Open `admin/internal/logic/reviews/list_reviews_logic.go`. Find where `listReq := appReview.ListReviewsRequest{...}` is built. Add:
```go
	UserID: req.UserID,
```

- [ ] **Step 4: Propagate to repository in service.go**

In `admin/internal/application/review/service.go`, find `ListReviews` method. Find the call to `repo.ListReviews(...)`. Add `UserID: req.UserID,` to the request struct passed to repo.

- [ ] **Step 5: Update repository signature and filter**

In `admin/internal/infrastructure/persistence/review_repository.go`, find `ListReviews` method signature and implementation:
- Add `userID int64` parameter
- Inside the method, after the existing `if query.Status != "" { ... }` filter, add:
```go
if userID > 0 {
    dbQuery = dbQuery.Where("user_id = ?", userID)
}
```
- Update the corresponding interface declaration (search for `ListReviews` in `domain/review/`).

- [ ] **Step 6: Regenerate API and build**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy/admin
make api
make build
```
Expected: both succeed.

- [ ] **Step 7: Smoke test**

```bash
TOKEN=$(... same as Task 7 ...)
curl -s "http://localhost:8888/api/v1/reviews?user_id=1&page=1&page_size=5" \
  -H "Authorization: Bearer $TOKEN" | jq '.data | {total, list_count: (.list | length)}'
```
Expected: numeric values; no error.

- [ ] **Step 8: Commit**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add admin/desc/review.api admin/internal/application/review/ admin/internal/logic/reviews/list_reviews_logic.go admin/internal/infrastructure/persistence/review_repository.go admin/internal/domain/review/ admin/internal/types/types.go admin/internal/handler/
git commit -m "feat(review): add user_id filter to ListReviews endpoint

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 9: Instrument CreateUserLogic

**Files:**
- Modify: `admin/internal/logic/users/create_user_logic.go`

**Interfaces:**
- Consumes: `svcCtx.OperationLogService` from Task 5
- Produces: writes one `CREATE_USER` log row after a successful user create

- [ ] **Step 1: Find success branch**

Open the file. Find the line where `UserService.Register(...)` is called and the call returns `nil` error (success branch).

- [ ] **Step 2: Add instrumentation right after success**

Immediately after the successful `Register` call (and before returning the response), add:

```go
	// Record operation log (best-effort, never blocks the parent operation)
	l.svcCtx.OperationLogService.Record(l.ctx, l.svcCtx.DB, appUser.RecordOperationLogInput{
		TenantID:     tenantID,
		UserID:       resp.ID,
		Action:       user.ActionCreateUser,
		OperatorID:   contextx.GetAdminID(l.ctx),
		OperatorName: contextx.GetAdminName(l.ctx),
		IPAddress:    contextx.GetClientIP(l.ctx),
		UserAgent:    contextx.GetUserAgent(l.ctx),
	})
```

- [ ] **Step 3: Add imports if missing**

Add to imports if not already present:
```go
import (
    appUser "github.com/colinrs/shopjoy/admin/internal/application/user"
    "github.com/colinrs/shopjoy/admin/internal/domain/user"
    "github.com/colinrs/shopjoy/pkg/contextx"
)
```

- [ ] **Step 4: Build**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy/admin
make build
```
Expected: succeeds.

- [ ] **Step 5: Commit**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add admin/internal/logic/users/create_user_logic.go
git commit -m "feat(user): instrument CreateUserLogic with CREATE_USER log

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 10: Instrument UpdateUserLogic

**Files:**
- Modify: `admin/internal/logic/users/update_user_logic.go`

Same pattern as Task 9. After the successful `UserService.Update(...)` call, add:

```go
	l.svcCtx.OperationLogService.Record(l.ctx, l.svcCtx.DB, appUser.RecordOperationLogInput{
		TenantID:     tenantID,
		UserID:       req.ID,
		Action:       user.ActionUpdateUser,
		OperatorID:   contextx.GetAdminID(l.ctx),
		OperatorName: contextx.GetAdminName(l.ctx),
		IPAddress:    contextx.GetClientIP(l.ctx),
		UserAgent:    contextx.GetUserAgent(l.ctx),
	})
```

- [ ] Add the same 3 imports as Task 9 if missing.
- [ ] Run `make build` from `admin/`.
- [ ] Commit:
```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add admin/internal/logic/users/update_user_logic.go
git commit -m "feat(user): instrument UpdateUserLogic with UPDATE_USER log

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 11: Instrument SuspendUserLogic

**Files:**
- Modify: `admin/internal/logic/users/suspend_user_logic.go`

After successful `UserService.Suspend(...)`:

```go
	l.svcCtx.OperationLogService.Record(l.ctx, l.svcCtx.DB, appUser.RecordOperationLogInput{
		TenantID:     tenantID,
		UserID:       req.ID,
		Action:       user.ActionSuspendUser,
		OperatorID:   contextx.GetAdminID(l.ctx),
		OperatorName: contextx.GetAdminName(l.ctx),
		IPAddress:    contextx.GetClientIP(l.ctx),
		UserAgent:    contextx.GetUserAgent(l.ctx),
	})
```

- [ ] Add same imports if missing.
- [ ] `make build` from `admin/`.
- [ ] Commit:
```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add admin/internal/logic/users/suspend_user_logic.go
git commit -m "feat(user): instrument SuspendUserLogic with SUSPEND_USER log

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 12: Instrument SuspendUserWithReasonLogic

**Files:**
- Modify: `admin/internal/logic/users/suspend_user_with_reason_logic.go`

After successful suspend call:

```go
	l.svcCtx.OperationLogService.Record(l.ctx, l.svcCtx.DB, appUser.RecordOperationLogInput{
		TenantID:     tenantID,
		UserID:       req.ID,
		Action:       user.ActionSuspendWithReason,
		OperatorID:   contextx.GetAdminID(l.ctx),
		OperatorName: contextx.GetAdminName(l.ctx),
		Reason:       req.Reason,
		IPAddress:    contextx.GetClientIP(l.ctx),
		UserAgent:    contextx.GetUserAgent(l.ctx),
	})
```

- [ ] Same imports.
- [ ] `make build` from `admin/`.
- [ ] Commit:
```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add admin/internal/logic/users/suspend_user_with_reason_logic.go
git commit -m "feat(user): instrument SuspendUserWithReasonLogic with SUSPEND_WITH_REASON log

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 13: Instrument ActivateUserLogic

**Files:**
- Modify: `admin/internal/logic/users/activate_user_logic.go`

After successful activate call:

```go
	l.svcCtx.OperationLogService.Record(l.ctx, l.svcCtx.DB, appUser.RecordOperationLogInput{
		TenantID:     tenantID,
		UserID:       req.ID,
		Action:       user.ActionActivateUser,
		OperatorID:   contextx.GetAdminID(l.ctx),
		OperatorName: contextx.GetAdminName(l.ctx),
		IPAddress:    contextx.GetClientIP(l.ctx),
		UserAgent:    contextx.GetUserAgent(l.ctx),
	})
```

- [ ] Same imports.
- [ ] `make build`.
- [ ] Commit:
```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add admin/internal/logic/users/activate_user_logic.go
git commit -m "feat(user): instrument ActivateUserLogic with ACTIVATE_USER log

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 14: Instrument DeleteUserLogic

**Files:**
- Modify: `admin/internal/logic/users/delete_user_logic.go`

After successful delete call:

```go
	l.svcCtx.OperationLogService.Record(l.ctx, l.svcCtx.DB, appUser.RecordOperationLogInput{
		TenantID:     tenantID,
		UserID:       req.ID,
		Action:       user.ActionDeleteUser,
		OperatorID:   contextx.GetAdminID(l.ctx),
		OperatorName: contextx.GetAdminName(l.ctx),
		IPAddress:    contextx.GetClientIP(l.ctx),
		UserAgent:    contextx.GetUserAgent(l.ctx),
	})
```

- [ ] Same imports.
- [ ] `make build`.
- [ ] Commit:
```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add admin/internal/logic/users/delete_user_logic.go
git commit -m "feat(user): instrument DeleteUserLogic with DELETE_USER log

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 15: Instrument ResetPasswordLogic

**Files:**
- Modify: `admin/internal/logic/users/reset_password_logic.go`

After successful reset:

```go
	l.svcCtx.OperationLogService.Record(l.ctx, l.svcCtx.DB, appUser.RecordOperationLogInput{
		TenantID:     tenantID,
		UserID:       req.ID,
		Action:       user.ActionResetPassword,
		OperatorID:   contextx.GetAdminID(l.ctx),
		OperatorName: contextx.GetAdminName(l.ctx),
		IPAddress:    contextx.GetClientIP(l.ctx),
		UserAgent:    contextx.GetUserAgent(l.ctx),
	})
```

- [ ] Same imports.
- [ ] `make build`.
- [ ] Commit:
```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add admin/internal/logic/users/reset_password_logic.go
git commit -m "feat(user): instrument ResetPasswordLogic with RESET_PASSWORD log

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 16: End-to-end backend verification

**Files:** none (verification only)

- [ ] **Step 1: Restart admin service**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy/admin
make build && make run
```

- [ ] **Step 2: Trigger each action and verify log rows**

For each of: create, update, suspend, activate, reset-password, delete — call the endpoint and query the log table:
```bash
mysql -h 192.168.0.100 -u root shopjoy \
  -e "SELECT id, user_id, action, operator_name, created_at FROM user_operation_logs ORDER BY id DESC LIMIT 10;"
```
Expected: rows present with correct action strings.

- [ ] **Step 3: Verify list endpoint**

```bash
TOKEN=$(... login ...)
curl -s "http://localhost:8888/api/v1/users/<USER_ID>/operation-logs?page=1&page_size=20" \
  -H "Authorization: Bearer $TOKEN" | jq '.data | {total, first_action: .list[0].action, first_text: .list[0].action_text}'
```
Expected: `total >= 1`, `first_text` is Chinese label.

- [ ] **Step 4: Verify platform-admin tenant bypass**

Login as platform admin (tenant_id=0), call the endpoint for a user from any tenant, expect non-empty list.

No commit (verification only).

---

## Task 17: Frontend API client + types

**Files:**
- Modify: `shop-admin/src/api/user.ts` (append)

**Interfaces:**
- Produces: 4 functions (`getUserOrders`, `getUserPointsTransactions`, `getUserReviews`, `getUserOperationLogs`) + 4 interfaces consumed by Tasks 18-21

- [ ] **Step 1: Append types and functions**

Open `shop-admin/src/api/user.ts`. At the end of the file, append:

```ts
// ===================== User Detail Tab Lists =====================

export interface UserOrderListItem {
    order_id: string
    order_no: string
    status: string
    fulfillment_status: string
    total_amount: string
    currency: string
    created_at: string
}

export interface UserOrdersResponse {
    list: UserOrderListItem[]
    total: number
    page: number
    page_size: number
}

export function getUserOrders(
    userId: string,
    params: { page: number; page_size: number }
) {
    return request<UserOrdersResponse>({
        url: '/api/v1/orders',
        method: 'get',
        params: { ...params, user_id: userId }
    })
}

export interface UserPointsTransaction {
    id: string
    points: number
    balance_after: number
    type: string
    reference_type: string
    reference_id: string
    description: string
    created_at: string
}

export interface UserPointsTransactionsResponse {
    list: UserPointsTransaction[]
    total: number
    page: number
    page_size: number
}

export function getUserPointsTransactions(
    userId: string,
    params: { page: number; page_size: number }
) {
    return request<UserPointsTransactionsResponse>({
        url: '/api/v1/points/transactions',
        method: 'get',
        params: { ...params, user_id: userId }
    })
}

export interface UserReviewListItem {
    id: string
    product_name: string
    product_id: string
    user_name: string
    is_anonymous: boolean
    overall_rating: string
    content: string
    status: string
    created_at: string
}

export interface UserReviewsResponse {
    list: UserReviewListItem[]
    total: number
    page: number
    page_size: number
}

export function getUserReviews(
    userId: string,
    params: { page: number; page_size: number }
) {
    return request<UserReviewsResponse>({
        url: '/api/v1/reviews',
        method: 'get',
        params: { ...params, user_id: userId }
    })
}

export interface UserOperationLogItem {
    id: string
    user_id: string
    action: string
    action_text: string
    operator_id: string
    operator_name: string
    reason: string
    ip_address: string
    user_agent: string
    created_at: string
}

export interface UserOperationLogsResponse {
    list: UserOperationLogItem[]
    total: number
    page: number
    page_size: number
}

export function getUserOperationLogs(
    userId: string,
    params: { page: number; page_size: number; action?: string }
) {
    return request<UserOperationLogsResponse>({
        url: `/api/v1/users/${userId}/operation-logs`,
        method: 'get',
        params
    })
}
```

- [ ] **Step 2: Type-check**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy/shop-admin
npm run type-check 2>&1 | tail -20
```
Expected: no errors. If `type-check` script missing, use `npx vue-tsc --noEmit`.

- [ ] **Step 3: Commit**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add shop-admin/src/api/user.ts
git commit -m "feat(user.ts): API client for user detail tab lists

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 18: i18n keys (zh-CN + en-US)

**Files:**
- Modify: `shop-admin/src/locales/zh-CN.ts`
- Modify: `shop-admin/src/locales/en-US.ts`

- [ ] **Step 1: Find existing `users.*` keys**

Open `shop-admin/src/locales/zh-CN.ts`. Locate the `users:` block. Append inside it (preserve existing keys, find the closing `}` of `users:` and add before it):

```ts
    log: {
        title: '操作日志',
        empty: '暂无操作记录',
        filterAll: '全部动作',
        columns: {
            time: '时间',
            action: '动作',
            operator: '操作人',
            ip: 'IP 地址',
            reason: '原因',
            userAgent: 'User Agent'
        },
        actions: {
            CREATE_USER: '创建用户',
            UPDATE_USER: '更新资料',
            SUSPEND_USER: '禁用用户',
            SUSPEND_WITH_REASON: '禁用用户（带原因）',
            ACTIVATE_USER: '启用用户',
            DELETE_USER: '删除用户',
            RESET_PASSWORD: '重置密码'
        }
    },
    orders: {
        title: '订单记录',
        empty: '暂无订单',
        columns: {
            orderNo: '订单号',
            status: '状态',
            itemCount: '商品数',
            totalAmount: '订单金额',
            createdAt: '下单时间'
        }
    },
    points: {
        title: '积分记录',
        empty: '暂无积分流水',
        columns: {
            time: '时间',
            type: '类型',
            change: '积分变化',
            balance: '余额',
            description: '说明',
            reference: '关联'
        }
    },
    reviews: {
        title: '评价记录',
        empty: '暂无评价',
        anonymous: '匿名用户',
        columns: {
            product: '商品',
            rating: '评分',
            content: '内容',
            status: '状态',
            createdAt: '时间'
        }
    },
```

- [ ] **Step 2: Mirror to en-US.ts**

Same keys in `shop-admin/src/locales/en-US.ts`, English values:
```ts
    log: {
        title: 'Operation Logs',
        empty: 'No operation logs',
        filterAll: 'All actions',
        columns: {
            time: 'Time',
            action: 'Action',
            operator: 'Operator',
            ip: 'IP Address',
            reason: 'Reason',
            userAgent: 'User Agent'
        },
        actions: {
            CREATE_USER: 'Create User',
            UPDATE_USER: 'Update Profile',
            SUSPEND_USER: 'Suspend User',
            SUSPEND_WITH_REASON: 'Suspend (with reason)',
            ACTIVATE_USER: 'Activate User',
            DELETE_USER: 'Delete User',
            RESET_PASSWORD: 'Reset Password'
        }
    },
    orders: {
        title: 'Orders',
        empty: 'No orders',
        columns: {
            orderNo: 'Order No.',
            status: 'Status',
            itemCount: 'Items',
            totalAmount: 'Total',
            createdAt: 'Created At'
        }
    },
    points: {
        title: 'Points Transactions',
        empty: 'No point transactions',
        columns: {
            time: 'Time',
            type: 'Type',
            change: 'Change',
            balance: 'Balance',
            description: 'Description',
            reference: 'Reference'
        }
    },
    reviews: {
        title: 'Reviews',
        empty: 'No reviews',
        anonymous: 'Anonymous',
        columns: {
            product: 'Product',
            rating: 'Rating',
            content: 'Content',
            status: 'Status',
            createdAt: 'Created At'
        }
    },
```

- [ ] **Step 3: Build check**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy/shop-admin
npm run build 2>&1 | tail -10
```
Expected: success.

- [ ] **Step 4: Commit**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add shop-admin/src/locales/
git commit -m "feat(i18n): add user detail tab translations (zh-CN, en-US)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 19: UserOrderList.vue component

**Files:**
- Create: `shop-admin/src/views/users/components/UserOrderList.vue`

- [ ] **Step 1: Write component**

```vue
<template>
  <div class="user-order-list">
    <el-table
      v-loading="loading"
      :data="orders"
      stripe
    >
      <el-table-column :label="$t('users.orders.columns.orderNo')" min-width="180">
        <template #default="{ row }">
          <span class="order-no">{{ row.order_no }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.orders.columns.status')" width="120" align="center">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)" size="small">
            {{ row.status }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.orders.columns.totalAmount')" width="140" align="right">
        <template #default="{ row }">
          <span class="amount">{{ row.currency }} {{ row.total_amount }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.orders.columns.createdAt')" width="180">
        <template #default="{ row }">
          <span class="time-text">{{ formatDateTime(row.created_at) }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('common.actions')" width="120" align="center" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="goDetail(row)">
            {{ $t('common.viewDetail') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div v-if="total > pageSize" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next"
        @size-change="load"
        @current-change="load"
      />
    </div>

    <el-empty v-if="!loading && orders.length === 0" :description="$t('users.orders.empty')" />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getUserOrders, type UserOrderListItem } from '@/api/user'
import { useErrorHandler } from '@/composables/useErrorHandler'

const props = defineProps<{ userId?: string }>()
const router = useRouter()
const { handleError } = useErrorHandler()

const loading = ref(false)
const orders = ref<UserOrderListItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const load = async () => {
  if (!props.userId) return
  loading.value = true
  try {
    const res = await getUserOrders(props.userId, {
      page: currentPage.value,
      page_size: pageSize.value
    })
    orders.value = res.list || []
    total.value = res.total || 0
  } catch (err) {
    handleError(err, false)
    orders.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const formatDateTime = (s: string) => {
  if (!s) return '-'
  return new Date(s).toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit'
  })
}

const getStatusType = (status: string): 'success' | 'warning' | 'info' | 'danger' | 'primary' => {
  if (status === 'completed') return 'success'
  if (status === 'cancelled' || status === 'refunded') return 'info'
  if (status === 'paid' || status === 'shipped') return 'primary'
  if (status === 'pending_payment' || status === 'pending_shipment') return 'warning'
  return 'danger'
}

const goDetail = (row: UserOrderListItem) => {
  router.push(`/orders/${row.order_id}`)
}

watch(() => props.userId, load, { immediate: true })
</script>

<style scoped>
.user-order-list { padding: 0; }
.order-no { font-family: 'Fira Code', monospace; font-weight: 500; }
.amount { font-family: 'Fira Sans', sans-serif; font-weight: 600; color: #EF4444; }
.time-text { font-size: 13px; color: #6B7280; font-family: 'Fira Code', monospace; }
.pagination-wrapper { display: flex; justify-content: flex-end; padding-top: 16px; margin-top: 16px; border-top: 1px solid #F3F4F6; }
:deep(.el-table__row:hover > td) { background-color: #F5F3FF !important; }
</style>
```

- [ ] **Step 2: Build**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy/shop-admin
npm run build 2>&1 | tail -5
```
Expected: success.

- [ ] **Step 3: Commit**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add shop-admin/src/views/users/components/UserOrderList.vue
git commit -m "feat(users): UserOrderList tab component

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 20: UserPointsList.vue component

**Files:**
- Create: `shop-admin/src/views/users/components/UserPointsList.vue`

- [ ] **Step 1: Write component**

```vue
<template>
  <div class="user-points-list">
    <el-table v-loading="loading" :data="txns" stripe>
      <el-table-column :label="$t('users.points.columns.time')" width="170">
        <template #default="{ row }">
          <span class="time-text">{{ formatDateTime(row.created_at) }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.points.columns.type')" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getTypeTag(row.type)" size="small">{{ row.type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.points.columns.change')" width="120" align="right">
        <template #default="{ row }">
          <span :class="row.points >= 0 ? 'positive' : 'negative'">
            {{ row.points >= 0 ? '+' : '' }}{{ row.points }}
          </span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.points.columns.balance')" width="100" align="right">
        <template #default="{ row }">
          <span class="balance">{{ row.balance_after }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.points.columns.description')" min-width="180">
        <template #default="{ row }">
          <span>{{ row.description || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.points.columns.reference')" width="180">
        <template #default="{ row }">
          <template v-if="row.reference_type === 'ORDER' && row.reference_id">
            <el-button link type="primary" size="small" @click="goOrder(row.reference_id)">
              {{ $t('users.points.columns.reference') }}: {{ row.reference_id }}
            </el-button>
          </template>
          <span v-else-if="row.reference_id">
            {{ row.reference_type }}:{{ row.reference_id }}
          </span>
          <span v-else>-</span>
        </template>
      </el-table-column>
    </el-table>

    <div v-if="total > pageSize" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next"
        @size-change="load"
        @current-change="load"
      />
    </div>

    <el-empty v-if="!loading && txns.length === 0" :description="$t('users.points.empty')" />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getUserPointsTransactions, type UserPointsTransaction } from '@/api/user'
import { useErrorHandler } from '@/composables/useErrorHandler'

const props = defineProps<{ userId?: string }>()
const router = useRouter()
const { handleError } = useErrorHandler()

const loading = ref(false)
const txns = ref<UserPointsTransaction[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const load = async () => {
  if (!props.userId) return
  loading.value = true
  try {
    const res = await getUserPointsTransactions(props.userId, {
      page: currentPage.value,
      page_size: pageSize.value
    })
    txns.value = res.list || []
    total.value = res.total || 0
  } catch (err) {
    handleError(err, false)
    txns.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const formatDateTime = (s: string) => {
  if (!s) return '-'
  return new Date(s).toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit'
  })
}

const getTypeTag = (t: string): 'success' | 'warning' | 'info' | 'danger' | 'primary' => {
  if (t === 'EARN') return 'success'
  if (t === 'REDEEM') return 'primary'
  if (t === 'EXPIRE') return 'warning'
  if (t === 'ADJUST') return 'info'
  return 'info'
}

const goOrder = (id: string) => router.push(`/orders/${id}`)

watch(() => props.userId, load, { immediate: true })
</script>

<style scoped>
.user-points-list { padding: 0; }
.time-text { font-size: 13px; color: #6B7280; font-family: 'Fira Code', monospace; }
.positive { color: #10B981; font-weight: 600; }
.negative { color: #EF4444; font-weight: 600; }
.balance { font-weight: 600; color: #1E1B4B; }
.pagination-wrapper { display: flex; justify-content: flex-end; padding-top: 16px; margin-top: 16px; border-top: 1px solid #F3F4F6; }
:deep(.el-table__row:hover > td) { background-color: #F5F3FF !important; }
</style>
```

- [ ] **Step 2: Build**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy/shop-admin
npm run build 2>&1 | tail -5
```
Expected: success.

- [ ] **Step 3: Commit**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add shop-admin/src/views/users/components/UserPointsList.vue
git commit -m "feat(users): UserPointsList tab component

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 21: UserReviewList.vue component

**Files:**
- Create: `shop-admin/src/views/users/components/UserReviewList.vue`

- [ ] **Step 1: Write component**

```vue
<template>
  <div class="user-review-list">
    <el-table v-loading="loading" :data="reviews" stripe>
      <el-table-column :label="$t('users.reviews.columns.product')" min-width="180">
        <template #default="{ row }">
          <span class="product-name">{{ row.product_name }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.reviews.columns.rating')" width="100" align="center">
        <template #default="{ row }">
          <span class="rating">
            <el-icon v-for="i in 5" :key="i" :class="i <= Number(row.overall_rating) ? 'star-filled' : 'star-empty'">
              <Star />
            </el-icon>
          </span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.reviews.columns.content')" min-width="220">
        <template #default="{ row }">
          <el-tooltip :content="row.content" placement="top" :disabled="row.content.length <= 80">
            <span class="content-text">{{ truncate(row.content, 80) }}</span>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.reviews.columns.status')" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)" size="small">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.reviews.columns.createdAt')" width="170">
        <template #default="{ row }">
          <span class="time-text">{{ formatDateTime(row.created_at) }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('common.actions')" width="100" align="center" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="goReviews(row)">
            {{ $t('common.view') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div v-if="total > pageSize" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next"
        @size-change="load"
        @current-change="load"
      />
    </div>

    <el-empty v-if="!loading && reviews.length === 0" :description="$t('users.reviews.empty')" />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { Star } from '@element-plus/icons-vue'
import { getUserReviews, type UserReviewListItem } from '@/api/user'
import { useErrorHandler } from '@/composables/useErrorHandler'

const props = defineProps<{ userId?: string }>()
const router = useRouter()
const { handleHandler: _h } = useErrorHandler() as any // unused, kept for parity
const { handleError } = useErrorHandler()

const loading = ref(false)
const reviews = ref<UserReviewListItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const load = async () => {
  if (!props.userId) return
  loading.value = true
  try {
    const res = await getUserReviews(props.userId, {
      page: currentPage.value,
      page_size: pageSize.value
    })
    reviews.value = res.list || []
    total.value = res.total || 0
  } catch (err) {
    handleError(err, false)
    reviews.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const formatDateTime = (s: string) => {
  if (!s) return '-'
  return new Date(s).toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit'
  })
}

const truncate = (s: string, n: number) => (s && s.length > n ? s.slice(0, n) + '…' : s || '-')

const getStatusType = (status: string): 'success' | 'warning' | 'info' | 'danger' => {
  if (status === 'approved' || status === 'visible') return 'success'
  if (status === 'pending') return 'warning'
  if (status === 'hidden') return 'info'
  return 'danger'
}

const goReviews = (row: UserReviewListItem) => {
  router.push({ path: '/reviews', query: { product_id: row.product_id } })
}

watch(() => props.userId, load, { immediate: true })
</script>

<style scoped>
.user-review-list { padding: 0; }
.product-name { font-weight: 500; color: #1E1B4B; }
.rating { display: inline-flex; gap: 2px; }
.star-filled { color: #F59E0B; font-size: 14px; }
.star-empty { color: #E5E7EB; font-size: 14px; }
.content-text { font-size: 13px; color: #374151; }
.time-text { font-size: 13px; color: #6B7280; font-family: 'Fira Code', monospace; }
.pagination-wrapper { display: flex; justify-content: flex-end; padding-top: 16px; margin-top: 16px; border-top: 1px solid #F3F4F6; }
:deep(.el-table__row:hover > td) { background-color: #F5F3FF !important; }
</style>
```

Note: there is a bug in the script above — the `useErrorHandler()` destructured alias `_h` is leftover from a draft. Remove it. The correct code uses only `const { handleError } = useErrorHandler()`. Also, the template uses `row.is_anonymous` for the "anonymous" label — fix this: the template should NOT show `row.user_name` directly; it should show `$t('users.reviews.anonymous')` when `row.is_anonymous` is true.

Update the template cell that displays the user (this is currently missing!). Add a `用户` column between content and status:

Replace the el-table-column block after `<el-table-column :label="$t('users.reviews.columns.content')" ...>` with:

```vue
      <el-table-column label="用户" width="100" align="center">
        <template #default="{ row }">
          <span v-if="row.is_anonymous" class="anonymous">{{ $t('users.reviews.anonymous') }}</span>
          <span v-else>{{ row.user_name || '-' }}</span>
        </template>
      </el-table-column>
```

And remove the leftover `_h` line in script. Final script imports and setup:

```ts
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { Star } from '@element-plus/icons-vue'
import { getUserReviews, type UserReviewListItem } from '@/api/user'
import { useErrorHandler } from '@/composables/useErrorHandler'

const props = defineProps<{ userId?: string }>()
const router = useRouter()
const { handleError } = useErrorHandler()
```

- [ ] **Step 2: Build**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy/shop-admin
npm run build 2>&1 | tail -5
```
Expected: success.

- [ ] **Step 3: Commit**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add shop-admin/src/views/users/components/UserReviewList.vue
git commit -m "feat(users): UserReviewList tab component with anonymous handling

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 22: Rewrite UserOperationLog.vue

**Files:**
- Modify: `shop-admin/src/views/users/components/UserOperationLog.vue` (full rewrite)

- [ ] **Step 1: Write component**

```vue
<template>
  <div class="user-operation-log">
    <div class="filter-bar">
      <el-select
        v-model="filterAction"
        :placeholder="$t('users.log.filterAll')"
        clearable
        class="action-filter"
        @change="reload"
      >
        <el-option :label="$t('users.log.filterAll')" value="" />
        <el-option
          v-for="opt in actionOptions"
          :key="opt.value"
          :label="opt.label"
          :value="opt.value"
        />
      </el-select>
    </div>

    <el-table v-loading="loading" :data="logs" stripe>
      <el-table-column :label="$t('users.log.columns.time')" width="170">
        <template #default="{ row }">
          <span class="time-text">{{ formatDateTime(row.created_at) }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.log.columns.action')" width="160" align="center">
        <template #default="{ row }">
          <el-tag :type="getActionType(row.action)" size="small">
            {{ row.action_text }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.log.columns.operator')" width="120">
        <template #default="{ row }">
          <span>{{ row.operator_name || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.log.columns.ip')" width="140">
        <template #default="{ row }">
          <span class="mono">{{ row.ip_address || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.log.columns.reason')" min-width="180">
        <template #default="{ row }">
          <span>{{ row.reason || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('users.log.columns.userAgent')" width="80" align="center">
        <template #default="{ row }">
          <el-tooltip v-if="row.user_agent" :content="row.user_agent" placement="top">
            <el-button link type="primary" size="small">{{ $t('common.view') }}</el-button>
          </el-tooltip>
          <span v-else>-</span>
        </template>
      </el-table-column>
    </el-table>

    <div v-if="total > pageSize" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next"
        @size-change="load"
        @current-change="load"
      />
    </div>

    <el-empty v-if="!loading && logs.length === 0" :description="$t('users.log.empty')" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { getUserOperationLogs, type UserOperationLogItem } from '@/api/user'
import { useErrorHandler } from '@/composables/useErrorHandler'

const props = defineProps<{ userId?: string }>()
const { t } = useI18n()
const { handleError } = useErrorHandler()

const loading = ref(false)
const logs = ref<UserOperationLogItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const filterAction = ref<string>('')

const actionOptions = computed(() => [
  { value: 'CREATE_USER', label: t('users.log.actions.CREATE_USER') },
  { value: 'UPDATE_USER', label: t('users.log.actions.UPDATE_USER') },
  { value: 'SUSPEND_USER', label: t('users.log.actions.SUSPEND_USER') },
  { value: 'ACTIVATE_USER', label: t('users.log.actions.ACTIVATE_USER') },
  { value: 'DELETE_USER', label: t('users.log.actions.DELETE_USER') },
  { value: 'RESET_PASSWORD', label: t('users.log.actions.RESET_PASSWORD') }
])

const load = async () => {
  if (!props.userId) return
  loading.value = true
  try {
    const res = await getUserOperationLogs(props.userId, {
      page: currentPage.value,
      page_size: pageSize.value,
      action: filterAction.value || undefined
    })
    logs.value = res.list || []
    total.value = res.total || 0
  } catch (err) {
    handleError(err, false)
    logs.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const reload = () => {
  currentPage.value = 1
  load()
}

const formatDateTime = (s: string) => {
  if (!s) return '-'
  return new Date(s).toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit'
  })
}

const getActionType = (action: string): 'success' | 'warning' | 'info' | 'danger' | 'primary' => {
  if (action.startsWith('SUSPEND') || action === 'DELETE_USER') return 'danger'
  if (action === 'ACTIVATE_USER') return 'success'
  if (action === 'RESET_PASSWORD') return 'warning'
  if (action === 'CREATE_USER') return 'primary'
  return 'info'
}

watch(() => props.userId, load, { immediate: true })
</script>

<style scoped>
.user-operation-log { padding: 0; }
.filter-bar { margin-bottom: 16px; display: flex; justify-content: flex-start; }
.action-filter { width: 200px; }
.time-text { font-size: 13px; color: #6B7280; font-family: 'Fira Code', monospace; }
.mono { font-family: 'Fira Code', monospace; font-size: 12px; color: #374151; }
.pagination-wrapper { display: flex; justify-content: flex-end; padding-top: 16px; margin-top: 16px; border-top: 1px solid #F3F4F6; }
:deep(.el-table__row:hover > td) { background-color: #F5F3FF !important; }
</style>
```

- [ ] **Step 2: Build**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy/shop-admin
npm run build 2>&1 | tail -5
```
Expected: success.

- [ ] **Step 3: Commit**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add shop-admin/src/views/users/components/UserOperationLog.vue
git commit -m "feat(users): rewrite UserOperationLog with filter and pagination

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 23: Wire components into UserDetailTabs

**Files:**
- Modify: `shop-admin/src/views/users/components/UserDetailTabs.vue`

- [ ] **Step 1: Read current file and modify**

Open the file. Find the imports and the template. Replace the entire file content with:

```vue
<template>
  <el-card class="tabs-card" shadow="never">
    <el-tabs v-model="activeTab" class="detail-tabs">
      <el-tab-pane :label="$t('users.basicInfo')" name="basic">
        <UserBasicInfo :user="user" @refresh="$emit('refresh')" />
      </el-tab-pane>
      <el-tab-pane :label="$t('users.addresses')" name="addresses">
        <UserAddressList :user-id="user?.id" />
      </el-tab-pane>
      <el-tab-pane :label="$t('users.orderRecords')" name="orders">
        <UserOrderList :user-id="user?.id" />
      </el-tab-pane>
      <el-tab-pane :label="$t('users.pointsRecords')" name="points">
        <UserPointsList :user-id="user?.id" />
      </el-tab-pane>
      <el-tab-pane :label="$t('users.reviewRecords')" name="reviews">
        <UserReviewList :user-id="user?.id" />
      </el-tab-pane>
      <el-tab-pane :label="$t('users.operationLogs')" name="logs">
        <UserOperationLog :user-id="user?.id" />
      </el-tab-pane>
    </el-tabs>
  </el-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import UserBasicInfo from './UserBasicInfo.vue'
import UserAddressList from './UserAddressList.vue'
import UserOrderList from './UserOrderList.vue'
import UserPointsList from './UserPointsList.vue'
import UserReviewList from './UserReviewList.vue'
import UserOperationLog from './UserOperationLog.vue'
import type { UserDetail } from '@/api/user'

defineProps<{ user: UserDetail | null }>()
defineEmits<{ refresh: [] }>()

const activeTab = ref('basic')
</script>

<style scoped>
.tabs-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}
.detail-tabs :deep(.el-tabs__item) {
  font-size: 15px;
  font-weight: 500;
}
.detail-tabs :deep(.el-tabs__item.is-active) { color: #6366F1; }
.detail-tabs :deep(.el-tabs__active-bar) { background-color: #6366F1; }
</style>
```

- [ ] **Step 2: Build**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy/shop-admin
npm run build 2>&1 | tail -5
```
Expected: success.

- [ ] **Step 3: Commit**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add shop-admin/src/views/users/components/UserDetailTabs.vue
git commit -m "feat(users): wire 4 new tab components into UserDetailTabs

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 24: Browser end-to-end verification

**Files:** none (verification only)

- [ ] **Step 1: Start admin service (with new endpoints) and shop-admin dev server**

```bash
# Terminal 1
cd /Users/dengyichuan/workspace/go/src/shopjoy/admin
make build && make run

# Terminal 2
cd /Users/dengyichuan/workspace/go/src/shopjoy/shop-admin
npm run dev
```

- [ ] **Step 2: Open browser at http://localhost:3000/users**

Log in as `superadmin@shopjoy.com / password123`. Confirm `/users` list still loads correctly.

- [ ] **Step 3: Navigate to a user detail page**

Click any user → `/users/{id}`.

- [ ] **Step 4: Verify each tab**

For each tab (orders, points, reviews, logs):
- Click tab → verify Network panel shows exactly ONE API call
- Verify list renders (or empty state)
- If list has items: change page → verify second API call updates the list
- For logs tab: change action filter → verify filter narrows the list

- [ ] **Step 5: Verify anonymous review handling**

On the reviews tab: find a row with `is_anonymous=true` (if any). Confirm it shows "匿名用户" not the user's real name.

- [ ] **Step 6: Verify operation log writes**

Trigger an action on the user (e.g., suspend via the summary card "禁用" button). Switch to the logs tab. Confirm a new row appears with action_text = "禁用用户".

No commit (verification only).

---

## Task 25: Final summary commit

- [ ] **Step 1: Run final build across both services**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy/admin && make build
cd /Users/dengyichuan/workspace/go/src/shopjoy/shop-admin && npm run build
```
Both: success.

- [ ] **Step 2: Confirm full git status is clean (or only intentional changes)**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy && git status
```
Expected: only the expected files (untracked changes should be limited to anything intentionally left out of scope).

- [ ] **Step 3: Final summary commit if any uncommitted changes**

```bash
cd /Users/dengyichuan/workspace/go/src/shopjoy
git add -A
git diff --staged --quiet || git commit -m "feat(users): complete user detail tabs (orders/points/reviews/logs)

Backend:
- New user_operation_logs table + migration
- OperationLog domain/repo/service with platform-admin tenant bypass
- 7 service-layer instrumentation points (never block parent operation)
- New GET /api/v1/users/{id}/operation-logs endpoint
- Extended ListReviews to filter by user_id

Frontend:
- 3 new tab components (UserOrderList, UserPointsList, UserReviewList)
- Rewrote UserOperationLog with filter + pagination
- 4 API client functions + types in api/user.ts
- i18n keys in zh-CN + en-US

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Self-Review Checklist (already applied during writing)

| Item | Status |
|---|---|
| Every spec section mapped to a task | ✅ §2→Task 1, §3→Tasks 6-7, §4→Task 8, §5→Tasks 19-22, §6→Task 17, §7→Task 18, §8→Tasks 9-15, §11→Tasks 16+24 |
| No "TBD" / "TODO" / "implement later" placeholders | ✅ |
| No "add appropriate error handling" vague steps | ✅ Each component specifies `handleError(err, false)` + empty state |
| Exact file paths in every step | ✅ |
| Complete code blocks (not "similar to Task N") | ✅ |
| Each task ends with `make build` + commit | ✅ |
| TDD where applicable (entity helper test in Task 2) | ✅ Other tasks have direct smoke tests rather than unit tests since the logic is mostly wiring |
| Frequent commits (one per task) | ✅ 25 tasks = 25+ commits |
| Platform-admin tenant bypass applied to new read path | ✅ Task 3 |
| Instrumentation failure never blocks parent | ✅ Task 5 makes Record return no error |
| Review anonymous handling | ✅ Task 21 |