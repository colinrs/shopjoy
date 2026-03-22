# Error Rules

Rules for error handling and custom error definitions.

## MUST

| # | Rule | Rationale |
|---|------|-----------|
| 1 | All business errors must be defined in `pkg/code/code.go` | Centralized error management, consistent HTTP codes |
| 2 | Each module must use its assigned error code range | Avoid conflicts, easier debugging |
| 3 | Errors must include proper HTTP status code | Correct API response semantics |
| 4 | Use `errors.Is()` for error comparison | Support error wrapping and comparison |
| 5 | Return `code.ErrXxx` directly, do not wrap with custom messages | Consistent error responses |

## SHOULD

| # | Rule | Rationale |
|---|------|-----------|
| 6 | Include i18n message key in error definition | Support multi-language error messages |
| 7 | Use descriptive error names starting with `Err` | Clear naming convention |
| 8 | Group related errors with consecutive codes | Logical organization |

## FORBIDDEN

| # | Rule | Consequence |
|---|------|-------------|
| 9 | Using `errors.New()` in application/domain layers | Inconsistent error handling, missing HTTP codes |
| 10 | Creating local error variables in packages | Decentralized errors, hard to maintain |
| 11 | Using `fmt.Errorf()` for business errors | Lost error code, incorrect HTTP status |
| 12 | Defining errors outside `pkg/code/code.go` | Fragmented error management |

## Code Examples

### Error Definition

```go
// GOOD: Define in pkg/code/code.go with HTTP code and business code
var (
    // Order Module (40xxx)
    ErrOrderNotFound      = &Err{HTTPCode: http.StatusNotFound, Code: 40001, Msg: "order not found"}
    ErrOrderInvalidStatus = &Err{HTTPCode: http.StatusBadRequest, Code: 40002, Msg: "invalid order status"}
    ErrOrderAlreadyPaid   = &Err{HTTPCode: http.StatusBadRequest, Code: 40003, Msg: "order already paid"}
)

// BAD: Local error definition
var ErrOrderNotFound = errors.New("order not found")
```

### Error Usage

```go
// GOOD: Return code error directly
func (s *OrderService) GetOrder(ctx context.Context, id int64) (*Order, error) {
    order, err := s.repo.FindByID(ctx, s.db, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, code.ErrOrderNotFound
        }
        return nil, err
    }
    return order, nil
}

// GOOD: Error comparison with errors.Is
func (h *Handler) HandleOrder(w http.ResponseWriter, r *http.Request) {
    order, err := s.GetOrder(ctx, id)
    if err != nil {
        if errors.Is(err, code.ErrOrderNotFound) {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

// BAD: Using errors.New()
func (s *OrderService) GetOrder(ctx context.Context, id int64) (*Order, error) {
    return nil, errors.New("order not found")  // No HTTP code, inconsistent
}

// BAD: Creating local error
var ErrInvalidOrder = errors.New("invalid order")  // Should be in pkg/code
```

### Error Code Ranges

```
Module          Range       Example
----------------------------------------
Admin User      10xxx       10001, 10002
User            11xxx       11001, 11002
Product         30xxx       30001, 30002
Category        301xx       30101, 30102
Order           40xxx       40001, 40002
Payment         50xxx       50001, 50002
Cart            60xxx       60001, 60002
Coupon          70xxx       70001, 70002
UserCoupon      701xx       70101, 70102
Promotion       80xxx       80001, 80002
Tenant          90xxx       90001, 90002
Role            100xxx      100001, 100002
Shop            110xxx      110001, 110002
Fulfillment     120xxx      120001, 120002
  - Shipment    120xxx      120001-120099
  - Refund      1201xx      120101-120199
  - Carrier     1202xx      120201-120299
  - RefundReason 1203xx     120301-120399
Shared          200xxx      200001, 200002
Auth            130xxx      130001, 130002
Cache           140xxx      140001, 140002
Market          150xxx      150001, 150002
ProductMarket   160xxx      160001, 160002
Inventory       170xxx      170001, 170002
Brand           180xxx      180001, 180002
SKU             190xxx      190001, 190002
```

## Checklist

- [ ] All business errors defined in `pkg/code/code.go`
- [ ] No `errors.New()` in application/domain layers
- [ ] No local error variable definitions
- [ ] Error code within assigned module range
- [ ] Proper HTTP status code set
- [ ] Use `errors.Is()` for comparison