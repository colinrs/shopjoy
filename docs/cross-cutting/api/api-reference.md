# ShopJoy API Reference

> **Version:** 1.0
> **Last Updated:** 2026-03-27
> **API Base URL:** `http://localhost:8888/api/v1` (Admin API)

---

## Overview

The ShopJoy API is a RESTful API built with go-zero framework. It provides comprehensive endpoints for managing a multi-tenant e-commerce platform.

### Base URL

| Environment | Admin API | Shop API |
|-------------|-----------|----------|
| Development | `http://localhost:8888/api/v1` | `http://localhost:8889/api/v1` |
| Production | `https://api-admin.shopjoy.com/api/v1` | `https://api.shopjoy.com/api/v1` |

### Authentication

Most endpoints require JWT authentication via the `Authorization` header:

```
Authorization: Bearer <access_token>
```

### Common Headers

| Header | Required | Description |
|--------|----------|-------------|
| `Content-Type` | Yes | `application/json` |
| `Authorization` | Yes | Bearer token (except auth endpoints) |
| `X-Tenant-ID` | No | Tenant identifier for multi-tenancy |

### Error Response Format

All errors follow a consistent format:

```json
{
  "code": 30012,
  "message": "product not found"
}
```

See [Error Codes Reference](../reference/2026-03-22-error-codes.md) for complete list.

---

## API Groups

### 1. Authentication (`/auth`)

Authentication endpoints do not require authorization.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/login` | Admin login |
| POST | `/auth/register` | Register tenant admin |

#### POST /auth/login

**Request:**
```json
{
  "account": "admin@shopjoy.com",
  "password": "secure_password",
  "ip": "192.168.1.1"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 86400,
  "user": {
    "id": 1,
    "tenant_id": 1,
    "username": "admin",
    "email": "admin@shopjoy.com",
    "real_name": "Administrator",
    "type": 1,
    "type_text": "平台超管"
  }
}
```

#### POST /auth/register

**Request:**
```json
{
  "email": "newadmin@shopjoy.com",
  "real_name": "New Admin",
  "password": "secure_password",
  "tenant_id": 0
}
```

---

### 2. Admin Users (`/admin-users`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/admin-users` | List admin users |
| POST | `/admin-users` | Create admin user |
| GET | `/admin-users/:id` | Get admin user detail |
| PUT | `/admin-users/:id` | Update admin user |
| DELETE | `/admin-users/:id` | Delete admin user |
| POST | `/admin-users/:id/reset-password` | Reset admin password |
| PUT | `/admin-users/:id/roles` | Assign roles |
| POST | `/admin-users/:id/disable` | Disable admin |
| POST | `/admin-users/:id/enable` | Enable admin |
| PUT | `/admin-users/profile` | Update own profile |
| PUT | `/admin-users/password` | Change own password |

#### Admin User Types

| Type | Value | Description |
|------|-------|-------------|
| 平台超管 | 1 | Super admin (platform level) |
| 商家管理员 | 2 | Tenant admin |
| 商家子账号 | 3 | Sub-account |

#### GET /admin-users

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | int | 1 | Page number |
| page_size | int | 20 | Items per page |
| keyword | string | - | Search keyword |
| type | int | - | Filter by user type |
| status | int | - | Filter by status (1=normal, 2=disabled) |
| tenant_id | int64 | - | Filter by tenant (super admin only) |

**Response:**
```json
{
  "list": [
    {
      "id": 1,
      "tenant_id": 1,
      "username": "admin",
      "email": "admin@shopjoy.com",
      "mobile": "+86-13800138000",
      "real_name": "Administrator",
      "avatar": "https://...",
      "type": 1,
      "type_text": "平台超管",
      "status": 1,
      "created_at": "2026-03-01T10:00:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "page_size": 20
}
```

---

