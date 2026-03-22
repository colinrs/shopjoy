# Promotion Module Design Specification

## Document Information

| Item | Value |
|------|-------|
| Document Title | Promotion Module Design |
| Version | 1.0.0 |
| Status | Draft |
| Created | 2026-03-22 |
| Author | Claude |

---

## Executive Summary

This document specifies the technical design for implementing the Promotion module in ShopJoy, a multi-tenant e-commerce SaaS platform. The module enables merchants to create, manage, and apply promotional rules to drive sales and customer engagement.

### Key Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Scope | Full MVP (7 priorities) | Complete feature set for merchant value |
| Entity Strategy | Extend existing entities | Backward compatibility, incremental migration |
| Calculation Location | Shared domain package (`pkg/domain/promotion`) | Direct access from both admin and shop services |
| Scope Filtering | Full support (storewide, products, categories, brands) | Merchant flexibility |
| Database Migration | SQL migration files | Version control, rollback capability |

---

## Architecture Overview

### System Context

```
┌─────────────────────────────────────────────────────────────────┐
│                        ShopJoy Platform                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐     ┌─────────────────┐     ┌─────────────┐  │
│  │  Admin API  │────▶│ pkg/domain/     │◀────│  Shop API   │  │
│  │  (CRUD)     │     │ promotion       │     │  (Checkout) │  │
│  └─────────────┘     └─────────────────┘     └─────────────┘  │
│         │                    │                    │           │
│         │                    ▼                    │           │
│         │             ┌─────────────┐             │           │
│         │             │   MySQL     │             │           │
│         └────────────▶│  (Shared)   │◀────────────┘           │
│                       └─────────────┘                         │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Layered Architecture

The module follows Domain-Driven Design principles with clear layer separation:

```
pkg/domain/promotion/          ← Shared Domain Layer
    ├── entity.go              ← Entities, Value Objects
    ├── service.go             ← Domain Services (Calculation)
    ├── repository.go          ← Repository Interfaces
    └── scope.go               ← Scope Matching Logic

admin/internal/
├── infra/persistence/         ← Infrastructure Layer
│   └── *_repo_impl.go         ← Repository Implementations
├── application/promotion/     ← Application Layer
│   └── *_app.go               ← Application Services
└── desc/promotion.api         ← API Definition

admin/internal/logic/          ← Handler Layer (generated)
admin/internal/handler/        ← Routes (generated)
```

---

## Domain Model

### Entities

#### Promotion (Aggregate Root)

```go
type Promotion struct {
    ID          int64
    TenantID    shared.TenantID
    Name        string
    Description string
    Type        PromotionType     // FULL_REDUCE, FIXED_DISCOUNT, PERCENT_DISCOUNT
    Status      PromotionStatus   // DRAFT, ACTIVE, INACTIVE, ENDED
    Priority    int               // Higher value = higher priority
    StartAt     time.Time         // UTC
    EndAt       time.Time         // UTC
    Scope       PromotionScope    // Embedded value object
    Currency    string            // ISO 4217: CNY, USD, EUR
    Rules       []PromotionRule   // Child entities
    Audit       shared.AuditInfo  // Embedded
    DeletedAt   *time.Time        // Soft delete
}

type PromotionType string

const (
    PromotionTypeFullReduce     PromotionType = "FULL_REDUCE"
    PromotionTypeFixedDiscount  PromotionType = "FIXED_DISCOUNT"
    PromotionTypePercentDiscount PromotionType = "PERCENT_DISCOUNT"
)

type PromotionStatus int

const (
    PromotionStatusDraft    PromotionStatus = iota
    PromotionStatusActive
    PromotionStatusInactive
    PromotionStatusEnded
)
```

#### PromotionScope (Value Object)

```go
type PromotionScope struct {
    Type       ScopeType  // STOREWIDE, PRODUCTS, CATEGORIES, BRANDS
    IDs        []int64    // Product/Category/Brand IDs (empty if storewide)
    ExcludeIDs []int64    // Product IDs to exclude from promotion
}

type ScopeType string

