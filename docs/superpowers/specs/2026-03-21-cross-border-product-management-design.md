# Cross-Border Product Management System Design

**Date:** 2026-03-21
**Status:** Draft
**Scope:** Product management with multi-market support (Shopify Markets pattern)

## Current Codebase Status

### Backend (Go)

**Implemented:**
| Component | Status | Location |
|-----------|--------|----------|
| Market Entity | ✅ Done | `admin/internal/domain/market/entity.go` |
| Market Repository | ✅ Done | `admin/internal/infrastructure/persistence/market_repository.go` |
| Market API Handlers | ✅ Done | `admin/internal/handler/markets/` |
| Market API Definition | ✅ Done | `admin/desc/market.api` |
| Product Entity | ✅ Done (partial) | `admin/internal/domain/product/entity.go` |
| Product Repository | ✅ Done (partial) | `admin/internal/infrastructure/persistence/product_repository.go` |
| Product API Handlers | ✅ Done (basic) | `admin/internal/handler/products/` |
| Product API Definition | ✅ Done (basic) | `admin/desc/product.api` |

**Missing (to implement):**
| Component | Description |
|-----------|-------------|
| ProductMarket Entity | Market-specific product data |
| ProductMarket Repository | Persistence layer |
| ProductMarket API | New endpoints |
| ProductVariant Entity | Product variants/SKUs |
| ProductVariant Repository | Persistence layer |
| ProductVariant API | New endpoints |
| ProductLocalization Entity | Translations |
| ProductLocalization Repository | Persistence layer |
| ProductLocalization API | New endpoints |
| StockLog Entity | Inventory tracking |
| StockLog Repository | Persistence layer |
| StockLog API | New endpoints |

**Existing Product Entity Fields:**
The `Product` entity already includes compliance fields:
```go
// admin/internal/domain/product/entity.go
type Product struct {
    // ... existing fields ...
    HSCode         string          // ✅ Exists
    COO            string          // ✅ Exists (Country of Origin)
    Weight         decimal.Decimal // ✅ Exists
    WeightUnit     string          // ✅ Exists
    Dimensions     Dimensions      // ✅ Exists
    DangerousGoods []string        // ✅ Exists
}
```

**API Gap:** The `.api` definition does not expose these fields. Need to extend DTOs.

### Frontend (Vue 3)

**Implemented:**
| Page | Status | Location |
|------|--------|----------|
| Login | ✅ Done | `shop-admin/src/views/login/` |
| Dashboard | ✅ Done (mock) | `shop-admin/src/views/dashboard/` |
| Products List | ✅ Done (mock) | `shop-admin/src/views/products/` |
| Orders | ✅ Done (mock) | `shop-admin/src/views/orders/` |
| Users | ✅ Done (mock) | `shop-admin/src/views/users/` |
| Admin Users | ✅ Done (mock) | `shop-admin/src/views/admin-users/` |
| Promotions | ✅ Done (mock) | `shop-admin/src/views/promotions/` |
| Shop Settings | ✅ Done (mock) | `shop-admin/src/views/shop/` |

**Missing:**
- Market Management page
- Product Detail page (with 6 tabs)
- API integration (currently using mock data)

### API Extension Strategy

**Principle:** Extend existing APIs without breaking changes.

**Product API Changes (`admin/desc/product.api`):**

1. **Extend CreateProductReq:**
```go
type CreateProductReq {
    // Existing fields
    Name        string `json:"name"`
    Description string `json:"description,optional"`
    Price       int64  `json:"price"`
    Currency    string `json:"currency,optional"`
    CostPrice   int64  `json:"cost_price,optional"`
    CategoryID  int64  `json:"category_id"`

    // New fields (all optional for backward compatibility)
    SKU            string   `json:"sku,optional"`
    Brand          string   `json:"brand,optional"`
    Tags           []string `json:"tags,optional"`
    Images         []string `json:"images,optional"`
    IsMatrixProduct bool     `json:"is_matrix_product,optional"`

    // Compliance fields
    HSCode         string  `json:"hs_code,optional"`
    COO            string  `json:"coo,optional"`           // Country of Origin
    Weight         string  `json:"weight,optional"`        // decimal as string
    WeightUnit     string  `json:"weight_unit,optional"`
    Length         string  `json:"length,optional"`
    Width          string  `json:"width,optional"`
    Height         string  `json:"height,optional"`
    DangerousGoods []string `json:"dangerous_goods,optional"`
}
```