### 3. Users (`/users`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/users` | List users |
| GET | `/users/:id` | Get user detail |
| PUT | `/users/:id` | Update user |
| POST | `/users/:id/suspend` | Suspend user |
| POST | `/users/:id/activate` | Activate user |
| POST | `/users/:id/suspend-reason` | Suspend with reason |
| DELETE | `/users/:id` | Delete user |
| POST | `/users/:id/reset-password` | Reset password |
| GET | `/users/:id/detail` | Get enhanced user detail |
| GET | `/users/:id/addresses` | Get user addresses |
| GET | `/users/stats` | Get user statistics |
| GET | `/users/stats/enhanced` | Get enhanced statistics |
| GET | `/users/export` | Export users |
| GET | `/users/enhanced` | List enhanced users |

#### User Status

| Status | Value | Description |
|--------|-------|-------------|
| Inactive | 0 | Not active |
| Active | 1 | Active user |
| Suspended | 2 | Suspended user |
| Deleted | 3 | Soft deleted |

#### Extended User Response

```json
{
  "id": 1,
  "tenant_id": 1,
  "email": "user@example.com",
  "phone": "+86-13800138000",
  "name": "John Doe",
  "avatar": "https://...",
  "gender": 1,
  "gender_text": "Male",
  "birthday": "1990-01-01",
  "status": 1,
  "status_text": "Active",
  "points_balance": 5000,
  "points_frozen": 100,
  "total_earned": 10000,
  "total_redeemed": 5000,
  "order_count": 25,
  "total_spent": "2500.00",
  "review_count": 10,
  "last_login": "2026-03-26T10:00:00Z",
  "created_at": "2026-01-01T00:00:00Z"
}
```

---

### 4. Roles (`/roles`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/roles` | List roles |
| POST | `/roles` | Create role |
| GET | `/roles/:id` | Get role with permissions |
| PUT | `/roles/:id` | Update role |
| DELETE | `/roles/:id` | Delete role |
| PUT | `/roles/:id/status` | Update role status |
| PUT | `/roles/:id/permissions` | Update role permissions |
| GET | `/permissions` | List all permissions |

#### Permission Types

| Type | Value | Description |
|------|-------|-------------|
| Menu | 0 | Menu permission |
| Button | 1 | Button/Action permission |
| API | 2 | API endpoint permission |

---

### 5. Products (`/products`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/products` | Create product |
| PUT | `/products/:id` | Update product |
| GET | `/products/:id` | Get product detail |
| GET | `/products` | List products |
| POST | `/products/:id/on-sale` | Put on sale |
| POST | `/products/:id/off-sale` | Take off sale |
| PUT | `/products/:id/stock` | Update stock |

#### Product Status

| Status | Value | Description |
|--------|-------|-------------|
| draft | 0 | Draft (not visible) |
| on_sale | 1 | On sale |
| off_sale | 2 | Off sale |
| deleted | 3 | Soft deleted |

#### Product Request Example

```json
{
  "name": "Premium T-Shirt",
  "description": "High quality cotton t-shirt",
  "price": 2999,
  "currency": "USD",
  "category_id": 1,
  "sku": "TSHIRT-001",
  "brand": "ShopJoy",
  "tags": ["summer", "casual"],
  "images": ["https://example.com/tshirt1.jpg"],
  "is_matrix_product": false,
  "hs_code": "6109.10",
  "coo": "CN",
  "weight": "200",
  "weight_unit": "g",
  "length": "30",
  "width": "25",
  "height": "2"
}
```

#### Product Response Example

```json
{
  "id": 1,
  "name": "Premium T-Shirt",
  "description": "High quality cotton t-shirt",
  "price": 2999,
  "currency": "USD",
  "cost_price": 1500,
  "stock": 100,
  "status": "on_sale",
  "category_id": 1,
  "created_at": "2026-03-01T10:00:00Z",
  "updated_at": "2026-03-15T10:00:00Z",
  "sku": "TSHIRT-001",
  "brand": "ShopJoy",
  "tags": ["summer", "casual"],
  "images": ["https://example.com/tshirt1.jpg"],
  "is_matrix_product": false,
  "hs_code": "6109.10",
  "coo": "CN",
  "weight": "200",
  "weight_unit": "g",
  "length": "30",
  "width": "25",
  "height": "2",
  "dangerous_goods": [],
  "markets": [
    {
      "market_id": 1,
      "market_code": "US",
      "market_name": "United States",
      "is_enabled": true,
      "price": "35.99",
      "currency": "USD"
    }
  ]
}
```

