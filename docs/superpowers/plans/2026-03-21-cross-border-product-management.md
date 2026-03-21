# Cross-Border Product Management Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a cross-border product management system with multi-market support, allowing sellers to manage products across US, UK, DE, FR, AU markets with per-market pricing, visibility control, and localization.

**Architecture:** DDD layered architecture (domain → infrastructure → handler). Backend extends existing Product/Market APIs, adds ProductMarket API. Frontend adds Market Management page and Product Detail page with 6 tabs.

**Tech Stack:** Go + go-zero + GORM (backend), Vue 3 + TypeScript + Element Plus (frontend)

---

## File Structure Overview

### Backend Files to Create/Modify
```
admin/
├── desc/
│   ├── product.api                    # MODIFY: Add compliance fields to DTOs
│   └── product_market.api             # CREATE: New API endpoints
├── internal/
│   ├── types/
│   │   └── types.go                   # AUTO-GEN after make api
│   ├── handler/
│   │   ├── products/                  # MODIFY: Update handlers
│   │   └── product_markets/           # CREATE: New handlers
│   └── logic/
│       ├── products/                  # MODIFY: Update logic
│       └── product_markets/           # CREATE: New logic
```

### Frontend Files to Create/Modify
```
shop-admin/src/
├── api/
│   ├── product.ts                     # MODIFY: Add compliance fields
│   └── market.ts                      # CREATE: Market API
├── views/
│   ├── products/
│   │   ├── index.vue                  # MODIFY: Add market filter
│   │   └── [id]/
│   │       └── index.vue              # CREATE: Product detail page
│   └── settings/
│       └── markets/
│           └── index.vue              # CREATE: Market management
├── router/index.ts                    # MODIFY: Add routes
└── utils/types.ts                     # MODIFY: Add types
```

---

## Phase 1: Core Infrastructure

### Task 1.1: Extend Product API Definition

**Files:**
- Modify: `admin/desc/product.api`

- [ ] **Step 1: Add compliance fields to CreateProductReq**