const (
    ScopeTypeStorewide  ScopeType = "STOREWIDE"
    ScopeTypeProducts   ScopeType = "PRODUCTS"
    ScopeTypeCategories ScopeType = "CATEGORIES"
    ScopeTypeBrands     ScopeType = "BRANDS"
)
```

#### PromotionRule (Entity)

```go
type PromotionRule struct {
    ID             int64
    PromotionID    int64
    ConditionType  ConditionType   // MIN_AMOUNT, MIN_QUANTITY
    ConditionValue int64           // Threshold value (cents)
    ActionType     ActionType      // FIXED_AMOUNT, PERCENTAGE
    ActionValue    int64           // Discount value (cents or basis points)
    MaxDiscount    int64           // Cap for percentage discounts (cents)
    Currency       string          // ISO 4217
    SortOrder      int             // Order for tiered rules
}

type ConditionType string

const (
    ConditionMinAmount   ConditionType = "MIN_AMOUNT"
    ConditionMinQuantity ConditionType = "MIN_QUANTITY"
)

type ActionType string

const (
    ActionFixedAmount ActionType = "FIXED_AMOUNT"
    ActionPercentage  ActionType = "PERCENTAGE"
)
```

#### Coupon (Aggregate Root)

```go
type Coupon struct {
    ID           int64
    TenantID     shared.TenantID
    Name         string
    Code         string            // Unique coupon code
    Description  string
    Type         CouponType        // FIXED_AMOUNT, PERCENTAGE
    Value        int64             // Discount value (cents or basis points)
    MinAmount    int64             // Minimum spend (cents)
    MaxDiscount  int64             // Cap for percentage (cents)
    Currency     string            // ISO 4217
    TotalCount   int               // Total uses allowed (0 = unlimited)
    UsedCount    int               // Current usage count
    PerUserLimit int               // Uses per user (0 = unlimited)
    Status       CouponStatus      // INACTIVE, ACTIVE, EXPIRED, DEPLETED
    StartAt      time.Time         // UTC
    EndAt        time.Time         // UTC
    Scope        PromotionScope    // Same scope support as promotions
    Audit        shared.AuditInfo
    DeletedAt    *time.Time        // Soft delete
}

type CouponType string

const (
    CouponTypeFixedAmount CouponType = "FIXED_AMOUNT"
    CouponTypePercentage  CouponType = "PERCENTAGE"
)

type CouponStatus int

const (
    CouponStatusInactive CouponStatus = iota
    CouponStatusActive
    CouponStatusExpired
    CouponStatusDepleted
)
```

#### UserCoupon (Entity)

```go
type UserCoupon struct {
    ID         int64
    TenantID   shared.TenantID
    UserID     int64
    CouponID   int64
    Status     UserCouponStatus  // UNUSED, USED, EXPIRED
    UsedAt     *time.Time        // When coupon was used
    OrderID    string            // Order ID when used
    ReceivedAt time.Time         // When issued to user
    ExpireAt   time.Time         // Expiration timestamp
}

type UserCouponStatus int

const (
    UserCouponStatusUnused UserCouponStatus = iota
    UserCouponStatusUsed
    UserCouponStatusExpired
)
```

#### PromotionUsage (Entity)

```go
type PromotionUsage struct {
    ID             int64
    TenantID       shared.TenantID
    PromotionID    int64
    RuleID         *int64           // Rule that triggered discount
    OrderID        string           // Order reference
    UserID         int64
    DiscountAmount int64            // Actual discount applied (cents)
    Currency       string
    OriginalAmount int64            // Order total before discount
    FinalAmount    int64            // Order total after discount
    CouponID       *int64           // Coupon used (if any)
    CreatedAt      time.Time
}
```

### Domain Services

#### CalculationService

```go
// CalculationService handles discount calculation logic
type CalculationService struct {
    promotionRepo PromotionRepository
    couponRepo    CouponRepository
    userCouponRepo UserCouponRepository
}

// CalculateRequest - input for calculation
type CalculateRequest struct {
    TenantID   shared.TenantID
    UserID     int64
    CartItems  []CartItem
    Currency   string
    CouponCode string  // Optional
}

