# Order Module Product Requirements Document

## Document Information

| Item | Value |
|------|-------|
| Document Title | Fulfillment Module Enhancement PRD (Admin) |
| Version | 1.1.0 |
| Status | Draft |
| Created | 2026-03-22 |
| Updated | 2026-03-23 |
| Author | Product Team |
| Module Location | admin/internal/domain/fulfillment/ |

---

## Product Overview

### Summary

本 PRD 描述对现有 Fulfillment 模块的功能扩展，新增订单管理能力。Fulfillment 模块现有发货、退款、物流公司管理等核心功能，本次扩展将新增：商家备注、改价、订单导出、今日统计等功能，完善商家的订单运营体验。

**现有功能（Fulfillment 模块 v1.0）**：
- 发货单管理（创建、列表、详情、更新）
- 退款管理（列表、详情、批准、拒绝）
- 物流公司管理
- 订单履约视图（订单列表、订单详情、订单发货、履约摘要）

**本次新增功能**：
- 商家备注（内部备注，仅商家可见）
- 改价功能（待付款订单可调整金额）
- 订单导出（CSV 同步导出）
- 今日统计（今日订单数、今日 GMV）

### Problem Statement

Fulfillment 模块已提供订单履约视图，但商家在日常运营中仍有以下痛点：

- **缺少改价能力**：待付款订单无法调整金额（如首单优惠补充、运费减免等场景）
- **缺少内部备注**：无法记录订单相关的内部沟通信息
- **缺少今日数据**：无法快速查看今日订单数和 GMV
- **缺少导出能力**：无法导出订单数据进行离线分析

### Solution Overview

扩展现有 Fulfillment 模块，新增以下功能：

1. **商家备注** - 支持对订单添加内部备注（仅商家可见）
2. **改价功能** - 支持对待付款订单调整金额（限制 20% 以内，需填写原因）
3. **今日统计** - 新增今日订单数和今日 GMV 统计
4. **订单导出** - 支持按筛选条件导出订单 CSV（同步下载）

---

## Goals

### Business Goals

| ID | Goal | Metric | Target |
|----|------|--------|--------|
| BG-001 | 支持灵活定价 | 改价使用率 | 10% 待付款订单使用改价 |
| BG-002 | 提升运营效率 | 备注使用率 | 30% 订单有内部备注 |
| BG-003 | 数据可导出 | 导出使用率 | 20% 商家每周使用导出 |

### User Goals

| ID | Goal | User Type |
|----|------|-----------|
| UG-001 | 调整待付款订单价格 | Merchant Admin |
| UG-002 | 为订单添加内部备注 | Merchant Admin / Customer Service |
| UG-003 | 查看今日订单和 GMV | Merchant Admin / Store Manager |
| UG-004 | 导出订单数据 | Merchant Admin |

### Non-Goals

| ID | Non-Goal | Reason |
|----|----------|--------|
| NG-001 | 部分退款（商品级别） | Phase 2 feature |
| NG-002 | 异步导出 | MVP 使用同步导出，大数据量场景 Phase 2 考虑 |
| NG-003 | 批量改价 | Phase 2 feature |
| NG-004 | 买家端订单页面 | 属于 shop service 范围 |

---

## User Personas

### Key User Types

| Persona | Description | Primary Use Cases |
|---------|-------------|-------------------|
| Merchant Admin | 商家管理员 | 改价、备注、查看统计、导出订单 |
| Customer Service | 客服人员 | 备注订单、处理退款 |

### Role-Based Access

| Role | 改价 | 备注 | 导出 | 统计 |
|------|------|------|------|------|
| Tenant Admin | ✅ | ✅ | ✅ | ✅ |
| Tenant Operations Manager | ❌ | ✅ | ✅ | ✅ |
| Tenant Customer Service | ❌ | ✅ | ❌ | ✅ |

---

## Functional Requirements

### 现有功能（Fulfillment v1.0）

