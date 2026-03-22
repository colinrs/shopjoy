# ShopJoy API Reference

> Multi-tenant e-commerce SaaS platform API documentation

## Overview

ShopJoy is a multi-tenant e-commerce platform built with Go and the go-zero framework. The Admin API provides endpoints for managing products, categories, brands, markets, and inventory.

**Base URL:** `https://api.shopjoy.com`

**Authentication:** JWT Bearer Token (all endpoints except auth)

**Content-Type:** `application/json`

---

## Authentication

### Login

```
POST /api/v1/auth/login
```

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| email | string | Yes | User email |
| password | string | Yes | User password |

**Response:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expire": "2024-03-22T10:00:00Z"
}
```

### Register

```
POST /api/v1/auth/register
```

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| email | string | Yes | User email |
| password | string | Yes | User password |
| name | string | Yes | User name |

---

## Products

### Create Product

```
POST /api/v1/products
```

**Authorization:** Required

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | Product name |
| description | string | No | Product description |
| price | int64 | Yes | Price in cents (e.g., 129900 = $1,299.00) |
| currency | string | No | Currency code (default: CNY) |
| cost_price | int64 | No | Cost price in cents |
| category_id | int64 | Yes | Category ID |
| sku | string | No | SKU code |
| brand | string | No | Brand name |
| tags | []string | No | Product tags |
| images | []string | No | Image URLs |
| is_matrix_product | bool | No | Has variants |
| hs_code | string | No | Harmonized System code |
| coo | string | No | Country of origin |
| weight | string | No | Weight value |
| weight_unit | string | No | Weight unit (g, kg, lb) |
| length | string | No | Length (cm) |
| width | string | No | Width (cm) |
| height | string | No | Height (cm) |
| dangerous_goods | []string | No | Dangerous goods identifiers |

**Response:**

```json
{
  "id": 1
}
```

### Get Product

```
GET /api/v1/products/:id
```

**Authorization:** Required

**Response:**

```json
{
  "id": 1,
  "name": "Nike Air Max 270",
  "description": "Nike Air Max 270 运动鞋，舒适透气",
  "price": 129900,
  "currency": "CNY",
  "cost_price": 80000,
  "stock": 100,
  "status": "on_sale",
  "category_id": 2,
  "sku": "SKU-001",
  "brand": "Nike",
  "tags": ["运动", "跑步", "休闲"],
  "images": ["https://cdn.example.com/p1-1.jpg"],
  "is_matrix_product": true,
  "hs_code": "64041100",
  "coo": "CN",
  "weight": "450.00",
  "weight_unit": "g",
  "length": "28.00",
  "width": "18.00",
  "height": "12.00",
  "dangerous_goods": [],
  "markets": [
    {
      "market_id": 1,
      "market_code": "US",
      "market_name": "United States",
      "is_enabled": true,
      "price": "199.99",
      "currency": "USD"
    }
  ],
  "created_at": "2024-03-22T10:00:00Z",
  "updated_at": "2024-03-22T10:00:00Z"
}
```

### List Products

```
GET /api/v1/products
```

**Authorization:** Required

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| name | string | - | Filter by name (partial match) |
| category_id | int64 | - | Filter by category |
| status | string | - | Filter by status (draft/on_sale/off_sale) |
| min_price | int64 | - | Minimum price (cents) |
| max_price | int64 | - | Maximum price (cents) |
| market_id | int64 | - | Filter by market |
| page | int | 1 | Page number |
| page_size | int | 20 | Items per page (max 100) |

**Response:**

```json
{
  "list": [
    { "...": "ProductDetailResp" }
  ],
  "total": 100,
  "page": 1,
  "page_size": 20
}
```

### Update Product

```
PUT /api/v1/products/:id
```

**Authorization:** Required

**Request Body:** Same as Create Product, all fields are required.

### Put On Sale

```
POST /api/v1/products/:id/on-sale
```

**Authorization:** Required

Changes product status to `on_sale`.

### Take Off Sale

```
POST /api/v1/products/:id/off-sale
```

**Authorization:** Required

Changes product status to `off_sale`.

### Update Stock

```
PUT /api/v1/products/:id/stock
```

**Authorization:** Required

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| quantity | int | Yes | New stock quantity |

---

## SKUs (Product Variants)

### Create SKU

```
POST /api/v1/skus
```

**Authorization:** Required

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| product_id | int64 | Yes | Parent product ID |
| code | string | Yes | Unique SKU code |
| price | int64 | Yes | Price in cents |
| currency | string | No | Currency code |
| stock | int | No | Initial stock |
| safety_stock | int | No | Safety stock threshold |
| pre_sale_enabled | bool | No | Enable pre-sale |
| attributes | map | No | Variant attributes (e.g., color, size) |

**Response:**

```json
{
  "id": 1
}
```

### Get SKU

```
GET /api/v1/skus/:id
```

**Authorization:** Required

**Response:**

```json
{
  "id": 1,
  "product_id": 1,
  "code": "SKU-001-BLK-42",
  "price": 129900,
  "currency": "CNY",
  "stock": 30,
  "available_stock": 25,
  "locked_stock": 5,
  "safety_stock": 10,
  "pre_sale_enabled": false,
  "attributes": {
    "颜色": "黑色",
    "尺码": "42"
  },
  "status": "enabled",
  "is_low_stock": false,
  "created_at": "2024-03-22T10:00:00Z",
  "updated_at": "2024-03-22T10:00:00Z"
}
```

### List SKUs by Product

```
GET /api/v1/products/:product_id/skus
```

**Authorization:** Required

### Update SKU

```
PUT /api/v1/skus/:id
```

**Authorization:** Required

### Delete SKU

```
DELETE /api/v1/skus/:id
```

**Authorization:** Required

---

## Product Localizations

### Create Localization

```
POST /api/v1/product-localizations
```

**Authorization:** Required

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| product_id | int64 | Yes | Product ID |
| language_code | string | Yes | Language code (en, zh-CN, ja, de, fr) |
| name | string | No | Translated name |
| description | string | No | Translated description |

### Get Localization

```
GET /api/v1/product-localizations/:id
```

### List Localizations by Product

```
GET /api/v1/products/:product_id/localizations
```

### Update Localization

```
PUT /api/v1/product-localizations/:id
```

### Delete Localization

```
DELETE /api/v1/product-localizations/:id
```

---

## Categories

### Create Category

```
POST /api/v1/categories
```

**Authorization:** Required

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | Category name |
| parent_id | int64 | No | Parent category ID (0 for root) |
| code | string | No | Category code |
| icon | string | No | Icon URL |
| image | string | No | Image URL |
| seo_title | string | No | SEO title |
| seo_description | string | No | SEO description |
| sort | int | No | Sort order |

### Get Category

```
GET /api/v1/categories/:id
```

### List Categories

```
GET /api/v1/categories
```

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| parent_id | int64 | Filter by parent (0 for root) |

### Get Category Tree

```
GET /api/v1/categories/tree
```

Returns hierarchical tree structure with children.

### Update Category

```
PUT /api/v1/categories/:id
```

### Update Category Status

```
PUT /api/v1/categories/:id/status
```

**Request Body:**

| Field | Type | Description |
|-------|------|-------------|
| status | int8 | 0=disabled, 1=enabled |

### Delete Category

```
DELETE /api/v1/categories/:id
```

### Move Category

```
PUT /api/v1/categories/:id/move
```

**Request Body:**

| Field | Type | Description |
|-------|------|-------------|
| new_parent_id | int64 | New parent category ID |

### Update Category Sort

```
PUT /api/v1/categories/sort
```

**Request Body:**

```json
{
  "sorts": [
    { "id": 1, "sort": 1 },
    { "id": 2, "sort": 2 }
  ]
}
```

### Set Category Market Visibility

```
PUT /api/v1/categories/:id/market-visibility
```

**Request Body:**

| Field | Type | Description |
|-------|------|-------------|
| market_ids | []int64 | Market IDs |
| visible | bool | Visibility status |

---

## Brands

### Create Brand

```
POST /api/v1/brands
```

**Authorization:** Required

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | Brand name |
| logo | string | No | Logo URL |
| description | string | No | Brand description |
| website | string | No | Official website |
| trademark_number | string | No | Trademark registration number |
| trademark_country | string | No | Trademark country |
| enable_page | bool | No | Enable brand page |
| sort | int | No | Sort order |

### Get Brand

```
GET /api/v1/brands/:id
```

### List Brands

```
GET /api/v1/brands
```

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | int | 1 | Page number |
| page_size | int | 20 | Items per page |
| name | string | - | Filter by name |
| status | int8 | - | Filter by status |

### Update Brand

```
PUT /api/v1/brands/:id
```

### Update Brand Status

```
PUT /api/v1/brands/:id/status
```

### Delete Brand

```
DELETE /api/v1/brands/:id
```

### Toggle Brand Page

```
PUT /api/v1/brands/:id/toggle-page
```

### Set Brand Market Visibility

```
PUT /api/v1/brands/:id/market-visibility
```

---

## Markets

### Create Market

```
POST /api/v1/markets
```

**Authorization:** Required

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| code | string | Yes | Market code (US, UK, DE, FR, AU) |
| name | string | Yes | Market name |
| currency | string | Yes | Currency (USD, GBP, EUR, AUD) |
| default_language | string | No | Default language code |
| flag | string | No | Flag emoji or icon |
| tax_rules | TaxConfig | No | Tax configuration |

**TaxConfig:**

| Field | Type | Description |
|-------|------|-------------|
| vat_rate | string | VAT rate (e.g., "0.20") |
| gst_rate | string | GST rate (e.g., "0.10") |
| ioss_enabled | bool | IOSS enabled for EU |
| include_tax | bool | Prices include tax |

### Get Market

```
GET /api/v1/markets/:id
```

### List Markets

```
GET /api/v1/markets
```

### Update Market

```
PUT /api/v1/markets/:id
```

### Delete Market

```
DELETE /api/v1/markets/:id
```

---

## Inventory

### Warehouses

#### Create Warehouse

```
POST /api/v1/warehouses
```

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| code | string | Yes | Warehouse code |
| name | string | Yes | Warehouse name |
| country | string | No | Country code |
| address | string | No | Full address |
| is_default | bool | No | Set as default |

#### Get Warehouse

```
GET /api/v1/warehouses/:id
```

#### List Warehouses

```
GET /api/v1/warehouses
```

#### Update Warehouse

```
PUT /api/v1/warehouses/:id
```

#### Update Warehouse Status

```
PUT /api/v1/warehouses/:id/status
```

#### Set Default Warehouse

```
PUT /api/v1/warehouses/:id/default
```

#### Delete Warehouse

```
DELETE /api/v1/warehouses/:id
```

### Stock Management

#### Update SKU Stock

```
PUT /api/v1/inventory/stock
```

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| sku_code | string | Yes | SKU code |
| warehouse_id | int64 | No | Warehouse ID (0 = all) |
| available_stock | int | Yes | New available stock |
| remark | string | No | Remark |

#### Adjust Stock

```
POST /api/v1/inventory/adjust
```

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| sku_code | string | Yes | SKU code |
| warehouse_id | int64 | Yes | Warehouse ID |
| quantity | int | Yes | Adjustment (positive = increase, negative = decrease) |
| remark | string | No | Reason for adjustment |

#### Get SKU Inventory

```
GET /api/v1/inventory/sku/:sku_code
```

**Response:**

```json
{
  "sku_code": "SKU-001-BLK-42",
  "product_id": 1,
  "total_stock": 30,
  "available_stock": 25,
  "locked_stock": 5,
  "safety_stock": 10,
  "is_low_stock": false,
  "warehouses": [
    {
      "warehouse_id": 1,
      "warehouse_name": "Shanghai Warehouse",
      "available_stock": 20,
      "locked_stock": 2
    }
  ]
}
```

#### Get Inventory Logs

```
GET /api/v1/inventory/logs
```

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | int | 1 | Page number |
| page_size | int | 20 | Items per page |
| sku_code | string | - | Filter by SKU |
| product_id | int64 | - | Filter by product |
| type | string | - | Filter by type (manual, order, return, adjustment) |

#### Get Low Stock SKUs

```
GET /api/v1/inventory/low-stock
```

Returns SKUs with stock below safety threshold.

#### Batch Update Safety Stock

```
PUT /api/v1/inventory/safety-stock
```

**Request Body:**

```json
{
  "items": [
    { "sku_code": "SKU-001", "safety_stock": 10 },
    { "sku_code": "SKU-002", "safety_stock": 15 }
  ]
}
```

---

## Error Codes

| Code | Message | Description |
|------|---------|-------------|
| 1001 | Unauthorized | Missing or invalid token |
| 1002 | Forbidden | Insufficient permissions |
| 2001 | Product not found | Product ID does not exist |
| 2002 | Invalid product status | Invalid status transition |
| 3001 | Category not found | Category ID does not exist |
| 3002 | Category has children | Cannot delete category with children |
| 4001 | Brand not found | Brand ID does not exist |
| 5001 | Market not found | Market ID does not exist |
| 6001 | Warehouse not found | Warehouse ID does not exist |
| 6002 | Insufficient stock | Not enough stock available |
| 7001 | SKU not found | SKU code does not exist |

---

## Rate Limiting

- **Default limit:** 100 requests per minute per tenant
- **Headers:**
  - `X-RateLimit-Limit`: Maximum requests
  - `X-RateLimit-Remaining`: Remaining requests
  - `X-RateLimit-Reset`: Reset timestamp

---

## Pagination

List endpoints support pagination with the following response structure:

```json
{
  "list": [...],
  "total": 100,
  "page": 1,
  "page_size": 20
}
```

**Query Parameters:**

| Parameter | Default | Max |
|-----------|---------|-----|
| page | 1 | - |
| page_size | 20 | 100 |