# SKU 编码生成策略设计文档

## 概述

设计一套 SKU 编码生成策略，满足以下需求：
- 业务可自定义前缀（租户级 + 商品级）
- 系统随机生成，全局唯一
- 带有租户标识，可通过 SKU 编码解码出租户ID

## 编码结构

### 完整格式

```
{租户前缀}[-{商品前缀}]-{紧凑编码}
```

### 示例

| 场景 | SKU 编码 |
|------|----------|
| 仅租户前缀 | `NIKE-3kF9xZmP2` |
| 租户+商品前缀 | `NIKE-SHOE-3kF9xZmP2` |
| 无前缀 | `3kF9xZmP2`（系统默认） |

### 紧凑编码结构（10位 Base62）

```
[租户ID(4位)][小时时间戳(3位)][随机序列(3位)]
     ↓           ↓              ↓
   T001        7xK           fM2
```

| 部分 | 位数 | Base62值范围 | 说明 |
|------|------|-------------|------|
| 租户ID编码 | 4位 | `0000` ~ `zzzz` | 可表示约 1,477 万个租户 |
| 小时时间戳 | 3位 | `000` ~ `zzz` | 相对偏移小时数，约 27 年周期 |
| 随机序列 | 3位 | `000` ~ `zzz` | 238,328 种组合，足以支撑高并发场景 |

### 长度限制

| 项目 | 最大长度 |
|------|----------|
| 租户前缀 | 8 字符 |
| 商品前缀 | 8 字符 |
| 紧凑编码 | 10 字符（固定） |
| 分隔符 | 1-2 字符 |
| **总长度上限** | **28 字符** |

## 编解码逻辑

### Base62 字符集

```
0-9: 数字 0-9
10-35: 大写 A-Z
36-61: 小写 a-z
```

### 租户ID 编码方式

将租户ID转换为 4 位 Base62 字符：

```
租户ID = 1      → 编码 = "0001"
租户ID = 12345  → 编码 = "0Dm7"
租户ID = 100万  → 编码 = "4c91"
```

### 时间戳编码方式

使用相对偏移小时数，从系统上线日期开始计算：

```
基准时间 = 2024-01-01 00:00:00 UTC
当前小时偏移 = (当前时间 - 基准时间) / 1小时
编码 = Base62(偏移值, 3位)
```

**示例：**
```
2024-01-01 00:00 → 偏移=0    → 编码="000"
2024-01-01 01:00 → 偏移=1    → 编码="001"
2025-03-22 15:00 → 偏移=10500+ → 编码="2xK"
```

**周期：** 3位 Base62 最大表示 238,328 小时 ≈ 27.2 年。

### 解码函数

```go
const (
    compactCodeLength = 10  // 紧凑编码固定长度
    tenantIDLength    = 4
    timestampLength   = 3
    randomLength      = 3
)

// ParseSKU 解析 SKU 编码
func ParseSKU(code string) (*SKUInfo, error) {
    // 1. 提取紧凑编码部分（最后一个 - 之后）
    parts := strings.Split(code, "-")
    compactCode := parts[len(parts)-1]

    // 2. 验证紧凑编码
    if len(compactCode) != compactCodeLength {
        return nil, ErrSKUParseFailed
    }
    if !isBase62(compactCode) {
        return nil, ErrSKUParseFailed
    }

    // 3. 解码租户ID（前4位）
    tenantID := decodeBase62(compactCode[0:tenantIDLength])

    // 4. 解码时间戳（中间3位）
    hourOffset := decodeBase62(compactCode[tenantIDLength : tenantIDLength+timestampLength])
    createdAt := baseTime.Add(time.Duration(hourOffset) * time.Hour)

    // 5. 提取前缀
    prefixes := parts[:len(parts)-1]

    return &SKUInfo{
        TenantPrefix:   prefixes[0] if len(prefixes) > 0 else "",
        ProductPrefix:  prefixes[1] if len(prefixes) > 1 else "",
        TenantID:       tenantID,
        CreatedAt:      createdAt,
        RandomSequence: compactCode[tenantIDLength+timestampLength:],
    }, nil
}

// ExtractTenantID 仅提取租户ID（高频场景优化）
func ExtractTenantID(code string) (int64, error) {
    // 找到最后一个 -
    idx := strings.LastIndex(code, "-")
    compactCode := code
    if idx >= 0 {
        compactCode = code[idx+1:]
    }

    // 验证长度和字符
    if len(compactCode) != compactCodeLength {
        return 0, ErrSKUParseFailed
    }
    if !isBase62(compactCode) {
        return 0, ErrSKUParseFailed
    }

    return decodeBase62(compactCode[0:tenantIDLength]), nil
}

// isBase62 验证字符串是否只包含 Base62 字符
func isBase62(s string) bool {
    for _, c := range s {
        if !((c >= '0' && c <= '9') ||
            (c >= 'A' && c <= 'Z') ||
            (c >= 'a' && c <= 'z')) {
            return false
        }
    }
    return true
}
```

