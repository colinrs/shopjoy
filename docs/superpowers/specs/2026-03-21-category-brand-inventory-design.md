# Category, Brand, and Inventory Management Design

**Date:** 2026-03-21
**Status:** Draft
**Scope:** Admin Backend (Merchant Management)

---

## 1. Category Management

### 1.1 Core Business Positioning

Categories are the skeleton structure for merchants to organize products, help buyers navigate, and enable search/recommendation/SEO. In cross-border scenarios, categories serve both **platform search optimization** and **merchant internal operations**.

### 1.2 Data Model

#### 1.2.1 Category Entity

```go
type Category struct {
    ID          int64          `gorm:"column:id;primaryKey"`
    TenantID    int64          `gorm:"column:tenant_id;not null;index"`
    ParentID    int64          `gorm:"column:parent_id;default:0;index"`  // 0 = root level
    Name        string         `gorm:"column:name;type:varchar(100);not null"`
    Code        string         `gorm:"column:code;type:varchar(50);uniqueIndex:idx_tenant_code"`
    Level       int            `gorm:"column:level;not null;default:1"`   // 1, 2, 3... (soft limit: 3)
    Sort        int            `gorm:"column:sort;default:0"`
    Icon        string         `gorm:"column:icon;type:varchar(500)"`     // Icon URL
    Image       string         `gorm:"column:image;type:varchar(500)"`    // Category image URL

    // SEO Fields
    SeoTitle       string      `gorm:"column:seo_title;type:varchar(200)"`
    SeoDescription string      `gorm:"column:seo_description;type:varchar(500)"`

    // Status
    Status       int8           `gorm:"column:status;not null;default:1"`  // 1=enabled, 0=disabled

    // Audit
    CreatedAt    time.Time      `gorm:"column:created_at;not null"`
    UpdatedAt    time.Time      `gorm:"column:updated_at;not null"`
    DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
}
```

#### 1.2.2 Category Market Visibility

```go
type CategoryMarket struct {
    ID          int64     `gorm:"column:id;primaryKey"`
    TenantID    int64     `gorm:"column:tenant_id;not null;index"`
    CategoryID  int64     `gorm:"column:category_id;not null;index"`
    MarketID    int64     `gorm:"column:market_id;not null;index"`
    IsVisible   bool      `gorm:"column:is_visible;not null;default:true"`
    CreatedAt   time.Time `gorm:"column:created_at;not null"`
    UpdatedAt   time.Time `gorm:"column:updated_at;not null"`

    // Unique constraint: (tenant_id, category_id, market_id)
}
```

### 1.3 Business Rules

#### 1.3.1 Category Structure

| Rule | Description |
|------|-------------|
| Tree structure | Supports parent-child hierarchy via `parent_id` |
| Level limit | Soft limit of 3 levels (configurable, allows more) |
| Leaf category | Category with no children; only leaf categories can have products |

#### 1.3.2 Create/Edit Rules

| Rule | Description |
|------|-------------|
| Unique name | Sibling categories (same `parent_id`) cannot have duplicate names |
| Tenant isolation | Each tenant has independent category tree |
| Level calculation | Auto-calculated from parent level + 1 |
| Code uniqueness | `code` is unique within tenant |

#### 1.3.3 Category Status

| Action | Effect |
|--------|--------|
| Disable category | Category and products hidden on storefront; still manageable in Admin |
| Enable category | Category and products visible again on storefront |
| Delete category | Products auto-migrated to parent category |

#### 1.3.4 Leaf Category Migration

When a category becomes non-leaf (new children added):

1. System detects existing products on the category
2. Merchant is prompted to migrate products
3. Migration options:
   - Migrate all to specific child category
   - Keep on parent (non-leaf, but allowed for backward compatibility)

#### 1.3.5 Market Visibility Rules

| Rule | Description |
|------|-------------|
| Visibility scope | Only hides category navigation; products still searchable |
| Product override | Products can override category's market visibility setting |
| Inheritance | Child categories inherit parent's market settings (cannot override) |

