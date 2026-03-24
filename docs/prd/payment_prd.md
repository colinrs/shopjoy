# Payment Module Product Requirements Document

## Document Information

| Item | Value |
|------|-------|
| Document Title | Payment Module PRD (Admin) |
| Version | 1.1.0 |
| Status | Draft |
| Created | 2026-03-23 |
| Updated | 2026-03-23 |
| Author | Product Team |
| Module Location | admin/internal/domain/payment/ |

---

## Product Overview

### Summary

本 PRD 描述 ShopJoy Admin 支付集成模块的功能需求。支付模块为商家后台提供支付信息查看、退款管理、支付统计等核心能力，帮助商家全面掌握订单支付状态和资金流水。

**核心功能**：
- 订单支付信息展示（支付渠道、流水号、金额、手续费、状态）
- 退款管理（发起退款、退款列表、退款详情）
- 支付统计看板（实收金额、渠道分布、退款金额、退款率）
- 支付流水查询（按订单号/流水号/时间查询）

**Phase 1 范围**：
- 只做 admin 端查看 + 退款发起能力
- 支付发起（买家支付页/二维码/H5）放在 shop/joy 模块，后续迭代
- 先支持 Stripe 渠道，支付宝/微信可二期扩展
- Stripe 渠道支持的货币：USD, EUR, GBP, JPY, SGD, HKD（CNY 等二期 Alipay/WeChat 接入）

### Problem Statement

商家在日常运营中，对于订单支付信息的管理存在以下痛点：

- **支付信息不透明**：无法在订单详情中清晰看到支付渠道、流水号、手续费等信息
- **退款操作分散**：退款申请和退款记录查看分散在不同页面，操作不便
- **缺少支付统计**：无法快速查看实收金额、渠道分布、退款率等关键指标
- **流水查询困难**：支付流水与订单关联不清晰，排查问题效率低

### Solution Overview

在 admin 服务中新增支付模块，提供以下能力：

1. **支付信息展示** - 在订单详情页集成支付信息区块
2. **退款管理** - 支持商家主动发起退款，查看退款记录和状态
3. **支付统计** - Dashboard 展示关键支付指标
4. **流水查询** - 支持多维度查询支付流水

---

## Goals

### Business Goals

| ID | Goal | Metric | Target |
|----|------|--------|--------|
| BG-001 | 提升退款处理效率 | 退款 API 响应时间 P95 | < 500ms |
| BG-002 | 降低支付异常率 | 系统错误率 | < 0.1% |
| BG-003 | 提高财务对账效率 | 对账时间节省 | 50% |

### User Goals

| ID | Goal | User Type |
|----|------|-----------|
| UG-001 | 查看订单支付详情 | Merchant Admin / Customer Service |
| UG-002 | 发起订单退款 | Merchant Admin |
| UG-003 | 查看支付统计 | Merchant Admin / Store Manager |
| UG-004 | 查询支付流水 | Merchant Admin / Finance |

### Non-Goals

| ID | Non-Goal | Reason |
|----|----------|--------|
| NG-001 | 买家端支付页面 | 属于 shop service 范围 |
| NG-002 | 支付渠道配置 | Phase 2 feature |
| NG-003 | 批量退款 | Phase 2 feature |
| NG-004 | 自动对账 | Phase 2 feature |
| NG-005 | 支付宝/微信渠道 | Phase 2，先支持 Stripe |
| NG-006 | CNY 货币支持 | Stripe 不支持 CNY，等 Phase 2 Alipay/WeChat |
| NG-007 | 跨货币退款 | 不支持，退款货币必须与原支付一致 |

---

## User Personas

### Key User Types

| Persona | Description | Primary Use Cases |
|---------|-------------|-------------------|
| Merchant Admin | 商家管理员 | 查看支付信息、发起退款、查看统计 |
| Customer Service | 客服人员 | 查看支付信息、处理退款申请 |
| Finance | 财务人员 | 查看支付流水、对账 |

### Role-Based Access

