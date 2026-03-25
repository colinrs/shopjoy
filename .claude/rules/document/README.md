# Project Directory Convention

项目目录规范，确保文档、SQL、代码三合一对应。

## 一、业务领域划分

项目共 10 个业务领域：

| 领域 | 英文名 | 说明 |
|-----|--------|------|
| 用户与权限 | user | 用户、租户、角色、管理员 |
| 商品目录 | product | 商品、分类、品牌、市场 |
| 订单 | order | 订单、购物车 |
| 促销 | promotion | 促销活动、优惠券 |
| 积分 | points | 积分系统 |
| 店铺设置 | shop | 店铺配置 |
| 店铺装修 | storefront | 主题、页面装修 |
| 履约 | fulfillment | 发货、物流、退款 |
| 支付 | payment | 支付、回调 |
| 评价 | review | 评价、回复 |

## 二、文档规范 (`docs/`)

### 目录结构

```
docs/
├── README.md                      # 文档索引
├── ARCHITECTURE.md                # 架构文档
├── {YYYY-MM-DD}-product-roadmap.md
│
├── domains/                       # 业务领域文档
│   ├── user/
│   │   ├── {日期}-user-prd.md
│   │   ├── {日期}-user-schema.md
│   │   └── design/
│   │       └── {日期}-user-ui-design.md
│   ├── product/
│   ├── order/
│   ├── promotion/
│   ├── points/
│   ├── shop/
│   ├── storefront/
│   ├── fulfillment/
│   ├── payment/
│   └── review/
│
├── cross-cutting/                 # 跨领域文档
│   ├── api/                       # API 文档
│   │   ├── README.md
│   │   ├── {日期}-api-reference.md
│   │   └── openapi.yaml
│   └── tech-design/               # 技术设计
│       └── {日期}-{名称}-design.md
│
├── guides/                        # 开发指南
│   ├── {日期}-onboarding.md
│   ├── {日期}-developer-guide.md
│   └── {日期}-user-guide.md
│
├── reference/                     # 参考资料
│   ├── {日期}-error-codes.md
│   ├── {日期}-database-overview.md
│   └── {日期}-code-documentation.md
│
├── plans/                         # 开发计划
├── superpowers/                   # AI Agent 文档
├── _archive/                      # 归档文档
└── _templates/                    # 文档模板
```

### 命名规范

#### 基本格式

```
{YYYY-MM-DD}-{领域}-{文档类型}.md
```

#### 文档类型对照表

| 文档类型 | 命名格式 | 示例 |
|---------|---------|------|
| PRD | `{日期}-{领域}-prd.md` | `2026-03-24-order-prd.md` |
| Schema 文档 | `{日期}-{领域}-schema.md` | `2026-03-24-user-schema.md` |
| UI 设计 | `{日期}-{领域}-ui-design.md` | `2026-03-24-payment-ui-design.md` |
| 领域设计 | `{日期}-{领域}-domain-design.md` | `2026-03-22-fulfillment-domain-design.md` |
| 子功能 PRD | `{日期}-{领域}-{功能}-prd.md` | `2026-03-25-fulfillment-shipping-prd.md` |
| 技术设计（跨领域） | `{日期}-{名称}-design.md` | `2026-03-22-sku-code-generation-design.md` |
| 指南 | `{日期}-{名称}.md` | `2026-03-18-onboarding.md` |
| 参考资料 | `{日期}-{名称}.md` | `2026-03-21-error-codes.md` |

### 文档位置规则

| 文档类型 | 放置位置 |
|---------|---------|
| 领域 PRD | `docs/domains/{领域}/` |
| 领域设计 | `docs/domains/{领域}/` |
| UI 设计 | `docs/domains/{领域}/design/` |
| 跨领域技术设计 | `docs/cross-cutting/tech-design/` |
| API 文档 | `docs/cross-cutting/api/` |
| 开发指南 | `docs/guides/` |
| 参考资料 | `docs/reference/` |
| 开发计划 | `docs/plans/` |

## 三、SQL 规范 (`sql/`)

### 目录结构

```
sql/
├── init.sql                       # 初始化入口
│
├── user/                          # 用户与权限领域
│   ├── schema.sql                 # 表结构定义（合并）
│   └── migrations/                # 迁移脚本
│       ├── 2026032401_create_user_addresses.sql
│       └── ...
│
├── product/                       # 商品领域
│   ├── schema.sql
│   └── migrations/
│
├── order/                         # 订单领域
├── promotion/                     # 促销领域
├── points/                        # 积分领域
├── shop/                          # 店铺设置领域
├── storefront/                    # 店铺装修领域
├── fulfillment/                   # 履约领域
├── payment/                       # 支付领域
└── review/                        # 评价领域
```

### Schema 文件规范

- 每个领域一个 `schema.sql` 文件
- 合并该领域所有相关表的定义
- 按表依赖顺序排列

