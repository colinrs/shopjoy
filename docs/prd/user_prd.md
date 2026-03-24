# User Management Module Product Requirements Document

## Document Information

| Item | Value |
|------|-------|
| Document Title | User Management Module PRD (Admin) |
| Version | 1.0.0 |
| Status | Draft |
| Created | 2026-03-24 |
| Updated | 2026-03-24 |
| Author | Product Team |
| Module Location | admin/internal/domain/user/, admin/internal/domain/admin_user/ |

---

## Product Overview

### Summary

本 PRD 描述 ShopJoy Admin 用户管理模块的功能需求。该模块为平台和商家提供安全、高效的用户管理体系，支持多租户隔离，满足运营、客服、风控等日常需求。

**核心功能**：
- C端用户（买家）管理 - 查看、搜索、冻结、详情
- 商家管理员管理 - 新增、编辑、权限分配、禁用
- 平台超级管理员管理 - 全平台权限管理
- 用户地址管理
- 用户统计与导出

### Problem Statement

ShopJoy 平台当前用户管理存在以下痛点：

- **用户类型混乱**：C端用户和管理员账号未明确区分
- **权限管理不足**：商家子账号权限粒度不够
- **平台管理缺失**：缺少平台级别的超级管理员体系
- **用户详情不完整**：无法查看用户关联的订单、积分、地址等信息
- **操作审计薄弱**：对用户的关键操作缺乏日志记录

### Solution Overview

在 admin 服务中完善用户管理模块，提供以下能力：

1. **三層用户架构** - 平台超管、商家管理员、C端用户分离管理
2. **完整用户画像** - 基本信息、地址、订单、积分、评论整合展示
3. **精细化权限** - 基于角色的权限控制（RBAC）
4. **操作审计** - 关键操作日志记录

---

## Goals

### Business Goals

| ID | Goal | Metric | Target |
|----|------|--------|--------|
| BG-001 | 提升用户管理效率 | 用户查询平均耗时 | < 30秒 |
| BG-002 | 加强平台安全 | 越权操作次数 | 0次 |
| BG-003 | 提升客服响应速度 | 用户问题定位时间 | < 1分钟 |

### User Goals

| ID | Goal | User Type |
|----|------|-----------|
| UG-001 | 管理C端用户信息 | Merchant Admin |
| UG-002 | 管理商家子账号 | Merchant Admin |
| UG-003 | 管理平台超管账号 | Platform Super Admin |
| UG-004 | 查看用户完整画像 | Customer Service |

### Non-Goals

| ID | Non-Goal | Reason |
|----|----------|--------|
| NG-001 | 用户自助注册/登录 | 属于 shop 服务 |
| NG-002 | 用户积分操作 | 已在 Points 模块实现 |
| NG-003 | 批量导入用户 | Phase 2 feature |
| NG-004 | 用户行为分析 | Phase 2 feature |

---

## User Personas

### Key User Types

| Persona | Description | Primary Use Cases |
|---------|-------------|-------------------|
| Platform Super Admin | ShopJoy 平台最高权限管理员 | 管理平台超管、查看全平台商家、系统配置 |
| Merchant Admin | 商家店长/运营负责人 | 管理C端用户、管理子账号、查看数据 |
| Merchant Sub Account | 商家运营/客服人员 | 查看用户信息、处理用户问题 |
| Customer Service | 客服人员 | 查询用户、处理投诉、调整积分 |

### Role-Based Access

| 功能模块 | 平台超管 | 商家主账号 | 商家子账号 |
|----------|---------|-----------|-----------|
| C端用户列表 | ✅ 全平台 | ✅ 本租户 | ✅ 本租户 |
| C端用户详情 | ✅ 全平台 | ✅ 本租户 | ✅ 本租户 |
| 冻结/解冻用户 | ✅ 全平台 | ✅ 本租户 | ❌ |
| 调整用户积分 | ✅ 全平台 | ✅ 本租户 | ❌ |
| 商家管理员管理 | ✅ | ❌ | ❌ |
| 平台超管管理 | ✅ | ❌ | ❌ |

---

## Functional Requirements

