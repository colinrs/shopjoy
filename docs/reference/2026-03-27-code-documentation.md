# ShopJoy Code Documentation

> **Version:** 1.0
> **Last Updated:** 2026-03-27

---

## Overview

This document provides detailed documentation for the core domain entities and shared components in the ShopJoy codebase.

---

## Table of Contents

1. [Shared Domain Components](#shared-domain-components)
2. [User Domain](#user-domain)
3. [Product Domain](#product-domain)
4. [Admin User Domain](#admin-user-domain)
5. [Role Domain](#role-domain)
6. [Tenant Domain](#tenant-domain)
7. [Payment Domain](#payment-domain)
8. [Fulfillment Domain](#fulfillment-domain)
9. [Promotion Domain](#promotion-domain)
10. [Points Domain](#points-domain)
11. [Storefront Domain](#storefront-domain)
12. [Review Domain](#review-domain)

---

## Shared Domain Components

Located in `pkg/domain/shared/`

### TenantID

```go
// TenantID is a value object representing a tenant identifier
type TenantID int64

// Int64 converts TenantID to int64
func (t TenantID) Int64() int64

// String converts TenantID to string
func (t TenantID) String() string

// IsValid checks if TenantID is valid (non-zero)
func (t TenantID) IsValid() bool
```

### Money

```go
// Money represents a monetary value with currency
type Money struct {
    Amount   int64   // Amount in cents/smallest unit
    Currency string  // ISO 4217 currency code
}

// NewMoney creates a new Money value object
func NewMoney(amount int64, currency string) Money

// Add adds two Money values (must have same currency)
func (m Money) Add(other Money) (Money, error)

// Subtract subtracts one Money from another
func (m Money) Subtract(other Money) (Money, error)

// Multiply multiplies Money by a factor
func (m Money) Multiply(factor int64) Money

// Equals checks if two Money values are equal
func (m Money) Equals(other Money) bool

// IsZero checks if amount is zero
func (m Money) IsZero() bool

// IsPositive checks if amount is positive
func (m Money) IsPositive() bool
```

### AuditInfo

```go
// AuditInfo embedded struct for tracking creation and modification
type AuditInfo struct {
    CreatedAt time.Time
    UpdatedAt time.Time
    CreatedBy int64
    UpdatedBy int64
}
```

### UnixTime

```go
// UnixTime wraps time.Time for database storage as Unix timestamp (milliseconds)
type UnixTime struct {
    time.Time
}

// NewUnixTime creates UnixTime from time.Time
func NewUnixTime(t time.Time) UnixTime

// Scan implements sql.Scanner for database retrieval
func (u *UnixTime) Scan(value interface{}) error

// Value implements driver.Valuer for database storage
func (u UnixTime) Value() (driver.Value, error)
```

---

## User Domain

Located in `admin/internal/domain/user/`

### User Entity

```go
// User represents a customer account
type User struct {
    ID        int64
    TenantID  shared.TenantID
    Email     string
    Phone     string
    Password  string        // Hashed password
    Name      string
    Avatar    string
    Gender    Gender        // 0=unknown, 1=male, 2=female, 3=other
    Birthday  *shared.UnixTime
    Status    Status        // 0=inactive, 1=active, 2=suspended
    LastLogin *shared.UnixTime
    DeletedAt *int64        // Soft delete timestamp
    Audit     shared.AuditInfo
}

// TableName returns the database table name
func (u *User) TableName() string

// SetPassword hashes and sets the user's password
func (u *User) SetPassword(plainPassword string) error

// CheckPassword verifies the password against stored hash
func (u *User) CheckPassword(plainPassword string) bool

// CanLogin checks if user can authenticate
func (u *User) CanLogin() bool

// UpdateLastLogin updates the last login timestamp
func (u *User) UpdateLastLogin()

// Suspend suspends the user account
func (u *User) Suspend() error

// Activate activates the user account
func (u *User) Activate() error

// SoftDelete marks the user as deleted
func (u *User) SoftDelete() error

// IsDeleted checks if user is soft deleted
func (u *User) IsDeleted() bool
```

### User Status Constants

```go
const (
    StatusInactive Status = iota  // 0
    StatusActive                   // 1
    StatusSuspended                // 2
    StatusDeleted                  // 3
)
```

### User Repository Interface

```go
type Repository interface {
    Create(ctx context.Context, db *gorm.DB, user *User) error
    Update(ctx context.Context, db *gorm.DB, user *User) error
    Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
    FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*User, error)
    FindByEmail(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, email string) (*User, error)
    FindByPhone(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, phone string) (*User, error)
    FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query Query) ([]*User, int64, error)
    Exists(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, email, phone string) (bool, error)
    GetStats(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (*Stats, error)
}
```

### User Address Entity

```go
// Address represents a user's shipping address
type Address struct {
    ID          int64
    UserID      int64
    TenantID    shared.TenantID
    Name        string
    Phone       string
    Country     string
    Province    string
    City        string
    District    string
    Address     string
    PostalCode  string
    IsDefault   bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

---

## Product Domain

Located in `admin/internal/domain/product/`

### Product Entity

```go
// Product represents a product in the catalog
type Product struct {
    ID              int64
    TenantID        shared.TenantID
    SKU             string
    Name            string
    Description     string
    Price           Money           // Embedded value object
    CostPrice       Money           // Embedded value object
    Stock           int
    Status          Status
    CategoryID      int64
    Brand           string
    SKUPrefix       string
    Tags            []string
    Images          []string
    IsMatrixProduct bool

    // Compliance fields (cross-border)
    HSCode         string
    COO            string          // Country of Origin
    Weight         decimal.Decimal
    WeightUnit     string
    Dimensions     Dimensions      // Embedded value object
    DangerousGoods []string

    DeletedAt *int64
    CreatedAt time.Time
    UpdatedAt time.Time
}

// TableName returns the database table name
func (p *Product) TableName() string

// NewProduct creates a new product (factory method)
func NewProduct(id int64, tenantID shared.TenantID, name, description string, price Money, categoryID int64) (*Product, error)

// NewProductWithCompliance creates a product with compliance info
func NewProductWithCompliance(...) (*Product, error)

// PutOnSale transitions product to on_sale status
func (p *Product) PutOnSale() error

// TakeOffSale transitions product to off_sale status
func (p *Product) TakeOffSale() error

// UpdateStock sets the stock quantity
func (p *Product) UpdateStock(quantity int) error

// DeductStock reduces stock by specified quantity
func (p *Product) DeductStock(quantity int) error

// UpdatePrice updates the product price
func (p *Product) UpdatePrice(newPrice Money) error

// SoftDelete marks product as deleted
func (p *Product) SoftDelete() error

// IsDeleted checks if product is deleted
func (p *Product) IsDeleted() bool

// IsOnSale checks if product is currently on sale
func (p *Product) IsOnSale() bool

// HasComplianceInfo checks if compliance info is set
func (p *Product) HasComplianceInfo() bool

// IsDangerousGoods checks if product is hazardous
func (p *Product) IsDangerousGoods() bool
```

### Product Status

```go
const (
    StatusDraft Status = iota  // 0
    StatusOnSale               // 1
    StatusOffSale              // 2
    StatusDeleted              // 3
)

// String returns the status as string
func (s Status) String() string

// IsValid checks if status value is valid
func (s Status) IsValid() bool

// CanTransitionTo checks valid status transitions
func (s Status) CanTransitionTo(target Status) bool
```

**Valid Status Transitions:**
- `draft` -> `on_sale`, `off_sale`
- `on_sale` -> `off_sale`
- `off_sale` -> `on_sale`

### Product Repository Interface

```go
type Repository interface {
    Create(ctx context.Context, db *gorm.DB, product *Product) error
    Update(ctx context.Context, db *gorm.DB, product *Product) error
    Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
    FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Product, error)
    FindByIDs(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, ids []int64) ([]*Product, error)
    FindList(ctx context.Context, db *gorm.DB, query Query) ([]*Product, int64, error)
    UpdateStock(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, delta int) error
    Exists(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (bool, error)
}
```

### SKU Entity

```go
// SKU represents a product variant (size, color, etc.)
type SKU struct {
    ID             int64
    ProductID      int64
    Code           string
    Price          Money
    Stock          int
    LockedStock    int
    SafetyStock    int
    PreSaleEnabled bool
    Attributes     map[string]string  // {color: "red", size: "L"}
    Status         Status
    CreatedAt      time.Time
    UpdatedAt      time.Time
}
```

### Category Entity

```go
// Category represents a product category (hierarchical)
type Category struct {
    ID             int64
    ParentID       int64
    TenantID       shared.TenantID
    Name           string
    Code           string
    Level          int         // 1-3 max depth
    Sort           int
    Icon           string
    Image          string
    SeoTitle       string
    SeoDescription string
    Status         Status
    CreatedAt      time.Time
    UpdatedAt      time.Time
}
```

### Market Entity

```go
// Market represents a sales region/country
type Market struct {
    ID               int64
    Code             string    // US, UK, DE, AU
    Name             string
    Currency         string    // USD, GBP, EUR, AUD
    DefaultLanguage  string
    Flag             string    // Flag emoji
    IsActive         bool
    IsDefault        bool
    VatRate          string
    GstRate          string
    IossEnabled      bool
    IncludeTax       bool
    CreatedAt        time.Time
    UpdatedAt        time.Time
}
```

### ProductMarket Entity

```go
// ProductMarket links products to markets with pricing
type ProductMarket struct {
    ID                   int64
    ProductID            int64
    MarketID             int64
    IsEnabled            bool
    Price                string     // Decimal as string
    CompareAtPrice       string
    Currency             string
    StockAlertThreshold  int
    PublishedAt          *int64
    CreatedAt           time.Time
    UpdatedAt           time.Time
}
```

### Inventory Entity

```go
// Inventory tracks stock per SKU per warehouse
type Inventory struct {
    ID              int64
    SKUID           int64
    WarehouseID     int64
    AvailableStock  int
    LockedStock     int
    SafetyStock     int
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

---

## Admin User Domain

Located in `admin/internal/domain/adminuser/`

### AdminUser Entity

```go
// AdminUser represents an admin account
type AdminUser struct {
    ID        int64
    TenantID  shared.TenantID  // 0 for platform admins
    Username  string
    Email     string
    Mobile    string
    Password  string          // Hashed
    RealName  string
    Avatar    string
    Type      int             // 1=platform, 2=tenant_admin, 3=sub_account
    Status    int             // 1=normal, 2=disabled
    LastLogin *shared.UnixTime
    Audit     shared.AuditInfo
}

// TableName returns the database table name
func (a *AdminUser) TableName() string

// CanLogin checks if admin can authenticate
func (a *AdminUser) CanLogin() bool

// IsPlatformAdmin checks if user is platform super admin
func (a *AdminUser) IsPlatformAdmin() bool

// IsTenantAdmin checks if user is tenant admin
func (a *AdminUser) IsTenantAdmin() bool
```

### AdminUser Types

```go
const (
    AdminTypePlatform AdminType = iota + 1  // 1
    AdminTypeTenant                        // 2
    AdminTypeSubAccount                    // 3
)

const (
    AdminStatusNormal   AdminStatus = 1
    AdminStatusDisabled AdminStatus = 2
)
```

---

## Role Domain

Located in `admin/internal/domain/role/`

### Role Entity

```go
// Role represents a role with permissions
type Role struct {
    ID          int64
    TenantID    shared.TenantID
    Name        string
    Code        string
    Description string
    Status      Status
    IsSystem    bool      // System roles cannot be modified
    Permissions []*Permission
    Audit       shared.AuditInfo
}

// TableName returns the database table name
func (r *Role) TableName() string
```

### Permission Entity

```go
// Permission represents an individual permission
type Permission struct {
    ID       int64
    Name     string
    Code     string
    Type     PermissionType  // 0=menu, 1=button, 2=api
    ParentID int64
    Path     string
    Icon     string
    Sort     int
}

type PermissionType int8

const (
    PermissionTypeMenu PermissionType = iota  // 0
    PermissionTypeButton                     // 1
    PermissionTypeAPI                        // 2
)
```

---

## Tenant Domain

Located in `admin/internal/domain/tenant/`

### Tenant Entity

```go
// Tenant represents a merchant/organization
type Tenant struct {
    ID          int64
    Name        string
    Code        string
    Domain      string
    Logo        string
    ContactName string
    ContactEmail string
    ContactPhone string
    Status      TenantStatus
    Plan        int        // Subscription plan
    ExpireAt    *int64
    Audit       shared.AuditInfo
}

type TenantStatus int

const (
    TenantStatusActive TenantStatus = iota  // 0
    TenantStatusSuspended                   // 1
    TenantStatusExpired                     // 2
)
```

---

## Payment Domain

Located in `admin/internal/domain/payment/`

### Payment Entity

```go
// Payment represents a payment transaction
type Payment struct {
    ID                   int64
    PaymentNo            string
    TenantID             shared.TenantID
    OrderID              int64
    UserID               int64
    Amount               Money
    Status               PaymentStatus
    Channel              string         // stripe, paypal, etc.
    ChannelTransactionID string
    PaidAt               *shared.UnixTime
    Audit                shared.AuditInfo
}

type PaymentStatus int8

const (
    PaymentStatusPending   PaymentStatus = 0
    PaymentStatusSucceeded PaymentStatus = 1
    PaymentStatusFailed    PaymentStatus = 2
)
```

### Refund Entity

```go
// Refund represents a refund transaction
type Refund struct {
    ID              int64
    RefundNo        string
    PaymentID       int64
    OrderID         int64
    UserID          int64
    Type            RefundType     // 1=full, 2=partial
    Amount          Money
    Status          RefundStatus
    ReasonType      string
    Reason          string
    RejectReason    string
    ApprovedAt      *shared.UnixTime
    ApprovedBy      int64
    CompletedAt     *shared.UnixTime
    ChannelRefundID string
    Audit           shared.AuditInfo
}

type RefundStatus int8

const (
    RefundStatusPending   RefundStatus = 0
    RefundStatusApproved  RefundStatus = 1
    RefundStatusRejected  RefundStatus = 2
    RefundStatusCompleted RefundStatus = 3
    RefundStatusCancelled RefundStatus = 4
)
```

---

## Fulfillment Domain

Located in `admin/internal/domain/fulfillment/`

### Shipment Entity

```go
// Shipment represents a package shipment
type Shipment struct {
    ID            int64
    ShipmentNo    string
    OrderID       int64
    CarrierCode   string
    CarrierName   string
    TrackingNo    string
    TrackingURL   string
    ShippingCost  Money
    Weight        decimal.Decimal
    Status        ShipmentStatus
    ShippedAt     *shared.UnixTime
    DeliveredAt   *shared.UnixTime
    Remark        string
    CreatedBy     int64
    Audit         shared.AuditInfo
    Items         []*ShipmentItem
}

type ShipmentStatus int8

const (
    ShipmentStatusPending   ShipmentStatus = 0
    ShipmentStatusShipped   ShipmentStatus = 1
    ShipmentStatusInTransit ShipmentStatus = 2
    ShipmentStatusDelivered ShipmentStatus = 3
    ShipmentStatusFailed    ShipmentStatus = 4
)
```

### Carrier Entity

```go
// Carrier represents a shipping carrier
type Carrier struct {
    ID          int64
    Code        string
    Name        string
    TrackingURL string
    IsActive    bool
    Sort        int
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

---

## Promotion Domain

Located in `admin/internal/domain/promotion/`

### Promotion Entity

```go
// Promotion represents a marketing promotion
type Promotion struct {
    ID              int64
    TenantID        shared.TenantID
    Name            string
    Description     string
    Type            string        // discount, coupon, flash_sale, bundle
    Status          string        // draft, active, inactive, expired
    DiscountType    string        // percentage, fixed_amount, buy_x_get_y
    DiscountValue   Money
    MinOrderAmount Money
    MaxDiscount    Money
    UsageLimit     int
    UsedCount      int
    PerUserLimit   int
    StartTime      *shared.UnixTime
    EndTime        *shared.UnixTime
    ProductIDs     []int64
    CategoryIDs    []int64
    MarketIDs      []int64
    Tags           []string
    Audit          shared.AuditInfo
}
```

### Coupon Entity

```go
// Coupon represents a coupon code
type Coupon struct {
    ID              int64
    TenantID        shared.TenantID
    Code            string
    Name            string
    Description     string
    Type            string        // fixed_amount, percentage, free_shipping
    DiscountValue   Money
    MinOrderAmount Money
    MaxDiscount    Money
    UsageLimit     int
    UsedCount      int
    PerUserLimit   int
    StartTime      *shared.UnixTime
    EndTime        *shared.UnixTime
    Status         CouponStatus
    ProductIDs     string         // JSON array as string
    CategoryIDs    string
    MarketIDs      string
    Audit          shared.AuditInfo
}

type CouponStatus string

const (
    CouponStatusDraft    CouponStatus = "draft"
    CouponStatusActive   CouponStatus = "active"
    CouponStatusInactive CouponStatus = "inactive"
    CouponStatusExpired  CouponStatus = "expired"
)
```

### UserCoupon Entity

```go
// UserCoupon represents a coupon issued to a user
type UserCoupon struct {
    ID        int64
    UserID    int64
    CouponID  int64
    OrderID   *int64
    Status    UserCouponStatus
    UsedAt    *shared.UnixTime
    Audit     shared.AuditInfo
}

type UserCouponStatus string

const (
    UserCouponStatusAvailable UserCouponStatus = "available"
    UserCouponStatusUsed     UserCouponStatus = "used"
    UserCouponStatusExpired  UserCouponStatus = "expired"
)
```

---

## Points Domain

Located in `admin/internal/domain/points/`

### PointsAccount Entity

```go
// PointsAccount represents a user's points wallet
type PointsAccount struct {
    ID             int64
    UserID         int64
    TenantID       shared.TenantID
    Balance        int64
    FrozenBalance  int64
    TotalEarned    int64
    TotalRedeemed  int64
    TotalExpired   int64
    Audit          shared.AuditInfo
}
```

### EarnRule Entity

```go
// EarnRule defines how points are earned
type EarnRule struct {
    ID               int64
    TenantID         shared.TenantID
    Name             string
    Description      string
    Scenario         string        // ORDER_PAYMENT, SIGN_IN, etc.
    CalculationType  string        // FIXED, RATIO, TIERED
    FixedPoints      int64
    Ratio            string        // Points per currency unit
    Tiers            []*TierConfig
    ConditionType    string        // NONE, NEW_USER, etc.
    ConditionValue   string        // JSON
    ExpirationMonths int
    Status           string        // draft, active, inactive
    Priority         int
    StartAt          *shared.UnixTime
    EndAt            *shared.UnixTime
    Audit            shared.AuditInfo
}

type TierConfig struct {
    Threshold *int64  // nil for unlimited tier
    Ratio     string
}
```

### RedeemRule Entity

```go
// RedeemRule defines how points can be redeemed
type RedeemRule struct {
    ID             int64
    TenantID       shared.TenantID
    Name           string
    Description    string
    CouponID       int64
    CouponName     string
    PointsRequired int64
    TotalStock     int64
    UsedStock      int64
    PerUserLimit   int
    Status         string        // inactive, active
    StartAt        *shared.UnixTime
    EndAt          *shared.UnixTime
    Audit          shared.AuditInfo
}
```

### PointsTransaction Entity

```go
// PointsTransaction records points changes
type PointsTransaction struct {
    ID            int64
    UserID        int64
    AccountID     int64
    Points        int64          // Positive=earn, negative=deduct
    BalanceAfter  int64
    Type          string         // EARN, REDEEM, ADJUST, EXPIRE, FREEZE, UNFREEZE
    ReferenceType string         // order, review, manual, etc.
    ReferenceID   string
    Description   string
    ExpiresAt     *shared.UnixTime
    CreatedAt     time.Time
}
```

---

## Storefront Domain

Located in `admin/internal/domain/storefront/`

### Theme Entity

```go
// Theme represents a storefront theme
type Theme struct {
    ID            int64
    Code          string
    Name          string
    Description   string
    PreviewImage  string
    Thumbnail     string
    IsPreset      bool
    IsCurrent     bool
    DefaultConfig string     // JSON
    ConfigSchema  string     // JSON schema
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

### Page Entity

```go
// Page represents a storefront page
type Page struct {
    ID          int64
    TenantID    shared.TenantID
    ThemeID     int64
    PageType    string
    Name        string
    Slug        string
    IsPublished bool
    Version     int
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### Decoration Entity

```go
// Decoration represents a block on a page
type Decoration struct {
    ID          int64
    PageID      int64
    BlockType   string    // hero, product_grid, banner, etc.
    BlockConfig string    // JSON configuration
    SortOrder   int
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

---

## Review Domain

Located in `admin/internal/domain/review/`

### Review Entity

```go
// Review represents a product review
type Review struct {
    ID             int64
    TenantID       shared.TenantID
    OrderID        int64
    ProductID      int64
    UserID         int64
    QualityRating  int       // 1-5
    ValueRating    int       // 1-5
    OverallRating  decimal.Decimal  // Calculated average
    Content        string
    Images         []string
    IsAnonymous    bool
    IsVerified     bool      // Verified purchase
    Status         ReviewStatus
    IsFeatured     bool
    HelpfulCount   int
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

type ReviewStatus string

const (
    ReviewStatusPending  ReviewStatus = "pending"
    ReviewStatusApproved ReviewStatus = "approved"
    ReviewStatusHidden   ReviewStatus = "hidden"
    ReviewStatusDeleted  ReviewStatus = "deleted"
)
```

### ReviewReply Entity

```go
// ReviewReply represents an admin reply to a review
type ReviewReply struct {
    ID        int64
    ReviewID  int64
    Content   string
    AdminID   int64
    AdminName string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

---

## Document History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-03-27 | Technical Team | Initial code documentation |
