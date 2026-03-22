// Package product implements the Product domain aggregate.
// This package contains the core business entities and repository interfaces
// for managing products in the ShopJoy e-commerce platform.
//
// # Domain Model
//
// The product domain follows Domain-Driven Design principles with:
//   - Product aggregate root
//   - SKU entity (product variant)
//   - Category entity (hierarchical)
//   - Brand entity
//   - Value objects: Status, Money, Dimensions
//
// # Product Status Lifecycle
//
// Products follow a state machine pattern:
//
//	draft ──┬──> on_sale ──┬──> off_sale ──┬──> on_sale
//	        │              │               │
//	        └──> deleted   └──> deleted    └──> deleted
//
// # Example Usage
//
// Creating a new product:
//
//	price := product.NewMoney(129900, "CNY") // ¥1,299.00
//	prod, err := product.NewProduct(1, tenantID, "Nike Air Max 270", "Description", price, 1)
//	if err != nil {
//	    return err
//	}
//	prod.PutOnSale() // Transition to on_sale status
//
// # Cross-Border Compliance
//
// Products support cross-border e-commerce compliance fields:
//   - HS Code (Harmonized System code for customs)
//   - Country of Origin (ISO country code)
//   - Weight and dimensions for shipping
//   - Dangerous goods declaration
package product

// Status represents the product lifecycle status.
// Products transition through defined states following business rules.
//
// Status transitions:
//   - draft -> on_sale, deleted
//   - on_sale -> off_sale, deleted
//   - off_sale -> on_sale, deleted
type Status int

const (
	// StatusDraft indicates a product in draft state.
	// Draft products are not visible to customers.
	// This is the initial state for new products.
	StatusDraft Status = iota

	// StatusOnSale indicates a product available for purchase.
	// Products must have stock > 0 to be put on sale.
	// Only on_sale products are visible in the shop.
	StatusOnSale

	// StatusOffSale indicates a product temporarily unavailable.
	// Product data is preserved but not purchasable.
	// Can be put back on sale.
	StatusOffSale

	// StatusDeleted indicates a soft-deleted product.
	// This is a terminal state - deleted products cannot be restored.
	StatusDeleted
)

// String returns the string representation of the status.
// Returns "draft", "on_sale", "off_sale", "deleted", or "unknown".
func (s Status) String() string

// IsValid checks if the status value is within valid range.
func (s Status) IsValid() bool

// CanTransitionTo checks if a status transition is allowed.
// Use this to validate state changes before applying them.
//
// Example:
//
//	if !product.Status.CanTransitionTo(product.StatusOnSale) {
//	    return errors.New("invalid status transition")
//	}
func (s Status) CanTransitionTo(target Status) bool

// Money represents a monetary value with currency binding.
// Amount is stored as int64 in cents to avoid floating-point precision issues.
// All monetary calculations should use this type.
//
// # Precision
//
// Amount is in smallest currency unit (cents):
//   - 129900 cents = $1,299.00 USD
//   - 129900 cents = ¥1,299.00 CNY
//
// # Currency Validation
//
// All arithmetic operations validate currency matching:
//
//	total, err := price1.Add(price2) // Returns error if currencies differ
type Money struct {
	// Amount in smallest currency unit (cents).
	// For $99.99, use 9999.
	Amount int64

	// Currency is the ISO 4217 currency code.
	// Examples: "CNY", "USD", "EUR", "GBP", "AUD"
	Currency string
}

// NewMoney creates a new Money value object.
// If currency is empty, defaults to "CNY".
//
// Example:
//
//	price := product.NewMoney(9900, "USD") // $99.00 USD
func NewMoney(amount int64, currency string) Money

// Add returns the sum of two money values.
// Returns an error if currencies don't match.
//
// Example:
//
//	total, err := price1.Add(price2)
//	if err != nil {
//	    // Handle currency mismatch
//	}
func (m Money) Add(other Money) (Money, error)

// Subtract returns the difference of two money values.
// Returns an error if currencies don't match or result would be negative.
func (m Money) Subtract(other Money) (Money, error)