### Feature 1: C端用户管理

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-001 | 用户列表 | 分页展示C端用户，支持筛选和搜索 | P0 |
| FR-002 | 用户搜索 | 按邮箱、手机号、昵称搜索 | P0 |
| FR-003 | 用户筛选 | 按状态、注册时间、积分区间筛选 | P0 |
| FR-004 | 用户详情 | 查看用户完整信息 | P0 |
| FR-005 | 冻结用户 | 禁止用户登录和下单 | P0 |
| FR-006 | 解冻用户 | 恢复用户正常状态 | P0 |
| FR-007 | 删除用户 | 软删除用户账号 | P1 |
| FR-008 | 用户导出 | 导出用户列表CSV | P1 |
| FR-009 | 用户统计 | 展示用户数量统计 | P0 |

### Feature 2: 用户详情页

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-010 | 基本信息 Tab | 展示用户基本资料 | P0 |
| FR-011 | 地址管理 Tab | 展示用户收货地址列表 | P0 |
| FR-012 | 订单记录 Tab | 展示用户订单历史 | P0 |
| FR-013 | 积分明细 Tab | 展示用户积分变动记录 | P0 |
| FR-014 | 评论记录 Tab | 展示用户发表的评论 | P1 |
| FR-015 | 操作日志 Tab | 展示管理员对该用户的操作记录 | P1 |

### Feature 3: 商家管理员管理

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-020 | 管理员列表 | 分页展示商家管理员 | P0 |
| FR-021 | 新增管理员 | 创建商家管理员账号 | P0 |
| FR-022 | 编辑管理员 | 修改管理员信息 | P0 |
| FR-023 | 禁用管理员 | 禁止管理员登录后台 | P0 |
| FR-024 | 启用管理员 | 恢复管理员登录权限 | P0 |
| FR-025 | 删除管理员 | 软删除管理员账号 | P1 |
| FR-026 | 重置密码 | 重置管理员密码 | P0 |
| FR-027 | 角色分配 | 为管理员分配角色 | P0 |

### Feature 4: 平台超级管理员管理

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-030 | 超管列表 | 展示平台超管账号列表 | P0 |
| FR-031 | 新增超管 | 创建平台超管账号 | P0 |
| FR-032 | 编辑超管 | 修改超管信息 | P0 |
| FR-033 | 禁用超管 | 禁止超管登录 | P0 |
| FR-034 | 启用超管 | 恢复超管登录权限 | P0 |
| FR-035 | 删除超管 | 软删除超管账号 | P1 |
| FR-036 | 重置密码 | 重置超管密码 | P0 |

### Feature 5: 用户地址管理

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-040 | 地址列表 | 查看用户收货地址 | P0 |
| FR-041 | 地址详情 | 查看单个地址详情 | P0 |

---

## User Experience

### C端用户查询流程

1. 商家客服进入「顾客管理」页面
2. 在搜索框输入用户邮箱或手机号
3. 系统实时返回匹配结果
4. 点击用户行进入详情页
5. 在详情页查看用户完整信息

### 冻结用户流程

1. 在用户详情页点击「冻结」按钮
2. 系统弹出确认对话框，要求输入冻结原因
3. 确认后用户状态变为「冻结」
4. 用户无法登录前台，无法下单
5. 操作记录到日志

### 新增商家管理员流程

1. 平台超管进入「管理员管理」页面
2. 点击「新增管理员」按钮
3. 填写用户名、邮箱、密码等信息
4. 选择账号类型（主账号/子账号）
5. 分配角色权限
6. 保存后管理员可登录后台

---

## Narrative

商家运营小明需要查看一个顾客的购物记录。他进入「顾客管理」页面，输入顾客的手机号进行搜索，找到了该顾客。点击进入详情页后，他看到顾客的基本信息、积分余额 500 分、共有 12 笔订单、总消费 $1,200。

顾客反映最近一笔订单没有收到积分，小明切换到「积分明细」Tab，发现该订单确实未产生积分记录。他点击「调整积分」按钮，选择「增加积分」，输入 100 分，填写原因「订单 #12345 积分补发」，确认后积分立即到账。

之后，小明需要为新入职的客服创建后台账号。他请求店长在「管理员管理」页面创建了子账号，分配了「客服」角色，新同事即可登录后台处理客户问题。

---

## Success Metrics

### User-Centric Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| 用户查询成功率 | 首次搜索找到目标用户占比 | > 95% |
| 详情页加载时间 | 用户详情页 P95 加载时间 | < 500ms |
| 操作成功率 | 冻结/解冻操作成功率 | 100% |

