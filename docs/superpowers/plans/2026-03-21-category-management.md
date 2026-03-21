# Category Management Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement Category management with tree structure, SEO fields, and market visibility for the ShopJoy admin backend.

**Architecture:** Extend existing Category entity in `admin/internal/domain/product/category.go`, add repository implementation, create go-zero API handlers, and build Vue frontend for tree management.

**Tech Stack:** Go 1.21+, go-zero framework, GORM, MySQL, Vue 3 + TypeScript + Element Plus

**Spec:** `docs/superpowers/specs/2026-03-21-category-brand-inventory-design.md` (Section 1)

---

## File Structure

### New Files

```
admin/
├── desc/category.api                          # API definitions
├── internal/
│   ├── infrastructure/persistence/
│   │   ├── category_repository.go             # Repository implementation
│   │   └── category_market_repository.go      # Market visibility repo
│   ├── handler/categories/                    # Auto-generated
│   └── logic/categories/                     # Business logic
sql/
└── 20260321_category.sql                      # Migration script
shop-admin/src/
├── api/category.ts                            # API client
└── views/categories/index.vue                 # Management page
```

### Modified Files

```
admin/internal/domain/product/category.go      # Add SEO fields, CategoryMarket
admin/internal/svc/service_context.go          # Add CategoryRepo
admin/desc/admin.api                           # Import category.api
```

---

## Task 1: Database Schema

**Files:**
- Create: `sql/20260321_category.sql`

- [ ] **Step 1: Create migration file**

```sql
-- sql/20260321_category.sql
-- Category management schema changes

-- Add SEO fields to categories table
ALTER TABLE `categories`
    ADD COLUMN `seo_title` VARCHAR(200) DEFAULT '' COMMENT 'SEO标题' AFTER `image`,
    ADD COLUMN `seo_description` VARCHAR(500) DEFAULT '' COMMENT 'SEO描述' AFTER `seo_title`;

-- Create category_markets table
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

- [ ] **Step 2: Run migration**

Run: `mysql -u root -p shopjoy < sql/20260321_category.sql`

Expected: Query OK messages for ALTER and CREATE

- [ ] **Step 3: Commit**

```bash
git add sql/20260321_category.sql
git commit -m "feat(db): add category SEO fields and category_markets table"
```

---

## Task 2: Update Category Domain Entity

**Files:**
- Modify: `admin/internal/domain/product/category.go`

- [ ] **Step 1: Read current file**

Run: Read `admin/internal/domain/product/category.go`

- [ ] **Step 2: Add SEO fields to Category struct**

After the `Image` field, add:

```go
    // SEO Fields
    SeoTitle       string
    SeoDescription string
```

- [ ] **Step 3: Add CategoryMarket entity after Brand struct**

```go
// CategoryMarket represents category visibility in specific markets
type CategoryMarket struct {
    ID         int64
    TenantID   shared.TenantID
    CategoryID int64
    MarketID   int64
    IsVisible  bool
    Audit      shared.AuditInfo `gorm:"embedded"`
}

func (cm *CategoryMarket) TableName() string {
    return "category_markets"
}
```

- [ ] **Step 4: Add CategorySort struct**

```go
// CategorySort for batch sort updates
type CategorySort struct {
    ID   int64
    Sort int
}
```

- [ ] **Step 5: Add methods to CategoryRepository interface**

```go
// Add to CategoryRepository interface:
    FindByCode(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, code string) (*Category, error)
    GetProductCount(ctx context.Context, db *gorm.DB, categoryID int64) (int64, error)
    UpdateSort(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, sorts []CategorySort) error
    Move(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, newParentID int64) error
