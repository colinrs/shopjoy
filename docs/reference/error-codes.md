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
| 40009 | 400 | 当前状态无法改价 | Cannot adjust price in current status |
| 40010 | 400 | 改价金额超出限制 | Price adjustment exceeds limit |
| 40011 | 400 | 导出数量超出限制 | Export count exceeds limit |
| 40012 | 400 | 改价原因不能为空 | Price adjustment reason is required |
| 40013 | 409 | 订单已被修改，请刷新后重试 | Order version conflict |
| 40014 | 400 | order cannot be cancelled in current status | Order cancellation not allowed |
| 40015 | 400 | cancel reason is required | Cancel reason is required |
| 40016 | 429 | payment reminder already sent recently | Too many payment reminders |
| 40017 | 400 | order already paid, cannot send reminder | Cannot remind paid order |
| 40018 | 400 | order cannot be reminded in current status | Cannot remind order |

---

## Payment Module (50xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 50001 | 404 | payment not found | Payment record not found |
| 50002 | 400 | invalid payment amount | Invalid payment amount |
| 50003 | 402 | payment failed | Payment processing failed |
| 50004 | 400 | payment already completed | Payment already completed |
| 50005 | 400 | payment expired | Payment session expired |
| 50006 | 400 | order not paid, cannot refund | Order must be paid before refund |
| 50007 | 400 | refund amount exceeds refundable | Refund amount too high |
| 50008 | 400 | refund reason is required | Refund reason required |
| 50009 | 500 | channel refund failed | Payment channel refund failed |
| 50010 | 404 | order not found | Order not found for payment |
| 50011 | 400 | order already fully refunded | Order already fully refunded |
| 50012 | 404 | transaction not found | Transaction not found |
| 50013 | 503 | payment channel unavailable | Payment channel unavailable |
| 50014 | 400 | refund not supported for this channel | Channel does not support refund |
| 50015 | 400 | currency not supported by channel | Currency not supported |
| 50016 | 404 | payment refund not found | Refund record not found |
| 50017 | 202 | payment requires additional action | Additional action required |
| 50018 | 400 | refund currency must match payment currency | Currency mismatch |
| 50019 | 409 | duplicate idempotency key | Idempotency key already used |
| 50020 | 409 | dispute created for this charge | Dispute created |
| 50021 | 400 | export limit exceeded | Export limit exceeded |
| 50022 | 401 | invalid stripe webhook signature | Webhook signature invalid |

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

## Fulfillment Module (120xxx)

### Shipment Errors (120xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 120001 | 404 | shipment not found | Shipment not found |
| 120002 | 400 | invalid tracking number | Invalid tracking number |
| 120003 | 400 | order already shipped | Order already shipped |
| 120004 | 400 | carrier is required | Carrier is required |
| 120005 | 400 | tracking number is required | Tracking number required |
| 120006 | 400 | invalid shipment status transition | Invalid status transition |
| 120007 | 400 | cannot cancel delivered shipment | Cannot cancel delivered shipment |
| 120008 | 404 | shipment item not found | Shipment item not found |
| 120009 | 404 | order not found | Order not found |
| 120010 | 400 | order cannot be shipped | Order cannot be shipped |
| 120011 | 409 | tracking number already exists | Duplicate tracking number |
| 120012 | 400 | shipment items are required | Shipment items required |
| 120013 | 400 | shipment item quantity exceeded order quantity | Quantity exceeds order |
| 120014 | 400 | invalid shipment items | Invalid shipment items |
| 120015 | 400 | export shipment count exceeds limit | Export limit exceeded |
| 120016 | 400 | shipment already cancelled | Shipment already cancelled |

### Refund Errors (1201xx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 120101 | 404 | refund not found | Refund not found |
| 120102 | 400 | invalid refund status | Invalid refund status |
| 120103 | 400 | cannot cancel refund in current status | Cannot cancel refund |
| 120104 | 409 | order already has pending refund | Duplicate pending refund |
| 120105 | 404 | order not found | Order not found |
| 120106 | 400 | order cannot be refunded | Order cannot be refunded |
| 120107 | 400 | refund period has expired | Refund period expired |
| 120108 | 400 | refund reason is required | Refund reason required |
| 120109 | 400 | reject reason is required | Reject reason required |
| 120110 | 400 | refund amount exceeds order amount | Refund amount too high |
| 120111 | 400 | order is not paid | Order is not paid |

### Carrier Errors (1202xx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 120201 | 404 | carrier not found | Carrier not found |
| 120202 | 409 | carrier code already exists | Duplicate carrier code |
| 120203 | 400 | carrier is inactive | Carrier is inactive |