### Technical Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| API 响应时间 | 列表查询 P95 | < 200ms |
| 数据隔离 | 租户数据泄漏次数 | 0 |
| 系统可用性 | 用户管理服务可用性 | 99.9% |

---

## Technical Considerations

### Architecture

用户管理模块采用 DDD 架构，**复用现有 AdminUser 实体支持三种用户类型**：

```
admin/internal/
├── domain/
│   ├── user/           # C端用户领域
│   │   ├── entity.go
│   │   ├── repository.go
│   │   └── address_entity.go  # 新增
│   ├── adminuser/      # 管理员领域（现有，扩展）
│   │   ├── entity.go   # 复用：Type=1 平台超管, Type=2 商家主账号, Type=3 商家子账号
│   │   └── repository.go
│   └── role/           # 角色领域（现有）
│       ├── entity.go
│       └── repository.go
├── application/
│   └── user/
│       └── service.go
└── handler/
    └── users/
```

**重要架构决策**：

现有 `admin_users` 表已支持三种管理员类型（见 `admin/internal/domain/adminuser/entity.go`）：
- `TypeSuperAdmin (1)` + `TenantID == 0`：平台超级管理员
- `TypeTenantAdmin (2)` + `TenantID > 0`：商家主账号
- `TypeTenantSub (3)` + `TenantID > 0`：商家子账号

**因此不再创建独立的 `platform_admins` 表**，通过现有实体区分用户类型。

### 现有 API 增强说明

本 PRD 是对现有 API 的**增强**，而非完全重新设计：

| API 文件 | 现有状态 | PRD 增强 |
|----------|----------|----------|
| `admin/desc/user.api` | 已有基础 CRUD | 增加 `/addresses`, `/export`, 搜索增强 |
| `admin/desc/admin_user.api` | 已有列表、禁用/启用 | 增加创建、删除、角色分配、密码重置 |
| `admin/internal/domain/user/entity.go` | 已有 User 实体 | 无需修改 |
| `admin/internal/domain/adminuser/entity.go` | 已支持三种类型 | 无需修改 |

**新增内容**：
1. `user_addresses` 表及相关实体
2. 用户详情聚合 API（积分、订单、评论等）
3. 管理员创建/删除/角色分配 API
4. 前端完整页面

### Database Design

#### admin_users 表（现有，扩展）

**说明**：该表已存在，支持三种管理员类型。平台超级管理员使用 `Type=1` 和 `TenantID=0`。

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | 主键 |
| tenant_id | BIGINT | NO | 0 | 租户ID（平台超管为0） |
| username | VARCHAR(64) | NO | | 用户名 |
| email | VARCHAR(128) | NO | | 邮箱 |
| password | VARCHAR(255) | NO | | 密码哈希 |
| mobile | VARCHAR(20) | NO | '' | 手机号 |
| real_name | VARCHAR(32) | NO | '' | 真实姓名 |
| avatar | VARCHAR(255) | NO | '' | 头像URL |
| type | TINYINT | NO | 1 | 1=平台超管 2=商家主账号 3=商家子账号 |
| status | TINYINT | NO | 1 | 1=启用 2=禁用 3=已删除 |
| last_login_at | TIMESTAMP | YES | NULL | 最后登录时间（UTC） |
| last_login_ip | VARCHAR(45) | NO | '' | 最后登录IP |
| deleted_at | TIMESTAMP | YES | NULL | 删除时间（UTC） |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 创建时间（UTC） |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 更新时间（UTC） |

**Indexes:**
- `uk_username` UNIQUE (username)
- `uk_email` UNIQUE (email)
- `uk_mobile` UNIQUE (mobile)
- `idx_tenant_id` (tenant_id)
- `idx_type` (type)
- `idx_status` (status)
- `idx_deleted_at` (deleted_at)

#### users 表（现有，扩展）

**说明**：该表已存在，用于C端用户（买家）。

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | 主键 |
| tenant_id | BIGINT | NO | | 租户ID |
| email | VARCHAR(100) | NO | '' | 邮箱 |
| phone | VARCHAR(20) | NO | '' | 手机号 |
| password | VARCHAR(255) | NO | '' | 密码哈希 |
| name | VARCHAR(100) | NO | '' | 昵称 |
| avatar | VARCHAR(500) | NO | '' | 头像URL |
| gender | TINYINT | NO | 0 | 0=未知 1=男 2=女 3=其他 |
| birthday | TIMESTAMP | YES | NULL | 生日 |
| status | TINYINT | NO | 1 | 0=未激活 1=正常 2=冻结 |
| last_login | TIMESTAMP | YES | NULL | 最后登录（UTC） |
| deleted_at | TIMESTAMP | YES | NULL | 删除时间（UTC） |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 创建时间（UTC） |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 更新时间（UTC） |

