# SKU 搜索下拉框改造设计

## 背景

库存管理页面（`/inventory`）中多处需要用户输入 SKU 编码：

1. **库存日志 Tab**：顶部筛选条件的 SKU 编码输入框
2. **库存调整 Tab**：表单中的 SKU 编码输入框
3. **批量安全库存弹窗**：多选 SKU 编码选择器

当前实现均为普通输入框或基于低库存列表的下拉框，体验较差且容易输错。本设计将这三处统一改造为「远程搜索 + 下拉选择」的形式，并新增后端 SKU 搜索接口支撑。

## 目标

- 用户在库存相关页面选择 SKU 时，通过搜索下拉框快速定位
- 下拉选项仅展示 SKU 编码，保持简洁
- 批量安全库存弹窗的已选 SKU 明细表保留商品名称和当前安全库存
- 默认状态下展示全部数据（不因为未选 SKU 而空白）

## 方案概述

采用「可复用 SKU 搜索组件」方案：

1. 后端新增通用 SKU 搜索接口
2. 前端封装 `SkuSearchSelect.vue` 组件
3. 在库存日志、库存调整、批量安全库存三处复用该组件

## 后端设计

### 新增接口

```
GET /api/v1/skus/search?keyword={keyword}&page={page}&page_size={page_size}
```

### 请求参数

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 否 | SKU 编码子串搜索关键词 |
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页条数，默认 20，最大 50 |

### 响应结构

```json
{
  "list": [
    {
      "sku_code": "SKU-RED-001",
      "product_id": "10001",
      "product_name": "示例商品",
      "safety_stock": 50
    }
  ],
  "total": 128
}
```

### 实现要点

- 在 `admin/desc/product.api` 中新增 `SearchSKUsReq` / `SearchSKUsResp` / `SearchSKUs` 路由
- 在 `admin/internal/logic/products/` 新增 `search_s_k_us_logic.go`
- 在 `admin/internal/infrastructure/persistence/sku_repository.go` 新增 `Search` 方法
- 查询条件：
  - `tenant_id` 按当前租户过滤
  - `status = 1`（仅启用状态 SKU）
  - `code LIKE '%{keyword}%'`（子串匹配）
- 关联 `products` 表获取 `product_name`
- 按 `created_at DESC` 排序

## 前端设计

### 新增组件 `SkuSearchSelect.vue`

位置：`shop-admin/src/components/SkuSearchSelect.vue`

#### Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| modelValue | string \| string[] | - | 选中值，支持单选或多选 |
| multiple | boolean | false | 是否多选 |
| placeholder | string | - | 占位文案 |
| clearable | boolean | true | 是否可清空 |
| style | string | - | 自定义样式 |

#### 行为

- 基于 Element Plus `el-select`（`filterable` + `remote`）
- `remote-method` 触发远程搜索，防抖 300ms
- 默认 `page_size=20`
- 关键词为空时不触发搜索，输入至少 1 个字符后才开始请求
- 下拉选项仅显示 `sku_code`
- 组件内部维护候选列表和 loading 状态
- 向外抛出 `update:modelValue` 和 `change` 事件

### 改造点

#### 1. 库存日志 Tab

文件：`shop-admin/src/views/inventory/index.vue`

- 将顶部筛选区的 `el-input v-model="logFilter.sku_code"` 替换为：

```vue
<SkuSearchSelect
  v-model="logFilter.sku_code"
  :placeholder="$t('inventory.skuCode')"
/>
```

- 未选择时 `logFilter.sku_code` 为空字符串，调用 `getInventoryLogs` 不传入 `sku_code`，展示全部日志

#### 2. 库存调整 Tab

文件：`shop-admin/src/views/inventory/index.vue`

- 将表单中的 `el-input v-model="adjustForm.sku_code"` 替换为：

```vue
<SkuSearchSelect
  v-model="adjustForm.sku_code"
  :placeholder="$t('inventory.skuCode')"
/>
```

#### 3. 批量安全库存弹窗

文件：`shop-admin/src/views/inventory/index.vue`

- 将弹窗中的多选 `el-select`（选项来自 `lowStockList`）替换为：

```vue
<SkuSearchSelect
  v-model="batchSafetyStockForm.sku_codes"
  multiple
  :placeholder="$t('inventory.enterSKUCode')"
  style="width: 100%"
/>
```

- 已选 SKU 明细表不再依赖 `lowStockList`，改为维护一个 `selectedSkuDetails` 映射
- 当用户选择 SKU 时，从搜索返回结果中提取 `product_name` 和 `safety_stock` 填充明细表
- 低库存 SKU 列表仍保留在页面顶部作为预警展示，但不再作为批量安全库存的选择来源

### API 层新增

文件：`shop-admin/src/api/product.ts`

新增 `searchSKUs` 函数：

```typescript
export interface SearchSKUItem {
  sku_code: string
  product_id: string
  product_name: string
  safety_stock: number
}

export function searchSKUs(params: { keyword?: string; page?: number; page_size?: number }) {
  return request<{ list: SearchSKUItem[]; total: number }>({
    url: '/api/v1/skus/search',
    method: 'get',
    params
  })
}
```

## 数据流

```
用户输入关键词
  └─> SkuSearchSelect 防抖 300ms
      └─> searchSKUs(keyword)
          └─> GET /api/v1/skus/search
              └─> 后端按 SKU 编码子串搜索启用状态 SKU
                  └─> 返回 sku_code / product_id / product_name / safety_stock
                      └─> 下拉框展示 sku_code
                          └─> 用户选择后回填到对应表单字段
```

## 错误处理

- 搜索接口失败：使用 `useErrorHandler` 统一提示，候选列表保持上一次结果或为空
- 关键词为空：组件可重新加载默认列表（可选，或保持空列表等待输入）
- 未选 SKU 提交批量更新：保持现有表单校验提示

## 测试要点

- 后端：
  - 子串搜索正确返回匹配的 SKU
  - 租户隔离生效
  - 禁用状态 SKU 不出现在结果中
  - 分页参数生效
- 前端：
  - 三处下拉框能正常搜索和选择
  - 批量安全库存弹窗已选列表展示正确
  - 未选择 SKU 时，库存日志 Tab 默认展示全部数据
  - 清空选择后恢复默认展示

## 影响范围

| 文件 | 变更 |
|------|------|
| `admin/desc/product.api` | 新增搜索接口定义 |
| `admin/internal/types/types.go` | 重新生成类型（自动） |
| `admin/internal/handler/products/search_s_k_us_handler.go` | 新增 handler（自动） |
| `admin/internal/logic/products/search_s_k_us_logic.go` | 新增业务逻辑 |
| `admin/internal/infrastructure/persistence/sku_repository.go` | 新增 Search 方法 |
| `shop-admin/src/api/product.ts` | 新增 searchSKUs |
| `shop-admin/src/components/SkuSearchSelect.vue` | 新增组件 |
| `shop-admin/src/views/inventory/index.vue` | 替换 3 处 SKU 输入/选择 |

## 备注

- 本设计遵循项目现有模式（参考 `CouponSelector.vue` 的远程搜索实现）
- 后端 SKU 搜索接口采用子串匹配，未来 SKU 数量极大时可考虑接入全文检索或前缀索引优化