#### 1.3.6 Sorting and Display

| Feature | Description |
|---------|-------------|
| Sorting | Unified `sort` field for same-level ordering |
| Drag-and-drop | Supports cross-level drag (auto-updates `level` and `parent_id`) |
| Icons/Images | Displayed in frontend category navigation |

### 1.4 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/categories` | Create category |
| PUT | `/api/v1/categories/:id` | Update category |
| DELETE | `/api/v1/categories/:id` | Delete category |
| PUT | `/api/v1/categories/:id/status` | Enable/disable category |
| GET | `/api/v1/categories/:id` | Get category detail |
| GET | `/api/v1/categories` | List categories (flat) |
| GET | `/api/v1/categories/tree` | Get category tree |
| PUT | `/api/v1/categories/sort` | Batch update sort order |
| POST | `/api/v1/categories/:id/migrate-products` | Migrate products to child category |
| GET | `/api/v1/categories/:id/product-count` | Get product count under category |
| PUT | `/api/v1/categories/:id/market-visibility` | Update market visibility |
| GET | `/api/v1/categories/:id/market-visibility` | Get market visibility settings |

### 1.5 Database Schema

```sql
CREATE TABLE categories (
    id              BIGINT          PRIMARY KEY,
    tenant_id       BIGINT          NOT NULL,
    parent_id       BIGINT          NOT NULL DEFAULT 0,
    name            VARCHAR(100)    NOT NULL,
    code            VARCHAR(50),
    level           INT             NOT NULL DEFAULT 1,
    sort            INT             NOT NULL DEFAULT 0,
    icon            VARCHAR(500),
    image           VARCHAR(500),
    seo_title       VARCHAR(200),
    seo_description VARCHAR(500),
    status          TINYINT         NOT NULL DEFAULT 1,
    created_at      TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP       NULL,

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_parent_id (parent_id),
    UNIQUE INDEX idx_tenant_code (tenant_id, code)
);

CREATE TABLE category_markets (
    id              BIGINT      PRIMARY KEY,
    tenant_id       BIGINT      NOT NULL,
    category_id     BIGINT      NOT NULL,
    market_id       BIGINT      NOT NULL,
    is_visible      BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_category_id (category_id),
    UNIQUE INDEX idx_tenant_category_market (tenant_id, category_id, market_id)
);
```

---

## 2. Brand Management

### 2.1 Core Business Positioning

Brands enhance product trust, enable brand pages, SEO, and advertising. In cross-border scenarios, brands also carry **compliance** and **intellectual property** responsibilities.

### 2.2 Data Model

#### 2.2.1 Brand Entity

```go
type Brand struct {
    ID          int64          `gorm:"column:id;primaryKey"`
    TenantID    int64          `gorm:"column:tenant_id;not null;index"`
    Name        string         `gorm:"column:name;type:varchar(100);not null;uniqueIndex:idx_tenant_name"`
    Logo        string         `gorm:"column:logo;type:varchar(500)"`
    Description string         `gorm:"column:description;type:text"`
    Website     string         `gorm:"column:website;type:varchar(500)"`
    Sort        int            `gorm:"column:sort;default:0"`

    // Brand Page
    EnablePage  bool           `gorm:"column:enable_page;default:false"`  // Enable brand page

    // Compliance (Cross-border)
    TrademarkNumber  string   `gorm:"column:trademark_number;type:varchar(100)"`
    TrademarkCountry string   `gorm:"column:trademark_country;type:varchar(10)"`  // ISO country code

    // Status
    Status      int8           `gorm:"column:status;not null;default:1"`  // 1=enabled, 0=disabled

    // Audit
    CreatedAt   time.Time      `gorm:"column:created_at;not null"`
    UpdatedAt   time.Time      `gorm:"column:updated_at;not null"`
    DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}
```

#### 2.2.2 Brand Market Visibility