**Indexes:**
- `uk_tenant_email` UNIQUE (tenant_id, email)
- `uk_tenant_phone` UNIQUE (tenant_id, phone)
- `idx_tenant_status` (tenant_id, status)
- `idx_created_at` (created_at)

#### user_addresses 表（新增）

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | 主键 |
| tenant_id | BIGINT | NO | | 租户ID |
| user_id | BIGINT | NO | | 用户ID |
| name | VARCHAR(100) | NO | | 收货人姓名 |
| phone | VARCHAR(20) | NO | | 收货人电话 |
| country | VARCHAR(50) | NO | '' | 国家代码 |
| province | VARCHAR(100) | NO | '' | 省份/州 |
| city | VARCHAR(100) | NO | '' | 城市 |
| district | VARCHAR(100) | NO | '' | 区/县 |
| address | VARCHAR(255) | NO | | 详细地址 |
| postal_code | VARCHAR(20) | NO | '' | 邮编 |
| is_default | TINYINT | NO | 0 | 是否默认地址 |
| deleted_at | TIMESTAMP | YES | NULL | 删除时间（UTC） |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 创建时间（UTC） |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 更新时间（UTC） |

**Indexes:**
- `idx_tenant_user` (tenant_id, user_id)
- `idx_user_id` (user_id)

---

## API Endpoints

### C端用户 API

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/users | 获取用户列表 |
| GET | /api/v1/users/:id | 获取用户详情 |
| PUT | /api/v1/users/:id | 更新用户信息 |
| POST | /api/v1/users/:id/suspend | 冻结用户 |
| POST | /api/v1/users/:id/activate | 启用用户 |
| DELETE | /api/v1/users/:id | 删除用户 |
| GET | /api/v1/users/:id/addresses | 获取用户地址列表 |
| GET | /api/v1/users/stats | 获取用户统计 |
| GET | /api/v1/users/export | 导出用户列表 |

### 商家管理员 API

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/admin-users | 获取管理员列表 |
| POST | /api/v1/admin-users | 创建管理员 |
| GET | /api/v1/admin-users/:id | 获取管理员详情 |
| PUT | /api/v1/admin-users/:id | 更新管理员信息 |
| DELETE | /api/v1/admin-users/:id | 删除管理员 |
| POST | /api/v1/admin-users/:id/enable | 启用管理员 |
| POST | /api/v1/admin-users/:id/disable | 禁用管理员 |
| POST | /api/v1/admin-users/:id/reset-password | 重置密码 |
| PUT | /api/v1/admin-users/:id/roles | 分配角色 |

### 平台超管 API

**说明**：平台超管复用 `/api/v1/admin-users` API，通过 `type=1` 和 `tenant_id=0` 区分。

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/admin-users?type=1 | 获取超管列表 |
| POST | /api/v1/admin-users | 创建超管（type=1, tenant_id=0） |
| GET | /api/v1/admin-users/:id | 获取超管详情 |
| PUT | /api/v1/admin-users/:id | 更新超管信息 |
| DELETE | /api/v1/admin-users/:id | 删除超管 |
| POST | /api/v1/admin-users/:id/enable | 启用超管 |
| POST | /api/v1/admin-users/:id/disable | 禁用超管 |
| POST | /api/v1/admin-users/:id/reset-password | 重置密码 |

**权限控制**：仅 `type=1`（平台超管）可管理其他平台超管账号。

### API 示例

#### GET /api/v1/users - 用户列表

**Request:**
```http
GET /api/v1/users?page=1&page_size=20&keyword=zhang&status=1
```

**Response:**
```json
{
  "list": [
    {
      "id": 10001,
      "tenant_id": 1,
      "email": "zhang@example.com",
      "phone": "138****8888",
      "name": "张三",
      "avatar": "https://...",
      "status": 1,
      "status_text": "正常",
      "points_balance": 500,
      "order_count": 12,
      "total_spent": "1200.00",
      "last_login": "2026-03-24T10:00:00Z",
      "created_at": "2025-06-10T08:00:00Z"
    }
  ],
  "total": 12500,
  "page": 1,
  "page_size": 20
}
```