```

- [ ] **Step 6: Add CategoryMarketRepository interface**

```go
// CategoryMarketRepository interface
type CategoryMarketRepository interface {
    Create(ctx context.Context, db *gorm.DB, cm *CategoryMarket) error
    Update(ctx context.Context, db *gorm.DB, cm *CategoryMarket) error
    FindByCategory(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, categoryID int64) ([]*CategoryMarket, error)
    DeleteByCategory(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, categoryID int64) error
}
```

- [ ] **Step 7: Run build**

Run: `cd admin && make build`

Expected: Build succeeded

- [ ] **Step 8: Commit**

```bash
git add admin/internal/domain/product/category.go
git commit -m "feat(domain): add SEO fields and CategoryMarket entity to Category"
```

---

## Task 3: Implement Category Repository

**Files:**
- Create: `admin/internal/infrastructure/persistence/category_repository.go`

- [ ] **Step 1: Create category_repository.go**

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

type categoryRepo struct{}

func NewCategoryRepository() product.CategoryRepository {
	return &categoryRepo{}
}

type categoryModel struct {
	ID             int64  `gorm:"column:id;primaryKey"`
	TenantID       int64  `gorm:"column:tenant_id;not null;index"`
	ParentID       int64  `gorm:"column:parent_id;default:0;index"`
	Name           string `gorm:"column:name;type:varchar(100);not null"`
	Code           string `gorm:"column:code;type:varchar(50)"`
	Level          int    `gorm:"column:level;not null;default:1"`
	Sort           int    `gorm:"column:sort;default:0"`
	Icon           string `gorm:"column:icon;type:varchar(500)"`
	Image          string `gorm:"column:image;type:varchar(500)"`
	SeoTitle       string `gorm:"column:seo_title;type:varchar(200)"`
	SeoDescription string `gorm:"column:seo_description;type:varchar(500)"`
	Status         int8   `gorm:"column:status;not null;default:1"`
	CreatedAt      int64  `gorm:"column:created_at;not null"`
	UpdatedAt      int64  `gorm:"column:updated_at;not null"`
	CreatedBy      int64  `gorm:"column:created_by"`
	UpdatedBy      int64  `gorm:"column:updated_by"`
	DeletedAt      *int64 `gorm:"column:deleted_at;index"`
}

func (categoryModel) TableName() string {
	return "categories"
}

func (m *categoryModel) toEntity() *product.Category {
	return &product.Category{
		ID:             m.ID,
		TenantID:       shared.TenantID(m.TenantID),
		ParentID:       m.ParentID,
		Name:           m.Name,
		Code:           m.Code,
		Level:          m.Level,
		Sort:           m.Sort,
		Icon:           m.Icon,
		Image:          m.Image,
		SeoTitle:       m.SeoTitle,
		SeoDescription: m.SeoDescription,
		Status:         product.CategoryStatus(m.Status),
		Audit: shared.AuditInfo{
			CreatedAt: time.Unix(m.CreatedAt, 0),
			UpdatedAt: time.Unix(m.UpdatedAt, 0),
			CreatedBy: m.CreatedBy,
			UpdatedBy: m.UpdatedBy,
		},
	}
}

func fromCategoryEntity(c *product.Category) *categoryModel {
	now := time.Now().Unix()
	createdAt := now
	updatedAt := now
	if !c.Audit.CreatedAt.IsZero() {
		createdAt = c.Audit.CreatedAt.Unix()
	}
	if !c.Audit.UpdatedAt.IsZero() {
		updatedAt = c.Audit.UpdatedAt.Unix()
	}
	return &categoryModel{
		ID:             c.ID,
		TenantID:       c.TenantID.Int64(),
		ParentID:       c.ParentID,
		Name:           c.Name,
		Code:           c.Code,
		Level:          c.Level,
		Sort:           c.Sort,
		Icon:           c.Icon,
		Image:          c.Image,
		SeoTitle:       c.SeoTitle,
		SeoDescription: c.SeoDescription,
		Status:         int8(c.Status),
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		CreatedBy:      c.Audit.CreatedBy,
		UpdatedBy:      c.Audit.UpdatedBy,
	}
}

func (r *categoryRepo) Create(ctx context.Context, db *gorm.DB, c *product.Category) error {
	model := fromCategoryEntity(c)
	return db.WithContext(ctx).Create(model).Error
}

func (r *categoryRepo) Update(ctx context.Context, db *gorm.DB, c *product.Category) error {
	model := fromCategoryEntity(c)
	return db.WithContext(ctx).Model(&categoryModel{}).
		Where("id = ? AND tenant_id = ?", c.ID, c.TenantID.Int64()).
		Updates(map[string]interface{}{
			"name":            model.Name,
			"code":            model.Code,
			"parent_id":       model.ParentID,
			"level":           model.Level,
			"sort":            model.Sort,
			"icon":            model.Icon,
			"image":           model.Image,
			"seo_title":       model.SeoTitle,
			"seo_description": model.SeoDescription,
			"status":          model.Status,
			"updated_at":      model.UpdatedAt,
			"updated_by":      model.UpdatedBy,
		}).Error
}

func (r *categoryRepo) Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	now := time.Now().Unix()
	return db.WithContext(ctx).Model(&categoryModel{}).
		Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).
		Update("deleted_at", now).Error
}

func (r *categoryRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*product.Category, error) {
	var model categoryModel
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

func (r *categoryRepo) FindByParentID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, parentID int64) ([]*product.Category, error) {
	var models []categoryModel
	err := db.WithContext(ctx).
		Where("parent_id = ? AND tenant_id = ? AND deleted_at IS NULL", parentID, tenantID.Int64()).
		Order("sort ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	categories := make([]*product.Category, len(models))
	for i, m := range models {
		categories[i] = m.toEntity()
	}
	return categories, nil
}

func (r *categoryRepo) FindAll(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) ([]*product.Category, error) {
	var models []categoryModel
	err := db.WithContext(ctx).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID.Int64()).
		Order("sort ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	categories := make([]*product.Category, len(models))
	for i, m := range models {
		categories[i] = m.toEntity()
	}
	return categories, nil
}

func (r *categoryRepo) FindTree(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) ([]*product.Category, error) {
	return r.FindAll(ctx, db, tenantID)
}

func (r *categoryRepo) FindByCode(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, code string) (*product.Category, error) {
	var model categoryModel
	err := db.WithContext(ctx).
		Where("code = ? AND tenant_id = ? AND deleted_at IS NULL", code, tenantID.Int64()).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *categoryRepo) GetProductCount(ctx context.Context, db *gorm.DB, categoryID int64) (int64, error) {
	var count int64
	err := db.WithContext(ctx).Table("products").
		Where("category_id = ?", categoryID).
		Count(&count).Error
	return count, err
}

func (r *categoryRepo) UpdateSort(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, sorts []product.CategorySort) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, s := range sorts {
			if err := tx.Model(&categoryModel{}).
				Where("id = ? AND tenant_id = ?", s.ID, tenantID.Int64()).
				Update("sort", s.Sort).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *categoryRepo) Move(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, newParentID int64) error {
	var parentLevel int = 0
	if newParentID > 0 {
		var parent categoryModel
		if err := db.WithContext(ctx).
			Where("id = ? AND tenant_id = ?", newParentID, tenantID.Int64()).
			First(&parent).Error; err != nil {
			return err
		}
		parentLevel = parent.Level
	}

	return db.WithContext(ctx).Model(&categoryModel{}).
		Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).
		Updates(map[string]interface{}{
			"parent_id": newParentID,
			"level":     parentLevel + 1,
		}).Error
}
```

