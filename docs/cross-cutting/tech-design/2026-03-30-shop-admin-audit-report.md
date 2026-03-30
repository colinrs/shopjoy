# ShopJoy Admin 项目审查报告

> **审查日期:** 2026-03-30
> **审查范围:** shop-admin (前端) + admin (后端 API)
> **审查专家:** 后端开发专家、前端开发专家、UI 设计专家

---

## 一、后端 API 问题

### 1.1 Swagger 生成 Bug (CRITICAL)

| 问题 | 位置 | 描述 |
|------|------|------|
| `/api/v1/carriers` 参数错误 | swagger admin.json:760-768 | `ListCarriersResp` 响应类型被错误当作查询参数 |
| `/api/v1/reviews/stats` 参数错误 | swagger admin.json:5257-5376 | `ReviewStatsResp` 字段被错误当作查询参数 |

**根因:** goctl swagger 工具在生成时将响应类型错误地当作请求参数处理。

**修复方案:** 检查 `.api` 文件定义，确保 request 和 response 类型正确分离，或在 `.api` 中明确定义空请求结构 `{}`。

### 1.2 后端符合规范检查

| 检查项 | 状态 | 说明 |
|--------|------|------|
| Error codes 定义 (pkg/code/code.go) | ✅ PASS | 所有错误码已正确定义，使用 `code.ErrXxx` 模式 |
| Monetary values (string 类型，元) | ✅ PASS | API 层使用 `string` 类型，如 `"1.99"` 表示 1.99 元 |
| Timestamps (RFC3339, time.Time) | ✅ PASS | 所有时间字段使用 `time.Time`，API 返回 RFC3339 格式字符串 |

**Error Code Ranges 验证:**
- Admin User: 10xxx
- User: 11xxx
- Product: 30xxx
- Category: 301xx
- Order: 40xxx
- Payment: 50xxx
- Cart: 60xxx
- Coupon: 70xxx
- Promotion: 80xxx
- Tenant: 90xxx
- Role: 100xxx
- Shop/Storefront: 110xxx
- Fulfillment: 120xxx
- Auth: 130xxx
- Market: 150xxx
- ProductMarket: 160xxx
- Inventory: 170xxx
- Brand: 180xxx
- SKU: 190xxx
- Review: 210xxx
- Dashboard: 220xxx
- Shipping: 230xxx
- Upload: 240xxx
- Points: 250xxx

---

## 二、前端问题汇总

### 2.1 关键 Bug (CRITICAL)

**价格格式化错误** - `shop-admin/src/views/products/index.vue:380`

```javascript
// 错误：除以 100
const formatPrice = (price: number) => (price / 100).toFixed(2)
```

```javascript
// 正确：API 返回的是字符串元格式 "1.99"，不需要除以 100
const formatPrice = (price: string) => price
```

**影响:** 所有商品价格会显示错误（如 "99.99" 元会显示为 "0.99" 元）

### 2.2 页面 Mock 数据清单

| 模块 | 页面 | 文件位置 | 问题行号 |
|------|------|----------|----------|
| Users | 用户详情 | `views/users/[id].vue` | Line 89 |
| Users | 用户地址列表 | `views/users/components/UserAddressList.vue` | Line 81 |
| Users | 用户操作日志 | `views/users/components/UserOperationLog.vue` | Line 91 |
| Admin Users | 管理员详情 | `views/admin-users/[id].vue` | Line 153 |
| **Payments** | 支付管理 | `views/payments/index.vue` | **Lines 221, 346, 378** |
| Points | 积分概览 | `views/points/dashboard/index.vue` | Lines 234, 255, 277, 297 |
| Points | 积分账户列表 | `views/points/accounts/index.vue` | Line 153 |
| Points | 积分账户详情 | `views/points/accounts/[id].vue` | Lines 149, 185 |
| Points | 积分赚取规则 | `views/points/earn-rules/index.vue` | Line 225 |
| Points | 积分兑换规则 | `views/points/redeem-rules/index.vue` | Line 223 |
| Points | 积分交易记录 | `views/points/transactions/index.vue` | Line 128 |
| Points | 积分兑换记录 | `views/points/redemptions/index.vue` | Line 183 |
| Fulfillment | 发货列表 | `views/fulfillment/shipments/index.vue` | Lines 292, 428, 460 |
| Fulfillment | 发货详情 | `views/fulfillment/shipments/[id]/index.vue` | Line 357 |
| Fulfillment | 退款列表 | `views/fulfillment/refunds/index.vue` | Lines 267, 441 |
| Fulfillment | 退款详情 | `views/fulfillment/refunds/[id]/index.vue` | Line 378 |
| Fulfillment | 履约统计 | `views/fulfillment/statistics/index.vue` | Line 389 |

