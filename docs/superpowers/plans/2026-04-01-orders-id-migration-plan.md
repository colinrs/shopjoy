# Orders ID Migration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Change orders.id from VARCHAR(64) to BIGINT AUTO_INCREMENT and update all foreign key references across the codebase.

**Architecture:** The primary key `orders.id` changes from VARCHAR(64) to BIGINT AUTO_INCREMENT. The `order_no` field remains VARCHAR(64) as the human-readable business identifier. All `order_id` foreign key fields in related tables change from VARCHAR(64) to BIGINT to maintain referential integrity.

**Tech Stack:** Go (go-zero framework), MySQL, TypeScript (Vue), goctl API generation

---

## Task 1: Update SQL Schema Files

**Files:**
- Modify: `sql/order/schema.sql:5-60` (orders table, order_items table)
- Modify: `sql/fulfillment/schema.sql:5-35` (shipments table)
- Modify: `sql/fulfillment/schema.sql:62-92` (refunds table)
- Modify: `sql/review/schema.sql:5-32` (reviews table)
- Modify: `sql/promotion/schema.sql:60-81` (promotion_usage table)
- Modify: `sql/promotion/schema.sql:123-142` (user_coupons table)

**Note:** `sql/payment/schema.sql` contains a deprecated `payments` table with VARCHAR order_id, but this table is marked as deprecated in comments. The active `order_payments` table already uses BIGINT. Only update if needed for consistency.

### Step 1: Update orders table in sql/order/schema.sql