| Role | 查看支付信息 | 发起退款 | 查看统计 | 导出流水 |
|------|------------|---------|---------|---------|
| Tenant Admin | ✅ | ✅ | ✅ | ✅ |
| Tenant Operations Manager | ✅ | ✅ | ✅ | ✅ |
| Tenant Customer Service | ✅ | ❌ | ✅ | ❌ |

---

## Functional Requirements

### Feature 1: 订单支付信息展示

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-001 | 支付信息区块 | 在订单详情页展示支付信息 | P0 |
| FR-002 | 支付渠道展示 | 显示支付渠道（Stripe/Alipay/WeChat） | P0 |
| FR-003 | 支付流水号 | 显示渠道支付流水号（Charge ID），支持复制 | P0 |
| FR-004 | PaymentIntent ID | 显示渠道 PaymentIntent ID（Stripe） | P0 |
| FR-005 | 支付金额 | 显示实际支付金额 | P0 |
| FR-006 | 手续费展示 | 显示渠道手续费 | P0 |
| FR-007 | 支付时间 | 显示支付成功时间 | P0 |
| FR-008 | 支付状态 | 显示支付状态（Tag 形式） | P0 |
| FR-009 | 退款记录 | 显示该订单的退款记录列表 | P0 |

### Feature 2: 退款管理

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-010 | 发起退款入口 | 在订单详情页提供退款按钮 | P0 |
| FR-011 | 全额退款 | 支持全额退款 | P0 |
| FR-012 | 部分退款 | 支持部分退款（输入金额） | P0 |
| FR-013 | 退款原因必填 | 选择退款原因 | P0 |
| FR-014 | 退款金额校验 | 退款金额 ≤ 已支付金额 - 已退款金额 | P0 |
| FR-015 | 幂等键支持 | 前端传入 idempotency_key 防止重复退款 | P0 |
| FR-016 | 退款列表页 | 展示所有退款记录 | P0 |
| FR-017 | 退款状态筛选 | 按退款状态筛选 | P0 |
| FR-018 | 退款详情 | 查看退款详情（含渠道退款流水号） | P0 |
| FR-019 | 渠道退款状态 | 显示渠道侧退款状态 | P0 |
| FR-020 | 货币一致性 | 退款货币自动使用原支付货币，不支持跨货币退款 | P0 |

### Feature 3: 支付统计

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-030 | 今日实收 | 今日已支付订单总金额 | P0 |
| FR-031 | 近7天实收 | 近7天已支付订单总金额 | P1 |
| FR-032 | 渠道分布 | 各支付渠道金额占比 | P1 |
| FR-033 | 退款金额 | 今日/近7天退款金额 | P1 |
| FR-034 | 退款率 | 退款金额 / 实收金额 | P1 |

### Feature 4: 支付流水查询

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-040 | 流水列表页 | 支付流水列表 | P1 |
| FR-041 | 多条件筛选 | 按订单号、流水号、渠道、状态、时间筛选 | P1 |
| FR-042 | 导出功能 | 导出支付流水 CSV（最大 10,000 条） | P2 |

### Feature 5: Webhook 处理

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-050 | Webhook 端点 | POST /api/v1/webhooks/stripe | P0 |
| FR-051 | 签名验证 | 验证 Stripe webhook 签名 | P0 |
| FR-052 | 支付成功事件 | 处理 payment_intent.succeeded | P0 |
| FR-053 | 支付失败事件 | 处理 payment_intent.payment_failed | P0 |
| FR-054 | 支付取消事件 | 处理 payment_intent.canceled | P0 |
| FR-055 | 3D Secure 事件 | 处理 payment_intent.requires_action | P0 |
| FR-056 | 退款成功事件 | 处理 charge.refunded | P0 |
| FR-057 | 退款更新事件 | 处理 charge.refund.updated | P0 |
| FR-058 | 幂等处理 | 同一事件多次发送只处理一次（使用 webhook_events 表） | P0 |
| FR-059 | 争议创建事件 | 处理 charge.dispute.created | P1 |
| FR-060 | 争议更新事件 | 处理 charge.dispute.updated | P1 |
| FR-061 | 争议关闭事件 | 处理 charge.dispute.closed | P1 |
| FR-062 | 时间戳验证 | 验证 webhook 时间戳在 5 分钟内，防止重放攻击 | P0 |