以下功能已实现，本次不做修改：

| 功能 | API | 说明 |
|------|-----|------|
| 订单列表 | GET /api/v1/orders | 支持订单状态、履约状态、退款状态筛选 |
| 订单详情 | GET /api/v1/orders/:id | 含订单信息、收货地址、商品、发货记录、退款记录 |
| 订单发货 | PUT /api/v1/orders/:id/ship | 创建发货单 |
| 履约摘要 | GET /api/v1/orders/fulfillment-summary | 待发货、部分发货、已发货、待退款等统计 |
| 退款列表 | GET /api/v1/refunds | 支持筛选 |
| 退款详情 | GET /api/v1/refunds/:id | 含订单信息、退款原因、凭证图片 |
| 批准退款 | PUT /api/v1/refunds/:id/approve | 批准退款申请 |
| 拒绝退款 | PUT /api/v1/refunds/:id/reject | 拒绝退款申请（需填写原因） |

### 新增功能：商家备注

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-001 | 更新商家备注 | PUT /api/v1/orders/:id/remark | P0 |
| FR-002 | 备注长度限制 | 最大 500 字符 | P0 |
| FR-003 | 备注可见性 | 仅商家后台可见，买家端不展示 | P0 |

### 新增功能：改价

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-010 | 改价接口 | PUT /api/v1/orders/:id/adjust-price | P0 |
| FR-011 | 订单状态限制 | 仅 `pending_payment` 状态可改价 | P0 |
| FR-012 | 改价金额限制 | 调整金额绝对值 ≤ 原金额 × 20% | P0 |
| FR-013 | 改价原因必填 | 需填写改价原因（审计要求） | P0 |
| FR-014 | 保留原始金额 | original_amount 存储改价前金额 | P0 |
| FR-015 | 乐观锁控制 | 使用 version 字段防止并发修改 | P0 |

### 新增功能：今日统计

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-020 | 今日订单数 | 今日创建的订单数量 | P0 |
| FR-021 | 今日 GMV | 今日已支付订单的总金额 | P0 |
| FR-022 | 统计接口扩展 | 扩展现有 fulfillment-summary 接口 | P0 |

### 新增功能：订单导出

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-030 | 导出接口 | GET /api/v1/orders/export | P0 |
| FR-031 | 筛选条件同步 | 导出时应用当前列表的筛选条件 | P0 |
| FR-032 | 导出格式 | CSV 格式，UTF-8 编码 | P0 |
| FR-033 | 导出限制 | 最大 10,000 条记录 | P0 |
| FR-034 | 同步下载 | 直接返回文件流，不存储 | P0 |

---

## User Experience

### 新增：改价操作

1. 在订单详情页，点击「改价」按钮
2. 弹出改价对话框：
   - 显示原订单金额
   - 输入调整金额（可正可负，单位：分）
   - 实时预览新金额
   - 输入改价原因（必填）
3. 点击「确认」
4. 系统校验：
   - 订单状态为 `pending_payment`
   - 调整金额不超过原金额 20%
5. 更新订单金额，创建审计记录
6. 显示成功提示

### 新增：添加备注

1. 在订单详情页，点击「备注」按钮
2. 弹出备注对话框：
   - 显示当前备注（如有）
   - 输入/编辑备注内容
3. 点击「保存」
4. 备注保存成功

### 新增：查看今日统计

1. 在订单列表页顶部，统计卡片区域
2. 显示今日订单数和今日 GMV
3. 点击卡片可跳转到今日订单列表

### 新增：导出订单

1. 在订单列表页，设置筛选条件
2. 点击「导出」按钮
3. 浏览器直接下载 CSV 文件
4. 文件名格式：`orders_YYYYMMDD_HHMMSS.csv`

---

## Narrative

商家小明早上打开后台，查看订单列表页的统计卡片，看到今日已有 45 笔订单，GMV 达到 12,890 元。

