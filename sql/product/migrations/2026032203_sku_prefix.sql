-- ============================================
-- Migration: SKU Prefix Fields
-- Date: 2026-03-22
-- Description: Add sku_prefix field to tenants and products tables
-- ============================================

-- ============================================
-- 1. Add sku_prefix to tenants table
-- ============================================
ALTER TABLE `tenants`
    ADD COLUMN `sku_prefix` VARCHAR(8) DEFAULT '' COMMENT 'SKU默认前缀' AFTER `address`;

-- ============================================
-- 2. Add sku_prefix to products table
-- ============================================
ALTER TABLE `products`
    ADD COLUMN `sku_prefix` VARCHAR(8) DEFAULT '' COMMENT '商品SKU前缀' AFTER `brand`;

-- ============================================
-- 3. Ensure unique constraint on skus table
-- ============================================
-- Check if index exists before creating
SET @exist := (SELECT COUNT(1) FROM information_schema.statistics
               WHERE table_schema = DATABASE()
               AND table_name = 'skus'
               AND index_name = 'uk_tenant_code');

SET @sqlstmt := IF(@exist = 0,
    'ALTER TABLE `skus` ADD UNIQUE INDEX `uk_tenant_code` (`tenant_id`, `code`)',
    'SELECT ''Index uk_tenant_code already exists''');

PREPARE stmt FROM @sqlstmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;