```go
type BrandMarket struct {
    ID          int64     `gorm:"column:id;primaryKey"`
    TenantID    int64     `gorm:"column:tenant_id;not null;index"`
    BrandID     int64     `gorm:"column:brand_id;not null;index"`
    MarketID    int64     `gorm:"column:market_id;not null;index"`
    IsVisible   bool      `gorm:"column:is_visible;not null;default:true"`
    CreatedAt   time.Time `gorm:"column:created_at;not null"`
    UpdatedAt   time.Time `gorm:"column:updated_at;not null"`
}
```

### 2.3 Business Rules

#### 2.3.1 Product-Brand Relationship

| Rule | Description |
|------|-------------|
| Storage | Product has `brand_id` foreign key |
| Optional | Products can have no brand (`brand_id` can be NULL) |
| One brand | Each product can only link to one brand |

#### 2.3.2 Brand Status

| Action | Effect |
|--------|--------|
| Disable brand | Brand name/logo hidden on storefront; products still visible |
| Enable brand | Brand info displayed normally |
| Delete brand | Products' `brand_id` set to NULL |

#### 2.3.3 Brand Page

| Rule | Description |
|------|-------------|
| Toggle | `enable_page` controls whether brand page is generated |
| URL format | `/brands/:id` or `/brands/:slug` |
| Content | Brand logo, description, products list |

#### 2.3.4 Market Visibility Rules

| Rule | Description |
|------|-------------|
| Visibility control | Brand can be hidden in specific markets |
| Product effect | When brand hidden in market, products don't show brand info |
| Brand page | Brand page returns 404 in hidden markets |

### 2.4 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/brands` | Create brand |
| PUT | `/api/v1/brands/:id` | Update brand |
| DELETE | `/api/v1/brands/:id` | Delete brand |
| PUT | `/api/v1/brands/:id/status` | Enable/disable brand |
| GET | `/api/v1/brands/:id` | Get brand detail |
| GET | `/api/v1/brands` | List brands (paginated) |
| GET | `/api/v1/brands/:id/product-count` | Get product count |
| PUT | `/api/v1/brands/:id/market-visibility` | Update market visibility |
| GET | `/api/v1/brands/:id/market-visibility` | Get market visibility settings |
| PUT | `/api/v1/brands/:id/toggle-page` | Toggle brand page |

### 2.5 Database Schema

```sql
CREATE TABLE brands (
    id                  BIGINT          PRIMARY KEY,
    tenant_id           BIGINT          NOT NULL,
    name                VARCHAR(100)    NOT NULL,
    logo                VARCHAR(500),
    description         TEXT,
    website             VARCHAR(500),
    sort                INT             NOT NULL DEFAULT 0,
    enable_page         BOOLEAN         NOT NULL DEFAULT FALSE,
    trademark_number    VARCHAR(100),
    trademark_country   VARCHAR(10),
    status              TINYINT         NOT NULL DEFAULT 1,
    created_at          TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at          TIMESTAMP       NULL,

    INDEX idx_tenant_id (tenant_id),
    UNIQUE INDEX idx_tenant_name (tenant_id, name)
);

CREATE TABLE brand_markets (
    id              BIGINT      PRIMARY KEY,
    tenant_id       BIGINT      NOT NULL,
    brand_id        BIGINT      NOT NULL,
    market_id       BIGINT      NOT NULL,
    is_visible      BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_brand_id (brand_id),
    UNIQUE INDEX idx_tenant_brand_market (tenant_id, brand_id, market_id)
);

-- Add brand_id to products table
ALTER TABLE products ADD COLUMN brand_id BIGINT NULL;
ALTER TABLE products ADD INDEX idx_brand_id (brand_id);
```

---

## 3. Inventory Management

### 3.1 Core Business Positioning

Inventory management tracks product availability, prevents overselling, and supports cross-border multi-warehouse scenarios.

### 3.2 Data Model

#### 3.2.1 Inventory Types

