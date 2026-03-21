# Brand Management Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement Brand management with brand page toggle, trademark compliance fields, and market visibility for the ShopJoy admin backend.

**Architecture:** Extend existing Brand entity in `admin/internal/domain/product/category.go`, add BrandMarket entity, implement repository, create go-zero API handlers, and build Vue frontend for brand management.

**Tech Stack:** Go 1.21+, go-zero framework, GORM, MySQL, Vue 3 + TypeScript + Element Plus

**Spec:** `docs/superpowers/specs/2026-03-21-category-brand-inventory-design.md` (Section 2)

---

## File Structure

### New Files

```
admin/
├── desc/brand.api                            # API definitions
├── internal/
│   ├── infrastructure/persistence/
│   │   ├── brand_repository.go               # Repository implementation
│   │   └── brand_market_repository.go        # Market visibility repo
│   ├── handler/brands/                       # Auto-generated
│   └── logic/brands/                         # Business logic
sql/
└── 20260321_brand.sql                        # Migration script
shop-admin/src/
├── api/brand.ts                              # API client
└── views/brands/index.vue                    # Management page
```

### Modified Files

```
admin/internal/domain/product/category.go     # Extend Brand, add BrandMarket
admin/internal/svc/service_context.go         # Add BrandRepo, BrandMarketRepo
admin/desc/admin.api                          # Import brand.api
```

---

## Task 1: Database Schema

**Files:**
- Create: `sql/20260321_brand.sql`

- [ ] **Step 1: Create migration file**

```sql
-- sql/20260321_brand.sql
-- Brand management schema changes

-- Add compliance and brand page fields to brands table
ALTER TABLE `brands`
    ADD COLUMN `enable_page` TINYINT NOT NULL DEFAULT 0 COMMENT '是否启用品牌专区' AFTER `sort`,
    ADD COLUMN `trademark_number` VARCHAR(100) DEFAULT '' COMMENT '商标号' AFTER `enable_page`,
    ADD COLUMN `trademark_country` VARCHAR(10) DEFAULT '' COMMENT '商标注册国家' AFTER `trademark_number`;

-- Create brand_markets table for market visibility
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
    KEY `idx_market_id` (`market_id`),
    UNIQUE KEY `idx_tenant_brand_market` (`tenant_id`, `brand_id`, `market_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌市场可见性';

-- Add brand_id to products table (migration from brand string)
ALTER TABLE `products`
    ADD COLUMN `brand_id` BIGINT NULL COMMENT '品牌ID' AFTER `brand`,
    ADD INDEX `idx_brand_id` (`brand_id`);
```

- [ ] **Step 2: Run migration**

Run: `mysql -u root -p shopjoy < sql/20260321_brand.sql`

Expected: Query OK messages for ALTER and CREATE

- [ ] **Step 3: Commit**

```bash
git add sql/20260321_brand.sql
git commit -m "feat(db): add brand compliance fields, brand_markets table, and brand_id to products"
```

---

## Task 2: Update Brand Domain Entity

**Files:**
- Modify: `admin/internal/domain/product/category.go`

- [ ] **Step 1: Read current file**

Run: Read `admin/internal/domain/product/category.go`

- [ ] **Step 2: Add fields to Brand struct**

Update the Brand struct to include new fields:

```go
type Brand struct {
	ID               int64
	TenantID         shared.TenantID
	Name             string
	Logo             string
	Description      string
	Website          string
	Sort             int
	EnablePage       bool              // Enable brand page
	TrademarkNumber  string            // Trademark registration number
	TrademarkCountry string            // ISO country code for trademark
	Status           shared.Status
	Audit            shared.AuditInfo `gorm:"embedded"`
}
```

- [ ] **Step 3: Add BrandMarket entity after Brand struct**

```go
// BrandMarket represents brand visibility in specific markets
type BrandMarket struct {
	ID         int64
	TenantID   shared.TenantID
	BrandID    int64
	MarketID   int64
	IsVisible  bool
	Audit      shared.AuditInfo `gorm:"embedded"`
}

func (bm *BrandMarket) TableName() string {
	return "brand_markets"
}
```

- [ ] **Step 4: Add methods to Brand**

```go
func (b *Brand) Enable() {
	b.Status = shared.StatusEnabled
}

func (b *Brand) Disable() {
	b.Status = shared.StatusDisabled
}

func (b *Brand) IsEnabled() bool {
	return b.Status == shared.StatusEnabled
}

func (b *Brand) TogglePage(enabled bool) {
	b.EnablePage = enabled
}
```

- [ ] **Step 5: Update BrandRepository interface**

Add new methods to the existing interface:

```go
type BrandRepository interface {
	Create(ctx context.Context, db *gorm.DB, brand *Brand) error
	Update(ctx context.Context, db *gorm.DB, brand *Brand) error
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Brand, error)
	FindByName(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, name string) (*Brand, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query BrandQuery) ([]*Brand, int64, error)
	GetProductCount(ctx context.Context, db *gorm.DB, brandID int64) (int64, error)
}
```

- [ ] **Step 6: Add BrandMarketRepository interface**

```go
// BrandMarketRepository interface
type BrandMarketRepository interface {
	Create(ctx context.Context, db *gorm.DB, bm *BrandMarket) error
	Update(ctx context.Context, db *gorm.DB, bm *BrandMarket) error
	FindByBrand(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, brandID int64) ([]*BrandMarket, error)
	DeleteByBrand(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, brandID int64) error
	BatchCreate(ctx context.Context, db *gorm.DB, items []*BrandMarket) error
}
```

- [ ] **Step 7: Run build**

Run: `cd admin && make build`

Expected: Build succeeded

- [ ] **Step 8: Commit**

```bash
git add admin/internal/domain/product/category.go
git commit -m "feat(domain): extend Brand with compliance fields and BrandMarket entity"
```

---

## Task 3: Implement Brand Repository

**Files:**
- Create: `admin/internal/infrastructure/persistence/brand_repository.go`

- [ ] **Step 1: Create brand_repository.go**

```go
package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type brandRepo struct{}