- [ ] **Step 2: Run build**

Run: `cd admin && make build`

Expected: Build succeeded

- [ ] **Step 3: Commit**

```bash
git add admin/internal/infrastructure/persistence/category_repository.go
git commit -m "feat(infra): implement CategoryRepository with SEO support"
```

---

## Task 4: Implement CategoryMarket Repository

**Files:**
- Create: `admin/internal/infrastructure/persistence/category_market_repository.go`

- [ ] **Step 1: Create category_market_repository.go**

```go
package persistence

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type categoryMarketRepo struct{}

func NewCategoryMarketRepository() product.CategoryMarketRepository {
	return &categoryMarketRepo{}
}

type categoryMarketModel struct {
	ID         int64 `gorm:"column:id;primaryKey"`
	TenantID   int64 `gorm:"column:tenant_id;not null;index"`
	CategoryID int64 `gorm:"column:category_id;not null;index"`
	MarketID   int64 `gorm:"column:market_id;not null;index"`
	IsVisible  bool  `gorm:"column:is_visible;not null;default:true"`
	CreatedAt  int64 `gorm:"column:created_at;not null"`
	UpdatedAt  int64 `gorm:"column:updated_at;not null"`
}

func (categoryMarketModel) TableName() string {
	return "category_markets"
}

func (m *categoryMarketModel) toEntity() *product.CategoryMarket {
	return &product.CategoryMarket{
		ID:         m.ID,
		TenantID:   shared.TenantID(m.TenantID),
		CategoryID: m.CategoryID,
		MarketID:   m.MarketID,
		IsVisible:  m.IsVisible,
		Audit: shared.AuditInfo{
			CreatedAt: time.Unix(m.CreatedAt, 0),
			UpdatedAt: time.Unix(m.UpdatedAt, 0),
		},
	}
}

func fromCategoryMarketEntity(cm *product.CategoryMarket) *categoryMarketModel {
	now := time.Now().Unix()
	createdAt := now
	updatedAt := now
	if !cm.Audit.CreatedAt.IsZero() {
		createdAt = cm.Audit.CreatedAt.Unix()
	}
	if !cm.Audit.UpdatedAt.IsZero() {
		updatedAt = cm.Audit.UpdatedAt.Unix()
	}
	return &categoryMarketModel{
		ID:         cm.ID,
		TenantID:   cm.TenantID.Int64(),
		CategoryID: cm.CategoryID,
		MarketID:   cm.MarketID,
		IsVisible:  cm.IsVisible,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

func (r *categoryMarketRepo) Create(ctx context.Context, db *gorm.DB, cm *product.CategoryMarket) error {
	model := fromCategoryMarketEntity(cm)
	return db.WithContext(ctx).Create(model).Error
}

func (r *categoryMarketRepo) Update(ctx context.Context, db *gorm.DB, cm *product.CategoryMarket) error {
	return db.WithContext(ctx).Model(&categoryMarketModel{}).
		Where("id = ?", cm.ID).
		Updates(map[string]interface{}{
			"is_visible":  cm.IsVisible,
			"updated_at":  time.Now().Unix(),
		}).Error
}

func (r *categoryMarketRepo) FindByCategory(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, categoryID int64) ([]*product.CategoryMarket, error) {
	var models []categoryMarketModel
	err := db.WithContext(ctx).
		Where("category_id = ? AND tenant_id = ?", categoryID, tenantID.Int64()).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	result := make([]*product.CategoryMarket, len(models))
	for i, m := range models {
		result[i] = m.toEntity()
	}
	return result, nil
}

func (r *categoryMarketRepo) DeleteByCategory(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, categoryID int64) error {
	return db.WithContext(ctx).
		Where("category_id = ? AND tenant_id = ?", categoryID, tenantID.Int64()).
		Delete(&categoryMarketModel{}).Error
}
```

- [ ] **Step 2: Run build**

Run: `cd admin && make build`

Expected: Build succeeded

- [ ] **Step 3: Commit**

```bash
git add admin/internal/infrastructure/persistence/category_market_repository.go
git commit -m "feat(infra): implement CategoryMarketRepository"
```

---

## Task 5: Create Category API Definition

**Files:**
- Create: `admin/desc/category.api`

- [ ] **Step 1: Create category.api**