---

### 6. SKU Management (`/skus`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/skus` | Create SKU |
| PUT | `/skus/:id` | Update SKU |
| GET | `/skus/:id` | Get SKU detail |
| DELETE | `/skus/:id` | Delete SKU |
| GET | `/products/:product_id/skus` | List SKUs by product |

#### SKU Response

```json
{
  "id": 1,
  "product_id": 1,
  "code": "TSHIRT-001-RED-L",
  "price": 2999,
  "currency": "USD",
  "stock": 50,
  "available_stock": 45,
  "locked_stock": 5,
  "safety_stock": 10,
  "pre_sale_enabled": false,
  "attributes": {
    "color": "Red",
    "size": "L"
  },
  "status": "on_sale",
  "is_low_stock": false,
  "created_at": "2026-03-01T10:00:00Z",
  "updated_at": "2026-03-15T10:00:00Z"
}
```

---

### 7. Product Localization (`/product-localizations`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/product-localizations` | Create localization |
| PUT | `/product-localizations/:id` | Update localization |
| GET | `/product-localizations/:id` | Get localization |
| GET | `/products/:product_id/localizations` | List localizations |
| DELETE | `/product-localizations/:id` | Delete localization |

---

### 8. Categories (`/categories`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/categories` | Create category |
| PUT | `/categories/:id` | Update category |
| GET | `/categories/:id` | Get category detail |
| GET | `/categories` | List categories |
| GET | `/categories/tree` | Get category tree |
| PUT | `/categories/:id/status` | Update status |
| DELETE | `/categories/:id` | Delete category |
| PUT | `/categories/sort` | Update sort order |
| PUT | `/categories/:id/move` | Move category |
| GET | `/categories/:id/product-count` | Get product count |
| PUT | `/categories/:id/market-visibility` | Set market visibility |
| GET | `/categories/:id/market-visibility` | Get market visibility |

#### Category Tree Response

```json
[
  {
    "id": 1,
    "parent_id": 0,
    "name": "Clothing",
    "code": "CLOTHING",
    "level": 1,
    "sort": 1,
    "icon": "https://...",
    "image": "https://...",
    "status": 1,
    "product_count": 150,
    "children": [
      {
        "id": 2,
        "parent_id": 1,
        "name": "T-Shirts",
        "code": "TSHIRTS",
        "level": 2,
        "sort": 1,
        "status": 1,
        "product_count": 50,
        "children": []
      }
    ]
  }
]
```

---

### 9. Brands (`/brands`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/brands` | Create brand |
| PUT | `/brands/:id` | Update brand |
| GET | `/brands/:id` | Get brand detail |
| GET | `/brands` | List brands |
| PUT | `/brands/:id/status` | Update status |
| DELETE | `/brands/:id` | Delete brand |
| PUT | `/brands/:id/toggle-page` | Toggle brand page |
| GET | `/brands/:id/product-count` | Get product count |
| PUT | `/brands/:id/market-visibility` | Set market visibility |
| GET | `/brands/:id/market-visibility` | Get market visibility |

---

### 10. Markets (`/markets`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/markets` | Create market |
| PUT | `/markets/:id` | Update market |
| GET | `/markets/:id` | Get market detail |
| GET | `/markets` | List markets |
| DELETE | `/markets/:id` | Delete market |

#### Market Tax Configuration

```json
{
  "vat_rate": "20.0",
  "gst_rate": "10.0",
  "ioss_enabled": true,
  "include_tax": false
}
```

---

