# Fulfillment Module Domain Design

## Document Information

| Item | Value |
|------|-------|
| Document Title | Fulfillment Domain Design |
| Version | 1.0.0 |
| Status | Draft |
| Created | 2026-03-22 |
| Module Location | `admin/internal/domain/fulfillment/` |

---

## 1. Overview

### 1.1 Purpose

This document defines the domain model for the Fulfillment bounded context within ShopJoy. The fulfillment module handles the complete order fulfillment lifecycle, from shipment creation and logistics tracking to refund processing.

### 1.2 Scope

**In Scope (MVP):**
- Shipment creation with carrier and tracking information
- Shipment status lifecycle (pending → shipped → in_transit → delivered/failed/cancelled)
- Refund application and approval workflow
- Integration with Order via domain events

**Out of Scope (Phase 2):**
- Third-party logistics API integration
- Return merchandise authorization (RMA)
- Automated shipping label generation
- Partial refund (line item level)
- Multi-warehouse routing

### 1.3 Key Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Shipment & Order relationship | Shipment is aggregate root | Supports split shipments naturally; Order is read model for fulfillment state |
| Refund & Order relationship | Refund is independent aggregate | Supports multiple refund attempts; clean separation |
| Status integration | Event-driven | Decouples Shipment from Order; aligns with DDD |
| Event granularity | Specific events per transition | Explicit handlers; clear intent; easy mapping to notifications |
| Carrier handling | Value object + custom fallback | Common carriers predefined; supports custom carriers |

---

## 2. Domain Model

### 2.1 Context Map

```
┌─────────────────────────────────────────────────────────────┐
│                    Fulfillment Context                       │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────┐         ┌─────────────────┐           │
│  │    Shipment     │         │     Refund      │           │
│  │  (Aggregate)    │         │  (Aggregate)    │           │
│  └────────┬────────┘         └────────┬────────┘           │
│           │                           │                     │
│           │ references                │ references          │
│           ▼                           ▼                     │
│  ┌─────────────────────────────────────────────┐           │
│  │              Order (in Sales Context)        │           │
│  │  - status, fulfillment_status, refund_status │           │
│  └─────────────────────────────────────────────┘           │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 Aggregates

#### Shipment Aggregate

```
Shipment (Aggregate Root)
├── ID: int64
├── TenantID: TenantID
├── OrderID: string
├── ShipmentNo: string
├── Status: ShipmentStatus
├── Carrier: Carrier (Value Object)
├── TrackingNo: string
├── Items: []ShipmentItem
├── Weight: decimal.Decimal
├── Cost: Money
├── ShippedAt: *time.Time
├── DeliveredAt: *time.Time
├── Remark: string
├── Audit: AuditInfo
└── events: []DomainEvent

ShipmentItem (Entity)
├── ID: int64
├── OrderItemID: int64
├── ProductID: int64
├── SKUId: int64
├── Quantity: int
├── ProductName: string (snapshot)
├── SKUName: string (snapshot)
└── Image: string (snapshot)
```

#### Refund Aggregate

```
Refund (Aggregate Root)
├── ID: int64
├── TenantID: TenantID
├── OrderID: string
├── RefundNo: string
├── UserID: int64
├── Status: RefundStatus
├── Type: RefundType
├── ReasonCode: string
├── Reason: string
├── Description: string
├── Images: []string
├── Amount: Money
├── RejectReason: string
├── ApprovedAt: *time.Time
├── ApprovedBy: int64
├── CompletedAt: *time.Time
├── Audit: AuditInfo
└── events: []DomainEvent
```

---

## 3. State Machines

### 3.1 Shipment State Machine

```
                    ┌─────────────┐
                    │   pending   │
                    └──────┬──────┘
                           │
           ┌───────────────┼───────────────┐
           │               │               │
           ▼               │               ▼
    ┌─────────────┐        │        ┌─────────────┐
    │  cancelled  │        │        │   shipped   │
    └─────────────┘        │        └──────┬──────┘
                           │               │
                           │       ┌───────┴───────┐
                           │       │               │
                           │       ▼               ▼
                           │ ┌─────────────┐ ┌─────────────┐
                           │ │   failed    │ │  in_transit │
                           │ └─────────────┘ └──────┬──────┘
                           │                        │
                           │                ┌───────┴───────┐
                           │                │               │
                           │                ▼               ▼
                           │         ┌─────────────┐ ┌─────────────┐
                           │         │  delivered  │ │   failed    │
                           │         └─────────────┘ └─────────────┘
