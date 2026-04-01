# Orders ID Migration: VARCHAR(64) to BIGINT

## Context

The current `orders` table uses `VARCHAR(64)` for the primary key `id`, storing human-readable strings like `'ORD202503010001'`. This is inconsistent with project conventions (which expect BIGINT for IDs) and creates type mismatches with related tables that reference `order_id`.

## Decision

- `orders.id` → BIGINT AUTO_INCREMENT (primary key)
- `order_no` remains VARCHAR(64) as the human-readable business identifier
- All foreign key `order_id` fields change from VARCHAR(64) to BIGINT

## Scope

### Database Schema Changes

| File | Table | Column | Change |
|------|-------|--------|--------|
| `sql/order/schema.sql` | orders | id | VARCHAR(64) → BIGINT |
| `sql/order/schema.sql` | order_items | order_id | VARCHAR(64) → BIGINT |
| `sql/fulfillment/schema.sql` | shipments | order_id | VARCHAR(64) → BIGINT |
| `sql/fulfillment/schema.sql` | refunds | order_id | VARCHAR(64) → BIGINT |
| `sql/review/schema.sql` | reviews | order_id | VARCHAR(64) → BIGINT |
| `sql/payment/schema.sql` | payments | order_id | VARCHAR(64) → BIGINT |
| `sql/promotion/schema.sql` | user_coupons | order_id | VARCHAR(64) → BIGINT |
| `sql/promotion/schema.sql` | promotion_usage | order_id | VARCHAR(64) → BIGINT |

Test data INSERT statements in `sql/order/schema.sql` must be updated to use numeric IDs.

### Domain Entity Changes

| File | Entity | Field | Type Change |
|------|--------|-------|-------------|
| `admin/internal/domain/order/entity.go` | Order | ID | string → int64 |
| `admin/internal/domain/order/entity.go` | OrderItem | OrderID | string → int64 |
| `admin/internal/domain/order/entity.go` | Repository | FindByID param | string → int64 |
| `shop/internal/domain/order/entity.go` | Order | ID | string → int64 |
| `shop/internal/domain/order/entity.go` | OrderItem | OrderID | string → int64 |
| `shop/internal/domain/order/entity.go` | Repository | FindByID param | string → int64 |
| `admin/internal/domain/fulfillment/entity.go` | Shipment | OrderID | string → int64 |
| `admin/internal/domain/fulfillment/entity.go` | Refund | OrderID | string → int64 |
| `admin/internal/domain/fulfillment/entity.go` | Query | OrderID | string → int64 |
| `pkg/domain/promotion/coupon.go` | UserCoupon | OrderID | string → int64 |
| `pkg/domain/promotion/coupon.go` | PromotionUsage | OrderID | string → int64 |

### API Definition Changes

**`admin/desc/fulfillment.api`** - Change `order_id` from `string` to `int64` in:
- CreateShipmentReq
- BatchShipmentItemReq
- BatchShipmentResultResp
- ShipmentDetailResp
- ListShipmentsReq
- GetOrderShipmentsReq
- RefundDetailResp
- ListRefundsReq
- OrderFulfillmentDetailResp
- All route path params using `:order_id`

After API changes, run `make api` in admin service to regenerate types.

### Frontend TypeScript Changes

| File | Changes |
|------|---------|
| `shop-admin/src/api/order.ts` | Order.order_id: string → number, API function params |
| `shop-admin/src/api/fulfillment.ts` | Order ID types in shipment/refund responses |
| `shop-admin/src/views/orders/**` | Components using order_id |
| `shop-admin/src/views/fulfillment/**` | Components using order_id |

### Repository Implementation Updates

Update all repository implementations referencing order_id as string:
- `admin/internal/infrastructure/persistence/order_repository.go`
- `admin/internal/infrastructure/persistence/shipment_repository.go`
- `admin/internal/infrastructure/persistence/refund_repository.go`
- `admin/internal/infrastructure/persistence/review_repository.go`
- `admin/internal/infrastructure/persistence/payment_repository.go`
- `admin/internal/infrastructure/persistence/promotion_usage_repository.go`
- `admin/internal/infrastructure/persistence/user_coupon_repository.go`
- `shop/internal/infrastructure/persistence/order_repository.go`

Also update application layer services, logic handlers, and any DTOs that carry order_id.

## Consequences

- All code referencing order_id must be updated consistently
- API contracts change (string → number for order_id)
- Frontend must update all order_id type annotations
- Test data must use numeric IDs
- No backward compatibility with existing VARCHAR order IDs

## Verification

After changes:
1. Run `make build` in admin and shop services
2. Run `make build` in shop-admin
3. All tests must pass