Change line 6 from:
```sql
`id` VARCHAR(64) NOT NULL COMMENT '订单ID',
```
To:
```sql
`id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '订单ID',
```

Change line 52 from:
```sql
PRIMARY KEY (`id`),
```
To:
```sql
PRIMARY KEY (`id`),
```

Also remove `UNIQUE KEY uk_order_no` constraint modification since order_no remains VARCHAR.

### Step 2: Update order_items table in sql/order/schema.sql

Change line 68 from:
```sql
`order_id` VARCHAR(64) NOT NULL COMMENT '订单ID',
```
To:
```sql
`order_id` BIGINT NOT NULL COMMENT '订单ID',
```

### Step 3: Update test data INSERT statements in sql/order/schema.sql

The INSERT statement at line 140 uses string IDs like `'ORD202503010001'`. Change to use numeric IDs (1, 2, 3, 4, 5, 6) matching BIGINT AUTO_INCREMENT.

Update all order_items INSERT to use numeric order_id values (1, 2, 3, etc.).

### Step 4: Update shipments table in sql/fulfillment/schema.sql

Change line 8 from:
```sql
`order_id` VARCHAR(64) NOT NULL COMMENT '订单ID',
```
To:
```sql
`order_id` BIGINT NOT NULL COMMENT '订单ID',
```

Update test data INSERT at line 178-186 to use numeric order_id values.

### Step 5: Update refunds table in sql/fulfillment/schema.sql

Change line 65 from:
```sql
`order_id` VARCHAR(64) NOT NULL COMMENT '订单ID',
```
To:
```sql
`order_id` BIGINT NOT NULL COMMENT '订单ID',
```

Update test data INSERT at line 197-202 to use numeric order_id values.

### Step 6: Update reviews table in sql/review/schema.sql

Change line 8 from:
```sql
`order_id` VARCHAR(64) NOT NULL,
```
To:
```sql
`order_id` BIGINT NOT NULL,
```

Update test data INSERT at line 82-92 to use numeric order_id values.

### Step 7: Update promotion_usage table in sql/promotion/schema.sql

Change line 65 from:
```sql
`order_id` VARCHAR(50) NOT NULL,
```
To:
```sql
`order_id` BIGINT NOT NULL,
```

### Step 8: Update user_coupons table in sql/promotion/schema.sql

Change line 130 from:
```sql
`order_id` VARCHAR(64) DEFAULT '' COMMENT '订单ID',
```
To:
```sql
`order_id` BIGINT DEFAULT 0 COMMENT '订单ID',
```

---

## Task 2: Update Go Domain Entities

**Files:**
- Modify: `admin/internal/domain/order/entity.go`
- Modify: `shop/internal/domain/order/entity.go`
- Modify: `admin/internal/domain/fulfillment/entity.go`
- Modify: `pkg/domain/promotion/coupon.go`

### Step 1: Update admin/internal/domain/order/entity.go

**Order struct (line ~51):**
Change:
```go
ID             string
```
To:
```go
ID             int64
```

**OrderItem struct (line ~142):**
Change:
```go
OrderID     string
```
To:
```go
OrderID     int64
```

**Repository interface methods (lines ~173-177):**
Change FindByID signature:
```go
FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id string) (*Order, error)
```
To:
```go
FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Order, error)
```

Change UpdateStatus signature:
```go
UpdateStatus(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id string, status Status) error
```
To:
```go
UpdateStatus(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, status Status) error
```

**Update GORM tags:**
For Order.ID gorm tag change `size:64` to nothing (BIGINT is default):
```go
ID             int64           `gorm:"column:id;primaryKey"`
```

For OrderItem.OrderID gorm tag change `size:64`:
```go
OrderID     int64           `gorm:"column:order_id;not null;index:idx_order_id"`
```

### Step 2: Update shop/internal/domain/order/entity.go

**Order struct (line ~146):**
Change:
```go
ID                  string          `gorm:"column:id;primaryKey;size:64"`
```
To:
```go
ID                  int64           `gorm:"column:id;primaryKey"`
```

**OrderItem struct (line ~296):**
Change:
```go
OrderID     string       `gorm:"column:order_id;not null;size:64;index:idx_order_id"`
```
To:
```go
OrderID     int64        `gorm:"column:order_id;not null;index:idx_order_id"`
```

**Repository interface methods:**
Change FindByID signature from string to int64.

### Step 3: Update admin/internal/domain/fulfillment/entity.go

**Shipment struct (line ~288):**
Change:
```go
OrderID          string          `gorm:"column:order_id;not null;index"`
```
To:
```go
OrderID          int64           `gorm:"column:order_id;not null;index"`
```

**Refund struct (line ~492):**
Change:
```go
OrderID      string         `gorm:"column:order_id;not null;index"`
```
To:
```go
OrderID      int64         `gorm:"column:order_id;not null;index"`
```

**Query structs:**
Change OrderID field type in ShipmentQuery, RefundQuery from string to int64.

### Step 4: Update pkg/domain/promotion/coupon.go

**UserCoupon struct (line ~142):**
Change:
```go
OrderID    string           `json:"order_id"`
```
To:
```go
OrderID    int64            `json:"order_id"`
```

**PromotionUsage struct (line ~170):**
Change:
```go
OrderID        string           `json:"order_id"`
```
To:
```go
OrderID        int64            `json:"order_id"`
```

---

## Task 3: Update Admin API Definitions

**Files:**
- Modify: `admin/desc/fulfillment.api`

### Step 1: Update fulfillment.api - Change order_id types from string to int64

**CreateShipmentReq (line ~18):**
Change:
```go
OrderID      string            `json:"order_id"`
```
To:
```go
OrderID      int64             `json:"order_id"`
```

**BatchShipmentItemReq (line ~35):**
Change:
```go
OrderID    string `json:"order_id"`
```
To:
```go
OrderID    int64  `json:"order_id"`
```

**BatchShipmentResultResp (line ~53):**
Change:
```go
OrderID    string `json:"order_id"`
```
To:
```go
OrderID    int64  `json:"order_id"`
```

**ShipmentDetailResp (line ~91):**
Change:
```go
OrderID       string              `json:"order_id"`
```
To:
```go
OrderID       int64               `json:"order_id"`
```

**ListShipmentsReq (line ~118):**
Change:
```go
OrderID           string `form:"order_id,optional"`
```
To:
```go
OrderID           int64  `form:"order_id,optional"`
```

**GetOrderShipmentsReq (line ~135):**
Change:
```go
OrderID string `path:"order_id"`
```
To:
```go
OrderID int64  `path:"order_id"`
```

**RefundDetailResp (line ~214):**
Change:
```go
OrderID        string   `json:"order_id"`
```
To:
```go
OrderID        int64    `json:"order_id"`
```

**ListRefundsReq (line ~249):**
Change:
```go
OrderID    string `form:"order_id,optional"`
```
To:
```go
OrderID    int64  `form:"order_id,optional"`
```

**OrderFulfillmentDetailResp (line ~417):**
Change:
```go
OrderID           string                      `json:"order_id"`
```
To:
```go
OrderID           int64                       `json:"order_id"`
```

**ShipOrderReq (line ~465):**
Change path param from ID int64 to proper path param:
```go
ID           int64             `path:"id"`
```
(Already int64, no change needed)

**Route paths (lines ~606, ~664, ~668, etc.):**
Routes like `get /api/v1/orders/:order_id/shipments` use `:order_id` as path param - these should be changed to `:id` to match the standard pattern since we're using int64 now.

Update:
- `get /api/v1/orders/:order_id/shipments` → `get /api/v1/orders/:id/shipments`
- Any other routes using `:order_id` as string path param

### Step 2: Regenerate API code

Run:
```bash
cd admin && make api
```

This regenerates `admin/internal/types/types.go` with updated type definitions.

---

## Task 4: Update Go Repository Implementations

**Files:**
- Modify: `admin/internal/infrastructure/persistence/order_repository.go`
- Modify: `admin/internal/infrastructure/persistence/shipment_repository.go`
- Modify: `admin/internal/infrastructure/persistence/refund_repository.go`
- Modify: `admin/internal/infrastructure/persistence/review_repository.go`
- Modify: `admin/internal/infrastructure/persistence/promotion_usage_repository.go`
- Modify: `admin/internal/infrastructure/persistence/user_coupon_repository.go`
- Modify: `shop/internal/infrastructure/persistence/order_repository.go`

### Step 1: Update admin/internal/infrastructure/persistence/order_repository.go

Change all method signatures and implementations that use `order_id` as string to use `int64`.

### Step 2: Update admin/internal/infrastructure/persistence/shipment_repository.go

Update `FindByOrderID` method signature and any other methods using `order_id`.

### Step 3: Update remaining repository files

Update refund_repository.go, review_repository.go, promotion_usage_repository.go, user_coupon_repository.go similarly.

### Step 4: Update shop/internal/infrastructure/persistence/order_repository.go

Update shop service order repository similarly.

---

## Task 5: Update Go Application/Logic Layer

**Files:**
- Modify: Files that call repository methods with order_id parameters
- Common locations: `admin/internal/application/`, `admin/internal/logic/`, `shop/internal/application/`

### Step 1: Update admin/application/fulfillment modules

Check and update any files that pass order_id as string to repository calls.

### Step 2: Update admin/application/payment service

Check `admin/internal/application/payment/service.go` and related files.

### Step 3: Update admin/application/review service

Check review service files.

### Step 4: Update admin/application/promotion coupon_app

Check `admin/internal/application/promotion/coupon_app.go` for UserCoupon operations.

---

## Task 6: Update Frontend TypeScript Code

**Files:**
- Modify: `shop-admin/src/api/order.ts`
- Modify: `shop-admin/src/api/fulfillment.ts`
- Modify: `shop-admin/src/views/orders/**` (Vue components)

### Step 1: Update shop-admin/src/api/order.ts

**Order interface (line ~68):**
Change:
```typescript
order_id: string
```
To:
```typescript
order_id: number
```

**API function parameters:**
Update `getOrderDetail`, `shipOrder`, `updateOrderRemark`, `adjustOrderPrice`, `cancelOrder`, `remindPayment` functions to use `orderId: number` instead of `orderId: string`.

### Step 2: Update shop-admin/src/api/fulfillment.ts

Update all OrderID related types from string to number.

### Step 3: Update Vue components

Search for components using order_id as string and update type annotations.

---

## Task 7: Build and Verify

### Step 1: Build admin service

Run:
```bash
cd admin && make build
```

Fix any compilation errors.

### Step 2: Build shop service

Run:
```bash
cd shop && make build
```

Fix any compilation errors.

### Step 3: Build frontend

Run:
```bash
cd shop-admin && npm run build
```

Fix any TypeScript errors.

### Step 4: Run tests

Run backend tests:
```bash
cd admin && go test ./...
cd shop && go test ./...
```

---

## Verification Checklist

- [ ] SQL schema files updated with BIGINT for order_id fields
- [ ] Test data INSERT statements updated to use numeric IDs
- [ ] Go domain entities updated (Order.ID, OrderItem.OrderID, Shipment.OrderID, Refund.OrderID, etc.)
- [ ] Go repository implementations updated
- [ ] Go application/logic layer updated
- [ ] API definition file updated (fulfillment.api)
- [ ] API code regenerated (`make api`)
- [ ] Frontend TypeScript types updated
- [ ] Frontend components using order_id updated
- [ ] All builds pass (admin, shop, shop-admin)
- [ ] All tests pass