---

## User Experience

### 支付信息查看

1. 商家在订单详情页看到「支付信息」区块
2. 展示支付渠道、流水号（可复制）、金额、手续费、支付时间
3. 支付状态以 Tag 形式展示（成功/失败/处理中/待验证）
4. 如有退款记录，展示退款历史表格

### 发起退款

1. 商家在订单详情页点击「退款」按钮
2. 弹出退款对话框：
   - 前端生成幂等键 `idempotency_key`
   - 选择全额退款 / 部分退款
   - 部分退款时输入金额（显示最大可退金额）
   - 选择退款原因
3. 点击「确认」
4. 系统校验：
   - 订单状态为已支付
   - 退款金额不超过可退金额
   - 检查幂等键是否已存在
5. 调用 Stripe 退款接口（使用 Charge ID）
6. 创建退款记录，更新订单状态
7. 显示成功提示

### 查看支付统计

1. 商家打开 Dashboard 页面
2. 看到支付统计卡片：
   - 今日实收金额（对比昨日）
   - 近7天实收金额
   - 渠道分布图
   - 退款金额和退款率

### 查询支付流水

1. 商家进入「支付流水」页面
2. 设置筛选条件（订单号/流水号/渠道/状态/时间）
3. 查看流水列表
4. 可导出 CSV 文件（最大 10,000 条）

---

## Narrative

商家小明早上打开后台，首先查看 Dashboard 的支付统计卡片，看到今日实收 USD 1,890.00，相比昨日增长 15%。渠道分布显示 Stripe 占比 100%。

一位客户来电询问订单支付状态，小明在订单详情页查看到该订单使用 Stripe 支付，PaymentIntent ID 为 `pi_3Oabc...`，Charge ID 为 `ch_3Oabc...`，支付时间 2026-03-23T10:30:00Z，手续费 2.9%。

客户因商品问题申请退款，小明点击「退款」按钮，选择全额退款，原因选择「DEFECTIVE（商品质量问题）」，确认后系统自动调用 Stripe 退款接口，退款预计 5-10 个工作日到账。

下午，小明在退款列表页查看今日处理的退款，状态显示「succeeded」，渠道退款流水号 `re_xxx` 也已同步显示。

下班前，小明导出今日支付流水，准备财务对账。

---

## Success Metrics

### User-Centric Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| 支付信息展示完整率 | 订单详情支付信息完整展示 | 100% |
| 退款发起成功率 | 退款 API 调用成功率 | > 99% |
| 统计卡片使用率 | Dashboard 支付统计查看率 | > 60% |

### Technical Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| 退款 API 响应时间 | P95 | < 500ms |
| 支付流水查询响应时间 | P95 | < 200ms |
| Webhook 处理延迟 | 从 Stripe 发送到处理完成 | < 3s |
| 系统错误率 | 5xx 错误比例 | < 0.1% |

---

## Technical Considerations

### 现有架构

支付模块采用 DDD 架构：

```
admin/internal/domain/payment/
├── entity.go          # 领域实体：Payment, PaymentTransaction, PaymentRefund
├── repository.go      # 仓储接口
└── service.go         # 领域服务
```

API 定义：`admin/desc/payment.api`

### 数据库设计

