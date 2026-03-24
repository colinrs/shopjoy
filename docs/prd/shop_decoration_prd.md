# Shop Decoration Module Product Requirements Document

## Document Information

| Item | Value |
|------|-------|
| Document Title | Shop Decoration Module PRD (Admin) |
| Version | 1.0.0 |
| Status | Draft |
| Created | 2026-03-24 |
| Updated | 2026-03-24 |
| Author | Product Team |
| Module Location | admin/internal/domain/storefront/ |

---

## Product Overview

### Summary

本 PRD 描述 ShopJoy Admin 页面装修模块的功能需求。该模块为商家后台提供主题切换、可视化页面装修、SEO 配置等核心能力，帮助商家打造个性化的店铺前台。

**核心功能**：
- 主题管理（浏览、切换、基础自定义）
- 可视化页面装修（拖拽式首页装修）
- SEO 设置（全局 + 页面级 SEO 配置）
- 版本历史与回滚
- 装修预览

**Phase 1 范围**：
- 仅支持首页的可视化装修
- 5 套预设主题，支持基础自定义（颜色、字体）
- 全局 SEO + 首页 SEO
- 保存草稿 / 发布 / 预览
- 版本历史（保留最近 10 个版本）

### Problem Statement

商家在店铺运营中，对于店铺外观和页面展示存在以下痛点：

- **无法个性化店铺**：所有店铺外观相同，难以体现品牌特色
- **页面修改依赖开发**：简单的页面调整也需要技术人员介入
- **SEO 配置分散**：SEO 设置散落在多处，难以统一管理
- **缺乏版本控制**：页面修改后无法回滚，误操作难以恢复
- **缺少实时预览**：修改后无法立即看到效果

### Solution Overview

在 admin 服务中新增页面装修模块，提供以下能力：

1. **主题管理** - 预设主题切换 + 基础自定义（颜色、字体）
2. **可视化页面装修** - 拖拽式页面构建，所见即所得
3. **SEO 配置** - 全局默认 + 页面级 SEO 设置
4. **版本管理** - 自动保存历史版本，支持一键回滚
5. **实时预览** - 侧边栏同步预览装修效果

---

## Goals

### Business Goals

| ID | Goal | Metric | Target |
|----|------|--------|--------|
| BG-001 | 提升店铺个性化能力 | 使用装修功能的商家占比 | > 60% |
| BG-002 | 降低页面修改成本 | 页面修改平均耗时 | < 10 分钟 |
| BG-003 | 提升店铺 SEO 效果 | SEO 配置完整率 | > 80% |

### User Goals

| ID | Goal | User Type |
|----|------|-----------|
| UG-001 | 快速切换店铺主题 | Merchant Admin |
| UG-002 | 可视化装修店铺首页 | Merchant Admin / Designer |
| UG-003 | 配置店铺 SEO 信息 | Merchant Admin |
| UG-004 | 预览装修效果 | Merchant Admin |
| UG-005 | 回滚页面版本 | Merchant Admin |

### Non-Goals

| ID | Non-Goal | Reason |
|----|----------|--------|
| NG-001 | 分类页/商品详情页装修 | Phase 2 feature |
| NG-002 | 自定义主题上传 | Phase 2 feature |
| NG-003 | 高级 SEO（Open Graph、结构化数据） | Phase 2 feature |
| NG-004 | 定时发布 | Phase 2 feature |
| NG-005 | 多语言装修 | Phase 2 feature |
| NG-006 | 移动端独立装修 | Phase 2 feature |

---

## User Personas

### Key User Types

| Persona | Description | Primary Use Cases |
|---------|-------------|-------------------|
| Merchant Admin | 商家管理员 | 主题切换、页面装修、SEO 配置 |
| Store Designer | 店铺设计师 | 可视化页面装修、组件配置 |
| Content Editor | 内容运营 | 文案编辑、SEO 优化 |

### Role-Based Access

| Role | 切换主题 | 页面装修 | 发布页面 | SEO 配置 | 版本回滚 |
|------|---------|---------|---------|---------|---------|
| Tenant Admin | ✅ | ✅ | ✅ | ✅ | ✅ |
| Tenant Operations Manager | ✅ | ✅ | ✅ | ✅ | ✅ |
| Tenant Content Editor | ❌ | ✅ | ❌ | ✅ | ❌ |

---

## Functional Requirements

### Feature 1: 主题管理

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-001 | 主题列表 | 展示所有可用主题（卡片形式） | P0 |
| FR-002 | 当前主题标识 | 高亮显示当前使用的主题 | P0 |
| FR-003 | 主题切换 | 一键切换到新主题 | P0 |
| FR-004 | 主题预览 | 预览主题效果（新窗口打开） | P0 |
| FR-005 | 主题自定义 | 自定义主题颜色（主色、辅色） | P0 |
| FR-006 | 字体选择 | 从预设字体中选择 | P0 |
| FR-007 | 按钮样式 | 选择按钮样式（圆角/方角/胶囊） | P1 |
| FR-008 | 重置主题配置 | 恢复主题默认设置 | P0 |