| Type | Description | Update Frequency |
|------|-------------|------------------|
| Available Stock | Real sellable inventory | Real-time |
| Locked Stock | Ordered but unpaid/in-transit | Real-time |

#### 3.2.2 Product-Level Inventory

```go
type ProductInventory struct {
    ID              int64     `gorm:"column:id;primaryKey"`
    TenantID        int64     `gorm:"column:tenant_id;not null;index"`
    ProductID       int64     `gorm:"column:product_id;not null;uniqueIndex:idx_product"`
    AvailableStock  int       `gorm:"column:available_stock;not null;default:0"`
    LockedStock     int       `gorm:"column:locked_stock;not null;default:0"`
    TotalStock      int       `gorm:"column:total_stock;not null;default:0"`  // Computed: available + locked
    CreatedAt       time.Time `gorm:"column:created_at;not null"`
    UpdatedAt       time.Time `gorm:"column:updated_at;not null"`
}
```

#### 3.2.3 SKU-Level Inventory

```go
type SKUInventory struct {
    ID              int64     `gorm:"column:id;primaryKey"`
    TenantID        int64     `gorm:"column:tenant_id;not null;index"`
    SKU             string    `gorm:"column:sku;not null;uniqueIndex:idx_tenant_sku"`
    ProductID       int64     `gorm:"column:product_id;not null;index"`
    AvailableStock  int       `gorm:"column:available_stock;not null;default:0"`
    LockedStock     int       `gorm:"column:locked_stock;not null;default:0"`
    TotalStock      int       `gorm:"column:total_stock;not null;default:0"`

    // Safety stock threshold
    SafetyStock     int       `gorm:"column:safety_stock;default:0"`

    // Pre-sale mode
    PreSaleEnabled  bool      `gorm:"column:presale_enabled;default:false"`

    CreatedAt       time.Time `gorm:"column:created_at;not null"`
    UpdatedAt       time.Time `gorm:"column:updated_at;not null"`
}
```

#### 3.2.4 Warehouse-Level Inventory

```go
type WarehouseInventory struct {
    ID              int64     `gorm:"column:id;primaryKey"`
    TenantID        int64     `gorm:"column:tenant_id;not null;index"`
    SKU             string    `gorm:"column:sku;not null;index"`
    WarehouseID     int64     `gorm:"column:warehouse_id;not null;index"`
    AvailableStock  int       `gorm:"column:available_stock;not null;default:0"`
    LockedStock     int       `gorm:"column:locked_stock;not null;default:0"`
    CreatedAt       time.Time `gorm:"column:created_at;not null"`
    UpdatedAt       time.Time `gorm:"column:updated_at;not null"`

    // Unique constraint: (tenant_id, sku, warehouse_id)
}
```

#### 3.2.5 Inventory Log

```go
type InventoryLog struct {
    ID              int64     `gorm:"column:id;primaryKey"`
    TenantID        int64     `gorm:"column:tenant_id;not null;index"`
    SKU             string    `gorm:"column:sku;not null;index"`
    ProductID       int64     `gorm:"column:product_id;not null;index"`
    WarehouseID     int64     `gorm:"column:warehouse_id;default:0"`  // 0 = default/total

    // Change details
    ChangeType      string    `gorm:"column:change_type;type:varchar(30);not null"`  // manual, order, return, adjustment
    ChangeQuantity  int       `gorm:"column:change_quantity;not null"`               // Positive = increase, negative = decrease
    BeforeStock     int       `gorm:"column:before_stock;not null"`
    AfterStock      int       `gorm:"column:after_stock;not null"`

    // Reference
    OrderNo         string    `gorm:"column:order_no;type:varchar(50)"`
    Remark          string    `gorm:"column:remark;type:varchar(500)"`

    // Audit
    OperatorID      int64     `gorm:"column:operator_id;not null"`
    CreatedAt       time.Time `gorm:"column:created_at;not null"`
}
```

### 3.3 Business Rules

