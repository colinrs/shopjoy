# Price Rules (Price / Monetary Values)

Rules for handling prices, amounts, and monetary calculations.

## MUST

| # | Rule | Rationale |
|---|------|-----------|
| 1 | All monetary calculations must use decimal-precise types | Avoid floating-point precision errors |
| 2 | Database monetary fields must use DECIMAL / NUMERIC | Exact storage, no rounding issues |
| 3 | Monetary values must originate from strings or precise types | Prevent float contamination |
| 4 | Monetary fields must define precision and decimal places | Clear contract, consistent handling |
| 5 | In multi-market systems, amounts must be bound to currency | Prevent currency mismatch errors |
| 6 | API monetary values must use string type representing yuan (元) | Frontend should not need to divide by 100; e.g., "1.99" means 1.99元, not 199分 |
| 7 | Internal monetary calculations may use different units | Convert at API boundary; domain layer uses decimal.Decimal |

## SHOULD

| # | Rule | Rationale |
|---|------|-----------|
| 6 | Use exact comparison methods for monetary values | Avoid epsilon comparison issues |
| 7 | Centralize currency conversion logic | Single source of truth, easier auditing |
| 8 | Unify monetary unit within system | Reduce conversion complexity |

## FORBIDDEN

| # | Rule | Consequence |
|---|------|-------------|
| 9 | Using float/double to calculate or store amounts | Precision loss (e.g., 0.1 + 0.2 ≠ 0.3) |
| 10 | Using round/math to fix floating-point errors | Masking the real problem, accumulating errors |
| 11 | Using FLOAT/DOUBLE for monetary database fields | Silent data corruption |
| 12 | Monetary amounts existing without currency | Ambiguous value, calculation errors |

## Code Examples

### Money Type Definition

```go
import "github.com/shopspring/decimal"

// GOOD: Money type with currency binding
type Money struct {
    Amount   decimal.Decimal `json:"amount"`
    Currency string          `json:"currency"` // ISO 4217: USD, EUR, SGD
}

// Constructor ensures proper initialization
func NewMoney(amount decimal.Decimal, currency string) Money {
    return Money{
        Amount:   amount,
        Currency: currency,
    }
}

// From string - safe parsing
func MoneyFromString(amountStr, currency string) (Money, error) {
    amount, err := decimal.NewFromString(amountStr)
    if err != nil {
        return Money{}, fmt.Errorf("invalid amount: %w", err)
    }
    return NewMoney(amount, currency), nil
}
```

### Monetary Calculations

```go
// BAD: Float arithmetic
func calculateTotal(price float64, quantity int) float64 {
    return price * float64(quantity) // Precision loss!
}

// GOOD: Decimal arithmetic
func CalculateTotal(unitPrice Money, quantity int) Money {
    total := unitPrice.Amount.Mul(decimal.NewFromInt(int64(quantity)))
    return NewMoney(total, unitPrice.Currency)
}

// GOOD: Safe addition with currency check
func (m Money) Add(other Money) (Money, error) {
    if m.Currency != other.Currency {
        return Money{}, errors.New("currency mismatch")
    }
    return NewMoney(m.Amount.Add(other.Amount), m.Currency), nil
}

// GOOD: Percentage calculation
func (m Money) ApplyDiscount(discountPercent decimal.Decimal) Money {
    discount := m.Amount.Mul(discountPercent).Div(decimal.NewFromInt(100))
    return NewMoney(m.Amount.Sub(discount), m.Currency)
}
```

### Monetary Comparison

```go
// BAD: Float comparison
func isGreater(a, b float64) bool {
    return a > b // Unreliable for money!
}

// GOOD: Decimal comparison
func (m Money) IsGreaterThan(other Money) bool {
    return m.Amount.GreaterThan(other.Amount)
}

func (m Money) Equals(other Money) bool {
    return m.Currency == other.Currency && m.Amount.Equal(other.Amount)
}

func (m Money) IsZero() bool {
    return m.Amount.IsZero()
}

func (m Money) IsPositive() bool {
    return m.Amount.IsPositive()
}
```

### Database Schema and Model

```sql
-- GOOD: DECIMAL with explicit precision
CREATE TABLE order_items (
    id              BIGINT          PRIMARY KEY,
    order_id        BIGINT          NOT NULL,
    product_id      BIGINT          NOT NULL,
    unit_price      DECIMAL(19,4)   NOT NULL,  -- Up to 15 digits + 4 decimals
    quantity        INT             NOT NULL,
    line_total      DECIMAL(19,4)   NOT NULL,
    currency        VARCHAR(3)      NOT NULL,  -- ISO 4217
    created_at      TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

```go
// GOOD: Model with decimal type
type OrderItem struct {
    ID         int64           `gorm:"column:id;primaryKey"`
    OrderID    int64           `gorm:"column:order_id;not null"`
    ProductID  int64           `gorm:"column:product_id;not null"`
    UnitPrice  decimal.Decimal `gorm:"column:unit_price;type:decimal(19,4);not null"`
    Quantity   int             `gorm:"column:quantity;not null"`
    LineTotal  decimal.Decimal `gorm:"column:line_total;type:decimal(19,4);not null"`
    Currency   string          `gorm:"column:currency;type:varchar(3);not null"`
    CreatedAt  time.Time       `gorm:"column:created_at"`
    UpdatedAt  time.Time       `gorm:"column:updated_at"`
}
```

### API Request/Response

```go
// GOOD: String for monetary values in API
type CreateOrderRequest struct {
    Items []OrderItemRequest `json:"items" binding:"required"`
}

type OrderItemRequest struct {
    ProductID string `json:"product_id" binding:"required"`
    UnitPrice string `json:"unit_price" binding:"required"` // "99.99"
    Quantity  int    `json:"quantity" binding:"required,min=1"`
    Currency  string `json:"currency" binding:"required,len=3"` // "USD"
}

// Parse string to decimal safely
func (r *OrderItemRequest) ToMoney() (Money, error) {
    return MoneyFromString(r.UnitPrice, r.Currency)
}
```

### Currency Conversion (Centralized)

```go
// GOOD: Centralized conversion service
type CurrencyConverter struct {
    rateProvider RateProvider
}

func (c *CurrencyConverter) Convert(ctx context.Context, amount Money, targetCurrency string) (Money, error) {
    if amount.Currency == targetCurrency {
        return amount, nil
    }

    rate, err := c.rateProvider.GetRate(ctx, amount.Currency, targetCurrency)
    if err != nil {
        return Money{}, fmt.Errorf("get rate: %w", err)
    }

    converted := amount.Amount.Mul(rate)
    return NewMoney(converted, targetCurrency), nil
}
```

## Checklist

- [ ] Monetary calculations use `decimal.Decimal`
- [ ] Database monetary fields are DECIMAL/NUMERIC
- [ ] API monetary values are strings (representing yuan, not cents)
- [ ] Money type includes currency
- [ ] No float/double for money anywhere
- [ ] Monetary comparisons use decimal methods (Equal, GreaterThan)
- [ ] Currency conversion is centralized
- [ ] Precision explicitly defined (e.g., DECIMAL(19,4))
- [ ] API accepts string prices like "1.99" (1.99元), not 199 (cents)