### 11. Product Markets (`/products/:id/markets`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/products/:id/markets` | List product markets |
| PUT | `/products/:id/markets/:market_id` | Update product market |
| POST | `/products/:id/push-to-market` | Push to markets |
| DELETE | `/products/:id/markets/:market_id` | Remove from market |

---

### 12. Inventory (`/inventory`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| PUT | `/inventory/stock` | Update SKU stock |
| POST | `/inventory/adjust` | Adjust stock |
| GET | `/inventory/sku/:sku_code` | Get SKU inventory |
| GET | `/inventory/logs` | Get inventory logs |
| GET | `/inventory/low-stock` | Get low stock SKUs |
| PUT | `/inventory/safety-stock` | Batch update safety stock |

#### Warehouse Management (`/warehouses`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/warehouses` | Create warehouse |
| PUT | `/warehouses/:id` | Update warehouse |
| GET | `/warehouses/:id` | Get warehouse detail |
| GET | `/warehouses` | List warehouses |
| PUT | `/warehouses/:id/status` | Update status |
| PUT | `/warehouses/:id/default` | Set default |
| DELETE | `/warehouses/:id` | Delete warehouse |

---

### 13. Promotions (`/promotions`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/promotions` | Create promotion |
| PUT | `/promotions/:id` | Update promotion |
| GET | `/promotions/:id` | Get promotion detail |
| GET | `/promotions` | List promotions |
| DELETE | `/promotions/:id` | Delete promotion |
| POST | `/promotions/:id/activate` | Activate promotion |
| POST | `/promotions/:id/deactivate` | Deactivate promotion |
| GET | `/promotions/:id/rules` | Get promotion rules |
| POST | `/promotions/:id/rules` | Add promotion rules |
| PUT | `/promotion-rules/:id` | Update rule |
| DELETE | `/promotion-rules/:id` | Delete rule |

#### Promotion Types

| Type | Description |
|------|-------------|
| discount | General discount |
| coupon | Coupon-based promotion |
| flash_sale | Flash sale event |
| bundle | Bundle offer |

#### Discount Types

| Type | Description |
|------|-------------|
| percentage | Percentage off |
| fixed_amount | Fixed amount off |
| buy_x_get_y | Buy X get Y free |

---

### 14. Coupons (`/coupons`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/coupons` | Create coupon |
| PUT | `/coupons/:id` | Update coupon |
| GET | `/coupons/:id` | Get coupon detail |
| GET | `/coupons` | List coupons |
| DELETE | `/coupons/:id` | Delete coupon |
| POST | `/coupons/generate` | Generate codes |
| GET | `/coupons/:id/usage` | Get usage records |

#### Coupon Types

| Type | Description |
|------|-------------|
| fixed_amount | Fixed amount discount |
| percentage | Percentage discount |
| free_shipping | Free shipping |

#### Coupon Status

| Status | Description |
|--------|-------------|
| draft | Not active |
| active | Active and usable |
| inactive | Manually deactivated |
| expired | Past end date |

---

### 15. User Coupons (`/user-coupons`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/user-coupons` | Issue coupon to user |
| GET | `/user-coupons` | List user coupons |

#### User Coupon Status

| Status | Description |
|--------|-------------|
| available | Can be used |
| used | Already used |
| expired | Past validity |

---

### 16. Payments (`/payments`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/payments/stats` | Get payment statistics |
| GET | `/payments/transactions` | List transactions |
| GET | `/payments/transactions/:id` | Get transaction detail |
| GET | `/orders/:id/payment` | Get order payment info |
| POST | `/orders/:id/refund` | Initiate refund |
| POST | `/webhooks/stripe` | Stripe webhook |

#### Payment Transaction Status

| Status | Value | Description |
|--------|-------|-------------|
| Pending | 0 | Awaiting payment |
| Succeeded | 1 | Payment successful |
| Failed | 2 | Payment failed |

#### Refund Status