### 2.3 硬编码问题

| 位置 | 问题描述 |
|------|----------|
| `products/index.vue:40-45, 227-230, 563-568` | 硬编码分类选项 (electronics, clothing, home, sports) 而非从 API 获取 |
| `fulfillment/shipments/index.vue:429-438` | 硬编码承运商列表 fallback |
| `fulfillment/refunds/index.vue:412-420` | 硬编码退款原因 fallback |
| `orders/index.vue:431-445` | `orderStats.pending_payment: 0` 硬编码 |

### 2.4 功能残缺 (Stub Implementations)

| 页面 | 功能 | 问题 |
|------|------|------|
| Dashboard | `查看全部` 按钮 | 无点击处理器 |
| Dashboard | `处理` 按钮 | 无点击处理器 |
| Products | `预览` | 仅显示消息，无实际预览 |
| Products | `复制` | Stub 实现 |
| Products | `批量删除` | Stub 实现 |
| Payments | `导出` | Stub 实现 |

### 2.5 语言不一致 (CRITICAL)

| 页面 | 语言使用 |
|------|----------|
| Dashboard | 混合 (统计中文，表格英文) |
| Orders | 全英文标签 |
| Users | 全中文 |
| Payments | 全英文 |
| Promotions | 混合 |
| Storefront | 混合 |

**问题:** 无 i18n 实现，所有文本硬编码。

### 2.6 缺失功能

| 问题 | 位置 |
|------|------|
| 无 i18n 实现 | 所有页面硬编码文本 |
| 订单状态筛选缺少 `completed` | `orders/index.vue:54` |
| 重复定义 `ShipmentStatus` | `shipments/index.vue` vs `fulfillment.ts` |

### 2.7 枚举不匹配

**Frontend OrderStatus** (`shop-admin/src/api/order.ts:6`):
```typescript
export type OrderStatus = 'pending_payment' | 'paid' | 'pending_shipment' | 'shipped' | 'completed' | 'cancelled' | 'refunding' | 'refunded'
```

**Backend Fulfillment Status** (from `fulfillment.api`):
- Shipment Status: `pending`, `shipped`, `in_transit`, `delivered`, `failed`, `cancelled`
- Fulfillment Status: `pending`, `partial_shipped`, `shipped`, `delivered`
- Refund Status: `pending`, `approved`, `rejected`, `completed`, `cancelled`

---

## 三、UI/UX 问题

### 3.1 语言混用
- 同一页面内中英文混杂
- 模块之间语言不统一

### 3.2 按钮无处理器
- Dashboard 的"查看全部"按钮
- Dashboard 的"处理"按钮

### 3.3 Mock 数据
- 多页面使用 Mock 而非真实 API

### 3.4 表单验证
- 错误处理缺失 (catch blocks 无用户提示)

### 3.5 类型不安全
- 混用 `any`
- 类型不匹配 (String vs Number for status switches)

---

## 四、API Client 文件清单