// CartItem - simplified cart item for calculation
type CartItem struct {
    ProductID  int64
    CategoryID int64
    BrandID    int64
    SKU        string
    Quantity   int
    UnitPrice  int64  // Cents
    LineTotal  int64  // Cents
}

// CalculateResult - output from calculation
type CalculateResult struct {
    OriginalTotal     int64
    PromotionDiscount int64
    CouponDiscount    int64
    FinalTotal        int64
    AppliedPromotions []AppliedPromotion
    AppliedCoupon     *AppliedCoupon
}

// AppliedPromotion - promotion applied to order
type AppliedPromotion struct {
    PromotionID    int64
    PromotionName  string
    RuleID         int64
    DiscountType   string
    DiscountAmount int64
}

// AppliedCoupon - coupon applied to order
type AppliedCoupon struct {
    CouponID      int64
    CouponName    string
    Code          string
    DiscountType  string
    DiscountAmount int64
}

func (s *CalculationService) CalculateDiscount(
    ctx context.Context,
    req *CalculateRequest,
) (*CalculateResult, error)
```

---

## Calculation Logic

### Flow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    CalculateDiscount()                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. FETCH ACTIVE PROMOTIONS                                     │
│     - Filter by tenant_id, status=ACTIVE                        │
│     - Filter by time: now BETWEEN start_at AND end_at           │
│     - Filter by currency match                                   │
│                                                                 │
│  2. SCOPE MATCHING                                              │
│     For each cart item:                                         │
│     - STOREWIDE: matches all (check exclude_ids)                │
│     - PRODUCTS: product_id IN scope_ids                         │
│     - CATEGORIES: category_id IN scope_ids                      │
│     - BRANDS: brand_id IN scope_ids                             │
│     - Exclude: skip if product_id IN exclude_ids                │
│                                                                 │
│  3. PRIORITY SORTING                                            │
│     - Sort promotions by priority DESC                          │
│     - Same priority: sort by created_at ASC (lower ID first)    │
│                                                                 │
│  4. PROMOTION DISCOUNT CALCULATION                              │
│     For each promotion (in priority order):                     │
│     - For tiered rules: find highest applicable tier            │
│     - Calculate discount based on action_type                   │
│       - FIXED_AMOUNT: use action_value directly                 │
│       - PERCENTAGE: line_total * (action_value / 10000)         │
│     - Apply max_discount cap if set                             │
│     - Record applied promotion                                  │
│                                                                 │
│  5. COUPON VALIDATION (if coupon_code provided)                 │
│     - Find coupon by code and tenant                            │
│     - Validate: status=ACTIVE, within time range                │
│     - Validate: usage limits (total_count, per_user_limit)      │
│     - Validate: min_amount threshold                            │
│     - Validate: currency match                                  │
│     - Apply scope filtering                                     │
│     - Calculate coupon discount                                 │
│                                                                 │
│  6. FINAL CALCULATION                                           │
│     - total_discount = promotion_discount + coupon_discount     │
│     - Ensure total_discount <= original_total (no negative)     │
│     - final_total = original_total - total_discount             │
│                                                                 │
│  7. RETURN RESULT                                               │
│     - Return breakdown with all applied discounts               │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Conflict Resolution Rules

| Rule | Description |
|------|-------------|
| Priority-based | Higher priority promotion wins for same scope |
| Same priority | Earlier created (lower ID) wins |
| Promotion + Coupon | Both apply (coupon stacks after promotion) |
| Currency mismatch | Skip promotion/coupon with different currency |
| Discount cap | Total discount cannot exceed cart total |

### Discount Calculation Examples

**Example 1: Tiered Full Reduction**

```
Rules:
- Spend $50, save $5
- Spend $100, save $15
- Spend $200, save $40

Cart: $120
Result: $15 discount (tier 2 applies, not tier 1)
Final: $105
```

**Example 2: Percentage with Cap**

```
Rule: 20% off, max $50

