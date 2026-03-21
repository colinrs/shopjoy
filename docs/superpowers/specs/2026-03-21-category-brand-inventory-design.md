# Category, Brand, and Inventory Management Design

**Date:** 2026-03-21
**Status:** Draft
**Scope:** Admin Backend (Merchant Management)

---

## 0. Existing System Context

This design extends the existing ShopJoy system. Key existing components:

### 0.1 Existing Tables (Already Implemented)

| Table | Description | Key Fields |
|-------|-------------|------------|
| `categories` | Category tree | `id`, `tenant_id`, `parent_id`, `name`, `code`, `level`, `sort`, `status` |
| `brands` | Brand info | `id`, `tenant_id`, `name`, `logo`, `description`, `website`, `status` |
| `products` | Product entity | `id`, `sku`, `name`, `price`, `stock`, `status`, `category_id`, `brand` (string) |
| `skus` | SKU variants | `id`, `product_id`, `code`, `price_amount`, `stock`, `attributes`, `status` |
| `markets` | Market config | `id`, `tenant_id`, `code`, `name`, `currency`, `is_active` |
| `product_markets` | Product-market relation | `id`, `product_id`, `market_id`, `is_enabled`, `price` |

### 0.2 Key Differences from Current Implementation

| Area | Current | Proposed | Migration Required |
|------|---------|----------|-------------------|
| Product.brand | `VARCHAR` string | `BIGINT brand_id` FK | Yes - data migration |
| Product.stock | Single `INT` field | Separate inventory tables | Yes - data migration |
| SKU.stock | Single `INT` field | Extended with locked_stock, safety_stock | Yes - schema change |
| Category SEO | Not present | Add `seo_title`, `seo_description` | Yes - schema change |
| Brand compliance | Not present | Add trademark fields | Yes - schema change |
| Soft delete | Products use `status = 3` | Categories/Brands use `deleted_at` | No - keep both patterns |

### 0.3 Design Principles

1. **Extend existing tables** rather than replace where possible
2. **Backward compatibility** for existing product/brand string field during migration
3. **Tenant isolation** enforced at repository layer (existing pattern)
4. **Unix timestamp storage** for time fields (existing pattern: BIGINT)

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

### 1.5 Database Schema Changes

**Note:** `categories` table already exists. Use ALTER statements to add new fields.

