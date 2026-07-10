# User Detail Tabs — Design Spec

**Date:** 2026-07-10
**Status:** Approved (brainstorming phase)
**Owner:** colinrs

## Goal

Complete the four placeholder tabs on `/users/{id}` detail page:

| Tab | Today | Target |
|---|---|---|
| 基本信息 | Wired via `getUserDetail` | No change |
| 收货地址 | Wired via `getUserAddresses` | No change |
| 订单记录 | Coming Soon | Real list, paged |
| 积分记录 | Coming Soon | Real list, paged |
| 评价记录 | Coming Soon | Real list, paged, anonymous handled |
| 操作日志 | Coming Soon | Real list, paged, action filter |

All lists must load lazily on tab activation, support pagination, handle empty/error states, and respect platform-admin tenant bypass (`tenantID == 0`).

## Non-Goals

- Server-side search/filter beyond what existing APIs already support
- Real-time updates (websocket / SSE)
- Action log write backfill for historical actions before 2026-07-10
- Bulk operations on any of the new lists

---

## §1 Architecture

```
┌──────────────────────────────────────────────────────────────┐
│  shop-admin/src/views/users/components/                      │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ UserDetailTabs.vue (orchestrator)                      │  │
│  │  ├─ basic     → UserBasicInfo.vue  (props only)        │  │
│  │  ├─ addresses → UserAddressList.vue (existing API)     │  │
│  │  ├─ orders    → UserOrderList.vue   [NEW]              │  │
│  │  ├─ points    → UserPointsList.vue  [NEW]              │  │
│  │  ├─ reviews   → UserReviewList.vue  [NEW]              │  │
│  │  └─ logs      → UserOperationLog.vue [REWRITE]         │  │
│  └────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
            │  lazy load on tab first activation
            ▼
┌──────────────────────────────────────────────────────────────┐
│  Backend (Go)                                                 │
│  ┌──────────────────────┐  ┌─────────────────────────────┐   │
│  │ GET /orders?uid=     │  │ GET /points/txns?uid=       │   │
│  │ (existing, reuse)    │  │ (existing, reuse)           │   │
│  └──────────────────────┘  └─────────────────────────────┘   │
│  ┌──────────────────────┐  ┌─────────────────────────────┐   │
│  │ GET /reviews?uid=    │  │ GET /users/{id}/op-logs NEW │   │
│  │ (extend req +1 field)│  │ (new table+model+repo+srv)  │   │
│  └──────────────────────┘  └─────────────────────────────┘   │
└──────────────────────────────────────────────────────────────┘
```

`UserDetailTabs.vue` is a thin container. Each tab component owns its own state, request lifecycle, and column definitions. No data is passed through props except `userId`.

---

## §2 Operation Log Data Model

### New table `user_operation_logs`

```sql
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

    INDEX idx_uol_user_id     (user_id, created_at),
    INDEX idx_uol_tenant      (tenant_id, created_at),
    INDEX idx_uol_action      (action)
);
```

Migration file: `sql/user/migrations/2026071001_create_user_operation_logs.sql`

Schema merge: also appended to `sql/user/schema.sql` per `CLAUDE.md` rules.

### Action enum (single source of truth)

| action string | Trigger location | i18n key |
|---|---|---|
| `CREATE_USER` | `CreateUserLogic` | `users.log.actions.CREATE_USER` |
| `UPDATE_USER` | `UpdateUserLogic` | `users.log.actions.UPDATE_USER` |
| `SUSPEND_USER` | `SuspendUserLogic` | `users.log.actions.SUSPEND_USER` |
| `SUSPEND_WITH_REASON` | `SuspendUserWithReasonLogic` | `users.log.actions.SUSPEND_USER` |
| `ACTIVATE_USER` | `ActivateUserLogic` | `users.log.actions.ACTIVATE_USER` |
| `DELETE_USER` | `DeleteUserLogic` | `users.log.actions.DELETE_USER` |
| `RESET_PASSWORD` | `ResetPasswordLogic` | `users.log.actions.RESET_PASSWORD` |

---

## §3 New Endpoint: `GET /api/v1/users/{id}/operation-logs`

Add to `admin/desc/user.api`:

```go
type (
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
)
```

Add to `users` server group:

```go
@doc "获取用户操作日志"
@handler ListUserOperationLogsHandler
get /api/v1/users/:id/operation-logs (ListUserOperationLogsReq) returns (ListUserOperationLogsResp)
```

Action text mapping happens server-side (zh-CN for now).

---

## §4 Review API Extension

In `admin/desc/review.api`, add one field to `ListReviewsReq`:

```go
ListReviewsReq {
    // ... existing fields ...
    UserID int64 `form:"user_id,optional"` // NEW
}
```

Repository: `if req.UserID > 0 { query = query.Where("user_id = ?", req.UserID) }`.

---

## §5 Frontend Components

### `UserOrderList.vue` (new)

- Columns: 订单号 / 状态 / 商品数 / 总金额 / 下单时间 / 操作（查看详情链接到 `/orders/{id}`）
- Pagination: `el-pagination`, page_size=10
- API: `GET /api/v1/orders?user_id={id}&page=&page_size=`
- Status badge colors per existing `getStatusType` helper

### `UserPointsList.vue` (new)

- Columns: 时间 / 类型 (EARN/REDEEM/EXPIRE/ADJUST/FREEZE/UNFREEZE) / 积分变化 (正绿负红, +x/−x) / 余额 / 说明 / 关联（reference_type:id）
- Pagination: page_size=10
- API: `GET /api/v1/points/transactions?user_id={id}&page=&page_size=`

### `UserReviewList.vue` (new)

- Columns: 商品 / 评分 (1-5 stars) / 内容 (truncate 80 chars) / 状态 / 时间 / 操作（查看链接到 `/reviews`）
- Anonymous reviews: replace `user_name` with `i18n.users.reviews.anonymous`
- Pagination: page_size=10
- API: `GET /api/v1/reviews?user_id={id}&page=&page_size=` (after backend extension)

### `UserOperationLog.vue` (rewrite)

- Columns: 时间 / 动作 (el-tag with color per action) / 操作人 / IP / 原因 / UA (tooltip)
- Top filter: action dropdown (全部 + 7 个 action options)
- Pagination: page_size=10
- API: `GET /api/v1/users/{id}/operation-logs?page=&page_size=&action=`

### Lazy loading

Each list component uses `watch(() => props.active, ...)` (or `onMounted`) to load only when its tab is first opened. Already-loaded data persists across tab switches — does not re-fetch unless pagination or filter changes.

---

## §6 API Client & Types

Append to `src/api/user.ts` (single file for now; can split later if it grows):

```ts
export interface UserOrderListItem {
    order_id: string
    order_no: string
    status: string
    fulfillment_status: string
    total_amount: string
    currency: string
    created_at: string
    item_count?: number
}

export interface UserPointsTransaction {
    id: string
    points: number
    balance_after: number
    type: string          // EARN, REDEEM, ADJUST, EXPIRE, FREEZE, UNFREEZE
    reference_type: string
    reference_id: string
    description: string
    created_at: string
}

export interface UserReviewListItem {
    id: string
    product_name: string
    user_name: string
    is_anonymous: boolean
    overall_rating: string
    content: string
    status: string
    created_at: string
}

export interface UserOperationLog {
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

export function getUserOrders(userId: string, params: { page: number; page_size: number })
export function getUserPointsTransactions(userId: string, params: { page: number; page_size: number })
export function getUserReviews(userId: string, params: { page: number; page_size: number })
export function getUserOperationLogs(userId: string, params: { page: number; page_size: number; action?: string })
```

---

## §7 i18n Keys

`zh-CN` additions under `users.*`. `en-US` mirror with English labels.

```
users:
  log:
    title: '操作日志'
    empty: '暂无操作记录'
    filter:
      all: '全部动作'
    columns:
      time: '时间'
      action: '动作'
      operator: '操作人'
      ip: 'IP 地址'
      reason: '原因'
      userAgent: 'User Agent'
    actions:
      CREATE_USER: '创建用户'
      UPDATE_USER: '更新资料'
      SUSPEND_USER: '禁用用户'
      ACTIVATE_USER: '启用用户'
      DELETE_USER: '删除用户'
      RESET_PASSWORD: '重置密码'
  orders:
    title: '订单记录'
    empty: '暂无订单'
    columns:
      orderNo: '订单号'
      status: '状态'
      itemCount: '商品数'
      totalAmount: '订单金额'
      createdAt: '下单时间'
  points:
    title: '积分记录'
    empty: '暂无积分流水'
    columns:
      time: '时间'
      type: '类型'
      change: '积分变化'
      balance: '余额'
      description: '说明'
      reference: '关联'
  reviews:
    title: '评价记录'
    empty: '暂无评价'
    anonymous: '匿名用户'
    columns:
      product: '商品'
      rating: '评分'
      content: '内容'
      status: '状态'
      createdAt: '时间'
```

---

## §8 Instrumentation Pattern (埋点)

In each user-related logic method that mutates state, after the main write succeeds:

```go
// SuspendUserLogic example
if err := l.svc.UserOperationLog.Create(ctx, l.svc.DB, &user.OperationLog{
    TenantID:     req.TenantID,
    UserID:       req.ID,
    Action:       "SUSPEND_USER",
    OperatorID:   contextx.GetAdminID(ctx),
    OperatorName: contextx.GetAdminName(ctx),
    IPAddress:    contextx.GetClientIP(ctx),
    UserAgent:    contextx.GetUserAgent(ctx),
    Reason:       "",
}); err != nil {
    logx.Errorf("record operation log failed: %v", err)
    // Intentional: do NOT fail the parent operation
}
return nil
```

**Key rule:** instrumentation failure never blocks the parent business operation. Log the error and continue.

Touch points (one per action):
- `CreateUserLogic` → CREATE_USER
- `UpdateUserLogic` → UPDATE_USER
- `SuspendUserLogic` → SUSPEND_USER
- `SuspendUserWithReasonLogic` → SUSPEND_WITH_REASON (carries `reason`)
- `ActivateUserLogic` → ACTIVATE_USER
- `DeleteUserLogic` → DELETE_USER
- `ResetPasswordLogic` → RESET_PASSWORD

---

## §9 Platform Admin Compatibility

Following the established convention in `user_repository.go`, all read paths must accept `tenantID == 0` as "platform admin, no tenant filter":

- New repository `OperationLogRepository.FindByUserID(tenantID, userID, query)`:
  ```go
  if tenantID != 0 {
      query = query.Where("tenant_id = ?", tenantID.Int64())
  }
  ```
- Review list repository extension applies same filter
- Frontend: no change — relies on existing auth/tenant propagation

---

## §10 Error Handling

### Backend
- Repository errors: propagate as-is, wrapped with operation context
- Logic errors: use standardized `code.ErrXxx` (per `error.md` rules)
- Instrumentation errors: log only, never block parent

### Frontend
- Each list component uses `useErrorHandler`
- On error: show `el-empty` with error description (not ElMessage toast — avoids spam when switching tabs)
- Empty state: dedicated message per tab

---

## §11 Testing & Acceptance

### Backend
- Migration replays cleanly via `make migrate` (or equivalent)
- `UserOperationLogRepository` unit tests: `Create`, `FindByUserID` with tenant filter
- Logic unit tests: each of the 7 touch points writes one log on happy path
- `make build` passes
- Manual `curl` smoke against running admin service: each endpoint returns expected shape

### Frontend
- `make build` passes (admin + shop-admin)
- Manual browser walkthrough with platform admin login:
  - Tab click triggers exactly one API call (verified in Network panel)
  - Pagination changes hit the API and update the list
  - Action filter on logs dropdown narrows the list
  - Empty states render when user has no orders / points / reviews / logs
  - Anonymous review shows "匿名用户" not the real name

---

## §12 File Touch List

### Backend (Go) — create
- `sql/user/migrations/2026071001_create_user_operation_logs.sql`
- `sql/user/schema.sql` (append table)
- `admin/internal/domain/user/operation_log_entity.go`
- `admin/internal/infrastructure/persistence/user_operation_log_repository.go`
- `admin/internal/application/user/operation_log_service.go`
- `admin/internal/logic/users/list_user_operation_logs_logic.go`
- `admin/internal/domain/user/operation_log_helper_test.go` (optional)

### Backend — modify
- `admin/desc/user.api` (add types + handler)
- `admin/internal/application/user/service_impl.go` (wire OperationLog service)
- `admin/internal/logic/users/{create,update,suspend,suspend_with_reason,activate,delete,reset_password}_user_logic.go` (埋点)
- `admin/internal/domain/review/review_repository.go` (add `user_id` filter)

### Frontend (Vue/TS) — create
- `shop-admin/src/views/users/components/UserOrderList.vue`
- `shop-admin/src/views/users/components/UserPointsList.vue`
- `shop-admin/src/views/users/components/UserReviewList.vue`

### Frontend — modify
- `shop-admin/src/views/users/components/UserDetailTabs.vue` (swap coming-soon for new components)
- `shop-admin/src/views/users/components/UserOperationLog.vue` (rewrite from placeholder)
- `shop-admin/src/api/user.ts` (append 4 functions + 4 interfaces)
- `shop-admin/src/locales/zh-CN.ts` (add keys)
- `shop-admin/src/locales/en-US.ts` (mirror)

---

## §13 Out of Scope / Future Work

- Bulk actions (delete multiple logs, etc.)
- Export logs to CSV
- Real-time log streaming via SSE
- Diff view showing before/after for UPDATE_USER actions
- Action log retention policy / auto-purge