```

**Transitions:**

| From | To | Trigger | Precondition |
|------|-----|---------|--------------|
| pending | shipped | `Ship()` | Carrier and TrackingNo provided |
| pending | cancelled | `Cancel()` | Status is pending |
| shipped | in_transit | `MarkInTransit()` | Status is shipped |
| shipped | failed | `MarkFailed()` | Status is shipped (carrier reports failure before in-transit) |
| in_transit | delivered | `MarkDelivered()` | Status is in_transit |
| in_transit | failed | `MarkFailed()` | Status is in_transit |

### 3.2 Refund State Machine

```
                    ┌─────────────┐
                    │   pending   │
                    └──────┬──────┘
           ┌───────────────┼───────────────┐
           │               │               │
           ▼               ▼               ▼
    ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
    │  rejected   │ │  cancelled  │ │  approved   │
    └─────────────┘ └─────────────┘ └──────┬──────┘
                                           │
                                           ▼
                                    ┌─────────────┐
                                    │ processing  │
                                    └──────┬──────┘
                                           │
                                           ▼
                                    ┌─────────────┐
                                    │  completed  │
                                    └─────────────┘
```

**Transitions:**

| From | To | Trigger | Precondition |
|------|-----|---------|--------------|
| pending | approved | `Approve()` | Status is pending |
| pending | rejected | `Reject()` | Status is pending |
| pending | cancelled | `Cancel()` | Status is pending |
| approved | processing | `StartProcessing()` | Status is approved |
| processing | completed | `Complete()` | Status is processing |

---

## 4. Domain Events

### 4.1 Shipment Events

```go
// Emitted when a shipment is created
type ShipmentCreated struct {
    TenantID    shared.TenantID
    ShipmentID  int64
    OrderID     string
    ShipmentNo  string
}

// Emitted when shipment status changes to shipped
type ShipmentShipped struct {
    TenantID    shared.TenantID
    ShipmentID  int64
    OrderID     string
    Carrier     Carrier
    TrackingNo  string
    ShippedAt   time.Time
}

// Emitted when shipment is in transit
type ShipmentInTransit struct {
    TenantID    shared.TenantID
    ShipmentID  int64
    OrderID     string
}

// Emitted when shipment is delivered
type ShipmentDelivered struct {
    TenantID    shared.TenantID
    ShipmentID  int64
    OrderID     string
    DeliveredAt time.Time
}

// Emitted when shipment delivery fails
type ShipmentFailed struct {
    TenantID    shared.TenantID
    ShipmentID  int64
    OrderID     string
    Reason      string
}

// Emitted when shipment is cancelled
type ShipmentCancelled struct {
    TenantID    shared.TenantID
    ShipmentID  int64
    OrderID     string
}
```

### 4.2 Refund Events

```go
// Emitted when buyer requests a refund
type RefundRequested struct {
    TenantID   shared.TenantID
    RefundID   int64
    OrderID    string
    UserID     int64
    Amount     shared.Money
    ReasonCode string
}

// Emitted when merchant approves refund
type RefundApproved struct {
    TenantID   shared.TenantID
    RefundID   int64
    OrderID    string
    ApprovedBy int64
    ApprovedAt time.Time
}

// Emitted when payment refund is initiated
type RefundProcessing struct {
    TenantID shared.TenantID
    RefundID int64
    OrderID  string
}

// Emitted when refund is completed
type RefundCompleted struct {
    TenantID    shared.TenantID
    RefundID    int64
    OrderID     string
    CompletedAt time.Time
}

// Emitted when refund is rejected
type RefundRejected struct {
    TenantID     shared.TenantID
    RefundID     int64
    OrderID      string
    RejectedBy   int64
    RejectReason string
}

// Emitted when buyer cancels refund application
type RefundCancelled struct {
    TenantID    shared.TenantID
    RefundID    int64
    OrderID     string
    CancelledBy int64
    CancelledAt time.Time
}
```

---

## 5. Value Objects

### 5.1 Carrier

```go
// Carrier is a value object representing a logistics carrier
type Carrier struct {
    code string
    name string
}

// Predefined carriers
var (
    CarrierSF  = Carrier{code: "SF", name: "顺丰速运"}
    CarrierYT  = Carrier{code: "YT", name: "圆通快递"}
    CarrierZT  = Carrier{code: "ZT", name: "中通快递"}
    CarrierST  = Carrier{code: "ST", name: "申通快递"}
    CarrierYD  = Carrier{code: "YD", name: "韵达快递"}
    CarrierEMS = Carrier{code: "EMS", name: "EMS"}
    CarrierJD  = Carrier{code: "JD", name: "京东物流"}
)

