# DB Rules (Database)

Rules for database schema design and operations.

## MUST

| # | Rule | Rationale |
|---|------|-----------|
| 1 | All database object names must use lowercase + underscores | Consistency, case-sensitivity issues across DBs |
| 2 | All tables must have a primary key that cannot be updated | Data integrity, replication support |
| 3 | Every table must include `created_at` / `updated_at` / `deleted_at` via `gorm.Model` | Audit trail, soft delete support |
| 4 | All database Model entities must embed `gorm.Model` (do NOT define time fields directly) | Consistency, following LeafAlloc pattern |
| 5 | Fields must explicitly define NOT NULL or default value | Prevent ambiguous NULL handling |
| 6 | Monetary fields must use DECIMAL / NUMERIC | Precision for financial calculations |
| 7 | Related field types must be exactly identical | Join performance, data integrity |
| 8 | DDL must be centrally managed and traceable | Change tracking, rollback capability |
| 9 | Production database changes must follow standard process | Safety, audit compliance |

## SHOULD

| # | Rule | Rationale |
|---|------|-----------|
| 9 | Single table fields ≤ 50 | Maintainability, query performance |
| 10 | Single table indexes ≤ 5 | Write performance, storage efficiency |
| 11 | SQL / index changes require review | Catch performance issues early |
| 12 | Design sharding strategy for large tables in advance | Avoid painful migrations later |
| 13 | Prefer database-level unique constraints | Stronger guarantee than application logic |

## FORBIDDEN

| # | Rule | Consequence |
|---|------|-------------|
| 14 | Cross-service JOIN | Service coupling, scaling issues |
| 15 | Foreign key constraints | Deployment complexity, performance impact |
| 16 | SELECT * | Unnecessary data transfer, schema coupling |
| 17 | Functions on indexed columns | Index bypass, full table scan |
| 18 | TEXT / BLOB in main tables | Query performance degradation |
| 19 | FLOAT / DOUBLE for monetary values | Precision loss, calculation errors |
| 20 | Manual production schema modifications | Untracked changes, potential data loss |

## Code Examples

### Table Schema

```sql
-- GOOD: Proper table structure with gorm.Model (time.Time timestamps)
CREATE TABLE orders (
    id              BIGINT          PRIMARY KEY,
    order_number    VARCHAR(32)     NOT NULL,
    customer_id     BIGINT          NOT NULL,
    total_amount    DECIMAL(19,4)   NOT NULL,  -- Precise decimal
    currency        VARCHAR(3)      NOT NULL,
    status          VARCHAR(32)     NOT NULL DEFAULT 'pending',
    created_at      TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP       NULL,                                        -- soft delete

    CONSTRAINT uk_order_number UNIQUE (order_number)
);

CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_orders_status ON orders(status);
```

### Query Patterns

```go
// BAD: SELECT *, function on indexed column
db.Raw("SELECT * FROM orders WHERE DATE(created_at) = ?", date)

// GOOD: Explicit columns, range query preserves index
db.Select("id", "order_number", "total_amount", "status").
    Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
    Find(&orders)
```

### Monetary Fields

```go
// BAD: Float for money
type Order struct {
    TotalAmount float64 `gorm:"column:total_amount"` // Precision loss!
}

// GOOD: Decimal type
import "github.com/shopspring/decimal"

type Order struct {
    TotalAmount decimal.Decimal `gorm:"column:total_amount;type:decimal(19,4)"`
    Currency    string          `gorm:"column:currency;type:varchar(3);not null"`
}
```

### Related Field Types

```go
// BAD: Type mismatch between tables
// orders.customer_id is BIGINT
// customers.id is INT  -- Type mismatch!

// GOOD: Exact type match
// orders.customer_id BIGINT references customers.id BIGINT
type Order struct {
    CustomerID int64 `gorm:"column:customer_id;type:bigint;not null"`
}

type Customer struct {
    ID int64 `gorm:"column:id;type:bigint;primaryKey"`
}
```

### gorm.Model Pattern

```go
import "gorm.io/gorm"

// GOOD: Model entity embeds gorm.Model (LeafAlloc pattern)
// gorm.Model includes: ID, CreatedAt, UpdatedAt, DeletedAt
type LeafAlloc struct {
    gorm.Model
}

func (*LeafAlloc) TableName() string {
    return "leaf_alloc"
}

// GOOD: Other entities should also embed gorm.Model
type ShippingTemplate struct {
    gorm.Model
    TenantID  int64  `gorm:"column:tenant_id;not null;index"`
    Name      string `gorm:"column:name;size:100;not null"`
    IsDefault bool   `gorm:"column:is_default;not null;default:false;index"`
    IsActive  bool   `gorm:"column:is_active;not null;default:true"`
}

// GOOD: Business time fields (PaidAt, ShippedAt, etc.) use *time.Time
type Order struct {
    gorm.Model
    Status     OrderStatus   `gorm:"column:status;not null;default:'pending_payment'"`
    PaidAt     *time.Time   `gorm:"column:paid_at"`
    ShippedAt  *time.Time   `gorm:"column:shipped_at"`
    DeliveredAt *time.Time  `gorm:"column:delivered_at"`
}
```

## Checklist

- [ ] Table/column names are lowercase with underscores
- [ ] Table has primary key
- [ ] Entity embeds `gorm.Model` (not defining time fields directly)
- [ ] Business time fields use `*time.Time`
- [ ] Fields have NOT NULL or default value
- [ ] Monetary fields use DECIMAL/NUMERIC
- [ ] Related fields have matching types
- [ ] No SELECT *
- [ ] No functions on indexed columns in WHERE
- [ ] No FLOAT/DOUBLE for money
- [ ] Changes tracked in migration files
