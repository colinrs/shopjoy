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

## SHOULD

| # | Rule | Rationale |
|---|------|-----------|
| 6 | Inter-system time transmission should use UTC (Z suffix) | Simplify parsing, avoid conversion errors |
| 7 | Database time field types should be consistent within system | Reduce conversion complexity |
| 8 | Frontend should only convert timezone for display | Separation of concerns |

## FORBIDDEN

| # | Rule | Consequence |
|---|------|-------------|
| 9 | Storing local time in database | Query complexity, DST issues |
| 10 | Mixing multiple timezone semantics internally | Bugs, incorrect calculations |
| 11 | Timestamp + comment to indicate timezone | Fragile, easily misread |
| 12 | Frontend participating in business time logic | Inconsistent behavior, manipulation risk |

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

## Checklist

- [ ] Time stored in database is UTC
- [ ] API time fields are RFC3339 strings
- [ ] Time format validated at API boundary
- [ ] Inter-service time includes timezone (prefer UTC)
- [ ] Business time logic is on backend only
- [ ] No local time stored in database
- [ ] No mixed timezone semantics in codebase
- [ ] Frontend only converts for display
