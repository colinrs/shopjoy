# ShopJoy Admin 前后端一致性审查报告

**审查日期**: 2026-03-27
**审查范围**: shop-admin (前端) + admin (后端)
**审查目的**: 识别前后端功能差异、字段不一致、未实现功能

---

## 一、关键问题汇总

### 🔴 严重问题 (Critical)

| # | 问题 | 模块 | 影响 |
|---|------|------|------|
| 1 | **Shop Settings 后端逻辑全部为空** | shop | 所有店铺设置功能无法使用 |
| 2 | **Role 模块前端未实现** | role | 角色管理完全缺失 |
| 3 | **Storefront 模块前端未实现** | storefront | 店铺装修、主题、页面、SEO全部缺失 |
| 4 | **多个功能使用 Mock 数据** | fulfillment, payments, points | 显示假数据 |
| 5 | **Stripe Webhook 无签名验证** | payment | 安全漏洞 |

### 🟠 高优先级问题 (High)

| # | 问题 | 模块 | 描述 |
|---|------|------|------|
| 6 | **删除商品 API 缺失** | product | 后端无 DELETE endpoint |
| 7 | **OrderCancel 路径参数不匹配** | order | 前端用 orderId string，后端用 id int64 |
| 8 | **UserAddress 字段名不匹配** | user | 前端 recipient_name vs 后端 name |
| 9 | **SaveDraftRequest 缺少 ID 字段** | storefront | 前端未传 page ID |
| 10 | **ShipOrderResponse 缺少 ID 字段** | fulfillment | 后端返回的 shipment_id 前端未接收 |
| 11 | **促销 product_ids 类型不匹配** | promotion | 前端 number[] vs 后端 JSON string |

### 🟡 中优先级问题 (Medium)

| # | 问题 | 模块 | 描述 |
|---|------|------|------|
| 12 | **多出 int/number 类型不匹配** | 多个 | price, stock 等字段 |
| 13 | **Points condition_value 类型不匹配** | points | 前端 object vs 后端 string |
| 14 | **Category ProductCount 始终为 0** | category | 占位符未实现 |
| 15 | **SKU Prefix 未从数据库获取** | product | 自动生成 SKU 可能不符合规范 |
| 16 | **Role 删除前未检查是否被使用** | role | 数据完整性风险 |
| 17 | **多个 Export 功能是占位符** | 多个 | 只有 ElMessage.success 无实际功能 |

---

## 二、前端未开发页面审查

| 页面/功能 | 状态 | 说明 |
|----------|------|------|
| **角色管理** (Role) | ❌ 完全缺失 | 后端有 8 个 API，前端 0 个 API |
| **店铺装修** (Storefront) | ❌ 完全缺失 | 主题/页面/SEO 均无 API 调用 |
| **Business Hours 设置** | ❌ 未实现 | 后端有 API，前端无 API |
| **店铺运费设置** | ⚠️ 部分实现 | 前端有 UI，后端逻辑空 |
| **店铺通知设置** | ⚠️ 部分实现 | 前端有 UI，后端逻辑空 |
| **店铺支付设置** | ⚠️ 部分实现 | 前端有 UI，后端逻辑空 |

---

## 三、后端 API 未被前端使用

| API 文件 | 后端 API 数 | 前端使用数 | 未使用数 |
|----------|------------|-----------|---------|
| role.api | 8 | 1 | **7** |
| storefront.api | 16 | 0 | **16** |
| shop.api | 9 | 7 | **2** (business-hours) |
| fulfillment.api | 21 | 20 | **1** (update shipment) |
| product.api | 18 | 17 | **1** (get single SKU) |

### 3.1 Role 模块详情

| Endpoint | Handler | 前端状态 |
|----------|---------|---------|
| `POST /api/v1/roles` | CreateRoleHandler | ❌ 未实现 |
| `GET /api/v1/roles/:id` | GetRoleHandler | ❌ 未实现 |
| `PUT /api/v1/roles/:id` | UpdateRoleHandler | ❌ 未实现 |
| `DELETE /api/v1/roles/:id` | DeleteRoleHandler | ❌ 未实现 |
| `PUT /api/v1/roles/:id/status` | UpdateRoleStatusHandler | ❌ 未实现 |
| `PUT /api/v1/roles/:id/permissions` | UpdateRolePermissionsHandler | ❌ 未实现 |
| `GET /api/v1/permissions` | ListPermissionsHandler | ❌ 未实现 |

### 3.2 Storefront 模块详情

整个模块 (16 个 API) 前端均未实现：

- Theme APIs: switch_theme, update_theme_config, list_themes, get_theme, create_theme, delete_theme
- Page APIs: list_pages, get_page, create_page, update_page, delete_page, save_draft
- Version APIs: list_versions, rollback_version
- SEO APIs: get_seo_config, update_seo_config