| Status | Value | Description |
|--------|-------|-------------|
| Pending | 0 | Awaiting approval |
| Approved | 1 | Approved, processing |
| Rejected | 2 | Rejected |
| Completed | 3 | Refund completed |
| Cancelled | 4 | Cancelled |

---

### 17. Fulfillment / Shipments (`/shipments`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/shipments` | Create shipment |
| POST | `/shipments/batch` | Batch create shipments |
| GET | `/shipments` | List shipments |
| GET | `/shipments/:id` | Get shipment detail |
| PUT | `/shipments/:id` | Update shipment |
| PUT | `/shipments/:id/status` | Update status |
| GET | `/orders/:order_id/shipments` | Get order shipments |
| GET | `/carriers` | List carriers |

#### Shipment Status

| Status | Value | Description |
|--------|-------|-------------|
| Pending | 0 | Created, not shipped |
| Shipped | 1 | Handed to carrier |
| In Transit | 2 | In transit |
| Delivered | 3 | Delivered |
| Failed | 4 | Delivery failed |

---

### 18. Refunds (`/refunds`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/refunds` | List refunds |
| GET | `/refunds/:id` | Get refund detail |
| PUT | `/refunds/:id/approve` | Approve refund |
| PUT | `/refunds/:id/reject` | Reject refund |
| GET | `/refund-reasons` | List reasons |
| GET | `/refunds/statistics` | Get statistics |

#### Refund Types

| Type | Value | Description |
|------|-------|-------------|
| Full Refund | 1 | Complete order refund |
| Partial Refund | 2 | Partial order refund |

---

### 19. Orders (`/orders`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/orders` | List orders (fulfillment view) |
| GET | `/orders/:id` | Get order fulfillment detail |
| PUT | `/orders/:id/ship` | Ship order |
| GET | `/orders/fulfillment-summary` | Get fulfillment summary |
| PUT | `/orders/:id/remark` | Update remark |
| PUT | `/orders/:id/adjust-price` | Adjust price |
| GET | `/orders/export` | Export orders |
| PUT | `/orders/:id/cancel` | Cancel order |
| POST | `/orders/:id/remind-payment` | Send payment reminder |

#### Order Fulfillment Status

| Status | Value | Description |
|--------|-------|-------------|
| Pending | 0 | Awaiting shipment |
| Partial | 1 | Partially shipped |
| Shipped | 2 | Fully shipped |
| Delivered | 3 | Delivered |

---

### 20. Reviews (`/reviews`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/reviews` | List reviews |
| GET | `/reviews/:id` | Get review detail |
| PUT | `/reviews/:id/approve` | Approve review |
| PUT | `/reviews/:id/hide` | Hide review |
| PUT | `/reviews/:id/show` | Show review |
| DELETE | `/reviews/:id` | Delete review |
| PUT | `/reviews/:id/featured` | Toggle featured |
| POST | `/reviews/:id/reply` | Create reply |
| PUT | `/reviews/:id/reply` | Update reply |
| DELETE | `/reviews/:id/reply` | Delete reply |
| POST | `/reviews/batch-approve` | Batch approve |
| POST | `/reviews/batch-hide` | Batch hide |
| GET | `/reviews/stats` | Get statistics |
| GET | `/reviews/product/:product_id/stats` | Get product stats |

#### Review Status

| Status | Description |
|--------|-------------|
| pending | Awaiting moderation |
| approved | Approved and visible |
| hidden | Hidden by admin |
| deleted | Soft deleted |

---

### 21. Points (`/points`)

#### Statistics (`/points/stats`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/points/stats` | Get points overview |
| GET | `/points/stats/trend` | Get trend data |
| GET | `/points/stats/top-users` | Get top users |
| GET | `/points/stats/expiring` | Get expiring points |

#### Earn Rules (`/points/earn-rules`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/points/earn-rules` | List earn rules |
| GET | `/points/earn-rules/:id` | Get earn rule detail |
| POST | `/points/earn-rules` | Create earn rule |
| PUT | `/points/earn-rules/:id` | Update earn rule |
| DELETE | `/points/earn-rules/:id` | Delete earn rule |
| POST | `/points/earn-rules/:id/activate` | Activate rule |
| POST | `/points/earn-rules/:id/deactivate` | Deactivate rule |