func NewBrandRepository() product.BrandRepository {
	return &brandRepo{}
}

type brandModel struct {
	ID               int64  `gorm:"column:id;primaryKey"`
	TenantID         int64  `gorm:"column:tenant_id;not null;index"`
	Name             string `gorm:"column:name;type:varchar(100);not null"`
	Logo             string `gorm:"column:logo;type:varchar(500)"`
	Description      string `gorm:"column:description;type:text"`
	Website          string `gorm:"column:website;type:varchar(500)"`
	Sort             int    `gorm:"column:sort;default:0"`
	EnablePage       bool   `gorm:"column:enable_page;default:false"`
	TrademarkNumber  string `gorm:"column:trademark_number;type:varchar(100)"`
	TrademarkCountry string `gorm:"column:trademark_country;type:varchar(10)"`
	Status           int8   `gorm:"column:status;not null;default:1"`
	CreatedAt        int64  `gorm:"column:created_at;not null"`
	UpdatedAt        int64  `gorm:"column:updated_at;not null"`
	CreatedBy        int64  `gorm:"column:created_by"`
	UpdatedBy        int64  `gorm:"column:updated_by"`
	DeletedAt        *int64 `gorm:"column:deleted_at;index"`
}

func (brandModel) TableName() string {
	return "brands"
}

func (m *brandModel) toEntity() *product.Brand {
	return &product.Brand{
		ID:               m.ID,
		TenantID:         shared.TenantID(m.TenantID),
		Name:             m.Name,
		Logo:             m.Logo,
		Description:      m.Description,
		Website:          m.Website,
		Sort:             m.Sort,
		EnablePage:       m.EnablePage,
		TrademarkNumber:  m.TrademarkNumber,
		TrademarkCountry: m.TrademarkCountry,
		Status:           shared.Status(m.Status),
		Audit: shared.AuditInfo{
			CreatedAt: time.Unix(m.CreatedAt, 0),
			UpdatedAt: time.Unix(m.UpdatedAt, 0),
			CreatedBy: m.CreatedBy,
			UpdatedBy: m.UpdatedBy,
		},
	}
}

func fromBrandEntity(b *product.Brand) *brandModel {
	now := time.Now().Unix()
	createdAt := now
	updatedAt := now
	if !b.Audit.CreatedAt.IsZero() {
		createdAt = b.Audit.CreatedAt.Unix()
	}
	if !b.Audit.UpdatedAt.IsZero() {
		updatedAt = b.Audit.UpdatedAt.Unix()
	}
	return &brandModel{
		ID:               b.ID,
		TenantID:         b.TenantID.Int64(),
		Name:             b.Name,
		Logo:             b.Logo,
		Description:      b.Description,
		Website:          b.Website,
		Sort:             b.Sort,
		EnablePage:       b.EnablePage,
		TrademarkNumber:  b.TrademarkNumber,
		TrademarkCountry: b.TrademarkCountry,
		Status:           int8(b.Status),
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		CreatedBy:        b.Audit.CreatedBy,
		UpdatedBy:        b.Audit.UpdatedBy,
	}
}

func (r *brandRepo) Create(ctx context.Context, db *gorm.DB, b *product.Brand) error {
	model := fromBrandEntity(b)
	return db.WithContext(ctx).Create(model).Error
}

func (r *brandRepo) Update(ctx context.Context, db *gorm.DB, b *product.Brand) error {
	model := fromBrandEntity(b)
	return db.WithContext(ctx).Model(&brandModel{}).
		Where("id = ? AND tenant_id = ?", b.ID, b.TenantID.Int64()).
		Updates(map[string]interface{}{
			"name":              model.Name,
			"logo":              model.Logo,
			"description":       model.Description,
			"website":           model.Website,
			"sort":              model.Sort,
			"enable_page":       model.EnablePage,
			"trademark_number":  model.TrademarkNumber,
			"trademark_country": model.TrademarkCountry,
			"status":            model.Status,
			"updated_at":        model.UpdatedAt,
			"updated_by":        model.UpdatedBy,
		}).Error
}

func (r *brandRepo) Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	now := time.Now().Unix()
	return db.WithContext(ctx).Model(&brandModel{}).
		Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).
		Update("deleted_at", now).Error
}

func (r *brandRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*product.Brand, error) {
	var model brandModel
	err := db.WithContext(ctx).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", id, tenantID.Int64()).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *brandRepo) FindByName(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, name string) (*product.Brand, error) {
	var model brandModel
	err := db.WithContext(ctx).
		Where("name = ? AND tenant_id = ? AND deleted_at IS NULL", name, tenantID.Int64()).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *brandRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query product.BrandQuery) ([]*product.Brand, int64, error) {
	var models []brandModel
	var total int64

	tx := db.WithContext(ctx).Model(&brandModel{}).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID.Int64())

	if query.Name != "" {
		tx = tx.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Status != 0 {
		tx = tx.Where("status = ?", query.Status)
	}

	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := tx.Order("sort ASC, id DESC").
		Offset(int(offset)).
		Limit(int(query.PageSize)).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	brands := make([]*product.Brand, len(models))
	for i, m := range models {
		brands[i] = m.toEntity()
	}
	return brands, total, nil
}

func (r *brandRepo) GetProductCount(ctx context.Context, db *gorm.DB, brandID int64) (int64, error) {
	var count int64
	err := db.WithContext(ctx).Table("products").
		Where("brand_id = ?", brandID).
		Count(&count).Error
	return count, err
}
```

- [ ] **Step 2: Run build**

Run: `cd admin && make build`

Expected: Build succeeded

- [ ] **Step 3: Commit**

```bash
git add admin/internal/infrastructure/persistence/brand_repository.go
git commit -m "feat(infra): implement BrandRepository with compliance fields"
```

---

## Task 4: Implement BrandMarket Repository

**Files:**
- Create: `admin/internal/infrastructure/persistence/brand_market_repository.go`

- [ ] **Step 1: Create brand_market_repository.go**

```go
package persistence

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type brandMarketRepo struct{}