一位老客户咨询是否可以给首单额外优惠，小明查看该客户的待付款订单，点击「改价」按钮，将订单金额减少 10 元，并在原因栏填写「首单额外优惠」。系统校验通过后，订单金额更新，客户完成支付。

下午，小明处理完发货后，为几个需要特别关注的订单添加了备注：「客户要求周末配送」「易碎品，注意包装」。

下班前，小明筛选出今天的订单，点击「导出」按钮，下载了 CSV 文件，准备在 Excel 中做进一步分析。

---

## Success Metrics

### User-Centric Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| 改价使用率 | 待付款订单使用改价功能的比例 | 10% |
| 备注使用率 | 订单添加备注的比例 | 30% |
| 导出使用率 | 商家每周使用导出的比例 | 20% |

### Technical Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| 改价接口响应时间 | P95 | <200ms |
| 导出接口响应时间 | 1000条记录 P95 | <2s |
| 错误率 | 改价/备注/导出操作失败率 | <0.1% |

---

## Technical Considerations

### 现有架构

Fulfillment 模块采用 DDD 架构：

```
admin/internal/domain/fulfillment/
├── entity.go          # 领域实体：Shipment, Refund, Carrier, RefundReason
├── repository.go      # 仓储接口
└── service.go         # 领域服务
```

现有 API 定义：`admin/desc/fulfillment.api`

### 数据库变更

#### orders 表新增字段

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| fulfillment_status | TINYINT | NO | 0 | 履约状态: 0-待发货, 1-部分发货, 2-已发货, 3-已送达 |
| refund_status | TINYINT | NO | 0 | 退款状态: 0-无, 1-待处理, 2-已批准, 3-已拒绝, 4-已完成 |
| merchant_remark | VARCHAR(500) | NO | '' | 商家内部备注 |
| original_amount | BIGINT | NO | 0 | 改价前原金额（分） |
| adjust_amount | BIGINT | NO | 0 | 改价金额（分），正数加价，负数减价 |
| adjust_reason | VARCHAR(200) | NO | '' | 改价原因 |
| adjusted_by | BIGINT | NO | 0 | 改价操作人 ID |
| adjusted_at | BIGINT | YES | NULL | 改价时间（Unix timestamp） |
| version | INT | NO | 1 | 乐观锁版本号 |
| payment_method | VARCHAR(32) | NO | '' | 支付方式: wechat, alipay |
| source | VARCHAR(32) | NO | '' | 订单来源: h5, mini_program, pc |

#### orders 表移除字段

| Column | Reason |
|--------|--------|
| tracking_no | 迁移到 shipments 表 |
| carrier | 迁移到 shipments 表 |

### API 扩展

在 `fulfillment.api` 中新增以下接口：

```go
// 更新商家备注
PUT /api/v1/orders/:id/remark
Request: { "remark": "string" }
Response: { "order_id": "string", "remark": "string", "updated_at": "string" }

// 改价
PUT /api/v1/orders/:id/adjust-price
Request: { "adjust_amount": -1000, "reason": "首单优惠补充" }
Response: { "order_id": "string", "original_amount": 10000, "adjust_amount": -1000, "new_pay_amount": 9000, "adjusted_at": "string" }

// 订单导出
GET /api/v1/orders/export?status=paid&start_time=...&end_time=...
Response: CSV 文件流

// 扩展履约摘要（新增今日统计）
GET /api/v1/orders/fulfillment-summary
Response: {
  // 现有字段
  "pending_shipment": 12,
  "partial_shipped": 2,
  "shipped": 50,
  "delivered": 30,
  "pending_refund": 3,
  // 新增字段
  "today_orders": 45,
  "today_gmv": 1289000
}
```

### 并发控制

改价操作使用乐观锁：

1. 读取订单时获取 version
2. 更新时检查 version 是否匹配
3. 不匹配则返回 `ErrOrderVersionConflict`
4. 前端提示用户刷新后重试

