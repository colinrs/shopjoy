# ShopJoy - 电商 SaaS 平台

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/Vue-3.4%2B-green)](https://vuejs.org/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

ShopJoy 是一个基于 **go-zero + DDD（领域驱动设计）** 的轻量级多租户独立电商 SaaS 平台，帮助中小企业 5 分钟开通独立店铺。

## 核心特性

- **多租户架构** - 支持多商家同时运营，数据严格隔离
- **DDD 领域驱动** - 清晰的限界上下文，易于维护和扩展
- **完整电商链路** - 商品、订单、支付、促销、店铺管理全覆盖
- **专业前端界面** - 现代化电商 UI，支持响应式设计
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
│   ├── src/views/            # 页面组件
│   │   ├── login/            # 登录页
│   │   ├── dashboard/        # 数据概览
│   │   ├── products/         # 商品管理
│   │   ├── orders/           # 订单管理
│   │   ├── users/            # 用户管理
│   │   ├── promotions/       # 促销管理
│   │   └── shop/             # 店铺设置
│   └── src/layouts/          # 布局组件
├── joy/                      # 商城前端 (Vue3 + Tailwind CSS)
│   ├── src/views/            # 页面组件
│   │   ├── home/             # 首页
│   │   ├── login/            # 登录页
│   │   ├── products/         # 商品列表/详情
│   │   ├── cart/             # 购物车
│   │   ├── checkout/         # 结算页
│   │   ├── orders/           # 订单列表
│   │   └── user/             # 用户中心
│   └── src/components/       # 公共组件
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

### 前端 - 管理后台 (shop-admin)
- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **UI 组件库**: Element Plus
- **图表库**: ECharts (数据可视化)
- **状态管理**: Pinia
- **路由**: Vue Router
- **图标**: @element-plus/icons-vue

### 前端 - 商城 (joy)
- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **CSS 框架**: Tailwind CSS
- **状态管理**: Pinia
- **路由**: Vue Router
- **图标**: Heroicons

## 功能模块

### 1. 身份认证 (Identity & Access)
- 用户注册/登录 (JWT + Refresh Token)
- RBAC 权限管理 (角色/权限)
- 多租户隔离

### 2. 商品管理 (Catalog)
- SPU + SKU 管理
- 分类、品牌管理 (支持三级分类)
- 库存管理 (Redis 预占)
- 商品状态机 (草稿→上架→下架)
- 商品多语言 (支持 en/zh-CN/ja/de/fr 等)
- 商品合规信息 (HS Code、原产地、危险品声明)

### 3. 多市场销售 (Multi-Market)
- 市场配置 (CN/US/UK/DE/FR/AU)
- 多币种定价
- 税务规则配置 (VAT/GST/IOSS)
- 分类/品牌市场可见性控制

### 4. 库存管理 (Inventory)
- 多仓库管理
- 库存预占/确认扣减
- 安全库存预警
- 库存变动日志

### 5. 订单系统 (Sales & Order)
- 购物车 (游客/登录用户)
- 订单创建与状态机
- 订单超时自动关闭
- 订单售后

### 6. 支付系统 (Payment)
- 多渠道支付 (支付宝/微信)
- 支付回调处理
- 退款管理

### 7. 促销系统 (Promotion)
- 优惠券管理
- 促销活动 (秒杀/满减/折扣)
- 促销规则引擎

### 8. 店铺管理 (Storefront)
- 店铺配置
- 主题切换
- 页面装修
- SEO 设置

## 前端页面清单

### 管理后台 (shop-admin)
| 页面 | 路径 | 功能描述 |
|------|------|----------|
| 登录 | `/login` | 分屏式登录/注册，支持 JWT 认证 |
| 数据概览 | `/dashboard` | 统计卡片、销售趋势图表、订单分布、热销商品 |
| 商品管理 | `/products` | 商品表格、批量操作、图片预览、价格格式化 |
| 订单管理 | `/orders` | 订单列表、状态筛选、展开详情、发货操作 |
| 用户管理 | `/users` | 用户卡片、角色标签、消费统计、状态切换 |
| 促销管理 | `/promotions` | 优惠券/活动双标签、进度条、创建编辑 |
| 店铺设置 | `/shop` | 基本信息、SEO设置、状态开关 |

### 商城前台 (joy)
| 页面 | 路径 | 功能描述 |
|------|------|----------|
| 首页 | `/` | Hero区块、分类、限时秒杀、精选商品、服务保障 |
| 登录 | `/login` | 分屏式设计、邮箱/密码、微信登录、手机号登录 |
| 商品列表 | `/products` | 侧边筛选、价格区间、评分筛选、网格/列表视图 |
| 商品详情 | `/products/:id` | 图片画廊、SKU选择、数量控制、评价、FAQ |
| 购物车 | `/cart` | 全选、数量调整、库存警告、优惠券、猜你喜欢 |
| 结算页 | `/checkout` | 地址选择、配送方式、优惠券、订单备注 |
| 订单列表 | `/orders` | 状态统计、筛选标签、订单卡片、操作按钮 |
| 用户中心 | `/user` | 用户信息、订单状态、服务网格、设置列表 |

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

## 前端开发指南

### shop-admin 开发

```bash
cd shop-admin

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

**依赖说明:**
- `echarts` - 数据可视化图表
- `@element-plus/icons-vue` - Element Plus 图标

### joy 商城开发

```bash
cd joy

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

**依赖说明:**
- `@heroicons/vue` - Heroicons 图标库
- Tailwind CSS 已通过 CDN 引入

## 文档

### 核心文档
- [架构设计文档](docs/ARCHITECTURE.md) - DDD 战略设计、分层架构
- [API 参考文档](docs/api-reference.md) - 完整 API 接口定义
- [API 简要文档](docs/API.md) - 接口快速入门
- [数据库设计](docs/database-schema.md) - 表结构设计说明
- [错误码参考](docs/ERROR_CODES.md) - 错误码定义与处理

### 开发指南
- [开发者指南](docs/developer-guide.md) - 开发环境搭建与工作流
- [新手上路](docs/guides/ONBOARDING.md) - 开发者入门指南
- [代码文档](docs/code-documentation.md) - 核心代码说明

### 用户指南
- [用户操作手册](docs/user-guide.md) - 管理后台使用教程

### 规划文档
- [业务开发计划](docs/plans/2026-03-18-business-development.md) - 功能开发路线图
- [产品需求文档](docs/prd/) - 功能需求详细说明
- [前端开发指南](AGENTS.md) - 前端开发规范和最佳实践

## 开发指南

### 生成 API 代码

```bash
cd admin && make api
cd shop && make api
```

### 运行测试

```bash
# 后端测试
go test ./...

# 前端测试 (shop-admin)
cd shop-admin && npm run lint

# 前端测试 (joy)
cd joy && npm run build
```

### 代码规范

```bash
# Go 代码检查
golangci-lint run --timeout=10m

# 前端代码格式化
cd shop-admin && npm run format
```

## 数据库表结构

### 核心业务表
| 表名 | 说明 |
|------|------|
| tenants | 租户表 |
| users | 用户表 |
| roles | 角色表 |
| user_roles | 用户角色关联 |

### 商品相关表
| 表名 | 说明 |
|------|------|
| products | 商品表 (SPU) |
| skus | SKU 表 (商品规格) |
| product_localizations | 商品多语言表 |
| product_markets | 商品市场关联表 |
| categories | 分类表 (支持三级树形结构) |
| category_markets | 分类市场可见性表 |
| brands | 品牌表 |
| brand_markets | 品牌市场可见性表 |

### 市场与库存表
| 表名 | 说明 |
|------|------|
| markets | 市场/区域表 (CN/US/UK/DE/FR/AU) |
| warehouses | 仓库表 |
| warehouse_inventories | 仓库库存表 |
| inventory_logs | 库存变动日志表 |

### 订单与支付表
| 表名 | 说明 |
|------|------|
| orders | 订单表 |
| order_items | 订单商品表 |
| carts | 购物车表 |
| cart_items | 购物车商品表 |
| payments | 支付表 |

### 促销表
| 表名 | 说明 |
|------|------|
| coupons | 优惠券表 |
| user_coupons | 用户优惠券表 |

### 店铺表
| 表名 | 说明 |
|------|------|
| shops | 店铺表 |
| themes | 主题表 |

## API 接口

### 认证接口
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/register` - 用户注册
- `PUT /api/v1/users/password` - 修改密码

### 商品接口
- `GET /api/v1/products` - 商品列表
- `POST /api/v1/products` - 创建商品
- `GET /api/v1/products/:id` - 商品详情
- `PUT /api/v1/products/:id` - 更新商品
- `POST /api/v1/products/:id/on-sale` - 上架
- `POST /api/v1/products/:id/off-sale` - 下架
- `PUT /api/v1/products/:id/stock` - 更新库存

### SKU 接口
- `GET /api/v1/skus/:id` - SKU 详情
- `POST /api/v1/skus` - 创建 SKU
- `PUT /api/v1/skus/:id` - 更新 SKU
- `DELETE /api/v1/skus/:id` - 删除 SKU
- `GET /api/v1/products/:product_id/skus` - 商品 SKU 列表

### 商品多语言接口
- `POST /api/v1/product-localizations` - 创建多语言内容
- `GET /api/v1/product-localizations/:id` - 获取多语言内容
- `GET /api/v1/products/:product_id/localizations` - 商品多语言列表
- `PUT /api/v1/product-localizations/:id` - 更新多语言内容
- `DELETE /api/v1/product-localizations/:id` - 删除多语言内容

### 分类接口
- `GET /api/v1/categories` - 分类列表
- `GET /api/v1/categories/tree` - 分类树
- `POST /api/v1/categories` - 创建分类
- `PUT /api/v1/categories/:id` - 更新分类
- `DELETE /api/v1/categories/:id` - 删除分类
- `PUT /api/v1/categories/:id/move` - 移动分类
- `PUT /api/v1/categories/sort` - 排序分类

### 品牌接口
- `GET /api/v1/brands` - 品牌列表
- `POST /api/v1/brands` - 创建品牌
- `PUT /api/v1/brands/:id` - 更新品牌
- `DELETE /api/v1/brands/:id` - 删除品牌

### 市场接口
- `GET /api/v1/markets` - 市场列表
- `POST /api/v1/markets` - 创建市场
- `PUT /api/v1/markets/:id` - 更新市场
- `DELETE /api/v1/markets/:id` - 删除市场

### 库存接口
- `GET /api/v1/warehouses` - 仓库列表
- `POST /api/v1/warehouses` - 创建仓库
- `PUT /api/v1/inventory/stock` - 更新库存
- `POST /api/v1/inventory/adjust` - 库存调整
- `GET /api/v1/inventory/sku/:sku_code` - SKU 库存详情
- `GET /api/v1/inventory/logs` - 库存日志
- `GET /api/v1/inventory/low-stock` - 低库存预警

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

完整 API 文档请查看 [docs/api-reference.md](docs/api-reference.md)

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

### 7 大限界上下文
1. **Identity & Access** - 身份认证与权限
2. **Catalog** - 商品目录
3. **Sales & Order** - 销售订单
4. **Promotion** - 促销活动
5. **Storefront** - 店铺前台
6. **Payment** - 支付
7. **Fulfillment** - 履约

## 设计资源

本项目使用专业设计系统进行 UI/UX 设计:
- **主色调**: `#059669` ( emerald-600 )
- **强调色**: `#F97316` ( orange-500 )
- **字体**: Fira Sans / Fira Code (shop-admin), Rubik / Nunito Sans (joy)
- **设计系统**: Data-Dense Dashboard (shop-admin), Vibrant Block-based (joy)

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
- [Tailwind CSS](https://tailwindcss.com/) - CSS 框架
- [Heroicons](https://heroicons.com/) - 图标库

---

**ShopJoy** - 让开店更简单
