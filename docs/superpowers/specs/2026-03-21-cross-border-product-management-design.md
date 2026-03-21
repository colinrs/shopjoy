# Cross-Border Product Management System Design

**Date:** 2026-03-21
**Status:** Draft
**Scope:** Product management with multi-market support (Shopify Markets pattern)

## Overview

Design a cross-border product management system that allows sellers to manage products across multiple markets (US, UK, DE, FR, AU) with per-market pricing, visibility control, and localization support.

### Core Principles

1. **Market as First-Class Entity** - Markets are configuration units with currency, language, and tax rules
2. **Product Data Sharing** - Core product data (name, images, compliance) is global; market-specific data (price, visibility, translations) overlays on top
3. **View as Market** - Product list can be filtered by market to see what's visible in each market
4. **Push to Market** - Enable products in new markets with pricing strategy options

## Target Markets (MVP)

| Market | Code | Currency | Language |
|--------|------|----------|----------|
| United States | US | USD | en |
| United Kingdom | UK | GBP | en |
| Germany | DE | EUR | de |
| France | FR | EUR | fr |
| Australia | AU | AUD | en |

**Languages Supported:** EN (base), DE, FR

## Data Model

### Market Entity

```go
type Market struct {
    ID             int64
    Code           string    // US, UK, DE, FR, AU
    Name           string    // "United States"
    Currency       string    // USD, GBP, EUR, AUD
    DefaultLanguage string   // en, de, fr
    Flag           string    // emoji or image URL
    IsActive       bool
    IsDefault      bool      // Primary Market
    TaxRules       TaxConfig `gorm:"type:json"`
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

type TaxConfig struct {
    VATRate      decimal.Decimal // For EU markets
    GSTRate      decimal.Decimal // For AU
    IOSSEnabled  bool            // Low value import scheme
    IncludeTax   bool            // Display prices tax-inclusive
}
```

### Product Entity

```go
type Product struct {
    ID             int64
    TenantID       int64
    SKU            string
    Name           string
    Description    string
    CategoryID     int64
    Brand          string
    Tags           []string `gorm:"type:json"`
    Images         []string `gorm:"type:json"`
    Status         ProductStatus // draft, on_sale, off_sale
    IsMatrixProduct bool         // Has variants or standalone

    // Compliance (Global)
    HSCode         string
    COO            string    // Country of Origin
    Weight         decimal.Decimal
    WeightUnit     string    // g, kg
    Dimensions     Dimensions `gorm:"embedded"`
    DangerousGoods []string `gorm:"type:json"` // battery, liquid, magnet, powder

    // Aggregated
    TotalStock     int

    CreatedAt      time.Time
    UpdatedAt      time.Time
}

type ProductStatus string

const (
    ProductStatusDraft  ProductStatus = "draft"
    ProductStatusOnSale ProductStatus = "on_sale"
    ProductStatusOffSale ProductStatus = "off_sale"
)

type Dimensions struct {
    Length decimal.Decimal
    Width  decimal.Decimal
    Height decimal.Decimal
    Unit   string // cm
}
```

### ProductVariant Entity

```go
type ProductVariant struct {
    ID           int64
    TenantID     int64
    ProductID    int64
    SKU          string
    Options      VariantOptions `gorm:"type:json"` // {size: "M", color: "Black"}
    IsStandalone bool           // true = independent SKU, false = matrix variant
    Images       []string       `gorm:"type:json"`
    SortOrder    int

    // Optional compliance override
    HSCodeOverride string
    WeightOverride decimal.Decimal

    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type VariantOptions map[string]string
```

### ProductMarket Entity (Market-Specific Data)

