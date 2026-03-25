-- ============================================
-- Migration: Category, Brand, and Inventory Management
-- Date: 2026-03-21
-- Description: Add SEO fields, compliance fields, inventory management,
--              and tenant_id to products/skus tables
-- ============================================

-- ============================================
-- 1. Add tenant_id to products table
-- ============================================
ALTER TABLE `products`
    ADD COLUMN `tenant_id` BIGINT NOT NULL DEFAULT 0 COMMENT '租户ID' AFTER `id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- ============================================
-- 2. Add tenant_id to skus table
-- ============================================
ALTER TABLE `skus`
    ADD COLUMN `tenant_id` BIGINT NOT NULL DEFAULT 0 COMMENT '租户ID' AFTER `id`,
    ADD INDEX `idx_tenant_id` (`tenant_id`);

-- ============================================
-- 3. Add SEO fields to categories table
-- ============================================
ALTER TABLE `categories`
    ADD COLUMN `seo_title` VARCHAR(200) DEFAULT '' COMMENT 'SEO标题' AFTER `image`,
    ADD COLUMN `seo_description` VARCHAR(500) DEFAULT '' COMMENT 'SEO描述' AFTER `seo_title`;

-- ============================================
-- 4. Add compliance fields to brands table
-- ============================================
ALTER TABLE `brands`
    ADD COLUMN `enable_page` TINYINT NOT NULL DEFAULT 0 COMMENT '是否启用品牌专区' AFTER `sort`,
    ADD COLUMN `trademark_number` VARCHAR(100) DEFAULT '' COMMENT '商标号' AFTER `enable_page`,
    ADD COLUMN `trademark_country` VARCHAR(10) DEFAULT '' COMMENT '商标注册国家' AFTER `trademark_number`;

-- ============================================
-- 5. Add brand_id to products table
-- ============================================
ALTER TABLE `products`
    ADD COLUMN `brand_id` BIGINT NULL COMMENT '品牌ID' AFTER `brand`,
    ADD INDEX `idx_brand_id` (`brand_id`);

-- ============================================
-- 6. Add inventory fields to skus table
-- ============================================
ALTER TABLE `skus`
    ADD COLUMN `available_stock` INT NOT NULL DEFAULT 0 COMMENT '可用库存' AFTER `stock`,
    ADD COLUMN `locked_stock` INT NOT NULL DEFAULT 0 COMMENT '锁定库存' AFTER `available_stock`,
    ADD COLUMN `safety_stock` INT NOT NULL DEFAULT 0 COMMENT '安全库存阈值' AFTER `locked_stock`,
    ADD COLUMN `presale_enabled` TINYINT NOT NULL DEFAULT 0 COMMENT '是否开启预售' AFTER `safety_stock`;

-- ============================================
-- 7. Create category_markets table
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
-- 8. Create brand_markets table
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
-- 9. Create warehouses table
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
-- 10. Create warehouse_inventories table
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
-- 11. Create inventory_logs table
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
-- 12. Add deleted_at to categories table (soft delete)
-- ============================================
ALTER TABLE `categories`
    ADD COLUMN `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间' AFTER `updated_by`,
    ADD INDEX `idx_deleted_at` (`deleted_at`);

-- ============================================
-- 13. Add deleted_at to brands table (soft delete)
-- ============================================
ALTER TABLE `brands`
    ADD COLUMN `deleted_at` BIGINT DEFAULT NULL COMMENT '删除时间' AFTER `updated_by`,
    ADD INDEX `idx_deleted_at` (`deleted_at`);