### 导出实现

同步导出，直接返回 CSV 文件流：

1. 应用当前筛选条件
2. 限制最大 10,000 条
3. 使用流式写入避免内存溢出
4. 设置响应头：
   - `Content-Type: text/csv; charset=utf-8`
   - `Content-Disposition: attachment; filename=orders_YYYYMMDD_HHMMSS.csv`

---

## Module Boundaries

### Fulfillment 模块职责（扩展后）

| Function | Responsibility |
|----------|---------------|
| 发货管理 | 创建、更新、查询发货单 |
| 退款管理 | 查询、批准、拒绝退款 |
| 物流公司管理 | 物流公司列表 |
| 订单履约视图 | 订单列表、订单详情（含履约信息） |
| **商家备注（新增）** | 更新订单内部备注 |
| **改价（新增）** | 调整待付款订单金额 |
| **今日统计（新增）** | 今日订单数、今日 GMV |
| **订单导出（新增）** | 导出订单 CSV |

### 不在范围内

| Function | 属于模块 |
|----------|----------|
| 订单创建 | Shop Service |
| 支付处理 | Payment Service |
| 促销计算 | Promotion Service |
| 库存扣减 | Inventory Service |

---

## Milestones and Sequencing

### Project Estimate

| Phase | Duration | Description |
|-------|----------|-------------|
| Phase 1: 数据库迁移 | 0.5 天 | orders 表字段变更 |
| Phase 2: 后端 API | 2 天 | 备注、改价、统计、导出接口 |
| Phase 3: 前端页面 | 1.5 天 | 改价对话框、备注对话框、统计卡片、导出按钮 |
| Phase 4: 联调测试 | 0.5 天 | E2E 测试、Bug 修复 |

**Total Estimate: 4.5 working days**

### Suggested Phases

#### Phase 1: 数据库迁移（Day 1 上午）

- 创建迁移脚本：orders 表新增字段
- 移除 tracking_no、carrier 字段（数据已迁移到 shipments 表）
- 回滚脚本

#### Phase 2: 后端 API（Day 1 下午 - Day 2）

- 扩展 fulfillment.api 定义
- 实现 UpdateRemarkHandler
- 实现 AdjustPriceHandler（含乐观锁）
- 扩展 GetFulfillmentSummaryHandler（新增今日统计）
- 实现 ExportOrdersHandler

#### Phase 3: 前端页面（Day 3 - Day 4 上午）

- 订单详情页：改价按钮、改价对话框
- 订单详情页：备注按钮、备注对话框
- 订单列表页：今日统计卡片
- 订单列表页：导出按钮

#### Phase 4: 联调测试（Day 4 下午）

- E2E 测试
- 并发改价测试
- 导出性能测试

---

## Database Schema

### orders 表变更

#### 新增字段

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| fulfillment_status | TINYINT | NO | 0 | 履约状态: 0-待发货, 1-部分发货, 2-已发货, 3-已送达 |
| refund_status | TINYINT | NO | 0 | 退款状态: 0-无, 1-待处理, 2-已批准, 3-已拒绝, 4-已完成 |
| merchant_remark | VARCHAR(500) | NO | '' | 商家内部备注 |
| original_amount | BIGINT | NO | 0 | 改价前原金额（分） |
| adjust_amount | BIGINT | NO | 0 | 改价金额（分），正数加价，负数减价 |
| adjust_reason | VARCHAR(200) | NO | '' | 改价原因 |
| adjusted_by | BIGINT | NO | 0 | 改价操作人 ID |
| adjusted_at | BIGINT | YES | NULL | 改价时间（Unix timestamp） |
| version | INT | NO | 1 | 乐观锁版本号 |
| payment_method | VARCHAR(32) | NO | '' | 支付方式: wechat, alipay |
| source | VARCHAR(32) | NO | '' | 订单来源: h5, mini_program, pc |