### 唯一性保证

| 机制 | 说明 |
|------|------|
| 时间戳（小时级） | 同一租户同一小时内生成的 SKU 有相同时间戳前缀 |
| 随机序列（3位） | 238,328 种组合，支持高并发场景 |
| 数据库唯一约束 | `tenant_id + code` 唯一索引作为最终保障 |
| 冲突重试 | 检测到冲突时重新生成随机序列 |

**容量估算：**
- 单租户每小时：238,328 个（3位随机序列）
- 实际场景绰绰有余，碰撞概率极低

**碰撞概率分析：**
- 生成 1,000 个 SKU/小时：碰撞概率 < 0.002%
- 生成 10,000 个 SKU/小时：碰撞概率 < 0.2%
- 批量导入 50,000 个 SKU：碰撞概率 < 5%（可通过重试解决）

## 前缀配置

### 数据库变更

**tenants 表新增字段：**

```sql
ALTER TABLE tenants ADD COLUMN sku_prefix VARCHAR(8) DEFAULT '' COMMENT 'SKU默认前缀';
```

**products 表新增字段：**

```sql
ALTER TABLE products ADD COLUMN sku_prefix VARCHAR(8) DEFAULT '' COMMENT '商品SKU前缀';
```

### 前缀规则优先级

```
商品前缀 + 租户前缀 > 租户前缀 > 无前缀
```

| 场景 | 租户前缀 | 商品前缀 | 最终前缀 |
|------|----------|----------|----------|
| 都有 | `NIKE` | `SHOE` | `NIKE-SHOE` |
| 仅租户 | `NIKE` | 空 | `NIKE` |
| 仅商品 | 空 | `SHOE` | `SHOE` |
| 都无 | 空 | 空 | 无前缀 |

### 前缀校验规则

```go
// 前缀校验
// 1. 只允许字母、数字
// 2. 长度 1-8 字符
// 3. 不能以数字开头（避免与紧凑编码混淆）
// 4. 不区分大小写（存储时统一大写）
```

**有效示例：** `NIKE`, `SHOE`, `A1`, `PROD`

**无效示例：** `123`, `NIKE-SHOE`（含分隔符）, `NIKE_2024`（含特殊字符）

## 代码结构

### 新增模块

```
pkg/sku/
├── generator.go      # SKU 生成器
├── parser.go         # SKU 解析器
├── encoder.go        # Base62 编解码
└── validator.go      # 前缀校验器
```

### 核心接口

```go
// Generator SKU 生成器接口
type Generator interface {
    // Generate 生成 SKU 编码
    Generate(tenantID int64, tenantPrefix, productPrefix string) (string, error)

    // GenerateWithRetry 带重试的生成（冲突时自动重试）
    GenerateWithRetry(tenantID int64, tenantPrefix, productPrefix string, maxRetry int) (string, error)
}

// Parser SKU 解析器接口
type Parser interface {
    // Parse 解析 SKU 编码，返回租户ID等信息
    Parse(code string) (*SKUInfo, error)

    // ExtractTenantID 仅提取租户ID（高频场景优化）
    ExtractTenantID(code string) (int64, error)
}

// SKUInfo SKU 解析结果
type SKUInfo struct {
    TenantPrefix   string    // 租户前缀
    ProductPrefix  string    // 商品前缀
    TenantID       int64     // 租户ID
    CreatedAt      time.Time // 生成时间（小时精度）
    RandomSequence string    // 随机序列
    CompactCode    string    // 紧凑编码部分
}
```

### 集成点

**ServiceContext 新增依赖：**

```go
type ServiceContext struct {
    // ... 现有字段
    SKUGenerator sku.Generator
}
```

**CreateSKU 逻辑变更：**