### Feature 2: 页面装修 - 基础功能

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-010 | 页面列表 | 展示可装修的页面列表（首页为主） | P0 |
| FR-011 | 区块库 | 左侧展示可拖拽的区块组件 | P0 |
| FR-012 | 拖拽排序 | 支持拖拽调整区块顺序 | P0 |
| FR-013 | 添加区块 | 从区块库拖拽添加到画布 | P0 |
| FR-014 | 删除区块 | 删除画布中的区块 | P0 |
| FR-015 | 区块编辑 | 右侧面板编辑选中区块的配置 | P0 |
| FR-016 | 侧边预览 | 右侧同步显示预览效果 | P0 |
| FR-017 | 撤销/重做 | 支持撤销和重做操作 | P1 |

### Feature 3: 页面装修 - 区块类型

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-020 | Banner/轮播图 | 支持多图轮播、自动播放、链接跳转 | P0 |
| FR-021 | 商品网格 | 展示商品列表，支持手动选择或自动推荐 | P0 |
| FR-022 | 商品轮播 | 横向滚动商品展示 | P0 |
| FR-023 | 特色商品 | 单个商品重点展示 | P1 |
| FR-024 | 分类导航 | 商品分类入口展示 | P0 |
| FR-025 | 富文本区块 | 支持 HTML 富文本编辑 | P0 |
| FR-026 | 图文区块 | 图片+文字组合布局 | P1 |
| FR-027 | 图片墙 | 多图片网格展示 | P1 |
| FR-028 | 视频播放器 | 嵌入视频播放 | P2 |
| FR-029 | 倒计时器 | 促销活动倒计时 | P1 |
| FR-030 | 促销横幅 | 促销信息横幅展示 | P1 |
| FR-031 | 优惠券展示 | 展示可用优惠券 | P2 |
| FR-032 | 社交链接 | 社交媒体链接入口 | P2 |
| FR-033 | Instagram Feed | Instagram 动态展示 | P2 |
| FR-034 | 搜索栏 | 商品搜索输入框 | P1 |
| FR-035 | 自定义菜单 | 自定义链接菜单 | P1 |
| FR-036 | 分割线 | 视觉分隔区块 | P1 |
| FR-037 | 容器区块 | 包含其他区块的容器 | P2 |

### Feature 4: 版本与发布

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-040 | 保存草稿 | 保存当前编辑状态（不发布） | P0 |
| FR-041 | 发布页面 | 将草稿发布为正式版本 | P0 |
| FR-042 | 取消发布 | 将已发布页面改为草稿状态 | P0 |
| FR-043 | 版本历史列表 | 展示历史版本（最近 10 个） | P0 |
| FR-044 | 版本预览 | 预览历史版本效果 | P0 |
| FR-045 | 版本回滚 | 回滚到指定历史版本 | P0 |
| FR-046 | 版本信息 | 显示版本创建时间、创建人 | P0 |

### Feature 5: SEO 设置

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| FR-050 | 全局 SEO 配置 | 设置默认标题、描述、关键词 | P0 |
| FR-051 | 首页 SEO 配置 | 单独设置首页 SEO 信息 | P0 |
| FR-052 | 标题字数提示 | 标题字数统计与提示（建议 60 字符内） | P0 |
| FR-053 | 描述字数提示 | 描述字数统计与提示（建议 160 字符内） | P0 |
| FR-054 | SEO 预览 | Google 搜索结果样式预览 | P0 |
| FR-055 | 关键词输入 | 多关键词逗号分隔输入 | P0 |

---

## User Experience

### 主题切换流程

1. 商家进入「店铺设置 → 主题管理」页面
2. 看到主题卡片网格，当前使用主题有「使用中」标签
3. 点击其他主题的「预览」按钮，新窗口打开预览效果
4. 确认后点击「应用」按钮
5. 系统切换主题，显示成功提示
6. 可在「主题自定义」区域调整颜色、字体

### 页面装修流程

1. 商家进入「店铺设置 → 页面装修」页面
2. 左侧看到区块库面板，中间是画布，右侧是预览面板
3. 从左侧拖拽区块到画布中
4. 点击画布中的区块，右侧显示属性编辑面板
5. 修改区块属性，预览面板实时更新
6. 拖拽区块调整顺序
7. 点击「保存草稿」保存当前进度
8. 点击「发布」使修改生效

### 版本回滚流程

1. 在页面装修页面点击「历史版本」
2. 查看版本列表（版本号、时间、操作人）
3. 点击「预览」查看历史版本效果
4. 点击「恢复」将页面恢复到该版本
5. 系统创建新版本（基于历史版本内容）

### SEO 配置流程

1. 商家进入「店铺设置 → SEO 设置」页面
2. 先配置全局 SEO（作为默认值）
3. 切换到「首页」标签，配置首页专属 SEO
4. 输入标题、描述、关键词
5. 实时看到 Google 搜索预览效果
6. 点击保存

---

## Narrative

商家小明经营一家服装店铺，想要为即将到来的夏季促销打造全新的店铺形象。

早上，小明打开后台，进入「店铺设置 → 主题管理」，浏览了 5 套预设主题。他喜欢「Modern」主题的简洁风格，点击「预览」确认效果后，点击「应用」切换主题。然后他在主题自定义中，将主色调改为品牌蓝色，字体改为 Poppins。

