# ShopJoy - 电商 SaaS 平台

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/Vue-3.4%2B-green)](https://vuejs.org/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

ShopJoy 是一个基于 **go-zero + DDD（领域驱动设计）** 的轻量级多租户独立电商 SaaS 平台，帮助中小企业 5 分钟开通独立店铺。

## 核心特性

- **多租户架构** - 支持多商家同时运营，数据严格隔离
- **DDD 领域驱动** - 清晰的限界上下文，易于维护和扩展
- **完整电商链路** - 商品、订单、支付、促销、店铺管理全覆盖
- **高性能** - 基于 go-zero 框架，支持高并发
- **前后端分离** - Vue3 前端，支持多主题

## 项目架构

```
shopjoy/
├── admin/                    # 管理后台 API (go-zero)
│   ├── internal/
│   │   ├── application/      # 应用层 - 用例编排
│   │   ├── domain/           # 领域层 - 实体/值对象/仓储接口
│   │   └── infrastructure/   # 基础设施层 - 仓储实现
│   └── desc/                 # API 定义文件
├── shop/                     # 商城 API (go-zero)
│   └── internal/
│       ├── application/      # 购物车、订单、支付应用服务
│       └── domain/           # 订单、购物车、支付领域
├── pkg/                      # 全局共享包
│   ├── domain/shared/        # Money, TenantID, DomainEvent
│   ├── tenant/               # 多租户上下文管理
│   ├── asyncq/               # 异步队列 (asynq)
│   └── auth/                 # JWT 认证
├── shop-admin/               # 管理后台前端 (Vue3 + Element Plus)
├── joy/                      # 商城前端 (Vue3 + Tailwind CSS)
└── migrations/               # 数据库迁移脚本
```

## 技术栈

### 后端
- **框架**: [go-zero](https://go-zero.dev/) - 微服务框架
- **ORM**: [GORM](https://gorm.io/) - 数据库操作
- **数据库**: MySQL 8.0 + Redis
- **缓存**: Redis (Cluster 模式)
- **队列**: asynq (基于 Redis)
- **认证**: JWT
- **JSON**: sonic (高性能)

### 前端
- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **UI 组件库**: Element Plus (admin) / Tailwind CSS (joy)
- **状态管理**: Pinia
- **路由**: Vue Router

## 功能模块

### 1. 身份认证 (Identity & Access)
- 用户注册/登录 (JWT + Refresh Token)
- RBAC 权限管理 (角色/权限)
- 多租户隔离

### 2. 商品管理 (Catalog)
- SPU + SKU 管理
- 分类、品牌管理
- 库存管理 (Redis 预占)
- 商品状态机 (草稿→上架→下架)

### 3. 订单系统 (Sales & Order)
- 购物车 (游客/登录用户)
- 订单创建与状态机
- 订单超时自动关闭
- 订单售后

### 4. 支付系统 (Payment)
- 多渠道支付 (支付宝/微信)
- 支付回调处理
- 退款管理

### 5. 促销系统 (Promotion)
- 优惠券管理
- 促销活动 (秒杀/满减/折扣)
- 促销规则引擎

### 6. 店铺管理 (Storefront)
- 店铺配置
- 主题切换
- 页面装修
- SEO 设置

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- MySQL 8.0
- Redis 7.0

### 1. 克隆项目

```bash
git clone https://github.com/colinrs/shopjoy.git
cd shopjoy
```

### 2. 启动基础设施

```bash
# 使用 Docker 启动 MySQL 和 Redis
docker-compose up -d mysql redis
```

### 3. 数据库迁移

```bash
# 运行迁移脚本
go run scripts/migrate.go -f admin/etc/admin-api.yaml
```

### 4. 启动后端服务

```bash
# 启动 Admin 服务 (端口 8888)
cd admin && go run admin.go -f etc/admin-api.yaml

# 启动 Shop 服务 (端口 8889)
cd shop && go run shop.go -f etc/shop-api.yaml
```

### 5. 启动前端

```bash
# 管理后台 (端口 3000)
cd shop-admin
npm install
npm run dev

# 商城前端 (端口 3001)
cd joy
npm install
npm run dev
```

### 6. Docker 一键部署

```bash
# 构建并启动所有服务
docker-compose up -d
```

访问地址:
- 管理后台: http://localhost:3000
- 商城前端: http://localhost:3001
- Admin API: http://localhost:8888
- Shop API: http://localhost:8889

## 文档

- [架构设计文档](docs/ARCHITECTURE.md) - DDD 战略设计、分层架构
- [API 文档](docs/API.md) - 接口定义和使用说明
- [产品需求文档](docs/prd/prd-1.md) - 功能需求详细说明

## 开发指南

### 生成 API 代码

```bash
cd admin && make api
cd shop && make api
```

### 运行测试

```bash
go test ./...
```

### 代码规范

```bash
golangci-lint run --timeout=10m
```

## 数据库表结构

| 表名 | 说明 |
|------|------|
| tenants | 租户表 |
| users | 用户表 |
| roles | 角色表 |
| user_roles | 用户角色关联 |
| products | 商品表 |
| skus | SKU 表 |
| categories | 分类表 |
| brands | 品牌表 |
| orders | 订单表 |
| order_items | 订单商品表 |
| carts | 购物车表 |
| cart_items | 购物车商品表 |
| payments | 支付表 |
| coupons | 优惠券表 |
| user_coupons | 用户优惠券表 |
| shops | 店铺表 |
| themes | 主题表 |

## API 接口

### 认证接口
- `POST /users/register` - 用户注册
- `POST /users/login` - 用户登录
- `PUT /users/password` - 修改密码

### 商品接口
- `GET /api/v1/products` - 商品列表
- `POST /api/v1/products` - 创建商品
- `GET /api/v1/products/:id` - 商品详情
- `PUT /api/v1/products/:id` - 更新商品
- `POST /api/v1/products/:id/on-sale` - 上架
- `POST /api/v1/products/:id/off-sale` - 下架

### 订单接口
- `GET /api/v1/orders` - 订单列表
- `POST /api/v1/orders` - 创建订单
- `GET /api/v1/orders/:id` - 订单详情
- `PUT /api/v1/orders/:id/cancel` - 取消订单

### 购物车接口
- `GET /api/v1/cart` - 获取购物车
- `POST /api/v1/cart/items` - 添加商品
- `PUT /api/v1/cart/items/:id` - 更新数量
- `DELETE /api/v1/cart/items/:id` - 删除商品

完整 API 文档请查看 [docs/API.md](docs/API.md)

## 架构特点

### DDD 分层架构
```
┌─────────────────────────────────────┐
│         Interface Layer             │  Handler / API
├─────────────────────────────────────┤
│        Application Layer            │  Service / DTO / Use Case
├─────────────────────────────────────┤
│          Domain Layer               │  Entity / Value Object / Repository Interface
├─────────────────────────────────────┤
│      Infrastructure Layer           │  Repository Impl / External Service
└─────────────────────────────────────┘
```

### 6 大限界上下文
1. **Identity & Access** - 身份认证与权限
2. **Catalog** - 商品目录
3. **Sales & Order** - 销售订单
4. **Promotion** - 促销活动
5. **Storefront** - 店铺前台
6. **Payment** - 支付

## 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## License

本项目基于 MIT 协议开源，详见 [LICENSE](LICENSE) 文件。

## 致谢

- [go-zero](https://go-zero.dev/) - 优秀的微服务框架
- [Vue.js](https://vuejs.org/) - 渐进式前端框架
- [Element Plus](https://element-plus.org/) - UI 组件库

---

**ShopJoy** - 让开店更简单