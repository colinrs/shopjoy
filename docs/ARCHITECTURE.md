# ShopJoy 电商SaaS平台 - 项目架构文档

## 项目概述

ShopJoy 是一个基于 **go-zero + DDD（领域驱动设计）** 的轻量级多租户独立电商SaaS平台。

## 项目结构

```
shopjoy/
├── admin/                          # 管理后台API（go-zero）
│   ├── internal/
│   │   ├── application/            # 应用层 - 用例编排
│   │   │   ├── product/            # 商品应用服务
│   │   │   └── user/               # 用户应用服务
│   │   ├── domain/                 # 领域层 - 核心业务逻辑
│   │   │   ├── product/            # 商品领域（实体、值对象、仓储接口）
│   │   │   ├── role/               # 角色权限领域
│   │   │   ├── tenant/             # 租户领域
│   │   │   ├── user/               # 用户领域
│   │   │   ├── category/           # 分类领域
│   │   │   ├── brand/              # 品牌领域
│   │   │   ├── promotion/          # 促销领域
│   │   │   ├── coupon/             # 优惠券领域
│   │   │   ├── storefront/         # 店铺领域
│   │   │   └── fulfillment/        # 履约领域
│   │   ├── infrastructure/         # 基础设施层
│   │   │   └── persistence/        # 仓储实现
│   │   └── ...
│   └── desc/                       # API定义文件
├── shop/                           # 商城买家端API（go-zero）
│   └── internal/
│       ├── domain/                 # 领域层
│       │   ├── order/              # 订单领域
│       │   ├── cart/               # 购物车领域
│       │   └── payment/            # 支付领域
│       └── application/            # 应用层
├── pkg/                            # 全局公共包
│   ├── domain/shared/              # 共享内核
│   │   ├── value_objects.go        # Money, TenantID, AuditInfo等
│   │   └── event.go                # DomainEvent, EventBus等
│   ├── tenant/                     # 多租户支持
│   │   ├── tenant.go               # Tenant实体和上下文
│   │   └── middleware.go           # 租户中间件
│   ├── asyncq/                     # 异步队列（asynq）
│   │   └── asyncq.go               # 任务队列客户端/服务端
│   ├── application/                # 应用层基础
│   │   └── base.go                 # CommandHandler, QueryHandler
│   ├── cache/                      # 缓存抽象
│   ├── code/                       # 错误码定义
│   ├── infra/                      # 基础设施
│   │   ├── db.go                   # 数据库连接
│   │   └── redis.go                # Redis连接
│   └── ...
├── shop-admin/                     # Vue3管理后台前端
│   ├── src/
│   │   ├── views/                  # 页面视图
│   │   ├── components/             # 组件
│   │   ├── stores/                 # Pinia状态管理
│   │   ├── router/                 # Vue Router
│   │   └── api/                    # API接口
│   └── package.json
├── joy/                            # Vue3商城前端（买家端）
│   ├── src/
│   │   ├── views/
│   │   ├── components/
│   │   ├── stores/
│   │   ├── router/
│   │   └── themes/                 # 多主题支持
│   └── package.json
└── docs/                           # 文档
    └── prd/
        └── prd-1.md                # 产品需求文档
```

## DDD战略设计 - 6大限界上下文

### 1. Identity & Access Context（身份与权限上下文）
- **位置**: `admin/internal/domain/{user,role,tenant}/`
- **核心实体**:
  - `Tenant` - 租户（多租户核心）
  - `User` - 用户
  - `Role` - 角色
  - `Permission` - 权限
- **关键特性**:
  - JWT + Refresh Token 双令牌机制
  - RBAC 权限体系
  - 租户隔离（tenant_id全局贯穿）

### 2. Catalog Context（商品目录上下文）
- **位置**: `admin/internal/domain/product/`
- **核心实体**:
  - `Product` - 商品（聚合根）
  - `SKU` - SKU规格
  - `Category` - 分类
  - `Brand` - 品牌
  - `Attribute` - 属性模板
- **关键特性**:
  - SPU + 多SKU管理
  - 库存管理（Redis预占+确认扣减）
  - 商品状态机（草稿→上架→下架）

### 3. Sales & Order Context（销售与订单上下文）
- **位置**: `shop/internal/domain/{order,cart}/`
- **核心实体**:
  - `Order` - 订单（聚合根）
  - `OrderItem` - 订单项
  - `Cart` - 购物车
  - `CartItem` - 购物车项