#### GET /api/v1/users/:id - 用户详情

**Response:**
```json
{
  "id": 10001,
  "tenant_id": 1,
  "email": "zhang@example.com",
  "phone": "13812345678",
  "name": "张三",
  "avatar": "https://...",
  "gender": 1,
  "gender_text": "男",
  "birthday": "1990-01-15",
  "status": 1,
  "status_text": "正常",
  "points_balance": 500,
  "points_frozen": 0,
  "total_earned": 850,
  "total_redeemed": 350,
  "order_count": 12,
  "total_spent": "1200.00",
  "review_count": 5,
  "last_login": "2026-03-24T10:00:00Z",
  "created_at": "2025-06-10T08:00:00Z",
  "updated_at": "2026-03-24T10:00:00Z"
}
```

**计算字段数据源说明**：

| 字段 | 数据源 | 计算方式 |
|------|--------|----------|
| points_balance | points_accounts 表 | 实时查询 |
| points_frozen | points_accounts 表 | 实时查询 |
| total_earned | points_accounts 表 | 实时查询 |
| total_redeemed | points_accounts 表 | 实时查询 |
| order_count | orders 表 | COUNT 聚合 |
| total_spent | orders 表 | SUM(paid_amount) |
| review_count | reviews 表 | COUNT 聚合 |

#### POST /api/v1/users/:id/suspend - 冻结用户

**Request:**
```json
{
  "reason": "违规操作，多次恶意下单"
}
```

**Response:**
```json
{
  "id": 10001,
  "status": 2,
  "status_text": "冻结",
  "updated_at": "2026-03-24T12:00:00Z"
}
```

#### PUT /api/v1/admin-users/:id/roles - 分配角色

**Request:**
```json
{
  "role_ids": [2, 3]
}
```

**Response:**
```json
{
  "id": 26,
  "roles": [
    {"id": 2, "name": "运营", "code": "operator"},
    {"id": 3, "name": "客服", "code": "customer_service"}
  ],
  "updated_at": "2026-03-24T12:00:00Z"
}
```

#### GET /api/v1/users/:id/addresses - 用户地址列表

**Response:**
```json
{
  "list": [
    {
      "id": 1,
      "name": "张三",
      "phone": "13812345678",
      "country": "CN",
      "province": "广东省",
      "city": "深圳市",
      "district": "南山区",
      "address": "科技园路100号",
      "postal_code": "518000",
      "is_default": true,
      "created_at": "2025-06-10T08:00:00Z"
    }
  ],
  "total": 2
}
```

#### GET /api/v1/users/export - 导出用户

**Request:**
```http
GET /api/v1/users/export?status=1&registered_start=2025-01-01&registered_end=2025-12-31
```

**Response:** CSV 文件下载

**导出限制**：最大 10,000 条记录

#### POST /api/v1/admin-users - 创建管理员

**Request:**
```json
{
  "tenant_id": 1,
  "username": "ops02",
  "email": "ops2@shop.com",
  "password": "SecurePass123",
  "mobile": "13900001111",
  "real_name": "运营小李",
  "type": 3,
  "role_ids": [2, 3]
}
```

**Response:**
```json
{
  "id": 26,
  "username": "ops02",
  "email": "ops2@shop.com",
  "type": 3,
  "type_text": "商家子账号",
  "created_at": "2026-03-24T12:00:00Z"
}
```

**说明**：`tenant_id` 字段仅平台超管可指定，商家管理员创建时自动使用当前租户。

---

## Entity Definitions

### UserAddress Entity (新增)