Cart: $100 → Discount: $20 → Final: $80
Cart: $400 → Discount: $80 → Capped at $50 → Final: $350
```

**Example 3: Promotion + Coupon Stack**

```
Promotion: Spend $100, save $10
Coupon: $5 off (min spend $50)

Cart: $150
Promotion discount: $10
Coupon discount: $5
Final: $135
```

---

## Repository Interfaces

```go
// PromotionRepository
type PromotionRepository interface {
    Create(ctx context.Context, db *gorm.DB, promotion *Promotion) error
    Update(ctx context.Context, db *gorm.DB, promotion *Promotion) error
    Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
    FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Promotion, error)
    FindActive(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, currency string) ([]*Promotion, error)
    FindList(ctx context.Context, db *gorm.DB, query PromotionQuery) ([]*Promotion, int64, error)

    // Rule management
    CreateRules(ctx context.Context, db *gorm.DB, rules []PromotionRule) error
    FindRulesByPromotionID(ctx context.Context, db *gorm.DB, promotionID int64) ([]PromotionRule, error)
    UpdateRule(ctx context.Context, db *gorm.DB, rule *PromotionRule) error
    DeleteRule(ctx context.Context, db *gorm.DB, id int64) error
}

// PromotionQuery
type PromotionQuery struct {
    shared.PageQuery
    TenantID shared.TenantID
    Name     string
    Status   *PromotionStatus
    Type     *PromotionType
}

// CouponRepository
type CouponRepository interface {
    Create(ctx context.Context, db *gorm.DB, coupon *Coupon) error
    Update(ctx context.Context, db *gorm.DB, coupon *Coupon) error
    Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
    FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Coupon, error)
    FindByCode(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, code string) (*Coupon, error)
    FindList(ctx context.Context, db *gorm.DB, query CouponQuery) ([]*Coupon, int64, error)
    IncrementUsage(ctx context.Context, db *gorm.DB, id int64) error
}

// CouponQuery
type CouponQuery struct {
    shared.PageQuery
    TenantID shared.TenantID
    Name     string
    Code     string
    Status   *CouponStatus
    Type     *CouponType
}

// UserCouponRepository
type UserCouponRepository interface {
    Create(ctx context.Context, db *gorm.DB, userCoupon *UserCoupon) error
    FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*UserCoupon, error)
    FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, status *UserCouponStatus) ([]*UserCoupon, error)
    FindByUserAndCoupon(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, couponID int64) ([]*UserCoupon, error)
    MarkUsed(ctx context.Context, db *gorm.DB, id int64, orderID string) error
}