- **关键特性**:
  - 订单状态机（待支付→已支付→待发货→已完成）
  - 购物车合并（游客+登录用户）
  - 订单超时自动关闭（asyncq）

### 4. Promotion Context（促销上下文）
- **位置**: `admin/internal/domain/{promotion,coupon}/`
- **核心实体**:
  - `Promotion` - 促销活动
  - `PromotionRule` - 促销规则
  - `Coupon` - 优惠券
  - `UserCoupon` - 用户优惠券
- **关键特性**:
  - 秒杀、满减、折扣
  - 促销规则引擎
  - 优惠券全生命周期管理

### 5. Storefront Context（店铺前台上下文）
- **位置**: `admin/internal/domain/storefront/`
- **核心实体**:
  - `Shop` - 店铺
  - `Theme` - 主题
  - `Page` - 页面装修
  - `Navigation` - 导航菜单
- **关键特性**:
  - 多主题切换
  - 页面装修
  - SEO配置

### 6. Fulfillment Context（履约上下文）
- **位置**: `admin/internal/domain/fulfillment/`
- **核心实体**:
  - `Shipment` - 发货单
  - `Refund` - 退款单
- **关键特性**:
  - 物流跟踪
  - 退款流程

### 7. Payment Context（支付上下文）
- **位置**: `shop/internal/domain/payment/`
- **核心实体**:
  - `Payment` - 支付单
  - `Refund` - 支付退款
- **关键特性**:
  - 多渠道支付（支付宝、微信等）
  - 支付回调处理

## 架构原则

### 1. 分层架构
```
┌─────────────────────────────────────────┐
│           Interface Layer               │
│  (Handler / Controller / API)           │
├─────────────────────────────────────────┤
│         Application Layer               │
│  (Service / DTO / Use Case)             │
├─────────────────────────────────────────┤
│           Domain Layer                  │
│  (Entity / Value Object / Repository)   │
├─────────────────────────────────────────┤
│       Infrastructure Layer              │
│  (RepositoryImpl / External Service)    │
└─────────────────────────────────────────┘
```

### 2. Repository模式
- 仓储接口定义在Domain层
- 仓储实现在Infrastructure层
- **DB/Tx作为参数传递，不存储在Repository中**

### 3. 多租户策略
- 所有聚合根包含 `TenantID`
- 通过 `pkg/tenant` middleware注入租户上下文
- 数据库采用 `tenant_id + 索引` 隔离

### 4. 共享内核 (Shared Kernel)
- `pkg/domain/shared/` 包含所有上下文共享的值对象
- `Money` - 金额（单位为分）
- `TenantID` - 租户ID
- `DomainEvent` - 领域事件
- `AuditInfo` - 审计信息

## 技术栈

### 后端
- **框架**: go-zero (微服务框架)
- **数据库**: PostgreSQL (PRD要求) / MySQL (当前)
- **缓存**: Redis (Cluster模式)
- **ORM**: GORM
- **异步队列**: asynq (基于Redis)
- **JSON**: sonic (高性能)
- **ID生成**: snowflake

### 前端
- **框架**: Vue 3 + TypeScript
- **状态管理**: Pinia
- **UI组件库**: Element Plus (admin)
- **样式**: Tailwind CSS (joy)
- **构建工具**: Vite

## 快速开始

### 安装依赖
```bash
# 后端
go mod download

# 前端 - 管理后台
cd shop-admin && npm install

# 前端 - 商城
cd joy && npm install
```

### 运行服务
```bash
# Admin API
cd admin && go run admin.go -f etc/admin-api.yaml

# Shop API  
cd shop && go run shop.go -f etc/shop-api.yaml

# 前端 - 管理后台
cd shop-admin && npm run dev

# 前端 - 商城
cd joy && npm run dev
```

### 代码生成
```bash
# 生成API代码（修改.desc文件后执行）
cd admin && make api
cd shop && make api
```

## 下一步开发计划

1. **API定义完善**: 使用go-zero的.api文件定义所有接口
2. **仓储实现**: 实现各个领域的GORM仓储
3. **应用服务实现**: 实现应用层服务逻辑
4. **认证授权**: 集成JWT和RBAC中间件
5. **前端页面**: 开发管理后台和商城页面
6. **支付集成**: 接入支付宝/微信支付
7. **物流集成**: 接入快递鸟等物流查询
8. **测试覆盖**: 单元测试和集成测试