#### 3.3.1 Product vs SKU Inventory

| Scenario | Behavior |
|----------|----------|
| Single SKU product | SKU inventory = product inventory |
| Multi-SKU product | Product inventory = sum of all SKU inventories (auto-computed) |
| Inventory deduction | Always at SKU level |

#### 3.3.2 Stock Operations

| Operation | Available | Locked | Description |
|-----------|-----------|--------|-------------|
| Create product | Initialize | 0 | Set initial stock |
| Manual adjust | +/- N | - | Direct adjustment |
| Lock stock (order) | -N | +N | Move from available to locked |
| Deduct stock (payment) | - | -N | Deduct from locked |
| Restore stock (cancel) | +N | -N | Move from locked to available |
| Unlock expired | +N | -N | Order timeout, restore stock |

#### 3.3.3 Deduction Timing Configuration

```go
type InventoryConfig struct {
    LockTiming       string // "on_order" or "on_payment"
    // on_order: Lock when order created
    // on_payment: Lock/deduct when payment successful
}
```

**Recommended flow (on_payment):**
1. Order created → No stock change
2. Payment successful → Deduct available stock directly
3. Order cancelled/timeout → No stock change (was never deducted)

**Alternative flow (on_order):**
1. Order created → Lock stock (available -N, locked +N)
2. Payment successful → Deduct locked stock (locked -N)
3. Order cancelled/timeout → Restore stock (available +N, locked -N)

#### 3.3.4 Sold-Out Logic

| Rule | Description |
|------|-------------|
| Auto off-sale | When available stock = 0, product auto-changes to "off_sale" status |
| Manual restore | After restocking, merchant manually puts product on sale |
| Pre-sale mode | When enabled, allows orders even with 0/negative stock |

#### 3.3.5 Safety Stock Alert

| Rule | Description |
|------|-------------|
| Threshold | SKU-level `safety_stock` field |
| Trigger | Available stock < safety_stock |
| Notification | In-app notification + email |
| Frequency | Once per day per SKU |

### 3.4 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/products/:id/inventory` | Get product inventory |
| PUT | `/api/v1/products/:id/inventory` | Update product inventory |
| GET | `/api/v1/skus/:sku/inventory` | Get SKU inventory |
| PUT | `/api/v1/skus/:sku/inventory` | Update SKU inventory |
| POST | `/api/v1/inventory/batch-update` | Batch update inventory |
| POST | `/api/v1/inventory/adjust` | Manual inventory adjustment |
| GET | `/api/v1/inventory/logs` | Get inventory logs |
| GET | `/api/v1/inventory/export` | Export inventory report |
| GET | `/api/v1/inventory/low-stock` | Get low-stock products |
| PUT | `/api/v1/skus/:sku/safety-stock` | Set safety stock threshold |
| PUT | `/api/v1/skus/:sku/presale` | Toggle pre-sale mode |
| GET | `/api/v1/warehouses/:id/inventory` | Get warehouse inventory |
| PUT | `/api/v1/warehouses/:id/inventory` | Set warehouse inventory |

### 3.5 Database Schema