#### order_payments 表（新增）

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | 主键 |
| tenant_id | BIGINT | NO | | 租户 ID |
| order_id | BIGINT | NO | | 订单 ID |
| payment_no | VARCHAR(32) | NO | | 支付单号 |
| payment_method | VARCHAR(32) | NO | | 支付方式: stripe, alipay, wechat |
| channel_intent_id | VARCHAR(64) | NO | '' | 渠道 PaymentIntent ID (Stripe: pi_xxx) |
| channel_payment_id | VARCHAR(64) | NO | '' | 渠道 Charge ID (Stripe: ch_xxx) |
| amount | BIGINT | NO | 0 | 支付金额（分） |
| currency | VARCHAR(3) | NO | 'USD' | 货币（Stripe: USD/EUR/GBP/JPY/SGD/HKD） |
| status | TINYINT | NO | 0 | 状态（见支付状态定义） |
| transaction_fee | BIGINT | NO | 0 | 手续费（分） |
| fee_currency | VARCHAR(3) | NO | 'USD' | 手续费货币 |
| paid_at | TIMESTAMP | YES | NULL | 支付成功时间（UTC） |
| failed_at | TIMESTAMP | YES | NULL | 失败时间（UTC） |
| failed_reason | VARCHAR(255) | NO | '' | 失败原因 |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 创建时间（UTC） |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 更新时间（UTC） |
| deleted_at | TIMESTAMP | YES | NULL | 删除时间（UTC） |

Indexes:
- `idx_tenant_order` (tenant_id, order_id)
- `idx_payment_no` (payment_no)
- `idx_channel_payment_id` (channel_payment_id)
- `idx_channel_intent_id` (channel_intent_id)
- `idx_status` (status)

#### payment_transactions 表（新增）

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | 主键 |
| tenant_id | BIGINT | NO | | 租户 ID |
| order_id | BIGINT | NO | | 订单 ID |
| payment_id | BIGINT | NO | | 支付单 ID |
| transaction_id | VARCHAR(64) | NO | | 交易流水号 |
| payment_method | VARCHAR(32) | NO | | 支付方式 |
| channel_transaction_id | VARCHAR(64) | NO | '' | 渠道交易流水号 |
| amount | BIGINT | NO | 0 | 交易金额（分） |
| currency | VARCHAR(3) | NO | 'USD' | 货币 |
| status | TINYINT | NO | 0 | 状态: 0=pending, 1=succeeded, 2=failed |
| transaction_fee | BIGINT | NO | 0 | 手续费（分） |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 创建时间（UTC） |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 更新时间（UTC） |
| paid_at | TIMESTAMP | YES | NULL | 支付成功时间（UTC） |
| failed_reason | VARCHAR(255) | NO | '' | 失败原因 |
| deleted_at | TIMESTAMP | YES | NULL | 删除时间（UTC） |

Indexes:
- `idx_tenant_order` (tenant_id, order_id)
- `idx_transaction_id` (transaction_id)
- `idx_channel_transaction_id` (channel_transaction_id)

#### payment_refunds 表（新增，独立于 fulfillment.refunds）

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | 主键 |
| tenant_id | BIGINT | NO | | 租户 ID |
| order_id | BIGINT | NO | | 订单 ID |
| payment_id | BIGINT | NO | | 支付单 ID |
| fulfillment_refund_id | BIGINT | YES | NULL | 关联的履约退款 ID（如有） |
| refund_no | VARCHAR(32) | NO | | 退款单号 |
| idempotency_key | VARCHAR(64) | NO | '' | 幂等键 |
| channel_refund_id | VARCHAR(64) | NO | '' | 渠道退款流水号（Stripe: re_xxx） |
| amount | BIGINT | NO | 0 | 退款金额（分） |
| currency | VARCHAR(3) | NO | 'USD' | 货币 |
| refund_fee | BIGINT | NO | 0 | 退款手续费（分） |
| status | TINYINT | NO | 0 | 渠道退款状态 |
| reason_type | VARCHAR(32) | NO | '' | 退款原因类型 |
| reason | VARCHAR(255) | NO | '' | 退款原因详情 |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 创建时间（UTC） |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 更新时间（UTC） |
| refunded_at | TIMESTAMP | YES | NULL | 退款完成时间（UTC） |
| created_by | BIGINT | NO | 0 | 创建人 ID |

Indexes:
- `idx_tenant_order` (tenant_id, order_id)
- `idx_refund_no` (refund_no)
- `idx_idempotency_key` (idempotency_key)
- `idx_channel_refund_id` (channel_refund_id)
- `idx_payment_id` (payment_id)

