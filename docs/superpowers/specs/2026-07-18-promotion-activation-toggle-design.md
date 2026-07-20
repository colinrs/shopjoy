# 促销活动状态列 Switch 开关设计

**日期：** 2026-07-18
**状态：** 已批准

## 1. 概述

在「促销活动」TAB 的状态列加入 `el-switch` 开关，使其与「优惠券」TAB 的交互保持一致——管理员无需进入操作列就能一键启用 / 停用促销活动。

后端 API、handler、logic、application 层均已就绪（`POST /api/v1/promotions/:id/activate` 与 `POST /api/v1/promotions/:id/deactivate`），本次仅做前端展示与交互改造。

## 2. 现状分析

### 2.1 后端（无需改动）

| 资源 | 路径 |
|------|------|
| API 定义 | `admin/desc/promotion.api`（已含 Activate / Deactivate 路由） |
| Handler | `admin/internal/handler/promotions/activate_promotion_handler.go` <br> `admin/internal/handler/promotions/deactivate_promotion_handler.go` |
| Logic | `admin/internal/logic/promotions/activate_promotion_logic.go` <br> `admin/internal/logic/promotions/deactivate_promotion_logic.go` |
| App | `admin/internal/application/promotion/promotion_app.go`（`ActivatePromotion` / `DeactivatePromotion` 方法） |

后端状态语义：
- `ActivatePromotion` → `StatusActive`（进行中）
- `DeactivatePromotion` → `StatusPaused`（已暂停）
- `ActivatePromotion` 在 `EndAt < now` 时返回 `ErrPromotionExpired`

### 2.2 前端（已部分就绪）

| 资源 | 路径 | 状态 |
|------|------|------|
| API 调用 | `shop-admin/src/api/promotion.ts` `activatePromotion(id)` / `deactivatePromotion(id)` | ✅ 已存在 |
| 模板按钮 | `shop-admin/src/views/promotions/index.vue` 操作列 509-525 行 | ⚠️ 存在但隐藏（位于最右侧「操作」列，仅 200px） |
| 处理函数 | `handleActivatePromotion` / `handleDeactivatePromotion` | ✅ 已存在 |

### 2.3 问题

- 状态列只显示静态 `el-tag`，无法直接交互
- 启用 / 停用按钮藏在操作列最右端，文字链样式（`type="success" link`）容易被忽略
- 与优惠券 TAB 体验不一致：优惠券状态列直接是 `el-switch` 单击切换

## 3. 设计方案

### 3.1 状态列改造

**文件：** `shop-admin/src/views/promotions/index.vue`

将状态列由「静态 tag」改为「按状态分支渲染」：

| 行状态 | 渲染 |
|--------|------|
| `active` | `el-switch`（ON） |
| `pending` | `el-switch`（OFF） |
| `paused` | `el-switch`（OFF） |
| `ended` | `el-tag`（不可切换） |
| 其他未知状态 | `el-tag`（防御性回退） |

**列宽度：** 由 `100` 调整为 `120`，容纳 switch 的 inline prompt 文本。

**Switch 配置（与优惠券 TAB 字段一致）：**

```vue
<el-switch
  v-if="row.status === 'active' || row.status === 'pending' || row.status === 'paused'"
  :model-value="row.status === 'active'"
  :loading="promotionToggleLoading[row.id] === true"
  :active-text="$t('promotions.activatedStatus')"
  :inactive-text="$t('promotions.unactivated')"
  inline-prompt
  @change="(val: boolean) => handleTogglePromotionStatus(row, val)"
/>
<el-tag
  v-else
  :type="getPromoStatusType(row.status)"
  effect="light"
  size="small"
>
  {{ getPromoStatusText(row.status) }}
</el-tag>
```

### 3.2 操作列清理

**删除** 操作列中已有的「启用」「停用」两个 link 按钮（index.vue 509-525 行），保留「编辑」「删除」。

**理由：** 状态列 switch 已经是单点击切换入口，操作列再保留按钮会造成重复交互、认知负担。与优惠券 TAB 保持一致。

### 3.3 处理函数

新增 `handleTogglePromotionStatus`，复用 `handleToggleCouponStatus` 的乐观更新 + 失败回滚模式：

