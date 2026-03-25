# ShopJoy API Reference

> Comprehensive API documentation with code examples

**Last Updated:** 2026-03-26
**API Version:** v1
**Base URL:** `https://api.shopjoy.com`

## Table of Contents

1. [Authentication](#authentication)
2. [Common Patterns](#common-patterns)
3. [Products API](#products-api)
4. [Orders API](#orders-api)
5. [Promotions API](#promotions-api)
6. [Coupons API](#coupons-api)
7. [Error Handling](#error-handling)

---

## Authentication

All API requests require JWT Bearer token authentication.

### Login

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@shopjoy.com",
  "password": "your-password"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expire": "2026-03-27T10:00:00Z"
}
```

### Using the Token

```http
GET /api/v1/products
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## Common Patterns

### Pagination

All list endpoints support pagination:

```http
GET /api/v1/products?page=1&page_size=20
```

**Response:**
```json
{
  "list": [...],
  "total": 150,
  "page": 1,
  "page_size": 20
}
```

### Filtering

Most list endpoints support filtering:

```http
GET /api/v1/products?status=on_sale&category_id=1&min_price=1000&max_price=50000
```

### Sorting

```http
GET /api/v1/products?sort=created_at&order=desc
```

---

## Products API

### List Products

```http
GET /api/v1/products?page=1&page_size=20&status=on_sale
Authorization: Bearer {token}
```

**Code Examples:**

```bash
# cURL
curl -X GET "https://api.shopjoy.com/api/v1/products?page=1&page_size=20" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"
```

```javascript
// JavaScript (fetch)
const response = await fetch('https://api.shopjoy.com/api/v1/products?page=1&page_size=20', {
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});
const data = await response.json();
console.log(data);
```

```go
// Go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type ProductListResponse struct {
    List     []*Product `json:"list"`
    Total    int64      `json:"total"`
    Page     int        `json:"page"`
    PageSize int        `json:"page_size"`
}

func listProducts(token string) (*ProductListResponse, error) {
    url := "https://api.shopjoy.com/api/v1/products?page=1&page_size=20"

    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", "Bearer "+token)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result ProductListResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &result, nil
}
```

```python
# Python
import requests

def list_products(token):
    headers = {'Authorization': f'Bearer {token}'}
    response = requests.get(
        'https://api.shopjoy.com/api/v1/products',
        params={'page': 1, 'page_size': 20},
        headers=headers
    )
    return response.json()
```

### Create Product

```http
POST /api/v1/products
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "Nike Air Max 270",
  "description": "Comfortable running shoes",
  "price": 129900,
  "currency": "CNY",
  "category_id": 1,
  "brand": "Nike",
  "tags": ["sports", "running"],
  "images": ["https://cdn.example.com/p1.jpg"],
  "is_matrix_product": true,
  "hs_code": "64041100",
  "coo": "CN",
  "weight": "450",
  "weight_unit": "g"
}
```

**Response:**
```json
{
  "id": 1001
}
```

### Get Product Details

```http
GET /api/v1/products/1001
Authorization: Bearer {token}
```

**Response:**
```json
{
  "id": 1001,
  "sku": "SKU-001",
  "name": "Nike Air Max 270",
  "description": "Comfortable running shoes",
  "price": 129900,
  "currency": "CNY",
  "stock": 100,
  "status": "on_sale",
  "category_id": 1,
  "brand": "Nike",
  "tags": ["sports", "running"],
  "images": ["https://cdn.example.com/p1.jpg"],
  "is_matrix_product": true,
  "markets": [
    {
      "market_id": 1,
      "market_code": "CN",
      "market_name": "中国大陆",
      "is_enabled": true,
      "price": "1299.00",
      "currency": "CNY"
    }
  ],
  "created_at": "2026-03-26T10:00:00Z",
  "updated_at": "2026-03-26T10:00:00Z"
}
```

### Update Product Status

```http
# Put on sale
POST /api/v1/products/1001/on-sale
Authorization: Bearer {token}

# Take off sale
POST /api/v1/products/1001/off-sale
Authorization: Bearer {token}
```

---

## Orders API

### Create Order

```http
POST /api/v1/orders
Authorization: Bearer {token}
Content-Type: application/json

{
  "items": [
    {
      "sku_id": 1,
      "quantity": 1
    }
  ],
  "shipping_address": {
    "name": "张三",
    "phone": "13800000001",
    "province": "北京市",
    "city": "北京市",
    "district": "朝阳区",
    "detail": "建国路88号院1号楼101"
  },
  "remark": "尽快发货",
  "coupon_code": "SAVE30"
}
```

**Response:**
```json
{
  "id": "ORD202603260001",
  "order_no": "ORD202603260001",
  "status": "pending_payment",
  "pay_amount": 126900,
  "expire_at": "2026-03-26T11:00:00Z",
  "payment_url": "https://pay.alipay.com/..."
}
```

### List Orders

```http
GET /api/v1/orders?status=pending_payment&page=1&page_size=20
Authorization: Bearer {token}
```

### Get Order Details

```http
GET /api/v1/orders/ORD202603260001
Authorization: Bearer {token}
```

**Response:**
```json
{
  "id": "ORD202603260001",
  "order_no": "ORD202603260001",
  "status": "pending_payment",
  "total_amount": 129900,
  "discount_amount": 3000,
  "freight_amount": 0,
  "pay_amount": 126900,
  "currency": "CNY",
  "items": [
    {
      "id": 1,
      "product_id": 1,
      "sku_id": 1,
      "product_name": "Nike Air Max 270",
      "sku_name": "黑色 42码",
      "image": "https://cdn.example.com/p1-1.jpg",
      "price": 129900,
      "quantity": 1,
      "total_amount": 129900
    }
  ],
  "shipping_address": {
    "name": "张三",
    "phone": "13800000001",
    "province": "北京市",
    "city": "北京市",
    "district": "朝阳区",
    "detail": "建国路88号院1号楼101"
  },
  "expire_at": "2026-03-26T11:00:00Z",
  "created_at": "2026-03-26T10:00:00Z"
}
```

### Cancel Order

```http
POST /api/v1/orders/ORD202603260001/cancel
Authorization: Bearer {token}
```

---

## Promotions API

### Create Promotion

```http
POST /api/v1/promotions
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "春季大促",
  "description": "春季全场促销活动",
  "type": "discount",
  "start_time": "2026-03-26T00:00:00Z",
  "end_time": "2026-04-26T00:00:00Z",
  "product_ids": [1, 2, 3],
  "market_ids": [1, 2]
}
```

### Add Promotion Rules

```http
POST /api/v1/promotions/1/rules
Authorization: Bearer {token}
Content-Type: application/json

{
  "rules": [
    {
      "rule_type": "amount",
      "operator": "gte",
      "value": "10000",
      "discount_type": "percentage",
      "discount_value": "10",
      "priority": 1
    },
    {
      "rule_type": "amount",
      "operator": "gte",
      "value": "30000",
      "discount_type": "percentage",
      "discount_value": "15",
      "priority": 2
    }
  ]
}
```

---

## Coupons API

### Create Coupon

```http
POST /api/v1/coupons
Authorization: Bearer {token}
Content-Type: application/json

{
  "code": "SPRING2026",
  "name": "春季优惠券",
  "description": "春季活动专享优惠",
  "type": "percentage",
  "discount_value": "10",
  "min_order_amount": "10000",
  "max_discount": "5000",
  "start_time": "2026-03-26T00:00:00Z",
  "end_time": "2026-04-26T00:00:00Z",
  "usage_limit": 1000,
  "per_user_limit": 1
}
```

### Generate Coupon Codes

```http
POST /api/v1/coupons/generate
Authorization: Bearer {token}
Content-Type: application/json

{
  "prefix": "VIP",
  "quantity": 100,
  "length": 8,
  "coupon_config": "{\"type\":\"percentage\",\"discount_value\":\"15\"}"
}
```

**Response:**
```json
{
  "codes": [
    "VIPABC1234",
    "VIPDEF5678",
    "..."
  ],
  "count": 100
}
```

### Issue Coupon to User

```http
POST /api/v1/user-coupons
Authorization: Bearer {token}
Content-Type: application/json

{
  "user_id": 1,
  "coupon_id": 1
}
```

---

## Error Handling

### Error Response Format

All errors follow a consistent format:

```json
{
  "code": 30001,
  "msg": "Product not found"
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| 10001 | 401 | Unauthorized |
| 10002 | 403 | Permission denied |
| 20001 | 400 | Invalid request parameter |
| 30001 | 404 | Product not found |
| 30002 | 400 | Invalid product status |
| 40001 | 404 | Order not found |
| 40002 | 400 | Invalid order status |
| 40003 | 400 | Order already paid |
| 50001 | 404 | Payment not found |
| 50002 | 400 | Payment failed |
| 70001 | 404 | Coupon not found |
| 70002 | 400 | Coupon expired |
| 70003 | 400 | Coupon usage limit reached |

### Error Handling Example

```javascript
// JavaScript
async function createOrder(orderData) {
  const response = await fetch('/api/v1/orders', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(orderData)
  });

  const data = await response.json();

  if (!response.ok) {
    switch (data.code) {
      case 40001:
        throw new Error('Order not found');
      case 40002:
        throw new Error('Invalid order status');
      case 40003:
        throw new Error('Order already paid');
      default:
        throw new Error(data.msg);
    }
  }

  return data;
}
```

```go
// Go
func createOrder(client *http.Client, token string, order *OrderRequest) (*Order, error) {
    body, _ := json.Marshal(order)
    req, _ := http.NewRequest("POST", "/api/v1/orders", bytes.NewReader(body))
    req.Header.Set("Authorization", "Bearer "+token)

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result struct {
        Code int    `json:"code"`
        Msg  string `json:"msg"`
        Data *Order `json:"data"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    if result.Code != 0 {
        return nil, fmt.Errorf("API error %d: %s", result.Code, result.Msg)
    }

    return result.Data, nil
}
```

---

## Rate Limiting

API requests are rate limited per tenant:

| Plan | Requests/min | Requests/day |
|------|--------------|--------------|
| Free | 100 | 10,000 |
| Pro | 1,000 | 100,000 |
| Enterprise | 10,000 | Unlimited |

Rate limit headers are included in responses:

```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1648281600
```

---

## Webhooks

Configure webhooks to receive event notifications:

### Supported Events

| Event | Description |
|-------|-------------|
| `order.created` | New order created |
| `order.paid` | Order paid successfully |
| `order.shipped` | Order shipped |
| `order.completed` | Order completed |
| `order.cancelled` | Order cancelled |
| `payment.success` | Payment successful |
| `payment.refunded` | Payment refunded |

### Webhook Payload

```json
{
  "event": "order.paid",
  "timestamp": "2026-03-26T10:05:00Z",
  "data": {
    "order_id": "ORD202603260001",
    "user_id": 1,
    "amount": 126900,
    "currency": "CNY"
  },
  "signature": "sha256=..."
}
```

---

## Related Documentation

- [OpenAPI Specification](./openapi.yaml)
- [Product Schema](../domains/product/2026-03-26-product-schema.md)
- [Order Schema](../domains/order/2026-03-26-order-schema.md)
- [Error Codes](../reference/2026-03-21-error-codes.md)