```go
syntax = "v1"

info (
    title:   "Category API"
    desc:    "分类管理相关接口"
    version: "v1"
)

type (
    CreateCategoryReq {
        Name           string `json:"name"`
        ParentID       int64  `json:"parent_id,optional"`
        Code           string `json:"code,optional"`
        Icon           string `json:"icon,optional"`
        Image          string `json:"image,optional"`
        SeoTitle       string `json:"seo_title,optional"`
        SeoDescription string `json:"seo_description,optional"`
        Sort           int    `json:"sort,optional"`
    }
    CreateCategoryResp {
        ID int64 `json:"id"`
    }
    UpdateCategoryReq {
        ID             int64  `path:"id"`
        Name           string `json:"name"`
        Code           string `json:"code,optional"`
        Icon           string `json:"icon,optional"`
        Image          string `json:"image,optional"`
        SeoTitle       string `json:"seo_title,optional"`
        SeoDescription string `json:"seo_description,optional"`
        Sort           int    `json:"sort,optional"`
    }
    CategoryDetailResp {
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
    CategoryTreeResp {
        CategoryDetailResp
        Children []*CategoryTreeResp `json:"children"`
    }
    GetCategoryReq {
        ID int64 `path:"id"`
    }
    ListCategoryReq {
        ParentID int64 `form:"parent_id,optional"`
    }
    ListCategoryResp {
        List []*CategoryDetailResp `json:"list"`
    }
    CategoryTreeReq {}
    UpdateCategoryStatusReq {
        ID     int64 `path:"id"`
        Status int8  `json:"status"`
    }
    UpdateCategorySortReq {
        Sorts []CategorySortItem `json:"sorts"`
    }
    CategorySortItem {
        ID   int64 `json:"id"`
        Sort int   `json:"sort"`
    }
    MoveCategoryReq {
        ID          int64 `path:"id"`
        NewParentID int64 `json:"new_parent_id"`
    }
    GetCategoryProductCountReq {
        ID int64 `path:"id"`
    }
    GetCategoryProductCountResp {
        Count int64 `json:"count"`
    }
    SetCategoryMarketVisibilityReq {
        CategoryID int64   `path:"id"`
        MarketIDs  []int64 `json:"market_ids"`
        Visible    bool    `json:"visible"`
    }
    GetCategoryMarketVisibilityReq {
        CategoryID int64 `path:"id"`
    }
    CategoryMarketVisibilityResp {
        CategoryID int64                   `json:"category_id"`
        Markets    []CategoryMarketItemResp `json:"markets"`
    }
    CategoryMarketItemResp {
        MarketID  int64 `json:"market_id"`
        IsVisible bool  `json:"is_visible"`
    }
)

@server (
    group:      categories
    middleware: AuthMiddleware
)
service admin-api {
    @doc "创建分类"
    @handler CreateCategoryHandler
    post /api/v1/categories (CreateCategoryReq) returns (CreateCategoryResp)

    @doc "更新分类"
    @handler UpdateCategoryHandler
    put /api/v1/categories/:id (UpdateCategoryReq) returns (CategoryDetailResp)

    @doc "获取分类详情"
    @handler GetCategoryHandler
    get /api/v1/categories/:id (GetCategoryReq) returns (CategoryDetailResp)

    @doc "获取分类列表"
    @handler ListCategoriesHandler
    get /api/v1/categories (ListCategoryReq) returns (ListCategoryResp)

    @doc "获取分类树"
    @handler GetCategoryTreeHandler
    get /api/v1/categories/tree (CategoryTreeReq) returns ([]CategoryTreeResp)

    @doc "更新分类状态"
    @handler UpdateCategoryStatusHandler
    put /api/v1/categories/:id/status (UpdateCategoryStatusReq) returns (CategoryDetailResp)

    @doc "删除分类"
    @handler DeleteCategoryHandler
    delete /api/v1/categories/:id (GetCategoryReq) returns (CreateCategoryResp)

    @doc "更新分类排序"
    @handler UpdateCategorySortHandler
    put /api/v1/categories/sort (UpdateCategorySortReq) returns (CreateCategoryResp)

    @doc "移动分类"
    @handler MoveCategoryHandler
    put /api/v1/categories/:id/move (MoveCategoryReq) returns (CategoryDetailResp)

    @doc "获取分类下商品数量"
    @handler GetCategoryProductCountHandler
    get /api/v1/categories/:id/product-count (GetCategoryProductCountReq) returns (GetCategoryProductCountResp)

    @doc "设置分类市场可见性"
    @handler SetCategoryMarketVisibilityHandler
    put /api/v1/categories/:id/market-visibility (SetCategoryMarketVisibilityReq) returns (CreateCategoryResp)

    @doc "获取分类市场可见性"
    @handler GetCategoryMarketVisibilityHandler
    get /api/v1/categories/:id/market-visibility (GetCategoryMarketVisibilityReq) returns (CategoryMarketVisibilityResp)
}
```

- [ ] **Step 2: Update admin.api to import category.api**

Read `admin/desc/admin.api` and add the import line at the end of the import block:

```go
import "category.api"
```

- [ ] **Step 3: Generate code**

Run: `cd admin && make api`

Expected: No errors, files generated in `internal/handler/categories/` and `internal/types/`

- [ ] **Step 4: Commit**

```bash
git add admin/desc/category.api admin/desc/admin.api admin/internal/handler/categories/ admin/internal/logic/categories/ admin/internal/types/types.go
git commit -m "feat(api): add Category API definitions and generate handlers"
```

---

## Task 6: Update ServiceContext

**Files:**
- Modify: `admin/internal/svc/service_context.go`

- [ ] **Step 1: Read current ServiceContext**

Run: Read `admin/internal/svc/service_context.go`

- [ ] **Step 2: Add CategoryRepo and CategoryMarketRepo fields**

Add to the `ServiceContext` struct:

```go
    CategoryRepo        product.CategoryRepository
    CategoryMarketRepo  product.CategoryMarketRepository
```

- [ ] **Step 3: Initialize repositories in NewServiceContext function**

Add after existing repository initializations:

```go
        CategoryRepo:       persistence.NewCategoryRepository(),
        CategoryMarketRepo: persistence.NewCategoryMarketRepository(),
```

- [ ] **Step 4: Run build**

Run: `cd admin && make build`

Expected: Build succeeded

- [ ] **Step 5: Commit**

```bash
git add admin/internal/svc/service_context.go
git commit -m "feat(svc): add CategoryRepo to ServiceContext"
```

---

## Task 7: Implement CreateCategory Logic

**Files:**
- Modify: `admin/internal/logic/categories/create_category_logic.go`

- [ ] **Step 1: Replace generated logic with implementation**