接下来，小明进入「页面装修」页面。他从左侧区块库拖拽一个「Banner」区块到画布顶部，上传了 3 张夏季促销海报图片，设置了 5 秒自动轮播。然后在 Banner 下方添加一个「商品网格」区块，标题设为「夏日新品」，手动选择了 8 款主打商品。

装修完成后，小明点击「保存草稿」，然后在右侧预览面板确认效果满意。点击「发布」后，店铺首页即刻更新。

最后，小明进入「SEO 设置」，为首页设置了专属的 SEO 标题「Summer Sale - Up to 50% Off | My Store」和描述，期待在搜索引擎中获得更好的曝光。

---

## Success Metrics

### User-Centric Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| 主题切换使用率 | 至少切换过一次主题的商家占比 | > 40% |
| 页面装修使用率 | 使用过页面装修的商家占比 | > 60% |
| SEO 配置完整率 | 首页 SEO 三项全部填写的商家占比 | > 80% |
| 版本回滚使用率 | 使用过版本回滚的商家占比 | < 10%（低代表误操作少） |

### Technical Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| 页面装修保存响应时间 | P95 | < 500ms |
| 页面发布响应时间 | P95 | < 300ms |
| 版本回滚响应时间 | P95 | < 500ms |
| 预览加载时间 | P95 | < 1s |

---

## Technical Considerations

### 现有架构

页面装修模块采用 DDD 架构：

```
admin/internal/domain/storefront/
├── entity.go          # 领域实体：Shop, Theme, Page, Decoration, SEOConfig
├── repository.go      # 仓储接口
└── service.go         # 领域服务
```

API 定义：`admin/desc/storefront.api`

### 数据库设计

#### themes 表

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | 主键 |
| tenant_id | BIGINT | NO | | 租户 ID |
| code | VARCHAR(50) | NO | | 主题编码（default, modern, minimal, bold, nature） |
| name | VARCHAR(100) | NO | | 主题名称 |
| preview_image | VARCHAR(500) | NO | '' | 预览图 URL |
| config_schema | JSON | NO | NULL | 主题可配置项 Schema |
| default_config | JSON | NO | NULL | 默认配置 |
| is_preset | TINYINT | NO | 1 | 是否为系统预设主题 |
| is_active | TINYINT | NO | 1 | 是否可用 |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 创建时间（UTC） |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 更新时间（UTC） |

Indexes:
- `idx_tenant_code` UNIQUE (tenant_id, code)
- `idx_tenant_active` (tenant_id, is_active)

#### shops 表扩展

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| current_theme_id | BIGINT | YES | NULL | 当前使用的主题 ID |
| theme_config | JSON | NO | '{}' | 主题个性化配置 |

#### pages 表

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | 主键 |
| tenant_id | BIGINT | NO | | 租户 ID |
| page_type | VARCHAR(30) | NO | | 页面类型：home, category, product, custom |
| name | VARCHAR(100) | NO | | 页面名称 |
| slug | VARCHAR(200) | NO | | URL 路径 |
| is_published | TINYINT | NO | 0 | 是否已发布 |
| published_at | TIMESTAMP | YES | NULL | 发布时间（UTC） |
| version | INT | NO | 1 | 当前版本号 |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 创建时间（UTC） |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 更新时间（UTC） |
| deleted_at | TIMESTAMP | YES | NULL | 删除时间（UTC） |

Indexes:
- `idx_tenant_slug` UNIQUE (tenant_id, slug)
- `idx_tenant_type` (tenant_id, page_type)

> **Note**: 现有 `admin/internal/domain/storefront/entity.go` 中的 `Page` 实体包含 `Content string` 字段。本设计采用「装饰器模式」，通过新增 `decorations` 表存储区块配置，而非使用 `Content` 字段。`Content` 字段保留用于简单页面场景，装饰型页面使用 `decorations` 关联。实施时需对现有 `pages` 表进行迁移。

#### decorations 表

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | 主键 |
| tenant_id | BIGINT | NO | | 租户 ID |
| page_id | BIGINT | NO | | 页面 ID |
| block_type | VARCHAR(50) | NO | | 区块类型 |
| block_config | JSON | NO | '{}' | 区块配置 |
| sort_order | INT | NO | 0 | 排序序号 |
| is_active | TINYINT | NO | 1 | 是否启用 |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 创建时间（UTC） |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 更新时间（UTC） |

Indexes:
- `idx_page_sort` (page_id, sort_order)
- `idx_tenant_page` (tenant_id, page_id)

#### page_versions 表

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | 主键 |
| tenant_id | BIGINT | NO | | 租户 ID |
| page_id | BIGINT | NO | | 页面 ID |
| version | INT | NO | | 版本号 |
| blocks | JSON | NO | '[]' | 区块快照 |
| created_by | BIGINT | NO | 0 | 创建人 ID |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 创建时间（UTC） |

Indexes:
- `idx_tenant_page_ver` UNIQUE (tenant_id, page_id, version)