#### Redeem Rules (`/points/redeem-rules`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/points/redeem-rules` | List redeem rules |
| GET | `/points/redeem-rules/:id` | Get redeem rule detail |
| POST | `/points/redeem-rules` | Create redeem rule |
| PUT | `/points/redeem-rules/:id` | Update redeem rule |
| DELETE | `/points/redeem-rules/:id` | Delete redeem rule |
| POST | `/points/redeem-rules/:id/activate` | Activate rule |
| POST | `/points/redeem-rules/:id/deactivate` | Deactivate rule |

#### Accounts (`/points/accounts`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/points/accounts` | List accounts |
| GET | `/points/accounts/:id` | Get account detail |
| GET | `/points/accounts/by-user/:userId` | Get by user ID |
| GET | `/points/accounts/:id/transactions` | Get transactions |
| POST | `/points/accounts/:id/adjust` | Manual adjustment |

#### Transactions (`/points/transactions`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/points/transactions` | List all transactions |
| GET | `/points/transactions/:id` | Get transaction detail |

#### Redemptions (`/points/redemptions`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/points/redemptions` | List redemptions |
| GET | `/points/redemptions/:id` | Get redemption detail |

#### Earn Rule Scenarios

| Scenario | Description |
|----------|-------------|
| ORDER_PAYMENT | Points earned on order payment |
| SIGN_IN | Points earned on daily login |
| PRODUCT_REVIEW | Points earned on reviewing product |
| FIRST_ORDER | Points for first order |

#### Points Transaction Types

| Type | Description |
|------|-------------|
| EARN | Points earned |
| REDEEM | Points redeemed |
| ADJUST | Manual adjustment |
| EXPIRE | Points expired |
| FREEZE | Points frozen |
| UNFREEZE | Points unfrozen |

---

### 22. Storefront / Themes (`/themes`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/themes` | List themes |
| GET | `/themes/current` | Get current theme |
| PUT | `/themes/switch` | Switch theme |
| PUT | `/themes/config` | Update theme config |
| GET | `/themes/audit-logs` | Get audit logs |

#### Pages (`/pages`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/pages` | List pages |
| GET | `/pages/:id` | Get page detail |
| PUT | `/pages/:id/draft` | Save draft |
| PUT | `/pages/:id/publish` | Publish page |
| PUT | `/pages/:id/unpublish` | Unpublish page |
| POST | `/pages/:id/decorations` | Add decoration |
| PUT | `/decorations/:id` | Update decoration |
| DELETE | `/decorations/:id` | Delete decoration |
| PUT | `/pages/:id/blocks/reorder` | Reorder blocks |

#### Versions (`/pages/:id/versions`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/pages/:id/versions` | List versions |
| GET | `/pages/:id/versions/:version` | Get version detail |
| PUT | `/pages/:id/restore` | Restore version |

#### SEO (`/seo`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/seo/global` | Get global SEO |
| PUT | `/seo/global` | Update global SEO |
| GET | `/seo/pages` | List page SEO configs |
| GET | `/seo/pages/:page_type` | Get page SEO |
| PUT | `/seo/pages/:page_type` | Update page SEO |

---

### 23. Shipping (`/shipping-templates`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/shipping-templates` | List templates |
| POST | `/shipping-templates` | Create template |
| GET | `/shipping-templates/:id` | Get template detail |
| PUT | `/shipping-templates/:id` | Update template |
| DELETE | `/shipping-templates/:id` | Delete template |
| PUT | `/shipping-templates/:id/set-default` | Set default |

#### Shipping Zones (`/shipping-templates/:template_id/zones`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/shipping-templates/:template_id/zones` | Create zone |
| PUT | `/shipping-zones/:id` | Update zone |
| DELETE | `/shipping-zones/:id` | Delete zone |
| PUT | `/shipping-templates/:template_id/zones/reorder` | Reorder zones |

