# Error Code Reference

This document lists all error codes used in the ShopJoy API.

## Error Response Format

```json
{
  "code": 30001,
  "msg": "商品名称不能为空"
}
```

## HTTP Status Codes

| Status | Meaning |
|--------|---------|
| 200 | Success |
| 400 | Bad Request - Invalid parameters |
| 401 | Unauthorized - Authentication required |
| 403 | Forbidden - No permission |
| 404 | Not Found - Resource doesn't exist |
| 409 | Conflict - Duplicate resource |
| 500 | Internal Server Error |

---

## General Errors (10xxx - 20xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 0 | 200 | success | Operation successful |
| 10001 | 400 | 参数有误 | Invalid parameters |
| 10003 | 400 | 参数有误 | Parameter error |
| 10004 | 200 | unknown error | Unknown error |
| 20001 | 200 | Validation failed | Request validation failed |
| 20002 | 200 | Database error | Database operation failed |
| 20003 | 200 | 业务已存在 | Business entity already exists |
| 20004 | 200 | 业务不存在 | Business entity not found |
| 20006 | 200 | http client error | HTTP client error |

---

## Authentication Errors (40xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 40100 | 401 | 未授权，请先登录 | Unauthorized - Please login first |
| 40101 | 401 | Token 已过期，请重新登录 | Token expired - Please login again |
| 40102 | 401 | 无效的 Token | Invalid token |
| 40300 | 403 | 没有权限访问该资源 | No permission to access resource |
| 40400 | 404 | 资源不存在 | Resource not found |
| 40500 | 405 | 请求方法不允许 | Method not allowed |
| 42900 | 429 | 请求过于频繁，请稍后再试 | Too many requests |
| 50000 | 500 | 服务内部错误 | Internal server error |
| 50300 | 503 | 服务暂时不可用 | Service unavailable |

---

## Admin User Module (10xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 10001 | 400 | invalid email format | Email format is invalid |
| 10002 | 400 | invalid phone format | Phone format is invalid |
| 10003 | 400 | password too weak | Password doesn't meet requirements |
| 10004 | 404 | admin user not found | Admin user not found |
| 10005 | 409 | duplicate admin user | Admin user already exists |
| 10006 | 401 | wrong password | Incorrect password |
| 10007 | 400 | cannot delete yourself | Cannot delete your own account |
| 10008 | 400 | user already deleted | User has already been deleted |
| 10009 | 403 | account disabled or deleted | Account is disabled |
| 10010 | 400 | passwords do not match | Password confirmation mismatch |
| 10011 | 403 | permission denied | Insufficient permissions |

---

## User Module (11xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 11001 | 400 | invalid email format | Email format is invalid |
| 11002 | 400 | invalid phone format | Phone format is invalid |
| 11003 | 400 | password too weak | Password doesn't meet requirements |
| 11004 | 404 | user not found | User not found |
| 11005 | 409 | duplicate user | User already exists |
| 11006 | 401 | wrong password | Incorrect password |
| 11007 | 400 | user already deleted | User has already been deleted |
| 11008 | 400 | passwords do not match | Password confirmation mismatch |

---

## Product Module (30xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 30001 | 400 | 商品名称不能为空 | Product name cannot be empty |
| 30002 | 400 | 商品价格必须大于0 | Product price must be greater than 0 |
| 30003 | 400 | 币种不匹配 | Currency mismatch |
| 30004 | 400 | 金额不足 | Insufficient amount |
| 30005 | 400 | 商品已删除 | Product has been deleted |
| 30006 | 400 | 无效的状态转换 | Invalid status transition |
| 30007 | 400 | 库存不能为0 | Stock cannot be 0 for on-sale product |
| 30008 | 400 | 库存不能为负数 | Stock cannot be negative |
| 30009 | 400 | 商品未上架 | Product is not on sale |
| 30010 | 400 | 无效的数量 | Invalid quantity |
| 30011 | 400 | 库存不足 | Insufficient stock |
| 30012 | 404 | 商品不存在 | Product not found |
| 30013 | 400 | invalid product id | Product ID is invalid |

### Category Errors (301xx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 30101 | 404 | category not found | Category not found |
| 30102 | 409 | duplicate category | Category already exists |
| 30103 | 400 | invalid category | Invalid category data |
| 30104 | 400 | category has children | Cannot delete category with children |
| 30105 | 400 | category has products | Cannot delete category with products |

---