#### seo_configs 表

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| id | BIGINT | NO | AUTO_INCREMENT | 主键 |
| tenant_id | BIGINT | NO | | 租户 ID |
| page_type | VARCHAR(30) | NO | | 页面类型：global, home, category, product |
| page_id | BIGINT | YES | NULL | 具体页面 ID（global 时为 NULL） |
| title | VARCHAR(200) | NO | '' | SEO 标题 |
| description | TEXT | NO | '' | SEO 描述 |
| keywords | VARCHAR(500) | NO | '' | SEO 关键词 |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 创建时间（UTC） |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 更新时间（UTC） |

Indexes:
- `idx_tenant_page_type` UNIQUE (tenant_id, page_type, COALESCE(page_id, 0))

### API Endpoints

#### Theme API

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/themes | 获取主题列表 |
| GET | /api/v1/themes/:id | 获取主题详情 |
| PUT | /api/v1/themes/switch | 切换当前主题 |
| PUT | /api/v1/themes/config | 更新主题配置 |
| GET | /api/v1/themes/current | 获取当前主题及配置 |

#### Page & Decoration API

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/pages | 获取页面列表 |
| POST | /api/v1/pages | 创建页面 |
| GET | /api/v1/pages/:id | 获取页面详情（含区块） |
| PUT | /api/v1/pages/:id | 更新页面基本信息 |
| DELETE | /api/v1/pages/:id | 删除页面 |
| PUT | /api/v1/pages/:id/publish | 发布页面 |
| PUT | /api/v1/pages/:id/unpublish | 取消发布 |
| GET | /api/v1/pages/:id/decorations | 获取页面区块列表 |
| POST | /api/v1/pages/:id/decorations | 添加区块 |
| PUT | /api/v1/decorations/:id | 更新区块配置 |
| DELETE | /api/v1/decorations/:id | 删除区块 |
| PUT | /api/v1/pages/:id/blocks/reorder | 批量调整区块顺序 |

#### Version API

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/pages/:id/versions | 获取版本历史 |
| GET | /api/v1/pages/:id/versions/:ver | 获取指定版本 |
| POST | /api/v1/pages/:id/versions | 创建新版本 |
| PUT | /api/v1/pages/:id/restore/:ver | 恢复到指定版本 |

#### SEO API

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/seo/global | 获取全局 SEO 配置 |
| PUT | /api/v1/seo/global | 更新全局 SEO |
| GET | /api/v1/seo/pages | 获取所有页面 SEO 配置 |
| GET | /api/v1/seo/pages/:pageType | 获取指定页面类型 SEO |
| PUT | /api/v1/seo/pages/:pageType | 更新页面 SEO |

#### Preview API

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/preview/page/:id | 获取页面预览数据（草稿版本） |

### API 示例

#### GET /api/v1/themes - 主题列表

**Response:**
```json
{
  "themes": [
    {
      "id": 1,
      "code": "classic",
      "name": "Classic",
      "preview_image": "https://cdn.../classic.png",
      "is_preset": true,
      "is_current": true
    },
    {
      "id": 2,
      "code": "modern",
      "name": "Modern",
      "preview_image": "https://cdn.../modern.png",
      "is_preset": true,
      "is_current": false
    }
  ]
}
```

#### PUT /api/v1/themes/switch - 切换主题

**Request:**
```json
{
  "theme_id": 2
}
```

**Response:**
```json
{
  "success": true,
  "theme": {
    "id": 2,
    "code": "modern",
    "name": "Modern"
  }
}
```

#### PUT /api/v1/themes/config - 更新主题配置

**Request:**
```json
{
  "config": {
    "primary_color": "#3B82F6",
    "secondary_color": "#1E40AF",
    "font_family": "poppins",
    "button_style": "rounded"
  }
}
```

#### GET /api/v1/pages/:id - 获取页面详情

**Response:**
```json
{
  "page": {
    "id": 1,
    "page_type": "home",
    "name": "Homepage",
    "slug": "/",
    "is_published": true,
    "version": 5
  },
  "decorations": [
    {
      "id": 1,
      "block_type": "banner",
      "block_config": {
        "images": [
          {"url": "https://...", "link": "/sale"}
        ],
        "height": 500,
        "autoplay": true,
        "interval": 5000
      },
      "sort_order": 1
    },
    {
      "id": 2,
      "block_type": "product_grid",
      "block_config": {
        "title": "Best Sellers",
        "layout": "grid",
        "columns": 4,
        "product_source": "manual",
        "product_ids": [101, 102, 103],
        "limit": 8
      },
      "sort_order": 2
    }
  ]
}
```

#### POST /api/v1/pages/:id/decorations - 添加区块

**Request:**
```json
{
  "block_type": "rich_text",
  "block_config": {
    "content": "<h2>Welcome</h2><p>Find the best products...</p>"
  },
  "sort_order": 3
}
```

**Response:**
```json
{
  "id": 3,
  "block_type": "rich_text",
  "block_config": {
    "content": "<h2>Welcome</h2><p>Find the best products...</p>"
  },
  "sort_order": 3
}
```

#### PUT /api/v1/pages/:id/blocks/reorder - 调整区块顺序

**Request:**
```json
{
  "block_orders": [
    {"id": 1, "sort_order": 2},
    {"id": 2, "sort_order": 1},
    {"id": 3, "sort_order": 3}
  ]
}
```

#### GET /api/v1/seo/global - 获取全局 SEO