#### webhook_events 表（新增，用于事件去重）

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | 主键 |
| tenant_id | BIGINT | NO | | 租户 ID |
| event_id | VARCHAR(64) | NO | | 渠道事件 ID（Stripe: evt_xxx） |
| event_type | VARCHAR(64) | NO | | 事件类型 |
| resource_id | VARCHAR(64) | NO | '' | 关联资源 ID（payment_intent_id / charge_id 等） |
| processed | TINYINT | NO | 0 | 是否已处理: 0=pending, 1=processed, 2=failed |
| raw_payload | TEXT | YES | NULL | 原始事件 JSON |
| error_message | VARCHAR(255) | NO | '' | 处理错误信息 |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 创建时间（UTC） |
| processed_at | TIMESTAMP | YES | NULL | 处理时间（UTC） |

Indexes:
- `idx_event_id` UNIQUE (event_id)
- `idx_tenant_event` (tenant_id, event_type)
- `idx_resource` (resource_id)

### API Endpoints

#### Payment API

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/payments/stats | 支付统计 |
| GET | /api/v1/payments/transactions | 支付流水列表 |
| GET | /api/v1/payments/transactions/:id | 支付流水详情 |
| GET | /api/v1/orders/:id/payment | 订单支付信息 |
| POST | /api/v1/orders/:id/refund | 发起退款 |
| POST | /api/v1/webhooks/stripe | Stripe Webhook 回调 |

#### GET /api/v1/payments/stats - 支付统计

**Query Parameters:**
```
period: 7d, 30d, 90d  // 统计周期
```

**Response:**
```json
{
  "today_received": 189000,
  "today_received_change": "+15%",
  "period_received": 1250000,
  "refund_amount": 12500,
  "refund_rate": "1.0%",
  "channel_distribution": [
    {"channel": "stripe", "amount": 1250000, "percentage": "100%"}
  ],
  "currency": "USD"
}
```

#### GET /api/v1/payments/transactions - 支付流水列表

**Query Parameters:**
```
page: 1
page_size: 20
order_no: ORD001             // optional
transaction_id: TXN001       // optional
payment_method: stripe       // optional
status: 1                    // optional: 0=pending, 1=succeeded, 2=failed
start_time: 2026-03-01T00:00:00Z
end_time: 2026-03-23T23:59:59Z
```

**Response:**
```json
{
  "list": [
    {
      "id": 1,
      "transaction_id": "TXN20260323001",
      "order_id": "12345",
      "order_no": "ORD20260323001",
      "payment_method": "stripe",
      "channel_transaction_id": "ch_3Oabc...",
      "amount": 29900,
      "currency": "USD",
      "transaction_fee": 867,
      "status": 1,
      "status_text": "Succeeded",
      "created_at": "2026-03-23T10:30:00Z",
      "paid_at": "2026-03-23T10:30:15Z"
    }
  ],
  "total": 100,
  "page": 1,
  "page_size": 20
}
```

#### GET /api/v1/orders/:id/payment - 订单支付信息

**Response:**
```json
{
  "payment_id": 1,
  "payment_no": "PAY20260323001",
  "payment_method": "stripe",
  "payment_method_text": "Stripe",
  "channel_intent_id": "pi_3Oabc...",
  "channel_payment_id": "ch_3Oabc...",
  "amount": 29900,
  "currency": "USD",
  "transaction_fee": 867,
  "fee_currency": "USD",
  "status": 2,
  "status_text": "Success",
  "paid_at": "2026-03-23T10:30:15Z",
  "refunded_amount": 0,
  "refunds": []
}
```

#### POST /api/v1/orders/:id/refund - 发起退款

**Request:**
```json
{
  "idempotency_key": "uuid-v4-generated-by-client",
  "amount": 29900,
  "reason_type": "DEFECTIVE",
  "reason": "商品质量问题"
}
```

**Response:**
```json
{
  "refund_id": 1,
  "refund_no": "REF20260323001",
  "amount": 29900,
  "currency": "USD",
  "status": 0,
  "status_text": "Pending",
  "channel_refund_id": "re_xxx"
}
```

#### POST /api/v1/webhooks/stripe - Stripe Webhook