```sql
-- Product-level inventory (computed from SKU inventory)
CREATE TABLE product_inventories (
    id              BIGINT      PRIMARY KEY,
    tenant_id       BIGINT      NOT NULL,
    product_id      BIGINT      NOT NULL,
    available_stock INT         NOT NULL DEFAULT 0,
    locked_stock    INT         NOT NULL DEFAULT 0,
    total_stock     INT         NOT NULL DEFAULT 0,
    created_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_id (tenant_id),
    UNIQUE INDEX idx_product_id (product_id)
);

-- SKU-level inventory
CREATE TABLE sku_inventories (
    id               BIGINT      PRIMARY KEY,
    tenant_id        BIGINT      NOT NULL,
    sku              VARCHAR(100) NOT NULL,
    product_id       BIGINT      NOT NULL,
    available_stock  INT         NOT NULL DEFAULT 0,
    locked_stock     INT         NOT NULL DEFAULT 0,
    total_stock      INT         NOT NULL DEFAULT 0,
    safety_stock     INT         NOT NULL DEFAULT 0,
    presale_enabled  BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_product_id (product_id),
    UNIQUE INDEX idx_tenant_sku (tenant_id, sku)
);

-- Warehouse-level inventory
CREATE TABLE warehouse_inventories (
    id              BIGINT      PRIMARY KEY,
    tenant_id       BIGINT      NOT NULL,
    sku             VARCHAR(100) NOT NULL,
    warehouse_id    BIGINT      NOT NULL,
    available_stock INT         NOT NULL DEFAULT 0,
    locked_stock    INT         NOT NULL DEFAULT 0,
    created_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_sku (sku),
    INDEX idx_warehouse_id (warehouse_id),
    UNIQUE INDEX idx_tenant_sku_warehouse (tenant_id, sku, warehouse_id)
);

-- Inventory change log
CREATE TABLE inventory_logs (
    id              BIGINT          PRIMARY KEY,
    tenant_id       BIGINT          NOT NULL,
    sku             VARCHAR(100)    NOT NULL,
    product_id      BIGINT          NOT NULL,
    warehouse_id    BIGINT          NOT NULL DEFAULT 0,
    change_type     VARCHAR(30)     NOT NULL,  -- manual, order, return, adjustment
    change_quantity INT             NOT NULL,
    before_stock    INT             NOT NULL,
    after_stock     INT             NOT NULL,
    order_no        VARCHAR(50),
    remark          VARCHAR(500),
    operator_id     BIGINT          NOT NULL,
    created_at      TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_sku (sku),
    INDEX idx_product_id (product_id),
    INDEX idx_created_at (created_at)
);
```

---

## 4. Domain Service Interactions

### 4.1 Category Service

```go
type CategoryService interface {
    // CRUD
    Create(ctx context.Context, req *CreateCategoryRequest) (*Category, error)
    Update(ctx context.Context, id int64, req *UpdateCategoryRequest) (*Category, error)
    Delete(ctx context.Context, id int64) error
    GetByID(ctx context.Context, id int64) (*Category, error)

    // Tree operations
    GetTree(ctx context.Context) ([]*CategoryNode, error)
    Move(ctx context.Context, id int64, newParentID int64) error
    UpdateSort(ctx context.Context, sorts []CategorySort) error

    // Status
    Enable(ctx context.Context, id int64) error
    Disable(ctx context.Context, id int64) error

    // Product migration
    MigrateProducts(ctx context.Context, fromCategoryID, toCategoryID int64) error
    GetProductCount(ctx context.Context, categoryID int64) (int64, error)

    // Market visibility
    SetMarketVisibility(ctx context.Context, categoryID int64, marketID int64, visible bool) error
    GetMarketVisibility(ctx context.Context, categoryID int64) ([]*CategoryMarket, error)

    // Helpers
    IsLeafCategory(ctx context.Context, id int64) (bool, error)
    GetChildren(ctx context.Context, parentID int64) ([]*Category, error)
    ValidateUniqueName(ctx context.Context, parentID int64, name string, excludeID int64) error
}
```

### 4.2 Brand Service

```go
type BrandService interface {
    // CRUD
    Create(ctx context.Context, req *CreateBrandRequest) (*Brand, error)
    Update(ctx context.Context, id int64, req *UpdateBrandRequest) (*Brand, error)
    Delete(ctx context.Context, id int64) error
    GetByID(ctx context.Context, id int64) (*Brand, error)
    List(ctx context.Context, query *BrandQuery) ([]*Brand, int64, error)

    // Status
    Enable(ctx context.Context, id int64) error
    Disable(ctx context.Context, id int64) error

    // Brand page
    TogglePage(ctx context.Context, id int64, enabled bool) error

    // Market visibility
    SetMarketVisibility(ctx context.Context, brandID int64, marketID int64, visible bool) error
    GetMarketVisibility(ctx context.Context, brandID int64) ([]*BrandMarket, error)

    // Statistics
    GetProductCount(ctx context.Context, brandID int64) (int64, error)
}
```