**Response:**
```json
{
  "title": "My Awesome Store - Shop Online",
  "description": "Discover great products at My Awesome Store. Quality products at great prices.",
  "keywords": "online store, shopping, products"
}
```

---

## Business Rules

### 主题规则

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-001 | 主题唯一性 | 每个租户同一时间只能使用一个主题 |
| BR-002 | 预设主题保护 | 系统预设主题不可删除，只能禁用 |
| BR-003 | 配置继承 | 主题配置只存储与默认值不同的部分 |
| BR-004 | 切换保留配置 | 切换主题后，原主题配置保留，再次切回可恢复 |

### 页面规则

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-010 | Slug 唯一 | 同一租户下页面 slug 不能重复 |
| BR-011 | 首页特殊处理 | 首页 slug 固定为 `/`，不可修改 |
| BR-012 | 发布状态约束 | 只有已发布的页面才对前台可见 |
| BR-013 | 区块排序 | 区块按 sort_order 升序排列 |

### 版本规则

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-020 | 自动版本 | 每次保存草稿自动创建新版本 |
| BR-021 | 版本上限 | 每个页面最多保留 10 个版本 |
| BR-022 | 版本清理 | 超过上限时删除最旧版本 |
| BR-023 | 回滚创建版本 | 回滚操作创建新版本，不修改历史版本 |

### SEO 规则

| Rule ID | Rule | Description |
|---------|------|-------------|
| BR-030 | 全局覆盖 | 页面级 SEO 覆盖全局 SEO |
| BR-031 | 标题长度建议 | 标题建议不超过 60 字符（超出显示警告） |
| BR-032 | 描述长度建议 | 描述建议不超过 160 字符（超出显示警告） |
| BR-033 | 关键词去重 | 关键词自动去除重复项 |

---

## Security Requirements

### 权限控制

| 要求 | 说明 |
|-----|------|
| 租户隔离 | 所有数据按租户隔离，API 强制校验 tenant_id |
| 角色权限 | 根据角色控制装修、发布、回滚等操作权限 |
| 操作审计 | 敏感操作（发布、回滚）记录审计日志 |

### 输入验证

| 要求 | 说明 |
|-----|------|
| 区块配置验证 | 服务端验证区块配置 Schema 合法性 |
| HTML 过滤 | 富文本区块过滤危险 HTML 标签和属性 |
| URL 验证 | 链接 URL 必须是合法 URL 格式 |
| 文件上传 | 图片上传通过现有文件服务，不直接处理 |

---

## Error Codes

Storefront Decoration 模块使用 110xxx 错误码范围（Storefront 模块已有 110001-110002 用于 Shop 基础错误，Decoration 模块使用 110100 及之后的子范围）：

| Constant | HTTP Status | Code | Message |
|----------|-------------|------|---------|
| ErrThemeNotFound | 404 | 110101 | theme not found |
| ErrThemeAlreadyActive | 400 | 110102 | theme is already active |
| ErrThemeConfigInvalid | 400 | 110103 | invalid theme configuration |
| ErrPageNotFound | 404 | 110201 | page not found |
| ErrPageSlugDuplicate | 400 | 110202 | page slug already exists |
| ErrPagePublishFailed | 500 | 110203 | failed to publish page |
| ErrDecorationNotFound | 404 | 110301 | decoration block not found |
| ErrInvalidBlockType | 400 | 110302 | invalid block type |
| ErrInvalidBlockConfig | 400 | 110303 | invalid block configuration |
| ErrBlockLimitExceeded | 400 | 110304 | block limit exceeded |
| ErrVersionNotFound | 404 | 110401 | version not found |
| ErrVersionRestoreFailed | 500 | 110402 | failed to restore version |
| ErrSEOConfigNotFound | 404 | 110501 | SEO config not found |

> **Note**: SEO 标题和描述长度超出建议值时，前端显示警告提示但不阻止保存。不定义硬性错误码，避免过度限制商家内容。

---

## User Stories

### US-001: 切换店铺主题

**Description**: 作为商家管理员，我想切换店铺主题，快速改变店铺外观。

**Acceptance Criteria**:
- Given 我在主题管理页面
- When 我点击某主题的「应用」按钮
- Then 系统切换到该主题
- And 店铺前台立即显示新主题样式
- And 该主题显示「使用中」标签

---

### US-002: 自定义主题颜色

**Description**: 作为商家管理员，我想自定义主题颜色，体现品牌特色。

**Acceptance Criteria**:
- Given 我已应用某主题
- When 我在主题自定义区域修改主色调
- Then 预览实时更新
- And 点击保存后，店铺前台应用新颜色

---

### US-003: 添加 Banner 区块

**Description**: 作为商家管理员，我想在首页添加轮播 Banner。

**Acceptance Criteria**:
- Given 我在页面装修页面
- When 我从区块库拖拽「Banner」区块到画布
- And 上传 3 张图片
- And 设置 5 秒自动轮播
- Then 画布显示 Banner 预览
- And 预览面板同步显示效果

---

### US-004: 添加商品展示区块

**Description**: 作为商家管理员，我想在首页展示精选商品。