2. **Extend ListProductReq (add market filter):**
```go
type ListProductReq {
    // Existing
    Name       string `form:"name,optional"`
    CategoryID int64  `form:"category_id,optional"`
    Status     string `form:"status,optional"`
    MinPrice   int64  `form:"min_price,optional"`
    MaxPrice   int64  `form:"max_price,optional"`
    Page       int    `form:"page,default=1"`
    PageSize   int    `form:"page_size,default=20"`

    // New
    MarketID   int64  `form:"market_id,optional"`   // Filter by market
}
```

3. **Extend ProductDetailResp:**
```go
type ProductDetailResp {
    // Existing fields ...

    // New fields
    SKU            string   `json:"sku"`
    Brand          string   `json:"brand"`
    Tags           []string `json:"tags"`
    Images         []string `json:"images"`
    IsMatrixProduct bool     `json:"is_matrix_product"`

    // Compliance
    HSCode         string   `json:"hs_code"`
    COO            string   `json:"coo"`
    Weight         string   `json:"weight"`
    WeightUnit     string   `json:"weight_unit"`
    Length         string   `json:"length"`
    Width          string   `json:"width"`
    Height         string   `json:"height"`
    DangerousGoods []string `json:"dangerous_goods"`

    // Market info (for list view)
    Markets        []ProductMarketInfo `json:"markets,optional"`
}

type ProductMarketInfo {
    MarketID   int64  `json:"market_id"`
    MarketCode string `json:"market_code"`
    IsEnabled  bool   `json:"is_enabled"`
    Price      string `json:"price"`
    Currency   string `json:"currency"`
}
```

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

**Backend:**
1. Extend `admin/desc/product.api`:
   - Add compliance fields to CreateProductReq/UpdateProductReq
   - Add compliance fields to ProductDetailResp
   - Add market_id filter to ListProductReq

2. Create ProductMarket domain entity (`admin/internal/domain/product/product_market.go`):
   - Already exists, verify fields match design

3. Create ProductMarket repository (`admin/internal/infrastructure/persistence/product_market_repository.go`):
   - CRUD operations
   - FindByProductID, FindByMarketID methods

4. Create ProductMarket API definition (`admin/desc/product_market.api`):
   - GET `/api/v1/products/:id/markets` - List product's market configs
   - PUT `/api/v1/products/:id/markets/:market_id` - Update market config
   - POST `/api/v1/products/:id/push-to-market` - Push to new markets
   - DELETE `/api/v1/products/:id/markets/:market_id` - Remove from market

5. Implement API handlers

**Frontend:**
1. Create API service (`shop-admin/src/api/market.ts`):
   - getMarkets, createMarket, updateMarket, deleteMarket

2. Create Market Management page (`shop-admin/src/views/settings/markets/`):
   - Market list table
   - Create/Edit dialog
   - Tax configuration form

3. Update product API service with compliance fields

**Database:**
- `product_markets` table already exists in SQL schema

---

### Phase 2: Product Management

**Backend:**
1. Update product list handler:
   - Support market_id filter
   - Return market visibility info per product

2. Create push-to-market handler:
   - Batch create ProductMarket records
   - Validate market exists and is active

**Frontend:**
1. Update Product List page:
   - Add market filter bar (button group: 全部 | US | UK | DE | FR | AU)
   - Add market visibility tags column
   - Show price in selected market currency

2. Create Product Detail page route (`/products/:id`):
   - Implement tab navigation (6 tabs)
   - Tab 1: Basic Info form with compliance fields

3. Implement Markets Tab:
   - Market enable/disable table
   - Push to Market dialog

**Push to Market Dialog:**
- Target market checkboxes
- Price input per market (manual entry, no auto FX)
- Option to copy from existing market

---

### Phase 3: Variants & Pricing

**Backend:**
1. Create ProductVariant entity (`admin/internal/domain/product/variant.go`)

2. Create ProductVariant repository

3. Create Variant API (`admin/desc/variant.api`):
   - GET `/api/v1/products/:id/variants`
   - POST `/api/v1/products/:id/variants`
   - PUT `/api/v1/products/:id/variants/:vid`
   - DELETE `/api/v1/products/:id/variants/:vid`

4. Update ProductMarket to support variant-level pricing

**Frontend:**
1. Implement Variants Tab:
   - Variant options configuration (Size, Color, etc.)
   - Variants table with SKU, options, stock
   - Generate variants from options

2. Implement Pricing Tab:
   - Market selector (button group)
   - Price table per variant
   - Compare at price
   - Copy from market action