---

## 四、字段定义不一致

### 4.1 类型不匹配 (int64 vs number)

| 模块 | 字段 | 前端类型 | 后端类型 |
|------|------|---------|---------|
| Product | price, cost_price | number | int64 |
| Product | stock | number | int |
| SKU | price | number | int64 |
| SKU | stock | number | int |
| Order | shipped_qty | number | int |
| Category | status | number | int8 |
| Points | threshold | number\|null | *int64 |

### 4.2 字段名不匹配

| 模块 | 场景 | 前端字段 | 后端字段 |
|------|------|---------|---------|
| UserAddress | 收货人 | recipient_name | name |
| Refund | 订单项 | order_items (前端期望) | 后端无此字段 |

### 4.3 序列化不匹配

| 模块 | 字段 | 前端类型 | 后端类型 |
|------|------|---------|---------|
| Promotion | product_ids | number[] | JSON string |
| Promotion | category_ids | number[] | JSON string |
| Points | condition_value | object | string |

### 4.4 API 路径不匹配

| 模块 | 功能 | 前端 URL | 问题 |
|------|------|---------|------|
| Order | cancelOrder | `/api/v1/orders/:orderId/cancel` | orderId 是 string，后端期望 int64 id |
| Order | batchCancelOrders | N/A | 未实现 |

---

## 五、枚举定义缺失

### 5.1 前端定义的枚举后端无约束

| 模块 | 字段 | 前端定义 | 后端定义 |
|------|------|---------|---------|
| Order | status | 'pending'\|'paid'\|'to_ship'\|'shipped'\|'completed'\|'cancelled' | string (无约束) |
| Order | fulfillment_status | 0\|1\|2\|3 | int8 (无约束) |
| Promotion | type | 'discount'\|'full_reduce'\|'flash_sale'\|'bundle' | string (无约束) |
| Promotion | status | 'draft'\|'active'\|'inactive'\|'expired' | string (无约束) |
| Promotion | discount_type | 'percentage'\|'fixed_amount' | string (无约束) |
| Promotion | coupon_type | 'fixed_amount'\|'percentage' | string (无约束) |
| Points | earn_rule_scenario | 'ORDER_PAYMENT'\|'SIGN_IN'\|'PRODUCT_REVIEW'\|'FIRST_ORDER' | string (无约束) |
| Points | calculation_type | 'FIXED'\|'RATIO'\|'TIERED' | string (无约束) |
| Refund | status | 0\|1\|2\|3\|4\|5 | int8 (无约束) |
| Review | status | 'pending'\|'approved'\|'hidden' | string (无约束) |
| Shipping | fee_type | 'fixed'\|'by_count'\|'by_weight'\|'free' | string (无约束) |

### 5.2 建议添加的枚举定义位置

枚举应添加到对应的 `.api` 文件中，使用 comments 或 go:generate 方式定义：

```go
// OrderStatus values: pending, paid, to_ship, shipped, completed, cancelled
// PromotionType values: discount, full_reduce, flash_sale, bundle
// 等等...
```

---

## 六、后端空实现占位符

### 6.1 Shop Settings 全部为空

| 文件 | 状态 |
|------|------|
| `internal/logic/shop/get_shop_settings_logic.go` | TODO 占位符 |
| `internal/logic/shop/update_shop_settings_logic.go` | TODO 占位符 |
| `internal/logic/shop/get_business_hours_logic.go` | TODO 占位符 |
| `internal/logic/shop/update_business_hours_logic.go` | TODO 占位符 |
| `internal/logic/shop/get_notification_settings_logic.go` | TODO 占位符 |
| `internal/logic/shop/update_notification_settings_logic.go` | TODO 占位符 |
| `internal/logic/shop/get_payment_settings_logic.go` | TODO 占位符 |
| `internal/logic/shop/update_payment_settings_logic.go` | TODO 占位符 |
| `internal/logic/shop/get_shipping_settings_logic.go` | TODO 占位符 |
| `internal/logic/shop/update_shipping_settings_logic.go` | TODO 占位符 |

### 6.2 其他占位符

| 文件 | 问题 |
|------|------|
| `internal/logic/categories/get_category_tree_logic.go` | getCategoryProductCount 始终返回 0 |
| `internal/logic/reviews/get_product_stats_logic.go` | ReplyCount, ReplyRate 硬编码为 0 |
| `internal/logic/products/create_sku_logic.go` | SKU Prefix 未从数据库获取 |
| `internal/logic/roles/delete_role_logic.go` | 删除前未检查 role 是否被使用 |
| `internal/logic/webhooks/stripe_webhook_logic.go` | 无签名验证 (安全漏洞) |
| `internal/logic/themes/switch_theme_logic.go` | UserName 为空 |

---

## 七、前端 Mock 数据和假功能