**Acceptance Criteria**:
- Given 我在页面装修页面
- When 我添加「商品网格」区块
- And 选择「手动选择商品」
- And 选择 8 个商品
- Then 画布显示商品网格预览

---

### US-005: 调整区块顺序

**Description**: 作为商家管理员，我想调整首页区块的显示顺序。

**Acceptance Criteria**:
- Given 画布中已有多个区块
- When 我拖拽某个区块到新位置
- Then 区块顺序更新
- And 预览面板同步更新显示顺序

---

### US-006: 发布页面

**Description**: 作为商家管理员，我想发布装修后的页面，使其对顾客可见。

**Acceptance Criteria**:
- Given 我已完成页面装修
- And 点击「保存草稿」
- When 我点击「发布」按钮
- Then 页面状态变为「已发布」
- And 店铺前台显示新版本

---

### US-007: 回滚到历史版本

**Description**: 作为商家管理员，我想回滚到之前的版本，撤销错误的修改。

**Acceptance Criteria**:
- Given 我在版本历史面板
- When 我选择某历史版本并点击「恢复」
- Then 当前页面恢复为该版本内容
- And 系统创建新版本记录此次回滚

---

### US-008: 配置首页 SEO

**Description**: 作为商家管理员，我想配置首页 SEO，提升搜索引擎排名。

**Acceptance Criteria**:
- Given 我在 SEO 设置页面
- When 我切换到「首页」标签
- And 输入标题、描述、关键词
- Then 我看到 Google 搜索预览效果
- And 保存后前台页面 meta 标签更新

---

## Frontend Structure

### 新增页面

```
shop-admin/src/
├── views/
│   └── shop/
│       ├── index.vue                  # 现有基本设置
│       ├── decoration/
│       │   ├── index.vue              # 页面装修主页面
│       │   ├── components/
│       │   │   ├── BlockLibrary.vue       # 左侧区块库面板
│       │   │   ├── BuilderCanvas.vue      # 中间画布区域
│       │   │   ├── PreviewPanel.vue       # 右侧预览面板
│       │   │   ├── BlockRenderer.vue      # 区块渲染组件
│       │   │   ├── BlockEditor.vue        # 区块属性编辑器
│       │   │   └── VersionHistory.vue     # 版本历史面板
│       │   └── blocks/                    # 区块组件目录
│       │       ├── BannerBlock.vue
│       │       ├── ProductGridBlock.vue
│       │       ├── ProductCarouselBlock.vue
│       │       ├── FeaturedProductBlock.vue
│       │       ├── CollectionProductsBlock.vue
│       │       ├── CategoryNavBlock.vue
│       │       ├── RichTextBlock.vue
│       │       ├── ImageTextBlock.vue
│       │       ├── ImageWallBlock.vue
│       │       ├── VideoPlayerBlock.vue
│       │       ├── CountdownTimerBlock.vue
│       │       ├── PromoBannerBlock.vue
│       │       ├── CouponDisplayBlock.vue
│       │       ├── SocialLinksBlock.vue
│       │       ├── InstagramFeedBlock.vue
│       │       ├── SearchBarBlock.vue
│       │       ├── MenuLinksBlock.vue
│       │       ├── DividerBlock.vue
│       │       └── ContainerBlock.vue
│       ├── themes/
│       │   └── index.vue              # 主题管理页面
│       └── seo/
│           └── index.vue              # SEO 设置页面
├── api/
│   └── storefront.ts                  # API 客户端
├── stores/
│   └── decoration.ts                  # Pinia 状态管理
└── composables/
    └── usePageBuilder.ts              # 拖拽逻辑封装
```

### 组件规格

#### BlockLibrary.vue

- 区块分类展示（Layout / Product / Navigation / Content / Marketing / Social）
- 拖拽手柄
- 区块图标和名称
- 搜索过滤

#### BuilderCanvas.vue

- 拖放接收区域
- 区块占位提示
- 拖拽排序支持
- 区块选中状态
- 空状态提示

#### PreviewPanel.vue

- 区块实时渲染
- 店铺前台样式模拟
- 响应式预览（可选）

#### BlockEditor.vue

- 根据区块类型动态渲染表单
- 表单验证
- 实时预览同步

---

## Block Types Specification

### Layout Blocks

#### Banner (banner)

| Config Field | Type | Description |
|--------------|------|-------------|
| images | Array | 轮播图片列表 |
| images[].url | String | 图片 URL |
| images[].alt | String | 图片 alt 文本 |
| images[].link | String | 点击跳转链接 |
| height | Number | 高度（px） |
| autoplay | Boolean | 是否自动播放 |
| interval | Number | 自动播放间隔（ms） |
| show_arrows | Boolean | 显示左右箭头 |
| show_dots | Boolean | 显示指示点 |
| overlay | Object | 覆盖层配置 |

#### Image Wall (image_wall)

| Config Field | Type | Description |
|--------------|------|-------------|
| images | Array | 图片列表 |
| columns | Number | 列数 |
| gap | Number | 间距（px） |
| aspect_ratio | String | 宽高比（16:9, 1:1, 4:3） |

#### Divider (divider)

| Config Field | Type | Description |
|--------------|------|-------------|
| style | String | 样式（solid, dashed, dotted） |
| color | String | 颜色 |
| margin_top | Number | 上边距 |
| margin_bottom | Number | 下边距 |