#### 移除字段

| Column | Reason |
|--------|--------|
| tracking_no | 已迁移到 shipments 表 |
| carrier | 已迁移到 shipments 表 |

#### 迁移脚本示例

```sql
-- 新增字段
ALTER TABLE `orders`
ADD COLUMN `fulfillment_status` TINYINT NOT NULL DEFAULT 0 COMMENT '履约状态' AFTER `status`,
ADD COLUMN `refund_status` TINYINT NOT NULL DEFAULT 0 COMMENT '退款状态' AFTER `fulfillment_status`,
ADD COLUMN `merchant_remark` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '商家备注' AFTER `remark`,
ADD COLUMN `original_amount` BIGINT NOT NULL DEFAULT 0 COMMENT '改价前原金额' AFTER `pay_amount`,
ADD COLUMN `adjust_amount` BIGINT NOT NULL DEFAULT 0 COMMENT '改价金额' AFTER `original_amount`,
ADD COLUMN `adjust_reason` VARCHAR(200) NOT NULL DEFAULT '' COMMENT '改价原因' AFTER `adjust_amount`,
ADD COLUMN `adjusted_by` BIGINT NOT NULL DEFAULT 0 COMMENT '改价操作人' AFTER `adjust_reason`,
ADD COLUMN `adjusted_at` BIGINT DEFAULT NULL COMMENT '改价时间' AFTER `adjusted_by`,
ADD COLUMN `version` INT NOT NULL DEFAULT 1 COMMENT '乐观锁版本' AFTER `adjusted_at`,
ADD COLUMN `payment_method` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '支付方式' AFTER `version`,
ADD COLUMN `source` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '订单来源' AFTER `payment_method`;

-- 初始化 original_amount 为现有 pay_amount
UPDATE `orders` SET `original_amount` = `pay_amount` WHERE `original_amount` = 0;

-- 移除已迁移字段（确认 shipments 表数据完整后执行）
-- ALTER TABLE `orders` DROP COLUMN `tracking_no`, DROP COLUMN `carrier`;

-- 添加索引
ALTER TABLE `orders` ADD INDEX `idx_fulfillment_status` (`fulfillment_status`);
ALTER TABLE `orders` ADD INDEX `idx_refund_status` (`refund_status`);
```

---

## API Endpoints

### 现有接口（Fulfillment v1.0）

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/orders | 订单列表（含履约状态筛选） |
| GET | /api/v1/orders/:id | 订单详情（含履约信息） |
| PUT | /api/v1/orders/:id/ship | 订单发货 |
| GET | /api/v1/orders/fulfillment-summary | 履约摘要统计 |
| GET | /api/v1/refunds | 退款列表 |
| GET | /api/v1/refunds/:id | 退款详情 |
| PUT | /api/v1/refunds/:id/approve | 批准退款 |
| PUT | /api/v1/refunds/:id/reject | 拒绝退款 |
| GET | /api/v1/carriers | 物流公司列表 |

### 新增接口

#### PUT /api/v1/orders/:id/remark - 更新商家备注

**Request:**
```json
{
  "remark": "客户要求周末配送，注意包装"
}
```

**Response:**
```json
{
  "order_id": "ORD20260322001",
  "remark": "客户要求周末配送，注意包装",
  "updated_at": "2026-03-22T15:00:00Z"
}
```

#### PUT /api/v1/orders/:id/adjust-price - 改价

**Request:**
```json
{
  "adjust_amount": -1000,
  "reason": "首单优惠补充"
}
```

**Response:**
```json
{
  "order_id": "ORD20260322001",
  "original_amount": 59800,
  "adjust_amount": -1000,
  "new_pay_amount": 58800,
  "adjusted_at": "2026-03-22T15:00:00Z"
}
```

**错误码:**
- `40006` - 当前状态无法改价（非 pending_payment）
- `40007` - 改价金额超出限制（超过 20%）
- `40010` - 改价原因不能为空
- `40011` - 订单已被修改，请刷新后重试（版本冲突）