```go
package categories

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCategoryLogic {
	return &CreateCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCategoryLogic) CreateCategory(req *types.CreateCategoryReq) (resp *types.CreateCategoryResp, err error) {
	tenantID := l.ctx.Value("tenant_id").(int64)
	userID := l.ctx.Value("user_id").(int64)

	// Calculate level from parent
	level := 1
	if req.ParentID > 0 {
		parent, err := l.svcCtx.CategoryRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ParentID)
		if err != nil {
			return nil, err
		}
		if parent != nil {
			level = parent.Level + 1
		}
	}

	// Generate ID
	id, err := l.svcCtx.IDGen.NextID(l.ctx)
	if err != nil {
		return nil, err
	}

	category := &product.Category{
		ID:             id,
		TenantID:       shared.TenantID(tenantID),
		ParentID:       req.ParentID,
		Name:           req.Name,
		Code:           req.Code,
		Level:          level,
		Sort:           req.Sort,
		Icon:           req.Icon,
		Image:          req.Image,
		SeoTitle:       req.SeoTitle,
		SeoDescription: req.SeoDescription,
		Status:         product.CategoryStatusEnabled,
		Audit:          shared.NewAuditInfo(userID),
	}

	if err := l.svcCtx.CategoryRepo.Create(l.ctx, l.svcCtx.DB, category); err != nil {
		return nil, err
	}

	return &types.CreateCategoryResp{ID: id}, nil
}
```

- [ ] **Step 2: Run build**

Run: `cd admin && make build`

Expected: Build succeeded

- [ ] **Step 3: Commit**

```bash
git add admin/internal/logic/categories/create_category_logic.go
git commit -m "feat(logic): implement CreateCategory logic"
```

---

## Task 8: Implement GetCategoryTree Logic

**Files:**
- Modify: `admin/internal/logic/categories/get_category_tree_logic.go`

- [ ] **Step 1: Replace with implementation**

```go
package categories

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoryTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCategoryTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryTreeLogic {
	return &GetCategoryTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCategoryTreeLogic) GetCategoryTree(req *types.CategoryTreeReq) (resp []*types.CategoryTreeResp, err error) {
	tenantID := l.ctx.Value("tenant_id").(int64)

	categories, err := l.svcCtx.CategoryRepo.FindAll(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID))
	if err != nil {
		return nil, err
	}

	// Build tree recursively
	tree := buildCategoryTree(categories, 0)
	return tree, nil
}

func buildCategoryTree(categories []*product.Category, parentID int64) []*types.CategoryTreeResp {
	var result []*types.CategoryTreeResp
	for _, cat := range categories {
		if cat.ParentID == parentID {
			node := &types.CategoryTreeResp{
				CategoryDetailResp: types.CategoryDetailResp{
					ID:             cat.ID,
					ParentID:       cat.ParentID,
					Name:           cat.Name,
					Code:           cat.Code,
					Level:          cat.Level,
					Sort:           cat.Sort,
					Icon:           cat.Icon,
					Image:          cat.Image,
					SeoTitle:       cat.SeoTitle,
					SeoDescription: cat.SeoDescription,
					Status:         int8(cat.Status),
					CreatedAt:      cat.Audit.CreatedAt.Format("2006-01-02 15:04:05"),
					UpdatedAt:      cat.Audit.UpdatedAt.Format("2006-01-02 15:04:05"),
				},
				Children: buildCategoryTree(categories, cat.ID),
			}
			result = append(result, node)
		}
	}
	return result
}
```

- [ ] **Step 2: Run build**

Run: `cd admin && make build`

Expected: Build succeeded

- [ ] **Step 3: Commit**

```bash
git add admin/internal/logic/categories/get_category_tree_logic.go
git commit -m "feat(logic): implement GetCategoryTree logic with tree building"
```

---

## Task 9: Implement Remaining Category Logic Handlers

**Files:**
- Modify: `admin/internal/logic/categories/*.go`

- [ ] **Step 1: Implement GetCategoryLogic**

```go
package categories

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryLogic {
	return &GetCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCategoryLogic) GetCategory(req *types.GetCategoryReq) (resp *types.CategoryDetailResp, err error) {
	tenantID := l.ctx.Value("tenant_id").(int64)

	cat, err := l.svcCtx.CategoryRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, types.NewCodeError(404, "category not found")
	}

	productCount, _ := l.svcCtx.CategoryRepo.GetProductCount(l.ctx, l.svcCtx.DB, req.ID)

	return &types.CategoryDetailResp{
		ID:             cat.ID,
		ParentID:       cat.ParentID,
		Name:           cat.Name,
		Code:           cat.Code,
		Level:          cat.Level,
		Sort:           cat.Sort,
		Icon:           cat.Icon,
		Image:          cat.Image,
		SeoTitle:       cat.SeoTitle,
		SeoDescription: cat.SeoDescription,
		Status:         int8(cat.Status),
		ProductCount:   productCount,
		CreatedAt:      cat.Audit.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      cat.Audit.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
```

- [ ] **Step 2: Implement UpdateCategoryLogic**

