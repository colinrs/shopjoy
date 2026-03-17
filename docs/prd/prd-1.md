**ShopGo SaaS 电商平台 PRD（V1.3 - go-zero + DDD 完整版）**

**版本**：V1.3  
**日期**：2026-03-17  
**作者**：电商 SaaS 产品专家（专注独立站 SaaS，精通 go-zero + DDD）

### 1. 项目概述

**产品名称**：shopjoy SaaS  
**产品定位**：轻量级、多租户独立电商 SaaS 平台，帮助中小企业 5 分钟开通独立店铺，支持自有域名、主题切换、商品管理、订单履约、营销促销全链路。  
**核心价值**：低代码开店、高性能、强扩展性、代码结构清晰长期可维护。  
**目标用户**：
- 商家（Tenant）：中小企业品牌、跨境卖家、个体商户
- 平台方：超级管理员
- 消费者：访客/注册买家

**MVP 目标**：3 个月上线，支持 500+ 商家同时运营，日订单峰值 5000 单。

**技术约束**：
- 后端统一使用 **go-zero** 框架（单体应用模式，后期可按 Bounded Context 拆微服务）
- 严格遵循 **DDD（领域驱动设计）** 指导业务与代码对齐
- 异步任务统一使用 **asyncq**（基于 Redis 的轻量异步队列）
- 项目分为：`admin`、`shop`、`pkg`、`shop-admin`、`joy`

### 2. 项目目录结构（go-zero + DDD 风格）

```bash
shopjoy/
├── admin/                          # 管理后台 API（go-zero）
│   ├── internal/
│   │   ├── config/
│   │   ├── handler/
│   │   ├── logic/                      # Application Layer
│   │   ├── domain/                     # DDD 领域层（按 Bounded Context 组织）
│   │   ├── repository/
│   │   ├── event/
│   │   ├── asyncq/                     # 异步任务统一入口
│   │   ├── types/
│   │   └── middleware/                 # tenant、auth、log 等
│   ├── etc/admin.yaml
│   ├── admin.go
│   └── desc/                           # *.api 文件
│
├── shop/                           # 商城买家端 API（go-zero）
│   ├── internal/ （结构同 admin，领域模型可复用 pkg）
│   ├── etc/shop.yaml
│   ├── shop.go
│   └── desc/
│
├── pkg/                                # 全局公共包
│   ├── tenant/                         # 多租户上下文、中间件
│   ├── domain/                         # 共享 Value Object、Event、Error
│   ├── repository/                     # 基础 Repository 接口
│   ├── asyncq/                         # 异步任务队列核心实现（Redis）
│   ├── errors/
│   ├── utils/
│   ├── cache/
│   └── middleware/
│
├── shop-admin/                         # Vue3 管理后台前端（Element Plus）
├── joy/                               # Vue3 商城前端（买家端，多主题）
├── deploy/
├── docs/                               # DDD 上下文地图、ER 图、API 文档
├── go.mod
└── README.md
```

### 3. DDD 战略设计 - Bounded Contexts（限界上下文）

本项目共划分 6 个核心限界上下文，每个上下文都有独立的 Ubiquitous Language、Aggregate、Repository 和 Domain Event。

1. **Identity & Access Context**（身份与权限上下文）
    - 租户（Tenant）、用户、角色、权限、RBAC
    - 多租户隔离核心（tenant_id 全局贯穿）

2. **Catalog Context**（商品目录上下文）
    - 商品（Product）、SKU、规格、分类、库存、属性
    - Aggregate 根：Product

3. **Sales & Order Context**（销售与订单上下文）
    - 购物车（Cart）、订单（Order）、订单项、支付、售后
    - Aggregate 根：Order（包含支付、优惠快照）

4. **Promotion Context**（促销上下文）※已修改
    - 促销活动（Promotion）、优惠券（Coupon）、促销规则（Rule）
    - 秒杀、满减、折扣、买赠等
    - Domain Event：CouponUsed、PromotionStarted 等

5. **Storefront Context**（店铺前台上下文）
    - 店铺（Shop）、主题（Theme）、页面装修、SEO、首页配置
    - 买家浏览、商品推荐

6. **Fulfillment Context**（履约上下文）
    - 发货、物流单、退款处理、售后订单

**共享内核（Shared Kernel）**：放在 `int/domain/shared/` 和 `domain/shared/`，包含 Money、Status、TenantID、DomainEvent 等。

**多租户策略**：所有 Aggregate 均包含 `TenantID`，通过 `pkg/tenant` middleware 在请求入口自动注入并校验，数据库采用 `tenant_id` + 索引隔离。

### 4. 各上下文核心功能清单（MVP 范围）

#### 4.1 Identity & Access Context
- 平台超级管理员管理
- 商家入驻申请、审核、开通店铺（Tenant）
- 用户登录（账号密码 + 手机号验证码）
- RBAC 权限体系（角色、菜单、按钮权限）
- 员工邀请与权限分配
- JWT + Refresh Token 双令牌机制

#### 4.2 Catalog Context
- 商品 SP U + 多 SKU 管理（支持多规格）
- 商品分类、品牌、属性模板
- 库存管理（Redis 预占 + 确认扣减）
- 商品上下架、批量导入导出
- 商品搜索（Elasticsearch，后续）

#### 4.3 Sales & Order Context
- 购物车（游客 + 登录用户）
- 下单流程（库存校验、促销计算、订单生成）
- 订单状态机管理（待支付 → 已支付 → 待发货 → 已完成 → 退款）
- 支付回调处理（异步）
- 订单列表、详情、售后申请
- 订单超时自动关闭（asyncq）

#### 4.4 Promotion Context（促销上下文）
- 促销活动创建（秒杀、满减、折扣）
- 优惠券全生命周期（创建、发放、核销、过期）
- 促销规则引擎（条件 + 优惠动作）
- 下单时最优促销自动匹配
- 促销数据统计与效果分析
- 异步任务：活动开始/结束通知、优惠券过期清理、统计任务

#### 4.5 Storefront Context
- 店铺基础设置（名称、LOGO、域名、公告）
- 主题切换（内置 3-5 个主题，支持 JSON 配置覆盖）
- 店铺首页、商品列表、详情页装修基础功能
- SEO 元数据配置
- 商品推荐（Redis 热点 + 简单规则）

#### 4.6 Fulfillment Context（MVP 简化版）
- 订单发货操作
- 物流单号录入与查询（对接快递鸟）
- 基础退款流程

### 5. 非功能性需求

- **性能**：商品详情接口 < 80ms，创建订单 < 300ms（P95）
- **并发**：支持单租户 200 QPS，平台整体 5000+ QPS
- **一致性**：库存扣减使用 Redis Lua 脚本 + asyncq 最终一致性
- **安全性**：HTTPS、参数校验、防 SQL 注入、敏感数据加密、多租户严格隔离
- **可扩展性**：清晰的 DDD 分层，便于后续拆微服务或新增上下文
- **监控**：Prometheus + Grafana + go-zero trace + Loki
- **异步处理**：所有耗时/非核心操作统一走 `asyncq`

### 6. 数据存储
- **主数据库**：PostgreSQL（tenant_id 分区索引）
- **缓存**：Redis（Cluster 模式，缓存商品、促销、购物车、Session）
- **文件**：MinIO / 阿里云 OSS
- **异步队列**：Redis（asyncq 实现）
- **搜索**：Elasticsearch（Phase 2）