```go
type ProductMarket struct {
    ID              int64
    TenantID        int64
    ProductID       int64
    VariantID       *int64    // NULL for products without variants; each variant has its own row for matrix products
    MarketID        int64

    // Visibility
    IsEnabled       bool      // Product visible in this market
    StatusOverride  *ProductStatus // Override global status per market

    // Pricing
    Price           decimal.Decimal
    CompareAtPrice  *decimal.Decimal // Original/sale comparison price

    // Inventory
    StockAlertThreshold int // Low stock alert for this market

    PublishedAt     *time.Time
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### ProductLocalization Entity

```go
type ProductLocalization struct {
    ID             int64
    TenantID       int64
    ProductID      int64
    VariantID      *int64
    LanguageCode   string    // en, de, fr

    // Localized Content
    Title          string
    SEOTitle       string
    MetaDescription string
    BulletPoints   []string `gorm:"type:json"`
    Description    string
    ImageAlts      map[string]string `gorm:"type:json"` // image_url -> alt text

    IsComplete     bool      // Translation complete flag
    CreatedAt      time.Time
    UpdatedAt      time.Time
}
```

### StockLog Entity

```go
type StockLog struct {
    ID              int64
    TenantID        int64
    ProductID       int64
    VariantID       *int64
    QuantityChange  int       // Positive = increase, negative = decrease
    Reason          string    // "sale", "restock", "adjustment", "transfer"
    ReferenceType   string    // "order", "manual", "inventory_check"
    ReferenceID     string    // Order ID or other reference
    CreatedAt       time.Time
}
```

## API Endpoints

### Product Management

```
GET    /api/v1/products                    # List products (with market filter)
POST   /api/v1/products                    # Create product
GET    /api/v1/products/:id                # Get product detail
PUT    /api/v1/products/:id                # Update product
DELETE /api/v1/products/:id                # Delete product