#### GET /api/v1/orders/export - 导出订单

**Query Parameters:**
```
status: paid                      // optional
fulfillment_status: pending       // optional
refund_status: none               // optional
start_time: 2026-03-01T00:00:00Z  // optional
end_time: 2026-03-22T23:59:59Z    // optional
```

**Response:**
- Content-Type: `text/csv; charset=utf-8`
- Content-Disposition: `attachment; filename=orders_20260322_150000.csv`
- Body: CSV 文件内容

**CSV 列:**
```
订单号,订单状态,履约状态,退款状态,商品总额,优惠金额,运费,实付金额,收货人,收货电话,收货地址,支付方式,创建时间,支付时间
```

**错误码:**
- `40009` - 导出数量超出限制（超过 10,000 条）

#### GET /api/v1/orders/fulfillment-summary - 履约摘要（扩展）

**Response:**
```json
{
  "pending_shipment": 12,
  "partial_shipped": 2,
  "shipped": 50,
  "delivered": 30,
  "pending_refund": 3,
  "refunding": 1,
  "total_orders": 98,
  "today_orders": 45,
  "today_gmv": 1289000
}
```
  "adjust_amount": -1000,
  "reason": "首单优惠补充"
}
```

**Response:**
```json
{
  "order_id": "12345",
  "original_amount": 59800,
  "adjust_amount": -1000,
  "new_pay_amount": 49000,
  "adjusted_at": "2026-03-22T15:00:00Z"
}
```

#### GET /api/v1/orders/stats

**Response:**
```json
{
  "pending_payment": 5,
  "pending_shipment": 12,
  "pending_refund": 3,
  "today_orders": 45,
  "today_gmv": 1289000,
  "currency": "CNY"
}
```

---

## Business Rules

### 改价规则

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-001 | 状态限制 | 仅 `pending_payment` 状态可改价 |
| BR-002 | 金额限制 | `|adjust_amount| ≤ original_amount × 20%` |
| BR-003 | 原因必填 | 改价原因必填，用于审计 |
| BR-004 | 原金额保留 | `original_amount` 存储改价前金额，不变 |
| BR-005 | 审计记录 | 记录改价操作人、时间、原因 |

### 备注规则

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-006 | 长度限制 | 最大 500 字符 |
| BR-007 | 可见性 | 仅商家后台可见，买家端不展示 |
| BR-008 | 可覆盖 | 多次更新备注，不保留历史版本 |

### 导出规则

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-009 | 数量限制 | 最大 10,000 条记录 |
| BR-010 | 筛选条件 | 导出时应用当前列表筛选条件 |
| BR-011 | 同步下载 | 直接返回文件流，不存储 |

### 并发控制

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-012 | 乐观锁 | 改价时检查 `version` 字段 |
| BR-013 | 版本递增 | 每次更新 `version` 递增 |
| BR-014 | 冲突处理 | 版本不匹配返回 `40011` 错误 |

---

## Error Codes

新增错误码（Order 模块使用 40xxx 范围）：

| Constant | HTTP Status | Code | Message |
|----------|-------------|------|---------|
| ErrOrderCannotAdjustPrice | 400 | 40006 | 当前状态无法改价 |
| ErrOrderAdjustAmountExceed | 400 | 40007 | 改价金额超出限制 |
| ErrOrderExportLimitExceed | 400 | 40009 | 导出数量超出限制 |
| ErrOrderAdjustReasonRequired | 400 | 40010 | 改价原因不能为空 |
| ErrOrderVersionConflict | 409 | 40011 | 订单已被修改，请刷新后重试 |

---

## User Stories

### US-001: 改价

**Description**: 作为商家管理员，我想调整待付款订单的价格，以便处理特殊场景。

**Acceptance Criteria**:
- Given 订单状态为 `pending_payment`
- When 我点击「改价」按钮
- And 输入调整金额 -500 分
- And 输入改价原因「首单优惠补充」
- And 点击「确认」
- Then 订单实付金额更新
- And 记录改价人、时间、原因

---

### US-002: 改价校验

**Description**: 作为系统，我要校验改价请求，防止违规操作。

**Acceptance Criteria**:
- Given 订单状态为 `paid`
- When 我尝试改价
- Then 返回错误「当前状态无法改价」
- Given 订单状态为 `pending_payment`，原金额 10000 分
- When 我尝试改价 -3000 分（30%）
- Then 返回错误「改价金额超出限制」

---

### US-003: 添加备注

**Description**: 作为商家管理员，我想为订单添加内部备注，方便团队协作。

**Acceptance Criteria**:
- Given 我正在查看订单详情
- When 我点击「备注」按钮
- And 输入「客户要求周末配送」
- And 点击「保存」
- Then 备注保存成功
- And 备注仅商家可见

---

### US-004: 查看今日统计

**Description**: 作为商家管理员，我想看到今日订单数和 GMV，了解业务状况。

**Acceptance Criteria**:
- Given 我打开订单列表页
- Then 我看到今日订单数和今日 GMV
- And 数据实时更新

---

### US-005: 导出订单

**Description**: 作为商家管理员，我想导出订单数据，以便离线分析。

**Acceptance Criteria**:
- Given 我设置了筛选条件
- When 我点击「导出」按钮
- Then 浏览器下载 CSV 文件
- And 文件包含符合筛选条件的订单

---

### US-006: 导出限制

**Description**: 作为系统，我要限制导出数量，保护系统性能。

**Acceptance Criteria**:
- Given 筛选结果超过 10,000 条
- When 我点击「导出」
- Then 返回错误「导出数量超出限制，请缩小筛选范围」

---

## Frontend Structure

### 新增组件

```
src/views/orders/
├── index.vue                    # 订单列表页（现有）
├── detail.vue                   # 订单详情页（现有）
└── components/
    ├── AdjustPriceDialog.vue    # 改价对话框（新增）
    ├── RemarkDialog.vue         # 备注对话框（新增）
    └── TodayStats.vue           # 今日统计卡片（新增）