// CustomCarrier creates a custom carrier for ones not in the predefined list
func CustomCarrier(name string) Carrier {
    return Carrier{code: "OTHER", name: name}
}

// CarrierFromCode returns a Carrier from code, with optional custom name
func CarrierFromCode(code string, customName string) Carrier

func (c Carrier) Code() string { return c.code }
func (c Carrier) Name() string { return c.name }
func (c Carrier) IsZero() bool { return c.code == "" }
```

### 5.2 RefundReason Codes

```go
// Reason codes - validated by application layer against reference table
const (
    ReasonDefective      = "DEFECTIVE"       // 商品质量问题
    ReasonWrongItem      = "WRONG_ITEM"      // 发错商品
    ReasonNotAsDescribed = "NOT_AS_DESCRIBED" // 商品与描述不符
    ReasonDamaged        = "DAMAGED"         // 运输损坏
    ReasonNoLongerNeeded = "NO_LONGER_NEEDED" // 不想要了
    ReasonLateDelivery   = "LATE_DELIVERY"   // 配送太慢
    ReasonOther          = "OTHER"           // 其他原因
)
```

### 5.3 Enums

```go
// ShipmentStatus defines the state of a shipment
type ShipmentStatus int

const (
    ShipmentStatusPending ShipmentStatus = iota
    ShipmentStatusShipped
    ShipmentStatusInTransit
    ShipmentStatusDelivered
    ShipmentStatusFailed
    ShipmentStatusCancelled
)

// RefundStatus defines the state of a refund
type RefundStatus int

const (
    RefundStatusPending RefundStatus = iota
    RefundStatusApproved
    RefundStatusRejected
    RefundStatusProcessing
    RefundStatusCompleted
    RefundStatusCancelled
)

// RefundType defines the type of refund
type RefundType int

const (
    RefundTypeFull RefundType = iota + 1
    // RefundTypePartial reserved for Phase 2
)

// FulfillmentStatus defines the derived fulfillment state of an order
// This enum is owned by the Fulfillment domain but stored on Order as a read model
type FulfillmentStatus int

const (
    FulfillmentStatusPending FulfillmentStatus = iota
    FulfillmentStatusPartialShipped
    FulfillmentStatusShipped
    FulfillmentStatusDelivered
)
```

---

## 6. Entities

### 6.1 Shipment

```go
// Shipment is the aggregate root for fulfillment
type Shipment struct {
    ID          int64
    TenantID    shared.TenantID
    OrderID     string
    ShipmentNo  string
    Status      ShipmentStatus
    Carrier     Carrier
    TrackingNo  string
    Items       []ShipmentItem
    Weight      decimal.Decimal
    Cost        shared.Money
    ShippedAt   *time.Time
    DeliveredAt *time.Time
    DeliveredBy int64              // User who marked as delivered
    Remark      string
    Audit       shared.AuditInfo
    Version     int                // For optimistic locking
    events      []shared.DomainEvent
}

// Constructor - enforces invariants at creation
func NewShipment(
    tenantID shared.TenantID,
    orderID string,
    items []ShipmentItem,
    createdBy int64,
) (*Shipment, error)

// Behavior methods
func (s *Shipment) Ship(carrier Carrier, trackingNo string, shippedBy int64) error
func (s *Shipment) MarkInTransit() error
func (s *Shipment) MarkDelivered(deliveredBy int64) error
func (s *Shipment) MarkFailed(reason string, updatedBy int64) error
func (s *Shipment) Cancel(cancelledBy int64) error

// Event access
func (s *Shipment) Events() []shared.DomainEvent
func (s *Shipment) ClearEvents()
```

### 6.2 ShipmentItem

```go
// ShipmentItem is an entity within Shipment aggregate
type ShipmentItem struct {
    ID          int64
    OrderItemID int64
    ProductID   int64
    SKUId       int64
    Quantity    int
    ProductName string  // snapshot
    SKUName     string  // snapshot
    Image       string  // snapshot
}