**Headers:**
```
Stripe-Signature: t=1234567890,v1=xxx,v0=xxx
```

**Request Body:** Stripe Event JSON

**处理的事件类型：**

| Event | 处理逻辑 |
|-------|---------|
| `payment_intent.succeeded` | 更新 payment 状态为 success，记录 channel_payment_id |
| `payment_intent.payment_failed` | 更新 payment 状态为 failed，记录失败原因 |
| `payment_intent.canceled` | 更新 payment 状态为 cancelled |
| `payment_intent.requires_action` | 更新 payment 状态为 requires_action |
| `charge.refunded` | 更新 refund 状态为 succeeded |
| `charge.refund.updated` | 同步退款状态 |

---

## Business Rules

### 支付状态规则

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-001 | 状态流转 | pending → processing → success/failed/cancelled |
| BR-002 | 退款状态 | success → refunded/partially_refunded |
| BR-003 | 状态同步 | 通过 webhook 同步 Stripe 状态 |
| BR-004 | 3D Secure | requires_action 状态需要额外验证 |

### 退款规则

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-010 | 金额校验 | 退款金额 ≤ 已支付金额 - 已退款金额 |
| BR-011 | 原因必填 | 退款原因必填 |
| BR-012 | 状态限制 | 仅 success 状态可发起退款 |
| BR-013 | 退款记录不可修改 | 退款一旦发起，只能查看 |
| BR-014 | 渠道同步 | 调用 Stripe 退款接口（使用 Charge ID），同步退款状态 |
| BR-015 | 幂等处理 | 使用 idempotency_key 防止重复退款 |
| BR-016 | 独立实体 | PaymentRefund 独立于 fulfillment.Refund |

### 金额规则

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-020 | 单位统一 | 所有金额使用「分」为单位 |
| BR-021 | 货币绑定 | 金额需绑定货币（ISO 4217） |
| BR-022 | 精确计算 | 使用 decimal 类型计算，避免浮点误差 |
| BR-023 | 货币一致 | 退款货币必须与原支付货币一致 |
| BR-024 | 不支持跨货币 | 不支持跨货币退款 |

### 货币规则

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-030 | Stripe 货币 | Phase 1 仅支持 USD, EUR, GBP, JPY, SGD, HKD |
| BR-031 | CNY 支持 | Phase 2 Alipay/WeChat 接入后支持 |

---

## Security Requirements

### PCI DSS Compliance

| 要求 | 说明 |
|-----|------|
| Stripe Elements | 使用 Stripe Elements 收集卡信息，符合 SAQ A 标准 |
| 禁止存储 | 不存储完整卡号、CVV、有效期 |

### API Key Management

| 要求 | 说明 |
|-----|------|
| 加密存储 | Stripe Secret Key 存储在加密配置中 |
| 最小权限 | 使用 Stripe Restricted Keys，仅授予必要权限 |
| 定期轮换 | 定期轮换 API Key |

### Webhook Security

| 要求 | 说明 |
|-----|------|
| 签名验证 | 验证 Stripe webhook 签名（使用 signing secret） |
| 时间戳验证 | 验证 webhook 时间戳在 5 分钟内，防止重放攻击 |
| 事件去重 | 使用 webhook_events 表记录已处理事件，防止重复处理 |
| 拒绝无效 | 签名无效或时间戳过期的请求返回 400 |

### Data Protection

| 要求 | 说明 |
|-----|------|
| 租户隔离 | 支付记录按租户隔离 |
| 操作审计 | 敏感操作记录审计日志 |

---

## Error Codes

Payment 模块使用 50xxx 错误码范围：