// PromotionUsageRepository
type PromotionUsageRepository interface {
    Create(ctx context.Context, db *gorm.DB, usage *PromotionUsage) error
    FindByOrderID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderID string) (*PromotionUsage, error)
    FindByPromotionID(ctx context.Context, db *gorm.DB, promotionID int64, query shared.PageQuery) ([]*PromotionUsage, int64, error)
    FindByCouponID(ctx context.Context, db *gorm.DB, couponID int64, query shared.PageQuery) ([]*PromotionUsage, int64, error)
}
```

---

## Database Schema

### promotions Table

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| tenant_id | BIGINT | NO | - | Tenant ID |
| name | VARCHAR(200) | NO | - | Promotion name |
| description | TEXT | YES | NULL | Description |
| type | VARCHAR(32) | NO | - | FULL_REDUCE, FIXED_DISCOUNT, PERCENT_DISCOUNT |
| status | TINYINT | NO | 0 | 0=draft, 1=active, 2=inactive, 3=ended |
| priority | INT | NO | 0 | Higher = more priority |
| start_at | BIGINT | NO | - | Start timestamp (UTC Unix) |
| end_at | BIGINT | NO | - | End timestamp (UTC Unix) |
| currency | VARCHAR(10) | NO | 'CNY' | ISO 4217 currency code |
| scope_type | VARCHAR(32) | NO | 'STOREWIDE' | STOREWIDE, PRODUCTS, CATEGORIES, BRANDS |
| scope_ids | JSON | YES | NULL | Array of scope IDs |
| exclude_ids | JSON | YES | NULL | Array of excluded product IDs |
| created_at | BIGINT | NO | - | Creation timestamp |
| updated_at | BIGINT | NO | - | Update timestamp |
| created_by | BIGINT | NO | 0 | Creator user ID |
| updated_by | BIGINT | NO | 0 | Updater user ID |
| deleted_at | BIGINT | YES | NULL | Soft delete timestamp |

**Indexes:**
- PRIMARY KEY (`id`)
- INDEX `idx_tenant_id` (`tenant_id`)
- INDEX `idx_status` (`status`)
- INDEX `idx_tenant_status_time` (`tenant_id`, `status`, `start_at`, `end_at`)
- INDEX `idx_deleted_at` (`deleted_at`)

### promotion_rules Table

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| promotion_id | BIGINT | NO | - | Parent promotion ID |
| condition_type | VARCHAR(32) | NO | - | MIN_AMOUNT, MIN_QUANTITY |
| condition_value | BIGINT | NO | 0 | Threshold value (cents) |
| action_type | VARCHAR(32) | NO | - | FIXED_AMOUNT, PERCENTAGE |
| action_value | BIGINT | NO | 0 | Discount value |
| max_discount | BIGINT | NO | 0 | Maximum discount cap (cents) |
| currency | VARCHAR(10) | NO | 'CNY' | Currency code |
| sort_order | INT | NO | 0 | Sort order for tiers |
| created_at | BIGINT | NO | - | Creation timestamp |
| updated_at | BIGINT | NO | - | Update timestamp |

**Indexes:**
- PRIMARY KEY (`id`)
- INDEX `idx_promotion_id` (`promotion_id`)

### coupons Table

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| tenant_id | BIGINT | NO | - | Tenant ID |
| name | VARCHAR(200) | NO | - | Coupon name |
| code | VARCHAR(50) | NO | - | Unique coupon code |
| description | TEXT | YES | NULL | Description |
| type | VARCHAR(32) | NO | - | FIXED_AMOUNT, PERCENTAGE |
| value | BIGINT | NO | 0 | Discount value |
| min_amount | BIGINT | NO | 0 | Minimum spend (cents) |
| max_discount | BIGINT | NO | 0 | Maximum discount cap |
| currency | VARCHAR(10) | NO | 'CNY' | Currency code |
| total_count | INT | NO | 0 | Total uses allowed (0=unlimited) |
| used_count | INT | NO | 0 | Current usage count |
| per_user_limit | INT | NO | 0 | Per user limit (0=unlimited) |
| status | TINYINT | NO | 0 | 0=inactive, 1=active, 2=expired, 3=depleted |
| start_at | BIGINT | NO | - | Start timestamp |
| end_at | BIGINT | NO | - | End timestamp |
| scope_type | VARCHAR(32) | NO | 'STOREWIDE' | Scope type |
| scope_ids | JSON | YES | NULL | Scope IDs |
| created_at | BIGINT | NO | - | Creation timestamp |
| updated_at | BIGINT | NO | - | Update timestamp |
| created_by | BIGINT | NO | 0 | Creator ID |
| updated_by | BIGINT | NO | 0 | Updater ID |
| deleted_at | BIGINT | YES | NULL | Soft delete timestamp |

**Indexes:**
- PRIMARY KEY (`id`)
- UNIQUE INDEX `uk_tenant_code` (`tenant_id`, `code`)
- INDEX `idx_tenant_status_time` (`tenant_id`, `status`, `start_at`, `end_at`)

### user_coupons Table

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| tenant_id | BIGINT | NO | - | Tenant ID |
| user_id | BIGINT | NO | - | User ID |
| coupon_id | BIGINT | NO | - | Coupon ID |
| status | TINYINT | NO | 0 | 0=unused, 1=used, 2=expired |
| used_at | BIGINT | YES | NULL | Usage timestamp |
| order_id | VARCHAR(50) | NO | '' | Order ID when used |
| received_at | BIGINT | NO | - | Issue timestamp |
| expire_at | BIGINT | NO | - | Expiration timestamp |
| created_at | BIGINT | NO | - | Creation timestamp |
| updated_at | BIGINT | NO | - | Update timestamp |

**Indexes:**
- PRIMARY KEY (`id`)
- INDEX `idx_user_id` (`user_id`)
- INDEX `idx_coupon_id` (`coupon_id`)
- INDEX `idx_user_status` (`user_id`, `status`)

### promotion_usage Table (NEW)

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | Primary key |
| tenant_id | BIGINT | NO | - | Tenant ID |
| promotion_id | BIGINT | NO | - | Promotion ID |
| rule_id | BIGINT | YES | NULL | Rule ID applied |
| order_id | VARCHAR(50) | NO | - | Order ID |
| user_id | BIGINT | NO | - | User ID |
| discount_amount | BIGINT | NO | 0 | Discount applied (cents) |
| currency | VARCHAR(10) | NO | 'CNY' | Currency code |
| original_amount | BIGINT | NO | 0 | Original order total |
| final_amount | BIGINT | NO | 0 | Final order total |
| coupon_id | BIGINT | YES | NULL | Coupon ID if used |
| created_at | BIGINT | NO | - | Creation timestamp |

**Indexes:**
- PRIMARY KEY (`id`)
- INDEX `idx_tenant_id` (`tenant_id`)
- INDEX `idx_promotion_id` (`promotion_id`)
- INDEX `idx_order_id` (`order_id`)
- INDEX `idx_user_id` (`user_id`)
- INDEX `idx_created_at` (`created_at`)

---

## API Endpoints

### Promotion Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/promotions | Create promotion |
| PUT | /api/v1/promotions/:id | Update promotion |
| GET | /api/v1/promotions/:id | Get promotion details |
| GET | /api/v1/promotions | List promotions |
| DELETE | /api/v1/promotions/:id | Delete promotion (soft) |
| POST | /api/v1/promotions/:id/activate | Activate promotion |
| POST | /api/v1/promotions/:id/deactivate | Deactivate promotion |

### Promotion Rules

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/promotions/:id/rules | Get promotion rules |
| POST | /api/v1/promotions/:id/rules | Add rule(s) |
| PUT | /api/v1/promotion-rules/:id | Update rule |
| DELETE | /api/v1/promotion-rules/:id | Delete rule |

### Coupon Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/coupons | Create coupon |
| PUT | /api/v1/coupons/:id | Update coupon |
| GET | /api/v1/coupons/:id | Get coupon details |
| GET | /api/v1/coupons | List coupons |
| DELETE | /api/v1/coupons/:id | Delete coupon (soft) |
| POST | /api/v1/coupons/generate | Generate batch codes |
| GET | /api/v1/coupons/:id/usage | Get usage history |

### User Coupons

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/user-coupons | Issue coupon to user |
| GET | /api/v1/user-coupons | List user's coupons |

### Request/Response Types

See `admin/desc/promotion.api` for complete type definitions.

---

## File Structure

### New Files

```
pkg/domain/promotion/
├── entity.go           # Promotion, PromotionRule, PromotionScope
├── coupon.go           # Coupon, UserCoupon, PromotionUsage
├── repository.go       # Repository interfaces
├── service.go          # CalculationService
├── scope.go            # Scope matching functions
└── errors.go           # Domain errors

