-- ============================================
-- 分类表 (categories)
-- ============================================

CREATE TABLE IF NOT EXISTS `categories` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '分类ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `parent_id` BIGINT NOT NULL DEFAULT 0 COMMENT '父分类ID',
    `name` VARCHAR(100) NOT NULL COMMENT '分类名称',
    `code` VARCHAR(100) DEFAULT '' COMMENT '分类代码',
    `level` TINYINT NOT NULL DEFAULT 1 COMMENT '层级',
    `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `icon` VARCHAR(255) DEFAULT '' COMMENT '图标',
    `image` VARCHAR(500) DEFAULT '' COMMENT '图片',
    `seo_title` VARCHAR(200) DEFAULT '' COMMENT 'SEO标题',
    `seo_description` VARCHAR(500) DEFAULT '' COMMENT 'SEO描述',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分类表';

-- ============================================
-- 品牌表 (brands)
-- ============================================

CREATE TABLE IF NOT EXISTS `brands` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '品牌ID',
    `tenant_id` BIGINT NOT NULL COMMENT '租户ID',
    `name` VARCHAR(100) NOT NULL COMMENT '品牌名称',
    `logo` VARCHAR(500) DEFAULT '' COMMENT 'Logo URL',
    `description` TEXT COMMENT '描述',
    `website` VARCHAR(255) DEFAULT '' COMMENT '官网',
    `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `enable_page` TINYINT NOT NULL DEFAULT 0 COMMENT '是否启用品牌专区',
    `trademark_number` VARCHAR(100) DEFAULT '' COMMENT '商标号',
    `trademark_country` VARCHAR(10) DEFAULT '' COMMENT '商标注册国家',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_name` (`name`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌表';

-- ============================================
-- 商品表 (products)
-- ============================================

CREATE TABLE IF NOT EXISTS `products` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '商品ID',
    `tenant_id` BIGINT NOT NULL DEFAULT 0 COMMENT '租户ID',
    `sku` VARCHAR(64) DEFAULT '' COMMENT 'SKU代码',
    `name` VARCHAR(200) NOT NULL COMMENT '商品名称',
    `description` TEXT COMMENT '商品描述',
    `price` BIGINT NOT NULL DEFAULT 0 COMMENT '售价(分)',
    `cost_price` BIGINT NOT NULL DEFAULT 0 COMMENT '成本价(分)',
    `currency` VARCHAR(10) NOT NULL DEFAULT 'CNY' COMMENT '货币',
    `stock` INT NOT NULL DEFAULT 0 COMMENT '库存',
    `status` INT NOT NULL DEFAULT 0 COMMENT '状态: 0-草稿, 1-上架, 2-下架, 3-已删除',
    `category_id` BIGINT NOT NULL DEFAULT 0 COMMENT '分类ID',
    `brand` VARCHAR(64) DEFAULT '' COMMENT '品牌',
    `brand_id` BIGINT NULL COMMENT '品牌ID',
    `sku_prefix` VARCHAR(8) DEFAULT '' COMMENT '商品SKU前缀',
    `tags` JSON COMMENT '标签',
    `images` JSON COMMENT '图片列表',
    `is_matrix_product` TINYINT NOT NULL DEFAULT 0 COMMENT '是否有变体',
    `hs_code` VARCHAR(20) DEFAULT '' COMMENT 'HS编码',
    `coo` VARCHAR(10) DEFAULT '' COMMENT '原产国',
    `weight` DECIMAL(10,2) DEFAULT 0.00 COMMENT '重量',
    `weight_unit` VARCHAR(10) DEFAULT 'g' COMMENT '重量单位',
    `length` DECIMAL(10,2) DEFAULT 0.00 COMMENT '长度(cm)',
    `width` DECIMAL(10,2) DEFAULT 0.00 COMMENT '宽度(cm)',
    `height` DECIMAL(10,2) DEFAULT 0.00 COMMENT '高度(cm)',
    `dangerous_goods` JSON COMMENT '危险品标识',
    `deleted_at` BIGINT NULL COMMENT '删除时间',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_sku` (`sku`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_name` (`name`),
    KEY `idx_category_id` (`category_id`),
    KEY `idx_brand_id` (`brand_id`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品表';

-- ============================================
-- SKU表 (skus)
-- ============================================

CREATE TABLE IF NOT EXISTS `skus` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'SKU ID',
    `tenant_id` BIGINT NOT NULL DEFAULT 0 COMMENT '租户ID',
    `product_id` BIGINT NOT NULL COMMENT '商品ID',
    `code` VARCHAR(100) NOT NULL COMMENT 'SKU代码',
    `price_amount` BIGINT NOT NULL DEFAULT 0 COMMENT '价格(分)',
    `price_currency` VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',
    `stock` INT NOT NULL DEFAULT 0 COMMENT '库存',
    `available_stock` INT NOT NULL DEFAULT 0 COMMENT '可用库存',
    `locked_stock` INT NOT NULL DEFAULT 0 COMMENT '锁定库存',
    `safety_stock` INT NOT NULL DEFAULT 0 COMMENT '安全库存阈值',
    `presale_enabled` TINYINT NOT NULL DEFAULT 0 COMMENT '是否开启预售',
    `attributes` JSON COMMENT '属性',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    UNIQUE KEY `uk_tenant_code` (`tenant_id`, `code`),
    KEY `idx_product_id` (`product_id`),
    KEY `idx_status` (`status`),
    KEY `idx_low_stock_alert` (`tenant_id`, `status`, `safety_stock`, `available_stock`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='SKU表';

-- ============================================
-- 商品市场关联表 (product_markets)
-- ============================================

CREATE TABLE IF NOT EXISTS `product_markets` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `tenant_id` BIGINT NOT NULL DEFAULT 0 COMMENT '租户ID',
    `product_id` BIGINT NOT NULL COMMENT '商品ID',
    `variant_id` BIGINT DEFAULT NULL COMMENT '变体ID',
    `market_id` BIGINT NOT NULL COMMENT '市场ID',
    `is_enabled` TINYINT NOT NULL DEFAULT 0 COMMENT '是否启用',
    `status_override` INT DEFAULT NULL COMMENT '状态覆盖',
    `price` DECIMAL(10,2) DEFAULT 0.00 COMMENT '市场专属价格',
    `compare_at_price` DECIMAL(10,2) DEFAULT NULL COMMENT '对比价格',
    `stock_alert_threshold` INT NOT NULL DEFAULT 0 COMMENT '库存预警阈值',
    `published_at` BIGINT DEFAULT NULL COMMENT '发布时间',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_product_id` (`product_id`),
    KEY `idx_market_id` (`market_id`),
    KEY `idx_variant_id` (`variant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品市场关联表';

-- ============================================
-- 商品多语言表 (product_localizations)
-- ============================================

CREATE TABLE IF NOT EXISTS `product_localizations` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL,
    `product_id` BIGINT NOT NULL,
    `language_code` VARCHAR(10) NOT NULL COMMENT '语言代码: en, zh-CN, ja',
    `name` VARCHAR(200) DEFAULT '' COMMENT '产品名称',
    `description` TEXT COMMENT '产品描述',
    `created_at` BIGINT NOT NULL DEFAULT 0,
    `updated_at` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_product_id` (`product_id`),
    UNIQUE KEY `idx_tenant_product_language` (`tenant_id`, `product_id`, `language_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='产品多语言表';

-- ============================================
-- 分类市场可见性表 (category_markets)
-- ============================================

CREATE TABLE IF NOT EXISTS `category_markets` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL,
    `category_id` BIGINT NOT NULL,
    `market_id` BIGINT NOT NULL,
    `is_visible` TINYINT NOT NULL DEFAULT 1 COMMENT '是否可见',
    `created_at` BIGINT NOT NULL DEFAULT 0,
    `updated_at` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_category_id` (`category_id`),
    KEY `idx_market_id` (`market_id`),
    UNIQUE KEY `idx_tenant_category_market` (`tenant_id`, `category_id`, `market_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分类市场可见性';

-- ============================================
-- 品牌市场可见性表 (brand_markets)
-- ============================================

CREATE TABLE IF NOT EXISTS `brand_markets` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL,
    `brand_id` BIGINT NOT NULL,
    `market_id` BIGINT NOT NULL,
    `is_visible` TINYINT NOT NULL DEFAULT 1 COMMENT '是否可见',
    `created_at` BIGINT NOT NULL DEFAULT 0,
    `updated_at` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_brand_id` (`brand_id`),
    KEY `idx_market_id` (`market_id`),
    UNIQUE KEY `idx_tenant_brand_market` (`tenant_id`, `brand_id`, `market_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌市场可见性';

-- ============================================
-- 仓库表 (warehouses)
-- ============================================

CREATE TABLE IF NOT EXISTS `warehouses` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL,
    `code` VARCHAR(50) NOT NULL COMMENT '仓库代码',
    `name` VARCHAR(100) NOT NULL COMMENT '仓库名称',
    `country` VARCHAR(10) DEFAULT '' COMMENT '所在国家',
    `address` VARCHAR(500) DEFAULT '' COMMENT '详细地址',
    `is_default` TINYINT NOT NULL DEFAULT 0 COMMENT '是否默认仓库',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
    `created_at` BIGINT NOT NULL DEFAULT 0,
    `updated_at` BIGINT NOT NULL DEFAULT 0,
    `deleted_at` BIGINT DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    UNIQUE KEY `idx_tenant_code` (`tenant_id`, `code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='仓库表';

-- ============================================
-- 仓库库存表 (warehouse_inventories)
-- ============================================

CREATE TABLE IF NOT EXISTS `warehouse_inventories` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL,
    `sku_code` VARCHAR(100) NOT NULL COMMENT 'SKU代码',
    `warehouse_id` BIGINT NOT NULL,
    `available_stock` INT NOT NULL DEFAULT 0,
    `locked_stock` INT NOT NULL DEFAULT 0,
    `created_at` BIGINT NOT NULL DEFAULT 0,
    `updated_at` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_sku_code` (`sku_code`),
    KEY `idx_warehouse_id` (`warehouse_id`),
    UNIQUE KEY `idx_tenant_sku_warehouse` (`tenant_id`, `sku_code`, `warehouse_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='仓库库存表';

-- ============================================
-- 库存变更日志表 (inventory_logs)
-- ============================================

CREATE TABLE IF NOT EXISTS `inventory_logs` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT NOT NULL,
    `sku_code` VARCHAR(100) NOT NULL,
    `product_id` BIGINT NOT NULL,
    `warehouse_id` BIGINT NOT NULL DEFAULT 0 COMMENT '0=汇总',
    `change_type` VARCHAR(30) NOT NULL COMMENT 'manual, order, return, adjustment',
    `change_quantity` INT NOT NULL COMMENT '正数增加，负数减少',
    `before_stock` INT NOT NULL,
    `after_stock` INT NOT NULL,
    `order_no` VARCHAR(50) DEFAULT '',
    `remark` VARCHAR(500) DEFAULT '',
    `operator_id` BIGINT NOT NULL,
    `created_at` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_sku_code` (`sku_code`),
    KEY `idx_product_id` (`product_id`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存变更日志';

-- ============================================
-- 市场表 (markets)
-- 跨境电商市场配置
-- ============================================

CREATE TABLE IF NOT EXISTS `markets` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '市场ID',
    `tenant_id` BIGINT NOT NULL DEFAULT 0 COMMENT '租户ID',
    `code` VARCHAR(10) NOT NULL COMMENT '市场代码: US, UK, DE, FR, AU',
    `name` VARCHAR(100) NOT NULL COMMENT '市场名称',
    `currency` VARCHAR(10) NOT NULL COMMENT '货币: USD, GBP, EUR, AUD',
    `default_language` VARCHAR(10) DEFAULT 'en' COMMENT '默认语言',
    `flag` VARCHAR(255) DEFAULT '' COMMENT '旗帜图标',
    `is_active` TINYINT NOT NULL DEFAULT 1 COMMENT '是否启用',
    `is_default` TINYINT NOT NULL DEFAULT 0 COMMENT '是否主市场',
    `tax_rules` JSON COMMENT '税务配置',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_code` (`tenant_id`, `code`),
    KEY `idx_code` (`code`),
    KEY `idx_is_active` (`is_active`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='市场表';

-- ============================================
-- 测试数据
-- ============================================

-- 分类数据 (Demo Shop)
INSERT INTO `categories` (`id`, `tenant_id`, `parent_id`, `name`, `code`, `level`, `sort`, `icon`, `image`, `status`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
(1, 1, 0, '服装', 'clothing', 1, 1, 'shirt', 'https://cdn.example.com/cat1.png', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(2, 1, 1, '男装', 'mens', 2, 1, '', '', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(3, 1, 1, '女装', 'womens', 2, 2, '', '', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(4, 1, 0, '电子产品', 'electronics', 1, 2, 'laptop', 'https://cdn.example.com/cat2.png', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(5, 1, 4, '手机配件', 'phone-accessories', 2, 1, '', '', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(6, 1, 4, '电脑配件', 'computer-accessories', 2, 2, '', '', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(7, 1, 0, '家居', 'home', 1, 3, 'home', 'https://cdn.example.com/cat3.png', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2);

-- 品牌数据 (Demo Shop)
INSERT INTO `brands` (`id`, `tenant_id`, `name`, `logo`, `description`, `website`, `sort`, `status`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
(1, 1, 'Nike', 'https://cdn.example.com/nike.png', 'Just Do It', 'https://www.nike.com', 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(2, 1, 'Adidas', 'https://cdn.example.com/adidas.png', 'Impossible is Nothing', 'https://www.adidas.com', 2, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(3, 1, 'Apple', 'https://cdn.example.com/apple.png', 'Think Different', 'https://www.apple.com', 3, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(4, 1, 'Samsung', 'https://cdn.example.com/samsung.png', 'Do What You Can''t', 'https://www.samsung.com', 4, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2);

-- 商品数据 (Demo Shop)
INSERT INTO `products` (`id`, `sku`, `name`, `description`, `price`, `cost_price`, `currency`, `stock`, `status`, `category_id`, `brand`, `tags`, `images`, `is_matrix_product`, `hs_code`, `coo`, `weight`, `weight_unit`, `length`, `width`, `height`, `dangerous_goods`, `created_at`, `updated_at`) VALUES
(1, 'SKU-001', 'Nike Air Max 270', 'Nike Air Max 270 运动鞋，舒适透气', 129900, 80000, 'CNY', 100, 1, 2, 'Nike', '["运动", "跑步", "休闲"]', '["https://cdn.example.com/p1-1.jpg", "https://cdn.example.com/p1-2.jpg"]', 1, '64041100', 'CN', 450.00, 'g', 28.00, 18.00, 12.00, '[]', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(2, 'SKU-002', 'Adidas Ultraboost 22', 'Adidas Ultraboost 22 跑步鞋', 159900, 95000, 'CNY', 80, 1, 2, 'Adidas', '["运动", "跑步"]', '["https://cdn.example.com/p2-1.jpg"]', 1, '64041100', 'VN', 380.00, 'g', 27.00, 17.00, 11.00, '[]', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(3, 'SKU-003', 'iPhone 15 手机壳', 'iPhone 15 硅胶保护壳', 9900, 3000, 'CNY', 500, 1, 5, 'Apple', '["手机配件", "保护壳"]', '["https://cdn.example.com/p3-1.jpg"]', 0, '39269010', 'CN', 35.00, 'g', 15.00, 8.00, 1.50, '[]', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(4, 'SKU-004', 'MacBook 充电器', 'MacBook Pro 16寸 充电器 140W', 79900, 45000, 'CNY', 50, 1, 6, 'Apple', '["电脑配件", "充电器"]', '["https://cdn.example.com/p4-1.jpg"]', 0, '85044014', 'CN', 480.00, 'g', 15.00, 15.00, 3.50, '[]', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(5, 'SKU-005', '简约台灯', '北欧简约风格LED台灯', 29900, 12000, 'CNY', 200, 1, 7, '', '["家居", "灯具"]', '["https://cdn.example.com/p5-1.jpg"]', 0, '94052100', 'CN', 850.00, 'g', 35.00, 15.00, 45.00, '[]', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- SKU数据 (商品变体)
INSERT INTO `skus` (`id`, `product_id`, `code`, `price_amount`, `price_currency`, `stock`, `attributes`, `status`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
(1, 1, 'SKU-001-BLK-42', 129900, 'CNY', 30, '{"颜色": "黑色", "尺码": "42"}', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(2, 1, 'SKU-001-BLK-43', 129900, 'CNY', 25, '{"颜色": "黑色", "尺码": "43"}', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(3, 1, 'SKU-001-WHT-42', 129900, 'CNY', 25, '{"颜色": "白色", "尺码": "42"}', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(4, 1, 'SKU-001-WHT-43', 129900, 'CNY', 20, '{"颜色": "白色", "尺码": "43"}', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(5, 2, 'SKU-002-BLK-42', 159900, 'CNY', 40, '{"颜色": "黑色", "尺码": "42"}', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(6, 2, 'SKU-002-BLK-43', 159900, 'CNY', 40, '{"颜色": "黑色", "尺码": "43"}', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(7, 3, 'SKU-003-BLK', 9900, 'CNY', 200, '{"颜色": "黑色"}', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(8, 3, 'SKU-003-WHT', 9900, 'CNY', 200, '{"颜色": "白色"}', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2),
(9, 3, 'SKU-003-BLU', 9900, 'CNY', 100, '{"颜色": "蓝色"}', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 2, 2);

-- ============================================
-- 测试数据
-- ============================================

-- 市场数据
INSERT INTO `markets` (`id`, `tenant_id`, `code`, `name`, `currency`, `default_language`, `flag`, `is_active`, `is_default`, `tax_rules`, `created_at`, `updated_at`) VALUES
-- Demo Shop 市场
(1, 1, 'CN', '中国大陆', 'CNY', 'zh-CN', '🇨🇳', 1, 1, '{"IncludeTax": true, "VATRate": 0.13}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(2, 1, 'US', '美国', 'USD', 'en', '🇺🇸', 1, 0, '{"IncludeTax": false, "GSTRate": 0}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(3, 1, 'UK', '英国', 'GBP', 'en-GB', '🇬🇧', 1, 0, '{"IncludeTax": true, "VATRate": 0.20, "IOSSEnabled": true}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(4, 1, 'DE', '德国', 'EUR', 'de', '🇩🇪', 1, 0, '{"IncludeTax": true, "VATRate": 0.19, "IOSSEnabled": true}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(5, 1, 'AU', '澳大利亚', 'AUD', 'en-AU', '🇦🇺', 1, 0, '{"IncludeTax": true, "GSTRate": 0.10}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- Test Store 市场
(6, 2, 'CN', '中国大陆', 'CNY', 'zh-CN', '🇨🇳', 1, 1, '{"IncludeTax": true, "VATRate": 0.13}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(7, 2, 'US', '美国', 'USD', 'en', '🇺🇸', 1, 0, '{"IncludeTax": false}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- Enterprise Corp 市场
(8, 3, 'CN', '中国大陆', 'CNY', 'zh-CN', '🇨🇳', 1, 1, '{"IncludeTax": true, "VATRate": 0.13}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(9, 3, 'US', '美国', 'USD', 'en', '🇺🇸', 1, 0, '{"IncludeTax": false}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(10, 3, 'UK', '英国', 'GBP', 'en-GB', '🇬🇧', 1, 0, '{"IncludeTax": true, "VATRate": 0.20, "IOSSEnabled": true}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(11, 3, 'DE', '德国', 'EUR', 'de', '🇩🇪', 1, 0, '{"IncludeTax": true, "VATRate": 0.19, "IOSSEnabled": true}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(12, 3, 'FR', '法国', 'EUR', 'fr', '🇫🇷', 1, 0, '{"IncludeTax": true, "VATRate": 0.20, "IOSSEnabled": true}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(13, 3, 'AU', '澳大利亚', 'AUD', 'en-AU', '🇦🇺', 0, 0, '{"IncludeTax": true, "GSTRate": 0.10}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());