POST   /api/v1/products/:id/push-to-market # Push product to new markets
POST   /api/v1/products/batch              # Batch operations
```

### Product Variants

```
GET    /api/v1/products/:id/variants       # List variants
POST   /api/v1/products/:id/variants       # Create variant
PUT    /api/v1/products/:id/variants/:vid  # Update variant
DELETE /api/v1/products/:id/variants/:vid  # Delete variant
```

### Product Market Operations

```
GET    /api/v1/products/:id/markets        # Get market configurations
PUT    /api/v1/products/:id/markets/:market_id  # Update market config
DELETE /api/v1/products/:id/markets/:market_id  # Remove from market
```

### Product Localization

```
GET    /api/v1/products/:id/localizations  # Get all localizations
PUT    /api/v1/products/:id/localizations/:lang  # Update localization
```

### Market Management (Settings)

```
GET    /api/v1/markets                     # List markets
POST   /api/v1/markets                     # Create market
PUT    /api/v1/markets/:id                 # Update market
DELETE /api/v1/markets/:id                 # Delete market
```

### Inventory

```
GET    /api/v1/products/:id/inventory      # Get inventory info
POST   /api/v1/products/:id/inventory/adjust  # Adjust stock
GET    /api/v1/products/:id/inventory/logs    # Stock logs
```

## Frontend Pages

### Product List Page (`/products`)

**Features:**
1. **Market Filter Bar** (top): "全部市场" | US | UK | DE | FR | AU
2. **Table Columns:**
   - Selection checkbox
   - Product info (image, name, SKU, tags)
   - Market visibility tags (US ✓ UK ✓ DE ✗ FR ✗ AU ✗)
   - Price (in selected market currency)
   - Stock
   - Status
   - Actions

3. **Batch Actions:**
   - Push to markets (enable visibility)
   - Remove from markets (disable visibility)
   - Status change (on/off sale)
   - Pricing adjustment (±% or fixed per market)
   - Compliance fill (batch HS Code, COO)
   - Delete

**View as Market Behavior:**
- Selecting a market filters the list to show products visible in that market
- Price column shows that market's currency and price
- Status reflects market-specific status if overridden

### Product Detail Page (`/products/:id`) - 6 Tabs

**Tab 1: Basic Info**
- Product name, SKU, category, brand, tags
- Images/videos upload and management
- Compliance section:
  - HS Code (6-digit base + market extensions)
  - Country of Origin (dropdown)
  - Weight, dimensions
  - Dangerous goods checkboxes (battery, liquid, magnet, powder, etc.)
- Warning banner if compliance fields missing

**Tab 2: Variants**
- Variant options configuration (e.g., Size: S/M/L, Color: Black/White)
- Variants table:
  - SKU, options, images, stock, actions
- Support for:
  - Matrix variants (shared base info)
  - Standalone SKUs (independent products)

**Tab 3: Markets**
- Table view with columns: Market | Status | Currency | Price Set | Actions
- Bulk enable/disable selected markets
- "Push to Market" button for adding new markets

**Tab 4: Pricing**
- Market selector (button group: US | UK | DE | FR | AU)
- Price table for selected market:
  - All variants with price input
  - Compare at price
- "Copy from [Market]" quick action

**Tab 5: Localization**
- Language selector (EN | DE | FR)
- Form fields for selected language:
  - Title, SEO Title, Meta Description
  - Bullet Points (5 items)
  - Description (rich text)
  - Image Alt texts
- "Copy from English" to initialize
- Translation completeness indicator

**Tab 6: Inventory**
- Total stock summary
- Per-market safety stock thresholds
- Stock adjustment form
- Stock movement log

### Market Management Page (`/settings/markets`)

- Market list with: Name, Code, Currency, Status, Actions
- Create/Edit market modal:
  - Country/Region selection
  - Currency
  - Default language
  - Tax configuration (VAT/GST rates, IOSS)
- Primary market designation
- Enable/disable toggle

## Push to Market Flow

1. User selects product(s) from list or detail page
2. Clicks "Push to Market"
3. Modal appears:
   - Target market selection (checkboxes)
   - Pricing strategy:
     - Based on [source market] price ±X% (manual conversion - seller enters the price)
     - Manual entry (seller sets price directly in target currency)
   - Copy localization option
4. System creates `ProductMarket` records
5. Products appear in target markets with "draft" status
6. Seller reviews and sets final price → Publishes

**Currency Handling:** Pricing is manual per market. No automatic FX rate conversion. Seller is responsible for setting appropriate prices in each market's currency.

## Compliance Validation

- **Soft validation only** - warnings, not blockers
- Warning banner on product if missing:
  - HS Code
  - Country of Origin
  - Weight/Dimensions
- Filter option: "Show products missing compliance data"

## Inventory Management

- **Hybrid model**: Shared global stock + per-market alerts
- Single stock pool shared across all markets
- Per-market safety stock threshold (alert when below)
- Stock changes logged with reason codes
- No allocation per market (simpler for MVP)

**Concurrent Orders:** Stock is deducted atomically at order creation. If multiple markets order simultaneously, database transactions ensure no overselling. Consider adding optimistic locking or stock reservation if race conditions become an issue at scale.

## Batch Operations

| Operation | Description |
|-----------|-------------|
| Push to Markets | Enable products in selected markets |
| Remove from Markets | Disable products in selected markets |
| Status Change | Set on_sale/off_sale status |
| Pricing Adjustment | Apply ±% or fixed amount per market |
| Compliance Fill | Batch set HS Code, COO |
| Export/Import | CSV with multi-market columns |

## Implementation Phases

### Phase 1: Core Infrastructure
- Market entity and CRUD
- Product entity updates (compliance fields)
- ProductMarket entity
- Basic API endpoints

### Phase 2: Product Management
- Product list with market filter
- Product detail - Basic Info tab
- Product detail - Markets tab
- Push to Market flow

### Phase 3: Variants & Pricing
- Variant management
- Pricing tab
- Per-market pricing

### Phase 4: Localization
- ProductLocalization entity
- Localization tab
- Language switching

### Phase 5: Inventory & Batch
- Inventory management
- Stock logs
- Batch operations

### Phase 6: Dashboard
- Multi-market sales comparison
- Compliance warnings
- Performance by market

## Database Migration Required

New tables:
- `markets`
- `product_markets`
- `product_localizations`
- `product_variants` (if not exists)
- `stock_logs`

Modified tables:
- `products` - Add compliance fields, is_matrix_product

## Dependencies

- `github.com/shopspring/decimal` for monetary calculations
- GORM for database operations
- Element Plus for admin UI components