---

### Phase 4: Localization

**Backend:**
1. Create ProductLocalization entity

2. Create ProductLocalization repository

3. Create Localization API:
   - GET `/api/v1/products/:id/localizations`
   - PUT `/api/v1/products/:id/localizations/:lang`

**Frontend:**
1. Implement Localization Tab:
   - Language selector (EN | DE | FR)
   - Translation form (title, SEO, bullets, description)
   - Copy from English action
   - Completeness indicator

---

### Phase 5: Inventory & Batch

**Backend:**
1. Create StockLog entity

2. Create StockLog repository

3. Create Inventory API:
   - GET `/api/v1/products/:id/inventory`
   - POST `/api/v1/products/:id/inventory/adjust`
   - GET `/api/v1/products/:id/inventory/logs`

4. Create Batch Operations API:
   - POST `/api/v1/products/batch/push-to-markets`
   - POST `/api/v1/products/batch/status`
   - POST `/api/v1/products/batch/compliance`

**Frontend:**
1. Implement Inventory Tab:
   - Stock summary
   - Adjustment form
   - Movement log table

2. Implement Batch Actions:
   - Push to Markets dialog (multi-select)
   - Batch status change
   - Batch compliance fill

---

### Phase 6: Dashboard

**Backend:**
1. Create Dashboard API:
   - GET `/api/v1/dashboard/market-stats`
   - GET `/api/v1/dashboard/compliance-warnings`

**Frontend:**
1. Update Dashboard page:
   - Multi-market sales comparison chart
   - Compliance warnings widget
   - Top products by market

---

## File Structure (New Files)

### Backend
```
admin/
├── desc/
│   ├── product.api          # Extended
│   ├── product_market.api   # New
│   ├── variant.api          # New
│   ├── localization.api     # New
│   └── inventory.api        # New
├── internal/
│   ├── domain/
│   │   └── product/
│   │       ├── entity.go            # Extended
│   │       ├── product_market.go    # Exists
│   │       ├── variant.go           # New
│   │       └── localization.go      # New
│   ├── infrastructure/
│   │   └── persistence/
│   │       ├── product_repository.go        # Extended
│   │       ├── product_market_repository.go # Exists
│   │       ├── variant_repository.go        # New
│   │       ├── localization_repository.go   # New
│   │       └── stock_log_repository.go      # New
│   └── handler/
│       ├── products/         # Extended
│       ├── product_markets/  # New
│       ├── variants/         # New
│       └── inventory/        # New
```

### Frontend
```
shop-admin/src/
├── api/
│   ├── product.ts       # Extended
│   ├── market.ts        # New
│   ├── variant.ts       # New
│   └── inventory.ts     # New
├── views/
│   ├── products/
│   │   ├── index.vue          # Updated
│   │   └── [id]/
│   │       └── index.vue      # New (detail page)
│   └── settings/
│       └── markets/
│           └── index.vue      # New
└── utils/
    └── types.ts        # Extended with new interfaces
```

## Database Migration Required

### Existing Tables (SQL already created in `sql/` directory)
| Table | Status | SQL File |
|-------|--------|----------|
| `markets` | ✅ Created | `sql/market.sql` |
| `products` | ✅ Created | `sql/product.sql` |
| `product_markets` | ✅ Created | `sql/product.sql` |
| `skus` | ✅ Created | `sql/product.sql` |

### New Tables Required
| Table | Description | Phase |
|-------|-------------|-------|
| `product_variants` | Product variants (if different from skus) | Phase 3 |
| `product_localizations` | Product translations | Phase 4 |
| `stock_logs` | Inventory movement logs | Phase 5 |

### SQL Files Location
All SQL files are in: `/sql/` directory

Run order:
1. `tenant.sql`
2. `admin_user.sql`
3. `user.sql`
4. `role.sql`
5. `market.sql` ← Market table
6. `product.sql` ← Products, SKUs, ProductMarkets
7. `storefront.sql`
8. `coupon.sql`
9. `promotion.sql`
10. `cart.sql`
11. `order.sql`
12. `payment.sql`
13. `fulfillment.sql`

## Dependencies

### Backend (Go)
- `github.com/shopspring/decimal` for monetary calculations
- GORM for database operations
- go-zero framework

### Frontend (Vue 3)
- Element Plus for UI components
- Vue Router for navigation
- Pinia for state management
- Axios for HTTP requests
- ECharts for dashboard charts