func NewShipmentItem(orderItemID, productID, skuID int64, quantity int) (ShipmentItem, error)
```

### 6.3 Refund

```go
// Refund is an aggregate root for refund processing
type Refund struct {
    ID           int64
    TenantID     shared.TenantID
    OrderID      string
    RefundNo     string
    UserID       int64
    Status       RefundStatus
    Type         RefundType
    ReasonCode   string
    Reason       string
    Description  string
    Images       []string
    Amount       shared.Money
    RejectReason string
    ApprovedAt   *time.Time
    ApprovedBy   int64
    CompletedAt  *time.Time
    Audit        shared.AuditInfo
    Version      int                // For optimistic locking
    events       []shared.DomainEvent
}

// Constructor
func NewRefund(
    tenantID shared.TenantID,
    orderID string,
    userID int64,
    amount shared.Money,
    reasonCode string,
    reason string,
    createdBy int64,
) (*Refund, error)

// Behavior methods
func (r *Refund) Approve(approvedBy int64) error
func (r *Refund) Reject(reason string, rejectedBy int64) error
func (r *Refund) StartProcessing() error
func (r *Refund) Complete() error
func (r *Refund) Cancel(cancelledBy int64) error

// Event access
func (r *Refund) Events() []shared.DomainEvent
func (r *Refund) ClearEvents()
```

---

## 7. Repository Interfaces

### 7.1 ShipmentRepository

```go
type ShipmentRepository interface {
    Create(ctx context.Context, db *gorm.DB, shipment *Shipment) error
    Update(ctx context.Context, db *gorm.DB, shipment *Shipment) error
    FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Shipment, error)
    FindByShipmentNo(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, shipmentNo string) (*Shipment, error)
    FindByOrderID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderID string) ([]*Shipment, error)
    FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query ShipmentQuery) ([]*Shipment, int64, error)
    ExistsByTrackingNo(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, trackingNo string) (bool, error)
}

type ShipmentQuery struct {
    shared.PageQuery
    Status    ShipmentStatus
    OrderID   string
    Carrier   string
    StartTime *time.Time
    EndTime   *time.Time
}
```

### 7.2 RefundRepository

```go
type RefundRepository interface {
    Create(ctx context.Context, db *gorm.DB, refund *Refund) error
    Update(ctx context.Context, db *gorm.DB, refund *Refund) error
    FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Refund, error)
    FindByRefundNo(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, refundNo string) (*Refund, error)
    FindByOrderID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderID string) ([]*Refund, error)
    FindActiveByOrderID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderID string) (*Refund, error)
    FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query RefundQuery) ([]*Refund, int64, error)
}