## Order Module (40xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 40001 | 404 | order not found | Order not found |
| 40002 | 400 | invalid order status | Invalid order status |
| 40003 | 400 | order already paid | Order has already been paid |
| 40004 | 400 | order not paid | Order has not been paid |
| 40005 | 400 | order expired | Order has expired |
| 40006 | 400 | insufficient stock | Insufficient stock for order |
| 40007 | 400 | invalid amount | Invalid order amount |
| 40008 | 400 | cart is empty | Cannot create order with empty cart |

---

## Payment Module (50xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 50001 | 404 | payment not found | Payment record not found |
| 50002 | 400 | invalid payment amount | Invalid payment amount |
| 50003 | 402 | payment failed | Payment processing failed |
| 50004 | 400 | payment already completed | Payment already completed |
| 50005 | 400 | payment expired | Payment session expired |

---

## Cart Module (60xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 60001 | 404 | cart item not found | Cart item not found |
| 60002 | 400 | invalid quantity | Invalid item quantity |
| 60003 | 400 | cart is empty | Cart is empty |

---

## Coupon Module (70xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 70001 | 404 | coupon not found | Coupon not found |
| 70002 | 400 | coupon expired | Coupon has expired |
| 70003 | 400 | coupon used up | Coupon usage limit reached |
| 70004 | 400 | coupon not started | Coupon campaign hasn't started |
| 70005 | 400 | coupon already used | User has already used this coupon |
| 70006 | 400 | invalid coupon code | Invalid coupon code |
| 70007 | 400 | cart amount below minimum | Cart amount below minimum requirement |

---

## Promotion Module (80xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 80001 | 404 | promotion not found | Promotion not found |
| 80002 | 400 | invalid promotion | Invalid promotion data |
| 80003 | 400 | promotion expired | Promotion has expired |
| 80004 | 400 | promotion not started | Promotion hasn't started yet |

---

## Tenant Module (90xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 90001 | 404 | tenant not found | Tenant not found |
| 90002 | 409 | duplicate tenant | Tenant already exists |
| 90003 | 400 | invalid domain | Invalid domain format |
| 90004 | 403 | tenant is inactive | Tenant account is suspended |
| 90005 | 400 | cannot suspend expired tenant | Cannot suspend already expired tenant |
| 90006 | 400 | invalid tenant id | Invalid tenant ID |
| 90007 | 400 | tenant name is required | Tenant name is required |
| 90008 | 400 | tenant code is required | Tenant code is required |

---

## Market Module (150xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 150001 | 404 | market not found | Market not found |
| 150002 | 409 | duplicate market code | Market code already exists |
| 150003 | 400 | market code is required | Market code is required |
| 150004 | 400 | market name is required | Market name is required |
| 150005 | 400 | market currency is required | Market currency is required |
| 150006 | 400 | market is inactive | Market is inactive |
| 150007 | 400 | market is already default | Market is already the default |
| 150008 | 400 | cannot delete default market | Cannot delete the default market |

---

## Product Market Module (160xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 160001 | 404 | product market not found | Product-market association not found |
| 160002 | 409 | product already in market | Product is already in this market |
| 160003 | 400 | price is required for market | Price is required when adding to market |

---

## Inventory Module (170xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 170001 | 400 | insufficient available stock | Not enough available stock |
| 170002 | 400 | insufficient locked stock | Not enough locked stock to release |
| 170003 | 404 | sku not found | SKU not found in inventory |
| 170004 | 404 | warehouse not found | Warehouse not found |
| 170005 | 409 | duplicate warehouse code | Warehouse code already exists |

---

## Brand Module (180xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 180001 | 404 | brand not found | Brand not found |
| 180002 | 409 | brand name already exists | Brand name already exists |
| 180003 | 400 | brand has products | Cannot delete brand with products |

---

## Handling Errors in Frontend

```typescript
// Example error handling in Vue component
try {
  const response = await createProduct(data)
  ElMessage.success('创建成功')
} catch (error: any) {
  const errorCode = error.response?.data?.code

  switch (errorCode) {
    case 30001:
      ElMessage.error('商品名称不能为空')
      break
    case 30002:
      ElMessage.error('商品价格必须大于0')
      break
    case 40100:
      // Redirect to login
      router.push('/login')
      break
    default:
      ElMessage.error(error.response?.data?.msg || '操作失败')
  }
}
```

## Best Practices

1. **Always check the `code` field**, not just HTTP status
2. **Handle authentication errors** (401xx) by redirecting to login
3. **Display user-friendly messages** from the `msg` field
4. **Log error codes** for debugging purposes
5. **Implement retry logic** for transient errors (50300)