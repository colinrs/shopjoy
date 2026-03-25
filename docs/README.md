# ShopJoy 文档中心

## 文档结构

```
docs/
├── ARCHITECTURE.md              # 项目架构文档
├── 2026-03-25-product-roadmap.md # 产品路线图
│
├── domains/                     # 业务领域文档
│   ├── user/                    # 用户与权限
│   ├── product/                 # 商品目录
│   ├── order/                   # 订单
│   ├── promotion/               # 促销
│   ├── points/                  # 积分
│   ├── shop/                    # 店铺设置
│   ├── storefront/              # 店铺装修
│   ├── fulfillment/             # 履约
│   ├── payment/                 # 支付
│   └── review/                  # 评价
│
├── cross-cutting/               # 跨领域文档
│   ├── tech-design/             # 技术设计
│   └── api/                     # API文档
│
├── guides/                      # 开发指南
├── reference/                   # 参考资料
├── plans/                       # 开发计划
├── superpowers/                 # AI Agent文档
│
├── _archive/                    # 归档文档
└── _templates/                  # 文档模板
```

## 业务领域

| 领域 | 说明 | 文档目录 |
|-----|------|---------|
| user | 用户与权限（用户、租户、角色） | [domains/user/](domains/user/) |
| product | 商品目录（商品、分类、品牌、市场） | [domains/product/](domains/product/) |
| order | 订单（订单、购物车） | [domains/order/](domains/order/) |
| promotion | 促销（促销活动、优惠券） | [domains/promotion/](domains/promotion/) |
| points | 积分系统 | [domains/points/](domains/points/) |
| shop | 店铺设置 | [domains/shop/](domains/shop/) |
| storefront | 店铺装修（主题、页面） | [domains/storefront/](domains/storefront/) |
| fulfillment | 履约（发货、物流、退款） | [domains/fulfillment/](domains/fulfillment/) |
| payment | 支付（支付、回调） | [domains/payment/](domains/payment/) |
| review | 评价（评价、回复） | [domains/review/](domains/review/) |

## 文档命名规范

```
{YYYY-MM-DD}-{领域}-{文档类型}.md

示例:
2026-03-24-order-prd.md          # 订单PRD
2026-03-24-user-schema.md        # 用户表结构
2026-03-24-payment-ui-design.md  # 支付UI设计
```

## 快速导航

### PRD文档
- [用户与权限 PRD](domains/user/2026-03-24-user-prd.md)
- [商品目录 PRD](domains/product/) - 待创建
- [订单 PRD](domains/order/2026-03-24-order-prd.md)
- [促销 PRD](domains/promotion/2026-03-22-promotion-prd.md)
- [积分 PRD](domains/points/2026-03-24-points-prd.md)
- [店铺设置 PRD](domains/shop/2026-03-25-shop-prd.md)
- [店铺装修 PRD](domains/storefront/2026-03-24-storefront-prd.md)
- [履约 PRD](domains/fulfillment/2026-03-22-fulfillment-prd.md)
- [物流 PRD](domains/fulfillment/2026-03-25-fulfillment-shipping-prd.md)
- [支付 PRD](domains/payment/2026-03-24-payment-prd.md)
- [评价 PRD](domains/review/2026-03-24-review-prd.md)

### Schema文档
- [商品目录 Schema](domains/product/2026-03-26-product-schema.md)
- [订单 Schema](domains/order/2026-03-26-order-schema.md)
- [促销 Schema](domains/promotion/2026-03-26-promotion-schema.md)
- [支付 Schema](domains/payment/2026-03-26-payment-schema.md)

### API文档
- [API设计规范](cross-cutting/api/README.md)
- [API参考文档](cross-cutting/api/2026-03-26-api-reference.md)
- [OpenAPI规范](cross-cutting/api/openapi.yaml)
- [缺失API设计](cross-cutting/api/2026-03-25-missing-apis-design.md)

### 开发指南
- [入职指南](guides/2026-03-18-onboarding.md)
- [开发者指南](guides/2026-03-22-developer-guide.md)
- [用户指南](guides/2026-03-22-user-guide.md)

### 参考资料
- [错误码](reference/2026-03-21-error-codes.md)
- [数据库概览](reference/2026-03-22-database-overview.md)
- [代码文档规范](reference/2026-03-22-code-documentation.md)
- [架构图](reference/2026-03-22-architecture-diagrams.md)

### 技术设计
- [促销设计](cross-cutting/tech-design/2026-03-22-promotion-design.md)
- [SKU编码生成设计](cross-cutting/tech-design/2026-03-22-sku-code-generation-design.md)

## SQL文件结构

```
sql/
├── init.sql                     # 初始化入口
├── user/                        # 用户与权限
│   ├── schema.sql
│   └── migrations/
├── product/                     # 商品目录
├── order/                       # 订单
├── promotion/                   # 促销
├── points/                      # 积分
├── shop/                        # 店铺设置
├── storefront/                  # 店铺装修
├── fulfillment/                 # 履约
├── payment/                     # 支付
└── review/                      # 评价
```

## 目录对应关系

| 领域 | 文档 | SQL | 代码 |
|-----|------|-----|------|
| user | docs/domains/user/ | sql/user/ | admin/internal/domain/{user,adminuser,role,tenant}/ |
| product | docs/domains/product/ | sql/product/ | admin/internal/domain/{product,market}/ |
| order | docs/domains/order/ | sql/order/ | shop/internal/domain/{order,cart}/ |
| promotion | docs/domains/promotion/ | sql/promotion/ | admin/internal/domain/{promotion,coupon}/ |
| points | docs/domains/points/ | sql/points/ | 待创建 |
| shop | docs/domains/shop/ | sql/shop/ | admin/internal/handler/shop/ |
| storefront | docs/domains/storefront/ | sql/storefront/ | admin/internal/domain/storefront/ |
| fulfillment | docs/domains/fulfillment/ | sql/fulfillment/ | admin/internal/domain/fulfillment/ |
| payment | docs/domains/payment/ | sql/payment/ | shop/internal/domain/payment/ |
| review | docs/domains/review/ | sql/review/ | admin/internal/domain/review/ |