```

### 组件规格

#### AdjustPriceDialog.vue

- 显示原订单金额
- 输入调整金额（数字，可负数）
- 实时预览新金额
- 输入改价原因（必填）
- 显示 20% 限制提示
- 校验失败显示错误信息

#### RemarkDialog.vue

- 显示当前备注内容（如有）
- 输入/编辑备注（最大 500 字符）
- 保存/取消按钮

#### TodayStats.vue

- 显示今日订单数
- 显示今日 GMV
- 与现有统计卡片样式一致

---

## Appendix

### 订单状态定义

| Status | 中文 | Description |
|--------|------|-------------|
| pending_payment | 待付款 | 订单创建，等待支付 |
| paid | 已支付 | 支付完成，等待发货 |
| shipped | 已发货 | 全部商品已发货 |
| completed | 已完成 | 订单完成 |
| canceled | 已取消 | 订单取消 |
| refunding | 退款中 | 退款申请处理中 |
| refunded | 已退款 | 退款完成 |

### 履约状态定义

| Status | 中文 | Description |
|--------|------|-------------|
| pending | 待发货 | 无商品发货 |
| partial_shipped | 部分发货 | 部分商品发货 |
| shipped | 已发货 | 全部商品发货 |
| delivered | 已送达 | 全部发货单已送达 |

### 退款状态定义

| Status | 中文 | Description |
|--------|------|-------------|
| none | 无 | 无退款活动 |
| pending | 待处理 | 退款申请待审核 |
| approved | 已同意 | 退款已批准 |
| rejected | 已拒绝 | 退款已拒绝 |
| completed | 已完成 | 退款已完成 |