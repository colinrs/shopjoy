# Time Rules

Rules for time handling, storage, and transmission.

## MUST

| # | Rule | Rationale |
|---|------|-----------|
| 1 | Backend time storage must be unified to UTC | Consistent queries, no timezone confusion |
| 2 | Inter-service time transmission must include timezone | Unambiguous interpretation |
| 3 | API time fields must use strings (ISO 8601 / RFC3339) | Human-readable, timezone-aware |
| 4 | API layer must validate time format | Prevent invalid data entry |
| 5 | All business time calculations must be done on backend | Single source of truth, consistent logic |
| 6 | Domain entity time fields must use `time.Time` type | Type safety, consistent handling |
| 7 | Database timestamp fields must use `TIMESTAMP` type | Proper timezone support |
| 8 | AuditInfo value object must use `time.Time` for CreatedAt/UpdatedAt | Consistent with domain model |

## SHOULD

| # | Rule | Rationale |
|---|------|-----------|
| 9 | Inter-system time transmission should use UTC (Z suffix) | Simplify parsing, avoid conversion errors |
| 10 | Database time field types should be consistent within system | Reduce conversion complexity |
| 11 | Frontend should only convert timezone for display | Separation of concerns |

## FORBIDDEN

| # | Rule | Consequence |
|---|------|-------------|
| 12 | Storing local time in database | Query complexity, DST issues |
| 13 | Mixing multiple timezone semantics internally | Bugs, incorrect calculations |
| 14 | Timestamp + comment to indicate timezone | Fragile, easily misread |
| 15 | Frontend participating in business time logic | Inconsistent behavior, manipulation risk |
| 16 | Using `int64` Unix timestamp for domain entity time fields | Type safety issues, confusion |
| 17 | Using `float` for time-related calculations | Precision loss |

## Code Examples

### Time Storage and Retrieval

```go
// GOOD: Store UTC, convert for display
type Event struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    StartTime time.Time `json:"start_time" gorm:"column:start_time"` // Stored as UTC
}

// Store: always UTC
func (r *EventRepo) Create(ctx context.Context, db *gorm.DB, event *Event) error {
    event.StartTime = event.StartTime.UTC() // Ensure UTC before save
    return db.WithContext(ctx).Create(event).Error
}
```

### API Time Format

```go
// GOOD: ISO 8601 / RFC3339 string format
type CreateEventRequest struct {
    Name      string `json:"name" binding:"required"`
    StartTime string `json:"start_time" binding:"required"` // "2024-03-18T10:00:00Z"
}

type EventResponse struct {
    ID        int64  `json:"id"`
    Name      string `json:"name"`
    StartTime string `json:"start_time"` // Always RFC3339
}

func (l *EventLogic) Create(ctx context.Context, req *CreateEventRequest) (*EventResponse, error) {
    // Parse and validate time format
    startTime, err := time.Parse(time.RFC3339, req.StartTime)
    if err != nil {
        return nil, errorx.ErrInvalidTimeFormat
    }

    event := &model.Event{
        Name:      req.Name,
        StartTime: startTime.UTC(), // Convert to UTC for storage
    }
    // ...
}

func toResponse(event *model.Event) *EventResponse {
    return &EventResponse{
        ID:        event.ID,
        Name:      event.Name,
        StartTime: event.StartTime.Format(time.RFC3339), // UTC string
    }
}
```

### Time Validation

```go
// GOOD: Validate time format at API boundary
func ValidateTimeFormat(timeStr string) (time.Time, error) {
    t, err := time.Parse(time.RFC3339, timeStr)
    if err != nil {
        return time.Time{}, fmt.Errorf("invalid time format, expected RFC3339: %w", err)
    }
    return t.UTC(), nil
}
```

### Business Time Calculations

```go
// BAD: Time logic scattered, mixed timezones
func isExpired(createdAt time.Time) bool {
    localNow := time.Now() // Local timezone!
    return localNow.Sub(createdAt) > 24*time.Hour
}

// GOOD: All calculations in UTC on backend
func (s *OrderService) IsExpired(ctx context.Context, order *Order) bool {
    now := time.Now().UTC()
    expirationTime := order.CreatedAt.Add(s.config.OrderExpirationDuration)
    return now.After(expirationTime)
}

// GOOD: Time range queries in UTC
func (r *OrderRepo) FindByDateRange(ctx context.Context, db *gorm.DB, start, end time.Time) ([]*Order, error) {
    var orders []*Order
    err := db.WithContext(ctx).
        Where("created_at >= ? AND created_at < ?", start.UTC(), end.UTC()).
        Find(&orders).Error
    return orders, err
}
```

### Database Schema