Open `admin/desc/product.api` and extend `CreateProductReq`:

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
    HSCode         string   `json:"hs_code,optional"`
    COO            string   `json:"coo,optional"`
    Weight         string   `json:"weight,optional"`
    WeightUnit     string   `json:"weight_unit,optional"`
    Length         string   `json:"length,optional"`
    Width          string   `json:"width,optional"`
    Height         string   `json:"height,optional"`
    DangerousGoods []string `json:"dangerous_goods,optional"`
}
```

- [ ] **Step 2: Extend UpdateProductReq with same fields**

Add the same new fields to `UpdateProductReq`:

```go
type UpdateProductReq {
    ID          int64  `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description,optional"`
    Price       int64  `json:"price"`
    Currency    string `json:"currency,optional"`
    CategoryID  int64  `json:"category_id"`

    // New fields
    SKU            string   `json:"sku,optional"`
    Brand          string   `json:"brand,optional"`
    Tags           []string `json:"tags,optional"`
    Images         []string `json:"images,optional"`
    IsMatrixProduct bool     `json:"is_matrix_product,optional"`

    // Compliance fields
    HSCode         string   `json:"hs_code,optional"`
    COO            string   `json:"coo,optional"`
    Weight         string   `json:"weight,optional"`
    WeightUnit     string   `json:"weight_unit,optional"`
    Length         string   `json:"length,optional"`
    Width          string   `json:"width,optional"`
    Height         string   `json:"height,optional"`
    DangerousGoods []string `json:"dangerous_goods,optional"`
}
```

- [ ] **Step 3: Extend ProductDetailResp**

Add new fields to `ProductDetailResp`:

```go
type ProductDetailResp {
    ID          int64  `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Price       int64  `json:"price"`
    Currency    string `json:"currency"`
    CostPrice   int64  `json:"cost_price"`
    Stock       int    `json:"stock"`
    Status      string `json:"status"`
    CategoryID  int64  `json:"category_id"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`

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

    // Market info
    Markets        []ProductMarketInfo `json:"markets,optional"`
}

type ProductMarketInfo {
    MarketID   int64  `json:"market_id"`
    MarketCode string `json:"market_code"`
    MarketName string `json:"market_name"`
    IsEnabled  bool   `json:"is_enabled"`
    Price      string `json:"price"`
    Currency   string `json:"currency"`
}
```

- [ ] **Step 4: Add market_id filter to ListProductReq**

```go
type ListProductReq {
    Name       string `form:"name,optional"`
    CategoryID int64  `form:"category_id,optional"`
    Status     string `form:"status,optional"`
    MinPrice   int64  `form:"min_price,optional"`
    MaxPrice   int64  `form:"max_price,optional"`
    Page       int    `form:"page,default=1"`
    PageSize   int    `form:"page_size,default=20"`

    // New
    MarketID   int64  `form:"market_id,optional"`
}
```

- [ ] **Step 5: Run make api to regenerate code**

```bash
cd admin && make api
```

Expected: No errors, types.go and routes.go updated

- [ ] **Step 6: Commit**

```bash
git add admin/desc/product.api admin/internal/types/ admin/internal/handler/routes.go
git commit -m "feat(api): extend product API with compliance and market fields"
```

---

### Task 1.2: Create ProductMarket API Definition

**Files:**
- Create: `admin/desc/product_market.api`

- [ ] **Step 1: Create product_market.api file**

Create `admin/desc/product_market.api`:

```go
syntax = "v1"

info (
    title:   "Product Market API"
    desc:    "商品市场关联管理接口"
    version: "v1"
)

type (
    // ProductMarket response
    ProductMarketResp {
        ID                  int64   `json:"id"`
        ProductID           int64   `json:"product_id"`
        MarketID            int64   `json:"market_id"`
        MarketCode          string  `json:"market_code"`
        MarketName          string  `json:"market_name"`
        IsEnabled           bool    `json:"is_enabled"`
        Price               string  `json:"price"`
        CompareAtPrice      string  `json:"compare_at_price,optional"`
        Currency            string  `json:"currency"`
        StockAlertThreshold int     `json:"stock_alert_threshold"`
        PublishedAt         string  `json:"published_at,optional"`
    }

    // List product markets
    ListProductMarketsResp {
        List []*ProductMarketResp `json:"list"`
    }

    // Update product market
    UpdateProductMarketReq {
        ProductID           int64   `path:"id"`
        MarketID            int64   `path:"market_id"`
        IsEnabled           *bool   `json:"is_enabled,optional"`
        Price               string  `json:"price,optional"`
        CompareAtPrice      string  `json:"compare_at_price,optional"`
        StockAlertThreshold int     `json:"stock_alert_threshold,optional"`
    }

    // Push to market request
    PushToMarketReq {
        ProductID  int64   `path:"id"`
        MarketIDs  []int64 `json:"market_ids"`
        Prices     []string `json:"prices"`  // Price per market, same order as market_ids
    }

    PushToMarketResp {
        Success  []int64 `json:"success"`   // Market IDs successfully added
        Failed   []int64 `json:"failed"`    // Market IDs that failed
    }

    // Remove from market
    RemoveFromMarketReq {
        ProductID int64 `path:"id"`
        MarketID  int64 `path:"market_id"`
    }
)

@server (
    group:      product_markets
    middleware: AuthMiddleware
)
service admin-api {
    @doc "获取商品市场配置列表"
    @handler ListProductMarketsHandler
    get /api/v1/products/:id/markets returns (ListProductMarketsResp)

    @doc "更新商品市场配置"
    @handler UpdateProductMarketHandler
    put /api/v1/products/:id/markets/:market_id (UpdateProductMarketReq) returns (ProductMarketResp)

    @doc "推送商品到市场"
    @handler PushToMarketHandler
    post /api/v1/products/:id/push-to-market (PushToMarketReq) returns (PushToMarketResp)

    @doc "从市场移除商品"
    @handler RemoveFromMarketHandler
    delete /api/v1/products/:id/markets/:market_id (RemoveFromMarketReq)
}
```

- [ ] **Step 2: Import in admin.api**

Open `admin/desc/admin.api` and add import:

```go
import "product_market.api"
```

- [ ] **Step 3: Run make api**

```bash
cd admin && make api
```

Expected: New handler files created in `admin/internal/handler/product_markets/`

- [ ] **Step 4: Commit**

```bash
git add admin/desc/product_market.api admin/desc/admin.api admin/internal/handler/product_markets/ admin/internal/types/
git commit -m "feat(api): add product market API definition"
```

---

### Task 1.3: Implement ProductMarket Handlers

**Files:**
- Create: `admin/internal/logic/product_markets/list_product_markets_logic.go`
- Create: `admin/internal/logic/product_markets/update_product_market_logic.go`
- Create: `admin/internal/logic/product_markets/push_to_market_logic.go`
- Create: `admin/internal/logic/product_markets/remove_from_market_logic.go`

- [ ] **Step 1: Implement ListProductMarketsLogic**

Create `admin/internal/logic/product_markets/list_product_markets_logic.go`:

```go
package product_markets

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

type ListProductMarketsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListProductMarketsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductMarketsLogic {
	return &ListProductMarketsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListProductMarketsLogic) ListProductMarkets(req *types.ListProductMarketsReq) (resp *types.ListProductMarketsResp, err error) {
	db := l.svcCtx.DB
	repo := persistence.NewProductMarketRepository()

	productMarkets, err := repo.FindByProductID(l.ctx, db, req.ProductID)
	if err != nil {
		return nil, err
	}

	// Get market info
	marketRepo := persistence.NewMarketRepository()
	markets, err := marketRepo.FindAll(l.ctx, db)
	if err != nil {
		return nil, err
	}

	marketMap := make(map[int64]*types.MarketResponse)
	for _, m := range markets {
		marketMap[m.ID] = &types.MarketResponse{
			ID:              m.ID,
			Code:            m.Code,
			Name:            m.Name,
			Currency:        m.Currency,
			DefaultLanguage: m.DefaultLanguage,
			Flag:            m.Flag,
			IsActive:        m.IsActive,
			IsDefault:       m.IsDefault,
		}
	}

	list := make([]*types.ProductMarketResp, 0, len(productMarkets))
	for _, pm := range productMarkets {
		market, ok := marketMap[pm.MarketID]
		if !ok {
			continue
		}

		var compareAtPrice string
		if pm.CompareAtPrice != nil {
			compareAtPrice = pm.CompareAtPrice.String()
		}

		var publishedAt string
		if pm.PublishedAt != nil {
			publishedAt = pm.PublishedAt.Format("2006-01-02 15:04:05")
		}

		list = append(list, &types.ProductMarketResp{
			ID:                  pm.ID,
			ProductID:           pm.ProductID,
			MarketID:            pm.MarketID,
			MarketCode:          market.Code,
			MarketName:          market.Name,
			IsEnabled:           pm.IsEnabled,
			Price:               pm.Price.String(),
			CompareAtPrice:      compareAtPrice,
			Currency:            market.Currency,
			StockAlertThreshold: pm.StockAlertThreshold,
			PublishedAt:         publishedAt,
		})
	}

	return &types.ListProductMarketsResp{List: list}, nil
}
```

- [ ] **Step 2: Implement PushToMarketLogic**

Create `admin/internal/logic/product_markets/push_to_market_logic.go`:

```go
package product_markets

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
)

type PushToMarketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPushToMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushToMarketLogic {
	return &PushToMarketLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PushToMarketLogic) PushToMarket(req *types.PushToMarketReq) (resp *types.PushToMarketResp, err error) {
	db := l.svcCtx.DB
	repo := persistence.NewProductMarketRepository()
	marketRepo := persistence.NewMarketRepository()

	// Get tenant ID from context
	tenantID := l.ctx.Value("tenant_id")
	var tid int64
	if tidVal, ok := tenantID.(int64); ok {
		tid = tidVal
	}

	var success, failed []int64

	for i, marketID := range req.MarketIDs {
		// Validate market exists
		market, err := marketRepo.FindByID(l.ctx, db, marketID)
		if err != nil || !market.IsActive {
			failed = append(failed, marketID)
			continue
		}

		// Check if already exists
		existing, _ := repo.FindByProductAndMarket(l.ctx, db, req.ProductID, marketID, nil)
		if existing != nil {
			failed = append(failed, marketID)
			continue
		}

		// Parse price
		var price decimal.Decimal
		if i < len(req.Prices) {
			price, _ = decimal.NewFromString(req.Prices[i])
		}

		// Create ProductMarket
		pm := &product.ProductMarket{
			TenantID:  tid,
			ProductID: req.ProductID,
			MarketID:  marketID,
			Price:     price,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := repo.Create(l.ctx, db, pm); err != nil {
			failed = append(failed, marketID)
			continue
		}

		success = append(success, marketID)
	}

	return &types.PushToMarketResp{
		Success: success,
		Failed:  failed,
	}, nil
}
```

- [ ] **Step 3: Implement UpdateProductMarketLogic**

Create `admin/internal/logic/product_markets/update_product_market_logic.go`:

```go
package product_markets

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/shopspring/decimal"
)

type UpdateProductMarketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProductMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductMarketLogic {
	return &UpdateProductMarketLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProductMarketLogic) UpdateProductMarket(req *types.UpdateProductMarketReq) (resp *types.ProductMarketResp, err error) {
	db := l.svcCtx.DB
	repo := persistence.NewProductMarketRepository()

	pm, err := repo.FindByProductAndMarket(l.ctx, db, req.ProductID, req.MarketID, nil)
	if err != nil {
		return nil, err
	}

	if req.IsEnabled != nil {
		pm.IsEnabled = *req.IsEnabled
	}

	if req.Price != "" {
		pm.Price, _ = decimal.NewFromString(req.Price)
	}

	if req.CompareAtPrice != "" {
		cap, _ := decimal.NewFromString(req.CompareAtPrice)
		pm.CompareAtPrice = &cap
	}

	pm.StockAlertThreshold = req.StockAlertThreshold
	pm.UpdatedAt = time.Now()

	if err := repo.Update(l.ctx, db, pm); err != nil {
		return nil, err
	}

	// Get market info for response
	marketRepo := persistence.NewMarketRepository()
	market, err := marketRepo.FindByID(l.ctx, db, pm.MarketID)
	if err != nil {
		return nil, err
	}

	var compareAtPrice string
	if pm.CompareAtPrice != nil {
		compareAtPrice = pm.CompareAtPrice.String()
	}

	var publishedAt string
	if pm.PublishedAt != nil {
		publishedAt = pm.PublishedAt.Format("2006-01-02 15:04:05")
	}

	return &types.ProductMarketResp{
		ID:                  pm.ID,
		ProductID:           pm.ProductID,
		MarketID:            pm.MarketID,
		MarketCode:          market.Code,
		MarketName:          market.Name,
		IsEnabled:           pm.IsEnabled,
		Price:               pm.Price.String(),
		CompareAtPrice:      compareAtPrice,
		Currency:            market.Currency,
		StockAlertThreshold: pm.StockAlertThreshold,
		PublishedAt:         publishedAt,
	}, nil
}
```

- [ ] **Step 4: Implement RemoveFromMarketLogic**

Create `admin/internal/logic/product_markets/remove_from_market_logic.go`:

```go
package product_markets

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

type RemoveFromMarketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveFromMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveFromMarketLogic {
	return &RemoveFromMarketLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveFromMarketLogic) RemoveFromMarket(req *types.RemoveFromMarketReq) error {
	db := l.svcCtx.DB
	repo := persistence.NewProductMarketRepository()

	pm, err := repo.FindByProductAndMarket(l.ctx, db, req.ProductID, req.MarketID, nil)
	if err != nil {
		return err
	}

	return repo.Delete(l.ctx, db, pm.ID)
}
```

- [ ] **Step 5: Build and verify**

```bash
cd admin && make build
```

Expected: Build successful

- [ ] **Step 6: Commit**

```bash
git add admin/internal/logic/product_markets/
git commit -m "feat(product-market): implement product market API handlers"
```

---

### Task 1.4: Create Market API Service (Frontend)

**Files:**
- Create: `shop-admin/src/api/market.ts`

- [ ] **Step 1: Create market.ts API file**

Create `shop-admin/src/api/market.ts`:

```typescript
import request from '@/utils/request'

export interface Market {
  id: number
  code: string
  name: string
  currency: string
  default_language: string
  flag: string
  is_active: boolean
  is_default: boolean
  tax_rules: {
    vat_rate: string
    gst_rate: string
    ioss_enabled: boolean
    include_tax: boolean
  }
  created_at: string
  updated_at: string
}

export interface ListMarketsResponse {
  list: Market[]
  total: number
}

export interface CreateMarketRequest {
  code: string
  name: string
  currency: string
  default_language?: string
  flag?: string
  tax_rules?: {
    vat_rate?: string
    gst_rate?: string
    ioss_enabled?: boolean
    include_tax?: boolean
  }
}

export interface UpdateMarketRequest {
  id: number
  name?: string
  is_active?: boolean
  tax_rules?: {
    vat_rate?: string
    gst_rate?: string
    ioss_enabled?: boolean
    include_tax?: boolean
  }
}

export function getMarkets() {
  return request<ListMarketsResponse>({
    url: '/api/v1/markets',
    method: 'get'
  })
}

export function getMarket(id: number) {
  return request<Market>({
    url: `/api/v1/markets/${id}`,
    method: 'get'
  })
}

export function createMarket(data: CreateMarketRequest) {
  return request<Market>({
    url: '/api/v1/markets',
    method: 'post',
    data
  })
}

export function updateMarket(data: UpdateMarketRequest) {
  return request<Market>({
    url: `/api/v1/markets/${data.id}`,
    method: 'put',
    data
  })
}

export function deleteMarket(id: number) {
  return request({
    url: `/api/v1/markets/${id}`,
    method: 'delete'
  })
}
```

- [ ] **Step 2: Commit**

```bash
git add shop-admin/src/api/market.ts
git commit -m "feat(frontend): add market API service"
```

---

### Task 1.5: Create Market Management Page (Frontend)

**Files:**
- Create: `shop-admin/src/views/settings/markets/index.vue`
- Modify: `shop-admin/src/router/index.ts`

- [ ] **Step 1: Create markets directory**

```bash
mkdir -p shop-admin/src/views/settings/markets
```

- [ ] **Step 2: Create Market Management page**

Create `shop-admin/src/views/settings/markets/index.vue`:

```vue
<template>
  <div class="markets-page">
    <!-- Header -->
    <el-card class="header-card" shadow="never">
      <div class="header-bar">
        <h2>市场管理</h2>
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>新增市场
        </el-button>
      </div>
    </el-card>

    <!-- Markets Table -->
    <el-card class="table-card" shadow="never">
      <el-table :data="marketList" v-loading="loading" stripe>
        <el-table-column label="市场" min-width="200">
          <template #default="{ row }">
            <div class="market-cell">
              <span class="market-flag">{{ row.flag }}</span>
              <div class="market-info">
                <p class="market-name">{{ row.name }}</p>
                <p class="market-code">{{ row.code }}</p>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="currency" label="货币" width="100" align="center" />
        <el-table-column prop="default_language" label="默认语言" width="120" align="center" />
        <el-table-column label="税配置" min-width="180">
          <template #default="{ row }">
            <div class="tax-info">
              <span v-if="row.tax_rules?.vat_rate">VAT: {{ row.tax_rules.vat_rate }}%</span>
              <span v-if="row.tax_rules?.gst_rate">GST: {{ row.tax_rules.gst_rate }}%</span>
              <span v-if="!row.tax_rules?.vat_rate && !row.tax_rules?.gst_rate">-</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.is_active ? 'success' : 'info'" size="small">
              {{ row.is_active ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="主市场" width="80" align="center">
          <template #default="{ row }">
            <el-icon v-if="row.is_default" class="primary-icon"><Star /></el-icon>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button
              type="danger"
              link
              size="small"
              @click="handleDelete(row)"
              :disabled="row.is_default"
            >删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Add/Edit Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑市场' : '新增市场'"
      width="600px"
      destroy-on-close
    >
      <el-form :model="marketForm" label-width="100px" :rules="formRules" ref="formRef">
        <el-form-item label="市场代码" prop="code">
          <el-select
            v-model="marketForm.code"
            placeholder="选择市场"
            style="width: 100%"
            :disabled="isEdit"
          >
            <el-option label="🇺🇸 美国 (US)" value="US" />
            <el-option label="🇬🇧 英国 (UK)" value="UK" />
            <el-option label="🇩🇪 德国 (DE)" value="DE" />
            <el-option label="🇫🇷 法国 (FR)" value="FR" />
            <el-option label="🇦🇺 澳大利亚 (AU)" value="AU" />
          </el-select>
        </el-form-item>
        <el-form-item label="市场名称" prop="name">
          <el-input v-model="marketForm.name" placeholder="如: United States" />
        </el-form-item>
        <el-form-item label="货币" prop="currency">
          <el-select v-model="marketForm.currency" placeholder="选择货币" style="width: 100%">
            <el-option label="USD - 美元" value="USD" />
            <el-option label="GBP - 英镑" value="GBP" />
            <el-option label="EUR - 欧元" value="EUR" />
            <el-option label="AUD - 澳元" value="AUD" />
            <el-option label="CNY - 人民币" value="CNY" />
          </el-select>
        </el-form-item>
        <el-form-item label="默认语言">
          <el-select v-model="marketForm.default_language" style="width: 100%">
            <el-option label="English" value="en" />
            <el-option label="Deutsch" value="de" />
            <el-option label="Français" value="fr" />
            <el-option label="中文" value="zh" />
          </el-select>
        </el-form-item>
        <el-form-item label="启用状态">
          <el-switch v-model="marketForm.is_active" />
        </el-form-item>

        <el-divider>税务配置</el-divider>

        <el-form-item label="VAT税率">
          <el-input-number
            v-model="marketForm.tax_rules.vat_rate"
            :min="0"
            :max="100"
            :precision="2"
            style="width: 200px"
          />
          <span class="input-suffix">%</span>
        </el-form-item>
        <el-form-item label="GST税率">
          <el-input-number
            v-model="marketForm.tax_rules.gst_rate"
            :min="0"
            :max="100"
            :precision="2"
            style="width: 200px"
          />
          <span class="input-suffix">%</span>
        </el-form-item>
        <el-form-item label="IOSS启用">
          <el-switch v-model="marketForm.tax_rules.ioss_enabled" />
          <span class="form-hint">低值商品进口方案</span>
        </el-form-item>
        <el-form-item label="含税显示">
          <el-switch v-model="marketForm.tax_rules.include_tax" />
          <span class="form-hint">价格是否含税显示</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="saveLoading">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Star } from '@element-plus/icons-vue'
import { getMarkets, createMarket, updateMarket, deleteMarket, type Market } from '@/api/market'

const loading = ref(false)
const saveLoading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()

const marketList = ref<Market[]>([])

const marketForm = reactive({
  id: 0,
  code: '',
  name: '',
  currency: 'USD',
  default_language: 'en',
  is_active: true,
  tax_rules: {
    vat_rate: 0,
    gst_rate: 0,
    ioss_enabled: false,
    include_tax: false
  }
})

const formRules = {
  code: [{ required: true, message: '请选择市场', trigger: 'change' }],
  name: [{ required: true, message: '请输入市场名称', trigger: 'blur' }],
  currency: [{ required: true, message: '请选择货币', trigger: 'change' }]
}

const loadMarkets = async () => {
  loading.value = true
  try {
    const res = await getMarkets()
    marketList.value = res.list || []
  } catch (error) {
    ElMessage.error('加载市场列表失败')
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  isEdit.value = false
  Object.assign(marketForm, {
    id: 0,
    code: '',
    name: '',
    currency: 'USD',
    default_language: 'en',
    is_active: true,
    tax_rules: {
      vat_rate: 0,
      gst_rate: 0,
      ioss_enabled: false,
      include_tax: false
    }
  })
  dialogVisible.value = true
}

const handleEdit = (row: Market) => {
  isEdit.value = true
  Object.assign(marketForm, {
    id: row.id,
    code: row.code,
    name: row.name,
    currency: row.currency,
    default_language: row.default_language,
    is_active: row.is_active,
    tax_rules: {
      vat_rate: parseFloat(row.tax_rules?.vat_rate || '0'),
      gst_rate: parseFloat(row.tax_rules?.gst_rate || '0'),
      ioss_enabled: row.tax_rules?.ioss_enabled || false,
      include_tax: row.tax_rules?.include_tax || false
    }
  })
  dialogVisible.value = true
}

const handleDelete = (row: Market) => {
  ElMessageBox.confirm(`确认删除市场 "${row.name}"?`, '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteMarket(row.id)
      ElMessage.success('删除成功')
      loadMarkets()
    } catch (error) {
      ElMessage.error('删除失败')
    }
  })
}

const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      saveLoading.value = true
      try {
        const data = {
          code: marketForm.code,
          name: marketForm.name,
          currency: marketForm.currency,
          default_language: marketForm.default_language,
          tax_rules: {
            vat_rate: marketForm.tax_rules.vat_rate.toString(),
            gst_rate: marketForm.tax_rules.gst_rate.toString(),
            ioss_enabled: marketForm.tax_rules.ioss_enabled,
            include_tax: marketForm.tax_rules.include_tax
          }
        }

        if (isEdit.value) {
          await updateMarket({ id: marketForm.id, ...data })
          ElMessage.success('更新成功')
        } else {
          await createMarket(data)
          ElMessage.success('创建成功')
        }

        dialogVisible.value = false
        loadMarkets()
      } catch (error) {
        ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
      } finally {
        saveLoading.value = false
      }
    }
  })
}

onMounted(() => {
  loadMarkets()
})
</script>

<style scoped>
.markets-page {
  padding: 0;
}

.header-card {
  margin-bottom: 20px;
}

.header-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-bar h2 {
  margin: 0;
  font-size: 20px;
}

.market-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.market-flag {
  font-size: 24px;
}

.market-info {
  flex: 1;
}

.market-name {
  font-weight: 500;
  margin: 0;
}

.market-code {
  font-size: 12px;
  color: #909399;
  margin: 0;
}

.tax-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 13px;
}

.primary-icon {
  color: #e6a23c;
  font-size: 18px;
}

.form-hint {
  margin-left: 12px;
  font-size: 12px;
  color: #909399;
}

.input-suffix {
  margin-left: 8px;
  color: #606266;
}
</style>
```

- [ ] **Step 3: Update router**

Add route to `shop-admin/src/router/index.ts`:

```typescript
// Add to children array:
{
  path: 'settings/markets',
  name: 'settings-markets',
  component: () => import('@/views/settings/markets/index.vue'),
  meta: { title: '市场管理' }
}
```

- [ ] **Step 4: Commit**

```bash
git add shop-admin/src/views/settings/ shop-admin/src/router/index.ts
git commit -m "feat(frontend): add market management page"
```

---

## Phase 2: Product Management

### Task 2.1: Update Product List Handler for Market Filter

**Files:**
- Modify: `admin/internal/logic/products/list_product_logic.go`

- [ ] **Step 1: Update ListProductLogic to support market filter**

Modify `admin/internal/logic/products/list_product_logic.go`:

Add import and update the logic to filter by market and include market info in response.

Key changes:
1. If `MarketID` is provided in request, join with `product_markets` table
2. Add `Markets` field to each product in response

- [ ] **Step 2: Build and test**

```bash
cd admin && make build
```

- [ ] **Step 3: Commit**

```bash
git add admin/internal/logic/products/list_product_logic.go
git commit -m "feat(product): add market filter to product list"
```

---

### Task 2.2: Update Product API Service (Frontend)

**Files:**
- Modify: `shop-admin/src/api/product.ts`

- [ ] **Step 1: Extend product API with compliance fields**

Update `shop-admin/src/api/product.ts`:

```typescript
import request from '@/utils/request'

export interface ProductMarketInfo {
  market_id: number
  market_code: string
  market_name: string
  is_enabled: boolean
  price: string
  currency: string
}

export interface Product {
  id: number
  name: string
  description: string
  price: number
  currency: string
  cost_price: number
  stock: number
  status: string
  category_id: number
  created_at: string
  updated_at: string
  // New fields
  sku: string
  brand: string
  tags: string[]
  images: string[]
  is_matrix_product: boolean
  // Compliance
  hs_code: string
  coo: string
  weight: string
  weight_unit: string
  length: string
  width: string
  height: string
  dangerous_goods: string[]
  // Markets
  markets: ProductMarketInfo[]
}

export interface CreateProductRequest {
  name: string
  description?: string
  price: number
  currency?: string
  cost_price?: number
  category_id: number
  // New fields
  sku?: string
  brand?: string
  tags?: string[]
  images?: string[]
  is_matrix_product?: boolean
  // Compliance
  hs_code?: string
  coo?: string
  weight?: string
  weight_unit?: string
  length?: string
  width?: string
  height?: string
  dangerous_goods?: string[]
}

export interface ListProductsParams {
  page: number
  page_size: number
  name?: string
  category_id?: number
  status?: string
  min_price?: number
  max_price?: number
  market_id?: number  // New
}

export function getProductList(params: ListProductsParams) {
  return request<{ list: Product[]; total: number; page: number; page_size: number }>({
    url: '/api/v1/products',
    method: 'get',
    params
  })
}

export function createProduct(data: CreateProductRequest) {
  return request<{ id: number }>({
    url: '/api/v1/products',
    method: 'post',
    data
  })
}

export function getProduct(id: number) {
  return request<Product>({
    url: `/api/v1/products/${id}`,
    method: 'get'
  })
}

export function updateProduct(id: number, data: Partial<CreateProductRequest>) {
  return request<Product>({
    url: `/api/v1/products/${id}`,
    method: 'put',
    data: { id, ...data }
  })
}

export function putOnSale(id: number) {
  return request<Product>({
    url: `/api/v1/products/${id}/on-sale`,
    method: 'post',
    data: {}
  })
}

export function takeOffSale(id: number) {
  return request<Product>({
    url: `/api/v1/products/${id}/off-sale`,
    method: 'post',
    data: {}
  })
}

export function updateStock(id: number, quantity: number) {
  return request<Product>({
    url: `/api/v1/products/${id}/stock`,
    method: 'put',
    data: { quantity }
  })
}

// Product Market APIs
export interface ProductMarket {
  id: number
  product_id: number
  market_id: number
  market_code: string
  market_name: string
  is_enabled: boolean
  price: string
  compare_at_price: string
  currency: string
  stock_alert_threshold: number
  published_at: string
}

export function getProductMarkets(productId: number) {
  return request<{ list: ProductMarket[] }>({
    url: `/api/v1/products/${productId}/markets`,
    method: 'get'
  })
}

export function updateProductMarket(productId: number, marketId: number, data: {
  is_enabled?: boolean
  price?: string
  compare_at_price?: string
  stock_alert_threshold?: number
}) {
  return request<ProductMarket>({
    url: `/api/v1/products/${productId}/markets/${marketId}`,
    method: 'put',
    data
  })
}

export function pushToMarket(productId: number, data: {
  market_ids: number[]
  prices: string[]
}) {
  return request<{ success: number[]; failed: number[] }>({
    url: `/api/v1/products/${productId}/push-to-market`,
    method: 'post',
    data
  })
}

export function removeFromMarket(productId: number, marketId: number) {
  return request({
    url: `/api/v1/products/${productId}/markets/${marketId}`,
    method: 'delete'
  })
}
```

- [ ] **Step 2: Commit**

```bash
git add shop-admin/src/api/product.ts
git commit -m "feat(frontend): extend product API with compliance and market fields"
```

---

### Task 2.3: Update Product List Page with Market Filter

**Files:**
- Modify: `shop-admin/src/views/products/index.vue`

- [ ] **Step 1: Add market filter bar**

Add above the search bar:

```vue
<!-- Market Filter Bar -->
<div class="market-filter-bar">
  <el-radio-group v-model="selectedMarket" @change="handleMarketChange">
    <el-radio-button :value="0">全部市场</el-radio-button>
    <el-radio-button
      v-for="market in markets"
      :key="market.id"
      :value="market.id"
    >
      {{ market.flag }} {{ market.code }}
    </el-radio-button>
  </el-radio-group>
</div>
```

- [ ] **Step 2: Add market visibility column**

Add to table columns:

```vue
<el-table-column label="市场" min-width="180">
  <template #default="{ row }">
    <div class="market-tags">
      <el-tag
        v-for="m in row.markets"
        :key="m.market_id"
        :type="m.is_enabled ? 'success' : 'info'"
        size="small"
        class="market-tag"
      >
        {{ m.market_code }} {{ m.is_enabled ? '✓' : '✗' }}
      </el-tag>
    </div>
  </template>
</el-table-column>
```

- [ ] **Step 3: Update price column to show market currency**

Modify price column to show price in selected market's currency when a market is selected.

- [ ] **Step 4: Add batch push to market action**

Add to bulk actions:

```vue
<el-button size="small" type="primary" @click="handleBatchPushToMarket">
  推送到市场
</el-button>
```

- [ ] **Step 5: Implement Push to Market dialog**

Create a dialog component for pushing products to markets with price input.

- [ ] **Step 6: Connect to real API**

Replace mock data with API calls.

- [ ] **Step 7: Commit**

```bash
git add shop-admin/src/views/products/index.vue
git commit -m "feat(frontend): add market filter and push to market to product list"
```

---

### Task 2.4: Create Product Detail Page

**Files:**
- Create: `shop-admin/src/views/products/[id]/index.vue`
- Modify: `shop-admin/src/router/index.ts`

- [ ] **Step 1: Create product detail directory**

```bash
mkdir -p shop-admin/src/views/products/[id]
```

- [ ] **Step 2: Create product detail page with tabs**

Create `shop-admin/src/views/products/[id]/index.vue` with:
- Tab navigation (6 tabs)
- Basic Info tab with compliance form
- Markets tab with enable/disable
- Skeleton for remaining tabs

- [ ] **Step 3: Add route**

```typescript
{
  path: 'products/:id',
  name: 'product-detail',
  component: () => import('@/views/products/[id]/index.vue'),
  meta: { title: '商品详情' }
}
```

- [ ] **Step 4: Commit**

```bash
git add shop-admin/src/views/products/ shop-admin/src/router/index.ts
git commit -m "feat(frontend): create product detail page with tabs"
```

---

## Phase 3-6 Summary

Due to the length of this plan, the remaining phases follow similar patterns:

### Phase 3: Variants & Pricing
- Create `ProductVariant` entity and repository
- Create `variant.api` definition
- Implement variant handlers
- Add Variants Tab and Pricing Tab to frontend

### Phase 4: Localization
- Create `ProductLocalization` entity and repository
- Create `localization.api` definition
- Implement localization handlers
- Add Localization Tab to frontend

### Phase 5: Inventory & Batch
- Create `StockLog` entity and repository
- Create `inventory.api` definition
- Implement inventory handlers and batch operations
- Add Inventory Tab and batch actions to frontend

### Phase 6: Dashboard
- Create dashboard statistics API
- Update Dashboard page with multi-market charts

---

## Testing Checklist

- [ ] Backend: All API endpoints return correct responses
- [ ] Backend: Market filter works in product list
- [ ] Backend: Push to market creates correct records
- [ ] Frontend: Market management page CRUD works
- [ ] Frontend: Product list market filter works
- [ ] Frontend: Product detail page loads and saves data
- [ ] Integration: Full push to market flow works end-to-end

---

## Notes

1. **Database**: SQL files already exist in `sql/` directory. Run `sql/market.sql` and `sql/product.sql` to create tables.

2. **Authentication**: All API endpoints use `AuthMiddleware`. Ensure tenant ID is passed in context.

3. **Currency Handling**: Prices are stored as decimals. No automatic FX conversion - sellers set prices manually per market.

4. **Backward Compatibility**: All new fields in request DTOs are optional to not break existing clients.