type RefundQuery struct {
    shared.PageQuery
    Status    RefundStatus
    OrderID   string
    UserID    int64
    StartTime *time.Time
    EndTime   *time.Time
}
```

---

## 8. Invariants

### 8.1 Domain Enforced

| Invariant | Location | Error |
|-----------|----------|-------|
| Shipment requires orderID | `NewShipment()` | `ErrOrderIDRequired` |
| Shipment requires at least one item | `NewShipment()` | `ErrShipmentItemsRequired` |
| ShipmentItem requires positive quantity | `NewShipmentItem()` | `ErrInvalidQuantity` |
| Ship requires carrier and trackingNo | `Ship()` | `ErrCarrierRequired`, `ErrTrackingNoRequired` |
| State transitions from valid state only | All state methods | `ErrShipmentCannotShip`, etc. |
| Refund requires orderID | `NewRefund()` | `ErrOrderIDRequired` |
| Refund requires positive amount | `NewRefund()` | `ErrInvalidRefundAmount` |
| Refund requires reasonCode | `NewRefund()` | `ErrReasonCodeRequired` |

### 8.2 Application Enforced

| Invariant | Check Location |
|-----------|----------------|
| Order must be paid before shipment | `CreateShipmentCmd` handler |
| Only one pending refund per order | `CreateRefundCmd` handler |
| Shipment items cannot exceed ordered quantity | `CreateShipmentCmd` handler |
| Refund time limit (configurable days) | `CreateRefundCmd` handler |
| Shipment number uniqueness | `GenerateShipmentNo()` with retry |
| Refund number uniqueness | `GenerateRefundNo()` with retry |

---

## 9. Event Handlers

### 9.1 Shipment Events

| Event | Handler Action |
|-------|----------------|
| `ShipmentCreated` | Log audit trail |
| `ShipmentShipped` | Update Order status to "shipped"; calculate fulfillment_status; send buyer notification |
| `ShipmentDelivered` | Update Order fulfillment_status; if all delivered, set Order status to "completed" |
| `ShipmentFailed` | Alert merchant; update Order fulfillment_status |
| `ShipmentCancelled` | Release reserved items; recalculate fulfillment_status |

### 9.2 Refund Events

| Event | Handler Action |
|-------|----------------|
| `RefundRequested` | Notify merchant of new refund request |
| `RefundApproved` | Update Order status to "refunding"; call payment service to initiate refund |
| `RefundProcessing` | Notify buyer "refund in progress" |
| `RefundCompleted` | Update Order status to "refunded"; release inventory; notify buyer |
| `RefundRejected` | Notify buyer with rejection reason; revert Order status |
| `RefundCancelled` | Revert Order status to previous state |

---

## 10. Fulfillment Status Calculation

Fulfillment status is derived from shipment states (shipment count based):

```go
func calculateFulfillmentStatus(shipments []*Shipment) FulfillmentStatus {
    if len(shipments) == 0 {
        return FulfillmentStatusPending
    }

    deliveredCount := 0
    shippedCount := 0
    failedCount := 0
    cancelledCount := 0
    pendingCount := 0

    for _, s := range shipments {
        switch s.Status {
        case ShipmentStatusDelivered:
            deliveredCount++
        case ShipmentStatusInTransit, ShipmentStatusShipped:
            shippedCount++
        case ShipmentStatusFailed:
            failedCount++
        case ShipmentStatusCancelled:
            cancelledCount++
        case ShipmentStatusPending:
            pendingCount++
        }
    }

    // All delivered
    if deliveredCount == len(shipments) {
        return FulfillmentStatusDelivered
    }

    // All cancelled - no active shipments, treat as pending
    if cancelledCount == len(shipments) {
        return FulfillmentStatusPending
    }

    // Count active shipments (non-cancelled)
    activeCount := len(shipments) - cancelledCount

    // If any delivered, consider partial/full shipped
    if deliveredCount > 0 {
        if deliveredCount == activeCount {
            return FulfillmentStatusDelivered
        }
        return FulfillmentStatusPartialShipped
    }

    // Check if all active shipments are shipped
    if pendingCount == 0 && shippedCount == activeCount {
        return FulfillmentStatusShipped
    }

    // Some pending, some shipped
    if pendingCount > 0 && shippedCount > 0 {
        return FulfillmentStatusPartialShipped
    }

    // All pending or mix of pending/cancelled
    return FulfillmentStatusPending
}
```

**Status Values:**

| Status | Meaning |
|--------|---------|
| `pending` | No shipments created, or all shipments cancelled |
| `partial_shipped` | Some shipments shipped/in_transit/delivered, some pending/cancelled |
| `shipped` | All active shipments are shipped or in transit (none delivered yet) |
| `delivered` | All active shipments delivered |

**Edge Cases:**

| Scenario | Result |
|----------|--------|
| All shipments cancelled | `pending` (reset to allow new shipments) |
| Mix of delivered + failed | `partial_shipped` (some succeeded) |
| Mix of delivered + cancelled | `partial_shipped` (delivered count > 0) |
| Mix of shipped + failed | `pending` (failed shipments need resolution) |

---

## 11. Number Generation

### 11.1 Shipment Number

Format: `SHP{YYYYMMDD}{sequence}`

Example: `SHP20260322001`

- Prefix: `SHP`
- Date: 8 digits (YYYYMMDD)
- Sequence: 3+ digits, resets daily per tenant

### 11.2 Refund Number

Format: `REF{YYYYMMDD}{sequence}`

Example: `REF20260322001`

- Prefix: `REF`
- Date: 8 digits (YYYYMMDD)
- Sequence: 3+ digits, resets daily per tenant

---

## 12. Database Schema

### 12.1 shipments

| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT | Primary key |
| tenant_id | BIGINT | Tenant ID |
| order_id | VARCHAR(64) | Order ID |
| shipment_no | VARCHAR(32) | Shipment number (unique per tenant) |
| status | TINYINT | 0=pending, 1=shipped, 2=in_transit, 3=delivered, 4=failed, 5=cancelled |
| carrier | VARCHAR(50) | Carrier name |
| carrier_code | VARCHAR(20) | Carrier code |
| tracking_no | VARCHAR(100) | Tracking number |
| weight | DECIMAL(10,2) | Weight in kg |
| cost_amount | BIGINT | Shipping cost in cents |
| cost_currency | VARCHAR(10) | Currency code |
| shipped_at | BIGINT | Shipment timestamp |
| delivered_at | BIGINT | Delivery timestamp |
| delivered_by | BIGINT | User who marked as delivered |
| remark | VARCHAR(500) | Remarks |
| version | INT | Version for optimistic locking |
| created_at | BIGINT | Creation timestamp |
| updated_at | BIGINT | Update timestamp |
| created_by | BIGINT | Creator ID |
| updated_by | BIGINT | Updater ID |
| deleted_at | BIGINT | Soft delete timestamp |

**Indexes:**
- PRIMARY KEY (`id`)
- UNIQUE INDEX `uk_shipment_no` (`tenant_id`, `shipment_no`)
- INDEX `idx_tenant_id` (`tenant_id`)
- INDEX `idx_order_id` (`order_id`)
- INDEX `idx_tracking_no` (`tracking_no`)
- INDEX `idx_status` (`status`)
- INDEX `idx_deleted_at` (`deleted_at`)

### 12.2 shipment_items

| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT | Primary key |
| tenant_id | BIGINT | Tenant ID |
| shipment_id | BIGINT | Shipment ID |
| order_item_id | BIGINT | Order item ID |
| product_id | BIGINT | Product ID |
| sku_id | BIGINT | SKU ID |
| product_name | VARCHAR(255) | Product name snapshot |
| sku_name | VARCHAR(255) | SKU name snapshot |
| image | VARCHAR(500) | Product image snapshot |
| quantity | INT | Quantity shipped |
| created_at | BIGINT | Creation timestamp |

### 12.3 refunds

| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT | Primary key |
| tenant_id | BIGINT | Tenant ID |
| order_id | VARCHAR(64) | Order ID |
| refund_no | VARCHAR(32) | Refund number (unique per tenant) |
| user_id | BIGINT | Buyer user ID |
| type | TINYINT | 1=full_refund |
| status | TINYINT | 0=pending, 1=approved, 2=rejected, 3=processing, 4=completed, 5=cancelled |
| reason_code | VARCHAR(50) | Reason code |
| reason | VARCHAR(500) | Reason text |
| description | TEXT | Detailed description |
| images | JSON | Evidence image URLs |
| amount | BIGINT | Refund amount in cents |
| currency | VARCHAR(10) | Currency code |
| reject_reason | VARCHAR(500) | Rejection reason |
| approved_at | BIGINT | Approval timestamp |
| approved_by | BIGINT | Approver ID |
| completed_at | BIGINT | Completion timestamp |
| version | INT | Version for optimistic locking |
| created_at | BIGINT | Creation timestamp |
| updated_at | BIGINT | Update timestamp |
| created_by | BIGINT | Creator ID |
| updated_by | BIGINT | Updater ID |
| deleted_at | BIGINT | Soft delete timestamp |

**Indexes:**
- PRIMARY KEY (`id`)
- UNIQUE INDEX `uk_refund_no` (`tenant_id`, `refund_no`)
- INDEX `idx_tenant_id` (`tenant_id`)
- INDEX `idx_order_id` (`order_id`)
- INDEX `idx_user_id` (`user_id`)
- INDEX `idx_status` (`status`)
- INDEX `idx_deleted_at` (`deleted_at`)

### 12.4 orders (extensions)

| Column | Type | Description |
|--------|------|-------------|
| fulfillment_status | TINYINT | 0=pending, 1=partial_shipped, 2=shipped, 3=delivered |
| refund_status | TINYINT | 0=none, 1=pending, 2=approved, 3=rejected, 4=processing, 5=completed |

---

## 13. File Structure

```
admin/internal/domain/fulfillment/
├── entity.go           # Shipment, Refund, ShipmentItem entities
├── status.go           # ShipmentStatus, RefundStatus enums
├── carrier.go          # Carrier value object
├── refund_reason.go    # Refund reason constants
├── events.go           # Domain event definitions
├── repository.go       # Repository interfaces
├── errors.go           # Domain errors

admin/internal/application/fulfillment/
├── commands.go         # Command handlers (CreateShipment, ShipOrder, etc.)
├── queries.go          # Query handlers (ListShipments, GetRefund, etc.)
├── handlers.go         # Event handlers
├── service.go          # Application service facade

admin/internal/infrastructure/persistence/fulfillment/
├── shipment_repository.go    # ShipmentRepository implementation
├── refund_repository.go      # RefundRepository implementation
├── models.go                 # GORM models
```

---

## 14. References

- PRD: `docs/prd/fulfillment_prd.md`
- Architecture: `docs/ARCHITECTURE.md`
- Shared Kernel: `pkg/domain/shared/`