### 4.3 Inventory Service

```go
type InventoryService interface {
    // Query
    GetProductInventory(ctx context.Context, productID int64) (*ProductInventory, error)
    GetSKUInventory(ctx context.Context, sku string) (*SKUInventory, error)
    GetWarehouseInventory(ctx context.Context, sku string, warehouseID int64) (*WarehouseInventory, error)

    // Update
    SetStock(ctx context.Context, sku string, availableStock int, operatorID int64, remark string) error
    AdjustStock(ctx context.Context, sku string, delta int, operatorID int64, remark string) error
    BatchUpdate(ctx context.Context, updates []StockUpdate, operatorID int64) error

    // Stock operations (called by order system)
    LockStock(ctx context.Context, sku string, quantity int, orderNo string) error
    DeductStock(ctx context.Context, sku string, quantity int, orderNo string) error
    RestoreStock(ctx context.Context, sku string, quantity int, orderNo string) error

    // Warehouse allocation
    SetWarehouseStock(ctx context.Context, sku string, warehouseID int64, stock int, operatorID int64) error

    // Settings
    SetSafetyStock(ctx context.Context, sku string, threshold int) error
    TogglePreSale(ctx context.Context, sku string, enabled bool) error

    // Logs
    GetLogs(ctx context.Context, query *InventoryLogQuery) ([]*InventoryLog, int64, error)

    // Export
    Export(ctx context.Context, query *InventoryExportQuery) ([]byte, error)

    // Alerts
    GetLowStockSKUs(ctx context.Context) ([]*SKUInventory, error)
    CheckAndNotifyLowStock(ctx context.Context) error
}
```

---

## 5. Frontend Implementation Notes

### 5.1 Category Management UI

- **Category tree**: Drag-and-drop tree component
- **Cascader selector**: For product edit page category selection
- **Market visibility**: Checkbox group for market selection
- **Product migration**: Modal with category selector

### 5.2 Brand Management UI

- **Brand list**: Table with logo, name, status, product count
- **Brand selector**: Searchable dropdown for product edit
- **Market visibility**: Checkbox group for market selection
- **Brand page toggle**: Switch component

### 5.3 Inventory Management UI

- **Product inventory**: Display total + per-SKU breakdown
- **SKU inventory table**: Editable cells for stock adjustment
- **Inventory log**: Timeline/table with filters
- **Low stock alerts**: Badge/notification in header
- **Batch operations**: Modal with CSV preview
- **Export**: Download button for inventory report

---

## 6. Migration Plan

### Phase 1: Database Schema
1. Create new tables: `categories`, `category_markets`, `brands`, `brand_markets`, inventory tables
2. Add `brand_id` to `products` table
3. Update `category_id` foreign key constraints

### Phase 2: Backend Implementation
1. Implement Category domain, repository, service, handlers
2. Implement Brand domain, repository, service, handlers
3. Implement Inventory domain, repository, service, handlers
4. Update Product service to integrate category and brand

### Phase 3: Frontend Implementation
1. Category management pages (list, tree, edit)
2. Brand management pages (list, edit)
3. Inventory management tab in product detail
4. Category/Brand selectors in product edit form

### Phase 4: Integration
1. Connect frontend to backend APIs
2. Implement inventory deduction in order flow
3. Implement low stock notifications

---

## 7. Open Questions

| # | Question | Status |
|---|----------|--------|
| 1 | Should category code be auto-generated or manual? | TBD |
| 2 | Max inventory log retention period default? | Suggest: 90 days |
| 3 | Low stock notification delivery time? | Suggest: Daily at 9:00 AM |
| 4 | Pre-sale mode auto-disable after restock? | Suggest: Manual only |