// Equals checks if two money values are identical.
// Both amount and currency must match.
func (m Money) Equals(other Money) bool

// Dimensions represents physical product dimensions.
// Used for shipping calculations and compliance.
type Dimensions struct {
	// Length in the specified unit.
	Length decimal.Decimal

	// Width in the specified unit.
	Width decimal.Decimal

	// Height in the specified unit.
	Height decimal.Decimal

	// Unit of measurement, typically "cm".
	Unit string
}

// Product is the aggregate root for the product domain.
// It represents a sellable item (SPU - Standard Product Unit).
//
// # Product vs SKU
//
// Product represents the abstract product (e.g., "Nike Air Max 270").
// SKU represents a specific variant (e.g., "Nike Air Max 270, Black, Size 42").
// Products with variants have is_matrix_product = true.
//
// # Multi-Tenancy
//
// All products are scoped to a tenant via tenant_id.
// Repository methods require tenant_id for data isolation.
//
// # Compliance Fields
//
// For cross-border e-commerce, products include:
//   - HS Code: Customs classification
//   - COO (Country of Origin): Manufacturing country
//   - Weight: For shipping calculations
//   - Dangerous Goods: Special handling requirements
type Product struct {
	// ID is the unique product identifier.
	ID int64

	// TenantID is the tenant this product belongs to.
	// Required for multi-tenant data isolation.
	TenantID shared.TenantID

	// SKU is the default SKU code for simple products.
	// For matrix products, this may be empty or a reference code.
	SKU string

	// Name is the product display name.
	// Required field, max 200 characters.
	Name string

	// Description is the product description.
	// Supports markdown formatting.
	Description string

	// Price is the selling price.
	// Stored as Money value object with currency.
	Price Money

	// CostPrice is the purchase/manufacturing cost.
	// Optional, used for profit calculations.
	CostPrice Money

	// Stock is the total inventory count.
	// Must be > 0 for products to be put on sale.
	Stock int

	// Status is the current product lifecycle state.
	// Use PutOnSale/TakeOffSale to change status.
	Status Status

	// CategoryID is the product's category.
	// References the categories table.
	CategoryID int64

	// Brand is the brand name (legacy field).
	// Prefer using brand_id for new implementations.
	Brand string

	// Tags are search/filter tags.
	// Stored as JSON array in database.
	Tags []string

	// Images are product image URLs.
	// First image is used as primary/thumbnail.
	Images []string

	// IsMatrixProduct indicates if product has variants.
	// If true, create SKUs for each variant.
	IsMatrixProduct bool

	// HSCode is the Harmonized System code.
	// Required for cross-border shipping.
	// Format: 6-10 digit code (e.g., "64041100")
	HSCode string

	// COO is the Country of Origin.
	// ISO 3166-1 alpha-2 code (e.g., "CN", "VN", "US")
	COO string

	// Weight of the product.
	// Use WeightUnit to specify measurement unit.
	Weight decimal.Decimal

	// WeightUnit is the weight measurement unit.
	// Common values: "g", "kg", "lb"
	WeightUnit string

	// Dimensions are the physical dimensions.
	// Used for shipping volume calculations.
	Dimensions Dimensions

	// DangerousGoods are hazardous material identifiers.
	// Examples: ["battery", "flammable", "magnetic"]
	DangerousGoods []string

	// CreatedAt is the product creation timestamp.
	CreatedAt time.Time

	// UpdatedAt is the last modification timestamp.
	UpdatedAt time.Time
}

// TableName returns the database table name for GORM.
func (p *Product) TableName() string

// NewProduct creates a new product with validation.
// Returns an error if name is empty or price is invalid.
//
// The product is created in StatusDraft state.
// Use PutOnSale to make it available for purchase.
//
// Example:
//
//	price := product.NewMoney(129900, "CNY")
//	prod, err := product.NewProduct(1, tenantID, "Product Name", "Description", price, 1)
//	if err != nil {
//	    return err
//	}
//	repo.Create(ctx, db, prod)
func NewProduct(id int64, tenantID shared.TenantID, name, description string, price Money, categoryID int64) (*Product, error)