```sql
-- Add SEO fields to existing categories table
ALTER TABLE `categories`
    ADD COLUMN `seo_title` VARCHAR(200) DEFAULT '' COMMENT 'SEO标题' AFTER `image`,
    ADD COLUMN `seo_description` VARCHAR(500) DEFAULT '' COMMENT 'SEO描述' AFTER `seo_title`;

-- Create new table for category market visibility
CREATE TABLE IF NOT EXISTS `category_markets` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL,
    `category_id` BIGINT NOT NULL,
    `market_id` BIGINT NOT NULL,
    `is_visible` TINYINT NOT NULL DEFAULT 1 COMMENT '是否可见',
    `created_at` BIGINT NOT NULL DEFAULT 0,
    `updated_at` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_category_id` (`category_id`),
    UNIQUE KEY `idx_tenant_category_market` (`tenant_id`, `category_id`, `market_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分类市场可见性';
```

### 1.6 Market Entity Reference

**Market** is already defined in `sql/market.sql`. Key fields:

| Field | Type | Description |
|-------|------|-------------|
| id | BIGINT | Primary key |
| tenant_id | BIGINT | Tenant ID |
| code | VARCHAR(10) | Market code (US, UK, DE, FR, AU, CN) |
| name | VARCHAR(100) | Market name |
| currency | VARCHAR(10) | Currency (USD, GBP, EUR, CNY) |
| is_active | TINYINT | Active status |

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

### 2.5 Database Schema Changes

**Note:** `brands` table already exists. Use ALTER statements to add new fields.

```sql
-- Add compliance and brand page fields to existing brands table
ALTER TABLE `brands`
    ADD COLUMN `enable_page` TINYINT NOT NULL DEFAULT 0 COMMENT '是否启用品牌专区' AFTER `sort`,
    ADD COLUMN `trademark_number` VARCHAR(100) DEFAULT '' COMMENT '商标号' AFTER `enable_page`,
    ADD COLUMN `trademark_country` VARCHAR(10) DEFAULT '' COMMENT '商标注册国家' AFTER `trademark_number`;

-- Create new table for brand market visibility
CREATE TABLE IF NOT EXISTS `brand_markets` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL,
    `brand_id` BIGINT NOT NULL,
    `market_id` BIGINT NOT NULL,
    `is_visible` TINYINT NOT NULL DEFAULT 1 COMMENT '是否可见',
    `created_at` BIGINT NOT NULL DEFAULT 0,
    `updated_at` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_brand_id` (`brand_id`),
    UNIQUE KEY `idx_tenant_brand_market` (`tenant_id`, `brand_id`, `market_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌市场可见性';

-- Add brand_id to products table (migrate from brand string)
ALTER TABLE `products`
    ADD COLUMN `brand_id` BIGINT NULL COMMENT '品牌ID' AFTER `brand`,
    ADD INDEX `idx_brand_id` (`brand_id`);
```

### 2.6 Brand Migration Strategy

**Migrating `brand` string field to `brand_id` FK:**

1. **Phase 1 - Schema**: Add `brand_id` column (nullable)
2. **Phase 2 - Data Migration**:
   - For each distinct `brand` value in products, create Brand record
   - Update products to set `brand_id` from newly created brands
3. **Phase 3 - Cleanup**: Remove `brand` string column (after verification)

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

**Note:** `skus` table already exists. Extend with inventory fields.

```go
// Extend existing SKU entity with inventory fields
type SKUInventory struct {
    ID              int64     `gorm:"column:id;primaryKey"`
    TenantID        int64     `gorm:"column:tenant_id;not null;index"`
    ProductID       int64     `gorm:"column:product_id;not null;index"`
    Code            string    `gorm:"column:code;not null;uniqueIndex"`  // Existing SKU code

    // Existing price fields
    PriceAmount     int64     `gorm:"column:price_amount"`
    PriceCurrency   string    `gorm:"column:price_currency"`

    // Extended inventory fields (new)
    AvailableStock  int       `gorm:"column:available_stock;not null;default:0"`
    LockedStock     int       `gorm:"column:locked_stock;not null;default:0"`
    SafetyStock     int       `gorm:"column:safety_stock;default:0"`
    PreSaleEnabled  bool      `gorm:"column:presale_enabled;default:false"`

    // Existing fields
    Attributes      JSON      `gorm:"column:attributes;type:json"`
    Status          int8      `gorm:"column:status;not null;default:1"`

    CreatedAt       time.Time `gorm:"column:created_at;not null"`
    UpdatedAt       time.Time `gorm:"column:updated_at;not null"`
}
```

#### 3.2.4 Warehouse Entity

```go
type Warehouse struct {
    ID          int64     `gorm:"column:id;primaryKey"`
    TenantID    int64     `gorm:"column:tenant_id;not null;index"`
    Code        string    `gorm:"column:code;type:varchar(50);not null;uniqueIndex:idx_tenant_code"`
    Name        string    `gorm:"column:name;type:varchar(100);not null"`
    Country     string    `gorm:"column:country;type:varchar(10)"`  // ISO country code
    Address     string    `gorm:"column:address;type:varchar(500)"`
    IsDefault   bool      `gorm:"column:is_default;default:false"`
    Status      int8      `gorm:"column:status;not null;default:1"`  // 1=active, 0=inactive
    CreatedAt   time.Time `gorm:"column:created_at;not null"`
    UpdatedAt   time.Time `gorm:"column:updated_at;not null"`
    DeletedAt   *time.Time `gorm:"column:deleted_at"`
}
```

#### 3.2.5 Warehouse-Level Inventory

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

### 3.5 Database Schema Changes

**Note:** `skus` table already exists. Use ALTER statements to extend inventory fields.

```sql
-- Add inventory fields to existing skus table
ALTER TABLE `skus`
    ADD COLUMN `available_stock` INT NOT NULL DEFAULT 0 COMMENT '可用库存' AFTER `stock`,
    ADD COLUMN `locked_stock` INT NOT NULL DEFAULT 0 COMMENT '锁定库存' AFTER `available_stock`,
    ADD COLUMN `safety_stock` INT NOT NULL DEFAULT 0 COMMENT '安全库存阈值' AFTER `locked_stock`,
    ADD COLUMN `presale_enabled` TINYINT NOT NULL DEFAULT 0 COMMENT '是否开启预售' AFTER `safety_stock`;

-- Migrate existing stock to available_stock
UPDATE `skus` SET `available_stock` = `stock` WHERE `available_stock` = 0;

-- Create warehouses table
CREATE TABLE IF NOT EXISTS `warehouses` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL,
    `code` VARCHAR(50) NOT NULL COMMENT '仓库代码',
    `name` VARCHAR(100) NOT NULL COMMENT '仓库名称',
    `country` VARCHAR(10) DEFAULT '' COMMENT '所在国家',
    `address` VARCHAR(500) DEFAULT '' COMMENT '详细地址',
    `is_default` TINYINT NOT NULL DEFAULT 0 COMMENT '是否默认仓库',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
    `created_at` BIGINT NOT NULL DEFAULT 0,
    `updated_at` BIGINT NOT NULL DEFAULT 0,
    `deleted_at` BIGINT DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    UNIQUE KEY `idx_tenant_code` (`tenant_id`, `code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='仓库表';

-- Create warehouse-level inventory table
CREATE TABLE IF NOT EXISTS `warehouse_inventories` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL,
    `sku_code` VARCHAR(100) NOT NULL COMMENT 'SKU代码',
    `warehouse_id` BIGINT NOT NULL,
    `available_stock` INT NOT NULL DEFAULT 0,
    `locked_stock` INT NOT NULL DEFAULT 0,
    `created_at` BIGINT NOT NULL DEFAULT 0,
    `updated_at` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_sku_code` (`sku_code`),
    KEY `idx_warehouse_id` (`warehouse_id`),
    UNIQUE KEY `idx_tenant_sku_warehouse` (`tenant_id`, `sku_code`, `warehouse_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT '仓库库存表';

-- Create inventory change log table
CREATE TABLE IF NOT EXISTS `inventory_logs` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL,
    `sku_code` VARCHAR(100) NOT NULL,
    `product_id` BIGINT NOT NULL,
    `warehouse_id` BIGINT NOT NULL DEFAULT 0 COMMENT '0=汇总',
    `change_type` VARCHAR(30) NOT NULL COMMENT 'manual, order, return, adjustment',
    `change_quantity` INT NOT NULL COMMENT '正数增加，负数减少',
    `before_stock` INT NOT NULL,
    `after_stock` INT NOT NULL,
    `order_no` VARCHAR(50) DEFAULT '',
    `remark` VARCHAR(500) DEFAULT '',
    `operator_id` BIGINT NOT NULL,
    `created_at` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_sku_code` (`sku_code`),
    KEY `idx_product_id` (`product_id`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存变更日志';
```

### 3.6 Concurrency Control

**Preventing overselling in high-concurrency scenarios:**

#### 3.6.1 Database-Level Locking

```go
// LockStock with row-level locking (FOR UPDATE)
func (r *inventoryRepo) LockStock(ctx context.Context, db *gorm.DB, skuCode string, quantity int, orderNo string) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // SELECT ... FOR UPDATE to acquire row lock
        var sku SKUInventory
        if err := tx.Set("gorm:query_option", "FOR UPDATE").
            Where("code = ? AND tenant_id = ?", skuCode, tenantID).
            First(&sku).Error; err != nil {
            return err
        }

        if sku.AvailableStock < quantity {
            return ErrInsufficientStock
        }

        // Update within transaction
        return tx.Model(&SKUInventory{}).
            Where("code = ?", skuCode).
            Updates(map[string]interface{}{
                "available_stock": gorm.Expr("available_stock - ?", quantity),
                "locked_stock":    gorm.Expr("locked_stock + ?", quantity),
                "updated_at":      time.Now().Unix(),
            }).Error
    })
}
```

#### 3.6.2 Optimistic Locking Alternative

```go
// Using version field for optimistic locking
type SKUInventory struct {
    // ... other fields
    Version     int       `gorm:"column:version;not null;default:0"`
}

func (r *inventoryRepo) LockStock(ctx context.Context, db *gorm.DB, skuCode string, quantity int, version int) error {
    result := db.Model(&SKUInventory{}).
        Where("code = ? AND version = ?", skuCode, version).
        Updates(map[string]interface{}{
            "available_stock": gorm.Expr("available_stock - ?", quantity),
            "locked_stock":    gorm.Expr("locked_stock + ?", quantity),
            "version":         gorm.Expr("version + 1"),
            "updated_at":      time.Now().Unix(),
        })

    if result.RowsAffected == 0 {
        return ErrConcurrentModification // Retry needed
    }
    return nil
}
```

#### 3.6.3 Redis Distributed Lock (Optional)

For extremely high concurrency, use Redis distributed lock before database operations:

```go
func (s *inventoryService) LockStock(ctx context.Context, skuCode string, quantity int) error {
    lockKey := fmt.Sprintf("inventory:lock:%s", skuCode)

    // Acquire distributed lock
    if !s.redis.SetNX(ctx, lockKey, "1", 5*time.Second) {
        return ErrLockConflict
    }
    defer s.redis.Del(ctx, lockKey)

    // Proceed with database operation
    return s.repo.LockStock(ctx, s.db, skuCode, quantity)
}
```

### 3.7 Inventory Configuration Storage

**Tenant-level inventory settings stored in `tenant_settings` table:**

```go
type InventorySettings struct {
    LockTiming          string `json:"lock_timing"`           // "on_order" or "on_payment"
    LogRetentionDays    int    `json:"log_retention_days"`    // Default: 90
    LowStockNotifyTime  string `json:"low_stock_notify_time"` // Default: "09:00"
    EnableLowStockAlert bool   `json:"enable_low_stock_alert"` // Default: true
}
```

Stored as JSON in existing `tenant_settings` table under key `inventory`.

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

### Phase 1: Database Schema Updates
1. Add SEO fields to `categories` table
2. Add compliance fields to `brands` table
3. Add `brand_id` to `products` table
4. Add inventory fields to `skus` table
5. Create `category_markets`, `brand_markets`, `warehouses`, `warehouse_inventories`, `inventory_logs` tables

### Phase 2: Data Migration
1. Migrate `products.brand` string to `brand_id`:
   - Extract distinct brand names from products
   - Create Brand records for each unique brand
   - Update products to set `brand_id`
2. Migrate `skus.stock` to `skus.available_stock`
3. Migrate `products.stock` to be computed from SKU inventories

### Phase 3: Backend Implementation
1. Implement Category domain, repository, service, handlers
2. Implement Brand domain, repository, service, handlers
3. Implement Inventory domain, repository, service, handlers
4. Update Product service to integrate category and brand
5. Add inventory operations to order flow

### Phase 4: Frontend Implementation
1. Category management pages (list, tree, edit)
2. Brand management pages (list, edit)
3. Inventory management tab in product detail
4. Category/Brand selectors in product edit form

### Phase 5: Integration
1. Connect frontend to backend APIs
2. Implement inventory deduction in order flow
3. Implement low stock notifications (asynq task)

---

## 7. API Request/Response DTOs

### 7.1 Category APIs

```go
// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
    Name           string `json:"name,optional"`                // 分类名称
    ParentID       int64  `json:"parent_id,optional"`           // 父分类ID
    Code           string `json:"code,optional"`                // 分类代码
    Icon           string `json:"icon,optional"`                // 图标URL
    Image          string `json:"image,optional"`               // 图片URL
    SeoTitle       string `json:"seo_title,optional"`           // SEO标题
    SeoDescription string `json:"seo_description,optional"`     // SEO描述
    Sort           int    `json:"sort,optional"`                // 排序
}

// CategoryResponse 分类响应
type CategoryResponse struct {
    ID             int64  `json:"id"`
    ParentID       int64  `json:"parent_id"`
    Name           string `json:"name"`
    Code           string `json:"code"`
    Level          int    `json:"level"`
    Sort           int    `json:"sort"`
    Icon           string `json:"icon"`
    Image          string `json:"image"`
    SeoTitle       string `json:"seo_title"`
    SeoDescription string `json:"seo_description"`
    Status         int8   `json:"status"`
    ProductCount   int64  `json:"product_count"`
    CreatedAt      string `json:"created_at"`
    UpdatedAt      string `json:"updated_at"`
}

// CategoryTreeResponse 分类树响应
type CategoryTreeResponse struct {
    CategoryResponse
    Children []*CategoryTreeResponse `json:"children"`
}
```

### 7.2 Brand APIs

```go
// CreateBrandRequest 创建品牌请求
type CreateBrandRequest struct {
    Name             string `json:"name"`                          // 品牌名称
    Logo             string `json:"logo,optional"`                 // Logo URL
    Description      string `json:"description,optional"`          // 描述
    Website          string `json:"website,optional"`              // 官网
    TrademarkNumber  string `json:"trademark_number,optional"`     // 商标号
    TrademarkCountry string `json:"trademark_country,optional"`    // 商标注册国家
    EnablePage       bool   `json:"enable_page,optional"`          // 启用品牌专区
    Sort             int    `json:"sort,optional"`                 // 排序
}

// BrandResponse 品牌响应
type BrandResponse struct {
    ID               int64  `json:"id"`
    Name             string `json:"name"`
    Logo             string `json:"logo"`
    Description      string `json:"description"`
    Website          string `json:"website"`
    TrademarkNumber  string `json:"trademark_number"`
    TrademarkCountry string `json:"trademark_country"`
    EnablePage       bool   `json:"enable_page"`
    Sort             int    `json:"sort"`
    Status           int8   `json:"status"`
    ProductCount     int64  `json:"product_count"`
    CreatedAt        string `json:"created_at"`
    UpdatedAt        string `json:"updated_at"`
}
```

### 7.3 Inventory APIs

```go
// UpdateStockRequest 更新库存请求
type UpdateStockRequest struct {
    AvailableStock int    `json:"available_stock"`    // 可用库存
    Remark         string `json:"remark,optional"`    // 备注
}

// AdjustStockRequest 库存调整请求
type AdjustStockRequest struct {
    Quantity int    `json:"quantity"`         // 调整数量（正数增加，负数减少）
    Remark   string `json:"remark,optional"`  // 备注
}

// SKUInventoryResponse SKU库存响应
type SKUInventoryResponse struct {
    SKUCode         string `json:"sku_code"`
    ProductID       int64  `json:"product_id"`
    ProductName     string `json:"product_name"`
    AvailableStock  int    `json:"available_stock"`
    LockedStock     int    `json:"locked_stock"`
    TotalStock      int    `json:"total_stock"`
    SafetyStock     int    `json:"safety_stock"`
    PreSaleEnabled  bool   `json:"presale_enabled"`
    Attributes      string `json:"attributes"`
    Status          int8   `json:"status"`
}

// InventoryLogResponse 库存日志响应
type InventoryLogResponse struct {
    ID             int64  `json:"id"`
    SKUCode        string `json:"sku_code"`
    ChangeType     string `json:"change_type"`
    ChangeQuantity int    `json:"change_quantity"`
    BeforeStock    int    `json:"before_stock"`
    AfterStock     int    `json:"after_stock"`
    OrderNo        string `json:"order_no"`
    Remark         string `json:"remark"`
    OperatorName   string `json:"operator_name"`
    CreatedAt      string `json:"created_at"`
}
```

---

## 8. Low Stock Notification Integration

### 8.1 Asynq Task Definition

```go
const TaskLowStockAlert = "inventory:low_stock_alert"

type LowStockAlertPayload struct {
    TenantID int64
}

func (s *inventoryService) CheckAndNotifyLowStock(ctx context.Context, tenantID int64) error {
    // 1. Get all SKUs where available_stock < safety_stock
    lowStockSKUs, err := s.repo.GetLowStockSKUs(ctx, tenantID)
    if err != nil {
        return err
    }

    if len(lowStockSKUs) == 0 {
        return nil
    }

    // 2. Get tenant settings for notification
    settings, err := s.settingsService.GetInventorySettings(ctx, tenantID)
    if err != nil {
        return err
    }

    // 3. Send notification (in-app + email)
    return s.notifyService.SendLowStockAlert(ctx, tenantID, lowStockSKUs, settings)
}
```

### 8.2 Scheduled Task

```go
// Register in main.go
scheduler := asynq.NewScheduler(...)
scheduler.Register("0 9 * * *", asynq.NewTask(TaskLowStockAlert, payload))
// Runs daily at 9:00 AM tenant local time
```

---

## 9. Resolved Questions

| # | Question | Resolution |
|---|----------|------------|
| 1 | Category code auto-generation | Manual entry, optional field. Auto-generate from name if empty. |
| 2 | Inventory log retention | Default 90 days, configurable in tenant settings. |
| 3 | Low stock notification time | Default 09:00, configurable in tenant settings. |
| 4 | Pre-sale mode auto-disable | Manual only. Merchant must disable after restock. |
| 5 | Soft delete pattern | Products use status=deleted (existing), Categories/Brands use deleted_at (new). |