#### Fee Types

| Type | Description |
|------|-------------|
| fixed | Fixed shipping fee |
| by_count | Fee by item count |
| by_weight | Fee by weight |
| free | Free shipping |

#### Template Mappings (`/shipping-template-mappings`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/shipping-templates/:id/mappings` | List mappings |
| POST | `/shipping-template-mappings` | Create mapping |
| DELETE | `/shipping-template-mappings/:id` | Delete mapping |

#### Shipping Calculator

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/shipping/calculate` | Calculate fee |

#### Regions

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/regions` | List regions |

---

### 24. Shop Settings (`/shop`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/shop` | Get shop settings |
| PUT | `/shop` | Update shop settings |
| GET | `/shop/business-hours` | Get business hours |
| PUT | `/shop/business-hours` | Update business hours |
| GET | `/shop/notifications` | Get notification settings |
| PUT | `/shop/notifications` | Update notifications |
| GET | `/shop/payment` | Get payment settings |
| PUT | `/shop/payment` | Update payment settings |
| GET | `/shop/shipping` | Get shipping settings |
| PUT | `/shop/shipping` | Update shipping settings |

---

### 25. Dashboard (`/dashboard`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/dashboard` | Get all dashboard data |
| GET | `/dashboard/overview` | Get overview stats |
| GET | `/dashboard/sales-trend` | Get sales trend |
| GET | `/dashboard/order-status` | Get order status distribution |
| GET | `/dashboard/top-products` | Get top products |
| GET | `/dashboard/pending-orders` | Get pending orders |
| GET | `/dashboard/activities` | Get recent activities |

---

### 26. Upload (`/uploads`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/uploads` | Upload file |
| DELETE | `/uploads/:id` | Delete file |

#### Upload Categories

| Category | Description |
|----------|-------------|
| product | Product images |
| banner | Banner images |
| avatar | User avatars |

---

## Pagination

All list endpoints support pagination:

| Parameter | Type | Default | Max | Description |
|-----------|------|---------|-----|-------------|
| page | int | 1 | - | Page number |
| page_size | int | 20 | 100 | Items per page |

**Response includes:**
```json
{
  "list": [...],
  "total": 100,
  "page": 1,
  "page_size": 20
}
```

---

## Filtering

List endpoints support various filter parameters:

### Products
- `name` - Product name (partial match)
- `category_id` - Category filter
- `status` - Status filter
- `min_price` - Minimum price
- `max_price` - Maximum price
- `market_id` - Market filter

### Orders
- `order_no` - Order number
- `user_id` - User filter
- `status` - Order status
- `fulfillment_status` - Fulfillment status
- `refund_status` - Refund status
- `start_time` - Start date (RFC3339)
- `end_time` - End date (RFC3339)

### Users
- `keyword` - Search keyword
- `status` - User status
- `register_start` - Registration start date
- `register_end` - Registration end date

---

## Time Format

All timestamps in API requests and responses use RFC3339 format:

```
2026-03-27T10:00:00Z
```

For date-only filters:
```
2026-03-27
```

---

## Monetary Values

Monetary values in API are integers representing the smallest currency unit (cents):

| Value | Represents |
|-------|------------|
| 2999 | $29.99 |
| 1000 | $10.00 |

Some APIs use string format for decimal precision (see specific endpoint docs).

---

## Rate Limiting

API requests are rate limited to protect the service:

| Endpoint Group | Limit |
|----------------|-------|
| Auth endpoints | 10/minute |
| Read endpoints | 100/minute |
| Write endpoints | 30/minute |

When rate limited, the API returns:
```json
{
  "code": 42900,
  "message": "请求过于频繁，请稍后再试"
}
```

---

## Versioning

Current API version: **v1**

All endpoints are prefixed with `/api/v1/`.

---

## Document History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-03-27 | Technical Team | Initial comprehensive API reference |