```go
// admin/internal/domain/user/address_entity.go
package user

import (
    "time"
    "github.com/colinrs/shopjoy/pkg/domain/shared"
)

type UserAddress struct {
    ID          int64            `gorm:"column:id;primaryKey"`
    TenantID    shared.TenantID  `gorm:"column:tenant_id;not null;index"`
    UserID      int64            `gorm:"column:user_id;not null;index"`
    Name        string           `gorm:"column:name;not null;size:100"`
    Phone       string           `gorm:"column:phone;not null;size:20"`
    Country     string           `gorm:"column:country;size:50"`
    Province    string           `gorm:"column:province;size:100"`
    City        string           `gorm:"column:city;size:100"`
    District    string           `gorm:"column:district;size:100"`
    Address     string           `gorm:"column:address;not null;size:255"`
    PostalCode  string           `gorm:"column:postal_code;size:20"`
    IsDefault   bool             `gorm:"column:is_default;default:false"`
    DeletedAt   *time.Time       `gorm:"column:deleted_at"`
    CreatedAt   time.Time        `gorm:"column:created_at;not null"`
    UpdatedAt   time.Time        `gorm:"column:updated_at;not null"`
}

func (a *UserAddress) TableName() string {
    return "user_addresses"
}

func (a *UserAddress) SetAsDefault() {
    a.IsDefault = true
}

func (a *UserAddress) UnsetDefault() {
    a.IsDefault = false
}
```

### UserAddress Repository Interface

```go
type AddressRepository interface {
    Create(ctx context.Context, db *gorm.DB, address *UserAddress) error
    Update(ctx context.Context, db *gorm.DB, address *UserAddress) error
    Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
    FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*UserAddress, error)
    FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64) ([]*UserAddress, error)
    SetDefault(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID, addressID int64) error
}
```

---

## Business Rules

### 用户规则

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-001 | 租户隔离 | 商家管理员只能查看本租户的C端用户 |
| BR-002 | 平台超管权限 | 平台超管可查看全平台所有用户 |
| BR-003 | 冻结限制 | 冻结用户需填写原因，记录到日志 |
| BR-004 | 软删除 | 用户删除使用软删除，保留数据 |
| BR-005 | 密码强度 | 管理员密码至少6位 |

### 管理员规则

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-010 | 用户名唯一 | 同一租户下管理员用户名唯一 |
| BR-011 | 邮箱唯一 | 同一租户下管理员邮箱唯一 |
| BR-012 | 主账号限制 | 每个租户只能有一个主账号 |
| BR-013 | 自我操作限制 | 管理员不能禁用/删除自己 |
| BR-014 | 重置密码 | 重置密码生成临时密码，首次登录强制修改 |

---

## Error Codes

用户管理模块使用现有错误码，新增错误码在未使用的范围内定义：

### 用户模块 (使用现有 11xxx)

| Constant | HTTP Status | Code | Message | Status |
|----------|-------------|------|---------|--------|
| ErrUserNotFound | 404 | 11004 | user not found | 现有 |
| ErrUserDuplicateUser | 409 | 11005 | duplicate user | 现有 |
| ErrUserPasswordTooWeak | 400 | 11003 | password too weak | 现有 |
| ErrUserAlreadyDeleted | 400 | 11007 | user already deleted | 现有 |
| ErrUserSuspended | 400 | 11009 | user is suspended | **新增** |
| ErrUserCannotSuspendSelf | 400 | 11010 | cannot suspend yourself | **新增** |
| ErrUserCannotDeleteSelf | 400 | 11011 | cannot delete yourself | **新增** |
| ErrAddressNotFound | 404 | 11012 | address not found | **新增** |

### 管理员模块 (使用现有 10xxx)

| Constant | HTTP Status | Code | Message | Status |
|----------|-------------|------|---------|--------|
| ErrAdminUserNotFound | 404 | 10004 | admin user not found | 现有 |
| ErrAdminDuplicateUser | 409 | 10005 | duplicate admin user | 现有 |
| ErrAdminPasswordTooWeak | 400 | 10003 | password too weak | 现有 |
| ErrAdminCannotDeleteSelf | 400 | 10007 | cannot delete yourself | 现有 |
| ErrAdminAlreadyDeleted | 400 | 10008 | user already deleted | 现有 |
| ErrAdminCannotDisableSelf | 400 | 10012 | cannot disable yourself | **新增** |
| ErrAdminInvalidType | 400 | 10013 | invalid admin type | **新增** |

### 平台超管模块 (使用现有 10xxx，复用 AdminUser)

平台超管复用 AdminUser 实体，使用 Type=1 和 TenantID=0 标识，无需独立错误码。

---

## User Stories

### US-001: 查看用户列表

**Description**: 作为商家管理员，我想查看C端用户列表，了解用户基本情况。