| Constant | HTTP Status | Code | Message |
|----------|-------------|------|---------|
| ErrPaymentNotFound | 404 | 50001 | payment not found |
| ErrPaymentInvalidAmount | 400 | 50002 | invalid payment amount |
| ErrPaymentFailed | 402 | 50003 | payment failed |
| ErrPaymentAlreadyPaid | 400 | 50004 | payment already completed |
| ErrPaymentExpired | 400 | 50005 | payment expired |
| ErrPaymentOrderNotPaid | 400 | 50006 | order not paid, cannot refund |
| ErrRefundAmountExceeded | 400 | 50007 | refund amount exceeds refundable |
| ErrPaymentRefundReasonRequired | 400 | 50008 | refund reason is required |
| ErrChannelRefundFailed | 500 | 50009 | channel refund failed |
| ErrPaymentOrderNotFound | 404 | 50010 | order not found |
| ErrPaymentOrderAlreadyRefunded | 400 | 50011 | order already fully refunded |
| ErrTransactionNotFound | 404 | 50012 | transaction not found |
| ErrPaymentChannelUnavailable | 503 | 50013 | payment channel unavailable |
| ErrRefundNotSupported | 400 | 50014 | refund not supported for this channel |
| ErrCurrencyNotSupported | 400 | 50015 | currency not supported by channel |
| ErrPaymentRefundNotFound | 404 | 50016 | payment refund not found |
| ErrPaymentRequiresAction | 202 | 50017 | payment requires additional action (3D Secure) |
| ErrRefundCurrencyMismatch | 400 | 50018 | refund currency must match payment currency |
| ErrIdempotencyKeyConflict | 409 | 50019 | duplicate idempotency key |
| ErrDisputeCreated | 409 | 50020 | dispute created for this charge |

---

## User Stories

### US-001: 查看订单支付信息

**Description**: 作为商家管理员，我想查看订单的支付详情，了解支付状态和金额。

**Acceptance Criteria**:
- Given 我正在查看订单详情
- When 订单已支付
- Then 我看到支付渠道、PaymentIntent ID、Charge ID、金额、手续费、支付时间
- And 流水号可复制

---

### US-002: 发起全额退款

**Description**: 作为商家管理员，我想为订单发起全额退款。

**Acceptance Criteria**:
- Given 订单状态为支付成功
- When 我点击「退款」按钮
- And 选择「全额退款」
- And 选择退款原因
- And 点击「确认」
- Then 系统使用 Charge ID 调用 Stripe 退款接口
- And 创建退款记录
- And 更新订单状态为已退款

---

### US-003: 发起部分退款

**Description**: 作为商家管理员，我想为订单发起部分退款。

**Acceptance Criteria**:
- Given 订单状态为支付成功
- When 我点击「退款」按钮
- And 选择「部分退款」
- And 输入退款金额
- And 选择退款原因
- And 点击「确认」
- Then 系统校验退款金额不超过可退金额
- And 使用 Charge ID 调用 Stripe 退款接口
- And 创建退款记录
- And 更新订单状态为部分退款

---

### US-004: 查看支付统计

**Description**: 作为商家管理员，我想查看支付统计数据。

**Acceptance Criteria**:
- Given 我打开 Dashboard
- Then 我看到今日实收、近7天实收、退款金额、退款率
- And 我看到渠道分布图

---

### US-005: 查询支付流水

**Description**: 作为商家管理员，我想查询支付流水记录。

**Acceptance Criteria**:
- Given 我进入支付流水页面
- When 我输入订单号或流水号
- Then 我看到匹配的支付流水列表

---

### US-006: 防止重复退款

**Description**: 作为系统，我要防止重复退款。

**Acceptance Criteria**:
- Given 用户发起退款请求
- When 同一 idempotency_key 已存在
- Then 返回已存在的退款记录
- And 不创建新退款

---

## Frontend Structure

### 新增页面

```
src/views/
├── payments/
│   ├── index.vue                    # 支付流水列表页
│   └── components/
│       ├── PaymentStatsCard.vue     # 支付统计卡片
│       └── RefundDialog.vue         # 退款对话框
└── orders/
    └── [id]/
        └── components/
            └── PaymentInfo.vue      # 订单支付信息区块
```

### 组件规格

#### PaymentStatsCard.vue

- 今日实收金额（对比昨日涨跌）
- 近7天实收金额
- 渠道分布（进度条）
- 退款金额
- 退款率

#### PaymentInfo.vue