```go
package categories

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCategoryLogic {
	return &UpdateCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCategoryLogic) UpdateCategory(req *types.UpdateCategoryReq) (resp *types.CategoryDetailResp, err error) {
	tenantID := l.ctx.Value("tenant_id").(int64)
	userID := l.ctx.Value("user_id").(int64)

	cat, err := l.svcCtx.CategoryRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, types.NewCodeError(404, "category not found")
	}

	cat.Name = req.Name
	cat.Code = req.Code
	cat.Icon = req.Icon
	cat.Image = req.Image
	cat.SeoTitle = req.SeoTitle
	cat.SeoDescription = req.SeoDescription
	cat.Sort = req.Sort
	cat.Audit.UpdatedAt = time.Now()
	cat.Audit.UpdatedBy = userID

	if err := l.svcCtx.CategoryRepo.Update(l.ctx, l.svcCtx.DB, cat); err != nil {
		return nil, err
	}

	productCount, _ := l.svcCtx.CategoryRepo.GetProductCount(l.ctx, l.svcCtx.DB, req.ID)

	return &types.CategoryDetailResp{
		ID:             cat.ID,
		ParentID:       cat.ParentID,
		Name:           cat.Name,
		Code:           cat.Code,
		Level:          cat.Level,
		Sort:           cat.Sort,
		Icon:           cat.Icon,
		Image:          cat.Image,
		SeoTitle:       cat.SeoTitle,
		SeoDescription: cat.SeoDescription,
		Status:         int8(cat.Status),
		ProductCount:   productCount,
		CreatedAt:      cat.Audit.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      cat.Audit.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
```

- [ ] **Step 3: Implement DeleteCategoryLogic**

```go
package categories

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCategoryLogic {
	return &DeleteCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCategoryLogic) DeleteCategory(req *types.GetCategoryReq) (resp *types.CreateCategoryResp, err error) {
	tenantID := l.ctx.Value("tenant_id").(int64)

	// Check if category has children
	children, err := l.svcCtx.CategoryRepo.FindByParentID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if len(children) > 0 {
		return nil, types.NewCodeError(400, "cannot delete category with children")
	}

	// Migrate products to parent category
	cat, err := l.svcCtx.CategoryRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if cat != nil {
		l.svcCtx.DB.Exec("UPDATE products SET category_id = ? WHERE category_id = ?", cat.ParentID, req.ID)
	}

	if err := l.svcCtx.CategoryRepo.Delete(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID); err != nil {
		return nil, err
	}

	return &types.CreateCategoryResp{ID: req.ID}, nil
}
```

- [ ] **Step 4: Implement UpdateCategoryStatusLogic**

```go
package categories

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCategoryStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCategoryStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCategoryStatusLogic {
	return &UpdateCategoryStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCategoryStatusLogic) UpdateCategoryStatus(req *types.UpdateCategoryStatusReq) (resp *types.CategoryDetailResp, err error) {
	tenantID := l.ctx.Value("tenant_id").(int64)
	userID := l.ctx.Value("user_id").(int64)

	cat, err := l.svcCtx.CategoryRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, types.NewCodeError(404, "category not found")
	}

	if req.Status == 1 {
		cat.Enable()
	} else {
		cat.Disable()
	}
	cat.Audit.UpdatedAt = time.Now()
	cat.Audit.UpdatedBy = userID

	if err := l.svcCtx.CategoryRepo.Update(l.ctx, l.svcCtx.DB, cat); err != nil {
		return nil, err
	}

	productCount, _ := l.svcCtx.CategoryRepo.GetProductCount(l.ctx, l.svcCtx.DB, req.ID)

	return &types.CategoryDetailResp{
		ID:             cat.ID,
		ParentID:       cat.ParentID,
		Name:           cat.Name,
		Code:           cat.Code,
		Level:          cat.Level,
		Sort:           cat.Sort,
		Icon:           cat.Icon,
		Image:          cat.Image,
		SeoTitle:       cat.SeoTitle,
		SeoDescription: cat.SeoDescription,
		Status:         int8(cat.Status),
		ProductCount:   productCount,
		CreatedAt:      cat.Audit.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      cat.Audit.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
```

- [ ] **Step 5: Implement UpdateCategorySortLogic**

```go
package categories

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCategorySortLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCategorySortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCategorySortLogic {
	return &UpdateCategorySortLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCategorySortLogic) UpdateCategorySort(req *types.UpdateCategorySortReq) (resp *types.CreateCategoryResp, err error) {
	tenantID := l.ctx.Value("tenant_id").(int64)

	sorts := make([]product.CategorySort, len(req.Sorts))
	for i, s := range req.Sorts {
		sorts[i] = product.CategorySort{ID: s.ID, Sort: s.Sort}
	}

	if err := l.svcCtx.CategoryRepo.UpdateSort(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), sorts); err != nil {
		return nil, err
	}

	return &types.CreateCategoryResp{ID: 0}, nil
}
```

- [ ] **Step 6: Implement MoveCategoryLogic**

```go
package categories

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type MoveCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMoveCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MoveCategoryLogic {
	return &MoveCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MoveCategoryLogic) MoveCategory(req *types.MoveCategoryReq) (resp *types.CategoryDetailResp, err error) {
	tenantID := l.ctx.Value("tenant_id").(int64)
	userID := l.ctx.Value("user_id").(int64)

	if err := l.svcCtx.CategoryRepo.Move(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID, req.NewParentID); err != nil {
		return nil, err
	}

	cat, err := l.svcCtx.CategoryRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}

	cat.Audit.UpdatedAt = time.Now()
	cat.Audit.UpdatedBy = userID
	l.svcCtx.CategoryRepo.Update(l.ctx, l.svcCtx.DB, cat)

	productCount, _ := l.svcCtx.CategoryRepo.GetProductCount(l.ctx, l.svcCtx.DB, req.ID)

	return &types.CategoryDetailResp{
		ID:             cat.ID,
		ParentID:       cat.ParentID,
		Name:           cat.Name,
		Code:           cat.Code,
		Level:          cat.Level,
		Sort:           cat.Sort,
		Icon:           cat.Icon,
		Image:          cat.Image,
		SeoTitle:       cat.SeoTitle,
		SeoDescription: cat.SeoDescription,
		Status:         int8(cat.Status),
		ProductCount:   productCount,
		CreatedAt:      cat.Audit.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      cat.Audit.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
```

