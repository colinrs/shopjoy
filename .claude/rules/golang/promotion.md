# Promotion Domain Rules

Rules specific to the promotion module (`admin/internal/logic/promotions/`, `admin/internal/application/promotion/`, `admin/internal/infrastructure/persistence/promotion_repository.go`, `pkg/domain/promotion/`).

## MUST

| # | Rule | Rationale |
|---|------|-----------|
| 1 | Every field on `UpdatePromotionReq` / `CreatePromotionReq` must have a corresponding assignment in the logic layer AND the application layer | Silent field loss is the #1 bug class in this module. The wire types carry all fields; the logic/app layers must pass them through. |
| 2 | Discount fields (`discount_type`, `discount_value`, `min_order_amount`, `max_discount`) are stored in `promotion_rules`, not on `promotions` | The `promotions` table has no discount columns. The logic layer must build a `CreatePromotionRuleRequest` from these wire fields and pass it in `CreatePromotionRequest.Rules` / `UpdatePromotionRequest.Rules`. |
| 3 | `convertPromotionToDetailResp()` must map ALL response fields from the entity/rules | This function is the single point of truth for the GET/LIST response. If a field exists on `PromotionDetailResp`, it must be populated here. |
| 4 | `promotionModel` column names must match the actual DB column names | The `gorm:"column:..."` tag must match `SHOW CREATE TABLE`. Example: `max_discount_amount` not `max_discount`. |
| 5 | Use `*Status` / `*Type` pointers for filter parameters in `QueryPromotionRequest` / `QueryCouponRequest` | The iota-zero values (`StatusPending=0`, `TypeDiscount=0`) collide with "no filter" if using value types with `!= 0` sentinel. Use nil pointer = no filter. |
| 6 | The wire status `"expired"` is derived from `EndAt`, not stored | Never map `"expired"` to a `Status` enum value. Use `ExpiredOnly bool` on the query and filter with `end_at <= NOW()` in the repository. |
| 7 | `buildPromotionScope()` / `buildCouponScope()` must normalize lowercase wire values to uppercase domain constants | Wire sends `"products"` but domain constant is `ScopeTypeProducts = "PRODUCTS"`. Use `strings.ToUpper()` before comparing. |

## FORBIDDEN

| # | Rule | Consequence |
|---|------|-------------|
| 8 | Using `!= 0` as a sentinel for "filter set" on iota-based enums | Silently drops filters for the zero-value enum member (pending status, discount type) |
| 9 | Adding fields to wire types without wiring them through logic + app + persistence layers | Creates the illusion of functionality while silently dropping data |
| 10 | Using `DESCRIBE` to verify schema matches | `DESCRIBE` does not show column comments. Use `SHOW CREATE TABLE` or `SHOW FULL COLUMNS`. |
| 11 | Leaving stale doc comments that say "no DB columns" after schema migration lands | Misleads future developers into thinking fields are no-ops |

## Checklist

Before committing any promotion/coupon change:

- [ ] Every field on the wire request type has a corresponding assignment in the logic layer
- [ ] Every field on the wire request type has a corresponding assignment in the application layer
- [ ] `convertPromotionToDetailResp()` / `convertCouponToDetailResp()` maps all response fields
- [ ] `promotionModel` / `couponModel` column names match `SHOW CREATE TABLE`
- [ ] Filter parameters use pointer types (`*Status`, `*Type`)
- [ ] `"expired"` status uses `ExpiredOnly bool` + `end_at <= NOW()`, not a status enum
- [ ] Scope type strings are normalized with `strings.ToUpper()` before domain comparison
- [ ] `go build ./...` passes after changes
- [ ] Doc comments are accurate (no "no DB columns" after columns exist)