- 支付渠道（Tag）
- PaymentIntent ID（可复制）
- Charge ID（可复制）
- 实付金额
- 手续费
- 支付时间
- 支付状态（Tag）
- 退款记录表格

#### RefundDialog.vue

- 幂等键（前端生成 UUID）
- 全额/部分退款切换
- 退款金额输入（部分退款时）
- 最大可退金额显示
- 退款原因选择
- 确认/取消按钮

---

## Milestones and Sequencing

### Project Estimate

| Phase | Duration | Description |
|-------|----------|-------------|
| Phase 1: 数据库设计 | 0.5 天 | payment、transaction、payment_refund 表 |
| Phase 2: 后端 API | 2 天 | 支付信息、退款、统计、流水、webhook 接口 |
| Phase 3: 前端页面 | 2 天 | 支付信息区块、退款对话框、统计卡片、流水列表 |
| Phase 4: Stripe 集成 | 1 天 | Stripe 退款接口调用、webhook 处理、签名验证 |
| Phase 5: 联调测试 | 0.5 天 | E2E 测试、Bug 修复 |

**Total Estimate: 6 working days**

---

## Appendix

### 支付状态定义

| Status | 值 | 中文 | English | Description |
|--------|---|------|---------|-------------|
| pending | 0 | 待支付 | Pending | 订单创建，等待支付 |
| processing | 1 | 支付中 | Processing | 支付发起，处理中 |
| success | 2 | 支付成功 | Success | 支付完成 |
| failed | 3 | 支付失败 | Failed | 支付失败 |
| cancelled | 4 | 已取消 | Cancelled | 订单取消 |
| refunded | 5 | 已退款 | Refunded | 全额退款完成 |
| partially_refunded | 6 | 部分退款 | Partially Refunded | 部分金额已退回 |
| requires_action | 7 | 待验证 | Requires Action | 3D Secure 验证待处理 |

### 退款状态定义（PaymentRefund）

| Status | 值 | 中文 | English | Description |
|--------|---|------|---------|-------------|
| pending | 0 | 待处理 | Pending | 退款已创建，等待渠道处理 |
| succeeded | 1 | 退款成功 | Succeeded | 渠道退款完成 |
| failed | 2 | 退款失败 | Failed | 渠道退款失败 |

### 支付渠道定义

| Code | Name | Supported Currencies |
|------|------|---------------------|
| stripe | Stripe | USD, EUR, GBP, JPY, SGD, HKD |
| alipay | 支付宝 | CNY (Phase 2) |
| wechat | 微信支付 | CNY (Phase 2) |

### 退款原因类型

| Code | 中文 | English |
|------|------|---------|
| DEFECTIVE | 商品质量问题 | Product Defective |
| WRONG_ITEM | 发错商品 | Wrong Item Received |
| NOT_AS_DESCRIBED | 商品与描述不符 | Not As Described |
| DAMAGED | 商品损坏 | Damaged in Transit |
| NO_LONGER_NEEDED | 不再需要 | No Longer Needed |
| LATE_DELIVERY | 配送延迟 | Late Delivery |
| OTHER | 其他原因 | Other |

### Stripe 事件映射

| Stripe Event | 系统动作 |
|--------------|---------|
| `payment_intent.succeeded` | payment.status = success, 记录 channel_payment_id |
| `payment_intent.payment_failed` | payment.status = failed |
| `payment_intent.canceled` | payment.status = cancelled |
| `payment_intent.requires_action` | payment.status = requires_action |
| `charge.refunded` | payment_refund.status = succeeded |
| `charge.refund.updated` | 同步 payment_refund.status |
| `charge.dispute.created` | 创建争议记录，通知商家 |
| `charge.dispute.updated` | 更新争议状态 |
| `charge.dispute.closed` | 关闭争议，记录结果（won/lost） |

### 零小数货币处理

| 货币 | 处理方式 |
|-----|---------|
| JPY | 金额单位为日元整数，无需除以 100 转换 |
| USD/EUR/GBP/SGD/HKD | 金额单位为分，Stripe 传参时保持不变 |