func NewBrandMarketRepository() product.BrandMarketRepository {
	return &brandMarketRepo{}
}

type brandMarketModel struct {
	ID         int64 `gorm:"column:id;primaryKey"`
	TenantID   int64 `gorm:"column:tenant_id;not null;index"`
	BrandID    int64 `gorm:"column:brand_id;not null;index"`
	MarketID   int64 `gorm:"column:market_id;not null;index"`
	IsVisible  bool  `gorm:"column:is_visible;not null;default:true"`
	CreatedAt  int64 `gorm:"column:created_at;not null"`
	UpdatedAt  int64 `gorm:"column:updated_at;not null"`
}

func (brandMarketModel) TableName() string {
	return "brand_markets"
}

func (m *brandMarketModel) toEntity() *product.BrandMarket {
	return &product.BrandMarket{
		ID:        m.ID,
		TenantID:  shared.TenantID(m.TenantID),
		BrandID:   m.BrandID,
		MarketID:  m.MarketID,
		IsVisible: m.IsVisible,
		Audit: shared.AuditInfo{
			CreatedAt: time.Unix(m.CreatedAt, 0),
			UpdatedAt: time.Unix(m.UpdatedAt, 0),
		},
	}
}

func fromBrandMarketEntity(bm *product.BrandMarket) *brandMarketModel {
	now := time.Now().Unix()
	createdAt := now
	updatedAt := now
	if !bm.Audit.CreatedAt.IsZero() {
		createdAt = bm.Audit.CreatedAt.Unix()
	}
	if !bm.Audit.UpdatedAt.IsZero() {
		updatedAt = bm.Audit.UpdatedAt.Unix()
	}
	return &brandMarketModel{
		ID:        bm.ID,
		TenantID:  bm.TenantID.Int64(),
		BrandID:   bm.BrandID,
		MarketID:  bm.MarketID,
		IsVisible: bm.IsVisible,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (r *brandMarketRepo) Create(ctx context.Context, db *gorm.DB, bm *product.BrandMarket) error {
	model := fromBrandMarketEntity(bm)
	return db.WithContext(ctx).Create(model).Error
}

func (r *brandMarketRepo) Update(ctx context.Context, db *gorm.DB, bm *product.BrandMarket) error {
	return db.WithContext(ctx).Model(&brandMarketModel{}).
		Where("id = ?", bm.ID).
		Updates(map[string]interface{}{
			"is_visible": bm.IsVisible,
			"updated_at": time.Now().Unix(),
		}).Error
}

func (r *brandMarketRepo) FindByBrand(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, brandID int64) ([]*product.BrandMarket, error) {
	var models []brandMarketModel
	err := db.WithContext(ctx).
		Where("brand_id = ? AND tenant_id = ?", brandID, tenantID.Int64()).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	result := make([]*product.BrandMarket, len(models))
	for i, m := range models {
		result[i] = m.toEntity()
	}
	return result, nil
}

func (r *brandMarketRepo) DeleteByBrand(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, brandID int64) error {
	return db.WithContext(ctx).
		Where("brand_id = ? AND tenant_id = ?", brandID, tenantID.Int64()).
		Delete(&brandMarketModel{}).Error
}

func (r *brandMarketRepo) BatchCreate(ctx context.Context, db *gorm.DB, items []*product.BrandMarket) error {
	if len(items) == 0 {
		return nil
	}
	models := make([]*brandMarketModel, len(items))
	for i, item := range items {
		models[i] = fromBrandMarketEntity(item)
	}
	return db.WithContext(ctx).CreateInBatches(models, 100).Error
}
```

- [ ] **Step 2: Run build**

Run: `cd admin && make build`

Expected: Build succeeded

- [ ] **Step 3: Commit**

```bash
git add admin/internal/infrastructure/persistence/brand_market_repository.go
git commit -m "feat(infra): implement BrandMarketRepository"
```

---

## Task 5: Create Brand API Definition

**Files:**
- Create: `admin/desc/brand.api`

- [ ] **Step 1: Create brand.api**

```go
syntax = "v1"

info (
    title:   "Brand API"
    desc:    "品牌管理相关接口"
    version: "v1"
)

type (
    CreateBrandReq {
        Name             string `json:"name"`
        Logo             string `json:"logo,optional"`
        Description      string `json:"description,optional"`
        Website          string `json:"website,optional"`
        TrademarkNumber  string `json:"trademark_number,optional"`
        TrademarkCountry string `json:"trademark_country,optional"`
        EnablePage       bool   `json:"enable_page,optional"`
        Sort             int    `json:"sort,optional"`
    }
    CreateBrandResp {
        ID int64 `json:"id"`
    }
    UpdateBrandReq {
        ID               int64  `path:"id"`
        Name             string `json:"name"`
        Logo             string `json:"logo,optional"`
        Description      string `json:"description,optional"`
        Website          string `json:"website,optional"`
        TrademarkNumber  string `json:"trademark_number,optional"`
        TrademarkCountry string `json:"trademark_country,optional"`
        EnablePage       bool   `json:"enable_page,optional"`
        Sort             int    `json:"sort,optional"`
    }
    BrandDetailResp {
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
    GetBrandReq {
        ID int64 `path:"id"`
    }
    ListBrandReq {
        Page     int    `form:"page,default=1"`
        PageSize int    `form:"page_size,default=20"`
        Name     string `form:"name,optional"`
        Status   int8   `form:"status,optional"`
    }
    ListBrandResp {
        List  []BrandDetailResp `json:"list"`
        Total int64             `json:"total"`
    }
    UpdateBrandStatusReq {
        ID     int64 `path:"id"`
        Status int8  `json:"status"`
    }
    ToggleBrandPageReq {
        ID      int64 `path:"id"`
        Enabled bool  `json:"enabled"`
    }
    GetBrandProductCountReq {
        ID int64 `path:"id"`
    }
    GetBrandProductCountResp {
        Count int64 `json:"count"`
    }
    SetBrandMarketVisibilityReq {
        BrandID   int64   `path:"id"`
        MarketIDs []int64 `json:"market_ids"`
        Visible   bool    `json:"visible"`
    }
    GetBrandMarketVisibilityReq {
        BrandID int64 `path:"id"`
    }
    BrandMarketVisibilityResp {
        BrandID int64                `json:"brand_id"`
        Markets []BrandMarketItemResp `json:"markets"`
    }
    BrandMarketItemResp {
        MarketID  int64 `json:"market_id"`
        IsVisible bool  `json:"is_visible"`
    }
)

@server (
    group:      brands
    middleware: AuthMiddleware
)
service admin-api {
    @doc "创建品牌"
    @handler CreateBrandHandler
    post /api/v1/brands (CreateBrandReq) returns (CreateBrandResp)

    @doc "更新品牌"
    @handler UpdateBrandHandler
    put /api/v1/brands/:id (UpdateBrandReq) returns (BrandDetailResp)

    @doc "获取品牌详情"
    @handler GetBrandHandler
    get /api/v1/brands/:id (GetBrandReq) returns (BrandDetailResp)

    @doc "获取品牌列表"
    @handler ListBrandsHandler
    get /api/v1/brands (ListBrandReq) returns (ListBrandResp)

    @doc "更新品牌状态"
    @handler UpdateBrandStatusHandler
    put /api/v1/brands/:id/status (UpdateBrandStatusReq) returns (BrandDetailResp)

    @doc "删除品牌"
    @handler DeleteBrandHandler
    delete /api/v1/brands/:id (GetBrandReq) returns (CreateBrandResp)

    @doc "切换品牌专区"
    @handler ToggleBrandPageHandler
    put /api/v1/brands/:id/toggle-page (ToggleBrandPageReq) returns (BrandDetailResp)

    @doc "获取品牌下商品数量"
    @handler GetBrandProductCountHandler
    get /api/v1/brands/:id/product-count (GetBrandProductCountReq) returns (GetBrandProductCountResp)

    @doc "设置品牌市场可见性"
    @handler SetBrandMarketVisibilityHandler
    put /api/v1/brands/:id/market-visibility (SetBrandMarketVisibilityReq) returns (CreateBrandResp)

    @doc "获取品牌市场可见性"
    @handler GetBrandMarketVisibilityHandler
    get /api/v1/brands/:id/market-visibility (GetBrandMarketVisibilityReq) returns (BrandMarketVisibilityResp)
}
```

- [ ] **Step 2: Update admin.api to import brand.api**

Read `admin/desc/admin.api` and add the import line at the end of the import block:

```go
import "brand.api"
```

- [ ] **Step 3: Generate code**

Run: `cd admin && make api`

Expected: No errors, files generated in `internal/handler/brands/` and `internal/types/`

- [ ] **Step 4: Commit**

```bash
git add admin/desc/brand.api admin/desc/admin.api admin/internal/handler/brands/ admin/internal/logic/brands/ admin/internal/types/types.go
git commit -m "feat(api): add Brand API definitions and generate handlers"
```

---

## Task 6: Update ServiceContext

**Files:**
- Modify: `admin/internal/svc/service_context.go`

- [ ] **Step 1: Read current ServiceContext**

Run: Read `admin/internal/svc/service_context.go`

- [ ] **Step 2: Add BrandRepo and BrandMarketRepo fields**

Add to the `ServiceContext` struct:

```go
    BrandRepo         product.BrandRepository
    BrandMarketRepo   product.BrandMarketRepository
```

- [ ] **Step 3: Initialize repositories in NewServiceContext function**

Add after existing repository initializations:

```go
        BrandRepo:        persistence.NewBrandRepository(),
        BrandMarketRepo:  persistence.NewBrandMarketRepository(),
```

- [ ] **Step 4: Run build**

Run: `cd admin && make build`

Expected: Build succeeded

- [ ] **Step 5: Commit**

```bash
git add admin/internal/svc/service_context.go
git commit -m "feat(svc): add BrandRepo to ServiceContext"
```

---

## Task 7: Implement CreateBrand Logic

**Files:**
- Modify: `admin/internal/logic/brands/create_brand_logic.go`

- [ ] **Step 1: Replace generated logic with implementation**

```go
package brands

import (
	"context"
	"errors"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateBrandLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateBrandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBrandLogic {
	return &CreateBrandLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateBrandLogic) CreateBrand(req *types.CreateBrandReq) (resp *types.CreateBrandResp, err error) {
	tenantID := l.ctx.Value("tenant_id").(int64)
	userID := l.ctx.Value("user_id").(int64)

	// Check if brand name already exists
	existing, err := l.svcCtx.BrandRepo.FindByName(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("brand name already exists")
	}

	// Generate ID
	id, err := l.svcCtx.IDGen.NextID(l.ctx)
	if err != nil {
		return nil, err
	}

	brand := &product.Brand{
		ID:               id,
		TenantID:         shared.TenantID(tenantID),
		Name:             req.Name,
		Logo:             req.Logo,
		Description:      req.Description,
		Website:          req.Website,
		Sort:             req.Sort,
		EnablePage:       req.EnablePage,
		TrademarkNumber:  req.TrademarkNumber,
		TrademarkCountry: req.TrademarkCountry,
		Status:           shared.StatusEnabled,
		Audit:            shared.NewAuditInfo(userID),
	}

	if err := l.svcCtx.BrandRepo.Create(l.ctx, l.svcCtx.DB, brand); err != nil {
		return nil, err
	}

	return &types.CreateBrandResp{ID: id}, nil
}
```

- [ ] **Step 2: Run build**

Run: `cd admin && make build`

Expected: Build succeeded

- [ ] **Step 3: Commit**

```bash
git add admin/internal/logic/brands/create_brand_logic.go
git commit -m "feat(logic): implement CreateBrand logic"
```

---

## Task 8: Implement ListBrands Logic

**Files:**
- Modify: `admin/internal/logic/brands/list_brands_logic.go`

- [ ] **Step 1: Replace with implementation**

```go
package brands

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListBrandsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListBrandsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBrandsLogic {
	return &ListBrandsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListBrandsLogic) ListBrands(req *types.ListBrandReq) (resp *types.ListBrandResp, err error) {
	tenantID := l.ctx.Value("tenant_id").(int64)

	query := product.BrandQuery{
		PageQuery: shared.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Name: req.Name,
	}
	if req.Status > 0 {
		query.Status = shared.Status(req.Status)
	}

	brands, total, err := l.svcCtx.BrandRepo.FindList(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), query)
	if err != nil {
		return nil, err
	}

	list := make([]types.BrandDetailResp, len(brands))
	for i, b := range brands {
		productCount, _ := l.svcCtx.BrandRepo.GetProductCount(l.ctx, l.svcCtx.DB, b.ID)
		list[i] = types.BrandDetailResp{
			ID:               b.ID,
			Name:             b.Name,
			Logo:             b.Logo,
			Description:      b.Description,
			Website:          b.Website,
			TrademarkNumber:  b.TrademarkNumber,
			TrademarkCountry: b.TrademarkCountry,
			EnablePage:       b.EnablePage,
			Sort:             b.Sort,
			Status:           int8(b.Status),
			ProductCount:     productCount,
			CreatedAt:        b.Audit.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:        b.Audit.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &types.ListBrandResp{
		List:  list,
		Total: total,
	}, nil
}
```

- [ ] **Step 2: Run build**

Run: `cd admin && make build`

Expected: Build succeeded

- [ ] **Step 3: Commit**

```bash
git add admin/internal/logic/brands/list_brands_logic.go
git commit -m "feat(logic): implement ListBrands logic"
```

---

## Task 9: Implement Remaining Brand Logic Handlers

**Files:**
- Modify: `admin/internal/logic/brands/*.go`

- [ ] **Step 1: Implement GetBrandLogic**

```go
package brands

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBrandLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBrandLogic {
	return &GetBrandLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandLogic) GetBrand(req *types.GetBrandReq) (resp *types.BrandDetailResp, err error) {
	tenantID := l.ctx.Value("tenant_id").(int64)

	brand, err := l.svcCtx.BrandRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, types.NewCodeError(404, "brand not found")
	}

	productCount, _ := l.svcCtx.BrandRepo.GetProductCount(l.ctx, l.svcCtx.DB, req.ID)

	return &types.BrandDetailResp{
		ID:               brand.ID,
		Name:             brand.Name,
		Logo:             brand.Logo,
		Description:      brand.Description,
		Website:          brand.Website,
		TrademarkNumber:  brand.TrademarkNumber,
		TrademarkCountry: brand.TrademarkCountry,
		EnablePage:       brand.EnablePage,
		Sort:             brand.Sort,
		Status:           int8(brand.Status),
		ProductCount:     productCount,
		CreatedAt:        brand.Audit.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        brand.Audit.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
```

- [ ] **Step 2: Implement UpdateBrandLogic**

```go
package brands

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBrandLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateBrandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBrandLogic {
	return &UpdateBrandLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBrandLogic) UpdateBrand(req *types.UpdateBrandReq) (resp *types.BrandDetailResp, err error) {
	tenantID := l.ctx.Value("tenant_id").(int64)
	userID := l.ctx.Value("user_id").(int64)

	brand, err := l.svcCtx.BrandRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, types.NewCodeError(404, "brand not found")
	}

	brand.Name = req.Name
	brand.Logo = req.Logo
	brand.Description = req.Description
	brand.Website = req.Website
	brand.TrademarkNumber = req.TrademarkNumber
	brand.TrademarkCountry = req.TrademarkCountry
	brand.EnablePage = req.EnablePage
	brand.Sort = req.Sort
	brand.Audit.UpdatedAt = time.Now()
	brand.Audit.UpdatedBy = userID

	if err := l.svcCtx.BrandRepo.Update(l.ctx, l.svcCtx.DB, brand); err != nil {
		return nil, err
	}

	productCount, _ := l.svcCtx.BrandRepo.GetProductCount(l.ctx, l.svcCtx.DB, req.ID)

	return &types.BrandDetailResp{
		ID:               brand.ID,
		Name:             brand.Name,
		Logo:             brand.Logo,
		Description:      brand.Description,
		Website:          brand.Website,
		TrademarkNumber:  brand.TrademarkNumber,
		TrademarkCountry: brand.TrademarkCountry,
		EnablePage:       brand.EnablePage,
		Sort:             brand.Sort,
		Status:           int8(brand.Status),
		ProductCount:     productCount,
		CreatedAt:        brand.Audit.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        brand.Audit.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
```

- [ ] **Step 3: Implement DeleteBrandLogic**

```go
package brands

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBrandLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteBrandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBrandLogic {
	return &DeleteBrandLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteBrandLogic) DeleteBrand(req *types.GetBrandReq) (resp *types.CreateBrandResp, err error) {
	tenantID := l.ctx.Value("tenant_id").(int64)

	// Check brand exists
	brand, err := l.svcCtx.BrandRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, types.NewCodeError(404, "brand not found")
	}

	// Set products' brand_id to NULL
	l.svcCtx.DB.Exec("UPDATE products SET brand_id = NULL WHERE brand_id = ?", req.ID)

	// Delete brand market visibility
	l.svcCtx.BrandMarketRepo.DeleteByBrand(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)

	// Soft delete brand
	if err := l.svcCtx.BrandRepo.Delete(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID); err != nil {
		return nil, err
	}

	return &types.CreateBrandResp{ID: req.ID}, nil
}
```

- [ ] **Step 4: Implement UpdateBrandStatusLogic**

```go
package brands

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBrandStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateBrandStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBrandStatusLogic {
	return &UpdateBrandStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBrandStatusLogic) UpdateBrandStatus(req *types.UpdateBrandStatusReq) (resp *types.BrandDetailResp, err error) {
	tenantID := l.ctx.Value("tenant_id").(int64)
	userID := l.ctx.Value("user_id").(int64)

	brand, err := l.svcCtx.BrandRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, types.NewCodeError(404, "brand not found")
	}

	if req.Status == 1 {
		brand.Enable()
	} else {
		brand.Disable()
	}
	brand.Audit.UpdatedAt = time.Now()
	brand.Audit.UpdatedBy = userID

	if err := l.svcCtx.BrandRepo.Update(l.ctx, l.svcCtx.DB, brand); err != nil {
		return nil, err
	}

	productCount, _ := l.svcCtx.BrandRepo.GetProductCount(l.ctx, l.svcCtx.DB, req.ID)

	return &types.BrandDetailResp{
		ID:               brand.ID,
		Name:             brand.Name,
		Logo:             brand.Logo,
		Description:      brand.Description,
		Website:          brand.Website,
		TrademarkNumber:  brand.TrademarkNumber,
		TrademarkCountry: brand.TrademarkCountry,
		EnablePage:       brand.EnablePage,
		Sort:             brand.Sort,
		Status:           int8(brand.Status),
		ProductCount:     productCount,
		CreatedAt:        brand.Audit.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        brand.Audit.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
```

- [ ] **Step 5: Implement ToggleBrandPageLogic**

```go
package brands

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type ToggleBrandPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewToggleBrandPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToggleBrandPageLogic {
	return &ToggleBrandPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToggleBrandPageLogic) ToggleBrandPage(req *types.ToggleBrandPageReq) (resp *types.BrandDetailResp, err error) {
	tenantID := l.ctx.Value("tenant_id").(int64)
	userID := l.ctx.Value("user_id").(int64)

	brand, err := l.svcCtx.BrandRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, types.NewCodeError(404, "brand not found")
	}

	brand.TogglePage(req.Enabled)
	brand.Audit.UpdatedAt = time.Now()
	brand.Audit.UpdatedBy = userID

	if err := l.svcCtx.BrandRepo.Update(l.ctx, l.svcCtx.DB, brand); err != nil {
		return nil, err
	}

	productCount, _ := l.svcCtx.BrandRepo.GetProductCount(l.ctx, l.svcCtx.DB, req.ID)

	return &types.BrandDetailResp{
		ID:               brand.ID,
		Name:             brand.Name,
		Logo:             brand.Logo,
		Description:      brand.Description,
		Website:          brand.Website,
		TrademarkNumber:  brand.TrademarkNumber,
		TrademarkCountry: brand.TrademarkCountry,
		EnablePage:       brand.EnablePage,
		Sort:             brand.Sort,
		Status:           int8(brand.Status),
		ProductCount:     productCount,
		CreatedAt:        brand.Audit.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        brand.Audit.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
```

- [ ] **Step 6: Implement Market Visibility Logic**

```go
// SetBrandMarketVisibilityLogic
package brands

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetBrandMarketVisibilityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetBrandMarketVisibilityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetBrandMarketVisibilityLogic {
	return &SetBrandMarketVisibilityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetBrandMarketVisibilityLogic) SetBrandMarketVisibility(req *types.SetBrandMarketVisibilityReq) (resp *types.CreateBrandResp, err error) {
	tenantID := l.ctx.Value("tenant_id").(int64)

	// Delete existing visibility for these markets
	for _, marketID := range req.MarketIDs {
		l.svcCtx.DB.Exec(
			"DELETE FROM brand_markets WHERE brand_id = ? AND market_id = ? AND tenant_id = ?",
			req.BrandID, marketID, tenantID,
		)
	}

	// Create new visibility records
	now := time.Now().Unix()
	items := make([]*product.BrandMarket, len(req.MarketIDs))
	for i, marketID := range req.MarketIDs {
		items[i] = &product.BrandMarket{
			TenantID:  shared.TenantID(tenantID),
			BrandID:   req.BrandID,
			MarketID:  marketID,
			IsVisible: req.Visible,
			Audit: shared.AuditInfo{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
	}

	if err := l.svcCtx.BrandMarketRepo.BatchCreate(l.ctx, l.svcCtx.DB, items); err != nil {
		return nil, err
	}

	return &types.CreateBrandResp{ID: req.BrandID}, nil
}
```

```go
// GetBrandMarketVisibilityLogic
package brands

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBrandMarketVisibilityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandMarketVisibilityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBrandMarketVisibilityLogic {
	return &GetBrandMarketVisibilityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandMarketVisibilityLogic) GetBrandMarketVisibility(req *types.GetBrandMarketVisibilityReq) (resp *types.BrandMarketVisibilityResp, err error) {
	tenantID := l.ctx.Value("tenant_id").(int64)

	markets, err := l.svcCtx.BrandMarketRepo.FindByBrand(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.BrandID)
	if err != nil {
		return nil, err
	}

	items := make([]types.BrandMarketItemResp, len(markets))
	for i, m := range markets {
		items[i] = types.BrandMarketItemResp{
			MarketID:  m.MarketID,
			IsVisible: m.IsVisible,
		}
	}

	return &types.BrandMarketVisibilityResp{
		BrandID: req.BrandID,
		Markets: items,
	}, nil
}
```

- [ ] **Step 7: Run build**

Run: `cd admin && make build`

Expected: Build succeeded

- [ ] **Step 8: Commit**

```bash
git add admin/internal/logic/brands/
git commit -m "feat(logic): implement all Brand CRUD logic handlers"
```

---

## Task 10: Create Brand API Client (Frontend)

**Files:**
- Create: `shop-admin/src/api/brand.ts`

- [ ] **Step 1: Create brand.ts**

```typescript
import request from '@/utils/request'

export interface Brand {
  id: number
  name: string
  logo: string
  description: string
  website: string
  trademark_number: string
  trademark_country: string
  enable_page: boolean
  sort: number
  status: number
  product_count: number
  created_at: string
  updated_at: string
}

export interface CreateBrandRequest {
  name: string
  logo?: string
  description?: string
  website?: string
  trademark_number?: string
  trademark_country?: string
  enable_page?: boolean
  sort?: number
}

export interface UpdateBrandRequest {
  id: number
  name: string
  logo?: string
  description?: string
  website?: string
  trademark_number?: string
  trademark_country?: string
  enable_page?: boolean
  sort?: number
}

export interface BrandListParams {
  page?: number
  page_size?: number
  name?: string
  status?: number
}

// Create brand
export function createBrand(data: CreateBrandRequest) {
  return request.post<{ id: number }>('/api/v1/brands', data)
}

// Update brand
export function updateBrand(data: UpdateBrandRequest) {
  return request.put<Brand>(`/api/v1/brands/${data.id}`, data)
}

// Get brand detail
export function getBrand(id: number) {
  return request.get<Brand>(`/api/v1/brands/${id}`)
}

// List brands
export function listBrands(params: BrandListParams) {
  return request.get<{ list: Brand[]; total: number }>('/api/v1/brands', { params })
}

// Update brand status
export function updateBrandStatus(id: number, status: number) {
  return request.put<Brand>(`/api/v1/brands/${id}/status`, { status })
}

// Delete brand
export function deleteBrand(id: number) {
  return request.delete<{ id: number }>(`/api/v1/brands/${id}`)
}

// Toggle brand page
export function toggleBrandPage(id: number, enabled: boolean) {
  return request.put<Brand>(`/api/v1/brands/${id}/toggle-page`, { enabled })
}

// Get product count
export function getBrandProductCount(id: number) {
  return request.get<{ count: number }>(`/api/v1/brands/${id}/product-count`)
}

// Set market visibility
export function setBrandMarketVisibility(brandId: number, marketIds: number[], visible: boolean) {
  return request.put<{ id: number }>(`/api/v1/brands/${brandId}/market-visibility`, {
    market_ids: marketIds,
    visible
  })
}

// Get market visibility
export function getBrandMarketVisibility(brandId: number) {
  return request.get<{
    brand_id: number
    markets: { market_id: number; is_visible: boolean }[]
  }>(`/api/v1/brands/${brandId}/market-visibility`)
}
```

- [ ] **Step 2: Commit**

```bash
git add shop-admin/src/api/brand.ts
git commit -m "feat(frontend): add Brand API client"
```

---

## Task 11: Create Brand Management Page (Frontend)

**Files:**
- Create: `shop-admin/src/views/brands/index.vue`

- [ ] **Step 1: Create index.vue**

```vue
<template>
  <div class="brand-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>Brand Management</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            Add Brand
          </el-button>
        </div>
      </template>

      <!-- Search -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="Name">
          <el-input v-model="searchForm.name" placeholder="Search by name" clearable @keyup.enter="loadList" />
        </el-form-item>
        <el-form-item label="Status">
          <el-select v-model="searchForm.status" placeholder="All" clearable>
            <el-option label="Enabled" :value="1" />
            <el-option label="Disabled" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadList">Search</el-button>
        </el-form-item>
      </el-form>

      <!-- Table -->
      <el-table :data="brandList" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="Logo" width="80">
          <template #default="{ row }">
            <el-image v-if="row.logo" :src="row.logo" style="width: 50px; height: 50px" fit="contain" />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="Name" min-width="120" />
        <el-table-column label="Brand Page" width="100">
          <template #default="{ row }">
            <el-switch :model-value="row.enable_page" @change="(val: boolean) => handleTogglePage(row, val)" />
          </template>
        </el-table-column>
        <el-table-column prop="product_count" label="Products" width="100" />
        <el-table-column label="Status" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? 'Enabled' : 'Disabled' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="Sort" width="80" />
        <el-table-column label="Trademark" min-width="150">
          <template #default="{ row }">
            <span v-if="row.trademark_number">{{ row.trademark_number }} ({{ row.trademark_country }})</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="Created" width="160" />
        <el-table-column label="Actions" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">Edit</el-button>
            <el-button link type="primary" @click="handleMarketVisibility(row)">Markets</el-button>
            <el-button link type="danger" @click="handleDelete(row)">Delete</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- Pagination -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next"
        @size-change="loadList"
        @current-change="loadList"
        class="pagination"
      />
    </el-card>

    <!-- Create/Edit Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="editingBrand ? 'Edit Brand' : 'Create Brand'"
      width="600px"
      destroy-on-close
    >
      <el-form :model="formData" label-width="140px" ref="formRef">
        <el-form-item label="Name" prop="name" :rules="[{ required: true, message: 'Name is required' }]">
          <el-input v-model="formData.name" placeholder="Enter brand name" />
        </el-form-item>
        <el-form-item label="Logo URL">
          <el-input v-model="formData.logo" placeholder="Enter logo URL" />
        </el-form-item>
        <el-form-item label="Description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="Enter description" />
        </el-form-item>
        <el-form-item label="Website">
          <el-input v-model="formData.website" placeholder="Enter website URL" />
        </el-form-item>
        <el-form-item label="Trademark Number">
          <el-input v-model="formData.trademark_number" placeholder="Enter trademark number" />
        </el-form-item>
        <el-form-item label="Trademark Country">
          <el-input v-model="formData.trademark_country" placeholder="e.g., US, CN, UK" maxlength="10" />
        </el-form-item>
        <el-form-item label="Enable Brand Page">
          <el-switch v-model="formData.enable_page" />
        </el-form-item>
        <el-form-item label="Sort">
          <el-input-number v-model="formData.sort" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">Submit</el-button>
      </template>
    </el-dialog>

    <!-- Market Visibility Dialog -->
    <el-dialog v-model="marketDialogVisible" title="Market Visibility" width="500px">
      <el-checkbox-group v-model="selectedMarkets">
        <el-checkbox v-for="market in markets" :key="market.id" :label="market.id">
          {{ market.name }}
        </el-checkbox>
      </el-checkbox-group>
      <template #footer>
        <el-button @click="marketDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="handleSaveMarketVisibility" :loading="submitting">Save</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  listBrands,
  createBrand,
  updateBrand,
  deleteBrand,
  toggleBrandPage,
  getBrandMarketVisibility,
  setBrandMarketVisibility,
  type Brand,
  type CreateBrandRequest
} from '@/api/brand'
import { listMarkets, type Market } from '@/api/market'

const brandList = ref<Brand[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const marketDialogVisible = ref(false)
const submitting = ref(false)
const editingBrand = ref<Brand | null>(null)
const formRef = ref()
const currentBrandId = ref<number>(0)
const markets = ref<Market[]>([])
const selectedMarkets = ref<number[]>([])

const searchForm = reactive({
  name: '',
  status: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const formData = reactive<CreateBrandRequest & { id?: number }>({
  name: '',
  logo: '',
  description: '',
  website: '',
  trademark_number: '',
  trademark_country: '',
  enable_page: false,
  sort: 0
})

const loadMarkets = async () => {
  try {
    const data = await listMarkets()
    markets.value = data.list || []
  } catch (error) {
    // Ignore if market API not available
  }
}

const loadList = async () => {
  loading.value = true
  try {
    const data = await listBrands({
      page: pagination.page,
      page_size: pagination.pageSize,
      name: searchForm.name || undefined,
      status: searchForm.status
    })
    brandList.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('Failed to load brands')
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  editingBrand.value = null
  Object.assign(formData, {
    name: '',
    logo: '',
    description: '',
    website: '',
    trademark_number: '',
    trademark_country: '',
    enable_page: false,
    sort: 0
  })
  dialogVisible.value = true
}

const handleEdit = (brand: Brand) => {
  editingBrand.value = brand
  Object.assign(formData, {
    id: brand.id,
    name: brand.name,
    logo: brand.logo,
    description: brand.description,
    website: brand.website,
    trademark_number: brand.trademark_number,
    trademark_country: brand.trademark_country,
    enable_page: brand.enable_page,
    sort: brand.sort
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await formRef.value?.validate()
  submitting.value = true
  try {
    if (editingBrand.value) {
      await updateBrand({ ...formData, id: editingBrand.value.id })
      ElMessage.success('Brand updated')
    } else {
      await createBrand(formData)
      ElMessage.success('Brand created')
    }
    dialogVisible.value = false
    loadList()
  } catch (error) {
    ElMessage.error('Failed to save brand')
  } finally {
    submitting.value = false
  }
}

const handleTogglePage = async (brand: Brand, enabled: boolean) => {
  try {
    await toggleBrandPage(brand.id, enabled)
    brand.enable_page = enabled
    ElMessage.success('Brand page updated')
  } catch (error) {
    ElMessage.error('Failed to update brand page')
  }
}

const handleDelete = async (brand: Brand) => {
  try {
    await ElMessageBox.confirm(
      `Delete "${brand.name}"? Products will have brand removed.`,
      'Confirm Delete',
      { type: 'warning' }
    )
    await deleteBrand(brand.id)
    ElMessage.success('Brand deleted')
    loadList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('Failed to delete brand')
    }
  }
}

const handleMarketVisibility = async (brand: Brand) => {
  currentBrandId.value = brand.id
  try {
    const data = await getBrandMarketVisibility(brand.id)
    selectedMarkets.value = data.markets
      .filter((m) => m.is_visible)
      .map((m) => m.market_id)
  } catch (error) {
    selectedMarkets.value = []
  }
  marketDialogVisible.value = true
}

const handleSaveMarketVisibility = async () => {
  submitting.value = true
  try {
    await setBrandMarketVisibility(currentBrandId.value, selectedMarkets.value, true)
    ElMessage.success('Market visibility updated')
    marketDialogVisible.value = false
  } catch (error) {
    ElMessage.error('Failed to update market visibility')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadList()
  loadMarkets()
})
</script>

<style scoped>
.brand-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  justify-content: flex-end;
}
</style>
```

- [ ] **Step 2: Add route (if needed)**

Check `shop-admin/src/router/index.ts` and add route for brands if not exists.

- [ ] **Step 3: Commit**

```bash
git add shop-admin/src/views/brands/index.vue
git commit -m "feat(frontend): add Brand management page with market visibility"
```

---

## Verification

- [ ] **Step 1: Run backend build**

Run: `cd admin && make build`

Expected: Build succeeded

- [ ] **Step 2: Start backend server**

Run: `cd admin && ./admin -f etc/admin.yaml`

Expected: Server starts on configured port

- [ ] **Step 3: Test API with curl**

```bash
# Create brand
curl -X POST http://localhost:8888/api/v1/brands \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name":"Apple","logo":"https://example.com/apple.png","trademark_number":"TM123456","trademark_country":"US"}'

# List brands
curl "http://localhost:8888/api/v1/brands?page=1&page_size=20" \
  -H "Authorization: Bearer <token>"
```

Expected: Successful responses

- [ ] **Step 4: Test frontend**

Run: `cd shop-admin && npm run dev`

Expected: Brand page loads at `/brands`, table displays correctly

---

## Rollback

```sql
-- Remove brand compliance fields
ALTER TABLE `brands`
    DROP COLUMN `enable_page`,
    DROP COLUMN `trademark_number`,
    DROP COLUMN `trademark_country`;

-- Drop brand_markets table
DROP TABLE IF EXISTS `brand_markets`;

-- Remove brand_id from products
ALTER TABLE `products`
    DROP COLUMN `brand_id`;
```