```sql
-- GOOD: Consistent TIMESTAMP type, stored as UTC
CREATE TABLE events (
    id          BIGINT      PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    start_time  TIMESTAMP   NOT NULL,  -- UTC
    end_time    TIMESTAMP   NOT NULL,  -- UTC
    created_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### AuditInfo Value Object (审计信息)

Domain entities use embedded `AuditInfo` for timestamp tracking:

```go
// AuditInfo 审计信息（值对象）
type AuditInfo struct {
    CreatedAt time.Time // 创建时间
    UpdatedAt time.Time // 更新时间
    CreatedBy int64     // 创建人
    UpdatedBy int64     // 更新人
}

// Domain Entity with embedded AuditInfo
type Order struct {
    application.Model              // 嵌入，包含 ID, CreatedAt, UpdatedAt, DeletedAt
    TenantID     shared.TenantID  `gorm:"column:tenant_id;not null;index"`
    OrderNo      string           `gorm:"column:order_no;not null;uniqueIndex"`
    Status       OrderStatus      `gorm:"column:status;not null;default:0"`
    PaidAt       *time.Time      `gorm:"column:paid_at"`        // 业务时间字段用 *time.Time
    ShippedAt    *time.Time      `gorm:"column:shipped_at"`     // 业务时间字段用 *time.Time
    Audit        AuditInfo        `gorm:"embedded"`              // 审计信息
}
```

**Key Points:**
- `AuditInfo.CreatedAt/UpdatedAt` 使用 `time.Time` 类型
- 业务时间字段（如 `PaidAt`, `ShippedAt`）也使用 `*time.Time`
- 数据库字段统一使用 `TIMESTAMP` 类型
- **GORM 自动处理 `time.Time` ↔ `TIMESTAMP` 转换，无需手动转换**

### Repository Layer

GORM 直接使用 domain entity，`time.Time` 与数据库 `TIMESTAMP` 自动转换：

```go
// Repository 直接操作 domain entity
type orderRepo struct{}

func (r *orderRepo) Create(ctx context.Context, db *gorm.DB, order *Order) error {
    // GORM 自动将 order.CreatedAt (time.Time) 转换为数据库 TIMESTAMP
    return db.WithContext(ctx).Create(order).Error
}

func (r *orderRepo) FindByID(ctx context.Context, db *gorm.DB, id int64) (*Order, error) {
    var order Order
    // GORM 自动将数据库 TIMESTAMP 转换为 order.CreatedAt (time.Time)
    err := db.WithContext(ctx).First(&order, id).Error
    return &order, err
}
```

### Application Layer Usage

```go
// Domain entity uses time.Time
type Order struct {
    Audit AuditInfo
}

// Creating new entity - use time.Now().UTC()
func NewOrder(tenantID shared.TenantID, orderNo string) *Order {
    now := time.Now().UTC()
    return &Order{
        OrderNo: orderNo,
        Audit: AuditInfo{
            CreatedAt: now,  // time.Time
            UpdatedAt: now,
        },
    }
}

// Application layer response - format to string
func toResponse(o *Order) *OrderResponse {
    return &OrderResponse{
        ID:        o.ID,
        OrderNo:   o.OrderNo,
        CreatedAt: o.Audit.CreatedAt.Format(time.RFC3339),  // → RFC3339 string
        UpdatedAt: o.Audit.UpdatedAt.Format(time.RFC3339),
    }
}
```

**Common Mistakes to Avoid:**

```go
// BAD: Using time.Now().Unix() for time.Time field
brand := &Brand{
    Audit: AuditInfo{
        CreatedAt: time.Now().Unix(),  // WRONG! int64 vs time.Time
    },
}

// GOOD: Using time.Now().UTC() for time.Time field
brand := &Brand{
    Audit: AuditInfo{
        CreatedAt: time.Now().UTC(),  // CORRECT
    },
}

// BAD: Using time.Unix() when already time.Time
createdAtStr := time.Unix(order.Audit.CreatedAt.Unix(), 0).Format(time.RFC3339)

// GOOD: Use time.Time directly, format as RFC3339
createdAtStr := order.Audit.CreatedAt.Format(time.RFC3339)
```

## Checklist

- [ ] Time stored in database is UTC
- [ ] API time fields are RFC3339 strings
- [ ] Time format validated at API boundary
- [ ] Inter-service time includes timezone (prefer UTC)
- [ ] Business time logic is on backend only
- [ ] No local time stored in database
- [ ] No mixed timezone semantics in codebase
- [ ] Frontend only converts for display
- [ ] Domain entity time fields use `time.Time`
- [ ] Database timestamp fields use `TIMESTAMP` type
- [ ] AuditInfo uses `time.Time` for CreatedAt/UpdatedAt
- [ ] Repository directly uses domain entity (GORM handles conversion)