示例：`sql/user/schema.sql` 包含 tenant、user、admin_user、role 等表。

### 迁移文件规范

#### 命名格式

```
{YYYYMMDD}{序号}_{动作}_{对象}.sql
```

#### 动作词

| 动作 | 说明 | 示例 |
|-----|------|------|
| create | 创建表 | `2026032401_create_reviews.sql` |
| alter | 修改表结构 | `2026032201_alter_promotions_add_scope.sql` |
| add | 添加字段/索引 | `2026032301_add_order_fields.sql` |
| drop | 删除字段/表 | `2026032401_drop_old_column.sql` |

#### 序号规则

- 每天从 01 开始
- 同一天多个迁移按执行顺序递增

## 四、代码目录规范

### 后端目录结构

```
admin/                             # 或 shop/
├── desc/                          # API 定义文件
│   ├── admin.api                  # 主入口
│   ├── user.api                   # 用户模块
│   └── product.api                # 商品模块
│
└── internal/
    ├── handler/                   # HTTP 处理器（自动生成）
    ├── logic/                     # 业务逻辑
    ├── types/                     # 类型定义（自动生成）
    ├── svc/                       # 服务上下文
    ├── config/                    # 配置
    │
    ├── domain/                    # 领域层
    │   ├── user/                  # 用户领域
    │   ├── product/               # 商品领域
    │   ├── order/                 # 订单领域
    │   └── ...
    │
    ├── application/               # 应用层
    │   ├── user/
    │   └── product/
    │
    └── infrastructure/            # 基础设施层
        └── persistence/           # 仓储实现
```

### 前端目录结构

```
shop-admin/                        # 或 joy/
├── src/
│   ├── views/                     # 页面组件
│   │   ├── login/
│   │   ├── products/
│   │   └── orders/
│   ├── components/                # 公共组件
│   ├── layouts/                   # 布局组件
│   ├── stores/                    # Pinia 状态管理
│   ├── api/                       # API 接口
│   └── router/                    # 路由配置
└── package.json
```

## 五、三合一对应关系

文档、SQL、代码目录按领域一一对应：

| 领域 | 文档目录 | SQL 目录 | 代码目录 |
|-----|---------|---------|---------|
| user | `docs/domains/user/` | `sql/user/` | `domain/{user,adminuser,role,tenant}/` |
| product | `docs/domains/product/` | `sql/product/` | `domain/{product,market}/` |
| order | `docs/domains/order/` | `sql/order/` | `domain/{order,cart}/` |
| promotion | `docs/domains/promotion/` | `sql/promotion/` | `domain/{promotion,coupon}/` |
| points | `docs/domains/points/` | `sql/points/` | 待创建 |
| shop | `docs/domains/shop/` | `sql/shop/` | `handler/shop/` |
| storefront | `docs/domains/storefront/` | `sql/storefront/` | `domain/storefront/` |
| fulfillment | `docs/domains/fulfillment/` | `sql/fulfillment/` | `domain/fulfillment/` |
| payment | `docs/domains/payment/` | `sql/payment/` | `domain/payment/` |
| review | `docs/domains/review/` | `sql/review/` | `domain/review/` |

## 六、核心原则

### MUST

| # | 规则 | 说明 |
|---|------|------|
| 1 | 三合一对应 | 文档、SQL、代码使用相同领域划分 |
| 2 | 日期前缀 | 文档和迁移文件使用日期前缀 |
| 3 | 领域隔离 | 每个领域独立的目录 |
| 4 | 跨领域放 cross-cutting | 跨领域的文档统一放置 |

### SHOULD

| # | 规则 | 说明 |
|---|------|------|
| 5 | 归档机制 | 过期文档移入 `_archive/` |
| 6 | 模板复用 | 使用 `_templates/` 中的模板 |
| 7 | 索引维护 | 更新 `docs/README.md` 索引 |

### FORBIDDEN

| # | 规则 | 说明 |
|---|------|------|
| 8 | 随意放置文档 | 必须按领域组织 |
| 9 | 省略日期前缀 | 文档必须有日期前缀 |
| 10 | 迁移放根目录 | 迁移必须放在领域 migrations 目录 |

## 七、示例

### 新增领域功能

假设要为 `order` 领域新增「订单导出」功能：

```
# 1. 创建 PRD 文档
docs/domains/order/2026-03-26-order-export-prd.md

# 2. 创建迁移文件（如需新表）
sql/order/migrations/2026032601_create_order_exports.sql

# 3. 创建代码
admin/internal/domain/order/export_entity.go
admin/internal/logic/orders/export_order_logic.go
```

### 新增跨领域功能

假设要新增「数据同步」功能（跨多个领域）：

```
# 创建技术设计文档
docs/cross-cutting/tech-design/2026-03-26-data-sync-design.md
```