admin/internal/infra/persistence/
├── promotion_repo_impl.go
├── coupon_repo_impl.go
├── user_coupon_repo_impl.go
└── promotion_usage_repo_impl.go

admin/internal/application/promotion/
├── promotion_app.go
└── coupon_app.go

admin/desc/
└── promotion.api

migrations/
├── 2026032201_alter_promotions_add_scope.sql
├── 2026032202_alter_promotion_rules_add_fields.sql
├── 2026032203_alter_coupons_add_scope.sql
└── 2026032204_create_promotion_usage.sql
```

### Modified/Deprecated Files

```
admin/internal/domain/promotion/entity.go  → MOVE to pkg/domain/promotion/
admin/internal/domain/coupon/entity.go     → MOVE to pkg/domain/promotion/
```

---

## Implementation Phases

### Phase 1: Foundation (Week 1-2)

1. Create `pkg/domain/promotion/` package with entities
2. Write and execute database migrations
3. Implement repository interfaces and GORM implementations
4. Implement `CalculationService` core logic
5. Implement scope matching functions

### Phase 2: Admin APIs (Week 3-4)

1. Define `admin/desc/promotion.api`
2. Run `make api` to generate handlers
3. Implement application services (`promotion_app.go`, `coupon_app.go`)
4. Implement logic handlers
5. Unit tests for calculation service

### Phase 3: Integration (Week 5-6)

1. Integrate calculation service with shop service
2. Implement promotion usage recording
3. End-to-end testing
4. Performance testing

### Phase 4: Polish (Week 7)

1. Error handling refinement
2. Logging and monitoring
3. Documentation
4. Production deployment

---

## Testing Strategy

### Unit Tests

- `service_test.go`: Calculation logic tests
- `scope_test.go`: Scope matching tests
- `entity_test.go`: Entity method tests

### Integration Tests

- Repository implementation tests
- API endpoint tests

### Test Cases

| Scenario | Input | Expected Output |
|----------|-------|-----------------|
| Tiered discount - tier 1 | Cart $60, tiers at $50/$100 | $5 discount |
| Tiered discount - tier 2 | Cart $120, tiers at $50/$100 | $15 discount |
| Percentage with cap | Cart $400, 20% max $50 | $50 discount |
| Priority conflict | 2 promotions, priority 10 vs 5 | Higher priority wins |
| Coupon + Promotion | Valid coupon + active promotion | Both apply |
| Expired promotion | End time passed | Not applied |
| Currency mismatch | Promotion USD, cart EUR | Not applied |
| Scope - products | Scope products [1,2], cart has 3 | Only 1,2 discounted |
| Scope - exclude | Storewide, exclude [5] | All except 5 discounted |
| Discount cap | Cart $10, coupon $20 | $10 discount (capped) |

---

## Error Handling

### Domain Errors

```go
var (
    ErrPromotionNotFound      = errors.New("promotion not found")
    ErrPromotionExpired       = errors.New("promotion has expired")
    ErrPromotionNotActive     = errors.New("promotion is not active")
    ErrPromotionCurrencyMismatch = errors.New("promotion currency mismatch")

    ErrCouponNotFound         = errors.New("coupon not found")
    ErrCouponExpired          = errors.New("coupon has expired")
    ErrCouponDepleted         = errors.New("coupon usage limit reached")
    ErrCouponUserLimitReached = errors.New("user coupon limit reached")
    ErrCouponMinAmountNotMet  = errors.New("minimum spend not met")
    ErrCouponCurrencyMismatch = errors.New("coupon currency mismatch")

    ErrInvalidScopeType       = errors.New("invalid scope type")
    ErrInvalidPromotionType   = errors.New("invalid promotion type")
    ErrInvalidCouponType      = errors.New("invalid coupon type")
)
```

---

## Performance Considerations

1. **Active Promotion Caching**: Cache active promotions per tenant in Redis (5-minute TTL)
2. **Scope Lookup Optimization**: Pre-load product→category→brand mappings
3. **Database Indexing**: Proper indexes on frequently queried columns
4. **Pagination**: All list endpoints support pagination (max 100 per page)

---

## Security Considerations

1. **Tenant Isolation**: All queries filtered by `tenant_id`
2. **Soft Delete**: Deleted records preserved with `deleted_at` timestamp
3. **Rate Limiting**: Consider rate limiting coupon validation to prevent brute-force
4. **Input Validation**: All API inputs validated before processing

---

## References

- PRD: `docs/prd/promotion-prd.md`
- DDD Patterns: `docs/ARCHITECTURE.md`
- Shared Domain: `pkg/domain/shared/`
- Database Rules: `.claude/rules/common/db.md`
- Price Rules: `.claude/rules/golang/price.md`