**Acceptance Criteria**:
- Given 我在顾客管理页面
- When 页面加载完成
- Then 我看到用户列表，包含邮箱、手机、昵称、积分、订单数
- And 列表支持分页
- And 默认按注册时间倒序排列

---

### US-002: 搜索用户

**Description**: 作为客服人员，我想通过邮箱或手机号快速找到用户。

**Acceptance Criteria**:
- Given 我在用户列表页面
- When 我在搜索框输入邮箱或手机号
- Then 列表实时过滤显示匹配结果
- And 支持模糊搜索

---

### US-003: 冻结用户

**Description**: 作为商家管理员，我想冻结违规用户，阻止其登录和下单。

**Acceptance Criteria**:
- Given 我在用户详情页
- When 我点击「冻结」按钮
- And 输入冻结原因
- Then 用户状态变为「冻结」
- And 用户无法登录前台
- And 操作记录到日志

---

### US-004: 查看用户详情

**Description**: 作为客服人员，我想查看用户完整信息，帮助解决用户问题。

**Acceptance Criteria**:
- Given 我点击某个用户
- Then 我看到用户基本信息
- And 我可以切换到地址、订单、积分等Tab
- And 每个Tab内容正确加载

---

### US-005: 创建商家管理员

**Description**: 作为平台超管，我想为商家创建管理员账号。

**Acceptance Criteria**:
- Given 我在管理员管理页面
- When 我点击「新增管理员」
- And 填写用户名、邮箱、密码
- And 选择租户和类型
- And 分配角色
- Then 管理员创建成功
- And 管理员可以登录后台

---

### US-006: 分配角色

**Description**: 作为平台超管，我想为管理员分配不同的角色权限。

**Acceptance Criteria**:
- Given 我在管理员详情页
- When 我点击「分配角色」
- And 选择一个或多个角色
- Then 管理员权限立即更新
- And 下次登录时生效

---

### US-007: 创建平台超管

**Description**: 作为平台超管，我想创建新的平台超管账号。

**Acceptance Criteria**:
- Given 我在平台管理员页面
- When 我点击「新增超管」
- And 填写用户名、邮箱、密码
- Then 超管创建成功
- And 超管拥有全平台权限

---

## Frontend Structure

### 页面结构

```
shop-admin/src/views/
├── users/                          # C端用户管理
│   ├── index.vue                   # 用户列表
│   ├── [id].vue                    # 用户详情
│   └── components/
│       ├── UserFilter.vue          # 筛选组件
│       ├── UserStatsCard.vue       # 统计卡片
│       ├── UserDetailTabs.vue      # 详情Tab容器
│       ├── UserBasicInfo.vue       # 基本信息Tab
│       ├── UserAddressList.vue     # 地址列表Tab
│       ├── UserOrderHistory.vue    # 订单记录Tab
│       ├── UserPointsHistory.vue   # 积分明细Tab
│       ├── UserReviewHistory.vue   # 评论记录Tab
│       └── UserOperationLog.vue    # 操作日志Tab
│
└── admin-users/                    # 管理员管理（商家 + 平台）
    ├── index.vue                   # 管理员列表（支持类型筛选）
    ├── [id].vue                    # 管理员详情/编辑
    └── components/
        ├── AdminUserForm.vue       # 新增/编辑表单
        ├── AdminUserFilter.vue     # 筛选组件（含类型筛选）
        └── RoleAssignDialog.vue    # 角色分配弹窗
```

**说明**：平台超管和商家管理员共用 `admin-users` 页面，通过类型筛选区分：
- 平台超管：`type=1`
- 商家主账号：`type=2`
- 商家子账号：`type=3`

平台超管登录时，可看到全部三种类型；商家管理员登录时，仅看到本租户的 `type=2,3`。

### API 模块