- [ ] **Step 7: Run build**

Run: `cd admin && make build`

Expected: Build succeeded

- [ ] **Step 8: Commit**

```bash
git add admin/internal/logic/categories/
git commit -m "feat(logic): implement all Category CRUD logic handlers"
```

---

## Task 10: Create Category API Client (Frontend)

**Files:**
- Create: `shop-admin/src/api/category.ts`

- [ ] **Step 1: Create category.ts**

```typescript
import request from '@/utils/request'

export interface Category {
  id: number
  parent_id: number
  name: string
  code: string
  level: number
  sort: number
  icon: string
  image: string
  seo_title: string
  seo_description: string
  status: number
  product_count: number
  created_at: string
  updated_at: string
}

export interface CategoryTree extends Category {
  children: CategoryTree[]
}

export interface CreateCategoryRequest {
  name: string
  parent_id?: number
  code?: string
  icon?: string
  image?: string
  seo_title?: string
  seo_description?: string
  sort?: number
}

export interface UpdateCategoryRequest {
  id: number
  name: string
  code?: string
  icon?: string
  image?: string
  seo_title?: string
  seo_description?: string
  sort?: number
}

export interface CategorySortItem {
  id: number
  sort: number
}

// Create category
export function createCategory(data: CreateCategoryRequest) {
  return request.post<{ id: number }>('/api/v1/categories', data)
}

// Update category
export function updateCategory(data: UpdateCategoryRequest) {
  return request.put<Category>(`/api/v1/categories/${data.id}`, data)
}

// Get category detail
export function getCategory(id: number) {
  return request.get<Category>(`/api/v1/categories/${id}`)
}

// List categories
export function listCategories(parentId?: number) {
  return request.get<{ list: Category[] }>('/api/v1/categories', {
    params: { parent_id: parentId }
  })
}

// Get category tree
export function getCategoryTree() {
  return request.get<CategoryTree[]>('/api/v1/categories/tree')
}

// Update category status
export function updateCategoryStatus(id: number, status: number) {
  return request.put<Category>(`/api/v1/categories/${id}/status`, { status })
}

// Delete category
export function deleteCategory(id: number) {
  return request.delete<{ id: number }>(`/api/v1/categories/${id}`)
}

// Update category sort
export function updateCategorySort(sorts: CategorySortItem[]) {
  return request.put<{ id: number }>('/api/v1/categories/sort', { sorts })
}

// Move category
export function moveCategory(id: number, newParentId: number) {
  return request.put<Category>(`/api/v1/categories/${id}/move`, { new_parent_id: newParentId })
}

// Get product count
export function getCategoryProductCount(id: number) {
  return request.get<{ count: number }>(`/api/v1/categories/${id}/product-count`)
}

// Set market visibility
export function setCategoryMarketVisibility(categoryId: number, marketIds: number[], visible: boolean) {
  return request.put<{ id: number }>(`/api/v1/categories/${categoryId}/market-visibility`, {
    market_ids: marketIds,
    visible
  })
}

// Get market visibility
export function getCategoryMarketVisibility(categoryId: number) {
  return request.get<{
    category_id: number
    markets: { market_id: number; is_visible: boolean }[]
  }>(`/api/v1/categories/${categoryId}/market-visibility`)
}
```

- [ ] **Step 2: Commit**

```bash
git add shop-admin/src/api/category.ts
git commit -m "feat(frontend): add Category API client"
```

---

## Task 11: Create Category Management Page (Frontend)

**Files:**
- Create: `shop-admin/src/views/categories/index.vue`

- [ ] **Step 1: Create index.vue**