| API 文件 | 功能数 | 状态 |
|----------|--------|------|
| `admin-user.ts` | 16 | ✅ |
| `brand.ts` | 10 | ✅ |
| `category.ts` | 12 | ✅ |
| `dashboard.ts` | 7 | ✅ |
| `fulfillment.ts` | 15 | ✅ |
| `inventory.ts` | - | 未审查 |
| `market.ts` | 5 | ✅ |
| `order.ts` | 12 | ✅ |
| `payment.ts` | - | 未审查 |
| `points.ts` | 30+ | ✅ |
| `product.ts` | 18 | ✅ |
| `promotion.ts` | 18 | ✅ |
| `review.ts` | - | 未审查 |
| `shipping.ts` | - | 未审查 |
| `shop.ts` | - | 未审查 |
| `storefront.ts` | - | 未审查 |
| `upload.ts` | - | 未审查 |
| `user.ts` | 14 | ✅ |

---

## 五、实施计划

### Phase 1: 紧急修复 (Critical Bugs)

| # | 任务 | 改动文件 | 优先级 |
|---|------|----------|--------|
| 1.1 | 修复价格格式化 bug | `shop-admin/src/views/products/index.vue:380` | P0 |
| 1.2 | 修复 Swagger `/api/v1/carriers` 生成问题 | `admin/desc/shipping.api` | P0 |
| 1.3 | 修复 Swagger `/api/v1/reviews/stats` 生成问题 | `admin/desc/review.api` | P0 |

### Phase 2: 后端 API 一致性

| # | 任务 | 改动文件 |
|---|------|----------|
| 2.1 | 审查所有 `.api` 文件与 swagger 一致性 | 21 `.api` files |
| 2.2 | 修复 `.api` 定义中的 swagger 生成问题 | affected `.api` files |
| 2.3 | 重新生成 swagger 和 types | auto-generated |

### Phase 3: 前端 API 对接 (Mock → Real)

| # | 模块 | 改动文件 |
|---|------|----------|
| 3.1 | Payments 模块 | `payments/index.vue`, `api/payment.ts` |
| 3.2 | Points 模块 | 10+ files in `points/` |
| 3.3 | Fulfillment 模块 | 6 files in `fulfillment/` |
| 3.4 | Users 模块 | 4 files in `users/` |
| 3.5 | Admin Users 模块 | 2 files in `admin-users/` |

### Phase 4: 前端功能完善

| # | 任务 | 改动文件 |
|---|------|----------|
| 4.1 | Dashboard 按钮处理器 | `dashboard/index.vue` |
| 4.2 | Products 预览/复制/批量删除 | `products/index.vue` |
| 4.3 | Orders stats API 对接 | `orders/index.vue` |
| 4.4 | 导出功能实现 | 各模块 |

### Phase 5: 前端类型与规范

| # | 任务 | 改动文件 |
|---|------|----------|
| 5.1 | 移除硬编码分类 - 改为 API 获取 | `products/index.vue` |
| 5.2 | 修复枚举不匹配问题 | types + views |
| 5.3 | 移除重复 `ShipmentStatus` 定义 | use from `fulfillment.ts` |
| 5.4 | 添加统一错误处理 | catch blocks |

### Phase 6: 语言标准化

| # | 任务 | 说明 |
|---|------|------|
| 6.1 | 建立语言标准 | 模块级统一中/英文 |
| 6.2 | 统一现有页面语言 | 30+ files |

---

## 六、执行顺序

```
Phase 1 (Critical Bugs)
    ↓
Phase 2 (Backend API consistency)
    ↓
Phase 3 (Frontend API integration) ─ 可并行执行
    ↓
Phase 4 (Frontend features)
    ↓
Phase 5 (Types & standards)
    ↓
Phase 6 (Language consistency)
```

---

## 七、审查文件清单

### 后端文件
- `/admin/desc/*.api` (21 files)
- `/admin/swagger/admin.json`
- `/admin/pkg/code/code.go`

### 前端文件
- `/shop-admin/src/views/**/*.vue` (30+ pages)
- `/shop-admin/src/api/*.ts` (18 files)
- `/shop-admin/src/components/**`

---

## 八、后续行动

- [ ] 用户确认本报告内容
- [ ] 确定 Phase 1 紧急修复启动时间
- [ ] 分配 Phase 3 各模块负责人
- [ ] 建立 i18n 实施标准