### Refund Reason Errors (1203xx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 120301 | 404 | refund reason not found | Refund reason not found |
| 120302 | 409 | refund reason code already exists | Duplicate refund reason code |

---

## Review Module (210xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 210001 | 404 | review not found | Review not found |
| 210002 | 400 | review already has reply | Review already replied |
| 210003 | 400 | cannot reply to hidden review | Cannot reply hidden review |
| 210004 | 400 | invalid review status | Invalid review status |
| 210005 | 400 | review content exceeds limit | Review too long |
| 210006 | 400 | reply content exceeds limit | Reply too long |
| 210007 | 404 | reply not found | Reply not found |
| 210008 | 400 | cannot approve review in current status | Cannot approve review |
| 210009 | 400 | cannot hide review in current status | Cannot hide review |
| 210010 | 400 | cannot show review in current status | Cannot show review |
| 210011 | 400 | can only feature approved reviews | Cannot feature unapproved review |
| 210012 | 400 | review already deleted | Review already deleted |
| 210013 | 400 | reply content cannot be empty | Reply cannot be empty |
| 210014 | 400 | rating must be between 1 and 5 | Invalid rating |
| 210015 | 400 | batch operation requires at least one review id | Batch list empty |
| 210016 | 400 | batch operation limited to 100 reviews | Batch size limit |

---

## Shipping Module (230xxx)

### Shipping Template Errors (230xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 230001 | 404 | shipping template not found | Template not found |
| 230002 | 400 | template name is required | Template name required |
| 230003 | 400 | cannot delete template with zones | Template has zones |
| 230004 | 400 | cannot delete default template | Cannot delete default |
| 230005 | 409 | template name already exists | Duplicate template name |

### Shipping Zone Errors (2301xx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 230101 | 404 | shipping zone not found | Zone not found |
| 230102 | 400 | zone name is required | Zone name required |
| 230103 | 400 | zone regions are required | Regions required |
| 230104 | 400 | invalid fee type | Invalid fee type |
| 230105 | 400 | fee configuration is required | Fee config required |
| 230106 | 409 | region already assigned to another zone | Duplicate region |

### Shipping Mapping Errors (2302xx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 230201 | 404 | shipping mapping not found | Mapping not found |
| 230202 | 409 | mapping already exists | Mapping already exists |
| 230203 | 400 | invalid target type | Invalid target type |

### Shipping Calculator Errors (2303xx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 230301 | 400 | no matching zone for address | No zone for address |
| 230302 | 400 | no default shipping template configured | No default template |
| 230303 | 400 | items are required | Items required |
| 230304 | 400 | address is required | Address required |
| 230305 | 400 | invalid quantity in items | Invalid quantity |
| 230306 | 400 | invalid weight in items | Invalid weight |
| 230307 | 400 | invalid price in items | Invalid price |

---

## Upload Module (240xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 240001 | 400 | unsupported file type | File type not allowed |
| 240002 | 400 | file size exceeded | File too large |
| 240003 | 400 | invalid category | Invalid category |
| 240004 | 500 | upload failed | Upload failed |
| 240005 | 404 | file not found | File not found |

---

## Points Module (250xxx)

### Earn Rule Errors (250xxx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 250001 | 404 | earn rule not found | Earn rule not found |
| 250002 | 400 | invalid earn rule | Invalid earn rule |
| 250003 | 400 | earn rule is not active | Earn rule inactive |
| 250004 | 400 | earn rule is already active | Already active |
| 250005 | 400 | earn rule has expired | Earn rule expired |
| 250006 | 400 | earn rule has not started | Earn rule not started |

### Redeem Rule Errors (2501xx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 250101 | 404 | redeem rule not found | Redeem rule not found |
| 250102 | 400 | invalid redeem rule | Invalid redeem rule |
| 250103 | 400 | redeem rule is not active | Redeem rule inactive |
| 250104 | 400 | redeem rule is out of stock | Out of stock |
| 250105 | 400 | user has reached redemption limit | User limit reached |

### Points Account Errors (2502xx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 250201 | 404 | points account not found | Points account not found |
| 250202 | 400 | insufficient points balance | Insufficient points |
| 250203 | 400 | points are frozen | Points frozen |

### Points Transaction Errors (2503xx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 250301 | 404 | points transaction not found | Transaction not found |
| 250302 | 400 | export limit exceeded, maximum 10000 records | Export limit exceeded |

### Points Redemption Errors (2504xx)

| Code | HTTP Status | Message | Description |
|------|-------------|---------|-------------|
| 250401 | 404 | points redemption not found | Redemption not found |
| 250402 | 400 | invalid redemption status | Invalid redemption status |
| 250403 | 400 | cannot cancel redemption in current status | Cannot cancel redemption |

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