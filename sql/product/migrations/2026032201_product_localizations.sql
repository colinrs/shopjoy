-- ============================================
-- Migration: Product Localizations
-- Date: 2026-03-22
-- Description: Add product_localizations table for
--              multi-language product support
-- ============================================

-- ============================================
-- 1. Create product_localizations table
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