#### Container (container)

| Config Field | Type | Description |
|--------------|------|-------------|
| background_color | String | 背景色 |
| padding | Number | 内边距 |
| blocks | Array | 包含的区块 |

### Product Blocks

#### Product Grid (product_grid)

| Config Field | Type | Description |
|--------------|------|-------------|
| title | String | 标题 |
| layout | String | 布局（grid, list） |
| columns | Number | 列数（1-6） |
| product_source | String | 来源（manual, collection, bestseller, new） |
| product_ids | Array | 手动选择的商品 ID |
| collection_id | Number | 分类 ID |
| limit | Number | 显示数量 |
| show_price | Boolean | 显示价格 |
| show_add_to_cart | Boolean | 显示加购按钮 |
| show_rating | Boolean | 显示评分 |

#### Product Carousel (product_carousel)

| Config Field | Type | Description |
|--------------|------|-------------|
| title | String | 标题 |
| product_source | String | 来源 |
| product_ids | Array | 商品 ID |
| collection_id | Number | 分类 ID |
| limit | Number | 显示数量 |
| autoplay | Boolean | 自动滚动 |
| show_arrows | Boolean | 显示箭头 |

#### Featured Product (featured_product)

| Config Field | Type | Description |
|--------------|------|-------------|
| product_id | Number | 商品 ID |
| layout | String | 布局（left, right） |
| show_description | Boolean | 显示描述 |
| show_price | Boolean | 显示价格 |
| show_add_to_cart | Boolean | 显示加购按钮 |

#### Collection Products (collection_products)

| Config Field | Type | Description |
|--------------|------|-------------|
| collection_id | Number | 分类 ID |
| layout | String | 布局 |
| columns | Number | 列数 |
| limit | Number | 显示数量 |
| show_view_all | Boolean | 显示「查看全部」按钮 |

### Navigation Blocks

#### Category Navigation (category_nav)

| Config Field | Type | Description |
|--------------|------|-------------|
| title | String | 标题 |
| layout | String | 布局（grid, list） |
| columns | Number | 列数 |
| categories | Array | 分类列表 |
| show_images | Boolean | 显示分类图片 |
| show_product_count | Boolean | 显示商品数量 |

#### Menu Links (menu_links)

| Config Field | Type | Description |
|--------------|------|-------------|
| title | String | 标题 |
| links | Array | 链接列表 |
| links[].label | String | 链接文本 |
| links[].url | String | 链接地址 |
| links[].icon | String | 图标 |

#### Search Bar (search_bar)

| Config Field | Type | Description |
|--------------|------|-------------|
| placeholder | String | 占位文本 |
| show_category_filter | Boolean | 显示分类筛选 |
| popular_searches | Array | 热门搜索词 |

### Content Blocks

#### Rich Text (rich_text)

| Config Field | Type | Description |
|--------------|------|-------------|
| content | String | HTML 内容 |

#### Image + Text (image_text)

| Config Field | Type | Description |
|--------------|------|-------------|
| layout | String | 布局（left, right） |
| image | String | 图片 URL |
| title | String | 标题 |
| description | String | 描述 |
| link | String | 链接 |

#### Video Player (video_player)

| Config Field | Type | Description |
|--------------|------|-------------|
| video_url | String | 视频 URL |
| autoplay | Boolean | 自动播放 |
| controls | Boolean | 显示控制条 |
| poster | String | 封面图 |

#### Custom HTML (custom_html)

| Config Field | Type | Description |
|--------------|------|-------------|
| html | String | HTML 代码 |
| css | String | CSS 样式 |

### Marketing Blocks

#### Countdown Timer (countdown_timer)

| Config Field | Type | Description |
|--------------|------|-------------|
| title | String | 标题 |
| end_time | String | 结束时间（ISO 8601） |
| show_days | Boolean | 显示天数 |
| show_hours | Boolean | 显示小时 |
| show_minutes | Boolean | 显示分钟 |
| style | String | 样式 |
| background_color | String | 背景色 |

#### Promotional Banner (promo_banner)

| Config Field | Type | Description |
|--------------|------|-------------|
| title | String | 标题 |
| description | String | 描述 |
| button_text | String | 按钮文本 |
| button_link | String | 按钮链接 |
| background_color | String | 背景色 |
| text_color | String | 文字颜色 |

#### Coupon Display (coupon_display)

| Config Field | Type | Description |
|--------------|------|-------------|
| coupon_codes | Array | 优惠券代码列表 |
| layout | String | 布局（grid, carousel） |
| show_expiry | Boolean | 显示有效期 |

### Social Blocks

#### Social Links (social_links)

| Config Field | Type | Description |
|--------------|------|-------------|
| platforms | Array | 平台列表 |
| platforms[].platform | String | 平台名称（facebook, instagram, twitter, etc.） |
| platforms[].url | String | 链接地址 |
| platforms[].icon | String | 图标 |
| layout | String | 布局（horizontal, vertical） |

#### Instagram Feed (instagram_feed)

| Config Field | Type | Description |
|--------------|------|-------------|
| account | String | 账号名 |
| post_count | Number | 显示数量 |
| columns | Number | 列数 |