### 7.1 使用 Mock 数据的页面

| 文件 | 组件 | 说明 |
|------|------|------|
| `shop-admin/src/views/payments/index.vue` | Transaction list | API 失败时回退到 mock 数据 |
| `shop-admin/src/views/fulfillment/refunds/index.vue` | Refund list | 硬编码 mock 数据 |
| `shop-admin/src/views/fulfillment/shipments/index.vue` | Shipment list | 硬编码 mock 数据 |
| `shop-admin/src/views/points/transactions/index.vue` | Points transactions | API 失败时回退到 mock 数据 |
| `shop-admin/src/views/points/accounts/[id].vue` | Account detail | API 失败时回退到 mock 数据 |
| `shop-admin/src/views/admin-users/[id].vue` | Admin user detail | API 失败时回退到 mock 数据 |

### 7.2 假功能 (只有消息，无实际实现)

| 文件 | 功能 | 问题 |
|------|------|------|
| `products/index.vue` | handleDelete | 显示成功但未调用 API |
| `products/index.vue` | handleBatchDelete | 显示成功但未调用 API |
| `products/index.vue` | handleCommand(copy) | 只有成功消息 |
| `payments/index.vue` | handleExport | 只有成功消息 |
| `shipping/index.vue` | handleExport | 只有成功消息 |
| `fulfillment/refunds/index.vue` | handleExport | 只有成功消息 |
| `fulfillment/shipments/index.vue` | handleExport | 只有成功消息 |
| `points/transactions/index.vue` | handleExport | 显示"功能开发中" |
| `storefront/pages/index.vue` | createPage | 显示"即将上线" |

---

## 八、UI/UX 问题

### 8.1 表单验证问题

| 文件 | 问题 |
|------|------|
| `products/index.vue` | Push to Market 对话框有 required 但无 rules 定义 |
| `users/index.vue` | 编辑对话框无验证规则 |

### 8.2 操作不便问题

| 问题 | 位置 | 建议 |
|------|------|------|
| 表格无法排序 | products, orders, payments | 添加排序选项 |
| 无快速筛选 | orders, payments | 添加日期范围芯片、状态芯片 |
| 无键盘快捷键 | 全局 | 添加 Esc 关闭、Enter 搜索等 |
| 弹窗关闭无警告 | products, promotions | 添加未保存更改提示 |
| 分页无确认 | orders | 输入页码后需确认 |

### 8.3 响应式问题

| 文件 | 问题 |
|------|------|
| `products/index.vue` | search-input 固定 280px |
| `shipping/index.vue` | grid 布局在移动端适配问题 |

---

## 九、数据库 Schema 问题

| 问题 | 位置 | 说明 |
|------|------|------|
| Storefront schema 缺失 | `sql/storefront/schema.sql` 不存在 | storefront 表定义在 shop/schema.sql 中，违反项目规范 |

---

## 十、开发计划

### Phase 1: 核心功能修复 (Critical)

#### 1.1 Shop Settings 后端实现 (P0)

```
负责人: Go Backend
工时: 3-4 天
依赖: 无
任务:
- 实现 get/update_shop_settings_logic.go
- 实现 get/update_business_hours_logic.go
- 实现 get/update_notification_settings_logic.go
- 实现 get/update_payment_settings_logic.go
- 实现 get/update_shipping_settings_logic.go
```

#### 1.2 Role 模块前后端实现 (P0)

```
负责人: Fullstack
工时: 5-6 天
依赖: 无
任务:
- 前端: 创建 role.api.ts
- 前端: 创建 views/roles/ 页面 (list, create, edit, permissions)
- 后端: 确保 role.api 所有 8 个 endpoint 正常工作
- 后端: 添加删除前 role-in-use 检查
```

#### 1.3 Storefront 模块前后端实现 (P0)

```
负责人: Fullstack
工时: 8-10 天
依赖: shop settings
任务:
- 前端: 完善 storefront.api.ts
- 前端: 实现 themes/index.vue 主题管理
- 前端: 实现 pages/ 页面管理 (list, edit, preview)
- 前端: 实现 seo/ SEO 设置
- 前端: 实现 decorations 装修配置
- 后端: 修复 switch_theme_logic.go username TODO
```

### Phase 2: 数据一致性修复 (High)

#### 2.1 API 路径和字段修复

```
负责人: Fullstack
工时: 2-3 天
依赖: 无
任务:
- 修复 order cancel API 路径 (orderId string -> id int64)
- 修复 UserAddress recipient_name -> name
- 修复 SaveDraftRequest 添加 ID 字段
- 修复 ShipOrderResponse 添加 shipment_id 字段
- 添加 product delete API
```

#### 2.2 类型统一修复