```vue
<template>
  <div class="category-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>Category Management</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            Add Category
          </el-button>
        </div>
      </template>

      <el-tree
        ref="treeRef"
        :data="categoryTree"
        :props="{ label: 'name', children: 'children' }"
        node-key="id"
        default-expand-all
        draggable
        @node-drop="handleDrop"
      >
        <template #default="{ node, data }">
          <div class="tree-node">
            <div class="node-content">
              <el-icon v-if="data.icon"><component :is="data.icon" /></el-icon>
              <span class="node-name">{{ data.name }}</span>
              <el-tag :type="data.status === 1 ? 'success' : 'info'" size="small">
                {{ data.status === 1 ? 'Enabled' : 'Disabled' }}
              </el-tag>
              <span class="product-count">{{ data.product_count }} products</span>
            </div>
            <div class="node-actions">
              <el-button link type="primary" @click.stop="handleEdit(data)">Edit</el-button>
              <el-button link type="primary" @click.stop="handleAddChild(data)">Add Child</el-button>
              <el-switch
                :model-value="data.status === 1"
                @change="(val: boolean) => handleStatusChange(data, val)"
              />
              <el-button link type="danger" @click.stop="handleDelete(data)">Delete</el-button>
            </div>
          </div>
        </template>
      </el-tree>
    </el-card>

    <!-- Create/Edit Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="editingCategory ? 'Edit Category' : 'Create Category'"
      width="600px"
      destroy-on-close
    >
      <el-form :model="formData" label-width="120px" ref="formRef">
        <el-form-item label="Name" prop="name" :rules="[{ required: true, message: 'Name is required' }]">
          <el-input v-model="formData.name" placeholder="Enter category name" />
        </el-form-item>
        <el-form-item label="Code">
          <el-input v-model="formData.code" placeholder="Enter category code" />
        </el-form-item>
        <el-form-item label="Icon URL">
          <el-input v-model="formData.icon" placeholder="Enter icon URL" />
        </el-form-item>
        <el-form-item label="Image URL">
          <el-input v-model="formData.image" placeholder="Enter image URL" />
        </el-form-item>
        <el-form-item label="SEO Title">
          <el-input v-model="formData.seo_title" placeholder="Enter SEO title" />
        </el-form-item>
        <el-form-item label="SEO Description">
          <el-input v-model="formData.seo_description" type="textarea" :rows="3" placeholder="Enter SEO description" />
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  getCategoryTree,
  createCategory,
  updateCategory,
  deleteCategory,
  updateCategoryStatus,
  updateCategorySort,
  type CategoryTree,
  type CreateCategoryRequest
} from '@/api/category'

const treeRef = ref()
const categoryTree = ref<CategoryTree[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const submitting = ref(false)
const editingCategory = ref<CategoryTree | null>(null)
const formRef = ref()

const formData = reactive<CreateCategoryRequest & { id?: number }>({
  name: '',
  code: '',
  icon: '',
  image: '',
  seo_title: '',
  seo_description: '',
  sort: 0,
  parent_id: 0
})

const loadTree = async () => {
  loading.value = true
  try {
    const data = await getCategoryTree()
    categoryTree.value = data
  } catch (error) {
    ElMessage.error('Failed to load categories')
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  editingCategory.value = null
  Object.assign(formData, {
    name: '',
    code: '',
    icon: '',
    image: '',
    seo_title: '',
    seo_description: '',
    sort: 0,
    parent_id: 0
  })
  dialogVisible.value = true
}

const handleAddChild = (parent: CategoryTree) => {
  editingCategory.value = null
  Object.assign(formData, {
    name: '',
    code: '',
    icon: '',
    image: '',
    seo_title: '',
    seo_description: '',
    sort: 0,
    parent_id: parent.id
  })
  dialogVisible.value = true
}

const handleEdit = (category: CategoryTree) => {
  editingCategory.value = category
  Object.assign(formData, {
    id: category.id,
    name: category.name,
    code: category.code,
    icon: category.icon,
    image: category.image,
    seo_title: category.seo_title,
    seo_description: category.seo_description,
    sort: category.sort,
    parent_id: category.parent_id
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await formRef.value?.validate()
  submitting.value = true
  try {
    if (editingCategory.value) {
      await updateCategory({ ...formData, id: editingCategory.value.id })
      ElMessage.success('Category updated')
    } else {
      await createCategory(formData)
      ElMessage.success('Category created')
    }
    dialogVisible.value = false
    loadTree()
  } catch (error) {
    ElMessage.error('Failed to save category')
  } finally {
    submitting.value = false
  }
}

const handleStatusChange = async (category: CategoryTree, enabled: boolean) => {
  try {
    await updateCategoryStatus(category.id, enabled ? 1 : 0)
    category.status = enabled ? 1 : 0
    ElMessage.success('Status updated')
  } catch (error) {
    ElMessage.error('Failed to update status')
  }
}

const handleDelete = async (category: CategoryTree) => {
  try {
    await ElMessageBox.confirm(
      `Delete "${category.name}"? Products will be moved to parent category.`,
      'Confirm Delete',
      { type: 'warning' }
    )
    await deleteCategory(category.id)
    ElMessage.success('Category deleted')
    loadTree()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('Failed to delete category')
    }
  }
}

const handleDrop = async () => {
  const nodes = treeRef.value?.store?.nodesMap
  if (!nodes) return

  const sorts: { id: number; sort: number }[] = []
  const collectSorts = (items: CategoryTree[], order: number = 0) => {
    items.forEach((item, index) => {
      sorts.push({ id: item.id, sort: order + index })
      if (item.children?.length) {
        collectSorts(item.children, (order + index + 1) * 100)
      }
    })
  }
  collectSorts(categoryTree.value)

  try {
    await updateCategorySort(sorts)
    ElMessage.success('Sort order updated')
  } catch (error) {
    ElMessage.error('Failed to update sort order')
    loadTree()
  }
}

onMounted(loadTree)
</script>

<style scoped>
.category-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.tree-node {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-right: 10px;
}

.node-content {
  display: flex;
  align-items: center;
  gap: 10px;
}

.node-name {
  font-weight: 500;
}

.product-count {
  color: #909399;
  font-size: 12px;
}

.node-actions {
  display: flex;
  align-items: center;
  gap: 5px;
}
</style>
```

- [ ] **Step 2: Add route (if needed)**

Check `shop-admin/src/router/index.ts` and add route for categories if not exists.

- [ ] **Step 3: Commit**

```bash
git add shop-admin/src/views/categories/index.vue
git commit -m "feat(frontend): add Category management page with tree view"
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
# Create category
curl -X POST http://localhost:8888/api/v1/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name":"Electronics","code":"electronics"}'

# Get category tree
curl http://localhost:8888/api/v1/categories/tree \
  -H "Authorization: Bearer <token>"
```

Expected: Successful responses

- [ ] **Step 4: Test frontend**

Run: `cd shop-admin && npm run dev`

Expected: Category page loads at `/categories`, tree displays correctly

---

## Rollback

```sql
-- Remove SEO fields
ALTER TABLE `categories`
    DROP COLUMN `seo_title`,
    DROP COLUMN `seo_description`;

-- Drop category_markets table
DROP TABLE IF EXISTS `category_markets`;
```