---

## Theme Configuration Schema

### 预设主题

| Code | Name | Description | Default Colors |
|------|------|-------------|----------------|
| classic | Classic | 经典简洁，侧边栏导航 | Primary: #3B82F6, Secondary: #1E40AF |
| modern | Modern | 现代全宽，极简设计 | Primary: #10B981, Secondary: #059669 |
| minimal | Minimal | 极简白色，文字为主 | Primary: #000000, Secondary: #6B7280 |
| bold | Bold | 深色主题，鲜艳配色 | Primary: #8B5CF6, Secondary: #6D28D9 |
| nature | Nature | 自然风格，大地色系 | Primary: #059669, Secondary: #047857 |

### 可配置项

```json
{
  "fields": [
    {
      "key": "primary_color",
      "label": "Primary Color",
      "type": "color",
      "default": "#3B82F6"
    },
    {
      "key": "secondary_color",
      "label": "Secondary Color",
      "type": "color",
      "default": "#1E40AF"
    },
    {
      "key": "font_family",
      "label": "Font Family",
      "type": "select",
      "options": [
        {"value": "inter", "label": "Inter"},
        {"value": "roboto", "label": "Roboto"},
        {"value": "open-sans", "label": "Open Sans"},
        {"value": "lato", "label": "Lato"},
        {"value": "poppins", "label": "Poppins"}
      ],
      "default": "inter"
    },
    {
      "key": "button_style",
      "label": "Button Style",
      "type": "select",
      "options": [
        {"value": "rounded", "label": "Rounded"},
        {"value": "square", "label": "Square"},
        {"value": "pill", "label": "Pill"}
      ],
      "default": "rounded"
    }
  ]
}
```

---

## Milestones and Sequencing

### Phase 1: Backend Foundation (Week 1-2)

| Task | Duration | Description |
|------|----------|-------------|
| Database schema | 0.5 day | Create tables: themes, pages, decorations, page_versions, seo_configs |
| API definitions | 0.5 day | Create storefront.api with all endpoints |
| Domain entities | 1 day | Implement domain entities and repository interfaces |
| Repositories | 1.5 days | Implement GORM repositories |
| Business logic | 2 days | Implement logic for themes, pages, decorations, SEO |
| Unit tests | 1 day | Unit tests for business logic |
| Integration tests | 0.5 day | API integration tests |

### Phase 2: Frontend Builder (Week 2-3)

| Task | Duration | Description |
|------|----------|-------------|
| API client | 0.5 day | Create storefront.ts API client |
| Pinia store | 0.5 day | Create decoration.ts store |
| Block library | 1 day | Create BlockLibrary.vue with all block types |
| Builder canvas | 1.5 days | Create BuilderCanvas.vue with drag-drop |
| Block editor | 1.5 days | Create BlockEditor.vue forms for each block |
| Preview panel | 1 day | Create PreviewPanel.vue |
| Version history | 0.5 day | Create VersionHistory.vue |

### Phase 3: Theme & SEO (Week 3-4)

| Task | Duration | Description |
|------|----------|-------------|
| Theme management page | 1 day | Create themes/index.vue |
| Theme customization | 0.5 day | Color pickers, font selector |
| SEO settings page | 1 day | Create seo/index.vue |
| SEO preview | 0.5 day | Google search result preview |
| Router & menu | 0.5 day | Add routes and menu items |

### Phase 4: Polish & Testing (Week 4)

| Task | Duration | Description |
|------|----------|-------------|
| End-to-end testing | 1 day | Test complete workflows |
| Performance optimization | 0.5 day | Optimize preview rendering |
| Bug fixes | 1 day | Fix identified issues |
| Documentation | 0.5 day | API documentation |

**Total Estimate: 4 weeks**

---

## Appendix

### Block Type Mapping

| Block Type | Category | Priority |
|------------|----------|----------|
| banner | Layout | P0 |
| image_wall | Layout | P1 |
| divider | Layout | P1 |
| container | Layout | P2 |
| product_grid | Product | P0 |
| product_carousel | Product | P0 |
| featured_product | Product | P1 |
| collection_products | Product | P1 |
| category_nav | Navigation | P0 |
| menu_links | Navigation | P1 |
| search_bar | Navigation | P1 |
| rich_text | Content | P0 |
| image_text | Content | P1 |
| video_player | Content | P2 |
| custom_html | Content | P2 |
| countdown_timer | Marketing | P1 |
| promo_banner | Marketing | P1 |
| coupon_display | Marketing | P2 |
| social_links | Social | P2 |
| instagram_feed | Social | P2 |

### Page Types

| Page Type | Code | Description | SEO Available |
|-----------|------|-------------|---------------|
| Home | home | 首页 | ✅ |
| Category | category | 分类页（Phase 2） | ✅ (Phase 2) |
| Product | product | 商品详情页（Phase 2） | ✅ (Phase 2) |
| Custom | custom | 自定义页面（Phase 2） | ✅ (Phase 2) |

### SEO Character Limits

| Field | Recommended Max | Warning Threshold |
|-------|-----------------|-------------------|
| Title | 60 characters | > 60 |
| Description | 160 characters | > 160 |
| Keywords | 500 characters | > 500 |