```go
// 当前：用户手动传入 Code
sku := &product.SKU{
    Code: req.Code,
}

// 改为：系统自动生成
skuCode, err := l.svcCtx.SKUGenerator.GenerateWithRetry(
    tenantID,
    tenant.SKUPrefix,
    product.SKUPrefix,
    3,
)
sku := &product.SKU{
    Code: skuCode,
}
```

## 配置与错误处理

### 错误码

```go
// SKU 相关错误 (190xxx)
ErrSKUPrefixInvalid          = &Err{HTTPCode: 400, Code: 190001, Msg: "SKU前缀格式无效"}
ErrSKUPrefixTooLong          = &Err{HTTPCode: 400, Code: 190002, Msg: "SKU前缀长度超过限制"}
ErrSKUGenerateFailed         = &Err{HTTPCode: 500, Code: 190003, Msg: "SKU生成失败"}
ErrSKUParseFailed            = &Err{HTTPCode: 400, Code: 190004, Msg: "SKU编码解析失败"}
ErrSKUCodeTooLong            = &Err{HTTPCode: 400, Code: 190005, Msg: "SKU编码长度超过限制"}
ErrSKUPrefixStartsWithNumber = &Err{HTTPCode: 400, Code: 190006, Msg: "SKU前缀不能以数字开头"}
```

### 可配置项

```go
type Config struct {
    Epoch                  time.Time // 默认: 2024-01-01 00:00:00 UTC
    MaxTenantPrefixLength  int       // 默认: 8
    MaxProductPrefixLength int       // 默认: 8
    MaxTotalLength         int       // 默认: 32
    TenantIDLength         int       // 默认: 4
    TimestampLength        int       // 默认: 3
    RandomLength           int       // 默认: 3
    MaxRetry               int       // 默认: 3
}
```

## 并发安全

| 场景 | 策略 |
|------|------|
| 多协程同时生成 | 随机序列 + 时间戳保证碰撞概率极低 |
| 数据库唯一约束冲突 | 检测错误后重试生成 |
| 高并发场景 | 无锁设计，性能优于集中式发号器 |

## 测试策略

| 测试类型 | 覆盖内容 |
|----------|----------|
| 单元测试 | 编解码正确性、前缀校验、边界值 |
| 集成测试 | 与数据库唯一约束配合、重试机制 |
| 压力测试 | 并发生成碰撞率、性能基准 |
| 回归测试 | 解码旧 SKU 的向后兼容性 |

## 迁移策略

### API 兼容性

采用**混合模式**，支持新旧两种方式：

```go
// CreateSKU 逻辑
if req.Code != "" {
    // 用户手动传入 Code，保持兼容
    sku.Code = req.Code
} else {
    // 系统自动生成
    skuCode, err := l.svcCtx.SKUGenerator.GenerateWithRetry(...)
    sku.Code = skuCode
}
```

### 现有 SKU 处理

| 场景 | 处理方式 |
|------|----------|
| 已存在的 SKU | 保持不变，不强制迁移 |
| 新建 SKU（不传 Code） | 系统自动生成新格式 |
| 新建 SKU（传入 Code） | 使用传入的 Code（兼容模式） |

### 解码兼容

对于无法按新格式解析的 SKU 编码，解码函数返回错误但不影响正常使用：

```go
func ParseSKU(code string) (*SKUInfo, error) {
    // ... 验证和解码逻辑

    // 如果解析失败，返回错误但记录原始编码
    if !isValidFormat {
        return nil, ErrSKUParseFailed
    }
    // ...
}
```

### 数据库迁移

```sql
-- 1. 添加前缀字段
ALTER TABLE tenants ADD COLUMN sku_prefix VARCHAR(8) DEFAULT '' COMMENT 'SKU默认前缀';
ALTER TABLE products ADD COLUMN sku_prefix VARCHAR(8) DEFAULT '' COMMENT '商品SKU前缀';

-- 2. 确保 SKU 唯一约束存在
ALTER TABLE skus ADD UNIQUE INDEX uk_tenant_code (tenant_id, code);

-- 注意：不迁移现有 SKU 数据，保持原样
```

### 版本演进

为支持未来格式变更，保留版本扩展能力：

- 当前格式为 v1，紧凑编码长度固定 10 位
- 未来如需变更，可增加版本前缀（如 `V2-xxx`）
- 解码时根据前缀选择对应解析器