// NewProductWithCompliance creates a product with compliance fields.
// Use for cross-border products that need HS codes and origin info.
func NewProductWithCompliance(id int64, tenantID shared.TenantID, name, description, sku string, price Money, categoryID int64) (*Product, error)

// SetCompliance sets cross-border compliance information.
// Updates the product's HS code, origin, weight, and dimensions.
func (p *Product) SetCompliance(hsCode, coo string, weight decimal.Decimal, weightUnit string, dims Dimensions)

// HasComplianceInfo checks if product has required compliance data.
// Returns true if HS code, COO, and weight are all set.
func (p *Product) HasComplianceInfo() bool

// IsDangerousGoods checks if product requires special handling.
func (p *Product) IsDangerousGoods() bool

// PutOnSale transitions the product to on_sale status.
// Returns error if:
//   - Product is deleted
//   - Status transition is not allowed
//   - Stock is zero or negative
//
// Example:
//
//	if err := product.PutOnSale(); err != nil {
//	    return fmt.Errorf("cannot put on sale: %w", err)
//	}
func (p *Product) PutOnSale() error

// TakeOffSale transitions the product to off_sale status.
// Product becomes unavailable for purchase but is not deleted.
// Returns error if product is deleted or transition is not allowed.
func (p *Product) TakeOffSale() error

// UpdateStock sets the product stock to a new value.
// Use for manual stock updates (e.g., inventory counts).
// Returns error if product is deleted or quantity is negative.
func (p *Product) UpdateStock(quantity int) error

// DeductStock reduces stock by the specified quantity.
// Used by order processing when payment is confirmed.
// Returns error if product is not on sale or insufficient stock.
func (p *Product) DeductStock(quantity int) error

// UpdatePrice changes the product price.
// Returns error if product is deleted or new price is invalid.
func (p *Product) UpdatePrice(newPrice Money) error

// SoftDelete marks the product as deleted.
// This is a terminal state - the product cannot be restored.
func (p *Product) SoftDelete() error

// IsOnSale checks if the product is available for purchase.
// Returns true only if status is on_sale AND stock > 0.
func (p *Product) IsOnSale() bool

// Repository defines the interface for product persistence.
// Implementations should handle multi-tenant data isolation.
type Repository interface {
	// Create persists a new product.
	Create(ctx context.Context, db *gorm.DB, product *Product) error

	// Update modifies an existing product.
	Update(ctx context.Context, db *gorm.DB, product *Product) error

	// Delete soft-deletes a product by ID.
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error

	// FindByID retrieves a product by ID.
	// Returns error if not found or tenant mismatch.
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Product, error)

	// FindByIDs retrieves multiple products by IDs.
	FindByIDs(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, ids []int64) ([]*Product, error)

	// FindList retrieves products matching query criteria.
	FindList(ctx context.Context, db *gorm.DB, query Query) ([]*Product, int64, error)

	// UpdateStock atomically updates product stock.
	// Delta can be positive (add) or negative (subtract).
	UpdateStock(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, delta int) error

	// Exists checks if a product exists.
	Exists(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (bool, error)
}

// Query represents product search criteria.
// Use for filtering and pagination in list operations.
type Query struct {
	// TenantID for multi-tenant isolation (required).
	TenantID shared.TenantID

	// Name filters by partial name match.
	Name string

	// CategoryID filters by category.
	CategoryID int64

	// Status filters by product status.
	// nil means no status filter.
	Status *Status

	// MinPrice filters products above this price.
	MinPrice *int64

	// MaxPrice filters products below this price.
	MaxPrice *int64

	// MarketID filters products available in a market.
	MarketID int64

	// Page number (1-indexed).
	Page int

	// PageSize limits results per page.
	// Max: 100, Default: 20.
	PageSize int
}

// Validate normalizes query parameters.
// Ensures page >= 1 and page_size is within bounds.
func (q Query) Validate() error

// Offset calculates the database offset for pagination.
// Returns (Page - 1) * PageSize.
func (q Query) Offset() int

// Limit returns the page size for pagination.
func (q Query) Limit() int