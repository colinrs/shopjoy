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
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_status` (`status`)
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
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_name` (`name`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='品牌表';

-- ============================================
-- 商品表 (products)
-- ============================================

CREATE TABLE IF NOT EXISTS `products` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '商品ID',
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
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_sku` (`sku`),
    KEY `idx_name` (`name`),
    KEY `idx_category_id` (`category_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品表';

-- ============================================
-- SKU表 (skus)
-- ============================================

CREATE TABLE IF NOT EXISTS `skus` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'SKU ID',
    `product_id` BIGINT NOT NULL COMMENT '商品ID',
    `code` VARCHAR(100) NOT NULL COMMENT 'SKU代码',
    `price_amount` BIGINT NOT NULL DEFAULT 0 COMMENT '价格(分)',
    `price_currency` VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',
    `stock` INT NOT NULL DEFAULT 0 COMMENT '库存',
    `attributes` JSON COMMENT '属性',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
    `created_by` BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    `updated_by` BIGINT NOT NULL DEFAULT 0 COMMENT '更新人',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_product_id` (`product_id`),
    KEY `idx_status` (`status`)
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