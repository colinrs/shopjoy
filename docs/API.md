# ShopJoy API 文档

## 基础信息

- **Base URL**: `http://localhost:8888/api/v1` (Admin) / `http://localhost:8889/api/v1` (Shop)
- **Content-Type**: `application/json`
- **Authentication**: JWT Bearer Token

## 认证

### 登录

```http
POST /users/login
Content-Type: application/json
X-Tenant-ID: 1

{
  "email": "admin@example.com",
  "password": "123456"
}
```

**Response:**

```json
{
  "code": 0,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 3600,
    "user": {
      "id": 1,
      "email": "admin@example.com",
      "name": "Admin"
    }
  }
}
```

### 注册

```http
POST /users/register
Content-Type: application/json
X-Tenant-ID: 1

{
  "email": "user@example.com",
  "password": "123456",
  "name": "User"
}
```

## 商品管理

### 创建商品

```http
POST /products
Authorization: Bearer {token}
Content-Type: application/json
X-Tenant-ID: 1

{
  "name": "商品名称",
  "description": "商品描述",
  "price": 9900,
  "currency": "CNY",
  "category_id": 1
}
```

### 获取商品列表

```http
GET /products?page=1&page_size=20
Authorization: Bearer {token}
X-Tenant-ID: 1
```

### 获取商品详情

```http
GET /products/{id}
Authorization: Bearer {token}
X-Tenant-ID: 1
```

### 更新商品

```http
PUT /products/{id}
Authorization: Bearer {token}
Content-Type: application/json
X-Tenant-ID: 1

{
  "name": "更新后的名称",
  "price": 19900
}
```

### 上架商品

```http
POST /products/{id}/on-sale
Authorization: Bearer {token}
X-Tenant-ID: 1
```

### 下架商品

```http
POST /products/{id}/off-sale
Authorization: Bearer {token}
X-Tenant-ID: 1
```

## 订单管理

### 创建订单

```http
POST /orders
Authorization: Bearer {token}
Content-Type: application/json
X-Tenant-ID: 1

{
  "items": [
    {
      "product_id": 1,
      "sku_id": 1,
      "quantity": 2
    }
  ],
  "address": {
    "name": "收货人",
    "phone": "13800138000",
    "province": "广东省",
    "city": "深圳市",
    "district": "南山区",
    "address": "科技园"
  }
}
```

### 获取订单列表

```http
GET /orders?page=1&page_size=20&status=0
Authorization: Bearer {token}
X-Tenant-ID: 1
```

### 获取订单详情

```http
GET /orders/{id}
Authorization: Bearer {token}
X-Tenant-ID: 1
```

### 取消订单

```http
PUT /orders/{id}/cancel
Authorization: Bearer {token}
Content-Type: application/json
X-Tenant-ID: 1

{
  "reason": "不想买了"
}
```

## 购物车

### 加入购物车

```http
POST /cart/items
Authorization: Bearer {token}
Content-Type: application/json
X-Tenant-ID: 1

{
  "product_id": 1,
  "sku_id": 1,
  "quantity": 1
}
```

### 获取购物车

```http
GET /cart
Authorization: Bearer {token}
X-Tenant-ID: 1
```

### 更新购物车商品

```http
PUT /cart/items/{id}
Authorization: Bearer {token}
Content-Type: application/json
X-Tenant-ID: 1

{
  "quantity": 3
}
```

### 删除购物车商品

```http
DELETE /cart/items/{id}
Authorization: Bearer {token}
X-Tenant-ID: 1
```

## 错误码

| Code | Message | Description |
|------|---------|-------------|
| 0 | success | 成功 |
| 400 | bad request | 请求参数错误 |
| 401 | unauthorized | 未授权 |
| 403 | forbidden | 禁止访问 |
| 404 | not found | 资源不存在 |
| 500 | internal error | 服务器内部错误 |