```
负责人: Fullstack
工时: 2-3 天
依赖: 2.1
任务:
- 统一 price/cost_price 类型 (统一使用 int64 后端)
- 统一 promotion product_ids/category_ids 序列化
- 统一 points condition_value 序列化
- 添加所有枚举到 backend API 定义
```

### Phase 3: 枚举和文档完善 (Medium)

#### 3.1 后端 API 枚举定义

```
负责人: Go Backend
工时: 1-2 天
依赖: 无
任务:
- order.api: 添加 OrderStatus 枚举
- promotion.api: 添加 PromotionType, DiscountType 等枚举
- fulfillment.api: 添加 RefundStatus 等枚举
- review.api: 添加 ReviewStatus 枚举
- points.api: 添加 EarnRuleScenario 等枚举
```

#### 3.2 前端 TypeScript 类型对齐

```
负责人: Frontend
工时: 1-2 天
依赖: 3.1
任务:
- 更新所有 api/*.ts 使用后端定义的枚举
- 添加常量文件 constant.ts
```

### Phase 4: UI/UX 优化 (Medium)

#### 4.1 移除假功能

```
负责人: Frontend
工时: 1-2 天
依赖: Phase 2
任务:
- 实现 export 功能或移除按钮
- 实现 copy/duplicate 功能
- 实现 batch delete API 调用
- 移除 mock 数据 fallback
```

#### 4.2 操作体验改进

```
负责人: Frontend
工时: 2-3 天
依赖: 无
任务:
- 添加表格排序选项
- 添加快速筛选 chips (日期、状态)
- 添加键盘快捷键
- 添加表单未保存警告
- 改进分页 UX
```

### Phase 5: 遗留问题修复 (Low)

#### 5.1 后端占位符修复

```
负责人: Go Backend
工时: 1-2 天
依赖: 无
任务:
- 实现 getCategoryProductCount
- 实现 SKU Prefix 数据库获取
- 添加 Stripe webhook 签名验证
- 完善 review stats 计算
```

---

## 十一、优先级矩阵

| 优先级 | 任务 | 预计工时 | 风险 |
|--------|------|---------|------|
| P0 | Shop Settings 后端实现 | 3-4d | 低 |
| P0 | Role 模块前后端 | 5-6d | 中 |
| P0 | Storefront 前端实现 | 8-10d | 高 |
| P1 | API 路径/字段修复 | 2-3d | 低 |
| P1 | 类型统一修复 | 2-3d | 中 |
| P2 | 枚举定义完善 | 1-2d | 低 |
| P2 | UI 假功能移除 | 1-2d | 低 |
| P2 | 操作体验优化 | 2-3d | 低 |
| P3 | 遗留占位符修复 | 1-2d | 低 |

**总预计工时: 25-35 人/天**

---

## 十二、建议执行顺序

1. **第一周**: Shop Settings 后端 + API 修复
2. **第二周**: Role 模块实现
3. **第三-四周**: Storefront 前端实现 (最大工作量)
4. **第五周**: 类型统一 + 枚举定义
5. **第六周**: UI/UX 优化
6. **可选**: 遗留占位符修复

---

## 附录 A: 审查文件清单

### 后端 API 文件

```
admin/desc/
├── admin.api
├── admin_user.api
├── auth.api
├── brand.api
├── category.api
├── dashboard.api
├── fulfillment.api
├── inventory.api
├── market.api
├── payment.api
├── points.api
├── product.api
├── product_market.api
├── promotion.api
├── review.api
├── role.api
├── shipping.api
├── shop.api
├── storefront.api
├── upload.api
└── user.api
```

### 前端 API 文件

```
shop-admin/src/api/
├── admin-user.ts
├── brand.ts
├── category.ts
├── dashboard.ts
├── fulfillment.ts
├── inventory.ts
├── market.ts
├── order.ts
├── payment.ts
├── points.ts
├── product.ts
├── promotion.ts
├── review.ts
├── shipping.ts
├── shop.ts
├── storefront.ts
├── upload.ts
└── user.ts
```

### 前端页面文件

```
shop-admin/src/views/
├── admin-users/
├── brands/
├── categories/
├── dashboard/
├── fulfillment/
│   ├── refunds/
│   └── shipments/
├── inventory/
├── login/
├── orders/
├── payments/
├── points/
│   ├── accounts/
│   ├── earn-rules/
│   ├── redeem-rules/
│   ├── redemptions/
│   └── transactions/
├── products/
├── promotions/
├── reviews/
├── settings/
│   └── markets/
├── shipping/
├── shop/
├── storefront/
│   ├── pages/
│   └── themes/
└── users/
```

---

## 附录 B: 安全问题

### Stripe Webhook 签名验证缺失

**位置**: `admin/internal/logic/webhooks/stripe_webhook_logic.go`

**问题**: Webhook 处理未验证 Stripe 签名，存在安全风险

**建议**: 实现 Stripe 签名验证，参考 Stripe 官方文档