```typescript
const promotionToggleLoading = reactive<Record<string, boolean>>({})

const handleTogglePromotionStatus = async (row: Promotion, nextActive: boolean) => {
  const wasActive = row.status === 'active'
  // 乐观翻转，避免网络慢时开关回弹
  row.status = nextActive ? 'active' : 'paused'
  promotionToggleLoading[row.id] = true
  try {
    if (nextActive) {
      await activatePromotion(row.id)
      ElMessage.success(t('promotions.activateSuccess'))
    } else {
      await deactivatePromotion(row.id)
      ElMessage.success(t('promotions.deactivateSuccess'))
    }
    loadPromotions()
  } catch (error) {
    // 失败回滚，保证 UI 与服务端状态一致
    row.status = wasActive ? 'active' : 'paused'
    handleError(
      error,
      nextActive ? t('promotions.activatePromotionFailed') : t('promotions.deactivatePromotionFailed')
    )
  } finally {
    promotionToggleLoading[row.id] = false
  }
}
```

### 3.4 状态语义对齐

| 状态值 | Switch 位置 | 后端映射 |
|--------|------------|----------|
| `pending` | OFF | `activatePromotion` → `active` |
| `active` | ON | `deactivatePromotion` → `paused` |
| `paused` | OFF | `activatePromotion` → `active` |
| `ended` | tag（不可切换） | — |

注意：本地 `row.status` 在 `deactivatePromotion` 成功后被乐观地写成 `paused`，与后端 `DeactivatePromotion` 的实际行为一致（后端设置 `StatusPaused`）。

## 4. 边界情况与错误处理

| 场景 | 行为 |
|------|------|
| 后端返回 `ErrPromotionExpired`（活动已过 EndAt） | UI 状态回滚到原值，error 提示通过 `useErrorHandler` 展示 |
| 网络错误 / 接口 5xx | 同上，UI 回滚 |
| 用户对 `ended` 状态行操作 | 状态列渲染 tag 而非 switch，物理上无法触发切换 |
| 连续快速点击 | `:loading` 绑到 `promotionToggleLoading[row.id]`，防止重入 |
| 多 Tab 并发刷新 | 切换开关后调用 `loadPromotions()` 重新拉取列表 |

## 5. 国际化

无需新增键，复用现有：

| 用途 | 键 |
|------|-----|
| Switch ON 文案 | `promotions.activatedStatus`（"已激活"） |
| Switch OFF 文案 | `promotions.unactivated`（"未激活"） |
| 启用成功提示 | `promotions.activateSuccess` |
| 停用成功提示 | `promotions.deactivateSuccess` |
| 启用失败提示 | `promotions.activatePromotionFailed` |
| 停用失败提示 | `promotions.deactivatePromotionFailed` |

> 文案键与优惠券 TAB 完全相同，确保两边交互语言统一。

## 6. 验收标准

- [ ] 状态为 `active` 的行显示为开启的 switch
- [ ] 状态为 `pending` / `paused` 的行显示为关闭的 switch
- [ ] 状态为 `ended` 的行显示 tag，不可点击
- [ ] 关闭 switch → 调用 `deactivatePromotion`，成功后状态变为 `paused`
- [ ] 开启 switch → 调用 `activatePromotion`，成功后状态变为 `active`
- [ ] 后端返回错误时 UI 状态自动回滚
- [ ] 操作列不再显示「启用」「停用」按钮
- [ ] 与优惠券 TAB 的状态列交互保持一致
- [ ] `make build` 通过
- [ ] 手动验证三种状态（pending / active / paused）的切换

## 7. 改动文件清单

| 文件 | 改动类型 | 说明 |
|------|----------|------|
| `shop-admin/src/views/promotions/index.vue` | 修改 | 状态列 switch 化；操作列删除两个按钮；新增 `handleTogglePromotionStatus` 与 `promotionToggleLoading` |
| `shop-admin/src/locales/zh.json` | 无 | 复用现有键 |
| `shop-admin/src/locales/en.json` | 无 | 复用现有键 |
| `admin/...` | 无 | 后端不动 |
| `shop-admin/src/api/promotion.ts` | 无 | API 已有 |