```typescript
// src/api/user.ts
export const userApi = {
  list: (params: ListUsersRequest) => request.get('/api/v1/users', { params }),
  get: (id: number) => request.get(`/api/v1/users/${id}`),
  update: (id: number, data: UpdateUserRequest) => request.put(`/api/v1/users/${id}`, data),
  suspend: (id: number, reason: string) => request.post(`/api/v1/users/${id}/suspend`, { reason }),
  activate: (id: number) => request.post(`/api/v1/users/${id}/activate`),
  delete: (id: number) => request.delete(`/api/v1/users/${id}`),
  addresses: (id: number) => request.get(`/api/v1/users/${id}/addresses`),
  stats: () => request.get('/api/v1/users/stats'),
  export: (params: ExportUsersRequest) => request.get('/api/v1/users/export', { params, responseType: 'blob' }),
}

// src/api/adminUser.ts
export const adminUserApi = {
  list: (params: ListAdminUsersRequest) => request.get('/api/v1/admin-users', { params }),
  get: (id: number) => request.get(`/api/v1/admin-users/${id}`),
  create: (data: CreateAdminUserRequest) => request.post('/api/v1/admin-users', data),
  update: (id: number, data: UpdateAdminUserRequest) => request.put(`/api/v1/admin-users/${id}`, data),
  delete: (id: number) => request.delete(`/api/v1/admin-users/${id}`),
  enable: (id: number) => request.post(`/api/v1/admin-users/${id}/enable`),
  disable: (id: number) => request.post(`/api/v1/admin-users/${id}/disable`),
  resetPassword: (id: number) => request.post(`/api/v1/admin-users/${id}/reset-password`),
  assignRoles: (id: number, roleIds: number[]) => request.put(`/api/v1/admin-users/${id}/roles`, { role_ids: roleIds }),
}
```

---

## Specifications

### 搜索行为

| 字段 | 搜索方式 | 说明 |
|------|----------|------|
| keyword | 模糊搜索 | 同时匹配 email、phone、name |
| email | 精确匹配 | 完整邮箱地址 |
| phone | 精确匹配 | 完整手机号 |
| status | 精确匹配 | 状态值筛选 |

### 排序规则

| 列表 | 默认排序 | 支持排序字段 |
|------|----------|--------------|
| 用户列表 | created_at DESC | created_at, last_login, points_balance |
| 管理员列表 | created_at DESC | created_at, last_login_at, status |

### 导出限制

| 参数 | 限制 |
|------|------|
| 最大记录数 | 10,000 条 |
| 时间范围 | 最大 1 年 |
| 格式 | CSV (UTF-8 with BOM) |

### 密码规则

| 规则 | 要求 |
|------|------|
| 最小长度 | 6 位 |
| 字符要求 | 无强制要求 |
| 重置密码 | 生成随机 12 位密码，首次登录强制修改 |

---

## Milestones and Sequencing

### Phase 1: Backend Foundation (Week 1)

| Task | Duration | Description |
|------|----------|-------------|
| Database migration | 0.5 day | Create user_addresses table; extend admin_users API |
| Domain entities | 1 day | Implement UserAddress entity; extend User repository |
| Repositories | 1 day | Implement AddressRepository; extend existing repositories |
| Business logic | 1.5 days | Implement logic for user addresses, user details aggregation |
| Unit tests | 0.5 day | Unit tests for business logic |

### Phase 2: API Development (Week 2)

| Task | Duration | Description |
|------|----------|-------------|
| API definitions | 0.5 day | Extend user.api, admin_user.api |
| Handlers | 1.5 days | Implement HTTP handlers for new endpoints |
| Integration tests | 0.5 day | API integration tests |

### Phase 3: Frontend Development (Week 2-3)

| Task | Duration | Description |
|------|----------|-------------|
| API clients | 0.5 day | Create/extend API TypeScript modules |
| C端用户页面 | 2 days | 用户列表、详情、筛选组件 |
| 管理员页面 | 1.5 days | 管理员列表（含类型筛选）、表单、角色分配 |
| Router & menu | 0.5 day | 添加路由和菜单项 |

### Phase 4: Testing & Polish (Week 3)

| Task | Duration | Description |
|------|----------|-------------|
| End-to-end testing | 1 day | 测试完整用户流程 |
| Bug fixes | 0.5 day | 修复发现的问题 |
| Documentation | 0.5 day | API文档更新 |

**Total Estimate: 3 weeks**

---

## Appendix

### 用户状态

| Status | Value | Description |
|--------|-------|-------------|
| Inactive | 0 | 未激活 |
| Active | 1 | 正常 |
| Suspended | 2 | 冻结 |

### 管理员类型

| Type | Value | Description |
|------|-------|-------------|
| MerchantAdmin | 2 | 商家主账号 |
| MerchantSubAccount | 3 | 商家子账号 |

### 管理员状态

| Status | Value | Description |
|--------|-------|-------------|
| Active | 1 | 启用 |
| Disabled | 2 | 禁用 |